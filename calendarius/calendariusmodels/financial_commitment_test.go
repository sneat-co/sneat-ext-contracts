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
	if err := (FinancialCommitmentPage{Facts: []FinancialCommitmentFact{fact}}).Validate(); err != nil {
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
	duplicates := FinancialCommitmentPage{Facts: []FinancialCommitmentFact{fact, fact}}
	if duplicates.Validate() == nil {
		t.Fatal("accepted duplicate source facts")
	}
	unbounded := FinancialCommitmentPage{Facts: make([]FinancialCommitmentFact, MaxFinancialCommitmentPageSize+1)}
	if unbounded.Validate() == nil {
		t.Fatal("accepted more facts than the page bound")
	}
	badCursor := FinancialCommitmentPage{Facts: []FinancialCommitmentFact{}, HasMore: true, NextCursor: " padded "}
	if badCursor.Validate() == nil {
		t.Fatal("accepted a padded response cursor")
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

func TestFinancialCommitment_AdditionalCoverage(t *testing.T) {
	// 1. FinancialCommitmentPage
	// HasMore != (NextCursor != "")
	p1 := FinancialCommitmentPage{Facts: []FinancialCommitmentFact{}, HasMore: true, NextCursor: ""}
	if err := p1.Validate(); err == nil {
		t.Fatal("expected error for HasMore without NextCursor")
	}
	// IncompleteReason invalid
	p2 := FinancialCommitmentPage{Facts: []FinancialCommitmentFact{}, IncompleteReason: "invalid"}
	if err := p2.Validate(); err == nil {
		t.Fatal("expected error for invalid IncompleteReason")
	}
	// Facts[i].Validate error
	p3 := FinancialCommitmentPage{Facts: []FinancialCommitmentFact{{SpaceID: ""}}}
	if err := p3.Validate(); err == nil {
		t.Fatal("expected error for invalid Fact in Page")
	}

	// 2. FinancialCommitmentFact
	validFact := FinancialCommitmentFact{
		SpaceID: "space1", HappeningID: "h1", Title: "Utility",
		Prices:       []FinancialCommitmentPriceFact{{PriceID: "p1", AmountMinor: 12000, Currency: "EUR", ExpenseQuantity: 1, Term: FinancialCommitmentTerm{Unit: "month", Length: 1}, Periods: []FinancialCommitmentPeriodFact{}}},
		Occurrences:  []FinancialCommitmentOccurrenceFact{{OccurrenceID: "o1", ScheduledDate: "2026-09-01", EffectiveDate: "2026-09-01"}},
	}
	// invalid spaceID
	f := validFact
	f.SpaceID = ""
	if err := f.Validate(); err == nil {
		t.Fatal("expected error for empty spaceID")
	}
	// empty title
	f = validFact
	f.Title = "   "
	if err := f.Validate(); err == nil {
		t.Fatal("expected error for empty title")
	}
	// negative priceRevision
	f = validFact
	f.PriceRevision = -1
	if err := f.Validate(); err == nil {
		t.Fatal("expected error for negative priceRevision")
	}
	// occurrences exceeds max
	f = validFact
	f.Occurrences = make([]FinancialCommitmentOccurrenceFact, MaxFinancialCommitmentOccurrencesPerFact+1)
	if err := f.Validate(); err == nil {
		t.Fatal("expected error for exceeding max occurrences")
	}
	// occurrencesIncompleteReason invalid
	f = validFact
	f.Occurrences = make([]FinancialCommitmentOccurrenceFact, MaxFinancialCommitmentOccurrencesPerFact)
	f.OccurrencesIncompleteReason = "invalid_reason"
	if err := f.Validate(); err == nil {
		t.Fatal("expected error for invalid OccurrencesIncompleteReason")
	}
	// activeFromISO > activeToISO
	f = validFact
	f.ActiveFromISO = "2026-10-01"
	f.ActiveToISO = "2026-09-01"
	if err := f.Validate(); err == nil {
		t.Fatal("expected error for activeFrom > activeTo")
	}
	// price.ExpenseQuantity <= 0
	f = validFact
	f.Prices = []FinancialCommitmentPriceFact{{PriceID: "p1", ExpenseQuantity: 0}}
	if err := f.Validate(); err == nil {
		t.Fatal("expected error for non-positive expenseQuantity")
	}
	// duplicate priceID
	f = validFact
	pGood := FinancialCommitmentPriceFact{PriceID: "p1", AmountMinor: 100, Currency: "EUR", ExpenseQuantity: 1, Term: FinancialCommitmentTerm{Unit: "month", Length: 1}, Periods: []FinancialCommitmentPeriodFact{}}
	f.Prices = []FinancialCommitmentPriceFact{pGood, pGood}
	if err := f.Validate(); err == nil {
		t.Fatal("expected error for duplicate priceID")
	}
	// price.Validate() fails
	f = validFact
	f.Prices = []FinancialCommitmentPriceFact{{PriceID: "p1", AmountMinor: -1, Currency: "EUR", ExpenseQuantity: 1, Term: FinancialCommitmentTerm{Unit: "month", Length: 1}, Periods: []FinancialCommitmentPeriodFact{}}}
	if err := f.Validate(); err == nil {
		t.Fatal("expected error for invalid price")
	}
	// occurrenceID invalid
	f = validFact
	f.Occurrences = []FinancialCommitmentOccurrenceFact{{OccurrenceID: ""}}
	if err := f.Validate(); err == nil {
		t.Fatal("expected error for empty occurrenceID")
	}
	// duplicate occurrenceID
	f = validFact
	occGood := FinancialCommitmentOccurrenceFact{OccurrenceID: "o1", ScheduledDate: "2026-09-01", EffectiveDate: "2026-09-01"}
	f.Occurrences = []FinancialCommitmentOccurrenceFact{occGood, occGood}
	if err := f.Validate(); err == nil {
		t.Fatal("expected error for duplicate occurrenceID")
	}
	// occurrence contactLinks invalid
	f = validFact
	occWithBadLink := occGood
	occWithBadLink.ContactLinks = []FinancialCommitmentContactLink{{ContactID: ""}}
	f.Occurrences = []FinancialCommitmentOccurrenceFact{occWithBadLink}
	if err := f.Validate(); err == nil {
		t.Fatal("expected error for invalid contactLinks in occurrence")
	}

	// 3. FinancialCommitmentPriceFact
	// invalid priceID
	pf := FinancialCommitmentPriceFact{PriceID: ""}
	if err := pf.Validate(); err == nil {
		t.Fatal("expected error for empty priceID")
	}
	// invalid fields
	pf = FinancialCommitmentPriceFact{PriceID: "p1", Currency: "", ExpenseQuantity: 1}
	if err := pf.Validate(); err == nil {
		t.Fatal("expected error for empty currency")
	}
	// nil periods
	pf = FinancialCommitmentPriceFact{PriceID: "p1", AmountMinor: 100, Currency: "EUR", ExpenseQuantity: 1, Term: FinancialCommitmentTerm{Unit: "month", Length: 1}, Periods: nil}
	if err := pf.Validate(); err == nil {
		t.Fatal("expected error for nil periods")
	}
	// periods exceeds max
	pf = FinancialCommitmentPriceFact{PriceID: "p1", AmountMinor: 100, Currency: "EUR", ExpenseQuantity: 1, Term: FinancialCommitmentTerm{Unit: "month", Length: 1}, Periods: make([]FinancialCommitmentPeriodFact, MaxFinancialCommitmentPeriodsPerPrice+1)}
	if err := pf.Validate(); err == nil {
		t.Fatal("expected error for exceeding max periods")
	}
	// period.Validate() fails
	badPeriod := FinancialCommitmentPeriodFact{OccurrenceID: "", PeriodStartDate: "2026-09-01", PeriodEndDate: "2026-09-30"}
	pf = FinancialCommitmentPriceFact{PriceID: "p1", AmountMinor: 100, Currency: "EUR", ExpenseQuantity: 1, Term: FinancialCommitmentTerm{Unit: "month", Length: 1}, Periods: []FinancialCommitmentPeriodFact{badPeriod}}
	if err := pf.Validate(); err == nil {
		t.Fatal("expected error for invalid period in price")
	}

	// 4. FinancialCommitmentPeriodFact
	// invalid dates
	period := FinancialCommitmentPeriodFact{OccurrenceID: "p1", PeriodStartDate: "2026-9-1", PeriodEndDate: "2026-09-30"}
	if err := period.Validate(); err == nil {
		t.Fatal("expected error for invalid date format")
	}

	// 5. Helpers: validateUniqueIDs with invalid ID
	if err := validateUniqueIDs("ids", []string{""}); err == nil {
		t.Fatal("expected error for empty ID in validateUniqueIDs")
	}
	// role untrimmed or empty
	links := []FinancialCommitmentContactLink{{ContactID: "c1", Roles: []string{""}}}
	if err := validateCommitmentContactLinks("links", links); err == nil {
		t.Fatal("expected error for empty role")
	}
	// duplicate role
	links = []FinancialCommitmentContactLink{{ContactID: "c1", Roles: []string{"admin", "admin"}}}
	if err := validateCommitmentContactLinks("links", links); err == nil {
		t.Fatal("expected error for duplicate role")
	}
}
