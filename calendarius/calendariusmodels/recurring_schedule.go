package calendariusmodels

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"
)

// WeekdayCode is the two-letter Calendarius weekday code used by weekly slots.
type WeekdayCode string

const (
	Monday    WeekdayCode = "mo"
	Tuesday   WeekdayCode = "tu"
	Wednesday WeekdayCode = "we"
	Thursday  WeekdayCode = "th"
	Friday    WeekdayCode = "fr"
	Saturday  WeekdayCode = "sa"
	Sunday    WeekdayCode = "su"
)

// GoWeekday maps the code to time.Weekday; ok is false for an unknown code.
func (v WeekdayCode) GoWeekday() (weekday time.Weekday, ok bool) {
	switch v {
	case Monday:
		return time.Monday, true
	case Tuesday:
		return time.Tuesday, true
	case Wednesday:
		return time.Wednesday, true
	case Thursday:
		return time.Thursday, true
	case Friday:
		return time.Friday, true
	case Saturday:
		return time.Saturday, true
	case Sunday:
		return time.Sunday, true
	}
	return 0, false
}

// Recurring Happening kinds accepted by CreateRecurringHappening. They are the
// Calendarius HappeningKind values that make sense for a weekly schedule.
const (
	RecurringKindActivity     = "activity"
	RecurringKindAppointment  = "appointment"
	RecurringKindTask         = "task"
	RecurringKindAvailability = "availability"
)

// Limits of the recurring-schedule contract (UTF-8 byte limits).
const (
	RecurringRequestIDMaxBytes     = 200
	RecurringTitleMaxBytes         = 100
	RecurringReasonMaxBytes        = 10000
	RecurringSlotsMax              = 50
	RecurringParticipantsMax       = 100
	RecurringParticipantIDMaxBytes = 200
	RecurringSlotIDMaxBytes        = 50
	// RecurringDurationMaxMinutes bounds one occurrence (or availability
	// window) to a single day.
	RecurringDurationMaxMinutes = 24 * 60
	// OccurrenceRescheduleMaxShift bounds how far a rescheduled occurrence may
	// move from its original start. Providers rely on it to find moved
	// occurrences when expanding a window.
	OccurrenceRescheduleMaxShift = 31 * 24 * time.Hour
)

// ErrInvalidRecurringRequest is wrapped by every request validation failure.
var ErrInvalidRecurringRequest = errors.New("invalid recurring happening request")

func invalidRecurring(format string, args ...any) error {
	return fmt.Errorf("%w: %s", ErrInvalidRecurringRequest, fmt.Sprintf(format, args...))
}

// RecurringWeeklySlot is one weekly slot: every listed weekday at Start local
// wall-clock time (in the happening TimeZone) for DurationMinutes.
//
// ID is optional. When empty the provider assigns strconv.Itoa(slot index), so
// a consumer that addresses slots by index can derive the slot ID itself.
type RecurringWeeklySlot struct {
	ID              string
	Weekdays        []WeekdayCode
	Start           string // HH:MM, 24-hour
	DurationMinutes int
}

// CreateRecurringHappeningRequest creates one recurring weekly Happening.
// RequestID is the idempotency key: a replay with an identical payload returns
// the originally created Happening; another payload under the same RequestID
// (and the same principal and Space) returns ErrRequestIDConflict.
type CreateRecurringHappeningRequest struct {
	RequestID string
	Kind      string
	Title     string
	// TimeZone is the IANA TZDB name every slot's wall-clock time is in.
	TimeZone string
	// StartDate (YYYY-MM-DD, local) is the first eligible occurrence date.
	StartDate string
	Slots     []RecurringWeeklySlot
	// Participants are opaque consumer principal IDs (contact or user IDs).
	Participants []string
	// Extensions are extension-owned JSON objects stored on the Happening.
	Extensions map[string]json.RawMessage
}

// SlotID returns the effective ID of the slot at index i.
func (v CreateRecurringHappeningRequest) SlotID(i int) string {
	if id := v.Slots[i].ID; id != "" {
		return id
	}
	return fmt.Sprintf("%d", i)
}

// Validate rejects malformed, unbounded or ambiguous input. Every error wraps
// ErrInvalidRecurringRequest.
func (v CreateRecurringHappeningRequest) Validate() error {
	if err := ValidateRecurringRequestID(v.RequestID); err != nil {
		return err
	}
	switch v.Kind {
	case RecurringKindActivity, RecurringKindAppointment, RecurringKindTask, RecurringKindAvailability:
	default:
		return invalidRecurring("kind %q is not a recurring kind", v.Kind)
	}
	if strings.TrimSpace(v.Title) == "" || len(v.Title) > RecurringTitleMaxBytes {
		return invalidRecurring("title is required and must be at most %d bytes", RecurringTitleMaxBytes)
	}
	if _, err := LoadRecurringTimeZone(v.TimeZone); err != nil {
		return err
	}
	if _, err := time.Parse(time.DateOnly, v.StartDate); err != nil {
		return invalidRecurring("startDate must be YYYY-MM-DD: %v", err)
	}
	if len(v.Slots) == 0 || len(v.Slots) > RecurringSlotsMax {
		return invalidRecurring("slots must contain 1..%d items", RecurringSlotsMax)
	}
	ids := make(map[string]bool, len(v.Slots))
	for i, slot := range v.Slots {
		id := v.SlotID(i)
		if strings.TrimSpace(id) != id || len(id) > RecurringSlotIDMaxBytes || strings.ContainsAny(id, "/.") {
			return invalidRecurring("slots[%d].id %q is invalid", i, id)
		}
		if ids[id] {
			return invalidRecurring("slots[%d].id %q is duplicated", i, id)
		}
		ids[id] = true
		if len(slot.Weekdays) == 0 {
			return invalidRecurring("slots[%d].weekdays is required", i)
		}
		for j, weekday := range slot.Weekdays {
			if _, ok := weekday.GoWeekday(); !ok {
				return invalidRecurring("slots[%d].weekdays[%d] %q is unknown", i, j, weekday)
			}
			if slices.Contains(slot.Weekdays[:j], weekday) {
				return invalidRecurring("slots[%d].weekdays[%d] %q is duplicated", i, j, weekday)
			}
		}
		if _, err := time.Parse("15:04", slot.Start); err != nil || len(slot.Start) != 5 {
			return invalidRecurring("slots[%d].start must be HH:MM", i)
		}
		if slot.DurationMinutes <= 0 || slot.DurationMinutes > RecurringDurationMaxMinutes {
			return invalidRecurring("slots[%d].durationMinutes must be 1..%d", i, RecurringDurationMaxMinutes)
		}
	}
	if err := ValidateRecurringParticipants(v.Participants); err != nil {
		return err
	}
	for extensionID, data := range v.Extensions {
		if strings.TrimSpace(extensionID) == "" {
			return invalidRecurring("extensions has an empty extension ID")
		}
		var object map[string]any
		if err := json.Unmarshal(data, &object); err != nil || object == nil {
			return invalidRecurring("extensions[%s] must be a JSON object", extensionID)
		}
	}
	return nil
}

// Fingerprint is a stable hash of the full payload, used to tell an identical
// replay from a RequestID conflict. Call it on a validated request.
func (v CreateRecurringHappeningRequest) Fingerprint() string {
	canonical := v
	canonical.Extensions = nil
	extensions := make(map[string]any, len(v.Extensions))
	for id, data := range v.Extensions {
		var decoded any
		_ = json.Unmarshal(data, &decoded)
		extensions[id] = decoded
	}
	encoded, _ := json.Marshal(struct {
		Request    CreateRecurringHappeningRequest
		Extensions map[string]any
	}{canonical, extensions})
	sum := sha256.Sum256(encoded)
	return hex.EncodeToString(sum[:])
}

// ValidateRecurringRequestID validates an idempotency key.
func ValidateRecurringRequestID(requestID string) error {
	if strings.TrimSpace(requestID) == "" || len(requestID) > RecurringRequestIDMaxBytes {
		return invalidRecurring("requestID is required and must be at most %d bytes", RecurringRequestIDMaxBytes)
	}
	return nil
}

// ValidateRecurringParticipants validates a participant ID list (may be empty).
func ValidateRecurringParticipants(participants []string) error {
	if len(participants) > RecurringParticipantsMax {
		return invalidRecurring("participants must contain at most %d items", RecurringParticipantsMax)
	}
	for i, id := range participants {
		if strings.TrimSpace(id) == "" || len(id) > RecurringParticipantIDMaxBytes {
			return invalidRecurring("participants[%d] is empty or too long", i)
		}
		if slices.Contains(participants[:i], id) {
			return invalidRecurring("participants[%d] %q is duplicated", i, id)
		}
	}
	return nil
}

// LoadRecurringTimeZone loads an explicit IANA zone. Empty, "Local" and unknown
// names are rejected (never silently treated as UTC).
func LoadRecurringTimeZone(name string) (*time.Location, error) {
	if name == "" || name == "Local" {
		return nil, invalidRecurring("timeZone must be an explicit IANA TZDB name")
	}
	location, err := time.LoadLocation(name)
	if err != nil {
		return nil, invalidRecurring("timeZone %q is not an IANA TZDB name: %v", name, err)
	}
	return location, nil
}

// ResolveLocalTime returns the instant of a local wall-clock time in loc with
// the documented DST policy: a time inside a DST gap moves forward by the
// size of the gap (02:30 on a spring-forward night becomes 03:30), and an
// ambiguous time during a fall-back fold takes the earlier instant.
func ResolveLocalTime(year int, month time.Month, day, hour, minute int, loc *time.Location) time.Time {
	naive := time.Date(year, month, day, hour, minute, 0, 0, time.UTC)
	probe := time.Date(year, month, day, 12, 0, 0, 0, loc)
	_, offsetBefore := probe.Add(-36 * time.Hour).Zone()
	_, offsetAfter := probe.Add(36 * time.Hour).Zone()
	// At a fall-back fold the pre-transition offset is the larger one, so its
	// candidate is the earlier instant: the first match wins.
	for _, offset := range []int{offsetBefore, offsetAfter} {
		candidate := naive.Add(-time.Duration(offset) * time.Second).In(loc)
		if sameWallClock(candidate, year, month, day, hour, minute) {
			return candidate
		}
	}
	// DST gap: interpret with the pre-transition offset, which moves forward.
	return naive.Add(-time.Duration(offsetBefore) * time.Second).In(loc)
}

func sameWallClock(t time.Time, year int, month time.Month, day, hour, minute int) bool {
	y, m, d := t.Date()
	return y == year && m == month && d == day && t.Hour() == hour && t.Minute() == minute
}

// RecurringHappeningDisposition reports whether a create made a new Happening.
type RecurringHappeningDisposition string

const (
	RecurringHappeningCreated RecurringHappeningDisposition = "created"
	RecurringHappeningReused  RecurringHappeningDisposition = "reused"
)

// RecurringHappeningMutation is the result of CreateRecurringHappening.
type RecurringHappeningMutation struct {
	HappeningID string
	// SlotIDs are the effective slot IDs in request order.
	SlotIDs     []string
	Disposition RecurringHappeningDisposition
}
