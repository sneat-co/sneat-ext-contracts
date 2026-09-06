package calendariusmodels

import "testing"

func TestFinancialAgreementChargeRejectsNormalizedReconciliation(t *testing.T) {
	amount := int64(4000)
	fact := FinancialAgreementChargeFact{OwnerSpaceID: "space1", ReportingSpaceID: "space1", AgreementID: "agreement1", EnrollmentID: "enrollment1", HappeningID: "happening1", AgreementRevision: 1, TermsRevision: 1, ChargeID: "charge1", Direction: FinancialChargeDirectionExpense, AmountMinor: &amount, Currency: "EUR", EconomicPeriod: FinancialChargePeriod{StartDate: "2026-09-01", EndDate: "2026-09-30"}, TemporalBasis: FinancialChargeBasisNormalized, BillingTiming: FinancialChargeBillingUnknown, OwnerTimezone: "UTC", InvoiceReconciliationEligible: true, AcceptedPrice: AcceptedPriceTerms{PriceID: "price1", PriceRevision: 0, AmountMinor: 4000, Currency: "EUR", Quantity: 1, Term: FinancialCommitmentTerm{Unit: "year", Length: 1}}, ContactAttributions: []FinancialAttribution{}, AssetAttributions: []FinancialAttribution{}, Status: FinancialChargeStatusAvailable, Diagnostics: []string{}}
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
	fact.Status, fact.AmountMinor, fact.Direction, fact.Diagnostics = FinancialChargeStatusUnavailable, nil, "", []string{"partial_billing_period_policy_unknown"}
	if err := fact.Validate(); err != nil {
		t.Fatal(err)
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
	if err := (FinancialAgreementChargePage{Charges: []FinancialAgreementChargeFact{}, SnapshotConsistent: true}).Validate(); err != nil {
		t.Fatal(err)
	}
	if err := (FinancialAgreementChargePage{Charges: []FinancialAgreementChargeFact{}, HasMore: true, NextCursor: "next", SnapshotConsistent: true}).Validate(); err == nil {
		t.Fatal("multi-page result claimed snapshot consistency")
	}
}
