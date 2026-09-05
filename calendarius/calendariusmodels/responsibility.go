package calendariusmodels

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	ResponsibilityRosterMax       = 50
	ResponsibilityContactIDMaxLen = 200
	ResponsibilityRequestIDMaxLen = 200
)

type ResponsibilityAssignmentMode string

const (
	ResponsibilityAssignmentFixed    ResponsibilityAssignmentMode = "fixed"
	ResponsibilityAssignmentRotating ResponsibilityAssignmentMode = "rotating"
)

type ResponsibilityAssignmentPolicy struct {
	Mode             ResponsibilityAssignmentMode `json:"mode" firestore:"mode"`
	RosterContactIDs []string                     `json:"rosterContactIDs" firestore:"rosterContactIDs"`
}

func (v ResponsibilityAssignmentPolicy) Validate() error {
	if v.Mode != ResponsibilityAssignmentFixed && v.Mode != ResponsibilityAssignmentRotating {
		return fmt.Errorf("assignment.mode must be fixed or rotating")
	}
	if len(v.RosterContactIDs) == 0 || len(v.RosterContactIDs) > ResponsibilityRosterMax {
		return fmt.Errorf("assignment.rosterContactIDs must contain 1..%d contacts", ResponsibilityRosterMax)
	}
	if v.Mode == ResponsibilityAssignmentFixed && len(v.RosterContactIDs) != 1 {
		return fmt.Errorf("fixed assignment requires exactly one contact")
	}
	seen := make(map[string]struct{}, len(v.RosterContactIDs))
	for i, id := range v.RosterContactIDs {
		if strings.TrimSpace(id) != id || id == "" || !utf8.ValidString(id) || len(id) > ResponsibilityContactIDMaxLen || strings.ContainsAny(id, "@/:") {
			return fmt.Errorf("assignment.rosterContactIDs[%d] is not a valid same-Space contact ID", i)
		}
		if _, ok := seen[id]; ok {
			return fmt.Errorf("assignment.rosterContactIDs[%d] is duplicated", i)
		}
		seen[id] = struct{}{}
	}
	return nil
}

type ScheduledResponsibilitySpec struct {
	Title           string                         `json:"title" firestore:"title"`
	Description     string                         `json:"description,omitempty" firestore:"description,omitempty"`
	TimeZone        string                         `json:"timeZone" firestore:"timeZone"`
	FirstDate       string                         `json:"firstDate" firestore:"firstDate"`
	Weekday         string                         `json:"weekday" firestore:"weekday"`
	StartTime       string                         `json:"startTime" firestore:"startTime"`
	DurationMinutes int                            `json:"durationMinutes,omitempty" firestore:"durationMinutes,omitempty"`
	Assignment      ResponsibilityAssignmentPolicy `json:"assignment" firestore:"assignment"`
}

func (v ScheduledResponsibilitySpec) Validate() error {
	if strings.TrimSpace(v.Title) == "" || len(v.Title) > EventHappeningTitleMaxBytes {
		return fmt.Errorf("title is required and must be at most %d bytes", EventHappeningTitleMaxBytes)
	}
	if len(v.Description) > EventHappeningDescriptionMaxBytes {
		return fmt.Errorf("description exceeds %d bytes", EventHappeningDescriptionMaxBytes)
	}
	if v.TimeZone == "" || v.TimeZone == "Local" {
		return fmt.Errorf("timeZone must be an IANA timezone")
	}
	if _, err := time.LoadLocation(v.TimeZone); err != nil {
		return fmt.Errorf("timeZone must be an IANA timezone")
	}
	first, err := time.Parse(time.DateOnly, v.FirstDate)
	if err != nil {
		return fmt.Errorf("firstDate must be YYYY-MM-DD")
	}
	weekday, ok := responsibilityWeekdays[v.Weekday]
	if !ok || first.Weekday() != weekday {
		return fmt.Errorf("weekday must match firstDate")
	}
	if _, err = time.Parse("15:04", v.StartTime); err != nil {
		return fmt.Errorf("startTime must be HH:MM")
	}
	if v.DurationMinutes < 0 || v.DurationMinutes > EventHappeningDurationMaxMinutes {
		return fmt.Errorf("durationMinutes is outside the finite bound")
	}
	return v.Assignment.Validate()
}

var responsibilityWeekdays = map[string]time.Weekday{"mo": time.Monday, "tu": time.Tuesday, "we": time.Wednesday, "th": time.Thursday, "fr": time.Friday, "sa": time.Saturday, "su": time.Sunday}

type CreateScheduledResponsibilityRequest struct {
	RequestID       string                         `json:"requestID"`
	Spec            ScheduledResponsibilitySpec    `json:"spec"`
	HappeningFields *ResponsibilityHappeningFields `json:"happeningFields,omitempty"`
}

// ResponsibilityHappeningFields carries only generic extension composition.
// Schedule and responsibility fields remain derived from Spec by Calendarius.
type ResponsibilityHappeningFields struct {
	Ext     map[string]json.RawMessage `json:"ext,omitempty"`
	Related json.RawMessage            `json:"related,omitempty"`
}

func (v CreateScheduledResponsibilityRequest) Validate() error {
	if err := ValidateResponsibilityRequestID(v.RequestID); err != nil {
		return err
	}
	if err := v.Spec.Validate(); err != nil {
		return err
	}
	if v.HappeningFields != nil {
		return v.HappeningFields.Validate()
	}
	return nil
}

func (v ResponsibilityHappeningFields) Validate() error {
	for id, payload := range v.Ext {
		var object map[string]json.RawMessage
		if strings.TrimSpace(id) != id || id == "" || !json.Valid(payload) || json.Unmarshal(payload, &object) != nil || object == nil {
			return fmt.Errorf("happeningFields.ext contains an invalid extension payload")
		}
	}
	if len(v.Related) > 0 {
		var object map[string]json.RawMessage
		if !json.Valid(v.Related) || json.Unmarshal(v.Related, &object) != nil || object == nil {
			return fmt.Errorf("happeningFields.related must be a JSON object")
		}
	}
	return nil
}

type ResponsibilityOccurrenceRef struct {
	HappeningID string    `json:"happeningID" firestore:"happeningID"`
	SlotID      string    `json:"slotID" firestore:"slotID"`
	Date        string    `json:"date" firestore:"date"`
	Start       time.Time `json:"start" firestore:"start"`
	End         time.Time `json:"end,omitempty" firestore:"end,omitempty"`
}

func (v ResponsibilityOccurrenceRef) Key() string {
	h := sha256.New()
	for _, part := range []string{v.HappeningID, v.SlotID, v.Date} {
		var size [8]byte
		binary.BigEndian.PutUint64(size[:], uint64(len(part)))
		_, _ = h.Write(size[:])
		_, _ = h.Write([]byte(part))
	}
	sum := h.Sum(nil)
	return hex.EncodeToString(sum)
}

type ResponsibilityCompletion struct {
	OccurrenceKey     string    `json:"occurrenceKey" firestore:"occurrenceKey"`
	HappeningID       string    `json:"happeningID" firestore:"happeningID"`
	SlotID            string    `json:"slotID" firestore:"slotID"`
	Date              string    `json:"date" firestore:"date"`
	AssignedContactID string    `json:"assignedContactID" firestore:"assignedContactID"`
	CompletedBy       string    `json:"completedBy" firestore:"completedBy"`
	CompletedAt       time.Time `json:"completedAt" firestore:"completedAt"`
}
type ResponsibilityOccurrence struct {
	Ref               ResponsibilityOccurrenceRef `json:"ref"`
	AssignedContactID string                      `json:"assignedContactID,omitempty"`
	NeedsReassignment bool                        `json:"needsReassignment,omitempty"`
	Completion        *ResponsibilityCompletion   `json:"completion,omitempty"`
}
type ScheduledResponsibility struct {
	ID   string                      `json:"id"`
	Spec ScheduledResponsibilitySpec `json:"spec"`
}
type CompleteResponsibilityOccurrenceRequest struct {
	RequestID string                      `json:"requestID"`
	Ref       ResponsibilityOccurrenceRef `json:"ref"`
}

func (v CompleteResponsibilityOccurrenceRequest) Validate() error {
	if err := ValidateResponsibilityRequestID(v.RequestID); err != nil {
		return err
	}
	if v.Ref.HappeningID == "" || v.Ref.SlotID == "" {
		return fmt.Errorf("occurrence happeningID and slotID are required")
	}
	if _, err := time.Parse(time.DateOnly, v.Ref.Date); err != nil {
		return fmt.Errorf("occurrence date must be YYYY-MM-DD")
	}
	if v.Ref.Start.IsZero() || v.Ref.End.IsZero() || !v.Ref.End.After(v.Ref.Start) {
		return fmt.Errorf("occurrence start and end are required and must be ordered")
	}
	return nil
}

type ResponsibilityMutationDisposition string

const (
	ResponsibilityCreated   ResponsibilityMutationDisposition = "created"
	ResponsibilityCompleted ResponsibilityMutationDisposition = "completed"
	ResponsibilityUnchanged ResponsibilityMutationDisposition = "unchanged"
	ResponsibilityReused    ResponsibilityMutationDisposition = "reused"
)

type ScheduledResponsibilityMutation struct {
	Responsibility ScheduledResponsibility           `json:"responsibility"`
	Disposition    ResponsibilityMutationDisposition `json:"disposition"`
}
type ResponsibilityCompletionMutation struct {
	Completion  ResponsibilityCompletion          `json:"completion"`
	Disposition ResponsibilityMutationDisposition `json:"disposition"`
}

func ValidateResponsibilityRequestID(id string) error {
	if strings.TrimSpace(id) != id || id == "" || len(id) > ResponsibilityRequestIDMaxLen || !utf8.ValidString(id) {
		return fmt.Errorf("requestID is required and must be valid UTF-8 with at most %d bytes", ResponsibilityRequestIDMaxLen)
	}
	return nil
}
