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
		{SpaceID: "space1", FromMonthISO: "9999-12", Months: 2, PageSize: 1},
		{SpaceID: "space1", FromMonthISO: "2026-09", Months: 13, PageSize: 1},
		{SpaceID: "space1", FromMonthISO: "2026-09", Months: 1, PageSize: 65},
		{SpaceID: "space1", FromMonthISO: "2026-09", Months: 1, PageSize: 1, Cursor: "bad\nvalue"},
		{SpaceID: "unsafe/space", FromMonthISO: "2026-09", Months: 1, PageSize: 1},
		{SpaceID: "space\x00", FromMonthISO: "2026-09", Months: 1, PageSize: 1},
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
	fact := FinancialCommitmentFact{SpaceID: "space1", HappeningID: "h1", Title: "Utility", Prices: []FinancialCommitmentPriceFact{{PriceID: "p1", AmountMinor: 12000, Currency: "EUR", ExpenseQuantity: 1, Term: FinancialCommitmentTerm{Unit: "month", Length: 1}, Periods: []FinancialCommitmentPeriodFact{}}}, Occurrences: occurrences, OccurrencesIncompleteReason: "occurrence_limit"}
	if err := (FinancialCommitmentPage{Facts: []FinancialCommitmentFact{fact}, SnapshotConsistent: true}).Validate(); err != nil {
		t.Fatal(err)
	}
	fact.Occurrences = fact.Occurrences[:1]
	if err := fact.Validate(); err == nil {
		t.Fatal("accepted a partial fact below the declared cap")
	}
}

func TestFinancialCommitmentPageRejectsDuplicateAndUnboundedFacts(t *testing.T) {
	fact := FinancialCommitmentFact{
		SpaceID: "space1", HappeningID: "h1", Title: "Utility",
		Prices: []FinancialCommitmentPriceFact{}, Occurrences: []FinancialCommitmentOccurrenceFact{},
	}
	duplicates := FinancialCommitmentPage{Facts: []FinancialCommitmentFact{fact, fact}, SnapshotConsistent: true}
	if duplicates.Validate() == nil {
		t.Fatal("accepted duplicate source facts")
	}
	unbounded := FinancialCommitmentPage{Facts: make([]FinancialCommitmentFact, MaxFinancialCommitmentPageSize+1), SnapshotConsistent: true}
	if unbounded.Validate() == nil {
		t.Fatal("accepted more facts than the page bound")
	}
	badCursor := FinancialCommitmentPage{Facts: []FinancialCommitmentFact{}, HasMore: true, NextCursor: " padded ", IncompleteReason: "source_collection_mutable_between_pages"}
	if badCursor.Validate() == nil {
		t.Fatal("accepted a padded response cursor")
	}
	if (FinancialCommitmentPage{Facts: []FinancialCommitmentFact{}}).Validate() == nil {
		t.Fatal("accepted a false snapshot claim without an explicit reason")
	}
	if err := (FinancialCommitmentPage{Facts: []FinancialCommitmentFact{}, HasMore: true, NextCursor: "opaque", IncompleteReason: "source_collection_mutable_between_pages"}).Validate(); err != nil {
		t.Fatalf("rejected an explicitly incomplete mutable page: %v", err)
	}
}

func TestFinancialCommitmentFactRejectsInvalidDatesAndDuplicateReferences(t *testing.T) {
	valid := FinancialCommitmentFact{
		SpaceID: "space1", HappeningID: "h1", Title: "Utility",
		ActiveFromISO: "2026-02-28", ActiveToISO: "2026-09-30",
		AssetIDs:     []string{"property1"},
		ContactLinks: []FinancialCommitmentContactLink{{ContactID: "alice", Roles: []string{"participant"}}},
		Prices:       []FinancialCommitmentPriceFact{{PriceID: "p1", AmountMinor: 12000, Currency: "EUR", ExpenseQuantity: 1, Term: FinancialCommitmentTerm{Unit: "month", Length: 1}, Periods: []FinancialCommitmentPeriodFact{{OccurrenceID: "month:2026-09", PeriodStartDate: "2026-09-01", PeriodEndDate: "2026-09-30"}}}},
		Occurrences:  []FinancialCommitmentOccurrenceFact{{OccurrenceID: "o1", ScheduledDate: "2026-09-01", EffectiveDate: "2026-09-01"}},
	}
	if err := valid.Validate(); err != nil {
		t.Fatal(err)
	}
	invalid := valid
	invalid.HappeningID = "unsafe/happening"
	if invalid.Validate() == nil {
		t.Fatal("accepted an unsafe happening ID")
	}
	invalid = valid
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
	invalid = valid
	invalid.Occurrences = []FinancialCommitmentOccurrenceFact{{OccurrenceID: "o1", ScheduledDate: "2026-09-01", EffectiveDate: "2026-09-01", CancellationFinancialEffect: "waived"}}
	if invalid.Validate() == nil {
		t.Fatal("accepted an invented cancellation financial effect")
	}
	invalid = valid
	invalid.Prices = nil
	if invalid.Validate() == nil {
		t.Fatal("accepted null prices")
	}
	invalid = valid
	invalid.Occurrences = nil
	if invalid.Validate() == nil {
		t.Fatal("accepted null occurrences")
	}
}

func TestFinancialCommitmentPeriodsAreBoundedOwnerEconomicMonths(t *testing.T) {
	valid := FinancialCommitmentPeriodFact{OccurrenceID: "month:2026-02", PeriodStartDate: "2026-02-01", PeriodEndDate: "2026-02-28"}
	if err := valid.Validate(); err != nil {
		t.Fatal(err)
	}
	for _, invalid := range []FinancialCommitmentPeriodFact{
		{OccurrenceID: "month:2026-02", PeriodStartDate: "2026-02-02", PeriodEndDate: "2026-02-28"},
		{OccurrenceID: "month:2026-02", PeriodStartDate: "2026-02-01", PeriodEndDate: "2026-03-01"},
		{OccurrenceID: "unsafe/id", PeriodStartDate: "2026-02-01", PeriodEndDate: "2026-02-28"},
	} {
		if invalid.Validate() == nil {
			t.Fatalf("accepted invalid period: %+v", invalid)
		}
	}
	price := FinancialCommitmentPriceFact{PriceID: "p1", Currency: "EUR", ExpenseQuantity: 1, Term: FinancialCommitmentTerm{Unit: "month", Length: 1}, Periods: []FinancialCommitmentPeriodFact{valid, valid}}
	if price.Validate() == nil {
		t.Fatal("accepted duplicate owner period identities")
	}
	price.Periods = []FinancialCommitmentPeriodFact{
		{OccurrenceID: "month:2026-03", PeriodStartDate: "2026-03-01", PeriodEndDate: "2026-03-31"},
		{OccurrenceID: "another:2026-03", PeriodStartDate: "2026-03-01", PeriodEndDate: "2026-03-31"},
	}
	if price.Validate() == nil {
		t.Fatal("accepted two owner identities for the same economic month")
	}
	price.Periods = []FinancialCommitmentPeriodFact{
		{OccurrenceID: "month:2026-04", PeriodStartDate: "2026-04-01", PeriodEndDate: "2026-04-30"},
		{OccurrenceID: "month:2026-03", PeriodStartDate: "2026-03-01", PeriodEndDate: "2026-03-31"},
	}
	if price.Validate() == nil {
		t.Fatal("accepted economic months out of canonical order")
	}
}
