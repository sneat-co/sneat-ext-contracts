import { IWithCreated } from '@sneat/dto';

export type BookiusTargetKind =
  | 'person'
  | 'appointment'
  | 'consultation'
  | 'meeting-room'
  | 'facility'
  | 'asset'
  | 'equipment'
  | 'service'
  | 'event-session'
  | 'competition-entry'
  | 'custom';

export type BookiusConfirmationMode = 'automatic' | 'manual' | 'request';

export type BookiusBookingState =
  | 'draft'
  | 'held'
  | 'requested'
  | 'confirmed'
  | 'rescheduled'
  | 'cancelled'
  | 'expired'
  | 'completed'
  | 'no-show';

export interface BookiusExtensionRef {
  readonly ext: string;
  readonly collection?: string;
  readonly id: string;
}

export interface BookiusPublicPath {
  readonly handle?: string;
  readonly businessSlug?: string;
  readonly category?: string;
  readonly resourceSlug?: string;
}

export interface BookiusAvailabilitySourceRef extends BookiusExtensionRef {
  readonly ext: 'calendarius';
}

export interface BookiusContactRef extends BookiusExtensionRef {
  readonly ext: 'contactus';
}

export interface BookiusLocation {
  readonly mode: 'online' | 'phone' | 'physical' | 'ask-host' | 'ask-visitor';
  readonly label?: string;
  readonly url?: string;
  readonly address?: string;
}

export interface BookiusSlot {
  readonly start: string;
  readonly end: string;
  readonly timezone?: string;
}

export interface BookiusBookingTypeBrief {
  readonly id: string;
  readonly title: string;
  readonly slug: string;
  readonly durationMinutes: number;
  readonly targetKind: BookiusTargetKind;
  readonly confirmationMode: BookiusConfirmationMode;
}

export interface BookiusBookingTypeDbo
  extends Omit<BookiusBookingTypeBrief, 'id'>,
    IWithCreated {
  readonly description?: string;
  readonly visibility?: 'public' | 'hidden' | 'disabled';
  readonly availabilitySourceRef?: BookiusAvailabilitySourceRef;
  readonly targetRef?: BookiusExtensionRef;
  readonly location?: BookiusLocation;
  readonly intakeFields?: readonly string[];
}

export interface BookiusBookingPageDbo extends IWithCreated {
  readonly title: string;
  readonly slug: string;
  readonly path: BookiusPublicPath;
  readonly bookingTypeIDs: readonly string[];
  readonly visibility?: 'public' | 'hidden' | 'disabled';
  readonly defaultContext?: Record<string, string>;
}

export interface BookiusBookingTransition {
  readonly at: string;
  readonly by?: string;
  readonly from?: BookiusBookingState;
  readonly to: BookiusBookingState;
  readonly reason?: string;
}

export interface BookiusBookingDbo extends IWithCreated {
  readonly bookingTypeID: string;
  readonly bookingPageID?: string;
  readonly state: BookiusBookingState;
  readonly requestedSlot?: BookiusSlot;
  readonly confirmedSlot?: BookiusSlot;
  readonly visitorName?: string;
  readonly visitorEmail?: string;
  readonly visitorPhone?: string;
  readonly subject?: string;
  readonly message?: string;
  readonly hostRef?: BookiusContactRef;
  readonly inviteeContactRef?: BookiusContactRef;
  readonly resourceRef?: BookiusExtensionRef;
  readonly calendarCommitmentRef?: BookiusAvailabilitySourceRef;
  readonly transitionLog?: readonly BookiusBookingTransition[];
  readonly confirmedAt?: string;
  readonly cancelledAt?: string;
}

export interface BookiusBookingBrief {
  readonly id: string;
  readonly bookingTypeID: string;
  readonly state: BookiusBookingState;
  readonly requestedSlot?: BookiusSlot;
  readonly confirmedSlot?: BookiusSlot;
  readonly visitorName?: string;
  readonly subject?: string;
}

export interface BookiusCreateBookingTypeRequest {
  readonly spaceID: string;
  readonly bookingType: BookiusBookingTypeDbo;
}

export interface BookiusCreateBookingRequest {
  readonly spaceID?: string;
  readonly bookingTypeID: string;
  readonly bookingPageID?: string;
  readonly requestedSlot: BookiusSlot;
  readonly visitorName: string;
  readonly visitorEmail: string;
  readonly visitorPhone?: string;
  readonly subject?: string;
  readonly message?: string;
}

export interface BookiusCompetitionEntryTarget {
  readonly extensionID: string;
  readonly eventID: string;
  readonly tournamentID: string;
  readonly competitionID: string;
  readonly targetVersion: number;
}

/** No amount or currency: Bookius resolves the offer server-to-server. */
export interface BookiusCompetitionEntryReservationRequest {
  readonly requestID: string;
  readonly target: BookiusCompetitionEntryTarget;
  readonly participantReference: string;
  readonly entryReference: string;
}

export type BookiusCompetitionEntryReservationState =
  | 'requested'
  | 'waitlisted'
  | 'held'
  | 'checkout-ready'
  | 'confirmed'
  | 'failed'
  | 'expired'
  | 'cancelled'
  | 'refunded';

export interface BookiusCompetitionEntryReservation {
  readonly id: string;
  readonly requestID: string;
  readonly target: BookiusCompetitionEntryTarget;
  readonly participantReference: string;
  readonly entryReference: string;
  readonly bookingRevision: number;
  readonly state: BookiusCompetitionEntryReservationState;
  /** Server-authored hold terms. */
  readonly amountMinor: number;
  readonly currency: string;
  readonly offerReference: string;
  readonly offerVersion: number;
  readonly offerChecksum: string;
  readonly paymentState:
    | 'free'
    | 'not_started'
    | 'checkout_open'
    | 'paid'
    | 'refund_pending'
    | 'refunded'
    | 'failed';
  readonly confirmationEvidence?: BookiusCompetitionEntryConfirmationEvidence;
  readonly checkoutOperation: 'none' | 'pending' | 'ready' | 'failed';
  readonly confirmationDelivery: BookiusCompetitionEntryDeliveryState;
  readonly refundDelivery: BookiusCompetitionEntryDeliveryState;
  readonly expiresAt?: string;
}

export interface BookiusCompetitionEntrySettlementNotification {
  readonly settlementID: string;
  readonly settlementReference: string;
  readonly refundReference: string;
  readonly reservationID: string;
  readonly status: 'paid' | 'failed' | 'refunded';
  readonly amountMinor: number;
  readonly currency: string;
  readonly occurredAt: string;
}

export interface BookiusCompetitionEntryConfirmationEvidence {
  readonly kind: 'free' | 'settled';
  readonly settlementReference?: string;
}

export type BookiusCompetitionEntryDeliveryState =
  | 'none'
  | 'pending'
  | 'delivered'
  | 'failed';

/**
 * Public booking-type shape retained from the initial standalone contract.
 * The richer DBO above is used by the implementation; this remains useful to
 * lightweight consumers that do not need persistence metadata.
 */
export interface BookiusBookingType {
  readonly id?: string;
  readonly title: string;
  readonly slug: string;
  readonly description?: string;
  readonly durationMinutes: number;
  readonly targetKind: BookiusTargetKind;
  readonly targetRef?: BookiusExtensionRef;
  readonly availabilitySourceRef?: BookiusAvailabilitySourceRef;
  readonly location?: BookiusLocation;
  readonly confirmationMode: BookiusConfirmationMode;
  readonly visibility?: 'public' | 'hidden' | 'disabled';
}

export interface BookiusBookingPage {
  readonly id?: string;
  readonly title: string;
  readonly slug: string;
  readonly handle?: string;
  readonly businessSlug?: string;
  readonly category?: string;
  readonly resourceSlug?: string;
  readonly bookingTypeIDs: readonly string[];
  readonly visibility?: 'public' | 'hidden' | 'disabled';
}

export interface BookiusBooking {
  readonly id?: string;
  readonly bookingTypeID: string;
  readonly bookingPageID?: string;
  readonly state: BookiusBookingState;
  readonly requestedSlot?: BookiusSlot;
  readonly confirmedSlot?: BookiusSlot;
  readonly visitorName?: string;
  readonly visitorEmail?: string;
  readonly visitorPhone?: string;
  readonly subject?: string;
  readonly message?: string;
  readonly hostRef?: BookiusContactRef;
  readonly inviteeContactRef?: BookiusContactRef;
  readonly resourceRef?: BookiusExtensionRef;
  readonly calendarCommitmentRef?: BookiusAvailabilitySourceRef;
}

export interface BookiusCompetitionEntryLifecycleCommand {
  readonly commandID: string;
  readonly reservationID: string;
}

export type BookiusApproveCompetitionEntryReservationRequest =
  BookiusCompetitionEntryLifecycleCommand;

export type BookiusPromoteCompetitionEntryWaitlistRequest =
  BookiusCompetitionEntryLifecycleCommand;

export type BookiusExpireCompetitionEntryReservationRequest =
  BookiusCompetitionEntryLifecycleCommand;

export interface BookiusCancelCompetitionEntryReservationRequest
  extends BookiusCompetitionEntryLifecycleCommand {
  readonly origin: 'participant' | 'organiser';
  readonly actorReference: string;
  readonly authorityEvidence: string;
  readonly reason: string;
}

export interface BookiusCompetitionEntryCancellationValidation {
  readonly authorized: boolean;
  /** False only for a locked participant cancellation without an organiser refund override. */
  readonly refundAuthorized: boolean;
  readonly currentTournamentVersion: number;
  readonly registrationLocked: boolean;
  readonly authoriserReference?: string;
  readonly authorityEvidence: string;
  readonly validatedAt: string;
}
