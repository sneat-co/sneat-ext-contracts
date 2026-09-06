// Copyright (c) 2026 Sneat.co

package calendariusmodels

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

const (
	MaxFinancialOccurrenceWindowDays = 93
	MaxFinancialOccurrenceSources    = 64
	MaxFinancialOccurrenceResults    = 256
	maxFinancialSpaceIDBytes         = 30
	MaxFinancialOccurrenceIDBytes    = 256
)

type FinancialPricingAvailability string

const (
	FinancialPricingAvailable   FinancialPricingAvailability = "available"
	FinancialPricingAmbiguous   FinancialPricingAvailability = "ambiguous"
	FinancialPricingUnsupported FinancialPricingAvailability = "unsupported"
)

type CancellationFinancialEffect string

const CancellationFinancialEffectUnknown CancellationFinancialEffect = "unknown"

var (
	ErrFinancialOccurrenceNotFound      = errors.New("financial occurrence not found")
	ErrFinancialOccurrenceUnauthorized  = errors.New("financial occurrence is not authorized")
	ErrFinancialOccurrenceQueryTooBroad = errors.New("financial occurrence query is too broad")
)

type FinancialOccurrenceQuery struct {
	SpaceID     string `json:"spaceID"`
	HappeningID string `json:"happeningID"`
	FromDate    string `json:"fromDate"`
	ToDate      string `json:"toDate"`
	Timezone    string `json:"timezone,omitempty"`
}

func (v FinancialOccurrenceQuery) Validate() error {
	if err := validateFinancialSpaceID(v.SpaceID); err != nil {
		return err
	}
	if v.HappeningID != "" {
		if err := ValidateEventHappeningID(v.HappeningID); err != nil {
			return err
		}
		if v.HappeningID == "." || v.HappeningID == ".." || strings.ContainsAny(v.HappeningID, "/\\") || containsControl(v.HappeningID) {
			return errors.New("happeningID is not safe for a same-Space storage key")
		}
	}
	from, err := time.Parse("2006-01-02", v.FromDate)
	if err != nil || from.Format("2006-01-02") != v.FromDate {
		return errors.New("fromDate must be a real ISO date")
	}
	to, err := time.Parse("2006-01-02", v.ToDate)
	if err != nil || to.Format("2006-01-02") != v.ToDate || to.Before(from) {
		return errors.New("toDate must be a real ISO date on or after fromDate")
	}
	if days := int(to.Sub(from).Hours()/24) + 1; days > MaxFinancialOccurrenceWindowDays {
		return fmt.Errorf("financial occurrence window exceeds %d days", MaxFinancialOccurrenceWindowDays)
	}
	if v.Timezone != "" {
		if _, err = time.LoadLocation(v.Timezone); err != nil {
			return errors.New("timezone must be a valid IANA location")
		}
	}
	return nil
}

// ValidateForResolve requires the source identity used by an atomic consumer
// mutation. List queries may omit it to discover bounded same-Space candidates.
func (v FinancialOccurrenceQuery) ValidateForResolve() error {
	if err := v.Validate(); err != nil {
		return err
	}
	if v.HappeningID == "" {
		return errors.New("happeningID is required to resolve an occurrence")
	}
	return nil
}

type FinancialOccurrenceCandidate struct {
	HappeningID  string `json:"happeningID"`
	OccurrenceID string `json:"occurrenceID"`
	SlotID       string `json:"slotID,omitempty"`
	// ScheduledDate is the calendar occurrence date. It is not an invoice,
	// payment, posting, or cash-due date. Monthly economic costs omit it.
	ScheduledDate string `json:"scheduledDate,omitempty"`
	// PeriodStartDate/PeriodEndDate describe the economic source period. They
	// are never inferred invoice, due, posting, or payment dates.
	PeriodStartDate string `json:"periodStartDate"`
	PeriodEndDate   string `json:"periodEndDate"`
	Title           string `json:"title"`
	PriceID         string `json:"priceID,omitempty"`
	PriceRevision   int64  `json:"priceRevision,omitempty"`
	Currency        string `json:"currency,omitempty"`
	ExpectedMinor   *int64 `json:"expectedMinor,omitempty"`
	// AssetIDs are same-Space Assetus records linked by the source Happening.
	// They preserve attribution context; they do not change the expected amount.
	AssetIDs                 []string                     `json:"assetIDs,omitempty"`
	PricingAvailability      FinancialPricingAvailability `json:"pricingAvailability"`
	PricingUnavailableReason string                       `json:"pricingUnavailableReason,omitempty"`
	// CancellationFinancialEffect=unknown means the source occurrence is known
	// to be canceled while whether its charge was waived remains unresolved.
	CancellationFinancialEffect CancellationFinancialEffect `json:"cancellationFinancialEffect,omitempty"`
}

func ValidateFinancialOccurrenceID(value string) error {
	if value == "" || len(value) > MaxFinancialOccurrenceIDBytes || !utf8.ValidString(value) || strings.TrimSpace(value) != value || strings.ContainsAny(value, "/\\") || value == "." || value == ".." || containsControl(value) {
		return fmt.Errorf("occurrenceID must be a safe UTF-8 identifier with at most %d bytes", MaxFinancialOccurrenceIDBytes)
	}
	return nil
}

func validateFinancialSpaceID(value string) error {
	if value == "" || len(value) > maxFinancialSpaceIDBytes || !utf8.ValidString(value) {
		return errors.New("spaceID is required and must be at most 30 UTF-8 bytes")
	}
	for byteIndex, r := range value {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || (r == '~' && byteIndex == len("spot")) {
			continue
		}
		return errors.New("spaceID contains unsupported characters")
	}
	if strings.Contains(value, "~") && !strings.HasPrefix(value, "spot~") {
		return errors.New("spaceID contains unsupported characters")
	}
	return nil
}

func containsControl(value string) bool {
	for _, r := range value {
		if unicode.IsControl(r) {
			return true
		}
	}
	return false
}

type FinancialOccurrenceSource interface {
	// ListFinancialOccurrences returns verified source facts. A future AI or
	// automation caller must use this same authenticated port and must present
	// proposals separately from these verified facts; it must not bypass source
	// provenance, replay, or audit controls.
	// OccurrenceID is opaque to consumers and must come from this result; clients
	// must never construct or infer it from a slot or date.
	ListFinancialOccurrences(ctx context.Context, userID string, query FinancialOccurrenceQuery) ([]FinancialOccurrenceCandidate, error)
}
