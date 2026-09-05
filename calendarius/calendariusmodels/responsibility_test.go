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
func TestScheduledResponsibilityTimeZoneValidation(t *testing.T) {
	for _, zone := range []string{"UTC", "Europe/Dublin"} {
		spec := ScheduledResponsibilitySpec{Title: "Bins", TimeZone: zone, FirstDate: "2026-09-07", Weekday: "mo", StartTime: "19:00", Assignment: ResponsibilityAssignmentPolicy{Mode: ResponsibilityAssignmentFixed, RosterContactIDs: []string{"alice"}}}
		if err := spec.Validate(); err != nil {
			t.Fatalf("%s: %v", zone, err)
		}
	}
	spec := ScheduledResponsibilitySpec{Title: "Bins", TimeZone: "Local", FirstDate: "2026-09-07", Weekday: "mo", StartTime: "19:00", Assignment: ResponsibilityAssignmentPolicy{Mode: ResponsibilityAssignmentFixed, RosterContactIDs: []string{"alice"}}}
	if err := spec.Validate(); err == nil {
		t.Fatal("Local is process-dependent, not an accepted IANA timezone")
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
func TestResponsibilityOccurrenceKeyLengthPrefixesTupleParts(t *testing.T) {
	a := ResponsibilityOccurrenceRef{HappeningID: "a\x00b", SlotID: "c", Date: "2026-09-07"}
	b := ResponsibilityOccurrenceRef{HappeningID: "a", SlotID: "b\x00c", Date: "2026-09-07"}
	if a.Key() == b.Key() {
		t.Fatal("distinct tuple parts collided")
	}
}
