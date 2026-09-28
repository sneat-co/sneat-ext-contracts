package calendariusmodels

import (
	"encoding/json"
	"fmt"
	"strings"
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

func TestResponsibility_AdditionalCoverage(t *testing.T) {
	// 1. ResponsibilityAssignmentPolicy
	// invalid mode
	p := ResponsibilityAssignmentPolicy{Mode: "invalid", RosterContactIDs: []string{"c1"}}
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for invalid mode")
	}
	// empty roster
	p = ResponsibilityAssignmentPolicy{Mode: ResponsibilityAssignmentFixed, RosterContactIDs: nil}
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for empty roster")
	}
	// roster > ResponsibilityRosterMax (50)
	tooMany := make([]string, ResponsibilityRosterMax+1)
	for i := range tooMany {
		tooMany[i] = fmt.Sprintf("c%d", i)
	}
	p = ResponsibilityAssignmentPolicy{Mode: ResponsibilityAssignmentRotating, RosterContactIDs: tooMany}
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for roster exceeding max")
	}
	// fixed with > 1 contacts
	p = ResponsibilityAssignmentPolicy{Mode: ResponsibilityAssignmentFixed, RosterContactIDs: []string{"c1", "c2"}}
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for fixed with >1 contacts")
	}
	// contact ID validation checks: untrimmed, empty, non-utf8, > 200 chars, contains @/:
	for _, badID := range []string{" c1", "c1 ", "", "\xff", strings.Repeat("a", ResponsibilityContactIDMaxLen+1), "user@host", "space/contact", "urn:contact"} {
		p = ResponsibilityAssignmentPolicy{Mode: ResponsibilityAssignmentFixed, RosterContactIDs: []string{badID}}
		if err := p.Validate(); err == nil {
			t.Fatalf("expected error for bad contact ID: %q", badID)
		}
	}

	// 2. ScheduledResponsibilitySpec
	validSpec := ScheduledResponsibilitySpec{
		Title: "Bins", Description: "Take bins out", TimeZone: "Europe/Dublin",
		FirstDate: "2026-09-07", Weekday: "mo", StartTime: "19:00", DurationMinutes: 30,
		Assignment: ResponsibilityAssignmentPolicy{Mode: ResponsibilityAssignmentFixed, RosterContactIDs: []string{"alice"}},
	}
	if err := validSpec.Validate(); err != nil {
		t.Fatalf("validSpec failed: %v", err)
	}
	// title checks: empty, untrimmed, > EventHappeningTitleMaxBytes (100)
	s := validSpec
	s.Title = "   "
	if err := s.Validate(); err == nil {
		t.Fatal("expected error for blank title")
	}
	s = validSpec
	s.Title = strings.Repeat("x", EventHappeningTitleMaxBytes+1)
	if err := s.Validate(); err == nil {
		t.Fatal("expected error for title exceeding max bytes")
	}
	// description > EventHappeningDescriptionMaxBytes (5000)
	s = validSpec
	s.Description = strings.Repeat("x", EventHappeningDescriptionMaxBytes+1)
	if err := s.Validate(); err == nil {
		t.Fatal("expected error for description exceeding max bytes")
	}
	// empty timeZone
	s = validSpec
	s.TimeZone = ""
	if err := s.Validate(); err == nil {
		t.Fatal("expected error for empty timezone")
	}
	// invalid timeZone location
	s = validSpec
	s.TimeZone = "Mars/Olympus"
	if err := s.Validate(); err == nil {
		t.Fatal("expected error for invalid timezone location")
	}
	// invalid firstDate
	s = validSpec
	s.FirstDate = "not-a-date"
	if err := s.Validate(); err == nil {
		t.Fatal("expected error for invalid firstDate")
	}
	// unknown weekday or mismatching weekday
	s = validSpec
	s.Weekday = "unknown"
	if err := s.Validate(); err == nil {
		t.Fatal("expected error for unknown weekday")
	}
	s = validSpec
	s.Weekday = "tu" // 2026-09-07 is Monday
	if err := s.Validate(); err == nil {
		t.Fatal("expected error for mismatching weekday")
	}
	// invalid startTime
	s = validSpec
	s.StartTime = "bad-time"
	if err := s.Validate(); err == nil {
		t.Fatal("expected error for invalid startTime")
	}
	// durationMinutes < 0 or > EventHappeningDurationMaxMinutes
	s = validSpec
	s.DurationMinutes = -1
	if err := s.Validate(); err == nil {
		t.Fatal("expected error for negative durationMinutes")
	}
	s = validSpec
	s.DurationMinutes = EventHappeningDurationMaxMinutes + 1
	if err := s.Validate(); err == nil {
		t.Fatal("expected error for durationMinutes exceeding max")
	}

	// 3. CreateScheduledResponsibilityRequest
	req := CreateScheduledResponsibilityRequest{
		RequestID: "req1",
		Spec:      validSpec,
	}
	if err := req.Validate(); err != nil {
		t.Fatalf("req.Validate failed: %v", err)
	}
	// invalid requestID
	badReq := req
	badReq.RequestID = ""
	if err := badReq.Validate(); err == nil {
		t.Fatal("expected error for empty requestID")
	}
	// invalid spec
	badReq = req
	badReq.Spec.Title = ""
	if err := badReq.Validate(); err == nil {
		t.Fatal("expected error for invalid spec in request")
	}
	// with valid HappeningFields
	hfReq := req
	hfReq.HappeningFields = &ResponsibilityHappeningFields{}
	if err := hfReq.Validate(); err != nil {
		t.Fatalf("hfReq.Validate failed: %v", err)
	}
	// with invalid HappeningFields
	badHfReq := req
	badHfReq.HappeningFields = &ResponsibilityHappeningFields{Related: json.RawMessage(`{`)}
	if err := badHfReq.Validate(); err == nil {
		t.Fatal("expected error for invalid HappeningFields in request")
	}

	// 4. CompleteResponsibilityOccurrenceRequest
	start := time.Date(2026, 9, 7, 19, 0, 0, 0, time.UTC)
	compReq := CompleteResponsibilityOccurrenceRequest{
		RequestID: "req1",
		Ref: ResponsibilityOccurrenceRef{
			HappeningID: "h1",
			SlotID:      "s1",
			Date:        "2026-09-07",
			Start:       start,
			End:         start.Add(time.Hour),
		},
	}
	if err := compReq.Validate(); err != nil {
		t.Fatalf("compReq.Validate failed: %v", err)
	}
	// invalid requestID
	cReq := compReq
	cReq.RequestID = ""
	if err := cReq.Validate(); err == nil {
		t.Fatal("expected error for invalid requestID")
	}
	// missing happeningID or slotID
	cReq = compReq
	cReq.Ref.HappeningID = ""
	if err := cReq.Validate(); err == nil {
		t.Fatal("expected error for empty happeningID")
	}
	cReq = compReq
	cReq.Ref.SlotID = ""
	if err := cReq.Validate(); err == nil {
		t.Fatal("expected error for empty slotID")
	}
	// invalid date
	cReq = compReq
	cReq.Ref.Date = "invalid-date"
	if err := cReq.Validate(); err == nil {
		t.Fatal("expected error for invalid date")
	}
	// zero start, zero end, or end <= start
	cReq = compReq
	cReq.Ref.Start = time.Time{}
	if err := cReq.Validate(); err == nil {
		t.Fatal("expected error for zero start")
	}
	cReq = compReq
	cReq.Ref.End = time.Time{}
	if err := cReq.Validate(); err == nil {
		t.Fatal("expected error for zero end")
	}
	cReq = compReq
	cReq.Ref.End = start.Add(-time.Hour)
	if err := cReq.Validate(); err == nil {
		t.Fatal("expected error for end before start")
	}

	// 5. ValidateResponsibilityRequestID
	for _, badID := range []string{"", " leading", "trailing ", "\xff", strings.Repeat("x", ResponsibilityRequestIDMaxLen+1)} {
		if err := ValidateResponsibilityRequestID(badID); err == nil {
			t.Fatalf("expected error for bad requestID: %q", badID)
		}
	}
	if err := ValidateResponsibilityRequestID("valid-request-id-123"); err != nil {
		t.Fatalf("unexpected error for valid requestID: %v", err)
	}
}

