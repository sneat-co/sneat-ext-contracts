package calendariusmodels

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"
)

func validRecurringRequest() CreateRecurringHappeningRequest {
	return CreateRecurringHappeningRequest{
		RequestID: "r1",
		Kind:      RecurringKindAvailability,
		Title:     "Availability",
		TimeZone:  "America/New_York",
		StartDate: "2026-03-01",
		Slots: []RecurringWeeklySlot{
			{Weekdays: []WeekdayCode{Monday, Wednesday}, Start: "09:00", DurationMinutes: 120},
			{ID: "evening", Weekdays: []WeekdayCode{Friday}, Start: "18:30", DurationMinutes: 20},
		},
		Participants: []string{"m1", "u1"},
		Extensions:   map[string]json.RawMessage{"commitius": json.RawMessage(`{"a":1}`)},
	}
}

func TestCreateRecurringHappeningRequest_Validate(t *testing.T) {
	if err := validRecurringRequest().Validate(); err != nil {
		t.Fatalf("valid request: %v", err)
	}
	if got := validRecurringRequest().SlotID(0); got != "0" {
		t.Fatalf("SlotID(0) = %q", got)
	}
	if got := validRecurringRequest().SlotID(1); got != "evening" {
		t.Fatalf("SlotID(1) = %q", got)
	}
	cases := map[string]func(r *CreateRecurringHappeningRequest){
		"request ID":        func(r *CreateRecurringHappeningRequest) { r.RequestID = " " },
		"kind":              func(r *CreateRecurringHappeningRequest) { r.Kind = "event" },
		"title":             func(r *CreateRecurringHappeningRequest) { r.Title = strings.Repeat("x", RecurringTitleMaxBytes+1) },
		"empty zone":        func(r *CreateRecurringHappeningRequest) { r.TimeZone = "" },
		"local zone":        func(r *CreateRecurringHappeningRequest) { r.TimeZone = "Local" },
		"unknown zone":      func(r *CreateRecurringHappeningRequest) { r.TimeZone = "Mars/Olympus" },
		"start date":        func(r *CreateRecurringHappeningRequest) { r.StartDate = "2026-13-01" },
		"no slots":          func(r *CreateRecurringHappeningRequest) { r.Slots = nil },
		"slot ID":           func(r *CreateRecurringHappeningRequest) { r.Slots[1].ID = "a/b" },
		"duplicate slot ID": func(r *CreateRecurringHappeningRequest) { r.Slots[1].ID = "0" },
		"no weekdays":       func(r *CreateRecurringHappeningRequest) { r.Slots[0].Weekdays = nil },
		"unknown weekday":   func(r *CreateRecurringHappeningRequest) { r.Slots[0].Weekdays = []WeekdayCode{"xx"} },
		"duplicate weekday": func(r *CreateRecurringHappeningRequest) { r.Slots[0].Weekdays = []WeekdayCode{Monday, Monday} },
		"start":             func(r *CreateRecurringHappeningRequest) { r.Slots[0].Start = "9:00" },
		"start hour":        func(r *CreateRecurringHappeningRequest) { r.Slots[0].Start = "25:00" },
		"duration":          func(r *CreateRecurringHappeningRequest) { r.Slots[0].DurationMinutes = 0 },
		"long duration":     func(r *CreateRecurringHappeningRequest) { r.Slots[0].DurationMinutes = RecurringDurationMaxMinutes + 1 },
		"participant":       func(r *CreateRecurringHappeningRequest) { r.Participants = []string{""} },
		"dup participant":   func(r *CreateRecurringHappeningRequest) { r.Participants = []string{"a", "a"} },
		"too many":          func(r *CreateRecurringHappeningRequest) { r.Participants = make([]string, RecurringParticipantsMax+1) },
		"extension ID": func(r *CreateRecurringHappeningRequest) {
			r.Extensions = map[string]json.RawMessage{"": json.RawMessage(`{}`)}
		},
		"extension not object": func(r *CreateRecurringHappeningRequest) {
			r.Extensions = map[string]json.RawMessage{"x": json.RawMessage(`[1]`)}
		},
	}
	for name, mutate := range cases {
		request := validRecurringRequest()
		mutate(&request)
		if err := request.Validate(); !errors.Is(err, ErrInvalidRecurringRequest) {
			t.Errorf("%s: error = %v, want ErrInvalidRecurringRequest", name, err)
		}
	}
}

func TestCreateRecurringHappeningRequest_Fingerprint(t *testing.T) {
	a, b := validRecurringRequest(), validRecurringRequest()
	b.Extensions = map[string]json.RawMessage{"commitius": json.RawMessage(`{ "a" : 1 }`)}
	if a.Fingerprint() != b.Fingerprint() {
		t.Fatal("whitespace-only extension differences must not change the fingerprint")
	}
	b.Title = "Other"
	if a.Fingerprint() == b.Fingerprint() {
		t.Fatal("a changed title must change the fingerprint")
	}
}

func TestWeekdayCode_GoWeekday(t *testing.T) {
	want := map[WeekdayCode]time.Weekday{
		Monday: time.Monday, Tuesday: time.Tuesday, Wednesday: time.Wednesday, Thursday: time.Thursday,
		Friday: time.Friday, Saturday: time.Saturday, Sunday: time.Sunday,
	}
	for code, weekday := range want {
		if got, ok := code.GoWeekday(); !ok || got != weekday {
			t.Errorf("%s = %v, %v", code, got, ok)
		}
	}
	if _, ok := WeekdayCode("xx").GoWeekday(); ok {
		t.Error("unknown code must not map")
	}
}

func TestResolveLocalTime_DST(t *testing.T) {
	newYork, _ := time.LoadLocation("America/New_York")
	dublin, _ := time.LoadLocation("Europe/Dublin")
	cases := []struct {
		name    string
		loc     *time.Location
		y       int
		m       time.Month
		d, h, n int
		wantUTC time.Time
	}{
		{"NY normal winter", newYork, 2026, time.January, 5, 9, 0, time.Date(2026, time.January, 5, 14, 0, 0, 0, time.UTC)},
		{"NY spring gap moves forward", newYork, 2026, time.March, 8, 2, 30, time.Date(2026, time.March, 8, 7, 30, 0, 0, time.UTC)},
		{"NY after spring", newYork, 2026, time.March, 8, 9, 0, time.Date(2026, time.March, 8, 13, 0, 0, 0, time.UTC)},
		{"NY fall fold takes earlier", newYork, 2026, time.November, 1, 1, 30, time.Date(2026, time.November, 1, 5, 30, 0, 0, time.UTC)},
		{"NY after fall", newYork, 2026, time.November, 1, 9, 0, time.Date(2026, time.November, 1, 14, 0, 0, 0, time.UTC)},
		{"Dublin spring gap moves forward", dublin, 2026, time.March, 29, 1, 30, time.Date(2026, time.March, 29, 1, 30, 0, 0, time.UTC)},
		{"Dublin summer", dublin, 2026, time.July, 2, 18, 0, time.Date(2026, time.July, 2, 17, 0, 0, 0, time.UTC)},
		{"Dublin fall fold takes earlier", dublin, 2026, time.October, 25, 1, 30, time.Date(2026, time.October, 25, 0, 30, 0, 0, time.UTC)},
		{"Dublin winter", dublin, 2026, time.October, 29, 18, 0, time.Date(2026, time.October, 29, 18, 0, 0, 0, time.UTC)},
	}
	for _, c := range cases {
		got := ResolveLocalTime(c.y, c.m, c.d, c.h, c.n, c.loc)
		if !got.Equal(c.wantUTC) {
			t.Errorf("%s: got %v (%v), want %v", c.name, got.UTC(), got, c.wantUTC)
		}
		if got.Location() != c.loc {
			t.Errorf("%s: location = %v", c.name, got.Location())
		}
	}
}

func TestRecurringValidators(t *testing.T) {
	if err := ValidateRecurringRequestID(strings.Repeat("x", RecurringRequestIDMaxBytes+1)); !errors.Is(err, ErrInvalidRecurringRequest) {
		t.Fatalf("long request ID: %v", err)
	}
	if err := ValidateRecurringParticipants(nil); err != nil {
		t.Fatalf("empty participants: %v", err)
	}
	if loc, err := LoadRecurringTimeZone("Europe/Dublin"); err != nil || loc.String() != "Europe/Dublin" {
		t.Fatalf("LoadRecurringTimeZone = %v, %v", loc, err)
	}
}
