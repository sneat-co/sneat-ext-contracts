// Package contract4splitus defines storage-neutral Splitus provider contracts.
package contract4splitus

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	BillContractVersion     = 1
	MaxBillParticipants     = 256
	MaxBillListPageSize     = 100
	MaxObligationIDsPerLine = 256
	SettlementRoute         = "debtus.source-obligations"
	DebtusSourceNamespace   = "splitus"
)

var ErrInvalidRequest = errors.New("invalid Splitus bill request")

// ExactDecimalString is a canonical, non-negative major-unit amount with
// exactly two fraction digits. It is a string on the wire by design.
type ExactDecimalString string

func (a *ExactDecimalString) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return fmt.Errorf("%w: amount must be a JSON string: %v", ErrInvalidRequest, err)
	}
	parsed := ExactDecimalString(value)
	if _, err := parsed.MinorUnits(); err != nil {
		return err
	}
	*a = parsed
	return nil
}

// MinorUnits converts the canonical two-fraction-digit representation without
// floating point and rejects values above signed 64-bit minor units.
func (a ExactDecimalString) MinorUnits() (int64, error) {
	value := string(a)
	dot := strings.IndexByte(value, '.')
	if dot < 1 || len(value)-dot != 3 || value[0] == '-' {
		return 0, fmt.Errorf("%w: amount must have exactly two fraction digits", ErrInvalidRequest)
	}
	whole, fraction := value[:dot], value[dot+1:]
	if whole == "" || len(whole) > 1 && whole[0] == '0' {
		return 0, fmt.Errorf("%w: amount is not canonical", ErrInvalidRequest)
	}
	for _, digit := range whole + fraction {
		if digit < '0' || digit > '9' {
			return 0, fmt.Errorf("%w: amount contains a non-decimal digit", ErrInvalidRequest)
		}
	}
	minor, err := strconv.ParseInt(whole+fraction, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%w: amount exceeds signed 64-bit minor units", ErrInvalidRequest)
	}
	return minor, nil
}

type BillKind string

const (
	BillKindGeneral BillKind = "general"
	BillKindUtility BillKind = "utility"
)

type CurrencyCode string

const (
	CurrencyEUR CurrencyCode = "EUR"
	CurrencyGBP CurrencyCode = "GBP"
	CurrencyUSD CurrencyCode = "USD"
)

type UtilityKind string

const (
	UtilityKindElectricity UtilityKind = "electricity"
	UtilityKindGas         UtilityKind = "gas"
	UtilityKindWater       UtilityKind = "water"
	UtilityKindInternet    UtilityKind = "internet"
	UtilityKindOther       UtilityKind = "other"
)

type PostingStatus string

const (
	PostingStatusPending   PostingStatus = "pending"
	PostingStatusPosting   PostingStatus = "posting"
	PostingStatusApplied   PostingStatus = "applied"
	PostingStatusAttention PostingStatus = "attention"
)

type DebtusSettlementStatus string

const (
	DebtusSettlementStatusUnsettled   DebtusSettlementStatus = "unsettled"
	DebtusSettlementStatusPartSettled DebtusSettlementStatus = "part_settled"
	DebtusSettlementStatusSettled     DebtusSettlementStatus = "settled"
)

type ExpectedActualComparison string

const (
	ExpectedActualNotAvailable ExpectedActualComparison = "not_available"
	ExpectedActualMatches      ExpectedActualComparison = "matches"
	ExpectedActualIncreased    ExpectedActualComparison = "increased"
	ExpectedActualDecreased    ExpectedActualComparison = "decreased"
)

type BillAttentionCode string

const (
	BillAttentionAuthorizationChanged BillAttentionCode = "authorization_changed"
	BillAttentionSourceConflict       BillAttentionCode = "source_conflict"
	BillAttentionProviderRejected     BillAttentionCode = "provider_rejected"
	BillAttentionInvalidReceipt       BillAttentionCode = "invalid_provider_receipt"
	BillAttentionOperatorAction       BillAttentionCode = "operator_action_required"
)

type BillAllocationV1 struct {
	AllocationID string `json:"allocationID"`
	// ContactID is resolved through Contactus by the host. The opaque value is
	// never a user-facing label.
	ContactID string             `json:"contactID"`
	Amount    ExactDecimalString `json:"amount"`
}

type BillingPeriodV1 struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

type UtilityDetailsV1 struct {
	UtilityKind  UtilityKind     `json:"utilityKind"`
	ProviderName *string         `json:"providerName,omitempty"`
	Period       BillingPeriodV1 `json:"period"`
}

type RecurringOccurrenceV1 struct {
	HappeningID          string                    `json:"happeningID"`
	OccurrenceID         string                    `json:"occurrenceID"`
	ExpectedAmount       *ExactDecimalString       `json:"expectedAmount,omitempty"`
	StandingChargeAmount *ExactDecimalString       `json:"standingChargeAmount,omitempty"`
	ExpectedComparison   ExpectedActualComparison  `json:"expectedComparison"`
	PreviousComparable   *PreviousComparableBillV1 `json:"previousComparable,omitempty"`
}

type PreviousComparableBillV1 struct {
	BillID       string                   `json:"billID"`
	ActualAmount ExactDecimalString       `json:"actualAmount"`
	Comparison   ExpectedActualComparison `json:"comparison"`
}

type CreateBillV1Request struct {
	ContractVersion int    `json:"contractVersion"`
	SpaceID         string `json:"spaceID"`
	// BillID is client-stable across duplicate submission and lost responses.
	// Reusing it with changed allocations is a provider conflict.
	BillID              string                 `json:"billID"`
	RecorderUserID      string                 `json:"recorderUserID"`
	Title               *string                `json:"title,omitempty"`
	BillKind            BillKind               `json:"billKind"`
	Currency            CurrencyCode           `json:"currency"`
	ActualAmount        ExactDecimalString     `json:"actualAmount"`
	PaidAllocations     []BillAllocationV1     `json:"paidAllocations"`
	OwedAllocations     []BillAllocationV1     `json:"owedAllocations"`
	Utility             *UtilityDetailsV1      `json:"utility,omitempty"`
	RecurringOccurrence *RecurringOccurrenceV1 `json:"recurringOccurrence,omitempty"`
}

// BillServiceV1 is the storage-neutral port implemented by Splitus and wired
// by a host. authenticatedUserID comes from trusted host authentication.
type BillServiceV1 interface {
	CreateBill(context.Context, string, CreateBillV1Request) (CreateBillV1Response, error)
	GetBill(context.Context, string, GetBillV1Request) (GetBillV1Response, error)
	ListBills(context.Context, string, ListBillsV1Request) (ListBillsV1Response, error)
}

// ValidateAuthenticatedRecorder binds the request's audit claim to the
// trusted host authentication context. RecorderUserID grants no authority and
// never implies a paid or owed allocation.
func (r CreateBillV1Request) ValidateAuthenticatedRecorder(authenticatedUserID string) error {
	if err := validateStorageID("authenticatedUserID", authenticatedUserID); err != nil {
		return err
	}
	if r.RecorderUserID != authenticatedUserID {
		return invalid("recorderUserID must match the trusted authenticated identity")
	}
	return nil
}

func DecodeCreateBillV1Request(reader io.Reader) (CreateBillV1Request, error) {
	decoder := json.NewDecoder(reader)
	var request CreateBillV1Request
	if err := decoder.Decode(&request); err != nil {
		return CreateBillV1Request{}, fmt.Errorf("%w: decode: %v", ErrInvalidRequest, err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		if err == nil {
			err = errors.New("multiple JSON values")
		}
		return CreateBillV1Request{}, fmt.Errorf("%w: trailing JSON: %v", ErrInvalidRequest, err)
	}
	if err := request.Validate(); err != nil {
		return CreateBillV1Request{}, err
	}
	return request, nil
}

func (r CreateBillV1Request) Validate() error {
	if r.ContractVersion != BillContractVersion {
		return invalid("unsupported Splitus bill contract version")
	}
	for name, value := range map[string]string{
		"spaceID": r.SpaceID, "billID": r.BillID, "recorderUserID": r.RecorderUserID,
	} {
		if err := validateStorageID(name, value); err != nil {
			return err
		}
	}
	if r.Title != nil {
		if err := validateText("title", *r.Title, 256); err != nil {
			return err
		}
	}
	if r.BillKind != BillKindGeneral && r.BillKind != BillKindUtility {
		return invalid("unsupported billKind %q", r.BillKind)
	}
	if !isCurrencyCode(r.Currency) {
		return invalid("unsupported Splitus currency %q", r.Currency)
	}
	actual, err := positiveMinorUnits("actualAmount", r.ActualAmount)
	if err != nil {
		return err
	}
	if (r.BillKind == BillKindUtility) != (r.Utility != nil) {
		return invalid("utility details are required only for utility bills")
	}
	if r.Utility != nil {
		if err := r.Utility.Validate(); err != nil {
			return err
		}
	}
	if r.RecurringOccurrence != nil {
		if err := r.RecurringOccurrence.Validate(actual, r.BillID); err != nil {
			return err
		}
	}
	paid, paidContacts, err := validateAllocations("paidAllocations", r.PaidAllocations)
	if err != nil {
		return err
	}
	owed, owedContacts, err := validateAllocations("owedAllocations", r.OwedAllocations)
	if err != nil {
		return err
	}
	participants := make(map[string]struct{}, len(paidContacts)+len(owedContacts))
	for contactID := range paidContacts {
		participants[contactID] = struct{}{}
	}
	for contactID := range owedContacts {
		participants[contactID] = struct{}{}
	}
	if len(participants) < 2 || len(participants) > MaxBillParticipants {
		return invalid("bill must have 2 to %d participants", MaxBillParticipants)
	}
	if paid != actual {
		return invalid("paid allocations must reconcile to actualAmount")
	}
	if owed != actual {
		return invalid("owed allocations must reconcile to actualAmount")
	}
	return nil
}

func validateAllocations(name string, allocations []BillAllocationV1) (int64, map[string]struct{}, error) {
	if len(allocations) < 1 || len(allocations) > MaxBillParticipants {
		return 0, nil, invalid("%s must contain 1 to %d allocations", name, MaxBillParticipants)
	}
	allocationIDs := make(map[string]struct{}, len(allocations))
	contactIDs := make(map[string]struct{}, len(allocations))
	var total int64
	for i, allocation := range allocations {
		if err := validateStorageID(fmt.Sprintf("%s[%d].allocationID", name, i), allocation.AllocationID); err != nil {
			return 0, nil, err
		}
		if err := validateStorageID(fmt.Sprintf("%s[%d].contactID", name, i), allocation.ContactID); err != nil {
			return 0, nil, err
		}
		if _, exists := allocationIDs[allocation.AllocationID]; exists {
			return 0, nil, invalid("%s contains a duplicate allocation", name)
		}
		if _, exists := contactIDs[allocation.ContactID]; exists {
			return 0, nil, invalid("%s contains a duplicate contact", name)
		}
		allocationIDs[allocation.AllocationID] = struct{}{}
		contactIDs[allocation.ContactID] = struct{}{}
		amount, err := positiveMinorUnits(name, allocation.Amount)
		if err != nil {
			return 0, nil, err
		}
		if total > math.MaxInt64-amount {
			return 0, nil, invalid("%s total overflows signed 64-bit minor units", name)
		}
		total += amount
	}
	return total, contactIDs, nil
}

func (u UtilityDetailsV1) Validate() error {
	switch u.UtilityKind {
	case UtilityKindElectricity, UtilityKindGas, UtilityKindWater, UtilityKindInternet, UtilityKindOther:
	default:
		return invalid("unsupported utilityKind %q", u.UtilityKind)
	}
	if u.ProviderName != nil {
		if err := validateText("providerName", *u.ProviderName, 256); err != nil {
			return err
		}
	}
	return u.Period.Validate()
}

func (p BillingPeriodV1) Validate() error {
	start, err := time.Parse("2006-01-02", p.StartDate)
	if err != nil || start.Format("2006-01-02") != p.StartDate {
		return invalid("period.startDate must be a real ISO calendar date")
	}
	end, err := time.Parse("2006-01-02", p.EndDate)
	if err != nil || end.Format("2006-01-02") != p.EndDate {
		return invalid("period.endDate must be a real ISO calendar date")
	}
	if end.Before(start) {
		return invalid("period.endDate must not precede period.startDate")
	}
	return nil
}

func (r RecurringOccurrenceV1) Validate(actual int64, billID string) error {
	if err := validateStorageID("happeningID", r.HappeningID); err != nil {
		return err
	}
	if err := validateStorageID("occurrenceID", r.OccurrenceID); err != nil {
		return err
	}
	if (r.ExpectedAmount == nil) != (r.ExpectedComparison == ExpectedActualNotAvailable) {
		return invalid("comparison must be not_available exactly when expectedAmount is absent")
	}
	if r.ExpectedAmount != nil {
		expected, err := positiveMinorUnits("expectedAmount", *r.ExpectedAmount)
		if err != nil {
			return err
		}
		want := ExpectedActualMatches
		if actual > expected {
			want = ExpectedActualIncreased
		} else if actual < expected {
			want = ExpectedActualDecreased
		}
		if r.ExpectedComparison != want {
			return invalid("comparison does not match expected and actual amounts")
		}
	}
	if r.StandingChargeAmount != nil {
		charge, err := positiveMinorUnits("standingChargeAmount", *r.StandingChargeAmount)
		if err != nil {
			return err
		}
		if charge > actual {
			return invalid("standing charge cannot exceed actualAmount")
		}
	}
	if r.PreviousComparable != nil {
		if err := r.PreviousComparable.Validate(actual, billID); err != nil {
			return err
		}
	}
	return nil
}

func (p PreviousComparableBillV1) Validate(actual int64, billID string) error {
	if err := validateStorageID("previousComparable.billID", p.BillID); err != nil {
		return err
	}
	if p.BillID == billID {
		return invalid("previous comparable bill must have a different billID")
	}
	previous, err := positiveMinorUnits("previousComparable.actualAmount", p.ActualAmount)
	if err != nil {
		return err
	}
	want := ExpectedActualMatches
	if actual > previous {
		want = ExpectedActualIncreased
	} else if actual < previous {
		want = ExpectedActualDecreased
	}
	if p.Comparison != want {
		return invalid("previous comparison does not match current and prior actual amounts")
	}
	return nil
}

type DebtusReceiptLineV1 struct {
	LineID        string   `json:"lineID"`
	ObligationIDs []string `json:"obligationIDs"`
}

type BillPostingReceiptV1 struct {
	ReceiptID       string                `json:"receiptID"`
	OperationKey    string                `json:"operationKey"`
	InputDigest     string                `json:"inputDigest"`
	Revision        string                `json:"revision"`
	ObligationLines []DebtusReceiptLineV1 `json:"obligationLines"`
}

type BillPostingV1 struct {
	Status        PostingStatus         `json:"status"`
	OperationKey  string                `json:"operationKey"`
	InputDigest   string                `json:"inputDigest"`
	Receipt       *BillPostingReceiptV1 `json:"receipt,omitempty"`
	AttentionCode *BillAttentionCode    `json:"attentionCode,omitempty"`
}

func (p BillPostingV1) Validate(revision string) error {
	if err := validateStorageID("posting.operationKey", p.OperationKey); err != nil {
		return err
	}
	if err := validateDigest("posting.inputDigest", p.InputDigest); err != nil {
		return err
	}
	switch p.Status {
	case PostingStatusPending, PostingStatusPosting:
		if p.Receipt != nil || p.AttentionCode != nil {
			return invalid("%s posting cannot contain receipt or attentionCode", p.Status)
		}
	case PostingStatusApplied:
		if p.Receipt == nil || p.AttentionCode != nil {
			return invalid("applied posting requires a receipt and no attentionCode")
		}
		if err := p.Receipt.Validate(revision, p.OperationKey, p.InputDigest); err != nil {
			return err
		}
	case PostingStatusAttention:
		if p.Receipt != nil || p.AttentionCode == nil || !validAttentionCode(*p.AttentionCode) {
			return invalid("attention posting requires one supported attentionCode and no receipt")
		}
	default:
		return invalid("unsupported posting status %q", p.Status)
	}
	return nil
}

func (r BillPostingReceiptV1) Validate(revision, operationKey, inputDigest string) error {
	if err := validateStorageID("posting.receipt.receiptID", r.ReceiptID); err != nil {
		return err
	}
	if r.OperationKey != operationKey || r.InputDigest != inputDigest || r.Revision != revision {
		return invalid("posting receipt does not match the accepted bill revision")
	}
	if len(r.ObligationLines) > MaxBillParticipants {
		return invalid("posting receipt obligationLines exceed %d", MaxBillParticipants)
	}
	lineIDs := make(map[string]struct{}, len(r.ObligationLines))
	obligationIDs := make(map[string]struct{})
	for i, line := range r.ObligationLines {
		if err := validateStorageID(fmt.Sprintf("obligationLines[%d].lineID", i), line.LineID); err != nil {
			return err
		}
		if _, exists := lineIDs[line.LineID]; exists {
			return invalid("posting receipt repeats lineID %q", line.LineID)
		}
		lineIDs[line.LineID] = struct{}{}
		if len(line.ObligationIDs) < 1 || len(line.ObligationIDs) > MaxObligationIDsPerLine {
			return invalid("obligationLines[%d].obligationIDs must contain 1 to %d IDs", i, MaxObligationIDsPerLine)
		}
		for _, id := range line.ObligationIDs {
			if err := validateStorageID("obligationID", id); err != nil {
				return err
			}
			if _, exists := obligationIDs[id]; exists {
				return invalid("posting receipt repeats obligationID %q", id)
			}
			obligationIDs[id] = struct{}{}
		}
	}
	return nil
}

type DebtusSettlementTargetV1 struct {
	Route           string  `json:"route"`
	SpaceID         string  `json:"spaceID"`
	SourceNamespace string  `json:"sourceNamespace"`
	SourceRecordID  string  `json:"sourceRecordID"`
	LineID          *string `json:"lineID,omitempty"`
}

func (t DebtusSettlementTargetV1) Validate(spaceID, billID string, lineID *string) error {
	if t.Route != SettlementRoute || t.SourceNamespace != DebtusSourceNamespace || t.SpaceID != spaceID || t.SourceRecordID != billID {
		return invalid("settlement target does not match the Splitus bill source")
	}
	if lineID == nil {
		if t.LineID != nil {
			return invalid("bill settlement target cannot name an obligation line")
		}
		return nil
	}
	if t.LineID == nil || *t.LineID != *lineID {
		return invalid("obligation settlement target has a mismatched lineID")
	}
	return validateStorageID("settlementTarget.lineID", *t.LineID)
}

type DebtusObligationV1 struct {
	LineID            string                   `json:"lineID"`
	ObligationIDs     []string                 `json:"obligationIDs"`
	DebtorContactID   string                   `json:"debtorContactID"`
	CreditorContactID string                   `json:"creditorContactID"`
	Currency          CurrencyCode             `json:"currency"`
	PrincipalAmount   ExactDecimalString       `json:"principalAmount"`
	OutstandingAmount ExactDecimalString       `json:"outstandingAmount"`
	RepaidAmount      ExactDecimalString       `json:"repaidAmount"`
	CreditAmount      ExactDecimalString       `json:"creditAmount"`
	Status            DebtusSettlementStatus   `json:"status"`
	SettlementTarget  DebtusSettlementTargetV1 `json:"settlementTarget"`
}

type DebtusStatusV1 struct {
	Status           DebtusSettlementStatus   `json:"status"`
	Obligations      []DebtusObligationV1     `json:"obligations"`
	SettlementTarget DebtusSettlementTargetV1 `json:"settlementTarget"`
}

func (s DebtusStatusV1) Validate(spaceID, billID string, billCurrency CurrencyCode) error {
	if !validSettlementStatus(s.Status) {
		return invalid("unsupported Debtus settlement status %q", s.Status)
	}
	if len(s.Obligations) > MaxBillParticipants {
		return invalid("Debtus obligations exceed %d", MaxBillParticipants)
	}
	if err := s.SettlementTarget.Validate(spaceID, billID, nil); err != nil {
		return err
	}
	lineIDs := make(map[string]struct{}, len(s.Obligations))
	for i, obligation := range s.Obligations {
		if _, exists := lineIDs[obligation.LineID]; exists {
			return invalid("Debtus status repeats lineID %q", obligation.LineID)
		}
		lineIDs[obligation.LineID] = struct{}{}
		if err := obligation.Validate(spaceID, billID, billCurrency); err != nil {
			return fmt.Errorf("obligation %d: %w", i, err)
		}
	}
	return nil
}

func (o DebtusObligationV1) Validate(spaceID, billID string, billCurrency CurrencyCode) error {
	if err := validateStorageID("lineID", o.LineID); err != nil {
		return err
	}
	if len(o.ObligationIDs) < 1 || len(o.ObligationIDs) > MaxObligationIDsPerLine {
		return invalid("obligationIDs must contain 1 to %d IDs", MaxObligationIDsPerLine)
	}
	seen := make(map[string]struct{}, len(o.ObligationIDs))
	for _, id := range o.ObligationIDs {
		if err := validateStorageID("obligationID", id); err != nil {
			return err
		}
		if _, exists := seen[id]; exists {
			return invalid("duplicate obligationID %q", id)
		}
		seen[id] = struct{}{}
	}
	if err := validateStorageID("debtorContactID", o.DebtorContactID); err != nil {
		return err
	}
	if err := validateStorageID("creditorContactID", o.CreditorContactID); err != nil {
		return err
	}
	if o.DebtorContactID == o.CreditorContactID {
		return invalid("debtor and creditor contacts must differ")
	}
	if o.Currency != billCurrency || !isCurrencyCode(o.Currency) {
		return invalid("Debtus obligation currency differs from the bill")
	}
	if _, err := positiveMinorUnits("principalAmount", o.PrincipalAmount); err != nil {
		return err
	}
	for name, amount := range map[string]ExactDecimalString{
		"outstandingAmount": o.OutstandingAmount,
		"repaidAmount":      o.RepaidAmount,
		"creditAmount":      o.CreditAmount,
	} {
		if _, err := nonNegativeMinorUnits(name, amount); err != nil {
			return err
		}
	}
	if !validSettlementStatus(o.Status) {
		return invalid("unsupported obligation settlement status %q", o.Status)
	}
	return o.SettlementTarget.Validate(spaceID, billID, &o.LineID)
}

type BillV1 struct {
	CreateBillV1Request
	Revision  string          `json:"revision"`
	Posting   BillPostingV1   `json:"posting"`
	Debtus    *DebtusStatusV1 `json:"debtus,omitempty"`
	CreatedAt string          `json:"createdAt"`
	UpdatedAt string          `json:"updatedAt"`
}

func (b BillV1) Validate() error {
	if err := b.CreateBillV1Request.Validate(); err != nil {
		return err
	}
	if err := validatePositiveIntegerString("revision", b.Revision); err != nil {
		return err
	}
	createdAt, err := validateTimestamp("createdAt", b.CreatedAt)
	if err != nil {
		return err
	}
	updatedAt, err := validateTimestamp("updatedAt", b.UpdatedAt)
	if err != nil {
		return err
	}
	if updatedAt.Before(createdAt) {
		return invalid("updatedAt must not precede createdAt")
	}
	if err := b.Posting.Validate(b.Revision); err != nil {
		return err
	}
	if b.Debtus != nil {
		if b.Posting.Status != PostingStatusApplied {
			return invalid("Debtus status requires applied posting")
		}
		if err := b.Debtus.Validate(b.SpaceID, b.BillID, b.Currency); err != nil {
			return err
		}
		if b.Posting.Receipt == nil {
			return invalid("Debtus status requires an applied posting receipt")
		}
		if err := validateDebtusMatchesReceipt(*b.Posting.Receipt, *b.Debtus); err != nil {
			return err
		}
	}
	return nil
}

func validateDebtusMatchesReceipt(receipt BillPostingReceiptV1, debtus DebtusStatusV1) error {
	if len(receipt.ObligationLines) != len(debtus.Obligations) {
		return invalid("Debtus obligations do not exactly match the applied posting receipt")
	}
	receiptLines := make(map[string]map[string]struct{}, len(receipt.ObligationLines))
	for _, line := range receipt.ObligationLines {
		ids := make(map[string]struct{}, len(line.ObligationIDs))
		for _, id := range line.ObligationIDs {
			ids[id] = struct{}{}
		}
		receiptLines[line.LineID] = ids
	}
	for _, obligation := range debtus.Obligations {
		receiptIDs, exists := receiptLines[obligation.LineID]
		if !exists || len(receiptIDs) != len(obligation.ObligationIDs) {
			return invalid("Debtus obligations do not exactly match the applied posting receipt")
		}
		for _, id := range obligation.ObligationIDs {
			if _, exists := receiptIDs[id]; !exists {
				return invalid("Debtus obligations do not exactly match the applied posting receipt")
			}
		}
	}
	return nil
}

type CreateBillV1Response struct {
	ContractVersion int    `json:"contractVersion"`
	Bill            BillV1 `json:"bill"`
}

func (r CreateBillV1Response) Validate() error {
	if r.ContractVersion != BillContractVersion || r.Bill.ContractVersion != BillContractVersion {
		return invalid("unsupported Splitus bill contract version")
	}
	return r.Bill.Validate()
}

type GetBillV1Request struct {
	ContractVersion int    `json:"contractVersion"`
	SpaceID         string `json:"spaceID"`
	BillID          string `json:"billID"`
}

func (r GetBillV1Request) Validate() error {
	if r.ContractVersion != BillContractVersion {
		return invalid("unsupported Splitus bill contract version")
	}
	if err := validateStorageID("spaceID", r.SpaceID); err != nil {
		return err
	}
	return validateStorageID("billID", r.BillID)
}

type GetBillV1Response struct {
	ContractVersion int    `json:"contractVersion"`
	Bill            BillV1 `json:"bill"`
}

func (r GetBillV1Response) Validate() error {
	if r.ContractVersion != BillContractVersion || r.Bill.ContractVersion != BillContractVersion {
		return invalid("unsupported Splitus bill contract version")
	}
	return r.Bill.Validate()
}

type BillListItemV1 struct {
	ContractVersion        int                     `json:"contractVersion"`
	SpaceID                string                  `json:"spaceID"`
	BillID                 string                  `json:"billID"`
	Title                  *string                 `json:"title,omitempty"`
	BillKind               BillKind                `json:"billKind"`
	UtilityKind            *UtilityKind            `json:"utilityKind,omitempty"`
	Period                 *BillingPeriodV1        `json:"period,omitempty"`
	Currency               CurrencyCode            `json:"currency"`
	ActualAmount           ExactDecimalString      `json:"actualAmount"`
	OwnPaidAmount          ExactDecimalString      `json:"ownPaidAmount"`
	OwnOwedAmount          ExactDecimalString      `json:"ownOwedAmount"`
	PostingStatus          PostingStatus           `json:"postingStatus"`
	DebtusSettlementStatus *DebtusSettlementStatus `json:"debtusSettlementStatus,omitempty"`
	CreatedAt              string                  `json:"createdAt"`
}

func (b BillListItemV1) Validate() error {
	if b.ContractVersion != BillContractVersion {
		return invalid("unsupported Splitus bill contract version")
	}
	if err := validateStorageID("spaceID", b.SpaceID); err != nil {
		return err
	}
	if err := validateStorageID("billID", b.BillID); err != nil {
		return err
	}
	if b.Title != nil {
		if err := validateText("title", *b.Title, 256); err != nil {
			return err
		}
	}
	if b.BillKind != BillKindGeneral && b.BillKind != BillKindUtility {
		return invalid("unsupported billKind %q", b.BillKind)
	}
	if (b.BillKind == BillKindUtility) != (b.UtilityKind != nil && b.Period != nil) {
		return invalid("utility list items require utilityKind and period")
	}
	if b.UtilityKind != nil {
		u := UtilityDetailsV1{UtilityKind: *b.UtilityKind, Period: *b.Period}
		if err := u.Validate(); err != nil {
			return err
		}
	}
	if !isCurrencyCode(b.Currency) {
		return invalid("unsupported Splitus currency %q", b.Currency)
	}
	if _, err := positiveMinorUnits("actualAmount", b.ActualAmount); err != nil {
		return err
	}
	if _, err := nonNegativeMinorUnits("ownPaidAmount", b.OwnPaidAmount); err != nil {
		return err
	}
	if _, err := nonNegativeMinorUnits("ownOwedAmount", b.OwnOwedAmount); err != nil {
		return err
	}
	if !validPostingStatus(b.PostingStatus) {
		return invalid("unsupported postingStatus %q", b.PostingStatus)
	}
	if b.DebtusSettlementStatus != nil && !validSettlementStatus(*b.DebtusSettlementStatus) {
		return invalid("unsupported Debtus settlement status %q", *b.DebtusSettlementStatus)
	}
	if b.DebtusSettlementStatus != nil && b.PostingStatus != PostingStatusApplied {
		return invalid("Debtus settlement status requires applied posting")
	}
	_, err := validateTimestamp("createdAt", b.CreatedAt)
	return err
}

type ListBillsV1Request struct {
	ContractVersion int              `json:"contractVersion"`
	SpaceID         string           `json:"spaceID"`
	PageSize        int              `json:"pageSize"`
	Cursor          *string          `json:"cursor,omitempty"`
	UtilityKind     *UtilityKind     `json:"utilityKind,omitempty"`
	Period          *BillingPeriodV1 `json:"period,omitempty"`
}

func (r ListBillsV1Request) Validate() error {
	if r.ContractVersion != BillContractVersion {
		return invalid("unsupported Splitus bill contract version")
	}
	if err := validateStorageID("spaceID", r.SpaceID); err != nil {
		return err
	}
	if r.PageSize < 1 || r.PageSize > MaxBillListPageSize {
		return invalid("pageSize must be from 1 to %d", MaxBillListPageSize)
	}
	if r.Cursor != nil {
		if err := validateText("cursor", *r.Cursor, 2048); err != nil {
			return err
		}
	}
	if r.UtilityKind != nil {
		u := UtilityDetailsV1{UtilityKind: *r.UtilityKind, Period: BillingPeriodV1{StartDate: "2000-01-01", EndDate: "2000-01-01"}}
		if err := u.Validate(); err != nil {
			return err
		}
	}
	if r.Period != nil {
		return r.Period.Validate()
	}
	return nil
}

type ListBillsV1Response struct {
	ContractVersion int              `json:"contractVersion"`
	PageSize        int              `json:"pageSize"`
	Items           []BillListItemV1 `json:"items"`
	NextCursor      *string          `json:"nextCursor,omitempty"`
}

func (r ListBillsV1Response) Validate(requestedPageSize int) error {
	if r.ContractVersion != BillContractVersion {
		return invalid("unsupported Splitus bill contract version")
	}
	if r.PageSize < 1 || r.PageSize > MaxBillListPageSize || r.PageSize != requestedPageSize {
		return invalid("response pageSize does not match a bounded request")
	}
	if len(r.Items) > r.PageSize {
		return invalid("response items exceed pageSize")
	}
	for i, item := range r.Items {
		if err := item.Validate(); err != nil {
			return fmt.Errorf("item %d: %w", i, err)
		}
	}
	if r.NextCursor != nil {
		return validateText("nextCursor", *r.NextCursor, 2048)
	}
	return nil
}

func positiveMinorUnits(name string, amount ExactDecimalString) (int64, error) {
	minor, err := amount.MinorUnits()
	if err != nil {
		return 0, err
	}
	if minor <= 0 {
		return 0, invalid("%s must be positive", name)
	}
	return minor, nil
}

func nonNegativeMinorUnits(name string, amount ExactDecimalString) (int64, error) {
	minor, err := amount.MinorUnits()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", name, err)
	}
	return minor, nil
}

func validatePositiveIntegerString(name, value string) error {
	if value == "" || value[0] == '0' {
		return invalid("%s must be a canonical positive integer string", name)
	}
	if _, err := strconv.ParseUint(value, 10, 64); err != nil {
		return invalid("%s must fit unsigned 64-bit range", name)
	}
	return nil
}

func validateTimestamp(name, value string) (time.Time, error) {
	parsed, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		return time.Time{}, invalid("%s must be an RFC 3339 timestamp", name)
	}
	return parsed, nil
}

func validateDigest(name, value string) error {
	if len(value) != 64 {
		return invalid("%s must be a lowercase SHA-256 digest", name)
	}
	for _, char := range value {
		if !(char >= '0' && char <= '9' || char >= 'a' && char <= 'f') {
			return invalid("%s must be a lowercase SHA-256 digest", name)
		}
	}
	return nil
}

func validPostingStatus(status PostingStatus) bool {
	return status == PostingStatusPending || status == PostingStatusPosting || status == PostingStatusApplied || status == PostingStatusAttention
}

func validSettlementStatus(status DebtusSettlementStatus) bool {
	return status == DebtusSettlementStatusUnsettled || status == DebtusSettlementStatusPartSettled || status == DebtusSettlementStatusSettled
}

func validAttentionCode(code BillAttentionCode) bool {
	switch code {
	case BillAttentionAuthorizationChanged, BillAttentionSourceConflict, BillAttentionProviderRejected, BillAttentionInvalidReceipt, BillAttentionOperatorAction:
		return true
	default:
		return false
	}
}

func validateStorageID(name, value string) error {
	if !utf8.ValidString(value) || value == "" || strings.TrimSpace(value) != value || len(value) > 512 {
		return invalid("%s is empty, padded, too long, or invalid UTF-8", name)
	}
	for _, r := range value {
		if r < 0x20 || r == 0x7f {
			return invalid("%s contains a control character", name)
		}
	}
	if strings.Contains(value, "/") || value == "." || value == ".." || strings.HasPrefix(value, "__") && strings.HasSuffix(value, "__") {
		return invalid("%s is not a safe identifier", name)
	}
	return nil
}

func validateText(name, value string, maxBytes int) error {
	if !utf8.ValidString(value) || value == "" || strings.TrimSpace(value) != value || len(value) > maxBytes {
		return invalid("%s is empty, padded, too long, or invalid UTF-8", name)
	}
	for _, r := range value {
		if r < 0x20 || r == 0x7f {
			return invalid("%s contains a control character", name)
		}
	}
	return nil
}

func isCurrencyCode(value CurrencyCode) bool {
	return value == CurrencyEUR || value == CurrencyGBP || value == CurrencyUSD
}

func invalid(format string, args ...any) error {
	return fmt.Errorf("%w: %s", ErrInvalidRequest, fmt.Sprintf(format, args...))
}
