package calendariusmodels

import (
	"context"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	MaxFinancialCommitmentMonths             = 12
	MaxFinancialCommitmentPageSize           = 64
	MaxFinancialCommitmentOccurrencesPerFact = 4096
	MaxFinancialCommitmentCursorBytes        = 512
	MaxJavaScriptSafeInteger                 = 9_007_199_254_740_991
)

type FinancialCommitmentQuery struct {
	SpaceID      string `json:"spaceID"`
	FromMonthISO string `json:"fromMonthISO"`
	Months       int    `json:"months"`
	PageSize     int    `json:"pageSize"`
	Cursor       string `json:"cursor,omitempty"`
}

func (q FinancialCommitmentQuery) Validate() error {
	if err := validateEventHappeningText("spaceID", q.SpaceID, EventHappeningIDMaxBytes, true); err != nil {
		return err
	}
	if !validMonthISO(q.FromMonthISO) {
		return fmt.Errorf("fromMonthISO must be a real YYYY-MM month")
	}
	if q.Months < 1 || q.Months > MaxFinancialCommitmentMonths {
		return fmt.Errorf("months must be between 1 and %d", MaxFinancialCommitmentMonths)
	}
	fromMonth, _ := time.Parse("2006-01", q.FromMonthISO)
	if fromMonth.Year() > 9999 || fromMonth.AddDate(0, q.Months-1, 0).Year() > 9999 {
		return fmt.Errorf("query window exceeds year 9999")
	}
	if q.PageSize < 1 || q.PageSize > MaxFinancialCommitmentPageSize {
		return fmt.Errorf("pageSize must be between 1 and %d", MaxFinancialCommitmentPageSize)
	}
	if len(q.Cursor) > MaxFinancialCommitmentCursorBytes || !utf8.ValidString(q.Cursor) || strings.TrimSpace(q.Cursor) != q.Cursor || strings.ContainsAny(q.Cursor, "\r\n") {
		return fmt.Errorf("cursor must be an opaque finite single-line value")
	}
	return nil
}

type FinancialCommitmentPage struct {
	Facts            []FinancialCommitmentFact `json:"facts"`
	HasMore          bool                      `json:"hasMore"`
	NextCursor       string                    `json:"nextCursor,omitempty"`
	IncompleteReason string                    `json:"incompleteReason,omitempty"`
}

// Validate requires whole-fact pagination. A provider never splits one
// happening across pages: an occurrence cap is disclosed on that fact rather
// than duplicating or silently dropping the source behind a cursor.
func (p FinancialCommitmentPage) Validate() error {
	if p.Facts == nil {
		return fmt.Errorf("facts must be a non-nil array")
	}
	if len(p.Facts) > MaxFinancialCommitmentPageSize {
		return fmt.Errorf("facts exceeds maximum %d", MaxFinancialCommitmentPageSize)
	}
	if p.HasMore != (p.NextCursor != "") {
		return fmt.Errorf("hasMore and nextCursor are inconsistent")
	}
	if p.IncompleteReason != "" && p.IncompleteReason != "query_limit" {
		return fmt.Errorf("unsupported incompleteReason")
	}
	if len(p.NextCursor) > MaxFinancialCommitmentCursorBytes || !utf8.ValidString(p.NextCursor) || strings.TrimSpace(p.NextCursor) != p.NextCursor || strings.ContainsAny(p.NextCursor, "\r\n") {
		return fmt.Errorf("nextCursor is invalid")
	}
	seenFacts := make(map[string]struct{}, len(p.Facts))
	for i := range p.Facts {
		if err := p.Facts[i].Validate(); err != nil {
			return fmt.Errorf("facts[%d]: %w", i, err)
		}
		key := p.Facts[i].SpaceID + "\x00" + p.Facts[i].HappeningID
		if _, exists := seenFacts[key]; exists {
			return fmt.Errorf("facts contains duplicate Space and happening identity")
		}
		seenFacts[key] = struct{}{}
	}
	return nil
}

type FinancialCommitmentFact struct {
	SpaceID                     string                              `json:"spaceID"`
	HappeningID                 string                              `json:"happeningID"`
	Title                       string                              `json:"title"`
	PriceRevision               int                                 `json:"priceRevision"`
	ActiveFromISO               string                              `json:"activeFromISO,omitempty"`
	ActiveToISO                 string                              `json:"activeToISO,omitempty"`
	AssetIDs                    []string                            `json:"assetIDs,omitempty"`
	ContactLinks                []FinancialCommitmentContactLink    `json:"contactLinks,omitempty"`
	Prices                      []FinancialCommitmentPriceFact      `json:"prices"`
	Occurrences                 []FinancialCommitmentOccurrenceFact `json:"occurrences"`
	OccurrencesIncompleteReason string                              `json:"occurrencesIncompleteReason,omitempty"`
	Diagnostics                 []string                            `json:"diagnostics,omitempty"`
}

func (f FinancialCommitmentFact) Validate() error {
	if err := validateEventHappeningText("spaceID", f.SpaceID, EventHappeningIDMaxBytes, true); err != nil {
		return err
	}
	if err := validateEventHappeningText("happeningID", f.HappeningID, EventHappeningIDMaxBytes, true); err != nil {
		return err
	}
	if strings.TrimSpace(f.Title) == "" {
		return fmt.Errorf("title is required")
	}
	if f.PriceRevision < 0 || int64(f.PriceRevision) > MaxJavaScriptSafeInteger {
		return fmt.Errorf("priceRevision must be a non-negative JavaScript-safe integer")
	}
	if f.Prices == nil || f.Occurrences == nil {
		return fmt.Errorf("prices and occurrences must be non-nil arrays")
	}
	if len(f.Occurrences) > MaxFinancialCommitmentOccurrencesPerFact {
		return fmt.Errorf("occurrences exceeds maximum %d", MaxFinancialCommitmentOccurrencesPerFact)
	}
	if f.OccurrencesIncompleteReason != "" && len(f.Occurrences) != MaxFinancialCommitmentOccurrencesPerFact {
		return fmt.Errorf("occurrence incompleteness requires the exact per-fact cap")
	}
	if f.OccurrencesIncompleteReason != "" && f.OccurrencesIncompleteReason != "occurrence_limit" {
		return fmt.Errorf("unsupported occurrencesIncompleteReason")
	}
	if !validOptionalISODate(f.ActiveFromISO) || !validOptionalISODate(f.ActiveToISO) {
		return fmt.Errorf("active bounds must be real ISO dates")
	}
	if f.ActiveFromISO != "" && f.ActiveToISO != "" && f.ActiveFromISO > f.ActiveToISO {
		return fmt.Errorf("activeFromISO must not be after activeToISO")
	}
	if err := validateUniqueIDs("assetIDs", f.AssetIDs); err != nil {
		return err
	}
	if err := validateCommitmentContactLinks("contactLinks", f.ContactLinks); err != nil {
		return err
	}
	seen := map[string]struct{}{}
	for i, price := range f.Prices {
		if price.ExpenseQuantity <= 0 {
			return fmt.Errorf("prices[%d] is not an expense", i)
		}
		if _, ok := seen[price.PriceID]; ok {
			return fmt.Errorf("duplicate priceID %q", price.PriceID)
		}
		seen[price.PriceID] = struct{}{}
		if err := price.Validate(); err != nil {
			return fmt.Errorf("prices[%d]: %w", i, err)
		}
	}
	seenOccurrences := map[string]struct{}{}
	for i, occurrence := range f.Occurrences {
		if err := validateEventHappeningText("occurrenceID", occurrence.OccurrenceID, EventHappeningIDMaxBytes, true); err != nil {
			return fmt.Errorf("occurrences[%d]: %w", i, err)
		}
		if !validISODate(occurrence.ScheduledDate) || !validISODate(occurrence.EffectiveDate) {
			return fmt.Errorf("occurrences[%d] identity and dates are required", i)
		}
		if !validFinancialEffect(occurrence.CancellationFinancialEffect) || !validFinancialEffect(occurrence.DateAdjustmentFinancialEffect) {
			return fmt.Errorf("occurrences[%d] financial effect is invalid", i)
		}
		if _, exists := seenOccurrences[occurrence.OccurrenceID]; exists {
			return fmt.Errorf("duplicate occurrenceID %q", occurrence.OccurrenceID)
		}
		seenOccurrences[occurrence.OccurrenceID] = struct{}{}
		if err := validateCommitmentContactLinks(fmt.Sprintf("occurrences[%d].contactLinks", i), occurrence.ContactLinks); err != nil {
			return err
		}
	}
	return nil
}

type FinancialCommitmentContactLink struct {
	ContactID string   `json:"contactID"`
	Roles     []string `json:"roles,omitempty"`
}
type FinancialCommitmentTerm struct {
	Unit   string `json:"unit"`
	Length int64  `json:"length"`
}
type FinancialCommitmentPriceFact struct {
	PriceID         string                  `json:"priceID"`
	AmountMinor     int64                   `json:"amountMinor"`
	Currency        string                  `json:"currency"`
	ExpenseQuantity int64                   `json:"expenseQuantity"`
	Term            FinancialCommitmentTerm `json:"term"`
}

func (p FinancialCommitmentPriceFact) Validate() error {
	if err := validateEventHappeningText("priceID", p.PriceID, EventHappeningIDMaxBytes, true); err != nil {
		return err
	}
	if p.Currency == "" || p.AmountMinor < 0 || p.ExpenseQuantity <= 0 || p.Term.Unit == "" || p.Term.Length < 1 {
		return fmt.Errorf("price fields are invalid")
	}
	return nil
}

type FinancialCommitmentOccurrenceFact struct {
	OccurrenceID                  string                           `json:"occurrenceID"`
	SlotID                        string                           `json:"slotID,omitempty"`
	ScheduledDate                 string                           `json:"scheduledDate"`
	EffectiveDate                 string                           `json:"effectiveDate"`
	ContactLinks                  []FinancialCommitmentContactLink `json:"contactLinks,omitempty"`
	CancellationFinancialEffect   string                           `json:"cancellationFinancialEffect,omitempty"`
	DateAdjustmentFinancialEffect string                           `json:"dateAdjustmentFinancialEffect,omitempty"`
	Diagnostics                   []string                         `json:"diagnostics,omitempty"`
}

type FinancialCommitmentSource interface {
	ListFinancialCommitments(context.Context, string, FinancialCommitmentQuery) (FinancialCommitmentPage, error)
}

func validMonthISO(value string) bool {
	parsed, err := time.Parse("2006-01", value)
	return err == nil && parsed.Year() >= 1900 && parsed.Year() <= 9999 && parsed.Format("2006-01") == value
}

func validFinancialEffect(value string) bool { return value == "" || value == "unknown" }

func validISODate(value string) bool {
	if len(value) != len("2006-01-02") {
		return false
	}
	parsed, err := time.Parse("2006-01-02", value)
	return err == nil && parsed.Year() >= 1900 && parsed.Year() <= 9999
}

func validOptionalISODate(value string) bool { return value == "" || validISODate(value) }

func validateUniqueIDs(field string, values []string) error {
	seen := make(map[string]struct{}, len(values))
	for i, value := range values {
		if err := validateEventHappeningText(fmt.Sprintf("%s[%d]", field, i), value, EventHappeningIDMaxBytes, true); err != nil {
			return err
		}
		if _, exists := seen[value]; exists {
			return fmt.Errorf("%s contains duplicate ID %q", field, value)
		}
		seen[value] = struct{}{}
	}
	return nil
}

func validateCommitmentContactLinks(field string, links []FinancialCommitmentContactLink) error {
	seen := make(map[string]struct{}, len(links))
	for i, link := range links {
		if err := validateEventHappeningText(fmt.Sprintf("%s[%d].contactID", field, i), link.ContactID, EventHappeningIDMaxBytes, true); err != nil {
			return err
		}
		if _, exists := seen[link.ContactID]; exists {
			return fmt.Errorf("%s contains duplicate contactID %q", field, link.ContactID)
		}
		seen[link.ContactID] = struct{}{}
		roles := make(map[string]struct{}, len(link.Roles))
		for j, role := range link.Roles {
			if strings.TrimSpace(role) == "" || strings.TrimSpace(role) != role {
				return fmt.Errorf("%s[%d].roles[%d] is invalid", field, i, j)
			}
			if _, exists := roles[role]; exists {
				return fmt.Errorf("%s[%d].roles contains duplicate role %q", field, i, role)
			}
			roles[role] = struct{}{}
		}
	}
	return nil
}
