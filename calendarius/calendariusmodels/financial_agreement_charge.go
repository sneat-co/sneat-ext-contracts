package calendariusmodels

import (
	"context"
	"fmt"
	"time"
)

const (
	MaxFinancialAgreementChargePageSize = 256
	FinancialChargeDirectionExpense     = "expense"
	FinancialChargeDirectionIncome      = "income"
	FinancialChargeDirectionTransfer    = "transfer"
	FinancialChargeStatusAvailable      = "available"
	FinancialChargeStatusUnavailable    = "unavailable"
	FinancialChargeBasisOccurrence      = "occurrence_service_cost"
	FinancialChargeBasisContractPeriod  = "contract_period_cost"
	FinancialChargeBasisNormalized      = "normalized_service_cost"
	FinancialChargeBillingUnknown       = "unknown"
	FinancialChargeBillingKnownDue      = "known_due"
)

// FinancialAgreementChargeQuery asks the source owner for economic effects in
// whole calendar months. It does not ask the consumer to expand recurrence.
type FinancialAgreementChargeQuery struct {
	ReportingSpaceID string `json:"reportingSpaceID"`
	FromMonthISO     string `json:"fromMonthISO"`
	Months           int    `json:"months"`
	PageSize         int    `json:"pageSize"`
	Cursor           string `json:"cursor,omitempty"`
}

func (q FinancialAgreementChargeQuery) Validate() error {
	if err := validateFinancialCommitmentID("reportingSpaceID", q.ReportingSpaceID); err != nil {
		return err
	}
	from, err := time.Parse("2006-01", q.FromMonthISO)
	if err != nil || from.Year() < 1900 || q.Months < 1 || q.Months > 12 || from.AddDate(0, q.Months, 0).Year() > 9999 {
		return fmt.Errorf("financial agreement charge window is invalid")
	}
	if q.PageSize < 1 || q.PageSize > MaxFinancialAgreementChargePageSize || len(q.Cursor) > MaxFinancialAgreementCursorBytes {
		return fmt.Errorf("financial agreement charge page is invalid")
	}
	return nil
}

type FinancialChargePeriod struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

type FinancialAgreementChargeFact struct {
	OwnerSpaceID                  string                    `json:"ownerSpaceID"`
	ReportingSpaceID              string                    `json:"reportingSpaceID"`
	AgreementID                   string                    `json:"agreementID"`
	EnrollmentID                  string                    `json:"enrollmentID"`
	EnrollmentScope               *FinancialEnrollmentScope `json:"enrollmentScope,omitempty"`
	HappeningID                   string                    `json:"happeningID"`
	Title                         string                    `json:"title"`
	Regular                       bool                      `json:"regular"`
	AgreementRevision             int64                     `json:"agreementRevision"`
	TermsRevision                 int64                     `json:"termsRevision"`
	ChargeID                      string                    `json:"chargeID"`
	OccurrenceID                  string                    `json:"occurrenceID,omitempty"`
	Direction                     string                    `json:"direction,omitempty"`
	AmountMinor                   *int64                    `json:"amountMinor,omitempty"`
	Currency                      string                    `json:"currency"`
	EconomicPeriod                FinancialChargePeriod     `json:"economicPeriod"`
	TemporalBasis                 string                    `json:"temporalBasis"`
	BillingTiming                 string                    `json:"billingTiming"`
	DueDate                       string                    `json:"dueDate,omitempty"`
	OwnerTimezone                 string                    `json:"ownerTimezone"`
	InvoiceReconciliationEligible bool                      `json:"invoiceReconciliationEligible"`
	AcceptedPrice                 AcceptedPriceTerms        `json:"acceptedPrice"`
	ContactAttributions           []FinancialAttribution    `json:"contactAttributions"`
	AssetAttributions             []FinancialAttribution    `json:"assetAttributions"`
	Status                        string                    `json:"status"`
	Diagnostics                   []string                  `json:"diagnostics"`
}

func (f FinancialAgreementChargeFact) Validate() error {
	for name, value := range map[string]string{"ownerSpaceID": f.OwnerSpaceID, "reportingSpaceID": f.ReportingSpaceID, "agreementID": f.AgreementID, "enrollmentID": f.EnrollmentID, "happeningID": f.HappeningID, "chargeID": f.ChargeID} {
		if err := validateFinancialCommitmentID(name, value); err != nil {
			return err
		}
	}
	if f.AgreementRevision < 1 || f.TermsRevision < 1 || f.TermsRevision > f.AgreementRevision || len(f.Title) > 100 || f.ContactAttributions == nil || f.AssetAttributions == nil || f.Diagnostics == nil {
		return fmt.Errorf("financial agreement charge metadata is invalid")
	}
	if err := f.AcceptedPrice.Validate(); err != nil {
		return err
	}
	if f.EnrollmentScope != nil {
		if err := f.EnrollmentScope.Validate(); err != nil {
			return fmt.Errorf("financial agreement enrollment scope: %w", err)
		}
	}
	if f.Currency != f.AcceptedPrice.Currency || !validISODate(f.EconomicPeriod.StartDate) || !validISODate(f.EconomicPeriod.EndDate) || f.EconomicPeriod.EndDate < f.EconomicPeriod.StartDate {
		return fmt.Errorf("financial agreement charge currency or period is invalid")
	}
	if f.OwnerTimezone == "" {
		return fmt.Errorf("financial agreement owner timezone is required")
	}
	if _, err := time.LoadLocation(f.OwnerTimezone); err != nil {
		return fmt.Errorf("financial agreement owner timezone is invalid")
	}
	if f.BillingTiming != FinancialChargeBillingUnknown && f.BillingTiming != FinancialChargeBillingKnownDue {
		return fmt.Errorf("financial agreement billing timing is invalid")
	}
	if f.BillingTiming == FinancialChargeBillingKnownDue {
		if !validISODate(f.DueDate) {
			return fmt.Errorf("known billing timing requires dueDate")
		}
	} else if f.DueDate != "" {
		return fmt.Errorf("unknown billing timing cannot assert dueDate")
	}
	switch f.TemporalBasis {
	case FinancialChargeBasisOccurrence, FinancialChargeBasisContractPeriod:
	case FinancialChargeBasisNormalized:
		if f.InvoiceReconciliationEligible {
			return fmt.Errorf("normalized service cost cannot reconcile an invoice")
		}
	default:
		return fmt.Errorf("financial agreement temporal basis is invalid")
	}
	if f.Status == FinancialChargeStatusAvailable {
		if f.EnrollmentScope == nil || f.Title == "" || f.AmountMinor == nil || *f.AmountMinor < 0 || *f.AmountMinor > MaxJavaScriptSafeInteger || (*f.AmountMinor == 0 && f.TemporalBasis != FinancialChargeBasisNormalized) || (f.Direction != FinancialChargeDirectionExpense && f.Direction != FinancialChargeDirectionIncome && f.Direction != FinancialChargeDirectionTransfer) {
			return fmt.Errorf("available financial agreement charge is invalid")
		}
		if err := validateAttributions(*f.AmountMinor, "contactAttributions", f.ContactAttributions); err != nil {
			return err
		}
		if err := validateAttributions(*f.AmountMinor, "assetAttributions", f.AssetAttributions); err != nil {
			return err
		}
	} else if f.Status != FinancialChargeStatusUnavailable || f.AmountMinor != nil || f.Direction != "" || f.InvoiceReconciliationEligible || len(f.ContactAttributions) != 0 || len(f.AssetAttributions) != 0 || len(f.Diagnostics) == 0 {
		return fmt.Errorf("unavailable financial agreement charge is invalid")
	}
	return nil
}

type FinancialAgreementChargePage struct {
	Charges            []FinancialAgreementChargeFact `json:"charges"`
	HasMore            bool                           `json:"hasMore"`
	NextCursor         string                         `json:"nextCursor,omitempty"`
	SnapshotConsistent bool                           `json:"snapshotConsistent"`
	IncompleteReason   string                         `json:"incompleteReason,omitempty"`
}

func (p FinancialAgreementChargePage) Validate() error {
	if p.Charges == nil || len(p.Charges) > MaxFinancialAgreementChargePageSize || p.HasMore != (p.NextCursor != "") || p.HasMore && p.SnapshotConsistent || !p.SnapshotConsistent && p.IncompleteReason == "" {
		return fmt.Errorf("financial agreement charge page metadata is invalid")
	}
	for i, charge := range p.Charges {
		if err := charge.Validate(); err != nil {
			return fmt.Errorf("charges[%d]: %w", i, err)
		}
	}
	return nil
}

type FinancialAgreementChargeSource interface {
	ListFinancialAgreementCharges(ctx context.Context, actorUserID string, query FinancialAgreementChargeQuery) (FinancialAgreementChargePage, error)
}
