package calendariusmodels

import (
	"encoding/json"
	"testing"
	"time"
)

func TestResponsibilityHappeningFieldsValidation(t *testing.T) {
	valid := ResponsibilityHappeningFields{
		Ext:     map[string]json.RawMessage{"listus": json.RawMessage(`{"listTemplate":{"sourceListID":"do!regular"}}`)},
		Related: json.RawMessage(`{"listus":{"lists":{"do!regular":{}}}}`),
	}
	if err := valid.Validate(); err != nil {
		t.Fatal(err)
	}
	for name, fields := range map[string]ResponsibilityHappeningFields{
		"malformed related": {Related: json.RawMessage(`{`)},
		"null related":      {Related: json.RawMessage(`null`)},
		"array related":     {Related: json.RawMessage(`[]`)},
		"invalid ext json":  {Ext: map[string]json.RawMessage{"listus": json.RawMessage(`{`)}},
		"null ext payload":  {Ext: map[string]json.RawMessage{"listus": json.RawMessage(`null`)}},
		"invalid ext key":   {Ext: map[string]json.RawMessage{" listus": json.RawMessage(`{}`)}},
	} {
		t.Run(name, func(t *testing.T) {
			if err := fields.Validate(); err == nil {
				t.Fatal("invalid happening fields were accepted")
			}
		})
	}
}

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
