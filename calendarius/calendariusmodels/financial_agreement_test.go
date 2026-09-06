package calendariusmodels

import "testing"

func validFinancialAgreement() FinancialAgreementFact {
	return FinancialAgreementFact{
		OwnerSpaceID: "school", ReportingSpaceID: "family", AgreementID: "music-child-a", EnrollmentID: "child-a", HappeningID: "music",
		Revision: 2, TermsRevision: 1, State: FinancialAgreementStateRecordedExternal, ReportingSpaceSide: FinancialAgreementSidePayer,
		ReportingAmountMinor: 4000, ExternallyRecorded: true,
		Payers:              []FinancialPartyShare{{Party: FinancialPartyRef{Kind: FinancialPartyKindSpace, SpaceID: "family"}, AmountMinor: 4000}},
		Receivers:           []FinancialPartyShare{{Party: FinancialPartyRef{Kind: FinancialPartyKindSpace, SpaceID: "school"}, AmountMinor: 4000}},
		ContactAttributions: []FinancialAttribution{{ID: "child-a", AmountMinor: 4000}}, AssetAttributions: []FinancialAttribution{},
		AcceptedPrice:    AcceptedPriceTerms{PriceID: "monthly", PriceRevision: 0, AmountMinor: 4000, Currency: "EUR", Quantity: 1, Term: FinancialCommitmentTerm{Unit: "month", Length: 1}},
		EffectiveFromISO: "2026-09-01", Confirmations: []FinancialConfirmation{},
	}
}

func TestFinancialAgreementValidatesOwnerScopedReportingView(t *testing.T) {
	value := validFinancialAgreement()
	if err := value.Validate(); err != nil {
		t.Fatal(err)
	}
	value.Receivers = append(value.Receivers, FinancialPartyShare{Party: FinancialPartyRef{Kind: FinancialPartyKindSpace, SpaceID: "family"}, AmountMinor: 1})
	if err := value.Validate(); err == nil {
		t.Fatal("same reporting Space on both sides was accepted")
	}
}

func TestFinancialAgreementRequiresExactIndependentSideTotals(t *testing.T) {
	value := validFinancialAgreement()
	value.Receivers = []FinancialPartyShare{
		{Party: FinancialPartyRef{Kind: FinancialPartyKindContact, SpaceID: "school", ContactID: "teacher"}, AmountMinor: 3000},
		{Party: FinancialPartyRef{Kind: FinancialPartyKindSpace, SpaceID: "school"}, AmountMinor: 1000},
	}
	if err := value.Validate(); err != nil {
		t.Fatal(err)
	}
	value.Receivers[0].AmountMinor = 2999
	if err := value.Validate(); err == nil {
		t.Fatal("non-conserved receiver allocations were accepted")
	}
}

func TestFinancialAgreementAcceptedTermsAreStrictAndInitialRevisionIsValid(t *testing.T) {
	value := validFinancialAgreement()
	if err := value.AcceptedPrice.Validate(); err != nil {
		t.Fatal(err)
	}
	for _, mutate := range []func(*AcceptedPriceTerms){
		func(v *AcceptedPriceTerms) { v.Currency = "eur" },
		func(v *AcceptedPriceTerms) { v.Term.Unit = "fortnight" },
		func(v *AcceptedPriceTerms) { v.Term.Unit, v.Term.Length = "single", 2 },
		func(v *AcceptedPriceTerms) { v.AmountMinor, v.Quantity = MaxJavaScriptSafeInteger, 2 },
	} {
		invalid := value.AcceptedPrice
		mutate(&invalid)
		if err := invalid.Validate(); err == nil {
			t.Fatal("invalid accepted price was accepted")
		}
	}
}

func TestFinancialConfirmationSurvivesLaterRecordRevision(t *testing.T) {
	value := validFinancialAgreement()
	value.Revision = 3
	value.Confirmations = []FinancialConfirmation{{Party: value.Payers[0].Party, TermsRevision: 1, ActorUserID: "payer-user", ConfirmedAt: "2026-09-02T10:00:00Z"}}
	if err := value.Validate(); err != nil {
		t.Fatal(err)
	}
	value.Confirmations[0].TermsRevision = 2
	if err := value.Validate(); err == nil {
		t.Fatal("confirmation for another terms revision was accepted")
	}
}

func TestConfirmedAgreementRequiresActiveEvidenceFromEveryParty(t *testing.T) {
	value := validFinancialAgreement()
	value.State = FinancialAgreementStateConfirmed
	value.ExternallyRecorded = false
	value.Confirmations = []FinancialConfirmation{
		{Party: value.Payers[0].Party, TermsRevision: 1, ActorUserID: "payer-user", ConfirmedAt: "2026-09-02T10:00:00Z"},
	}
	if err := value.Validate(); err == nil {
		t.Fatal("confirmed agreement without receiver evidence was accepted")
	}
	value.Confirmations = append(value.Confirmations, FinancialConfirmation{Party: value.Receivers[0].Party, TermsRevision: 1, ActorUserID: "receiver-user", ConfirmedAt: "2026-09-02T11:00:00Z"})
	if err := value.Validate(); err != nil {
		t.Fatal(err)
	}
	value.Confirmations[1].RevokedAt = "2026-09-03T11:00:00Z"
	if err := value.Validate(); err == nil {
		t.Fatal("confirmed agreement with revoked receiver evidence was accepted")
	}
}

func TestFinancialAgreementQueryAndPageAreBoundedAndOwnerQualified(t *testing.T) {
	query := FinancialAgreementQuery{ReportingSpaceID: "family", FromMonthISO: "2026-09", Months: 12, PageSize: 64}
	if err := query.Validate(); err != nil {
		t.Fatal(err)
	}
	query.Cursor = " padded "
	if err := query.Validate(); err == nil {
		t.Fatal("invalid cursor accepted")
	}
	value := validFinancialAgreement()
	otherOwner := value
	otherOwner.OwnerSpaceID = "school-2"
	if err := (FinancialAgreementResult{Agreements: []FinancialAgreementFact{value, otherOwner}, SnapshotConsistent: true}).Validate(); err != nil {
		t.Fatal(err)
	}
	if err := (FinancialAgreementResult{Agreements: []FinancialAgreementFact{value, value}, SnapshotConsistent: true}).Validate(); err == nil {
		t.Fatal("duplicate owner identity accepted")
	}
	if err := (FinancialAgreementResult{Agreements: []FinancialAgreementFact{}, HasMore: true, NextCursor: "next", IncompleteReason: "source_collection_mutable_between_pages"}).Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestSaveFinancialAgreementUsesWeightsNotClientMoney(t *testing.T) {
	request := SaveFinancialAgreementRequest{
		OwnerSpaceID: "school", ReportingSpaceID: "family", AgreementID: "deal-1", EnrollmentID: "child-a", HappeningID: "music", PriceID: "monthly",
		Quantity: 1, OperationID: "op-1", State: FinancialAgreementStateRecordedExternal, ReportingSpaceSide: FinancialAgreementSidePayer,
		Payers:              []FinancialPartyWeight{{Party: FinancialPartyRef{Kind: FinancialPartyKindSpace, SpaceID: "family"}, Shares: 1}},
		Receivers:           []FinancialPartyWeight{{Party: FinancialPartyRef{Kind: FinancialPartyKindContact, SpaceID: "school", ContactID: "teacher"}, Shares: 3}, {Party: FinancialPartyRef{Kind: FinancialPartyKindSpace, SpaceID: "school"}, Shares: 1}},
		ContactAttributions: []FinancialAttributionWeight{{ID: "child-a", Shares: 1}}, EffectiveFromISO: "2026-09-01",
	}
	if err := request.Validate(); err != nil {
		t.Fatal(err)
	}
	request.Payers[0].Shares = 0
	if err := request.Validate(); err == nil {
		t.Fatal("zero allocation weight accepted")
	}
}

func TestConfirmationAndRevocationCommandsAreRevisionBounded(t *testing.T) {
	request := ConfirmFinancialAgreementRequest{ReportingSpaceID: "family", OwnerSpaceID: "school", AgreementID: "deal-1", ExpectedRevision: 2, TermsRevision: 1, OperationID: "op-2", Party: FinancialPartyRef{Kind: FinancialPartyKindSpace, SpaceID: "family"}}
	if err := request.Validate(); err != nil {
		t.Fatal(err)
	}
	revoke := RevokeFinancialAgreementGrantRequest{ReportingSpaceID: request.ReportingSpaceID, OwnerSpaceID: request.OwnerSpaceID, AgreementID: request.AgreementID, ExpectedRevision: request.ExpectedRevision, TermsRevision: request.TermsRevision, OperationID: "op-3", Party: request.Party}
	if err := revoke.Validate(); err != nil {
		t.Fatal(err)
	}
	revoke.Reason = string(make([]byte, 501))
	if err := revoke.Validate(); err == nil {
		t.Fatal("overlong revocation reason accepted")
	}
}
