package calendariusmodels

import "testing"

func validFinancialAgreement() FinancialAgreementFact {
	return FinancialAgreementFact{
		SpaceID: "family", AgreementID: "music-child-a", EnrollmentID: "child-a", HappeningID: "music", Revision: 1,
		State: FinancialAgreementStateRecordedExternal, OwnSpaceSide: FinancialAgreementSidePayer, ExternallyRecorded: true,
		Payers:           []FinancialPartyRef{{Kind: FinancialPartyKindSpace, SpaceID: "family"}},
		Receivers:        []FinancialPartyRef{{Kind: FinancialPartyKindContact, SpaceID: "family", ContactID: "teacher"}},
		AcceptedPrice:    AcceptedPriceTerms{PriceID: "monthly", PriceRevision: 2, AmountMinor: 4000, Currency: "EUR", Quantity: 1, Term: FinancialCommitmentTerm{Unit: "month", Length: 1}},
		EffectiveFromISO: "2026-09-01",
	}
}

func TestFinancialAgreementRequiresOneExplicitReportingSide(t *testing.T) {
	value := validFinancialAgreement()
	if err := value.Validate(); err != nil {
		t.Fatal(err)
	}
	value.Receivers = append(value.Receivers, FinancialPartyRef{Kind: FinancialPartyKindSpace, SpaceID: value.SpaceID})
	if err := value.Validate(); err == nil {
		t.Fatal("same reporting Space on payer and receiver sides was accepted")
	}
}

func TestFinancialAgreementDoesNotTreatContactStorageAsSpaceInterest(t *testing.T) {
	value := validFinancialAgreement()
	value.Payers = []FinancialPartyRef{{Kind: FinancialPartyKindContact, SpaceID: value.SpaceID, ContactID: "student"}}
	if err := value.Validate(); err == nil {
		t.Fatal("contact storage address was treated as reporting-Space economic interest")
	}
}

func TestFinancialAgreementRecordedExternalDisclosesProvenance(t *testing.T) {
	value := validFinancialAgreement()
	value.ExternallyRecorded = false
	if err := value.Validate(); err == nil {
		t.Fatal("externally recorded agreement omitted provenance")
	}
}

func TestFinancialAgreementResultRequiresUniqueNonNilFacts(t *testing.T) {
	value := validFinancialAgreement()
	if err := (FinancialAgreementResult{}).Validate(); err == nil {
		t.Fatal("nil agreement array accepted")
	}
	if err := (FinancialAgreementResult{Agreements: []FinancialAgreementFact{value, value}}).Validate(); err == nil {
		t.Fatal("duplicate agreement accepted")
	}
}

func TestFinancialAgreementQueryBoundsWindow(t *testing.T) {
	if err := (FinancialAgreementQuery{SpaceID: "family", FromMonthISO: "2026-09", Months: 12}).Validate(); err != nil {
		t.Fatal(err)
	}
	if err := (FinancialAgreementQuery{SpaceID: "family", FromMonthISO: "9999-12", Months: 2}).Validate(); err == nil {
		t.Fatal("year overflow accepted")
	}
}
