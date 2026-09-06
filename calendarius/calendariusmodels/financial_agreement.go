package calendariusmodels

import (
	"context"
	"fmt"
	"time"
)

const (
	MaxFinancialAgreementParties            = 32
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

// FinancialPartyRef identifies a contractual party. For a contact, SpaceID is
// storage addressing only; it does not assert that the contact economically
// belongs to that Space. A Space party is the explicit reporting interest.
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
	if p.Kind == FinancialPartyKindSpace && p.ContactID != "" {
		return fmt.Errorf("a Space financial party cannot name a contact")
	}
	if p.Kind == FinancialPartyKindContact {
		return validateFinancialCommitmentID("party.contactID", p.ContactID)
	}
	return nil
}

// AcceptedPriceTerms is the immutable terms snapshot selected from a happening
// price catalog. Later catalog corrections do not rewrite this agreement.
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
	if p.PriceRevision < 1 || p.PriceRevision > MaxJavaScriptSafeInteger || p.AmountMinor < 0 || p.AmountMinor > MaxJavaScriptSafeInteger || p.Quantity < 1 || p.Quantity > MaxJavaScriptSafeInteger {
		return fmt.Errorf("accepted price values are outside supported bounds")
	}
	if len(p.Currency) != 3 || p.Term.Unit == "" || p.Term.Length < 1 {
		return fmt.Errorf("accepted price currency and term are required")
	}
	return nil
}

// FinancialAgreementFact is one source-owned selected deal. Offered catalog
// alternatives are deliberately absent. OwnSpaceSide determines whether the
// same agreed terms project as expense or income for this authorized Space.
type FinancialAgreementFact struct {
	SpaceID            string              `json:"spaceID"`
	AgreementID        string              `json:"agreementID"`
	EnrollmentID       string              `json:"enrollmentID"`
	HappeningID        string              `json:"happeningID"`
	Revision           int64               `json:"revision"`
	State              string              `json:"state"`
	OwnSpaceSide       string              `json:"ownSpaceSide"`
	Payers             []FinancialPartyRef `json:"payers"`
	Receivers          []FinancialPartyRef `json:"receivers"`
	AcceptedPrice      AcceptedPriceTerms  `json:"acceptedPrice"`
	EffectiveFromISO   string              `json:"effectiveFromISO"`
	EffectiveToISO     string              `json:"effectiveToISO,omitempty"`
	ExternallyRecorded bool                `json:"externallyRecorded,omitempty"`
}

func (f FinancialAgreementFact) Validate() error {
	for name, value := range map[string]string{"spaceID": f.SpaceID, "agreementID": f.AgreementID, "enrollmentID": f.EnrollmentID, "happeningID": f.HappeningID} {
		if err := validateFinancialCommitmentID(name, value); err != nil {
			return err
		}
	}
	if f.Revision < 1 || f.Revision > MaxJavaScriptSafeInteger {
		return fmt.Errorf("agreement revision is invalid")
	}
	if !validFinancialAgreementState(f.State) {
		return fmt.Errorf("agreement state is invalid")
	}
	if f.OwnSpaceSide != FinancialAgreementSidePayer && f.OwnSpaceSide != FinancialAgreementSideReceiver {
		return fmt.Errorf("ownSpaceSide is invalid")
	}
	if !validISODate(f.EffectiveFromISO) || !validOptionalISODate(f.EffectiveToISO) || f.EffectiveToISO != "" && f.EffectiveToISO < f.EffectiveFromISO {
		return fmt.Errorf("agreement effective dates are invalid")
	}
	if err := f.AcceptedPrice.Validate(); err != nil {
		return err
	}
	if err := validateAgreementParties(f.SpaceID, f.OwnSpaceSide, f.Payers, f.Receivers); err != nil {
		return err
	}
	if f.State == FinancialAgreementStateRecordedExternal && !f.ExternallyRecorded {
		return fmt.Errorf("recorded external agreement must disclose provenance")
	}
	return nil
}

func validFinancialAgreementState(state string) bool {
	switch state {
	case FinancialAgreementStateOffered, FinancialAgreementStateRecordedExternal, FinancialAgreementStateConfirmed, FinancialAgreementStateEnded, FinancialAgreementStateCancelled:
		return true
	}
	return false
}

func validateAgreementParties(spaceID, ownSide string, payers, receivers []FinancialPartyRef) error {
	if len(payers) == 0 || len(receivers) == 0 || len(payers) > MaxFinancialAgreementParties || len(receivers) > MaxFinancialAgreementParties {
		return fmt.Errorf("payer and receiver parties are required and bounded")
	}
	seen := map[string]struct{}{}
	ownPayer, ownReceiver := false, false
	for _, side := range []struct {
		name   string
		values []FinancialPartyRef
	}{{"payer", payers}, {"receiver", receivers}} {
		for _, party := range side.values {
			if err := party.Validate(); err != nil {
				return fmt.Errorf("%s: %w", side.name, err)
			}
			key := party.Kind + "\x00" + party.SpaceID + "\x00" + party.ContactID
			if _, ok := seen[key]; ok {
				return fmt.Errorf("financial party appears more than once")
			}
			seen[key] = struct{}{}
			if party.Kind == FinancialPartyKindSpace && party.SpaceID == spaceID {
				if side.name == "payer" {
					ownPayer = true
				} else {
					ownReceiver = true
				}
			}
		}
	}
	if ownPayer && ownReceiver {
		return fmt.Errorf("the reporting Space cannot be on both sides of one fee")
	}
	if ownSide == FinancialAgreementSidePayer && !ownPayer || ownSide == FinancialAgreementSideReceiver && !ownReceiver {
		return fmt.Errorf("ownSpaceSide does not match an explicit Space party")
	}
	return nil
}

type FinancialAgreementQuery struct {
	SpaceID      string `json:"spaceID"`
	FromMonthISO string `json:"fromMonthISO"`
	Months       int    `json:"months"`
}

func (q FinancialAgreementQuery) Validate() error {
	if err := validateFinancialCommitmentID("spaceID", q.SpaceID); err != nil {
		return err
	}
	if !validMonthISO(q.FromMonthISO) || q.Months < 1 || q.Months > MaxFinancialCommitmentMonths {
		return fmt.Errorf("agreement query window is invalid")
	}
	from, _ := time.Parse("2006-01", q.FromMonthISO)
	if from.AddDate(0, q.Months-1, 0).Year() > 9999 {
		return fmt.Errorf("agreement query exceeds year 9999")
	}
	return nil
}

type FinancialAgreementResult struct {
	Agreements []FinancialAgreementFact `json:"agreements"`
}

func (r FinancialAgreementResult) Validate() error {
	if r.Agreements == nil {
		return fmt.Errorf("agreements must be a non-nil array")
	}
	seen := map[string]struct{}{}
	for i, agreement := range r.Agreements {
		if err := agreement.Validate(); err != nil {
			return fmt.Errorf("agreements[%d]: %w", i, err)
		}
		if _, ok := seen[agreement.AgreementID]; ok {
			return fmt.Errorf("duplicate agreementID")
		}
		seen[agreement.AgreementID] = struct{}{}
	}
	return nil
}

type FinancialAgreementSource interface {
	ListFinancialAgreements(context.Context, string, FinancialAgreementQuery) (FinancialAgreementResult, error)
}

// SaveFinancialAgreementRequest selects one catalog price for one enrollment.
// Accepted money/term values are resolved by Calendarius and are never supplied
// by the client. Confirmed is intentionally absent from this command: recording
// an external deal cannot fabricate the counterparty's digital acceptance.
type SaveFinancialAgreementRequest struct {
	SpaceID               string              `json:"spaceID"`
	AgreementID           string              `json:"agreementID"`
	EnrollmentID          string              `json:"enrollmentID"`
	HappeningID           string              `json:"happeningID"`
	PriceID               string              `json:"priceID"`
	ExpectedPriceRevision int64               `json:"expectedPriceRevision"`
	OperationID           string              `json:"operationID"`
	ExpectedRevision      int64               `json:"expectedRevision"`
	State                 string              `json:"state"`
	OwnSpaceSide          string              `json:"ownSpaceSide"`
	Payers                []FinancialPartyRef `json:"payers"`
	Receivers             []FinancialPartyRef `json:"receivers"`
	EffectiveFromISO      string              `json:"effectiveFromISO"`
	EffectiveToISO        string              `json:"effectiveToISO,omitempty"`
	Reason                string              `json:"reason,omitempty"`
}

func (r SaveFinancialAgreementRequest) Validate() error {
	for name, value := range map[string]string{"spaceID": r.SpaceID, "agreementID": r.AgreementID, "enrollmentID": r.EnrollmentID, "happeningID": r.HappeningID, "priceID": r.PriceID, "operationID": r.OperationID} {
		if err := validateFinancialCommitmentID(name, value); err != nil {
			return err
		}
	}
	if r.ExpectedRevision < 0 || r.ExpectedRevision > MaxJavaScriptSafeInteger || r.ExpectedPriceRevision < 1 || r.ExpectedPriceRevision > MaxJavaScriptSafeInteger {
		return fmt.Errorf("expected revisions are outside supported bounds")
	}
	if r.State != FinancialAgreementStateOffered && r.State != FinancialAgreementStateRecordedExternal {
		return fmt.Errorf("save agreement state must be offered or recorded_external")
	}
	if r.OwnSpaceSide != FinancialAgreementSidePayer && r.OwnSpaceSide != FinancialAgreementSideReceiver {
		return fmt.Errorf("ownSpaceSide is invalid")
	}
	if !validISODate(r.EffectiveFromISO) || !validOptionalISODate(r.EffectiveToISO) || r.EffectiveToISO != "" && r.EffectiveToISO < r.EffectiveFromISO {
		return fmt.Errorf("agreement effective dates are invalid")
	}
	if len(r.Reason) > 500 {
		return fmt.Errorf("reason exceeds 500 characters")
	}
	return validateAgreementParties(r.SpaceID, r.OwnSpaceSide, r.Payers, r.Receivers)
}

type SaveFinancialAgreementResponse struct {
	AgreementID string `json:"agreementID"`
	Revision    int64  `json:"revision"`
	UpdatedAt   string `json:"updatedAt"`
}
