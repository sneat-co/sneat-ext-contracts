package calendariusmodels

import (
	"fmt"
	"testing"
)

func TestFinancialCommitmentQueryBounds(t *testing.T) {
	valid := FinancialCommitmentQuery{SpaceID: "space1", FromMonthISO: "2026-09", Months: 12, PageSize: 64}
	if err := valid.Validate(); err != nil {
		t.Fatal(err)
	}
	for _, invalid := range []FinancialCommitmentQuery{
		{SpaceID: "space1", FromMonthISO: "2026-13", Months: 1, PageSize: 1},
		{SpaceID: "space1", FromMonthISO: "2026-09", Months: 13, PageSize: 1},
		{SpaceID: "space1", FromMonthISO: "2026-09", Months: 1, PageSize: 65},
		{SpaceID: "space1", FromMonthISO: "2026-09", Months: 1, PageSize: 1, Cursor: "bad\nvalue"},
	} {
		if invalid.Validate() == nil {
			t.Fatalf("accepted invalid query: %+v", invalid)
		}
	}
}

func TestFinancialCommitmentPageRequiresWholeFacts(t *testing.T) {
	if (FinancialCommitmentPage{}).Validate() == nil {
		t.Fatal("accepted nil facts")
	}
	occurrences := make([]FinancialCommitmentOccurrenceFact, MaxFinancialCommitmentOccurrencesPerFact)
	for i := range occurrences {
		occurrences[i] = FinancialCommitmentOccurrenceFact{OccurrenceID: fmt.Sprintf("o%d", i), ScheduledDate: "2026-09-01", EffectiveDate: "2026-09-01"}
	}
	fact := FinancialCommitmentFact{SpaceID: "space1", HappeningID: "h1", Title: "Utility", Prices: []FinancialCommitmentPriceFact{{PriceID: "p1", AmountMinor: 12000, Currency: "EUR", ExpenseQuantity: 1, Term: FinancialCommitmentTerm{Unit: "month", Length: 1}}}, Occurrences: occurrences, OccurrencesIncompleteReason: "occurrence_limit"}
	if err := (FinancialCommitmentPage{Facts: []FinancialCommitmentFact{fact}}).Validate(); err != nil {
		t.Fatal(err)
	}
	fact.Occurrences = fact.Occurrences[:1]
	if err := fact.Validate(); err == nil {
		t.Fatal("accepted a partial fact below the declared cap")
	}
}

func TestFinancialCommitmentFactRejectsInvalidDatesAndDuplicateReferences(t *testing.T) {
	valid := FinancialCommitmentFact{
		SpaceID: "space1", HappeningID: "h1", Title: "Utility",
		ActiveFromISO: "2026-02-28", ActiveToISO: "2026-09-30",
		AssetIDs:     []string{"property1"},
		ContactLinks: []FinancialCommitmentContactLink{{ContactID: "alice", Roles: []string{"participant"}}},
		Prices:       []FinancialCommitmentPriceFact{{PriceID: "p1", AmountMinor: 12000, Currency: "EUR", ExpenseQuantity: 1, Term: FinancialCommitmentTerm{Unit: "month", Length: 1}}},
		Occurrences:  []FinancialCommitmentOccurrenceFact{{OccurrenceID: "o1", ScheduledDate: "2026-09-01", EffectiveDate: "2026-09-01"}},
	}
	if err := valid.Validate(); err != nil {
		t.Fatal(err)
	}
	invalid := valid
	invalid.ActiveFromISO = "2026-02-31"
	if invalid.Validate() == nil {
		t.Fatal("accepted a rolled-over active date")
	}
	invalid = valid
	invalid.AssetIDs = []string{"property1", "property1"}
	if invalid.Validate() == nil {
		t.Fatal("accepted duplicate asset IDs")
	}
	invalid = valid
	invalid.ContactLinks = append(invalid.ContactLinks, invalid.ContactLinks[0])
	if invalid.Validate() == nil {
		t.Fatal("accepted duplicate contact links")
	}
	invalid = valid
	invalid.Occurrences[0].ScheduledDate = "2026-09-31"
	if invalid.Validate() == nil {
		t.Fatal("accepted an invalid occurrence date")
	}
}
