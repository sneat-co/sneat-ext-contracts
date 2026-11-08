package calendariusmodels

import (
	"encoding/json"
	"time"
)

// RecurringHappening is the public projection of a recurring happening needed
// by extensions that attach their own overlay data. Slots, adjustments and the
// persistence schema remain private to Calendarius.
type RecurringHappening struct {
	ID string
	// SpaceID is the Space that owns the Happening. Providers populate it so a
	// cross-Space listing can be followed by GetRecurring/ExpandOccurrences.
	SpaceID    string
	Active     bool
	Extensions map[string]json.RawMessage
	// Participants are the opaque participant IDs set through the
	// recurring-schedule contract (empty for Happenings that never used it).
	Participants []string
}

// Occurrence is one expanded recurring slot in an absolute time range.
//
// Identity: (HappeningID, SlotID, LocalDate) identifies an occurrence stably,
// LocalDate being the ORIGINAL local date in the slot's time zone. A reschedule
// or cancellation never changes the identity. Start/End are the effective
// instants after any reschedule; OriginalStart is the instant the recurrence
// rule produced.
type Occurrence struct {
	Start time.Time
	End   time.Time

	HappeningID   string
	SlotID        string
	LocalDate     string // YYYY-MM-DD, original local date
	OriginalStart time.Time
	// Cancelled occurrences are still returned (flagged) so callers can show
	// them and keep stable identities; callers that only need live
	// occurrences must skip them.
	Cancelled   bool
	Rescheduled bool
}
