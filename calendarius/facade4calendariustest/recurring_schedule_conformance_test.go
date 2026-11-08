package facade4calendariustest

import (
	"context"
	"encoding/json"
	"errors"
	"sort"
	"testing"
	"time"

	"github.com/sneat-co/sneat-ext-contracts/calendarius/calendariusmodels"
	"github.com/sneat-co/sneat-ext-contracts/calendarius/facade4calendarius"
)

// referenceRecurringFacade is a minimal in-memory provider proving the
// conformance suite is satisfiable.
type referenceRecurringFacade struct {
	happenings map[string]*referenceRecurring
}

type referenceRecurring struct {
	request     calendariusmodels.CreateRecurringHappeningRequest
	fingerprint string
	active      bool
	cancelled   map[string]bool      // slotID/date
	moved       map[string]time.Time // slotID/date
}

var _ facade4calendarius.RecurringScheduleFacade = (*referenceRecurringFacade)(nil)

func (f *referenceRecurringFacade) authorize(userID, spaceID string) error {
	if userID != RecurringScheduleConformanceUserID || spaceID != RecurringScheduleConformanceSpaceID {
		return facade4calendarius.ErrEventHappeningUnauthorized
	}
	return nil
}

func (f *referenceRecurringFacade) CreateRecurringHappening(_ context.Context, userID, spaceID string, request calendariusmodels.CreateRecurringHappeningRequest) (calendariusmodels.RecurringHappeningMutation, error) {
	if err := f.authorize(userID, spaceID); err != nil {
		return calendariusmodels.RecurringHappeningMutation{}, err
	}
	if err := request.Validate(); err != nil {
		return calendariusmodels.RecurringHappeningMutation{}, err
	}
	result := calendariusmodels.RecurringHappeningMutation{HappeningID: "h-" + request.RequestID, Disposition: calendariusmodels.RecurringHappeningCreated}
	for i := range request.Slots {
		result.SlotIDs = append(result.SlotIDs, request.SlotID(i))
	}
	if existing := f.happenings[result.HappeningID]; existing != nil {
		if existing.fingerprint != request.Fingerprint() {
			return calendariusmodels.RecurringHappeningMutation{}, facade4calendarius.ErrRequestIDConflict
		}
		result.Disposition = calendariusmodels.RecurringHappeningReused
		return result, nil
	}
	f.happenings[result.HappeningID] = &referenceRecurring{
		request: request, fingerprint: request.Fingerprint(), active: true,
		cancelled: map[string]bool{}, moved: map[string]time.Time{},
	}
	return result, nil
}

func (f *referenceRecurringFacade) occurrence(userID, spaceID, happeningID, slotID, localDate, requestID string) (*referenceRecurring, time.Time, error) {
	if err := f.authorize(userID, spaceID); err != nil {
		return nil, time.Time{}, err
	}
	if err := calendariusmodels.ValidateRecurringRequestID(requestID); err != nil {
		return nil, time.Time{}, err
	}
	date, err := time.Parse(time.DateOnly, localDate)
	if err != nil {
		return nil, time.Time{}, calendariusmodels.ErrInvalidRecurringRequest
	}
	h := f.happenings[happeningID]
	if h == nil {
		return nil, time.Time{}, facade4calendarius.ErrRecurringHappeningNotFound
	}
	for i, slot := range h.request.Slots {
		if h.request.SlotID(i) != slotID {
			continue
		}
		if !h.active {
			return nil, time.Time{}, facade4calendarius.ErrRecurringHappeningClosed
		}
		start, ok := referenceSlotStart(h.request, slot, date)
		if !ok {
			return nil, time.Time{}, facade4calendarius.ErrOccurrenceNotFound
		}
		return h, start, nil
	}
	return nil, time.Time{}, facade4calendarius.ErrRecurringHappeningNotFound
}

func referenceSlotStart(request calendariusmodels.CreateRecurringHappeningRequest, slot calendariusmodels.RecurringWeeklySlot, date time.Time) (time.Time, bool) {
	if date.Format(time.DateOnly) < request.StartDate {
		return time.Time{}, false
	}
	for _, code := range slot.Weekdays {
		if weekday, _ := code.GoWeekday(); weekday == date.Weekday() {
			loc, _ := calendariusmodels.LoadRecurringTimeZone(request.TimeZone)
			clock, _ := time.Parse("15:04", slot.Start)
			return calendariusmodels.ResolveLocalTime(date.Year(), date.Month(), date.Day(), clock.Hour(), clock.Minute(), loc), true
		}
	}
	return time.Time{}, false
}

func (f *referenceRecurringFacade) CancelOccurrence(_ context.Context, userID, spaceID, happeningID, slotID, localDate, _, requestID string) error {
	h, _, err := f.occurrence(userID, spaceID, happeningID, slotID, localDate, requestID)
	if err == nil {
		h.cancelled[slotID+"/"+localDate] = true
	}
	return err
}

func (f *referenceRecurringFacade) RescheduleOccurrence(_ context.Context, userID, spaceID, happeningID, slotID, localDate string, newStart time.Time, requestID string) error {
	h, original, err := f.occurrence(userID, spaceID, happeningID, slotID, localDate, requestID)
	if err != nil {
		return err
	}
	if shift := newStart.Sub(original); shift > calendariusmodels.OccurrenceRescheduleMaxShift || -shift > calendariusmodels.OccurrenceRescheduleMaxShift {
		return calendariusmodels.ErrInvalidRecurringRequest
	}
	delete(h.cancelled, slotID+"/"+localDate)
	h.moved[slotID+"/"+localDate] = newStart
	return nil
}

func (f *referenceRecurringFacade) SetParticipants(_ context.Context, userID, spaceID, happeningID string, participants []string, requestID string) error {
	if err := f.authorize(userID, spaceID); err != nil {
		return err
	}
	if err := calendariusmodels.ValidateRecurringParticipants(participants); err != nil {
		return err
	}
	h := f.happenings[happeningID]
	if h == nil {
		return facade4calendarius.ErrRecurringHappeningNotFound
	}
	h.request.Participants = append([]string(nil), participants...)
	return calendariusmodels.ValidateRecurringRequestID(requestID)
}

func (f *referenceRecurringFacade) CancelHappening(_ context.Context, userID, spaceID, happeningID, _, requestID string) error {
	if err := f.authorize(userID, spaceID); err != nil {
		return err
	}
	h := f.happenings[happeningID]
	if h == nil {
		return facade4calendarius.ErrRecurringHappeningNotFound
	}
	h.active = false
	return calendariusmodels.ValidateRecurringRequestID(requestID)
}

func (f *referenceRecurringFacade) ListActiveRecurring(context.Context, string) ([]calendariusmodels.RecurringHappening, error) {
	return nil, nil
}

func (f *referenceRecurringFacade) GetRecurring(_ context.Context, spaceID, happeningID string) (calendariusmodels.RecurringHappening, error) {
	h := f.happenings[happeningID]
	if h == nil {
		return calendariusmodels.RecurringHappening{}, facade4calendarius.ErrRecurringHappeningNotFound
	}
	return calendariusmodels.RecurringHappening{
		ID: happeningID, SpaceID: spaceID, Active: h.active,
		Extensions: h.request.Extensions, Participants: h.request.Participants,
	}, nil
}

func (f *referenceRecurringFacade) SetExtensionData(context.Context, string, string, string, json.RawMessage) error {
	return nil
}

func (f *referenceRecurringFacade) ExpandOccurrences(_ context.Context, _, happeningID, _ string, from, to time.Time) ([]calendariusmodels.Occurrence, error) {
	h := f.happenings[happeningID]
	var result []calendariusmodels.Occurrence
	for day := from.AddDate(0, 0, -32); day.Before(to.AddDate(0, 0, 32)); day = day.AddDate(0, 0, 1) {
		for i, slot := range h.request.Slots {
			original, ok := referenceSlotStart(h.request, slot, day)
			if !ok {
				continue
			}
			key := h.request.SlotID(i) + "/" + day.Format(time.DateOnly)
			start, moved := h.moved[key]
			if !moved {
				start = original
			}
			if start.Before(from) || !start.Before(to) {
				continue
			}
			result = append(result, calendariusmodels.Occurrence{
				Start: start, End: start.Add(time.Duration(slot.DurationMinutes) * time.Minute),
				HappeningID: happeningID, SlotID: h.request.SlotID(i), LocalDate: day.Format(time.DateOnly),
				OriginalStart: original, Cancelled: h.cancelled[key], Rescheduled: moved,
			})
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Start.Before(result[j].Start) })
	return result, nil
}

func TestRunRecurringScheduleFacadeConformance(t *testing.T) {
	RunRecurringScheduleFacadeConformance(t, func(*testing.T) facade4calendarius.RecurringScheduleFacade {
		return &referenceRecurringFacade{happenings: map[string]*referenceRecurring{}}
	})
}

func TestRecurringScheduleErrorSentinelsAreDistinct(t *testing.T) {
	sentinels := []error{
		facade4calendarius.ErrRecurringHappeningNotFound,
		facade4calendarius.ErrRecurringHappeningClosed,
		facade4calendarius.ErrOccurrenceNotFound,
	}
	for i, sentinel := range sentinels {
		for j, other := range sentinels {
			if i != j && errors.Is(sentinel, other) {
				t.Fatalf("sentinels %d and %d overlap", i, j)
			}
		}
	}
}
