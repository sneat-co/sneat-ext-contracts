import {
  BookiusBookingTypeBrief,
  BookiusCompetitionEntryReservationRequest,
  BookiusCancelCompetitionEntryReservationRequest,
  BookiusCompetitionEntryCancellationValidation,
  BookiusCompetitionEntryConfirmationEvidence,
  BookiusCreateBookingRequest,
} from './booking';
import { describe, expect, it } from 'vitest';

describe('Bookius booking DTOs', () => {
  it('supports generic booking type targets', () => {
    const bookingType: BookiusBookingTypeBrief = {
      id: 'office-meeting',
      title: 'Office Meeting',
      slug: 'office-meeting',
      durationMinutes: 60,
      targetKind: 'meeting-room',
      confirmationMode: 'request',
    };

    expect(bookingType.targetKind).toBe('meeting-room');
  });

  it('captures anonymous public booking requests', () => {
    const request: BookiusCreateBookingRequest = {
      bookingTypeID: 'investor-call',
      requestedSlot: {
        start: '2026-07-07T10:00:00Z',
        end: '2026-07-07T10:30:00Z',
        timezone: 'Europe/Dublin',
      },
      visitorName: 'Alex',
      visitorEmail: 'alex@example.com',
      subject: 'Sneat.co investment opportunity',
    };

    expect(request.visitorEmail).toBe('alex@example.com');
  });

  it('keeps price and currency out of a competition-entry browser request', () => {
    const request: BookiusCompetitionEntryReservationRequest = {
      requestID: 'request-1',
      target: {
        extensionID: 'competios',
        eventID: 'event-1',
        tournamentID: 'tournament-1',
        competitionID: 'competition-1',
        targetVersion: 1,
      },
      participantReference: 'participant-1',
      entryReference: 'entry-1',
    };

    expect(request.target.tournamentID).toBe('tournament-1');
    expect('amountMinor' in request).toBe(false);
    expect('currency' in request).toBe(false);
  });

  it('uses durable IDs for lifecycle commands and preserves an audited override reason', () => {
    const command: BookiusCancelCompetitionEntryReservationRequest = {
      commandID: 'cancel-1',
      reservationID: 'reservation-1',
      origin: 'organiser',
      actorReference: 'organiser-1',
      authorityEvidence: 'session-1',
      reason: 'organiser approved documented exception',
    };

    expect(command.commandID).toBe('cancel-1');
    expect(command.reason).toContain('approved');
  });

  it('keeps an explicit free confirmation distinct from a settled payment', () => {
    const evidence: BookiusCompetitionEntryConfirmationEvidence = {
      kind: 'settled',
      settlementReference: 'pi_1',
    };

    expect(evidence.settlementReference).toBe('pi_1');
  });

  it('represents a locked participant cancellation without granting a refund', () => {
    const validation: BookiusCompetitionEntryCancellationValidation = {
      authorized: true,
      refundAuthorized: false,
      currentTournamentVersion: 9,
      registrationLocked: true,
      authorityEvidence: 'competios-check-1',
      validatedAt: '2026-08-14T12:00:00Z',
    };

    expect(validation.refundAuthorized).toBe(false);
  });
});
