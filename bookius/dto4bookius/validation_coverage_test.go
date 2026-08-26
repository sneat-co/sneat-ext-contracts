package dto4bookius

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func requireInvalidCompetitionEntry(t *testing.T, err error) {
	t.Helper()
	if !errors.Is(err, ErrInvalidCompetitionEntry) {
		t.Fatalf("error = %v, want ErrInvalidCompetitionEntry", err)
	}
}

func TestReservationRequestValidationPaths(t *testing.T) {
	valid := CompetitionEntryReservationRequest{RequestID: "request", Target: CompetitionEntryTarget{ExtensionID: "competios", EventID: "event", TournamentID: "tournament", CompetitionID: "competition", TargetVersion: 1}, ParticipantReference: "participant", EntryReference: "entry"}
	mutations := []func(*CompetitionEntryReservationRequest){
		func(v *CompetitionEntryReservationRequest) { v.RequestID = " " },
		func(v *CompetitionEntryReservationRequest) { v.Target.ExtensionID = "" },
		func(v *CompetitionEntryReservationRequest) { v.Target.EventID = "" },
		func(v *CompetitionEntryReservationRequest) { v.Target.TournamentID = "" },
		func(v *CompetitionEntryReservationRequest) { v.Target.CompetitionID = "" },
		func(v *CompetitionEntryReservationRequest) { v.Target.TargetVersion = 0 },
		func(v *CompetitionEntryReservationRequest) { v.ParticipantReference = " " },
		func(v *CompetitionEntryReservationRequest) { v.EntryReference = " " },
	}
	for i, mutate := range mutations {
		value := valid
		mutate(&value)
		t.Run(string(rune('A'+i)), func(t *testing.T) {
			requireInvalidCompetitionEntry(t, ValidateCompetitionEntryReservationRequest(value))
		})
	}
}

func TestReservationValidationPaths(t *testing.T) {
	valid := validCompetitionReservation()
	invalidBase := []func(*CompetitionEntryReservation){
		func(v *CompetitionEntryReservation) { v.ID = "" }, func(v *CompetitionEntryReservation) { v.RequestID = " " },
		func(v *CompetitionEntryReservation) { v.Target.TargetVersion = 0 }, func(v *CompetitionEntryReservation) { v.ParticipantReference = " " },
		func(v *CompetitionEntryReservation) { v.EntryReference = " " }, func(v *CompetitionEntryReservation) { v.BookingRevision = 0 },
		func(v *CompetitionEntryReservation) { v.AmountMinor = -1 }, func(v *CompetitionEntryReservation) { v.Currency = "EU" },
		func(v *CompetitionEntryReservation) { v.Currency = "E1R" }, func(v *CompetitionEntryReservation) { v.OfferReference = " " },
		func(v *CompetitionEntryReservation) { v.OfferVersion = 0 }, func(v *CompetitionEntryReservation) { v.OfferChecksum = "bad" },
		func(v *CompetitionEntryReservation) { v.OfferChecksum = "sha256:" + strings.Repeat("g", 64) },
	}
	for i, mutate := range invalidBase {
		value := valid
		mutate(&value)
		t.Run("base-"+string(rune('A'+i)), func(t *testing.T) { requireInvalidCompetitionEntry(t, ValidateCompetitionEntryReservation(value)) })
	}
	for _, state := range []ReservationState{ReservationRequested, ReservationWaitlisted, ReservationHeld, ReservationCheckout, ReservationConfirmed, ReservationFailed, ReservationExpired, ReservationCancelled, ReservationRefunded} {
		value := valid
		value.State = state
		if state == ReservationConfirmed {
			value.ConfirmationEvidence = CompetitionEntryConfirmationEvidence{Kind: CompetitionEntryConfirmationEvidenceSettled, SettlementReference: "checkout"}
		}
		if err := ValidateCompetitionEntryReservation(value); err != nil {
			t.Fatalf("state %q: %v", state, err)
		}
	}
	value := valid
	value.State = "unknown"
	requireInvalidCompetitionEntry(t, ValidateCompetitionEntryReservation(value))
	value = valid
	value.State = ReservationCheckout
	value.ExpiresAt = nil
	requireInvalidCompetitionEntry(t, ValidateCompetitionEntryReservation(value))
	for _, payment := range []CompetitionEntryPaymentState{CompetitionEntryPaymentFree, CompetitionEntryPaymentNotStarted, CompetitionEntryPaymentCheckoutOpen, CompetitionEntryPaymentPaid, CompetitionEntryPaymentRefundPending, CompetitionEntryPaymentRefunded, CompetitionEntryPaymentFailed} {
		value = valid
		value.PaymentState = payment
		if err := ValidateCompetitionEntryReservation(value); err != nil {
			t.Fatalf("payment %q: %v", payment, err)
		}
	}
	value = valid
	value.PaymentState = "unknown"
	requireInvalidCompetitionEntry(t, ValidateCompetitionEntryReservation(value))
	for _, checkout := range []CompetitionEntryCheckoutOperationState{CompetitionEntryCheckoutNone, CompetitionEntryCheckoutPending, CompetitionEntryCheckoutReady, CompetitionEntryCheckoutFailed} {
		value = valid
		value.CheckoutOperation = checkout
		if err := ValidateCompetitionEntryReservation(value); err != nil {
			t.Fatalf("checkout %q: %v", checkout, err)
		}
	}
	value = valid
	value.CheckoutOperation = "unknown"
	requireInvalidCompetitionEntry(t, ValidateCompetitionEntryReservation(value))
	for _, delivery := range []CompetitionEntryDeliveryState{CompetitionEntryDeliveryNone, CompetitionEntryDeliveryPending, CompetitionEntryDeliveryDelivered, CompetitionEntryDeliveryFailed} {
		value = valid
		value.ConfirmationDelivery, value.RefundDelivery = delivery, delivery
		if err := ValidateCompetitionEntryReservation(value); err != nil {
			t.Fatalf("delivery %q: %v", delivery, err)
		}
	}
	value = valid
	value.ConfirmationDelivery = "unknown"
	requireInvalidCompetitionEntry(t, ValidateCompetitionEntryReservation(value))
	value = valid
	value.RefundDelivery = "unknown"
	requireInvalidCompetitionEntry(t, ValidateCompetitionEntryReservation(value))
	value = valid
	value.State = ReservationConfirmed
	value.ConfirmationEvidence = CompetitionEntryConfirmationEvidence{}
	requireInvalidCompetitionEntry(t, ValidateCompetitionEntryReservation(value))
	value = valid
	value.State = ReservationConfirmed
	value.AmountMinor = 0
	value.ConfirmationEvidence = CompetitionEntryConfirmationEvidence{Kind: CompetitionEntryConfirmationEvidenceFree, SettlementReference: "unexpected"}
	requireInvalidCompetitionEntry(t, ValidateCompetitionEntryReservation(value))
	value = valid
	value.State = ReservationConfirmed
	value.ConfirmationEvidence = CompetitionEntryConfirmationEvidence{Kind: CompetitionEntryConfirmationEvidenceSettled}
	requireInvalidCompetitionEntry(t, ValidateCompetitionEntryReservation(value))
}

func TestSettlementAndLifecycleValidationPaths(t *testing.T) {
	valid := SettlementNotification{SettlementID: "delivery", SettlementReference: "checkout", RefundReference: "refund", ReservationID: "reservation", Status: SettlementPaid, AmountMinor: 1, Currency: "EUR", OccurredAt: time.Now().UTC()}
	mutations := []func(*SettlementNotification){func(v *SettlementNotification) { v.SettlementID = "" }, func(v *SettlementNotification) { v.SettlementReference = " " }, func(v *SettlementNotification) { v.RefundReference = " " }, func(v *SettlementNotification) { v.ReservationID = "" }, func(v *SettlementNotification) { v.AmountMinor = -1 }, func(v *SettlementNotification) { v.Currency = "eur" }, func(v *SettlementNotification) { v.OccurredAt = time.Time{} }}
	for i, mutate := range mutations {
		value := valid
		mutate(&value)
		t.Run("settlement-"+string(rune('A'+i)), func(t *testing.T) { requireInvalidCompetitionEntry(t, ValidateSettlementNotification(value)) })
	}
	for _, status := range []SettlementStatus{SettlementPaid, SettlementFailed, SettlementRefunded} {
		value := valid
		value.Status = status
		if err := ValidateSettlementNotification(value); err != nil {
			t.Fatalf("status %q: %v", status, err)
		}
	}
	value := valid
	value.Status = "unknown"
	requireInvalidCompetitionEntry(t, ValidateSettlementNotification(value))
	reservationID := CompetitionEntryReservationID("reservation")
	requireInvalidCompetitionEntry(t, ValidateApproveCompetitionEntryReservationRequest(ApproveCompetitionEntryReservationRequest{ReservationID: reservationID}))
	requireInvalidCompetitionEntry(t, ValidateApproveCompetitionEntryReservationRequest(ApproveCompetitionEntryReservationRequest{CommandID: "command"}))
	requireInvalidCompetitionEntry(t, ValidatePromoteCompetitionEntryWaitlistRequest(PromoteCompetitionEntryWaitlistRequest{ReservationID: reservationID}))
	requireInvalidCompetitionEntry(t, ValidatePromoteCompetitionEntryWaitlistRequest(PromoteCompetitionEntryWaitlistRequest{CommandID: "command"}))
	requireInvalidCompetitionEntry(t, ValidateExpireCompetitionEntryReservationRequest(ExpireCompetitionEntryReservationRequest{ReservationID: reservationID}))
	requireInvalidCompetitionEntry(t, ValidateExpireCompetitionEntryReservationRequest(ExpireCompetitionEntryReservationRequest{CommandID: "command"}))
}

func TestCancellationValidationPaths(t *testing.T) {
	valid := CancelCompetitionEntryReservationRequest{CommandID: "command", ReservationID: "reservation", Origin: CompetitionEntryCancellationParticipant, ActorReference: "actor", AuthorityEvidence: "session", Reason: "reason"}
	mutations := []func(*CancelCompetitionEntryReservationRequest){func(v *CancelCompetitionEntryReservationRequest) { v.CommandID = " " }, func(v *CancelCompetitionEntryReservationRequest) { v.ReservationID = "" }, func(v *CancelCompetitionEntryReservationRequest) { v.ActorReference = " " }, func(v *CancelCompetitionEntryReservationRequest) { v.AuthorityEvidence = " " }, func(v *CancelCompetitionEntryReservationRequest) { v.Reason = " " }, func(v *CancelCompetitionEntryReservationRequest) { v.Origin = "unknown" }}
	for i, mutate := range mutations {
		value := valid
		mutate(&value)
		t.Run("request-"+string(rune('A'+i)), func(t *testing.T) {
			requireInvalidCompetitionEntry(t, ValidateCancelCompetitionEntryReservationRequest(value))
		})
	}
	validation := CompetitionEntryCancellationValidation{Authorized: true, RefundAuthorized: true, CurrentTournamentVersion: 1, AuthorityEvidence: "trusted", ValidatedAt: time.Now().UTC()}
	invalid := []func(*CompetitionEntryCancellationValidation){func(v *CompetitionEntryCancellationValidation) { v.Authorized = false }, func(v *CompetitionEntryCancellationValidation) { v.CurrentTournamentVersion = 0 }, func(v *CompetitionEntryCancellationValidation) { v.AuthorityEvidence = " " }, func(v *CompetitionEntryCancellationValidation) { v.ValidatedAt = time.Time{} }}
	for i, mutate := range invalid {
		value := validation
		mutate(&value)
		t.Run("validation-"+string(rune('A'+i)), func(t *testing.T) {
			requireInvalidCompetitionEntry(t, ValidateCompetitionEntryCancellationValidation(value, CompetitionEntryCancellationParticipant))
		})
	}
	value := validation
	value.RefundAuthorized = false
	requireInvalidCompetitionEntry(t, ValidateCompetitionEntryCancellationValidation(value, CompetitionEntryCancellationOrganiser))
	if err := ValidateCompetitionEntryCancellationValidation(validation, CompetitionEntryCancellationOrganiser); err != nil {
		t.Fatalf("organiser: %v", err)
	}
	if err := ValidateCompetitionEntryCancellationValidation(validation, CompetitionEntryCancellationParticipant); err != nil {
		t.Fatalf("unlocked participant: %v", err)
	}
	value = validation
	value.RefundAuthorized = false
	requireInvalidCompetitionEntry(t, ValidateCompetitionEntryCancellationValidation(value, CompetitionEntryCancellationParticipant))
	value = validation
	value.RegistrationLocked = true
	value.RefundAuthorized = false
	if err := ValidateCompetitionEntryCancellationValidation(value, CompetitionEntryCancellationParticipant); err != nil {
		t.Fatalf("locked no-refund participant: %v", err)
	}
	value.RefundAuthorized = true
	requireInvalidCompetitionEntry(t, ValidateCompetitionEntryCancellationValidation(value, CompetitionEntryCancellationParticipant))
	value.AuthoriserReference = "organiser"
	if err := ValidateCompetitionEntryCancellationValidation(value, CompetitionEntryCancellationParticipant); err != nil {
		t.Fatalf("locked participant override: %v", err)
	}
	requireInvalidCompetitionEntry(t, ValidateCompetitionEntryCancellationValidation(validation, "unknown"))
}
