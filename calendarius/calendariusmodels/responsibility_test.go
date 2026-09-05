package calendariusmodels

import (
	"testing"
	"time"
)

func TestScheduledResponsibilitySpecValidation(t *testing.T) {
	spec := ScheduledResponsibilitySpec{Title: "Bins", TimeZone: "Europe/Dublin", FirstDate: "2026-09-07", Weekday: "mo", StartTime: "19:00", Assignment: ResponsibilityAssignmentPolicy{Mode: ResponsibilityAssignmentRotating, RosterContactIDs: []string{"alice", "bob"}}}
	if err := spec.Validate(); err != nil {
		t.Fatal(err)
	}
	spec.Assignment.RosterContactIDs = []string{"alice", "alice"}
	if err := spec.Validate(); err == nil {
		t.Fatal("expected duplicate roster rejection")
	}
}
func TestResponsibilityOccurrenceKeyIncludesSlot(t *testing.T) {
	a := ResponsibilityOccurrenceRef{HappeningID: "h1", SlotID: "weekly", Date: "2026-09-07", Start: time.Now()}
	b := a
	b.SlotID = "other"
	if a.Key() == b.Key() {
		t.Fatal("different slots collided")
	}
}
