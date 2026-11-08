package facade4calendariustest

import (
	"context"
	"encoding/json"
	"errors"
	"slices"
	"testing"
	"time"

	"github.com/sneat-co/sneat-ext-contracts/calendarius/calendariusmodels"
	"github.com/sneat-co/sneat-ext-contracts/calendarius/facade4calendarius"
)

// RecurringScheduleConformanceUserID is a member of
// RecurringScheduleConformanceSpaceID; RecurringScheduleConformanceOutsiderID
// is not. A harness must seed them before returning the facade.
const (
	RecurringScheduleConformanceUserID     = "user1"
	RecurringScheduleConformanceOutsiderID = "outsider"
	RecurringScheduleConformanceSpaceID    = "space1"
)

// ConformanceRecurringRequest is the weekly Thursday 18:00 Europe/Dublin
// availability Happening every conformance subtest starts from.
func ConformanceRecurringRequest(requestID string) calendariusmodels.CreateRecurringHappeningRequest {
	return calendariusmodels.CreateRecurringHappeningRequest{
		RequestID:    requestID,
		Kind:         calendariusmodels.RecurringKindAvailability,
		Title:        "Mentor availability",
		TimeZone:     "Europe/Dublin",
		StartDate:    "2026-10-01",
		Slots:        []calendariusmodels.RecurringWeeklySlot{{Weekdays: []calendariusmodels.WeekdayCode{calendariusmodels.Thursday}, Start: "18:00", DurationMinutes: 60}},
		Participants: []string{"mentor1"},
		Extensions:   map[string]json.RawMessage{"commitius": json.RawMessage(`{"purpose":"availability"}`)},
	}
}

// RunRecurringScheduleFacadeConformance verifies the portable behaviour of a
// RecurringScheduleFacade. newFacade must return a fresh implementation for
// each subtest, with the conformance user/space/outsider fixtures seeded.
func RunRecurringScheduleFacadeConformance(
	t *testing.T,
	newFacade func(t *testing.T) facade4calendarius.RecurringScheduleFacade,
) {
	t.Helper()
	ctx := context.Background()
	const user, space = RecurringScheduleConformanceUserID, RecurringScheduleConformanceSpaceID
	dublin, _ := time.LoadLocation("Europe/Dublin")
	// Thursdays 2026-10-22 (IST, UTC+1) and 2026-10-29 (GMT, UTC+0) straddle
	// the 2026-10-25 Dublin DST end.
	from := time.Date(2026, time.October, 20, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, time.November, 3, 0, 0, 0, 0, time.UTC)

	create := func(t *testing.T, facade facade4calendarius.RecurringScheduleFacade, requestID string) string {
		t.Helper()
		created, err := facade.CreateRecurringHappening(ctx, user, space, ConformanceRecurringRequest(requestID))
		if err != nil {
			t.Fatalf("CreateRecurringHappening() error = %v", err)
		}
		if created.HappeningID == "" || created.Disposition != calendariusmodels.RecurringHappeningCreated ||
			!slices.Equal(created.SlotIDs, []string{"0"}) {
			t.Fatalf("CreateRecurringHappening() = %+v", created)
		}
		return created.HappeningID
	}
	expand := func(t *testing.T, facade facade4calendarius.RecurringScheduleFacade, happeningID string) []calendariusmodels.Occurrence {
		t.Helper()
		occurrences, err := facade.ExpandOccurrences(ctx, space, happeningID, "", from, to)
		if err != nil {
			t.Fatalf("ExpandOccurrences() error = %v", err)
		}
		return occurrences
	}

	t.Run("CreateIsIdempotentAndConflictsOnChangedPayload", func(t *testing.T) {
		facade := newFacade(t)
		id := create(t, facade, "create-1")
		replay, err := facade.CreateRecurringHappening(ctx, user, space, ConformanceRecurringRequest("create-1"))
		if err != nil || replay.HappeningID != id || replay.Disposition != calendariusmodels.RecurringHappeningReused {
			t.Fatalf("replay = (%+v, %v), want reused %s", replay, err, id)
		}
		changed := ConformanceRecurringRequest("create-1")
		changed.Title = "Changed"
		if _, err = facade.CreateRecurringHappening(ctx, user, space, changed); !errors.Is(err, facade4calendarius.ErrRequestIDConflict) {
			t.Fatalf("changed payload error = %v, want ErrRequestIDConflict", err)
		}
		got, err := facade.GetRecurring(ctx, space, id)
		if err != nil || !got.Active || got.SpaceID != space || !slices.Equal(got.Participants, []string{"mentor1"}) ||
			string(got.Extensions["commitius"]) != `{"purpose":"availability"}` {
			t.Fatalf("GetRecurring() = (%+v, %v)", got, err)
		}
	})

	t.Run("CreateRejectsInvalidTimeZoneAndOutsider", func(t *testing.T) {
		facade := newFacade(t)
		invalid := ConformanceRecurringRequest("create-invalid")
		invalid.TimeZone = "Mars/Olympus"
		if _, err := facade.CreateRecurringHappening(ctx, user, space, invalid); !errors.Is(err, calendariusmodels.ErrInvalidRecurringRequest) {
			t.Fatalf("invalid zone error = %v", err)
		}
		if _, err := facade.CreateRecurringHappening(ctx, RecurringScheduleConformanceOutsiderID, space, ConformanceRecurringRequest("create-outsider")); !errors.Is(err, facade4calendarius.ErrEventHappeningUnauthorized) {
			t.Fatalf("outsider error = %v, want ErrEventHappeningUnauthorized", err)
		}
	})

	t.Run("ExpansionKeepsWallClockAcrossDSTWithStableIdentity", func(t *testing.T) {
		facade := newFacade(t)
		id := create(t, facade, "create-dst")
		occurrences := expand(t, facade, id)
		if len(occurrences) != 2 {
			t.Fatalf("occurrences = %+v, want 2", occurrences)
		}
		for i, wantDate := range []string{"2026-10-22", "2026-10-29"} {
			o := occurrences[i]
			if o.HappeningID != id || o.SlotID != "0" || o.LocalDate != wantDate || o.Cancelled || o.Rescheduled {
				t.Fatalf("occurrence[%d] identity = %+v", i, o)
			}
			if local := o.Start.In(dublin); local.Hour() != 18 || local.Format(time.DateOnly) != wantDate ||
				!o.OriginalStart.Equal(o.Start) || o.End.Sub(o.Start) != time.Hour {
				t.Fatalf("occurrence[%d] timing = %+v", i, o)
			}
		}
		if occurrences[0].Start.UTC().Hour() != 17 || occurrences[1].Start.UTC().Hour() != 18 {
			t.Fatalf("UTC instants ignore DST: %v, %v", occurrences[0].Start, occurrences[1].Start)
		}
	})

	t.Run("CancelAndRescheduleOccurrences", func(t *testing.T) {
		facade := newFacade(t)
		id := create(t, facade, "create-adjust")
		for range 2 { // replay is a no-op
			if err := facade.CancelOccurrence(ctx, user, space, id, "0", "2026-10-22", "mentor away", "cancel-1"); err != nil {
				t.Fatalf("CancelOccurrence() error = %v", err)
			}
		}
		moved := time.Date(2026, time.October, 30, 9, 30, 0, 0, dublin)
		if err := facade.RescheduleOccurrence(ctx, user, space, id, "0", "2026-10-29", moved, "move-1"); err != nil {
			t.Fatalf("RescheduleOccurrence() error = %v", err)
		}
		occurrences := expand(t, facade, id)
		if len(occurrences) != 2 {
			t.Fatalf("occurrences = %+v, want 2", occurrences)
		}
		if o := occurrences[0]; o.LocalDate != "2026-10-22" || !o.Cancelled {
			t.Fatalf("cancelled occurrence = %+v", o)
		}
		if o := occurrences[1]; o.LocalDate != "2026-10-29" || !o.Rescheduled || o.Cancelled || !o.Start.Equal(moved) ||
			o.End.Sub(o.Start) != time.Hour || o.OriginalStart.In(dublin).Hour() != 18 {
			t.Fatalf("rescheduled occurrence = %+v", o)
		}
		// Rescheduling a cancelled occurrence reinstates it.
		reinstated := time.Date(2026, time.October, 22, 19, 0, 0, 0, dublin)
		if err := facade.RescheduleOccurrence(ctx, user, space, id, "0", "2026-10-22", reinstated, "move-2"); err != nil {
			t.Fatalf("RescheduleOccurrence(cancelled) error = %v", err)
		}
		if o := expand(t, facade, id)[0]; o.Cancelled || !o.Start.Equal(reinstated) {
			t.Fatalf("reinstated occurrence = %+v", o)
		}
		// A move out of the window drops the occurrence from it.
		if err := facade.RescheduleOccurrence(ctx, user, space, id, "0", "2026-10-29", to.Add(24*time.Hour), "move-3"); err != nil {
			t.Fatalf("RescheduleOccurrence(out) error = %v", err)
		}
		if got := expand(t, facade, id); len(got) != 1 || got[0].LocalDate != "2026-10-22" {
			t.Fatalf("after moving out = %+v", got)
		}
	})

	t.Run("OccurrenceMutationErrors", func(t *testing.T) {
		facade := newFacade(t)
		id := create(t, facade, "create-errors")
		checks := []struct {
			name string
			err  error
			want error
		}{
			{"unknown slot", facade.CancelOccurrence(ctx, user, space, id, "9", "2026-10-22", "", "r1"), facade4calendarius.ErrRecurringHappeningNotFound},
			{"unknown happening", facade.CancelOccurrence(ctx, user, space, "missing", "0", "2026-10-22", "", "r2"), facade4calendarius.ErrRecurringHappeningNotFound},
			{"not an occurrence", facade.CancelOccurrence(ctx, user, space, id, "0", "2026-10-23", "", "r3"), facade4calendarius.ErrOccurrenceNotFound},
			{"before start date", facade.CancelOccurrence(ctx, user, space, id, "0", "2026-09-24", "", "r4"), facade4calendarius.ErrOccurrenceNotFound},
			{"bad date", facade.CancelOccurrence(ctx, user, space, id, "0", "22-10-2026", "", "r5"), calendariusmodels.ErrInvalidRecurringRequest},
			{"no request ID", facade.CancelOccurrence(ctx, user, space, id, "0", "2026-10-22", "", ""), calendariusmodels.ErrInvalidRecurringRequest},
			{"too far", facade.RescheduleOccurrence(ctx, user, space, id, "0", "2026-10-22", from.AddDate(0, 3, 0), "r6"), calendariusmodels.ErrInvalidRecurringRequest},
			{"outsider", facade.CancelOccurrence(ctx, RecurringScheduleConformanceOutsiderID, space, id, "0", "2026-10-22", "", "r7"), facade4calendarius.ErrEventHappeningUnauthorized},
			{"duplicate participants", facade.SetParticipants(ctx, user, space, id, []string{"a", "a"}, "r8"), calendariusmodels.ErrInvalidRecurringRequest},
		}
		for _, check := range checks {
			if !errors.Is(check.err, check.want) {
				t.Errorf("%s: error = %v, want %v", check.name, check.err, check.want)
			}
		}
	})

	t.Run("SetParticipantsAndCancelHappening", func(t *testing.T) {
		facade := newFacade(t)
		id := create(t, facade, "create-participants")
		if err := facade.SetParticipants(ctx, user, space, id, []string{"mentor2", "member1"}, "participants-1"); err != nil {
			t.Fatalf("SetParticipants() error = %v", err)
		}
		if got, err := facade.GetRecurring(ctx, space, id); err != nil || !slices.Equal(got.Participants, []string{"mentor2", "member1"}) {
			t.Fatalf("GetRecurring() = (%+v, %v)", got, err)
		}
		for range 2 { // replay is a no-op
			if err := facade.CancelHappening(ctx, user, space, id, "group closed", "cancel-happening"); err != nil {
				t.Fatalf("CancelHappening() error = %v", err)
			}
		}
		if got, err := facade.GetRecurring(ctx, space, id); err != nil || got.Active {
			t.Fatalf("GetRecurring() after cancel = (%+v, %v)", got, err)
		}
		if err := facade.CancelOccurrence(ctx, user, space, id, "0", "2026-10-22", "", "late"); !errors.Is(err, facade4calendarius.ErrRecurringHappeningClosed) {
			t.Fatalf("CancelOccurrence(closed) error = %v", err)
		}
	})
}
