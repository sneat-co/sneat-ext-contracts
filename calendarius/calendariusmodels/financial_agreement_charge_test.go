package calendariusmodels

import "testing"

func TestFinancialAgreementChargeRejectsNormalizedReconciliation(t *testing.T) {
	amount := int64(4000)
	scope := FinancialEnrollmentScope{Kind: FinancialEnrollmentScopeWholeHappening, Subjects: []FinancialEnrollmentSubjectRef{}}
	fact := FinancialAgreementChargeFact{OwnerSpaceID: "space1", ReportingSpaceID: "space1", AgreementID: "agreement1", EnrollmentID: "enrollment1", EnrollmentScope: &scope, HappeningID: "happening1", Title: "Music lessons", Regular: true, AgreementRevision: 1, TermsRevision: 1, ChargeID: "charge1", Direction: FinancialChargeDirectionExpense, AmountMinor: &amount, Currency: "EUR", EconomicPeriod: FinancialChargePeriod{StartDate: "2026-09-01", EndDate: "2026-09-30"}, TemporalBasis: FinancialChargeBasisNormalized, BillingTiming: FinancialChargeBillingUnknown, OwnerTimezone: "UTC", InvoiceReconciliationEligible: true, AcceptedPrice: AcceptedPriceTerms{PriceID: "price1", PriceRevision: 0, AmountMinor: 4000, Currency: "EUR", Quantity: 1, Term: FinancialCommitmentTerm{Unit: "year", Length: 1}}, ContactAttributions: []FinancialAttribution{}, AssetAttributions: []FinancialAttribution{}, Status: FinancialChargeStatusAvailable, Diagnostics: []string{}}
	if fact.Validate() == nil {
		t.Fatal("normalized comparison accepted as invoice reconciliation candidate")
	}
	fact.InvoiceReconciliationEligible = false
	if err := fact.Validate(); err != nil {
		t.Fatal(err)
	}
	zero := int64(0)
	fact.AmountMinor = &zero
	if err := fact.Validate(); err != nil {
		t.Fatal(err)
	}
	fact.OwnerTimezone = ""
	if fact.Validate() == nil {
		t.Fatal("available charge with empty owner timezone accepted")
	}
	fact.OwnerTimezone = "UTC"
	fact.Status, fact.AmountMinor, fact.Direction, fact.Diagnostics = FinancialChargeStatusUnavailable, nil, "", []string{"partial_billing_period_policy_unknown"}
	if err := fact.Validate(); err != nil {
		t.Fatal(err)
	}
	fact.OwnerTimezone = ""
	if err := fact.Validate(); err != nil {
		t.Fatal(err)
	}
	fact.InvoiceReconciliationEligible = true
	if fact.Validate() == nil {
		t.Fatal("unavailable charge accepted as invoice reconciliation candidate")
	}
}

func TestFinancialAgreementChargeQueryAndPageBounds(t *testing.T) {
	if err := (FinancialAgreementChargeQuery{ReportingSpaceID: "space1", FromMonthISO: "2026-09", Months: 12, PageSize: 256}).Validate(); err != nil {
		t.Fatal(err)
	}
	for _, query := range []FinancialAgreementChargeQuery{{ReportingSpaceID: "space1", FromMonthISO: "2026-09", Months: 13, PageSize: 1}, {ReportingSpaceID: "space1", FromMonthISO: "9999-12", Months: 2, PageSize: 1}, {ReportingSpaceID: "space1", FromMonthISO: "2026-09", Months: 1, PageSize: 257}} {
		if query.Validate() == nil {
			t.Fatalf("accepted %+v", query)
		}
	}
	if err := (FinancialAgreementChargePage{Charges: []FinancialAgreementChargeFact{}, Coverages: []FinancialAgreementCoverageFact{}, SnapshotConsistent: true, SnapshotDigest: "digest"}).Validate(); err != nil {
		t.Fatal(err)
	}
	if err := (FinancialAgreementChargePage{Charges: []FinancialAgreementChargeFact{}, Coverages: []FinancialAgreementCoverageFact{}, HasMore: true, NextCursor: "next", SnapshotConsistent: true, SnapshotDigest: "digest"}).Validate(); err != nil {
		t.Fatal(err)
	}
	if err := (FinancialAgreementChargePage{Charges: []FinancialAgreementChargeFact{}, Coverages: []FinancialAgreementCoverageFact{}, HasMore: true, NextCursor: "next", SnapshotConsistent: true}).Validate(); err == nil {
		t.Fatal("snapshot consistency accepted without digest")
	}
}

func TestFinancialAgreementCharge_AdditionalCoverage(t *testing.T) {
	// Query: invalid reportingSpaceID
	q := FinancialAgreementChargeQuery{ReportingSpaceID: "", FromMonthISO: "2026-09", Months: 1, PageSize: 10}
	if err := q.Validate(); err == nil {
		t.Fatal("expected error for empty reportingSpaceID")
	}

	amount := int64(4000)
	scope := FinancialEnrollmentScope{Kind: FinancialEnrollmentScopeWholeHappening, Subjects: []FinancialEnrollmentSubjectRef{}}
	validFact := FinancialAgreementChargeFact{
		OwnerSpaceID: "space1", ReportingSpaceID: "space1", AgreementID: "agreement1", EnrollmentID: "enrollment1",
		EnrollmentScope: &scope, HappeningID: "happening1", Title: "Music lessons", Regular: true,
		AgreementRevision: 1, TermsRevision: 1, ChargeID: "charge1", Direction: FinancialChargeDirectionExpense,
		AmountMinor: &amount, Currency: "EUR", EconomicPeriod: FinancialChargePeriod{StartDate: "2026-09-01", EndDate: "2026-09-30"},
		TemporalBasis: FinancialChargeBasisOccurrence, BillingTiming: FinancialChargeBillingUnknown, OwnerTimezone: "UTC",
		AcceptedPrice: AcceptedPriceTerms{PriceID: "price1", PriceRevision: 0, AmountMinor: 4000, Currency: "EUR", Quantity: 1, Term: FinancialCommitmentTerm{Unit: "year", Length: 1}},
		ContactAttributions: []FinancialAttribution{}, AssetAttributions: []FinancialAttribution{}, Status: FinancialChargeStatusAvailable, Diagnostics: []string{},
	}
	if err := validFact.Validate(); err != nil {
		t.Fatalf("validFact failed: %v", err)
	}

	// 1. invalid ID in fact
	f := validFact
	f.OwnerSpaceID = ""
	if err := f.Validate(); err == nil {
		t.Fatal("expected error for empty ownerSpaceID")
	}

	// 2. invalid revision (termsRevision > agreementRevision)
	f = validFact
	f.TermsRevision = 2
	if err := f.Validate(); err == nil {
		t.Fatal("expected error for termsRevision > agreementRevision")
	}

	// 3. invalid accepted price
	f = validFact
	f.AcceptedPrice.Currency = ""
	if err := f.Validate(); err == nil {
		t.Fatal("expected error for invalid acceptedPrice")
	}

	// 4. invalid enrollment scope
	f = validFact
	badScope := FinancialEnrollmentScope{Kind: "invalid_kind"}
	f.EnrollmentScope = &badScope
	if err := f.Validate(); err == nil {
		t.Fatal("expected error for invalid enrollmentScope")
	}

	// 5. Currency mismatch
	f = validFact
	f.Currency = "USD"
	if err := f.Validate(); err == nil {
		t.Fatal("expected error for currency mismatch")
	}

	// 6. invalid timezone
	f = validFact
	f.OwnerTimezone = "not/a_valid_timezone"
	if err := f.Validate(); err == nil {
		t.Fatal("expected error for invalid timezone")
	}

	// 7. invalid billingTiming
	f = validFact
	f.BillingTiming = "invalid_timing"
	if err := f.Validate(); err == nil {
		t.Fatal("expected error for invalid billingTiming")
	}

	// 8. KnownDue billing timing without dueDate
	f = validFact
	f.BillingTiming = FinancialChargeBillingKnownDue
	f.DueDate = ""
	if err := f.Validate(); err == nil {
		t.Fatal("expected error for KnownDue without dueDate")
	}

	// 9. KnownDue billing timing with valid dueDate
	f = validFact
	f.BillingTiming = FinancialChargeBillingKnownDue
	f.DueDate = "2026-09-15"
	if err := f.Validate(); err != nil {
		t.Fatalf("unexpected error for KnownDue with valid dueDate: %v", err)
	}

	// 10. unknown billing timing with dueDate set
	f = validFact
	f.BillingTiming = FinancialChargeBillingUnknown
	f.DueDate = "2026-09-15"
	if err := f.Validate(); err == nil {
		t.Fatal("expected error for Unknown billingTiming with dueDate")
	}

	// 11. invalid temporalBasis
	f = validFact
	f.TemporalBasis = "invalid_basis"
	if err := f.Validate(); err == nil {
		t.Fatal("expected error for invalid temporalBasis")
	}

	// 12. ContractPeriod temporalBasis
	f = validFact
	f.TemporalBasis = FinancialChargeBasisContractPeriod
	if err := f.Validate(); err != nil {
		t.Fatalf("unexpected error for ContractPeriod: %v", err)
	}

	// 13. validateAttributions error
	f = validFact
	f.ContactAttributions = []FinancialAttribution{{ID: "", AmountMinor: 4000}}
	if err := f.Validate(); err == nil {
		t.Fatal("expected error for invalid contact attribution")
	}
	f = validFact
	f.AssetAttributions = []FinancialAttribution{{ID: "", AmountMinor: 4000}}
	if err := f.Validate(); err == nil {
		t.Fatal("expected error for invalid asset attribution")
	}
	f = validFact
	f.Status = "unknown_status"
	if err := f.Validate(); err == nil {
		t.Fatal("expected error for invalid status charge")
	}

	p := FinancialAgreementChargePage{
		Charges:             []FinancialAgreementChargeFact{{OwnerSpaceID: ""}},
		Coverages:           []FinancialAgreementCoverageFact{},
		SnapshotConsistent: true, SnapshotDigest: "d",
	}
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for invalid charge in page")
	}
	p = FinancialAgreementChargePage{
		Charges:             []FinancialAgreementChargeFact{},
		Coverages:           []FinancialAgreementCoverageFact{{OwnerSpaceID: ""}},
		SnapshotConsistent: true, SnapshotDigest: "d",
	}
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for invalid coverage in page")
	}

	// 15. FinancialAgreementCoverageFact
	validCov := FinancialAgreementCoverageFact{
		OwnerSpaceID: "space1", ReportingSpaceID: "space1", AgreementID: "ag1", EnrollmentID: "en1", HappeningID: "h1",
		EnrollmentScope: &scope, State: FinancialAgreementStateConfirmed, EffectiveFromISO: "2026-09-01", Verified: true,
	}
	if err := validCov.Validate(); err != nil {
		t.Fatalf("validCov failed: %v", err)
	}
	cov := validCov
	cov.OwnerSpaceID = ""
	if err := cov.Validate(); err == nil {
		t.Fatal("expected error for empty ownerSpaceID in coverage")
	}
	cov = validCov
	cov.Verified = false
	if err := cov.Validate(); err == nil {
		t.Fatal("expected error for unverified coverage")
	}
	cov = validCov
	cov.EnrollmentScope = &badScope
	if err := cov.Validate(); err == nil {
		t.Fatal("expected error for invalid scope in coverage")
	}
	cov = validCov
	cov.State = "invalid_state"
	if err := cov.Validate(); err == nil {
		t.Fatal("expected error for invalid state in coverage")
	}
	cov = validCov
	cov.State = FinancialAgreementStateRecordedExternal
	cov.EffectiveToISO = "2026-10-01"
	if err := cov.Validate(); err != nil {
		t.Fatalf("unexpected error for RecordedExternal with valid toISO: %v", err)
	}
}
