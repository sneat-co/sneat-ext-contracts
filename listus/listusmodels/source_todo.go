package listusmodels

import (
	"fmt"
	"strings"

	"github.com/sneat-co/sneat-go-core/coretypes"
	corevalidate "github.com/sneat-co/sneat-go-core/validate"
)

const SourceTodoMaxSafeInteger int64 = 9_007_199_254_740_991

type SourceTodoState string

const (
	SourceTodoActive    SourceTodoState = "active"
	SourceTodoCompleted SourceTodoState = "completed"
	SourceTodoCanceled  SourceTodoState = "canceled"
)

type SourceTodoCompletionDisposition string

const (
	SourceTodoNavigate      SourceTodoCompletionDisposition = "navigate"
	SourceTodoRequiresInput SourceTodoCompletionDisposition = "requires_input"
)

// SourceTodoSpec is server-authored intent for one embedded Listus item. Source
// and DueHappening are standard Sneat Linkage references; the stable source and
// purpose pair is the idempotent identity within the selected list.
type SourceTodoSpec struct {
	SpaceID               coretypes.SpaceID               `json:"spaceID"`
	ListID                string                          `json:"listID"`
	Source                coretypes.ItemRef               `json:"source"`
	Purpose               string                          `json:"purpose"`
	Title                 string                          `json:"title"`
	State                 SourceTodoState                 `json:"state"`
	DueHappening          coretypes.ItemRef               `json:"dueHappening"`
	DueTaskRevision       int64                           `json:"dueTaskRevision"`
	CompletionActionID    string                          `json:"completionActionID"`
	CompletionDisposition SourceTodoCompletionDisposition `json:"completionDisposition"`
}

func (v SourceTodoSpec) Validate() error {
	if err := coretypes.ValidateSpaceID(v.SpaceID); err != nil {
		return fmt.Errorf("spaceID: %w", err)
	}
	for name, value := range map[string]string{"listID": v.ListID, "purpose": v.Purpose, "completionActionID": v.CompletionActionID} {
		if value == "" || value != strings.TrimSpace(value) || len(value) > 100 {
			return fmt.Errorf("%s must be trimmed and 1..100 bytes", name)
		}
		if err := corevalidate.RecordID(value); err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}
	}
	if v.Title == "" || v.Title != strings.TrimSpace(v.Title) || len(v.Title) > 100 {
		return fmt.Errorf("title must be trimmed and 1..100 bytes")
	}
	if err := v.Source.Validate(); err != nil {
		return fmt.Errorf("source: %w", err)
	}
	if err := v.DueHappening.Validate(); err != nil {
		return fmt.Errorf("dueHappening: %w", err)
	}
	if v.DueTaskRevision < 1 || v.DueTaskRevision > SourceTodoMaxSafeInteger {
		return fmt.Errorf("dueTaskRevision must be a positive safe integer")
	}
	switch v.State {
	case SourceTodoActive, SourceTodoCompleted, SourceTodoCanceled:
	default:
		return fmt.Errorf("unsupported state %q", v.State)
	}
	switch v.CompletionDisposition {
	case SourceTodoNavigate, SourceTodoRequiresInput:
	default:
		return fmt.Errorf("unsupported completionDisposition %q", v.CompletionDisposition)
	}
	return nil
}

type SourceTodoPlanView struct {
	ListID       string            `json:"listID"`
	ItemID       string            `json:"itemID"`
	Item         coretypes.ItemRef `json:"item"`
	DueHappening coretypes.ItemRef `json:"dueHappening"`
	State        SourceTodoState   `json:"state"`
}

func (v SourceTodoPlanView) Validate() error {
	if v.ListID == "" || v.ItemID == "" {
		return fmt.Errorf("listID and itemID are required")
	}
	if err := v.Item.Validate(); err != nil {
		return fmt.Errorf("item: %w", err)
	}
	if err := v.DueHappening.Validate(); err != nil {
		return fmt.Errorf("dueHappening: %w", err)
	}
	switch v.State {
	case SourceTodoActive, SourceTodoCompleted, SourceTodoCanceled:
		return nil
	default:
		return fmt.Errorf("unsupported state %q", v.State)
	}
}
