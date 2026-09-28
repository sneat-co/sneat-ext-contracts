package calendariusmodels

import (
	"fmt"
	"strings"
	"testing"
)

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

func TestFinancialAgreementReportingAmountEqualsExplicitSpaceShare(t *testing.T) {
	value := validFinancialAgreement()
	value.Payers = []FinancialPartyShare{
		{Party: FinancialPartyRef{Kind: FinancialPartyKindSpace, SpaceID: "family"}, AmountMinor: 2000},
		{Party: FinancialPartyRef{Kind: FinancialPartyKindContact, SpaceID: "family", ContactID: "sponsor"}, AmountMinor: 2000},
	}
	if err := value.Validate(); err == nil {
		t.Fatal("wrong but in-range reporting amount accepted")
	}
	value.ReportingAmountMinor = 2000
	value.ContactAttributions = []FinancialAttribution{{ID: "child-a", AmountMinor: 2000}}
	if err := value.Validate(); err != nil {
		t.Fatal(err)
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

func TestFinancialAgreementRejectsNonpartyAndBackdatedRevocation(t *testing.T) {
	value := validFinancialAgreement()
	value.Confirmations = []FinancialConfirmation{{Party: FinancialPartyRef{Kind: FinancialPartyKindSpace, SpaceID: "stranger"}, TermsRevision: 1, ActorUserID: "user", ConfirmedAt: "2026-09-02T10:00:00Z"}}
	if err := value.Validate(); err == nil {
		t.Fatal("nonparty confirmation accepted")
	}
	value.Confirmations[0].Party = value.Payers[0].Party
	value.Confirmations[0].RevokedAt = "2026-09-01T10:00:00Z"
	if err := value.Validate(); err == nil {
		t.Fatal("revocation before confirmation accepted")
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

func TestFinancialAgreement_ComprehensiveCoverage(t *testing.T) {
	// 1. FinancialPartyRef
	ref := FinancialPartyRef{Kind: "invalid", SpaceID: "s1"}
	if err := ref.Validate(); err == nil {
		t.Fatal("expected error for invalid party kind")
	}
	ref = FinancialPartyRef{Kind: FinancialPartyKindSpace, SpaceID: " bad "}
	if err := ref.Validate(); err == nil {
		t.Fatal("expected error for untrimmed spaceID")
	}
	ref = FinancialPartyRef{Kind: FinancialPartyKindSpace, SpaceID: "s1", ContactID: "c1"}
	if err := ref.Validate(); err == nil {
		t.Fatal("expected error for Space party naming a contact")
	}
	ref = FinancialPartyRef{Kind: FinancialPartyKindContact, SpaceID: "s1", ContactID: " bad "}
	if err := ref.Validate(); err == nil {
		t.Fatal("expected error for invalid contactID")
	}
	ref = FinancialPartyRef{Kind: FinancialPartyKindContact, SpaceID: "s1", ContactID: "c1"}
	if err := ref.Validate(); err != nil {
		t.Fatalf("valid contact party failed: %v", err)
	}

	// 2. AcceptedPriceTerms
	price := AcceptedPriceTerms{PriceID: " bad ", PriceRevision: 0, AmountMinor: 100, Currency: "EUR", Quantity: 1, Term: FinancialCommitmentTerm{Unit: "month", Length: 1}}
	if err := price.Validate(); err == nil {
		t.Fatal("expected error for invalid priceID")
	}
	price = AcceptedPriceTerms{PriceID: "p1", PriceRevision: -1, AmountMinor: 100, Currency: "EUR", Quantity: 1, Term: FinancialCommitmentTerm{Unit: "month", Length: 1}}
	if err := price.Validate(); err == nil {
		t.Fatal("expected error for negative priceRevision")
	}
	price = AcceptedPriceTerms{PriceID: "p1", PriceRevision: 0, AmountMinor: 100, Currency: "EUR", Quantity: 0, Term: FinancialCommitmentTerm{Unit: "month", Length: 1}}
	if err := price.Validate(); err == nil {
		t.Fatal("expected error for quantity < 1")
	}
	price = AcceptedPriceTerms{PriceID: "p1", PriceRevision: 0, AmountMinor: 100, Currency: "EUR", Quantity: 1, Term: FinancialCommitmentTerm{Unit: "quarter", Length: 0}}
	if err := price.Validate(); err == nil {
		t.Fatal("expected error for term length < 1")
	}
	price = AcceptedPriceTerms{PriceID: "p1", PriceRevision: 0, AmountMinor: 100, Currency: "EUR", Quantity: 1, Term: FinancialCommitmentTerm{Unit: "year", Length: 1}}
	if err := price.Validate(); err != nil {
		t.Fatalf("quarter/year term failed: %v", err)
	}

	// 3. FinancialConfirmation
	conf := FinancialConfirmation{Party: FinancialPartyRef{Kind: "bad"}, TermsRevision: 1, ActorUserID: "u1", ConfirmedAt: "2026-09-01T10:00:00Z"}
	if err := conf.Validate(1); err == nil {
		t.Fatal("expected error for invalid party in confirmation")
	}
	conf = FinancialConfirmation{Party: FinancialPartyRef{Kind: FinancialPartyKindSpace, SpaceID: "s1"}, TermsRevision: 1, ActorUserID: " ", ConfirmedAt: "2026-09-01T10:00:00Z"}
	if err := conf.Validate(1); err == nil {
		t.Fatal("expected error for blank actorUserID")
	}
	conf = FinancialConfirmation{Party: FinancialPartyRef{Kind: FinancialPartyKindSpace, SpaceID: "s1"}, TermsRevision: 1, ActorUserID: "u1", ConfirmedAt: "bad-time"}
	if err := conf.Validate(1); err == nil {
		t.Fatal("expected error for invalid confirmedAt")
	}
	conf = FinancialConfirmation{Party: FinancialPartyRef{Kind: FinancialPartyKindSpace, SpaceID: "s1"}, TermsRevision: 1, ActorUserID: "u1", ConfirmedAt: "2026-09-01T10:00:00Z", RevokedAt: "bad-time"}
	if err := conf.Validate(1); err == nil {
		t.Fatal("expected error for invalid revokedAt")
	}

	// 4. FinancialAgreementFact
	fact := validFinancialAgreement()
	// bad ID
	badFact := fact
	badFact.OwnerSpaceID = " bad "
	if err := badFact.Validate(); err == nil {
		t.Fatal("expected error for invalid ownerSpaceID")
	}
	// bad revision
	badFact = fact
	badFact.Revision = 0
	if err := badFact.Validate(); err == nil {
		t.Fatal("expected error for revision < 1")
	}
	badFact = fact
	badFact.TermsRevision = 0
	if err := badFact.Validate(); err == nil {
		t.Fatal("expected error for termsRevision < 1")
	}
	badFact = fact
	badFact.TermsRevision = 3
	badFact.Revision = 2
	if err := badFact.Validate(); err == nil {
		t.Fatal("expected error for termsRevision > revision")
	}
	// bad state
	badFact = fact
	badFact.State = "invalid_state"
	if err := badFact.Validate(); err == nil {
		t.Fatal("expected error for invalid state")
	}
	// bad reportingSpaceSide
	badFact = fact
	badFact.ReportingSpaceSide = "neither"
	if err := badFact.Validate(); err == nil {
		t.Fatal("expected error for invalid reportingSpaceSide")
	}
	// bad dates
	badFact = fact
	badFact.EffectiveFromISO = "bad-date"
	if err := badFact.Validate(); err == nil {
		t.Fatal("expected error for bad EffectiveFromISO")
	}
	badFact = fact
	badFact.EffectiveToISO = "2026-08-01" // before From 2026-09-01
	if err := badFact.Validate(); err == nil {
		t.Fatal("expected error for to < from")
	}
	// bad acceptedPrice
	badFact = fact
	badFact.AcceptedPrice.PriceID = " bad "
	if err := badFact.Validate(); err == nil {
		t.Fatal("expected error for bad accepted price")
	}
	// bad reportingAmountMinor
	badFact = fact
	badFact.ReportingAmountMinor = 0
	if err := badFact.Validate(); err == nil {
		t.Fatal("expected error for reportingAmountMinor < 1")
	}
	badFact = fact
	badFact.ReportingAmountMinor = 5000 // total is 4000
	if err := badFact.Validate(); err == nil {
		t.Fatal("expected error for reportingAmountMinor > total")
	}
	// bad payers/receivers shares
	badFact = fact
	badFact.Payers = nil
	if err := badFact.Validate(); err == nil {
		t.Fatal("expected error for nil payers")
	}
	// ReportingSpaceSide == receiver
	recFact := validFinancialAgreement()
	recFact.ReportingSpaceID = "school"
	recFact.ReportingSpaceSide = FinancialAgreementSideReceiver
	recFact.ReportingAmountMinor = 4000
	recFact.Receivers = []FinancialPartyShare{{Party: FinancialPartyRef{Kind: FinancialPartyKindSpace, SpaceID: "school"}, AmountMinor: 4000}}
	recFact.ContactAttributions = []FinancialAttribution{{ID: "teacher", AmountMinor: 4000}}
	if err := recFact.Validate(); err != nil {
		t.Fatalf("receiver fact failed: %v", err)
	}
	// ReportingSpaceSide == receiver but reporting amount doesn't match space share
	recMismatch := recFact
	recMismatch.Receivers = []FinancialPartyShare{
		{Party: FinancialPartyRef{Kind: FinancialPartyKindSpace, SpaceID: "school"}, AmountMinor: 2000},
		{Party: FinancialPartyRef{Kind: FinancialPartyKindContact, SpaceID: "school", ContactID: "teacher"}, AmountMinor: 2000},
	}
	recMismatch.ContactAttributions = []FinancialAttribution{{ID: "teacher", AmountMinor: 4000}}
	if err := recMismatch.Validate(); err == nil {
		t.Fatal("expected error for reporting amount not matching receiver space share")
	}
	// bad contactAttributions / assetAttributions
	badFact = fact
	badFact.ContactAttributions = []FinancialAttribution{{ID: " bad ", AmountMinor: 4000}}
	if err := badFact.Validate(); err == nil {
		t.Fatal("expected error for bad contact attribution")
	}
	badFact = fact
	badFact.AssetAttributions = []FinancialAttribution{{ID: " bad ", AmountMinor: 4000}}
	if err := badFact.Validate(); err == nil {
		t.Fatal("expected error for bad asset attribution")
	}
	// nil confirmations
	badFact = fact
	badFact.Confirmations = nil
	if err := badFact.Validate(); err == nil {
		t.Fatal("expected error for nil confirmations")
	}
	// confirmations exceed agreement parties
	badFact = fact
	badFact.Confirmations = make([]FinancialConfirmation, len(fact.Payers)+len(fact.Receivers)+1)
	if err := badFact.Validate(); err == nil {
		t.Fatal("expected error for confirmations exceeding parties")
	}
	// duplicate confirmation party
	badFact = fact
	badFact.Confirmations = []FinancialConfirmation{
		{Party: fact.Payers[0].Party, TermsRevision: 1, ActorUserID: "u1", ConfirmedAt: "2026-09-02T10:00:00Z"},
		{Party: fact.Payers[0].Party, TermsRevision: 1, ActorUserID: "u2", ConfirmedAt: "2026-09-02T11:00:00Z"},
	}
	if err := badFact.Validate(); err == nil {
		t.Fatal("expected error for duplicate confirmation party")
	}
	// recorded external without ExternallyRecorded flag
	badFact = fact
	badFact.State = FinancialAgreementStateRecordedExternal
	badFact.ExternallyRecorded = false
	if err := badFact.Validate(); err == nil {
		t.Fatal("expected error for recorded external without ExternallyRecorded flag")
	}
	// offered, ended, cancelled states
	for _, st := range []string{FinancialAgreementStateOffered, FinancialAgreementStateEnded, FinancialAgreementStateCancelled} {
		sFact := fact
		sFact.State = st
		sFact.ExternallyRecorded = false
		if err := sFact.Validate(); err != nil {
			t.Fatalf("state %s failed: %v", st, err)
		}
	}

	// 5. validateAgreementShares
	// duplicate party in shares
	badShares := fact
	badShares.Payers = []FinancialPartyShare{
		{Party: fact.Payers[0].Party, AmountMinor: 2000},
		{Party: fact.Payers[0].Party, AmountMinor: 2000},
	}
	if err := badShares.Validate(); err == nil {
		t.Fatal("expected error for duplicate party in shares")
	}
	// invalid party in shares
	badParty := fact
	badParty.Payers = []FinancialPartyShare{
		{Party: FinancialPartyRef{Kind: "bad"}, AmountMinor: 4000},
	}
	if err := badParty.Validate(); err == nil {
		t.Fatal("expected error for invalid party in shares")
	}
	// amountMinor < 1
	badZero := fact
	badZero.Payers = []FinancialPartyShare{
		{Party: fact.Payers[0].Party, AmountMinor: 0},
	}
	if err := badZero.Validate(); err == nil {
		t.Fatal("expected error for zero share amount")
	}

	// 6. validateAttributions
	tooManyAttr := make([]FinancialAttribution, MaxFinancialAgreementParties+1)
	for i := range tooManyAttr {
		tooManyAttr[i] = FinancialAttribution{ID: fmt.Sprintf("a%d", i), AmountMinor: 1}
	}
	if err := validateAttributions(4000, "attr", tooManyAttr); err == nil {
		t.Fatal("expected error for too many attributions")
	}
	dupAttr := []FinancialAttribution{
		{ID: "a1", AmountMinor: 2000},
		{ID: "a1", AmountMinor: 2000},
	}
	if err := validateAttributions(4000, "attr", dupAttr); err == nil {
		t.Fatal("expected error for duplicate attribution ID")
	}
	badAmountAttr := []FinancialAttribution{
		{ID: "a1", AmountMinor: 5000},
	}
	if err := validateAttributions(4000, "attr", badAmountAttr); err == nil {
		t.Fatal("expected error for attribution amount exceeding total")
	}
	underAttr := []FinancialAttribution{
		{ID: "a1", AmountMinor: 2000},
	}
	if err := validateAttributions(4000, "attr", underAttr); err == nil {
		t.Fatal("expected error for attribution amount not conserving total")
	}

	// 7. validatePartyWeights
	tooManyPartyW := make([]FinancialPartyWeight, MaxFinancialAgreementParties+1)
	for i := range tooManyPartyW {
		tooManyPartyW[i] = FinancialPartyWeight{Party: FinancialPartyRef{Kind: FinancialPartyKindSpace, SpaceID: fmt.Sprintf("s%d", i)}, Shares: 1}
	}
	if err := validatePartyWeights("payers", tooManyPartyW); err == nil {
		t.Fatal("expected error for too many party weights")
	}
	if err := validatePartyWeights("payers", nil); err == nil {
		t.Fatal("expected error for empty party weights")
	}
	badPartyW := []FinancialPartyWeight{
		{Party: FinancialPartyRef{Kind: "bad"}, Shares: 1},
	}
	if err := validatePartyWeights("payers", badPartyW); err == nil {
		t.Fatal("expected error for invalid party in weights")
	}
	dupPartyW := []FinancialPartyWeight{
		{Party: FinancialPartyRef{Kind: FinancialPartyKindSpace, SpaceID: "s1"}, Shares: 1},
		{Party: FinancialPartyRef{Kind: FinancialPartyKindSpace, SpaceID: "s1"}, Shares: 1},
	}
	if err := validatePartyWeights("payers", dupPartyW); err == nil {
		t.Fatal("expected error for duplicate party in weights")
	}
	badSharesW := []FinancialPartyWeight{
		{Party: FinancialPartyRef{Kind: FinancialPartyKindSpace, SpaceID: "s1"}, Shares: 0},
	}
	if err := validatePartyWeights("payers", badSharesW); err == nil {
		t.Fatal("expected error for 0 shares in weights")
	}

	// 8. validateAttributionWeights
	tooManyAttrW := make([]FinancialAttributionWeight, MaxFinancialAgreementParties+1)
	for i := range tooManyAttrW {
		tooManyAttrW[i] = FinancialAttributionWeight{ID: fmt.Sprintf("a%d", i), Shares: 1}
	}
	if err := validateAttributionWeights("attr", tooManyAttrW); err == nil {
		t.Fatal("expected error for too many attribution weights")
	}
	badIDAttrW := []FinancialAttributionWeight{
		{ID: " bad ", Shares: 1},
	}
	if err := validateAttributionWeights("attr", badIDAttrW); err == nil {
		t.Fatal("expected error for bad ID in attribution weights")
	}
	dupAttrW := []FinancialAttributionWeight{
		{ID: "a1", Shares: 1},
		{ID: "a1", Shares: 1},
	}
	if err := validateAttributionWeights("attr", dupAttrW); err == nil {
		t.Fatal("expected error for duplicate ID in attribution weights")
	}
	badSharesAttrW := []FinancialAttributionWeight{
		{ID: "a1", Shares: 0},
	}
	if err := validateAttributionWeights("attr", badSharesAttrW); err == nil {
		t.Fatal("expected error for zero shares in attribution weights")
	}

	// 9. FinancialAgreementListRequest
	listReq := FinancialAgreementListRequest{OwnerSpaceID: "school", HappeningID: "music", PageSize: 50}
	if err := listReq.Validate(); err != nil {
		t.Fatalf("listReq failed: %v", err)
	}
	bListReq := listReq
	bListReq.OwnerSpaceID = " bad "
	if err := bListReq.Validate(); err == nil {
		t.Fatal("expected error for bad ownerSpaceID in list request")
	}
	bListReq = listReq
	bListReq.HappeningID = " bad "
	if err := bListReq.Validate(); err == nil {
		t.Fatal("expected error for bad happeningID in list request")
	}
	bListReq = listReq
	bListReq.PageSize = 0
	if err := bListReq.Validate(); err == nil {
		t.Fatal("expected error for pageSize 0 in list request")
	}
	bListReq = listReq
	bListReq.Cursor = " bad "
	if err := bListReq.Validate(); err == nil {
		t.Fatal("expected error for bad cursor in list request")
	}

	// 10. FinancialAgreementQuery
	query := FinancialAgreementQuery{ReportingSpaceID: "family", FromMonthISO: "2026-09", Months: 12, PageSize: 64}
	if err := query.Validate(); err != nil {
		t.Fatalf("query failed: %v", err)
	}
	bQuery := query
	bQuery.ReportingSpaceID = " bad "
	if err := bQuery.Validate(); err == nil {
		t.Fatal("expected error for bad reportingSpaceID in query")
	}
	bQuery = query
	bQuery.FromMonthISO = "bad-month"
	if err := bQuery.Validate(); err == nil {
		t.Fatal("expected error for bad fromMonthISO in query")
	}
	bQuery = query
	bQuery.Months = 0
	if err := bQuery.Validate(); err == nil {
		t.Fatal("expected error for months 0 in query")
	}
	bQuery = query
	bQuery.PageSize = 0
	if err := bQuery.Validate(); err == nil {
		t.Fatal("expected error for pageSize 0 in query")
	}
	bQuery = query
	bQuery.FromMonthISO = "9999-10"
	bQuery.Months = 12 // overflow 9999
	if err := bQuery.Validate(); err == nil {
		t.Fatal("expected error for query year > 9999")
	}

	// 11. validateAgreementCursor
	for _, badCursor := range []string{strings.Repeat("c", MaxFinancialAgreementCursorBytes+1), "\xff", " leading", "trailing ", "line\nbreak"} {
		if err := validateAgreementCursor(badCursor); err == nil {
			t.Fatalf("expected error for bad cursor: %q", badCursor)
		}
	}

	// 12. FinancialAgreementResult
	res := FinancialAgreementResult{
		Agreements:         []FinancialAgreementFact{validFinancialAgreement()},
		SnapshotConsistent: true,
	}
	if err := res.Validate(); err != nil {
		t.Fatalf("res failed: %v", err)
	}
	bRes := res
	bRes.Agreements = nil
	if err := bRes.Validate(); err == nil {
		t.Fatal("expected error for nil agreements in result")
	}
	bRes = res
	bRes.HasMore = true
	bRes.NextCursor = "" // mismatch
	if err := bRes.Validate(); err == nil {
		t.Fatal("expected error for hasMore with empty nextCursor")
	}
	bRes = res
	bRes.HasMore = true
	bRes.NextCursor = " bad "
	if err := bRes.Validate(); err == nil {
		t.Fatal("expected error for bad nextCursor")
	}
	bRes = res
	bRes.SnapshotConsistent = false
	bRes.IncompleteReason = "invalid_reason"
	if err := bRes.Validate(); err == nil {
		t.Fatal("expected error for invalid incompleteReason")
	}
	bRes = res
	bRes.SnapshotConsistent = true
	bRes.IncompleteReason = "query_limit" // consistency mismatch
	if err := bRes.Validate(); err == nil {
		t.Fatal("expected error for snapshotConsistent with incompleteReason")
	}
	bRes = res
	bRes.SnapshotConsistent = true
	bRes.HasMore = true
	bRes.NextCursor = "cursor1"
	if err := bRes.Validate(); err == nil {
		t.Fatal("expected error for multi-page claiming snapshotConsistent")
	}
	bRes = res
	bRes.Agreements = []FinancialAgreementFact{{OwnerSpaceID: " bad "}}
	if err := bRes.Validate(); err == nil {
		t.Fatal("expected error for invalid agreement in result")
	}

	// 13. SaveFinancialAgreementRequest
	validSaveReq := SaveFinancialAgreementRequest{
		OwnerSpaceID: "school", ReportingSpaceID: "family", AgreementID: "deal-1", EnrollmentID: "child-a", HappeningID: "music", PriceID: "monthly",
		ExpectedPriceRevision: 0, Quantity: 1, OperationID: "op-1", ExpectedRevision: 1, State: FinancialAgreementStateOffered, ReportingSpaceSide: FinancialAgreementSidePayer,
		Payers:              []FinancialPartyWeight{{Party: FinancialPartyRef{Kind: FinancialPartyKindSpace, SpaceID: "family"}, Shares: 1}},
		Receivers:           []FinancialPartyWeight{{Party: FinancialPartyRef{Kind: FinancialPartyKindSpace, SpaceID: "school"}, Shares: 1}},
		ContactAttributions: []FinancialAttributionWeight{{ID: "child-a", Shares: 1}}, EffectiveFromISO: "2026-09-01",
	}
	if err := validSaveReq.Validate(); err != nil {
		t.Fatalf("validSaveReq failed: %v", err)
	}
	bSave := validSaveReq
	bSave.OwnerSpaceID = " bad "
	if err := bSave.Validate(); err == nil {
		t.Fatal("expected error for bad ownerSpaceID in save request")
	}
	bSave = validSaveReq
	bSave.ExpectedRevision = -1
	if err := bSave.Validate(); err == nil {
		t.Fatal("expected error for negative expectedRevision in save request")
	}
	bSave = validSaveReq
	bSave.State = "invalid_state"
	if err := bSave.Validate(); err == nil {
		t.Fatal("expected error for invalid state in save request")
	}
	bSave = validSaveReq
	bSave.ReportingSpaceSide = "invalid_side"
	if err := bSave.Validate(); err == nil {
		t.Fatal("expected error for invalid reportingSpaceSide in save request")
	}
	bSave = validSaveReq
	bSave.EffectiveFromISO = "bad-date"
	if err := bSave.Validate(); err == nil {
		t.Fatal("expected error for bad effectiveFromISO in save request")
	}
	bSave = validSaveReq
	bSave.Payers = nil
	if err := bSave.Validate(); err == nil {
		t.Fatal("expected error for nil payers in save request")
	}
	bSave = validSaveReq
	bSave.Receivers = nil
	if err := bSave.Validate(); err == nil {
		t.Fatal("expected error for nil receivers in save request")
	}
	// party on both sides in save request
	bSave = validSaveReq
	bSave.Receivers = []FinancialPartyWeight{{Party: FinancialPartyRef{Kind: FinancialPartyKindSpace, SpaceID: "family"}, Shares: 1}}
	if err := bSave.Validate(); err == nil {
		t.Fatal("expected error for party on both sides in save request")
	}
	// reporting space not on reported side in save request
	bSave = validSaveReq
	bSave.ReportingSpaceSide = FinancialAgreementSideReceiver
	if err := bSave.Validate(); err == nil {
		t.Fatal("expected error for reporting space not on reported side in save request")
	}
	// bad contactAttributions / assetAttributions in save request
	bSave = validSaveReq
	bSave.ContactAttributions = []FinancialAttributionWeight{{ID: " bad ", Shares: 1}}
	if err := bSave.Validate(); err == nil {
		t.Fatal("expected error for bad contact attribution in save request")
	}
	bSave = validSaveReq
	bSave.AssetAttributions = []FinancialAttributionWeight{{ID: " bad ", Shares: 1}}
	if err := bSave.Validate(); err == nil {
		t.Fatal("expected error for bad asset attribution in save request")
	}

	// 14. ConfirmFinancialAgreementRequest
	confReq := ConfirmFinancialAgreementRequest{
		ReportingSpaceID: "family", OwnerSpaceID: "school", AgreementID: "deal-1", ExpectedRevision: 2, TermsRevision: 1, OperationID: "op-2",
		Party: FinancialPartyRef{Kind: FinancialPartyKindSpace, SpaceID: "family"},
	}
	if err := confReq.Validate(); err != nil {
		t.Fatalf("confReq failed: %v", err)
	}
	bConfReq := confReq
	bConfReq.ReportingSpaceID = " bad "
	if err := bConfReq.Validate(); err == nil {
		t.Fatal("expected error for bad reportingSpaceID in confirm request")
	}
	bConfReq = confReq
	bConfReq.ExpectedRevision = 0
	if err := bConfReq.Validate(); err == nil {
		t.Fatal("expected error for expectedRevision < 1 in confirm request")
	}
	bConfReq = confReq
	bConfReq.Party = FinancialPartyRef{Kind: "bad"}
	if err := bConfReq.Validate(); err == nil {
		t.Fatal("expected error for bad party in confirm request")
	}

	// 15. RevokeFinancialAgreementGrantRequest
	revokeReq := RevokeFinancialAgreementGrantRequest{
		ReportingSpaceID: "family", OwnerSpaceID: "school", AgreementID: "deal-1", ExpectedRevision: 2, TermsRevision: 1, OperationID: "op-2",
		Party: FinancialPartyRef{Kind: FinancialPartyKindSpace, SpaceID: "family"}, Reason: "no longer needed",
	}
	if err := revokeReq.Validate(); err != nil {
		t.Fatalf("revokeReq failed: %v", err)
	}
	bRevoke := revokeReq
	bRevoke.ReportingSpaceID = " bad "
	if err := bRevoke.Validate(); err == nil {
		t.Fatal("expected error for bad reportingSpaceID in revoke request")
	}

	// 16. validateReportingSpaceSide and confirmationRevoked coverage
	tPayers := []FinancialPartyShare{{Party: FinancialPartyRef{Kind: FinancialPartyKindSpace, SpaceID: "space1"}, AmountMinor: 100}}
	tReceivers := []FinancialPartyShare{{Party: FinancialPartyRef{Kind: FinancialPartyKindSpace, SpaceID: "space1"}, AmountMinor: 100}}
	if err := validateReportingSpaceSide("space1", FinancialAgreementSidePayer, tPayers, tReceivers); err == nil {
		t.Fatal("expected error for same reporting space on both sides")
	}
	if err := validateReportingSpaceSide("other", FinancialAgreementSidePayer, tPayers, nil); err == nil {
		t.Fatal("expected error for reporting space not on reported side")
	}

	// line 201: fact.Validate returning validateReportingSpaceSide error
	noSpacePayer := validFinancialAgreement()
	noSpacePayer.Payers = []FinancialPartyShare{
		{Party: FinancialPartyRef{Kind: FinancialPartyKindContact, SpaceID: "family", ContactID: "mom"}, AmountMinor: 4000},
	}
	if err := noSpacePayer.Validate(); err == nil {
		t.Fatal("expected error for space not explicit party on reported side")
	}

	// line 266: confirmationRevoked returning false
	if confirmationRevoked(nil, FinancialPartyRef{Kind: FinancialPartyKindSpace, SpaceID: "s1"}) {
		t.Fatal("expected false for empty confirmations")
	}
}


