package facade4calendarius

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/sneat-co/sneat-ext-contracts/calendarius/calendariusmodels"
)

// RecurringHappeningsFacade exposes recurring calendar behavior without
// leaking Calendarius DBOs, slot schemas or recurrence implementation types.
type RecurringHappeningsFacade interface {
	// ListActiveRecurring lists active recurring happenings of a kind across
	// every Space (trusted host use only).
	ListActiveRecurring(ctx context.Context, kind string) ([]calendariusmodels.RecurringHappening, error)
	GetRecurring(ctx context.Context, spaceID, happeningID string) (calendariusmodels.RecurringHappening, error)
	SetExtensionData(ctx context.Context, spaceID, happeningID, extensionID string, data json.RawMessage) error
	// ExpandOccurrences returns the occurrences whose effective start is in
	// [from, to), ordered by effective start then slot ID. Each slot is
	// expanded in its own stored IANA time zone; timezone is only a fallback
	// for legacy slots without one (empty = UTC). An invalid stored or
	// fallback zone is an error, never a silent UTC. Per-date adjustments are
	// honoured: cancelled occurrences are returned with Cancelled=true, and
	// rescheduled ones carry their moved Start/End with Rescheduled=true.
	// Wall-clock times are kept across DST changes; a time in a DST gap moves
	// forward and an ambiguous time takes the earlier instant
	// (calendariusmodels.ResolveLocalTime).
	ExpandOccurrences(ctx context.Context, spaceID, happeningID, timezone string, from, to time.Time) ([]calendariusmodels.Occurrence, error)
}

var (
	// ErrRecurringHappeningNotFound is returned when a mutation targets a
	// Happening (or slot) that does not exist or is not recurring.
	ErrRecurringHappeningNotFound = errors.New("recurring happening not found")
	// ErrRecurringHappeningClosed is returned when an occurrence mutation
	// targets a cancelled or deleted recurring Happening.
	ErrRecurringHappeningClosed = errors.New("recurring happening is closed")
	// ErrOccurrenceNotFound is returned when the (slot, local date) pair is
	// not an occurrence produced by the recurrence rule.
	ErrOccurrenceNotFound = errors.New("occurrence not found")
)

// RecurringScheduleFacade creates and mutates recurring weekly Happenings
// (group/personal schedules, availability windows). It embeds the read side.
//
// userID is the acting principal; providers authorize it against spaceID
// before revealing or mutating anything (facade.ErrUnauthorized-compatible
// errors are mapped to ErrEventHappeningUnauthorized). Validation failures
// wrap calendariusmodels.ErrInvalidRecurringRequest.
//
// Idempotency: CreateRecurringHappening derives the Happening ID from
// (spaceID, RequestID) and stores the payload fingerprint, so a replay returns
// Disposition=reused and a different payload returns ErrRequestIDConflict. The
// other mutations are state-setting (cancel, move-to, replace), so replaying
// the same request is a no-op; requestID is validated and recorded but a later
// different request wins.
type RecurringScheduleFacade interface {
	RecurringHappeningsFacade

	CreateRecurringHappening(
		ctx context.Context,
		userID, spaceID string,
		request calendariusmodels.CreateRecurringHappeningRequest,
	) (calendariusmodels.RecurringHappeningMutation, error)

	// CancelOccurrence cancels one occurrence identified by its ORIGINAL local
	// date in the slot's time zone.
	CancelOccurrence(
		ctx context.Context,
		userID, spaceID, happeningID, slotID, localDate, reason, requestID string,
	) error

	// RescheduleOccurrence moves one occurrence (identified by its original
	// local date) to newStart, keeping its duration. The move is limited to
	// calendariusmodels.OccurrenceRescheduleMaxShift; rescheduling a cancelled
	// occurrence reinstates it at the new time.
	RescheduleOccurrence(
		ctx context.Context,
		userID, spaceID, happeningID, slotID, localDate string,
		newStart time.Time,
		requestID string,
	) error

	// SetParticipants replaces the participant ID list.
	SetParticipants(
		ctx context.Context,
		userID, spaceID, happeningID string,
		participants []string,
		requestID string,
	) error

	// CancelHappening cancels the whole recurring Happening; already
	// cancelled is a no-op.
	CancelHappening(
		ctx context.Context,
		userID, spaceID, happeningID, reason, requestID string,
	) error
}
