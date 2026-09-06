package calendariusmodels

import (
	"context"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	MaxFinancialAgreementParties     = 32
	MaxFinancialAgreementPageSize    = 64
	MaxFinancialAgreementCursorBytes = 512

	FinancialAgreementStateOffered          = "offered"
	FinancialAgreementStateRecordedExternal = "recorded_external"
	FinancialAgreementStateConfirmed        = "confirmed"
	FinancialAgreementStateEnded            = "ended"
	FinancialAgreementStateCancelled        = "cancelled"
	FinancialAgreementSidePayer             = "payer"
	FinancialAgreementSideReceiver          = "receiver"
	FinancialPartyKindSpace                 = "space"
	FinancialPartyKindContact               = "contact"
)

// FinancialPartyRef identifies a contractual economic party. A contact's
// SpaceID is storage addressing only, not its economic-interest owner.
type FinancialPartyRef struct {
	Kind      string `json:"kind"`
	SpaceID   string `json:"spaceID"`
	ContactID string `json:"contactID,omitempty"`
}

func (p FinancialPartyRef) Validate() error {
	if p.Kind != FinancialPartyKindSpace && p.Kind != FinancialPartyKindContact {
		return fmt.Errorf("financial party kind must be space or contact")
	}
	if err := validateFinancialCommitmentID("party.spaceID", p.SpaceID); err != nil {
		return err
	}
	if p.Kind == FinancialPartyKindSpace {
		if p.ContactID != "" {
			return fmt.Errorf("a Space party cannot name a contact")
		}
		return nil
	}
	return validateFinancialCommitmentID("party.contactID", p.ContactID)
}

type FinancialPartyShare struct {
	Party       FinancialPartyRef `json:"party"`
	AmountMinor int64             `json:"amountMinor"`
}

// FinancialPartyWeight is untrusted allocation intent. The owner resolves it
// to exact minor-unit shares from the selected catalog terms.
type FinancialPartyWeight struct {
	Party  FinancialPartyRef `json:"party"`
	Shares int64             `json:"shares"`
}

type FinancialAttribution struct {
	ID          string `json:"id"`
	AmountMinor int64  `json:"amountMinor"`
}

type FinancialAttributionWeight struct {
	ID     string `json:"id"`
	Shares int64  `json:"shares"`
}

// AcceptedPriceTerms is the immutable snapshot of one selected catalog option.
// Quantity belongs to this agreement, never the whole happening.
type AcceptedPriceTerms struct {
	PriceID       string                  `json:"priceID"`
	PriceRevision int64                   `json:"priceRevision"`
	AmountMinor   int64                   `json:"amountMinor"`
	Currency      string                  `json:"currency"`
	Quantity      int64                   `json:"quantity"`
	Term          FinancialCommitmentTerm `json:"term"`
}

func (p AcceptedPriceTerms) Validate() error {
	if err := validateFinancialCommitmentID("priceID", p.PriceID); err != nil {
		return err
	}
	if p.PriceRevision < 0 || p.PriceRevision > MaxJavaScriptSafeInteger ||
		p.AmountMinor < 0 || p.AmountMinor > MaxJavaScriptSafeInteger ||
		p.Quantity < 1 || p.Quantity > MaxJavaScriptSafeInteger ||
		(p.AmountMinor > 0 && p.Quantity > MaxJavaScriptSafeInteger/p.AmountMinor) {
		return fmt.Errorf("accepted price values are outside supported bounds")
	}
	if p.Currency != "EUR" && p.Currency != "GBP" && p.Currency != "USD" {
		return fmt.Errorf("accepted price currency is unsupported")
	}
	if !validAgreementTerm(p.Term.Unit) || p.Term.Length < 1 || p.Term.Length > MaxJavaScriptSafeInteger {
		return fmt.Errorf("accepted price term is unsupported")
	}
	if p.Term.Unit == "single" && p.Term.Length != 1 {
		return fmt.Errorf("single term length must be one")
	}
	return nil
}

func validAgreementTerm(unit string) bool {
	switch unit {
	case "single", "second", "minute", "hour", "day", "week", "month", "quarter", "year":
		return true
	default:
		return false
	}
}

// FinancialConfirmation is one party's explicit evidence for an immutable
// terms revision. Later audit/reconciliation revisions do not invalidate it.
type FinancialConfirmation struct {
	Party         FinancialPartyRef `json:"party"`
	TermsRevision int64             `json:"termsRevision"`
	ActorUserID   string            `json:"actorUserID"`
	ConfirmedAt   string            `json:"confirmedAt"`
	RevokedAt     string            `json:"revokedAt,omitempty"`
}

func (c FinancialConfirmation) Validate(termsRevision int64) error {
	if err := c.Party.Validate(); err != nil {
		return err
	}
	if c.TermsRevision != termsRevision || strings.TrimSpace(c.ActorUserID) == "" {
		return fmt.Errorf("confirmation is not bound to this terms revision")
	}
	if _, err := time.Parse(time.RFC3339, c.ConfirmedAt); err != nil {
		return fmt.Errorf("confirmedAt is invalid")
	}
	if c.RevokedAt != "" {
		if _, err := time.Parse(time.RFC3339, c.RevokedAt); err != nil {
			return fmt.Errorf("revokedAt is invalid")
		}
	}
	return nil
}

// FinancialAgreementFact is an authorized reporting view of a source-owned
// selected deal. OwnerSpaceID scopes identity; ReportingSpaceID scopes access
// and totals. Cross-Space views must omit unrelated source facts.
type FinancialAgreementFact struct {
	OwnerSpaceID         string                  `json:"ownerSpaceID"`
	ReportingSpaceID     string                  `json:"reportingSpaceID"`
	AgreementID          string                  `json:"agreementID"`
	EnrollmentID         string                  `json:"enrollmentID"`
	HappeningID          string                  `json:"happeningID"`
	Revision             int64                   `json:"revision"`
	TermsRevision        int64                   `json:"termsRevision"`
	State                string                  `json:"state"`
	ReportingSpaceSide   string                  `json:"reportingSpaceSide"`
	ReportingAmountMinor int64                   `json:"reportingAmountMinor"`
	Payers               []FinancialPartyShare   `json:"payers"`
	Receivers            []FinancialPartyShare   `json:"receivers"`
	ContactAttributions  []FinancialAttribution  `json:"contactAttributions"`
	AssetAttributions    []FinancialAttribution  `json:"assetAttributions"`
	AcceptedPrice        AcceptedPriceTerms      `json:"acceptedPrice"`
	EffectiveFromISO     string                  `json:"effectiveFromISO"`
	EffectiveToISO       string                  `json:"effectiveToISO,omitempty"`
	ExternallyRecorded   bool                    `json:"externallyRecorded,omitempty"`
	Confirmations        []FinancialConfirmation `json:"confirmations"`
}

func (f FinancialAgreementFact) Validate() error {
	for name, value := range map[string]string{
		"ownerSpaceID": f.OwnerSpaceID, "reportingSpaceID": f.ReportingSpaceID,
		"agreementID": f.AgreementID, "enrollmentID": f.EnrollmentID, "happeningID": f.HappeningID,
	} {
		if err := validateFinancialCommitmentID(name, value); err != nil {
			return err
		}
	}
	if f.Revision < 1 || f.Revision > MaxJavaScriptSafeInteger ||
		f.TermsRevision < 1 || f.TermsRevision > f.Revision || !validAgreementState(f.State) {
		return fmt.Errorf("agreement revision or state is invalid")
	}
	if f.ReportingSpaceSide != FinancialAgreementSidePayer && f.ReportingSpaceSide != FinancialAgreementSideReceiver {
		return fmt.Errorf("reportingSpaceSide is invalid")
	}
	if !validISODate(f.EffectiveFromISO) || !validOptionalISODate(f.EffectiveToISO) ||
		(f.EffectiveToISO != "" && f.EffectiveToISO < f.EffectiveFromISO) {
		return fmt.Errorf("agreement effective dates are invalid")
	}
	if err := f.AcceptedPrice.Validate(); err != nil {
		return err
	}
	total := f.AcceptedPrice.AmountMinor * f.AcceptedPrice.Quantity
	if f.ReportingAmountMinor < 1 || f.ReportingAmountMinor > total {
		return fmt.Errorf("reporting amount is invalid")
	}
	if err := validateAgreementShares(total, f.Payers, f.Receivers); err != nil {
		return err
	}
	if err := validateReportingSpaceSide(f.ReportingSpaceID, f.ReportingSpaceSide, f.Payers, f.Receivers); err != nil {
		return err
	}
	if err := validateAttributions(f.ReportingAmountMinor, "contactAttributions", f.ContactAttributions); err != nil {
		return err
	}
	if err := validateAttributions(f.ReportingAmountMinor, "assetAttributions", f.AssetAttributions); err != nil {
		return err
	}
	if f.Confirmations == nil {
		return fmt.Errorf("confirmations must be a non-nil array")
	}
	seen := map[string]struct{}{}
	for i, confirmation := range f.Confirmations {
		if err := confirmation.Validate(f.TermsRevision); err != nil {
			return fmt.Errorf("confirmations[%d]: %w", i, err)
		}
		key := partyKey(confirmation.Party)
		if _, exists := seen[key]; exists {
			return fmt.Errorf("confirmations contains duplicate party")
		}
		seen[key] = struct{}{}
	}
	if f.State == FinancialAgreementStateConfirmed {
		for _, side := range [][]FinancialPartyShare{f.Payers, f.Receivers} {
			for _, share := range side {
				_, exists := seen[partyKey(share.Party)]
				if !exists || confirmationRevoked(f.Confirmations, share.Party) {
					return fmt.Errorf("confirmed agreement lacks active evidence for every party")
				}
			}
		}
	}
	if f.State == FinancialAgreementStateRecordedExternal && !f.ExternallyRecorded {
		return fmt.Errorf("recorded external agreement must disclose provenance")
	}
	return nil
}

func confirmationRevoked(confirmations []FinancialConfirmation, party FinancialPartyRef) bool {
	for _, confirmation := range confirmations {
		if partyKey(confirmation.Party) == partyKey(party) {
			return confirmation.RevokedAt != ""
		}
	}
	return false
}

func validAgreementState(state string) bool {
	switch state {
	case FinancialAgreementStateOffered, FinancialAgreementStateRecordedExternal,
		FinancialAgreementStateConfirmed, FinancialAgreementStateEnded, FinancialAgreementStateCancelled:
		return true
	default:
		return false
	}
}

func partyKey(p FinancialPartyRef) string { return p.Kind + "\x00" + p.SpaceID + "\x00" + p.ContactID }

func validateAgreementShares(total int64, payers, receivers []FinancialPartyShare) error {
	if len(payers) == 0 || len(receivers) == 0 || len(payers) > MaxFinancialAgreementParties || len(receivers) > MaxFinancialAgreementParties {
		return fmt.Errorf("payer and receiver parties are required and bounded")
	}
	seen := map[string]struct{}{}
	for _, side := range []struct {
		name string
		rows []FinancialPartyShare
	}{{"payer", payers}, {"receiver", receivers}} {
		allocated := int64(0)
		for _, row := range side.rows {
			if err := row.Party.Validate(); err != nil {
				return fmt.Errorf("%s: %w", side.name, err)
			}
			key := partyKey(row.Party)
			if _, exists := seen[key]; exists {
				return fmt.Errorf("financial party appears more than once")
			}
			seen[key] = struct{}{}
			if row.AmountMinor < 1 || row.AmountMinor > total-allocated {
				return fmt.Errorf("%s shares do not conserve the agreement total", side.name)
			}
			allocated += row.AmountMinor
		}
		if allocated != total {
			return fmt.Errorf("%s shares do not equal the agreement total", side.name)
		}
	}
	return nil
}

func validateReportingSpaceSide(spaceID, side string, payers, receivers []FinancialPartyShare) error {
	contains := func(rows []FinancialPartyShare) bool {
		for _, row := range rows {
			if row.Party.Kind == FinancialPartyKindSpace && row.Party.SpaceID == spaceID {
				return true
			}
		}
		return false
	}
	onPayer, onReceiver := contains(payers), contains(receivers)
	if onPayer && onReceiver {
		return fmt.Errorf("same reporting Space on both sides is internal movement")
	}
	if side == FinancialAgreementSidePayer && !onPayer || side == FinancialAgreementSideReceiver && !onReceiver {
		return fmt.Errorf("reporting Space is not an explicit party on its reported side")
	}
	return nil
}

func validateAttributions(total int64, name string, rows []FinancialAttribution) error {
	if len(rows) > MaxFinancialAgreementParties {
		return fmt.Errorf("%s exceeds maximum", name)
	}
	seen, allocated := map[string]struct{}{}, int64(0)
	for _, row := range rows {
		if err := validateFinancialCommitmentID(name+".id", row.ID); err != nil {
			return err
		}
		if _, exists := seen[row.ID]; exists {
			return fmt.Errorf("%s has duplicate ID", name)
		}
		seen[row.ID] = struct{}{}
		if row.AmountMinor < 1 || row.AmountMinor > total-allocated {
			return fmt.Errorf("%s amount is invalid", name)
		}
		allocated += row.AmountMinor
	}
	if len(rows) > 0 && allocated != total {
		return fmt.Errorf("%s does not conserve the reporting amount", name)
	}
	return nil
}

func validatePartyWeights(name string, rows []FinancialPartyWeight) error {
	if len(rows) == 0 || len(rows) > MaxFinancialAgreementParties {
		return fmt.Errorf("%s are required and bounded", name)
	}
	seen := map[string]struct{}{}
	for _, row := range rows {
		if err := row.Party.Validate(); err != nil {
			return err
		}
		key := partyKey(row.Party)
		if _, exists := seen[key]; exists {
			return fmt.Errorf("%s has a duplicate party", name)
		}
		seen[key] = struct{}{}
		if row.Shares < 1 || row.Shares > MaxJavaScriptSafeInteger {
			return fmt.Errorf("%s shares are invalid", name)
		}
	}
	return nil
}

func validateAttributionWeights(name string, rows []FinancialAttributionWeight) error {
	if len(rows) > MaxFinancialAgreementParties {
		return fmt.Errorf("%s exceeds maximum", name)
	}
	seen := map[string]struct{}{}
	for _, row := range rows {
		if err := validateFinancialCommitmentID(name+".id", row.ID); err != nil {
			return err
		}
		if _, exists := seen[row.ID]; exists {
			return fmt.Errorf("%s has duplicate ID", name)
		}
		seen[row.ID] = struct{}{}
		if row.Shares < 1 || row.Shares > MaxJavaScriptSafeInteger {
			return fmt.Errorf("%s shares are invalid", name)
		}
	}
	return nil
}

type FinancialAgreementQuery struct {
	ReportingSpaceID string `json:"reportingSpaceID"`
	FromMonthISO     string `json:"fromMonthISO"`
	Months           int    `json:"months"`
	PageSize         int    `json:"pageSize"`
	Cursor           string `json:"cursor,omitempty"`
}

func (q FinancialAgreementQuery) Validate() error {
	if err := validateFinancialCommitmentID("reportingSpaceID", q.ReportingSpaceID); err != nil {
		return err
	}
	if !validMonthISO(q.FromMonthISO) || q.Months < 1 || q.Months > MaxFinancialCommitmentMonths {
		return fmt.Errorf("agreement query window is invalid")
	}
	from, _ := time.Parse("2006-01", q.FromMonthISO)
	if from.AddDate(0, q.Months-1, 0).Year() > 9999 || q.PageSize < 1 || q.PageSize > MaxFinancialAgreementPageSize {
		return fmt.Errorf("agreement query window or page is invalid")
	}
	return validateAgreementCursor(q.Cursor)
}

func validateAgreementCursor(cursor string) error {
	if len(cursor) > MaxFinancialAgreementCursorBytes || !utf8.ValidString(cursor) || strings.TrimSpace(cursor) != cursor || strings.ContainsAny(cursor, "\r\n") {
		return fmt.Errorf("agreement cursor is invalid")
	}
	return nil
}

type FinancialAgreementResult struct {
	Agreements         []FinancialAgreementFact `json:"agreements"`
	HasMore            bool                     `json:"hasMore"`
	NextCursor         string                   `json:"nextCursor,omitempty"`
	SnapshotConsistent bool                     `json:"snapshotConsistent"`
	IncompleteReason   string                   `json:"incompleteReason,omitempty"`
}

func (r FinancialAgreementResult) Validate() error {
	if r.Agreements == nil || len(r.Agreements) > MaxFinancialAgreementPageSize || r.HasMore != (r.NextCursor != "") {
		return fmt.Errorf("agreement result paging is invalid")
	}
	if err := validateAgreementCursor(r.NextCursor); err != nil {
		return err
	}
	if r.IncompleteReason != "" && r.IncompleteReason != "query_limit" && r.IncompleteReason != "source_collection_mutable_between_pages" {
		return fmt.Errorf("agreement result incompleteReason is invalid")
	}
	if r.SnapshotConsistent == (r.IncompleteReason != "") {
		return fmt.Errorf("agreement result completeness is invalid")
	}
	seen := map[string]struct{}{}
	for i, agreement := range r.Agreements {
		if err := agreement.Validate(); err != nil {
			return fmt.Errorf("agreements[%d]: %w", i, err)
		}
		key := agreement.OwnerSpaceID + "\x00" + agreement.AgreementID
		if _, exists := seen[key]; exists {
			return fmt.Errorf("duplicate owner agreement identity")
		}
		seen[key] = struct{}{}
	}
	return nil
}

type FinancialAgreementSource interface {
	ListFinancialAgreements(context.Context, string, FinancialAgreementQuery) (FinancialAgreementResult, error)
}

// SaveFinancialAgreementRequest selects one option. Clients submit weights;
// the owner resolves all exact money inside the audited command.
type SaveFinancialAgreementRequest struct {
	OwnerSpaceID          string                       `json:"ownerSpaceID"`
	ReportingSpaceID      string                       `json:"reportingSpaceID"`
	AgreementID           string                       `json:"agreementID"`
	EnrollmentID          string                       `json:"enrollmentID"`
	HappeningID           string                       `json:"happeningID"`
	PriceID               string                       `json:"priceID"`
	ExpectedPriceRevision int64                        `json:"expectedPriceRevision"`
	Quantity              int64                        `json:"quantity"`
	OperationID           string                       `json:"operationID"`
	ExpectedRevision      int64                        `json:"expectedRevision"`
	State                 string                       `json:"state"`
	ReportingSpaceSide    string                       `json:"reportingSpaceSide"`
	Payers                []FinancialPartyWeight       `json:"payers"`
	Receivers             []FinancialPartyWeight       `json:"receivers"`
	ContactAttributions   []FinancialAttributionWeight `json:"contactAttributions"`
	AssetAttributions     []FinancialAttributionWeight `json:"assetAttributions"`
	EffectiveFromISO      string                       `json:"effectiveFromISO"`
	EffectiveToISO        string                       `json:"effectiveToISO,omitempty"`
	Reason                string                       `json:"reason,omitempty"`
}

func (r SaveFinancialAgreementRequest) Validate() error {
	for name, value := range map[string]string{
		"ownerSpaceID": r.OwnerSpaceID, "reportingSpaceID": r.ReportingSpaceID,
		"agreementID": r.AgreementID, "enrollmentID": r.EnrollmentID,
		"happeningID": r.HappeningID, "priceID": r.PriceID, "operationID": r.OperationID,
	} {
		if err := validateFinancialCommitmentID(name, value); err != nil {
			return err
		}
	}
	if r.ExpectedRevision < 0 || r.ExpectedRevision > MaxJavaScriptSafeInteger || r.ExpectedPriceRevision < 0 ||
		r.ExpectedPriceRevision > MaxJavaScriptSafeInteger || r.Quantity < 1 || r.Quantity > MaxJavaScriptSafeInteger {
		return fmt.Errorf("expected revisions or quantity are outside supported bounds")
	}
	if r.State != FinancialAgreementStateOffered && r.State != FinancialAgreementStateRecordedExternal {
		return fmt.Errorf("save agreement state must be offered or recorded_external")
	}
	if r.ReportingSpaceSide != FinancialAgreementSidePayer && r.ReportingSpaceSide != FinancialAgreementSideReceiver {
		return fmt.Errorf("reportingSpaceSide is invalid")
	}
	if !validISODate(r.EffectiveFromISO) || !validOptionalISODate(r.EffectiveToISO) ||
		(r.EffectiveToISO != "" && r.EffectiveToISO < r.EffectiveFromISO) || len(r.Reason) > 500 {
		return fmt.Errorf("agreement dates or reason are invalid")
	}
	if err := validatePartyWeights("payers", r.Payers); err != nil {
		return err
	}
	if err := validatePartyWeights("receivers", r.Receivers); err != nil {
		return err
	}
	seenParties := map[string]string{}
	for _, side := range []struct {
		name string
		rows []FinancialPartyWeight
	}{{FinancialAgreementSidePayer, r.Payers}, {FinancialAgreementSideReceiver, r.Receivers}} {
		for _, row := range side.rows {
			key := partyKey(row.Party)
			if existingSide, exists := seenParties[key]; exists {
				return fmt.Errorf("financial party appears on both %s and %s sides", existingSide, side.name)
			}
			seenParties[key] = side.name
		}
	}
	reportingSide, exists := seenParties[partyKey(FinancialPartyRef{Kind: FinancialPartyKindSpace, SpaceID: r.ReportingSpaceID})]
	if !exists || reportingSide != r.ReportingSpaceSide {
		return fmt.Errorf("reporting Space is not an explicit party on its reported side")
	}
	if err := validateAttributionWeights("contactAttributions", r.ContactAttributions); err != nil {
		return err
	}
	return validateAttributionWeights("assetAttributions", r.AssetAttributions)
}

type SaveFinancialAgreementResponse struct {
	AgreementID string `json:"agreementID"`
	Revision    int64  `json:"revision"`
	UpdatedAt   string `json:"updatedAt"`
}

type ConfirmFinancialAgreementRequest struct {
	ReportingSpaceID string            `json:"reportingSpaceID"`
	OwnerSpaceID     string            `json:"ownerSpaceID"`
	AgreementID      string            `json:"agreementID"`
	ExpectedRevision int64             `json:"expectedRevision"`
	TermsRevision    int64             `json:"termsRevision"`
	OperationID      string            `json:"operationID"`
	Party            FinancialPartyRef `json:"party"`
}

func (r ConfirmFinancialAgreementRequest) Validate() error {
	for name, value := range map[string]string{
		"reportingSpaceID": r.ReportingSpaceID, "ownerSpaceID": r.OwnerSpaceID,
		"agreementID": r.AgreementID, "operationID": r.OperationID,
	} {
		if err := validateFinancialCommitmentID(name, value); err != nil {
			return err
		}
	}
	if r.ExpectedRevision < 1 || r.ExpectedRevision > MaxJavaScriptSafeInteger ||
		r.TermsRevision < 1 || r.TermsRevision > r.ExpectedRevision {
		return fmt.Errorf("agreement or terms revision is invalid")
	}
	return r.Party.Validate()
}

type RevokeFinancialAgreementGrantRequest struct {
	ReportingSpaceID string            `json:"reportingSpaceID"`
	OwnerSpaceID     string            `json:"ownerSpaceID"`
	AgreementID      string            `json:"agreementID"`
	ExpectedRevision int64             `json:"expectedRevision"`
	TermsRevision    int64             `json:"termsRevision"`
	OperationID      string            `json:"operationID"`
	Party            FinancialPartyRef `json:"party"`
	Reason           string            `json:"reason,omitempty"`
}

func (r RevokeFinancialAgreementGrantRequest) Validate() error {
	confirmation := ConfirmFinancialAgreementRequest{
		ReportingSpaceID: r.ReportingSpaceID, OwnerSpaceID: r.OwnerSpaceID,
		AgreementID: r.AgreementID, ExpectedRevision: r.ExpectedRevision,
		TermsRevision: r.TermsRevision, OperationID: r.OperationID, Party: r.Party,
	}
	if err := confirmation.Validate(); err != nil {
		return err
	}
	if len(r.Reason) > 500 {
		return fmt.Errorf("reason exceeds maximum length")
	}
	return nil
}
