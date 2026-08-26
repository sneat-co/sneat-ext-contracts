package dto4bookius

import (
	"errors"
	"testing"
	"time"
)

func validCompetitionReservation() CompetitionEntryReservation {
	expiresAt := time.Now().UTC().Add(time.Minute)
	return CompetitionEntryReservation{
		ID: "reservation-1", RequestID: "request-1", Target: CompetitionEntryTarget{ExtensionID: "competios", EventID: "event-1", TournamentID: "tournament-1", CompetitionID: "competition-1", TargetVersion: 1},
		ParticipantReference: "participant-1", EntryReference: "entry-1", BookingRevision: 1, State: ReservationHeld,
		AmountMinor: 2500, Currency: "EUR", OfferReference: "offer-1", OfferVersion: 1, OfferChecksum: "sha256:1111111111111111111111111111111111111111111111111111111111111111", PaymentState: CompetitionEntryPaymentNotStarted, CheckoutOperation: CompetitionEntryCheckoutNone, ConfirmationDelivery: CompetitionEntryDeliveryNone, RefundDelivery: CompetitionEntryDeliveryNone, ExpiresAt: &expiresAt,
	}
}

func TestCompetitionEntryReservationNeverAcceptsBrowserPriceInput(t *testing.T) {
	request := CompetitionEntryReservationRequest{RequestID: "request-1", Target: CompetitionEntryTarget{ExtensionID: "competios", EventID: "event-1", TournamentID: "tournament-1", CompetitionID: "competition-1", TargetVersion: 1}, ParticipantReference: "participant-1", EntryReference: "entry-1"}
	if err := ValidateCompetitionEntryReservationRequest(request); err != nil {
		t.Fatalf("valid price-free request: %v", err)
	}
	if err := ValidateCompetitionEntryReservation(validCompetitionReservation()); err != nil {
		t.Fatalf("valid server reservation: %v", err)
	}
	for name, reservation := range map[string]CompetitionEntryReservation{
		"negative amount": func() CompetitionEntryReservation {
			value := validCompetitionReservation()
			value.AmountMinor = -1
			return value
		}(),
		"lowercase currency": func() CompetitionEntryReservation {
			value := validCompetitionReservation()
			value.Currency = "eur"
			return value
		}(),
		"hold without expiry": func() CompetitionEntryReservation {
			value := validCompetitionReservation()
			value.ExpiresAt = nil
			return value
		}(),
		"missing booking revision": func() CompetitionEntryReservation {
			value := validCompetitionReservation()
			value.BookingRevision = 0
			return value
		}(),
	} {
		if err := ValidateCompetitionEntryReservation(reservation); !errors.Is(err, ErrInvalidCompetitionEntry) {
			t.Errorf("%s error = %v, want ErrInvalidCompetitionEntry", name, err)
		}
	}
}

func TestCompetitionEntryConfirmationEvidenceDistinguishesFreeFromSettlement(t *testing.T) {
	free := validCompetitionReservation()
	free.AmountMinor, free.State, free.PaymentState = 0, ReservationConfirmed, CompetitionEntryPaymentFree
	free.ConfirmationEvidence = CompetitionEntryConfirmationEvidence{Kind: CompetitionEntryConfirmationEvidenceFree}
	if err := ValidateCompetitionEntryReservation(free); err != nil {
		t.Fatalf("free confirmation: %v", err)
	}
	paid := validCompetitionReservation()
	paid.State, paid.PaymentState = ReservationConfirmed, CompetitionEntryPaymentPaid
	paid.ConfirmationEvidence = CompetitionEntryConfirmationEvidence{Kind: CompetitionEntryConfirmationEvidenceSettled, SettlementReference: "pi_1"}
	if err := ValidateCompetitionEntryReservation(paid); err != nil {
		t.Fatalf("settled confirmation: %v", err)
	}
	paid.ConfirmationEvidence = CompetitionEntryConfirmationEvidence{Kind: CompetitionEntryConfirmationEvidenceFree}
	if err := ValidateCompetitionEntryReservation(paid); !errors.Is(err, ErrInvalidCompetitionEntry) {
		t.Fatalf("ambiguous paid confirmation: %v", err)
	}
	settlement := SettlementNotification{SettlementID: "delivery-1", SettlementReference: "checkout_1", RefundReference: "pi_1", ReservationID: paid.ID, Status: SettlementPaid, AmountMinor: paid.AmountMinor, Currency: paid.Currency, OccurredAt: time.Now().UTC()}
	if err := ValidateSettlementNotification(settlement); err != nil {
		t.Fatalf("verified settlement: %v", err)
	}
	settlement.SettlementReference = ""
	if err := ValidateSettlementNotification(settlement); !errors.Is(err, ErrInvalidCompetitionEntry) {
		t.Fatalf("settlement without provider reference: %v", err)
	}
}

func TestCompetitionEntryLifecycleCommandsRequireDurableIdentity(t *testing.T) {
	validID := CompetitionEntryReservationID("reservation-1")
	if err := ValidateApproveCompetitionEntryReservationRequest(ApproveCompetitionEntryReservationRequest{CommandID: "approve-1", ReservationID: validID}); err != nil {
		t.Fatalf("approve validation: %v", err)
	}
	if err := ValidatePromoteCompetitionEntryWaitlistRequest(PromoteCompetitionEntryWaitlistRequest{CommandID: "promote-1", ReservationID: validID}); err != nil {
		t.Fatalf("promote validation: %v", err)
	}
	if err := ValidateExpireCompetitionEntryReservationRequest(ExpireCompetitionEntryReservationRequest{CommandID: "expire-1", ReservationID: validID}); err != nil {
		t.Fatalf("expire validation: %v", err)
	}
	if err := ValidateCancelCompetitionEntryReservationRequest(CancelCompetitionEntryReservationRequest{CommandID: "cancel-1", ReservationID: validID, Origin: CompetitionEntryCancellationOrganiser, ActorReference: "organiser-1", AuthorityEvidence: "auth-1", Reason: "medical evidence"}); err != nil {
		t.Fatalf("cancel validation: %v", err)
	}
	if err := ValidateCancelCompetitionEntryReservationRequest(CancelCompetitionEntryReservationRequest{}); !errors.Is(err, ErrInvalidCompetitionEntry) {
		t.Fatalf("empty command = %v, want ErrInvalidCompetitionEntry", err)
	}
}

func TestCompetitionEntryCancellationRequiresTrustedCurrentLockEvidence(t *testing.T) {
	request := CancelCompetitionEntryReservationRequest{CommandID: "cancel-1", ReservationID: "reservation-1", Origin: CompetitionEntryCancellationParticipant, ActorReference: "participant-1", AuthorityEvidence: "session-1", Reason: "cannot attend"}
	if err := ValidateCancelCompetitionEntryReservationRequest(request); err != nil {
		t.Fatalf("request: %v", err)
	}
	validation := CompetitionEntryCancellationValidation{Authorized: true, RefundAuthorized: true, CurrentTournamentVersion: 9, AuthorityEvidence: "competios-check-1", ValidatedAt: time.Now().UTC()}
	if err := ValidateCompetitionEntryCancellationValidation(validation, request.Origin); err != nil {
		t.Fatalf("unlocked participant: %v", err)
	}
	validation.RegistrationLocked = true
	validation.RefundAuthorized = false
	if err := ValidateCompetitionEntryCancellationValidation(validation, request.Origin); err != nil {
		t.Fatalf("locked participant no-refund cancellation: %v", err)
	}
	validation.RefundAuthorized = true
	if err := ValidateCompetitionEntryCancellationValidation(validation, request.Origin); !errors.Is(err, ErrInvalidCompetitionEntry) {
		t.Fatalf("locked participant refund without override: %v", err)
	}
	validation.AuthoriserReference = "organiser-1"
	if err := ValidateCompetitionEntryCancellationValidation(validation, request.Origin); err != nil {
		t.Fatalf("locked participant refund override: %v", err)
	}
	request.Origin = CompetitionEntryCancellationOrganiser
	if err := ValidateCompetitionEntryCancellationValidation(validation, request.Origin); err != nil {
		t.Fatalf("organiser override: %v", err)
	}
}
