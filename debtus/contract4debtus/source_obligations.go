// Package contract4debtus defines storage-neutral Debtus provider contracts.
package contract4debtus

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	SourceNamespaceSplitus           = "splitus"
	SourceRepaymentContractVersion   = 1
	MaxSourceObligationIDs           = 256
	SourceRepaymentBrowserTimeLayout = "2006-01-02T15:04:05.000Z"
	SourceDueDateContractVersion     = 1

	// ReconcileSourceObligationsDigestEncoding domain-separates fingerprints
	// made by this contract and encoding from every other Debtus command.
	ReconcileSourceObligationsDigestEncoding = "sneat-ext-contracts/debtus:reconcile-source-obligations:encoding-v1"
)

// FormatSourceRepaymentBrowserTime emits the exact browser wire timestamp.
// time.Time's default JSON uses RFC3339Nano and omits .000 for a whole second,
// so browser adapters must call this instead of directly marshaling time.Time.
func FormatSourceRepaymentBrowserTime(value time.Time) (string, error) {
	if value.IsZero() {
		return "", fmt.Errorf("%w: repayment browser time is required", ErrInvalidRequest)
	}
	_, offset := value.Zone()
	if offset != 0 || value.Nanosecond()%int(time.Millisecond) != 0 {
		return "", fmt.Errorf("%w: repayment browser time must be UTC with millisecond precision", ErrInvalidRequest)
	}
	return value.Format(SourceRepaymentBrowserTimeLayout), nil
}

// ExactMinorAmountString is a canonical non-negative integer amount in a
// currency's minor units. It is encoded as a JSON string so browser clients do
// not lose precision. Use MinorUnits to obtain a checked int64 value.
type ExactMinorAmountString string

func (a *ExactMinorAmountString) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return fmt.Errorf("%w: amountMinor must be a JSON string: %v", ErrInvalidRequest, err)
	}
	parsed := ExactMinorAmountString(value)
	if _, err := parsed.MinorUnits(); err != nil {
		return err
	}
	*a = parsed
	return nil
}

// MinorUnits parses the exact representation and rejects non-canonical or
// signed-64-bit-overflow values.
func (a ExactMinorAmountString) MinorUnits() (int64, error) {
	value := string(a)
	if value == "" || len(value) > 19 || len(value) > 1 && value[0] == '0' {
		return 0, fmt.Errorf("%w: amountMinor must be a canonical non-negative integer string", ErrInvalidRequest)
	}
	for _, digit := range value {
		if digit < '0' || digit > '9' {
			return 0, fmt.Errorf("%w: amountMinor must be a canonical non-negative integer string", ErrInvalidRequest)
		}
	}
	minor, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%w: amountMinor exceeds signed 64-bit range", ErrInvalidRequest)
	}
	return minor, nil
}

var (
	ErrInvalidRequest = errors.New("invalid Debtus source obligation request")
	ErrDigestMismatch = errors.New("source obligation input digest does not match request")
	ErrRevisionStale  = errors.New("source obligation revision is stale")
	ErrOperationClash = errors.New("operation key was already used for different input")
)

// SourceRef is the stable identity of one source-owned financial record.
type SourceRef struct {
	Namespace string `json:"namespace"`
	SpaceID   string `json:"spaceID"`
	RecordID  string `json:"recordID"`
}

// ContactRef identifies a Debtus party through a contact in the source Space.
// It deliberately does not identify the authenticated recorder.
type ContactRef struct {
	SpaceID   string `json:"spaceID"`
	ContactID string `json:"contactID"`
}

// ObligationLine is one desired directed obligation. AmountMinor is an exact
// integer number of the currency's minor units; floating-point amounts are not
// part of this contract.
type ObligationLine struct {
	LineID      string     `json:"lineID"`
	Debtor      ContactRef `json:"debtor"`
	Creditor    ContactRef `json:"creditor"`
	Currency    string     `json:"currency"`
	AmountMinor int64      `json:"amountMinor"`
}

// ReconcileSourceObligationsRequest replaces the desired obligation lines for
// one source revision. RecorderUserID is audit identity, separate from both
// financial parties. A trusted provider must match RecorderUserID to the
// authenticated server context; this caller-supplied value grants no authority.
// InputDigest is the lowercase SHA-256 returned by CanonicalInputDigest.
type ReconcileSourceObligationsRequest struct {
	Source                   SourceRef        `json:"source"`
	RecorderUserID           string           `json:"recorderUserID"`
	ExpectedPreviousRevision uint64           `json:"expectedPreviousRevision"`
	NewRevision              uint64           `json:"newRevision"`
	OperationKey             string           `json:"operationKey"`
	InputDigest              string           `json:"inputDigest"`
	DesiredLines             []ObligationLine `json:"desiredLines"`
}

// Validate checks contract-level shape and deterministic replay identity.
// Provider authorization and current membership are runtime responsibilities.
func (r ReconcileSourceObligationsRequest) Validate() error {
	if err := validateToken("source namespace", r.Source.Namespace); err != nil {
		return err
	}
	if err := validateID("source spaceID", r.Source.SpaceID); err != nil {
		return err
	}
	if err := validateID("source recordID", r.Source.RecordID); err != nil {
		return err
	}
	if err := validateID("recorder userID", r.RecorderUserID); err != nil {
		return err
	}
	if err := validateID("operation key", r.OperationKey); err != nil {
		return err
	}
	if r.NewRevision == 0 || r.ExpectedPreviousRevision == math.MaxUint64 || r.NewRevision != r.ExpectedPreviousRevision+1 {
		return fmt.Errorf("%w: new revision must immediately follow expected previous revision", ErrInvalidRequest)
	}
	seen := make(map[string]struct{}, len(r.DesiredLines))
	totals := make(map[string]int64)
	for i, line := range r.DesiredLines {
		if err := validateLine(r.Source.SpaceID, line); err != nil {
			return fmt.Errorf("desired line %d: %w", i, err)
		}
		if _, ok := seen[line.LineID]; ok {
			return fmt.Errorf("%w: duplicate lineID %q", ErrInvalidRequest, line.LineID)
		}
		seen[line.LineID] = struct{}{}
		if totals[line.Currency] > math.MaxInt64-line.AmountMinor {
			return fmt.Errorf("%w: %s amount total overflows int64 minor units", ErrInvalidRequest, line.Currency)
		}
		totals[line.Currency] += line.AmountMinor
	}
	digest, err := r.CanonicalInputDigest()
	if err != nil {
		return err
	}
	if r.InputDigest != digest {
		return fmt.Errorf("%w: got %q, want %q", ErrDigestMismatch, r.InputDigest, digest)
	}
	return nil
}

func validateLine(sourceSpaceID string, line ObligationLine) error {
	if err := validateID("lineID", line.LineID); err != nil {
		return err
	}
	for name, party := range map[string]ContactRef{"debtor": line.Debtor, "creditor": line.Creditor} {
		if party.SpaceID != sourceSpaceID {
			return fmt.Errorf("%w: %s spaceID %q differs from source spaceID", ErrInvalidRequest, name, party.SpaceID)
		}
		if err := validateID(name+" contactID", party.ContactID); err != nil {
			return err
		}
	}
	if line.Debtor.ContactID == line.Creditor.ContactID {
		return fmt.Errorf("%w: debtor and creditor contactID must differ", ErrInvalidRequest)
	}
	if !isCurrencyCode(line.Currency) {
		return fmt.Errorf("%w: currency %q must be three uppercase ASCII letters", ErrInvalidRequest, line.Currency)
	}
	if line.AmountMinor <= 0 {
		return fmt.Errorf("%w: amountMinor must be positive", ErrInvalidRequest)
	}
	return nil
}

func validateID(name, value string) error {
	if !utf8.ValidString(value) || value == "" || strings.TrimSpace(value) != value || len(value) > 512 {
		return fmt.Errorf("%w: %s is empty, padded, or longer than 512 bytes", ErrInvalidRequest, name)
	}
	for _, r := range value {
		if r < 0x20 || r == 0x7f {
			return fmt.Errorf("%w: %s contains a control character", ErrInvalidRequest, name)
		}
	}
	return nil
}

func validateToken(name, value string) error {
	if err := validateID(name, value); err != nil {
		return err
	}
	for _, r := range value {
		if !(r >= 'a' && r <= 'z' || r >= '0' && r <= '9' || r == '-' || r == '_') {
			return fmt.Errorf("%w: %s must contain lowercase ASCII letters, digits, hyphens, or underscores", ErrInvalidRequest, name)
		}
	}
	return nil
}

func isCurrencyCode(value string) bool {
	if len(value) != 3 {
		return false
	}
	for i := range value {
		if value[i] < 'A' || value[i] > 'Z' {
			return false
		}
	}
	return true
}

// CanonicalInputDigest returns a deterministic lowercase SHA-256 fingerprint.
// Lines are ordered by LineID so transport ordering does not change identity.
func (r ReconcileSourceObligationsRequest) CanonicalInputDigest() (string, error) {
	lines := append([]ObligationLine(nil), r.DesiredLines...)
	sort.Slice(lines, func(i, j int) bool { return lines[i].LineID < lines[j].LineID })
	h := sha256.New()
	writeString := func(value string) {
		var size [8]byte
		binary.BigEndian.PutUint64(size[:], uint64(len(value)))
		_, _ = h.Write(size[:])
		_, _ = h.Write([]byte(value))
	}
	writeUint64 := func(value uint64) {
		var data [8]byte
		binary.BigEndian.PutUint64(data[:], value)
		_, _ = h.Write(data[:])
	}
	writeString(ReconcileSourceObligationsDigestEncoding)
	writeString(r.Source.Namespace)
	writeString(r.Source.SpaceID)
	writeString(r.Source.RecordID)
	writeString(r.RecorderUserID)
	writeUint64(r.ExpectedPreviousRevision)
	writeUint64(r.NewRevision)
	writeString(r.OperationKey)
	writeUint64(uint64(len(lines)))
	for _, line := range lines {
		writeString(line.LineID)
		writeString(line.Debtor.SpaceID)
		writeString(line.Debtor.ContactID)
		writeString(line.Creditor.SpaceID)
		writeString(line.Creditor.ContactID)
		writeString(line.Currency)
		writeUint64(uint64(line.AmountMinor))
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// ObligationResult maps a stable source line to Debtus-owned obligation IDs.
type ObligationResult struct {
	LineID        string   `json:"lineID"`
	ObligationIDs []string `json:"obligationIDs"`
}

// SourceObligationsReceipt is the immutable durable provider receipt, safe to
// return again for the same operation key and digest. Mutable posting and
// financial state are returned by GetSourceObligationsStatus.
type SourceObligationsReceipt struct {
	ReceiptID    string             `json:"receiptID"`
	Source       SourceRef          `json:"source"`
	Revision     uint64             `json:"revision"`
	OperationKey string             `json:"operationKey"`
	InputDigest  string             `json:"inputDigest"`
	Obligations  []ObligationResult `json:"obligations,omitempty"`
}

// ReconcileSourceObligationsResult reports current posting progress separately
// from the immutable applied receipt. Pending and posting are explicitly not
// proof that Debtus has applied any financial entry.
type ReconcileSourceObligationsResult struct {
	PostingStatus   PostingStatus             `json:"postingStatus"`
	AppliedReceipt  *SourceObligationsReceipt `json:"appliedReceipt,omitempty"`
	AttentionReason string                    `json:"attentionReason,omitempty"`
}

// Validate enforces the result envelope's proof boundary.
func (r ReconcileSourceObligationsResult) Validate() error {
	switch r.PostingStatus {
	case PostingStatusPending, PostingStatusPosting:
		if r.AppliedReceipt != nil {
			return fmt.Errorf("%w: %s result cannot contain an applied receipt", ErrInvalidRequest, r.PostingStatus)
		}
		if r.AttentionReason != "" {
			return fmt.Errorf("%w: %s result cannot contain an attention reason", ErrInvalidRequest, r.PostingStatus)
		}
	case PostingStatusApplied:
		if r.AppliedReceipt == nil {
			return fmt.Errorf("%w: applied result requires an immutable receipt", ErrInvalidRequest)
		}
		if r.AttentionReason != "" {
			return fmt.Errorf("%w: applied result cannot contain an attention reason", ErrInvalidRequest)
		}
	case PostingStatusAttention:
		if r.AppliedReceipt != nil {
			return fmt.Errorf("%w: attention result cannot contain an applied receipt", ErrInvalidRequest)
		}
		if strings.TrimSpace(r.AttentionReason) == "" {
			return fmt.Errorf("%w: attention result requires a reason", ErrInvalidRequest)
		}
	default:
		return fmt.Errorf("%w: unknown posting status %q", ErrInvalidRequest, r.PostingStatus)
	}
	return nil
}

// GetSourceObligationsStatusRequest identifies an authorized status read. A
// trusted provider must match ActorUserID to the authenticated server context
// and authorize that actor; this caller-supplied value grants no authority.
type GetSourceObligationsStatusRequest struct {
	Source      SourceRef `json:"source"`
	ActorUserID string    `json:"actorUserID"`
}

// PostingStatus describes whether Debtus has applied the accepted source
// revision. It is distinct from each obligation's financial settlement state.
type PostingStatus string

const (
	PostingStatusPending   PostingStatus = "pending"
	PostingStatusPosting   PostingStatus = "posting"
	PostingStatusApplied   PostingStatus = "applied"
	PostingStatusAttention PostingStatus = "attention"
)

// SettlementStatus is Debtus's current authoritative financial state for one
// source line.
type SettlementStatus string

const (
	SettlementStatusUnsettled   SettlementStatus = "unsettled"
	SettlementStatusPartSettled SettlementStatus = "part_settled"
	SettlementStatusSettled     SettlementStatus = "settled"
)

// SourceRepaymentUnavailableReason is a stable server-authored explanation for
// why a browser must not offer repayment for one source obligation. Clients
// render this value; they never infer permission from Contactus data.
type SourceRepaymentUnavailableReason string

const (
	SourceRepaymentReadOnlyRole            SourceRepaymentUnavailableReason = "readOnlyRole"
	SourceRepaymentNotParty                SourceRepaymentUnavailableReason = "notParty"
	SourceRepaymentNotOutstanding          SourceRepaymentUnavailableReason = "notOutstanding"
	SourceRepaymentLedgerPartyUnresolved   SourceRepaymentUnavailableReason = "ledgerPartyUnresolved"
	SourceRepaymentLinkedEventLimitReached SourceRepaymentUnavailableReason = "linkedEventLimitReached"
)

// SourceObligationRepaymentCapability is a render-time projection, not an
// authorization grant. A provider must reauthorize every repayment command.
type SourceObligationRepaymentCapability struct {
	CanRecordRepayment         bool                             `json:"canRecordRepayment"`
	RepaymentUnavailableReason SourceRepaymentUnavailableReason `json:"repaymentUnavailableReason,omitempty"`
}

func (c SourceObligationRepaymentCapability) Validate() error {
	if c.CanRecordRepayment {
		if c.RepaymentUnavailableReason != "" {
			return fmt.Errorf("%w: enabled repayment capability cannot have an unavailable reason", ErrInvalidRequest)
		}
		return nil
	}
	switch c.RepaymentUnavailableReason {
	case SourceRepaymentReadOnlyRole, SourceRepaymentNotParty, SourceRepaymentNotOutstanding,
		SourceRepaymentLedgerPartyUnresolved, SourceRepaymentLinkedEventLimitReached:
		return nil
	default:
		return fmt.Errorf("%w: disabled repayment capability requires a supported unavailable reason", ErrInvalidRequest)
	}
}

// SourceObligationStatus reports exact current Debtus amounts for one source
// line. PrincipalMinor is the accepted obligation principal, OutstandingMinor
// is unpaid liability, RepaidMinor is repayment retained in Debtus history,
// and CreditMinor is any current credit Debtus recognizes for the line. The
// presence of credit does not prescribe a refund or settlement path.
type SourceObligationStatus struct {
	LineID           string           `json:"lineID"`
	ObligationIDs    []string         `json:"obligationIDs"`
	Debtor           ContactRef       `json:"debtor"`
	Creditor         ContactRef       `json:"creditor"`
	Currency         string           `json:"currency"`
	PrincipalMinor   int64            `json:"principalMinor"`
	OutstandingMinor int64            `json:"outstandingMinor"`
	RepaidMinor      int64            `json:"repaidMinor"`
	CreditMinor      int64            `json:"creditMinor"`
	Status           SettlementStatus `json:"status"`
	// RepaymentCapability is optional for compatibility with readers created
	// before the source-scoped repayment contract. New mutation/read adapters
	// must populate it whenever they expose the repayment action.
	RepaymentCapability *SourceObligationRepaymentCapability `json:"repaymentCapability,omitempty"`
	DueDate             *SourceObligationDueDateResult       `json:"dueDateTask,omitempty"`
}

// Validate checks output invariants that consumers may rely on. The contract
// intentionally does not require principal minus repaid to equal outstanding:
// adjustments and credits make that equation incomplete.
func (s SourceObligationStatus) Validate() error {
	if err := validateLine(s.Debtor.SpaceID, ObligationLine{
		LineID: s.LineID, Debtor: s.Debtor, Creditor: s.Creditor,
		Currency: s.Currency, AmountMinor: 1,
	}); err != nil {
		return err
	}
	if len(s.ObligationIDs) == 0 || len(s.ObligationIDs) > MaxSourceObligationIDs {
		return fmt.Errorf("%w: financial status requires 1 to %d obligation IDs", ErrInvalidRequest, MaxSourceObligationIDs)
	}
	seenObligationIDs := make(map[string]struct{}, len(s.ObligationIDs))
	for _, obligationID := range s.ObligationIDs {
		if err := validateID("obligation ID", obligationID); err != nil {
			return err
		}
		if _, exists := seenObligationIDs[obligationID]; exists {
			return fmt.Errorf("%w: duplicate obligation ID %q", ErrInvalidRequest, obligationID)
		}
		seenObligationIDs[obligationID] = struct{}{}
	}
	if s.PrincipalMinor < 0 || s.OutstandingMinor < 0 || s.RepaidMinor < 0 || s.CreditMinor < 0 {
		return fmt.Errorf("%w: financial status amounts must be nonnegative minor units", ErrInvalidRequest)
	}
	switch s.Status {
	case SettlementStatusUnsettled, SettlementStatusPartSettled, SettlementStatusSettled:
	default:
		return fmt.Errorf("%w: unknown settlement status %q", ErrInvalidRequest, s.Status)
	}
	if s.RepaymentCapability != nil {
		if err := s.RepaymentCapability.Validate(); err != nil {
			return err
		}
	}
	if s.DueDate != nil {
		if err := s.DueDate.Validate(); err != nil {
			return err
		}
		if s.DueDate.LineID != s.LineID {
			return fmt.Errorf("%w: due-date task lineID differs from obligation", ErrInvalidRequest)
		}
	}
	return nil
}

// SourceObligationsStatus is mutable authoritative state read from Debtus. It
// is separate from the immutable reconciliation receipt because repayments,
// adjustments, and settlement can change after the source revision is posted.
type SourceObligationsStatus struct {
	LatestReceipt   SourceObligationsReceipt `json:"latestReceipt"`
	PostingStatus   PostingStatus            `json:"postingStatus"`
	AttentionReason string                   `json:"attentionReason,omitempty"`
	Obligations     []SourceObligationStatus `json:"obligations,omitempty"`
}

// SourceObligationActivityKind identifies a Debtus-authoritative historical
// event without exposing its persistence representation. Correction and credit
// kinds describe recorded facts; they do not prescribe approval or refund policy.
type SourceObligationActivityKind string

const (
	SourceActivityObligation   SourceObligationActivityKind = "obligation"
	SourceActivityRepayment    SourceObligationActivityKind = "repayment"
	SourceActivityAdjustment   SourceObligationActivityKind = "adjustment"
	SourceActivityCancellation SourceObligationActivityKind = "cancellation"
	SourceActivityCredit       SourceObligationActivityKind = "credit"
)

// SourceObligationActivity is one immutable, auditable Debtus activity. The
// stable RootActivityID links repayments, returns, and adjustments to their
// original financial root; LineIDs link the explanation back to source lines.
// AmountMinor is the exact nonnegative event magnitude in currency minor units.
type SourceObligationActivity struct {
	ActivityID         string                       `json:"activityID"`
	RootActivityID     string                       `json:"rootActivityID"`
	RelatedActivityIDs []string                     `json:"relatedActivityIDs,omitempty"`
	LineIDs            []string                     `json:"lineIDs"`
	Kind               SourceObligationActivityKind `json:"kind"`
	From               ContactRef                   `json:"from"`
	To                 ContactRef                   `json:"to"`
	Currency           string                       `json:"currency"`
	AmountMinor        int64                        `json:"amountMinor"`
	ActorUserID        string                       `json:"actorUserID"`
	OccurredAt         time.Time                    `json:"occurredAt"`
}

// Validate checks the stable, storage-neutral activity shape.
func (a SourceObligationActivity) Validate() error {
	for name, value := range map[string]string{
		"activityID": a.ActivityID, "rootActivityID": a.RootActivityID, "actorUserID": a.ActorUserID,
	} {
		if err := validateID(name, value); err != nil {
			return err
		}
	}
	if len(a.LineIDs) == 0 {
		return fmt.Errorf("%w: activity requires at least one source lineID", ErrInvalidRequest)
	}
	seen := make(map[string]struct{}, len(a.LineIDs))
	for _, lineID := range a.LineIDs {
		if err := validateID("activity lineID", lineID); err != nil {
			return err
		}
		if _, ok := seen[lineID]; ok {
			return fmt.Errorf("%w: duplicate activity lineID %q", ErrInvalidRequest, lineID)
		}
		seen[lineID] = struct{}{}
	}
	for name, party := range map[string]ContactRef{"from": a.From, "to": a.To} {
		if err := validateID(name+" spaceID", party.SpaceID); err != nil {
			return err
		}
		if err := validateID(name+" contactID", party.ContactID); err != nil {
			return err
		}
		if party.SpaceID != a.From.SpaceID {
			return fmt.Errorf("%w: activity parties must belong to the same source Space", ErrInvalidRequest)
		}
	}
	if a.From.ContactID == a.To.ContactID {
		return fmt.Errorf("%w: activity from and to contacts must differ", ErrInvalidRequest)
	}
	switch a.Kind {
	case SourceActivityObligation, SourceActivityRepayment, SourceActivityAdjustment, SourceActivityCancellation, SourceActivityCredit:
	default:
		return fmt.Errorf("%w: unknown source activity kind %q", ErrInvalidRequest, a.Kind)
	}
	if !isCurrencyCode(a.Currency) {
		return fmt.Errorf("%w: currency %q must be three uppercase ASCII letters", ErrInvalidRequest, a.Currency)
	}
	if a.AmountMinor <= 0 {
		return fmt.Errorf("%w: activity amountMinor must be positive", ErrInvalidRequest)
	}
	if a.OccurredAt.IsZero() {
		return fmt.Errorf("%w: activity occurredAt is required", ErrInvalidRequest)
	}
	return nil
}

// ListSourceObligationActivitiesRequest requests one bounded history page. A
// trusted provider must match ActorUserID to authenticated server context and
// authorize the source read. PageSize must be positive; providers may impose a
// documented lower maximum. Cursor is opaque to callers.
type ListSourceObligationActivitiesRequest struct {
	Source      SourceRef `json:"source"`
	ActorUserID string    `json:"actorUserID"`
	PageSize    uint16    `json:"pageSize"`
	Cursor      string    `json:"cursor,omitempty"`
}

// Validate checks that the caller requested a bounded page for one source.
func (r ListSourceObligationActivitiesRequest) Validate() error {
	if err := validateToken("source namespace", r.Source.Namespace); err != nil {
		return err
	}
	if err := validateID("source spaceID", r.Source.SpaceID); err != nil {
		return err
	}
	if err := validateID("source recordID", r.Source.RecordID); err != nil {
		return err
	}
	if err := validateID("actor userID", r.ActorUserID); err != nil {
		return err
	}
	if r.PageSize == 0 {
		return fmt.Errorf("%w: activity pageSize must be positive", ErrInvalidRequest)
	}
	if len(r.Cursor) > 4096 {
		return fmt.Errorf("%w: activity cursor is longer than 4096 bytes", ErrInvalidRequest)
	}
	return nil
}

// SourceObligationActivitiesPage is one bounded page of authoritative history.
type SourceObligationActivitiesPage struct {
	Activities []SourceObligationActivity `json:"activities"`
	NextCursor string                     `json:"nextCursor,omitempty"`
}

// RecordSourceObligationRepaymentRequest identifies one exact source-owned
// obligation and one retry-safe repayment. ActorUserID is supplied from the
// trusted authenticated host context and is deliberately absent from the
// browser request DTO. It grants no authority: the provider must require the
// actor to be the current debtor or creditor linked to this obligation and to
// have current Space authority before any write or idempotency claim.
type RecordSourceObligationRepaymentRequest struct {
	ContractVersion int                    `json:"contractVersion"`
	Source          SourceRef              `json:"source"`
	LineID          string                 `json:"lineID"`
	ObligationID    string                 `json:"obligationID"`
	Currency        string                 `json:"currency"`
	AmountMinor     ExactMinorAmountString `json:"amountMinor"`
	RepaidAt        time.Time              `json:"repaidAt"`
	OperationKey    string                 `json:"operationKey"`
	ActorUserID     string                 `json:"-"`
}

func (r RecordSourceObligationRepaymentRequest) Validate() error {
	if r.ContractVersion != SourceRepaymentContractVersion {
		return fmt.Errorf("%w: unsupported source repayment contract version %d", ErrInvalidRequest, r.ContractVersion)
	}
	if err := validateToken("source namespace", r.Source.Namespace); err != nil {
		return err
	}
	for name, value := range map[string]string{
		"source spaceID": r.Source.SpaceID, "source recordID": r.Source.RecordID,
		"lineID": r.LineID, "obligationID": r.ObligationID,
		"operation key": r.OperationKey, "actor userID": r.ActorUserID,
	} {
		if err := validateID(name, value); err != nil {
			return err
		}
	}
	if !isCurrencyCode(r.Currency) {
		return fmt.Errorf("%w: currency %q must be three uppercase ASCII letters", ErrInvalidRequest, r.Currency)
	}
	amount, err := r.AmountMinor.MinorUnits()
	if err != nil {
		return err
	}
	if amount <= 0 {
		return fmt.Errorf("%w: amountMinor must be positive", ErrInvalidRequest)
	}
	if r.RepaidAt.IsZero() {
		return fmt.Errorf("%w: repaidAt is required", ErrInvalidRequest)
	}
	_, offset := r.RepaidAt.Zone()
	if offset != 0 {
		return fmt.Errorf("%w: repaidAt must be UTC", ErrInvalidRequest)
	}
	if r.RepaidAt.Nanosecond()%int(time.Millisecond) != 0 {
		return fmt.Errorf("%w: repaidAt must have millisecond precision", ErrInvalidRequest)
	}
	return nil
}

// DecodeRecordSourceObligationRepaymentRequest decodes the actor-free browser
// DTO and adds the trusted authenticated actor supplied by the host. An
// actorUserID field in JSON is rejected as unknown rather than accepted as
// authority.
func DecodeRecordSourceObligationRepaymentRequest(reader io.Reader, authenticatedActorUserID string) (RecordSourceObligationRepaymentRequest, error) {
	decoder := json.NewDecoder(reader)
	decoder.DisallowUnknownFields()
	request := RecordSourceObligationRepaymentRequest{ActorUserID: authenticatedActorUserID}
	if err := decoder.Decode(&request); err != nil {
		return RecordSourceObligationRepaymentRequest{}, fmt.Errorf("%w: decode: %v", ErrInvalidRequest, err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		if err == nil {
			err = errors.New("multiple JSON values")
		}
		return RecordSourceObligationRepaymentRequest{}, fmt.Errorf("%w: trailing JSON: %v", ErrInvalidRequest, err)
	}
	if err := request.Validate(); err != nil {
		return RecordSourceObligationRepaymentRequest{}, err
	}
	return request, nil
}

// RecordSourceObligationRepaymentResult returns the updated authoritative
// obligation plus the immutable repayment activity. It reuses the existing
// source read models instead of defining a second financial projection. This
// is a trusted Go-port result: the browser adapter formats every minor amount
// as an exact string and omits SourceObligationActivity.ActorUserID.
type RecordSourceObligationRepaymentResult struct {
	ContractVersion int                      `json:"contractVersion"`
	Source          SourceRef                `json:"source"`
	OperationKey    string                   `json:"operationKey"`
	Obligation      SourceObligationStatus   `json:"obligation"`
	Repayment       SourceObligationActivity `json:"repayment"`
}

// ValidateFor binds a result to the exact accepted command. This catches a
// provider returning a valid-looking activity for another source line,
// obligation, actor, currency, amount, timestamp, or retry identity.
func (r RecordSourceObligationRepaymentResult) ValidateFor(request RecordSourceObligationRepaymentRequest) error {
	if err := request.Validate(); err != nil {
		return err
	}
	if r.ContractVersion != SourceRepaymentContractVersion {
		return fmt.Errorf("%w: unsupported source repayment result contract version %d", ErrInvalidRequest, r.ContractVersion)
	}
	if r.Source != request.Source {
		return fmt.Errorf("%w: repayment result source does not match request", ErrInvalidRequest)
	}
	if r.OperationKey != request.OperationKey {
		return fmt.Errorf("%w: repayment result operation key does not match request", ErrInvalidRequest)
	}
	if err := r.Obligation.Validate(); err != nil {
		return err
	}
	if r.Obligation.RepaymentCapability == nil {
		return fmt.Errorf("%w: repayment result requires server-authoritative capability", ErrInvalidRequest)
	}
	if r.Obligation.LineID != request.LineID || r.Obligation.Currency != request.Currency {
		return fmt.Errorf("%w: repayment result obligation does not match requested line and currency", ErrInvalidRequest)
	}
	foundObligation := false
	for _, obligationID := range r.Obligation.ObligationIDs {
		if obligationID == request.ObligationID {
			foundObligation = true
			break
		}
	}
	if !foundObligation {
		return fmt.Errorf("%w: repayment result does not contain requested obligationID", ErrInvalidRequest)
	}
	if err := r.Repayment.Validate(); err != nil {
		return err
	}
	if r.Repayment.From != r.Obligation.Debtor || r.Repayment.To != r.Obligation.Creditor {
		return fmt.Errorf("%w: repayment parties do not match obligation direction", ErrInvalidRequest)
	}
	amount, _ := request.AmountMinor.MinorUnits()
	if r.Repayment.Kind != SourceActivityRepayment || r.Repayment.RootActivityID != request.ObligationID ||
		len(r.Repayment.LineIDs) != 1 || r.Repayment.LineIDs[0] != request.LineID ||
		r.Repayment.Currency != request.Currency || r.Repayment.AmountMinor != amount ||
		r.Repayment.ActorUserID != request.ActorUserID || !r.Repayment.OccurredAt.Equal(request.RepaidAt) {
		return fmt.Errorf("%w: repayment activity does not match accepted command", ErrInvalidRequest)
	}
	return nil
}

// SourceObligationRepayments is the storage-neutral repayment port. Runtime
// composition supplies the trusted actor and reauthorizes every call.
type SourceObligationRepayments interface {
	RecordSourceObligationRepayment(context.Context, RecordSourceObligationRepaymentRequest) (RecordSourceObligationRepaymentResult, error)
}

type SetSourceObligationDueDateRequest struct {
	ContractVersion  int       `json:"contractVersion"`
	Source           SourceRef `json:"source"`
	LineID           string    `json:"lineID"`
	Title            string    `json:"title"`
	DueDate          string    `json:"dueDate,omitempty"`
	ExpectedRevision uint64    `json:"expectedRevision"`
	OperationKey     string    `json:"operationKey"`
	TodoListID       string    `json:"todoListID,omitempty"`
	ActorUserID      string    `json:"-"`
}

func (r SetSourceObligationDueDateRequest) Validate() error {
	if r.ContractVersion != SourceDueDateContractVersion {
		return fmt.Errorf("%w: unsupported source due-date contract version", ErrInvalidRequest)
	}
	if err := validateToken("source namespace", r.Source.Namespace); err != nil {
		return err
	}
	if err := validateID("source spaceID", r.Source.SpaceID); err != nil {
		return err
	}
	if err := validateID("source recordID", r.Source.RecordID); err != nil {
		return err
	}
	if err := validateID("lineID", r.LineID); err != nil {
		return err
	}
	if err := validateID("operation key", r.OperationKey); err != nil {
		return err
	}
	if r.Title == "" || r.Title != strings.TrimSpace(r.Title) || len(r.Title) > 100 {
		return fmt.Errorf("%w: title must be trimmed and 1..100 bytes", ErrInvalidRequest)
	}
	if r.DueDate != "" {
		parsed, err := time.Parse("2006-01-02", r.DueDate)
		if err != nil || parsed.Format("2006-01-02") != r.DueDate {
			return fmt.Errorf("%w: dueDate must be YYYY-MM-DD", ErrInvalidRequest)
		}
	}
	if r.TodoListID != "" {
		if err := validateID("todoListID", r.TodoListID); err != nil {
			return err
		}
	}
	if r.ExpectedRevision > 9_007_199_254_740_991 {
		return fmt.Errorf("%w: expectedRevision exceeds browser safe integer", ErrInvalidRequest)
	}
	if r.ActorUserID == "" {
		return fmt.Errorf("%w: authenticated actor is required", ErrInvalidRequest)
	}
	return nil
}

type SourceObligationDueDateState string

const (
	SourceObligationDueDateActive    SourceObligationDueDateState = "active"
	SourceObligationDueDateCompleted SourceObligationDueDateState = "completed"
	SourceObligationDueDateCanceled  SourceObligationDueDateState = "canceled"
)

type SourceObligationDueDateResult struct {
	ContractVersion int                          `json:"contractVersion"`
	Source          SourceRef                    `json:"source"`
	LineID          string                       `json:"lineID"`
	Revision        uint64                       `json:"revision"`
	Title           string                       `json:"title,omitempty"`
	DueDate         string                       `json:"dueDate,omitempty"`
	HappeningID     string                       `json:"happeningID"`
	State           SourceObligationDueDateState `json:"state"`
	TodoListID      string                       `json:"todoListID,omitempty"`
	TodoItemID      string                       `json:"todoItemID,omitempty"`
	UpdatedAt       time.Time                    `json:"updatedAt"`
	UpdatedBy       string                       `json:"updatedBy"`
}

func (r SourceObligationDueDateResult) Validate() error {
	if r.ContractVersion != SourceDueDateContractVersion || r.Revision == 0 || r.Revision > 9_007_199_254_740_991 {
		return fmt.Errorf("%w: invalid due-date result version or revision", ErrInvalidRequest)
	}
	if err := validateToken("source namespace", r.Source.Namespace); err != nil {
		return err
	}
	if err := validateID("source spaceID", r.Source.SpaceID); err != nil {
		return err
	}
	if err := validateID("source recordID", r.Source.RecordID); err != nil {
		return err
	}
	if err := validateID("lineID", r.LineID); err != nil {
		return err
	}
	if err := validateID("happeningID", r.HappeningID); err != nil {
		return err
	}
	if err := validateID("updatedBy", r.UpdatedBy); err != nil {
		return err
	}
	if r.Title != "" && (r.Title != strings.TrimSpace(r.Title) || len(r.Title) > 100) {
		return fmt.Errorf("%w: title must be trimmed and at most 100 bytes", ErrInvalidRequest)
	}
	if r.UpdatedAt.IsZero() {
		return fmt.Errorf("%w: updatedAt is required", ErrInvalidRequest)
	}
	if r.DueDate != "" {
		parsed, err := time.Parse("2006-01-02", r.DueDate)
		if err != nil || parsed.Format("2006-01-02") != r.DueDate {
			return fmt.Errorf("%w: dueDate must be YYYY-MM-DD", ErrInvalidRequest)
		}
	}
	switch r.State {
	case SourceObligationDueDateActive:
		if r.DueDate == "" {
			return fmt.Errorf("%w: active due-date task requires dueDate", ErrInvalidRequest)
		}
	case SourceObligationDueDateCompleted, SourceObligationDueDateCanceled:
	default:
		return fmt.Errorf("%w: unsupported due-date task state", ErrInvalidRequest)
	}
	return nil
}

type SourceObligationDueDates interface {
	SetSourceObligationDueDate(context.Context, SetSourceObligationDueDateRequest) (SourceObligationDueDateResult, error)
}

// SourceObligations provides the public Debtus reconciliation/read boundary.
// Runtime composition binds trusted authority evidence outside these DTOs.
type SourceObligations interface {
	ReconcileSourceObligations(context.Context, ReconcileSourceObligationsRequest) (ReconcileSourceObligationsResult, error)
	GetSourceObligationsStatus(context.Context, GetSourceObligationsStatusRequest) (SourceObligationsStatus, error)
	ListSourceObligationActivities(context.Context, ListSourceObligationActivitiesRequest) (SourceObligationActivitiesPage, error)
}
