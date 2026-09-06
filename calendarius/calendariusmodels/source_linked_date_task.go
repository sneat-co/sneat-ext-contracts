package calendariusmodels

import (
	"fmt"
	"strings"
	"time"

	"github.com/sneat-co/sneat-go-core/coretypes"
	corevalidate "github.com/sneat-co/sneat-go-core/validate"
)

const SourceLinkedDateTaskMaxSafeInteger int64 = 9_007_199_254_740_991

type SourceLinkedDateTaskState string

const (
	SourceLinkedDateTaskActive    SourceLinkedDateTaskState = "active"
	SourceLinkedDateTaskCompleted SourceLinkedDateTaskState = "completed"
	SourceLinkedDateTaskCanceled  SourceLinkedDateTaskState = "canceled"
)

type SourceLinkedDateTaskActionDisposition string

const (
	SourceLinkedDateTaskNavigate      SourceLinkedDateTaskActionDisposition = "navigate"
	SourceLinkedDateTaskRequiresInput SourceLinkedDateTaskActionDisposition = "requires_input"
	SourceLinkedDateTaskExecute       SourceLinkedDateTaskActionDisposition = "execute"
)

// SourceLinkedDateTaskRef is the owner-qualified identity of a source item.
// LineID distinguishes independently dated obligations within one record.
type SourceLinkedDateTaskRef struct {
	Namespace    string `json:"namespace"`
	OwnerSpaceID string `json:"ownerSpaceID"`
	RecordID     string `json:"recordID"`
	LineID       string `json:"lineID"`
}

// SourceLinkedDateTaskRelatedRef is persisted through standard Sneat Linkage.
// Stable source identity remains separate from this navigational relation.
type SourceLinkedDateTaskRelatedRef struct {
	ItemRef coretypes.ItemRef `json:"itemRef"`
	Role    string            `json:"role"`
}

func (v SourceLinkedDateTaskRelatedRef) Validate() error {
	if err := v.ItemRef.Validate(); err != nil {
		return fmt.Errorf("itemRef: %w", err)
	}
	for name, value := range map[string]string{"role": v.Role} {
		if value != strings.TrimSpace(value) || value == "" || len(value) > 100 {
			return fmt.Errorf("related %s must be a trimmed non-empty identifier of at most 100 bytes", name)
		}
	}
	return nil
}

func (v SourceLinkedDateTaskRef) Validate() error {
	if err := coretypes.ValidateSpaceID(coretypes.SpaceID(v.OwnerSpaceID)); err != nil {
		return fmt.Errorf("ownerSpaceID: %w", err)
	}
	for name, value := range map[string]string{"namespace": v.Namespace, "recordID": v.RecordID, "lineID": v.LineID} {
		if value != strings.TrimSpace(value) || value == "" || len(value) > 100 {
			return fmt.Errorf("%s must be a trimmed non-empty identifier of at most 100 bytes", name)
		}
	}
	if err := corevalidate.RecordID(v.RecordID); err != nil {
		return fmt.Errorf("recordID: %w", err)
	}
	if err := corevalidate.RecordID(v.LineID); err != nil {
		return fmt.Errorf("lineID: %w", err)
	}
	return nil
}

type SourceLinkedDateTask struct {
	HappeningID       string                                `json:"happeningID"`
	Revision          int64                                 `json:"revision"`
	Source            SourceLinkedDateTaskRef               `json:"source"`
	Title             string                                `json:"title"`
	DueDate           string                                `json:"dueDate,omitempty"`
	State             SourceLinkedDateTaskState             `json:"state"`
	ActionID          string                                `json:"actionID,omitempty"`
	ActionDisposition SourceLinkedDateTaskActionDisposition `json:"actionDisposition,omitempty"`
}

func (v SourceLinkedDateTask) Validate() error {
	if err := v.Source.Validate(); err != nil {
		return err
	}
	if v.HappeningID == "" || v.Revision < 1 || v.Revision > SourceLinkedDateTaskMaxSafeInteger {
		return fmt.Errorf("happeningID and positive revision are required")
	}
	if v.Title != strings.TrimSpace(v.Title) || v.Title == "" || len(v.Title) > 100 {
		return fmt.Errorf("title must be trimmed and 1..100 bytes")
	}
	if v.DueDate != "" {
		parsed, err := time.Parse("2006-01-02", v.DueDate)
		if err != nil || parsed.Format("2006-01-02") != v.DueDate {
			return fmt.Errorf("invalid dueDate: %w", err)
		}
	}
	switch v.State {
	case SourceLinkedDateTaskActive, SourceLinkedDateTaskCompleted, SourceLinkedDateTaskCanceled:
	default:
		return fmt.Errorf("unsupported state %q", v.State)
	}
	if v.State == SourceLinkedDateTaskActive && v.DueDate == "" {
		return fmt.Errorf("active task requires dueDate")
	}
	if v.ActionDisposition == "" {
		if v.ActionID != "" {
			return fmt.Errorf("actionID requires actionDisposition")
		}
		return nil
	}
	if v.ActionID == "" {
		return fmt.Errorf("actionDisposition requires actionID")
	}
	if v.ActionID != strings.TrimSpace(v.ActionID) || len(v.ActionID) > 100 {
		return fmt.Errorf("actionID must be trimmed and at most 100 bytes")
	}
	switch v.ActionDisposition {
	case SourceLinkedDateTaskNavigate, SourceLinkedDateTaskRequiresInput:
	default:
		return fmt.Errorf("unsupported actionDisposition %q", v.ActionDisposition)
	}
	return nil
}

type MutateSourceLinkedDateTaskRequest struct {
	OperationID       string                                `json:"operationID"`
	ExpectedRevision  int64                                 `json:"expectedRevision"`
	Source            SourceLinkedDateTaskRef               `json:"source"`
	Title             string                                `json:"title"`
	DueDate           string                                `json:"dueDate,omitempty"`
	State             SourceLinkedDateTaskState             `json:"state"`
	ActionID          string                                `json:"actionID,omitempty"`
	ActionDisposition SourceLinkedDateTaskActionDisposition `json:"actionDisposition,omitempty"`
	Related           []SourceLinkedDateTaskRelatedRef      `json:"related"`
}

func (v MutateSourceLinkedDateTaskRequest) Validate() error {
	if v.OperationID != strings.TrimSpace(v.OperationID) || v.OperationID == "" || len(v.OperationID) > 100 {
		return fmt.Errorf("operationID is required and must be at most 100 bytes")
	}
	if v.ExpectedRevision < 0 || v.ExpectedRevision > SourceLinkedDateTaskMaxSafeInteger {
		return fmt.Errorf("expectedRevision must not be negative")
	}
	probe := SourceLinkedDateTask{HappeningID: "planned", Revision: 1, Source: v.Source, Title: v.Title, DueDate: v.DueDate, State: v.State, ActionID: v.ActionID, ActionDisposition: v.ActionDisposition}
	if err := probe.Validate(); err != nil {
		return err
	}
	if v.Related == nil {
		return fmt.Errorf("related must be an array")
	}
	if len(v.Related) < 1 || len(v.Related) > 4 {
		return fmt.Errorf("related must have 1 to 4 items")
	}
	seen := map[string]bool{}
	for i, ref := range v.Related {
		if err := ref.Validate(); err != nil {
			return fmt.Errorf("related[%d]: %w", i, err)
		}
		key := ref.ItemRef.ID()
		if seen[key] {
			return fmt.Errorf("related[%d] duplicates an earlier reference", i)
		}
		seen[key] = true
	}
	return nil
}

type SourceLinkedDateTaskMutation struct {
	Task        SourceLinkedDateTask `json:"task"`
	Disposition string               `json:"disposition"`
}
