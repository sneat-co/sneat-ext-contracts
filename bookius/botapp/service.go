// Package botapp defines the narrow Bookius capabilities consumed by bot
// delivery. The types here are presentation-safe projections, not persistence
// envelopes: callers never receive secret references or DALgo records.
package botapp

import (
	"context"
	"errors"
	"time"
)

// VendorBotView is the vendor-bot projection needed by the owner management
// UI. It intentionally omits persistence records and secret references.
type VendorBotView struct {
	Code               string
	Title              string
	Enabled            bool
	OwnerUserID        string
	PaymentsConfigured bool
}

// VendorManagementService supplies the owner-facing vendor-bot management
// capability to a delivery adapter.
type VendorManagementService interface {
	GetVendorBotView(context.Context, string) (VendorBotView, error)
	SetProviderToken(context.Context, string, string) error
}

// Invoice is the server-authoritative amount and description to present in a
// messenger invoice. It is deliberately independent of payment persistence.
type Invoice struct {
	CommitmentID string
	AmountCents  int64
	Currency     string
	Description  string
}

// Validate ensures that an invoice is safe for a delivery adapter to present
// to a payer. Amounts and descriptions always originate server-side.
func (i Invoice) Validate() error {
	switch {
	case i.CommitmentID == "":
		return errors.New("bookius bot invoice: commitment ID is required")
	case i.AmountCents <= 0:
		return errors.New("bookius bot invoice: amount must be positive")
	case i.Currency == "":
		return errors.New("bookius bot invoice: currency is required")
	case i.Description == "":
		return errors.New("bookius bot invoice: description is required")
	default:
		return nil
	}
}

// Settlement reports the delivery-relevant outcome of an idempotent payment
// settlement operation.
type Settlement struct {
	AlreadyPaid bool
}

// PaymentService supplies the three capabilities needed by the hosted vendor
// bot payment flow. Implementations own commitment storage and settlement.
type PaymentService interface {
	CommitmentIDFromStartRef(string) (string, error)
	LoadInvoice(context.Context, string) (Invoice, error)
	SettlePayment(context.Context, string, string, time.Time) (Settlement, error)
}
