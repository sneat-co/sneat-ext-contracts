package dto4bookius

import (
	"context"
	"errors"
	"strings"
	"time"
)

const ExtensionID = "bookius"

type TargetKind string

const (
	TargetKindPerson       TargetKind = "person"
	TargetKindAppointment  TargetKind = "appointment"
	TargetKindConsultation TargetKind = "consultation"
	TargetKindMeetingRoom  TargetKind = "meeting-room"
	TargetKindFacility     TargetKind = "facility"
	TargetKindAsset        TargetKind = "asset"
	TargetKindEquipment    TargetKind = "equipment"
	TargetKindService      TargetKind = "service"
	TargetKindEventSession TargetKind = "event-session"
	// TargetKindCompetitionEntry is a paid entry into an externally-owned
	// competition. The provider-specific target identity remains opaque.
	TargetKindCompetitionEntry TargetKind = "competition-entry"
	TargetKindCustom           TargetKind = "custom"
)

type BookingState string

const (
	BookingStateDraft       BookingState = "draft"
	BookingStateHeld        BookingState = "held"
	BookingStateRequested   BookingState = "requested"
	BookingStateConfirmed   BookingState = "confirmed"
	BookingStateRescheduled BookingState = "rescheduled"
	BookingStateCancelled   BookingState = "cancelled"
	BookingStateExpired     BookingState = "expired"
	BookingStateCompleted   BookingState = "completed"
	BookingStateNoShow      BookingState = "no-show"
)

type ExtensionRef struct {
	Ext        string `json:"ext"`
	Collection string `json:"collection,omitempty"`
	ID         string `json:"id"`
}

type Slot struct {
	Start    string `json:"start"`
	End      string `json:"end"`
	Timezone string `json:"timezone,omitempty"`
}

type BookingType struct {
	ID                    string        `json:"id,omitempty"`
	Title                 string        `json:"title"`
	Slug                  string        `json:"slug"`
	Description           string        `json:"description,omitempty"`
	DurationMinutes       int           `json:"durationMinutes"`
	TargetKind            TargetKind    `json:"targetKind"`
	TargetRef             *ExtensionRef `json:"targetRef,omitempty"`
	AvailabilitySourceRef *ExtensionRef `json:"availabilitySourceRef,omitempty"`
	ConfirmationMode      string        `json:"confirmationMode"`
}

type CreateBookingRequest struct {
	SpaceID       string `json:"spaceID,omitempty"`
	BookingTypeID string `json:"bookingTypeID"`
	BookingPageID string `json:"bookingPageID,omitempty"`
	RequestedSlot Slot   `json:"requestedSlot"`
	VisitorName   string `json:"visitorName"`
	VisitorEmail  string `json:"visitorEmail"`
	VisitorPhone  string `json:"visitorPhone,omitempty"`
	Subject       string `json:"subject,omitempty"`
	Message       string `json:"message,omitempty"`
}

var ErrInvalidCompetitionEntry = errors.New("bookius: invalid competition entry value")

type CompetitionEntryReservationID string
type CheckoutID string
type SettlementID string

// CompetitionEntryTarget is an opaque external target. Bookius does not
// import Competios or assume a sport/game schema, but a Competition Entry must
// identify both an umbrella Event and a Tournament beneath it.
type CompetitionEntryTarget struct {
	ExtensionID   string `json:"extensionID"`
	EventID       string `json:"eventID"`
	TournamentID  string `json:"tournamentID"`
	CompetitionID string `json:"competitionID"`
	// TargetVersion binds a booking to the immutable external target revision.
	TargetVersion uint32 `json:"targetVersion"`
}

type ReservationState string

const (
	ReservationRequested  ReservationState = "requested"
	ReservationWaitlisted ReservationState = "waitlisted"
	ReservationHeld       ReservationState = "held"
	ReservationCheckout   ReservationState = "checkout-ready"
	ReservationConfirmed  ReservationState = "confirmed"
	ReservationFailed     ReservationState = "failed"
	ReservationExpired    ReservationState = "expired"
	ReservationCancelled  ReservationState = "cancelled"
	ReservationRefunded   ReservationState = "refunded"
)

// CompetitionEntryPaymentState is independent of booking state: for example,
// an active hold can be not_started while a confirmed booking is paid.
type CompetitionEntryPaymentState string

const (
	CompetitionEntryPaymentFree          CompetitionEntryPaymentState = "free"
	CompetitionEntryPaymentNotStarted    CompetitionEntryPaymentState = "not_started"
	CompetitionEntryPaymentCheckoutOpen  CompetitionEntryPaymentState = "checkout_open"
	CompetitionEntryPaymentPaid          CompetitionEntryPaymentState = "paid"
	CompetitionEntryPaymentRefundPending CompetitionEntryPaymentState = "refund_pending"
	CompetitionEntryPaymentRefunded      CompetitionEntryPaymentState = "refunded"
	CompetitionEntryPaymentFailed        CompetitionEntryPaymentState = "failed"
)

// CompetitionEntryReservationRequest intentionally contains no amount or
// currency. The target provider resolves the immutable offer server-to-server;
// browser input can only name the already-authorised participant/entry.
type CompetitionEntryReservationRequest struct {
	RequestID            string                 `json:"requestID"`
	Target               CompetitionEntryTarget `json:"target"`
	ParticipantReference string                 `json:"participantReference"`
	EntryReference       string                 `json:"entryReference"`
}

// CompetitionEntryReservation is Bookius' safe projection. Amount and
// currency are copied only from the server-authoritative offer/hold; callers
// must never construct this from a browser request.
type CompetitionEntryReservation struct {
	ID                   CompetitionEntryReservationID `json:"id"`
	RequestID            string                        `json:"requestID"`
	Target               CompetitionEntryTarget        `json:"target"`
	ParticipantReference string                        `json:"participantReference"`
	EntryReference       string                        `json:"entryReference"`
	// BookingRevision advances for each authoritative Bookius state change.
	BookingRevision uint32                       `json:"bookingRevision"`
	State           ReservationState             `json:"state"`
	AmountMinor     int64                        `json:"amountMinor"`
	Currency        string                       `json:"currency"`
	OfferReference  string                       `json:"offerReference"`
	OfferVersion    uint32                       `json:"offerVersion"`
	OfferChecksum   string                       `json:"offerChecksum"`
	PaymentState    CompetitionEntryPaymentState `json:"paymentState"`
	// ConfirmationEvidence is populated only when Bookius confirms the entry.
	// It makes a zero-price free confirmation distinguishable from a payment.
	ConfirmationEvidence CompetitionEntryConfirmationEvidence   `json:"confirmationEvidence,omitempty"`
	CheckoutOperation    CompetitionEntryCheckoutOperationState `json:"checkoutOperation"`
	ConfirmationDelivery CompetitionEntryDeliveryState          `json:"confirmationDelivery"`
	RefundDelivery       CompetitionEntryDeliveryState          `json:"refundDelivery"`
	ExpiresAt            *time.Time                             `json:"expiresAt,omitempty"`
}

type CompetitionEntryCheckoutOperationState string

const (
	CompetitionEntryCheckoutNone    CompetitionEntryCheckoutOperationState = "none"
	CompetitionEntryCheckoutPending CompetitionEntryCheckoutOperationState = "pending"
	CompetitionEntryCheckoutReady   CompetitionEntryCheckoutOperationState = "ready"
	CompetitionEntryCheckoutFailed  CompetitionEntryCheckoutOperationState = "failed"
)

type CompetitionEntryDeliveryState string

const (
	CompetitionEntryDeliveryNone      CompetitionEntryDeliveryState = "none"
	CompetitionEntryDeliveryPending   CompetitionEntryDeliveryState = "pending"
	CompetitionEntryDeliveryDelivered CompetitionEntryDeliveryState = "delivered"
	CompetitionEntryDeliveryFailed    CompetitionEntryDeliveryState = "failed"
)

type CheckoutProjection struct {
	CheckoutID    CheckoutID                    `json:"checkoutID"`
	ReservationID CompetitionEntryReservationID `json:"reservationID"`
	State         ReservationState              `json:"state"`
	CheckoutURL   string                        `json:"checkoutURL,omitempty"`
	ExpiresAt     time.Time                     `json:"expiresAt"`
}

type SettlementStatus string

const (
	SettlementPaid     SettlementStatus = "paid"
	SettlementFailed   SettlementStatus = "failed"
	SettlementRefunded SettlementStatus = "refunded"
)

// SettlementNotification is a server-to-server provider result. It repeats
// the held amount/currency so Bookius can fail closed on a mismatched payment;
// it is never browser-supplied input.
type SettlementNotification struct {
	SettlementID SettlementID `json:"settlementID"`
	// SettlementReference is the verified provider payment reference. It is
	// distinct from SettlementID, which only deduplicates webhook delivery.
	SettlementReference string `json:"settlementReference"`
	// RefundReference is the provider's opaque reversal handle. It can differ
	// from the checkout/settlement reference (for example Stripe Session vs
	// PaymentIntent) and is never used as payment-confirmation evidence.
	RefundReference string                        `json:"refundReference"`
	ReservationID   CompetitionEntryReservationID `json:"reservationID"`
	Status          SettlementStatus              `json:"status"`
	AmountMinor     int64                         `json:"amountMinor"`
	Currency        string                        `json:"currency"`
	OccurredAt      time.Time                     `json:"occurredAt"`
}

type CompetitionEntryConfirmationEvidenceKind string

const (
	CompetitionEntryConfirmationEvidenceFree    CompetitionEntryConfirmationEvidenceKind = "free"
	CompetitionEntryConfirmationEvidenceSettled CompetitionEntryConfirmationEvidenceKind = "settled"
)

// CompetitionEntryConfirmationEvidence is forwarded to the external Entry
// confirmation consumer. Free and settled confirmations are deliberately
// mutually exclusive, so a zero price is never mistaken for a missing payment.
type CompetitionEntryConfirmationEvidence struct {
	Kind                CompetitionEntryConfirmationEvidenceKind `json:"kind"`
	SettlementReference string                                   `json:"settlementReference,omitempty"`
}

// ApproveCompetitionEntryReservationRequest is an organiser command. Approval
// is deliberately separate from checkout: implementations create the hold
// first and only then may expose a payment link.
type ApproveCompetitionEntryReservationRequest struct {
	CommandID     string                        `json:"commandID"`
	ReservationID CompetitionEntryReservationID `json:"reservationID"`
}

// PromoteCompetitionEntryWaitlistRequest moves one existing waitlist record
// through the same approval/hold path. It never authorises a charge itself.
type PromoteCompetitionEntryWaitlistRequest struct {
	CommandID     string                        `json:"commandID"`
	ReservationID CompetitionEntryReservationID `json:"reservationID"`
}

// ExpireCompetitionEntryReservationRequest is the durable expiry command used
// by a scheduler. A replay is safe; an implementation must never expire an
// already confirmed booking.
type ExpireCompetitionEntryReservationRequest struct {
	CommandID     string                        `json:"commandID"`
	ReservationID CompetitionEntryReservationID `json:"reservationID"`
}

type CompetitionEntryCancellationOrigin string

const (
	CompetitionEntryCancellationParticipant CompetitionEntryCancellationOrigin = "participant"
	CompetitionEntryCancellationOrganiser   CompetitionEntryCancellationOrigin = "organiser"
)

// CancelCompetitionEntryReservationRequest is authenticated command evidence.
// AuthorityEvidence is minted by the protected delivery adapter; it is not a
// browser assertion. Current lock state is intentionally absent and must be
// resolved through CompetitionEntryCancellationValidator at command time.
type CancelCompetitionEntryReservationRequest struct {
	CommandID         string                             `json:"commandID"`
	ReservationID     CompetitionEntryReservationID      `json:"reservationID"`
	Origin            CompetitionEntryCancellationOrigin `json:"origin"`
	ActorReference    string                             `json:"actorReference"`
	AuthorityEvidence string                             `json:"authorityEvidence"`
	Reason            string                             `json:"reason"`
}

// CompetitionEntryCancellationValidation is trusted server-to-server output
// from the target owner. It binds the decision to the current Tournament
// version/lock state and records who authorised an organiser override.
type CompetitionEntryCancellationValidation struct {
	Authorized               bool      `json:"authorized"`
	RefundAuthorized         bool      `json:"refundAuthorized"`
	CurrentTournamentVersion uint32    `json:"currentTournamentVersion"`
	RegistrationLocked       bool      `json:"registrationLocked"`
	AuthoriserReference      string    `json:"authoriserReference,omitempty"`
	AuthorityEvidence        string    `json:"authorityEvidence"`
	ValidatedAt              time.Time `json:"validatedAt"`
}

type CompetitionEntryCancellationValidator interface {
	ValidateCompetitionEntryCancellation(context.Context, CompetitionEntryReservation, CancelCompetitionEntryReservationRequest) (CompetitionEntryCancellationValidation, error)
}

func ValidateCompetitionEntryReservationRequest(value CompetitionEntryReservationRequest) error {
	if strings.TrimSpace(value.RequestID) == "" || !validCompetitionEntryTarget(value.Target) || strings.TrimSpace(value.ParticipantReference) == "" || strings.TrimSpace(value.EntryReference) == "" {
		return ErrInvalidCompetitionEntry
	}
	return nil
}

func ValidateCompetitionEntryReservation(value CompetitionEntryReservation) error {
	if value.ID == "" || strings.TrimSpace(value.RequestID) == "" || !validCompetitionEntryTarget(value.Target) || strings.TrimSpace(value.ParticipantReference) == "" || strings.TrimSpace(value.EntryReference) == "" || value.BookingRevision == 0 || value.AmountMinor < 0 || !validISOCurrency(value.Currency) || strings.TrimSpace(value.OfferReference) == "" || value.OfferVersion == 0 || !validSHA256Checksum(value.OfferChecksum) {
		return ErrInvalidCompetitionEntry
	}
	switch value.State {
	case ReservationRequested, ReservationWaitlisted, ReservationHeld, ReservationCheckout, ReservationConfirmed, ReservationFailed, ReservationExpired, ReservationCancelled, ReservationRefunded:
	default:
		return ErrInvalidCompetitionEntry
	}
	if (value.State == ReservationHeld || value.State == ReservationCheckout) && (value.ExpiresAt == nil || value.ExpiresAt.IsZero()) {
		return ErrInvalidCompetitionEntry
	}
	switch value.PaymentState {
	case CompetitionEntryPaymentFree, CompetitionEntryPaymentNotStarted, CompetitionEntryPaymentCheckoutOpen, CompetitionEntryPaymentPaid, CompetitionEntryPaymentRefundPending, CompetitionEntryPaymentRefunded, CompetitionEntryPaymentFailed:
	default:
		return ErrInvalidCompetitionEntry
	}
	switch value.CheckoutOperation {
	case CompetitionEntryCheckoutNone, CompetitionEntryCheckoutPending, CompetitionEntryCheckoutReady, CompetitionEntryCheckoutFailed:
	default:
		return ErrInvalidCompetitionEntry
	}
	for _, delivery := range []CompetitionEntryDeliveryState{value.ConfirmationDelivery, value.RefundDelivery} {
		switch delivery {
		case CompetitionEntryDeliveryNone, CompetitionEntryDeliveryPending, CompetitionEntryDeliveryDelivered, CompetitionEntryDeliveryFailed:
		default:
			return ErrInvalidCompetitionEntry
		}
	}
	if value.State == ReservationConfirmed && !validCompetitionEntryConfirmationEvidence(value.ConfirmationEvidence, value.AmountMinor) {
		return ErrInvalidCompetitionEntry
	}
	return nil
}

func ValidateSettlementNotification(value SettlementNotification) error {
	if value.SettlementID == "" || strings.TrimSpace(value.SettlementReference) == "" || strings.TrimSpace(value.RefundReference) == "" || value.ReservationID == "" || value.AmountMinor < 0 || !validISOCurrency(value.Currency) || value.OccurredAt.IsZero() {
		return ErrInvalidCompetitionEntry
	}
	switch value.Status {
	case SettlementPaid, SettlementFailed, SettlementRefunded:
		return nil
	default:
		return ErrInvalidCompetitionEntry
	}
}

func ValidateApproveCompetitionEntryReservationRequest(value ApproveCompetitionEntryReservationRequest) error {
	if strings.TrimSpace(value.CommandID) == "" || value.ReservationID == "" {
		return ErrInvalidCompetitionEntry
	}
	return nil
}

func ValidatePromoteCompetitionEntryWaitlistRequest(value PromoteCompetitionEntryWaitlistRequest) error {
	if strings.TrimSpace(value.CommandID) == "" || value.ReservationID == "" {
		return ErrInvalidCompetitionEntry
	}
	return nil
}

func ValidateExpireCompetitionEntryReservationRequest(value ExpireCompetitionEntryReservationRequest) error {
	if strings.TrimSpace(value.CommandID) == "" || value.ReservationID == "" {
		return ErrInvalidCompetitionEntry
	}
	return nil
}

func ValidateCancelCompetitionEntryReservationRequest(value CancelCompetitionEntryReservationRequest) error {
	if strings.TrimSpace(value.CommandID) == "" || value.ReservationID == "" || strings.TrimSpace(value.ActorReference) == "" || strings.TrimSpace(value.AuthorityEvidence) == "" || strings.TrimSpace(value.Reason) == "" {
		return ErrInvalidCompetitionEntry
	}
	switch value.Origin {
	case CompetitionEntryCancellationParticipant, CompetitionEntryCancellationOrganiser:
		return nil
	default:
		return ErrInvalidCompetitionEntry
	}
}

func ValidateCompetitionEntryCancellationValidation(value CompetitionEntryCancellationValidation, origin CompetitionEntryCancellationOrigin) error {
	if !value.Authorized || value.CurrentTournamentVersion == 0 || strings.TrimSpace(value.AuthorityEvidence) == "" || value.ValidatedAt.IsZero() {
		return ErrInvalidCompetitionEntry
	}
	switch origin {
	case CompetitionEntryCancellationOrganiser:
		if !value.RefundAuthorized {
			return ErrInvalidCompetitionEntry
		}
	case CompetitionEntryCancellationParticipant:
		if !value.RegistrationLocked && !value.RefundAuthorized {
			return ErrInvalidCompetitionEntry
		}
		if value.RegistrationLocked && value.RefundAuthorized && strings.TrimSpace(value.AuthoriserReference) == "" {
			return ErrInvalidCompetitionEntry
		}
	default:
		return ErrInvalidCompetitionEntry
	}
	return nil
}

func validCompetitionEntryTarget(value CompetitionEntryTarget) bool {
	return value.ExtensionID != "" && value.EventID != "" && value.TournamentID != "" && value.CompetitionID != "" && value.TargetVersion > 0
}

func validCompetitionEntryConfirmationEvidence(value CompetitionEntryConfirmationEvidence, amountMinor int64) bool {
	switch value.Kind {
	case CompetitionEntryConfirmationEvidenceFree:
		return amountMinor == 0 && value.SettlementReference == ""
	case CompetitionEntryConfirmationEvidenceSettled:
		return amountMinor > 0 && strings.TrimSpace(value.SettlementReference) != ""
	default:
		return false
	}
}

func validISOCurrency(value string) bool {
	if len(value) != 3 {
		return false
	}
	for _, character := range value {
		if character < 'A' || character > 'Z' {
			return false
		}
	}
	return true
}

func validSHA256Checksum(value string) bool {
	if len(value) != len("sha256:")+64 || !strings.HasPrefix(value, "sha256:") {
		return false
	}
	for _, character := range value[len("sha256:"):] {
		if (character < '0' || character > '9') && (character < 'a' || character > 'f') {
			return false
		}
	}
	return true
}

// CompetitionEntryReservations is the narrow Bookius server port for paid
// competition participation. A provider resolves price and capacity when a
// hold is created; the browser cannot send amount/currency or settlement facts.
// Implementations must deduplicate request IDs and settlement IDs.
type CompetitionEntryReservations interface {
	ReserveCompetitionEntry(context.Context, CompetitionEntryReservationRequest) (CompetitionEntryReservation, error)
	ApproveCompetitionEntryReservation(context.Context, ApproveCompetitionEntryReservationRequest) (CompetitionEntryReservation, error)
	PromoteCompetitionEntryWaitlist(context.Context, PromoteCompetitionEntryWaitlistRequest) (CompetitionEntryReservation, error)
	BeginCompetitionEntryCheckout(context.Context, CompetitionEntryReservationID) (CheckoutProjection, error)
	RecordCompetitionEntrySettlement(context.Context, SettlementNotification) (CompetitionEntryReservation, error)
	ExpireCompetitionEntryReservation(context.Context, ExpireCompetitionEntryReservationRequest) (CompetitionEntryReservation, error)
	CancelCompetitionEntryReservation(context.Context, CancelCompetitionEntryReservationRequest) (CompetitionEntryReservation, error)
}
