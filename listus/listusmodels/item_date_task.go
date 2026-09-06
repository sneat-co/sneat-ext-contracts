package listusmodels

import (
	"fmt"
	"strings"
	"time"

	"github.com/sneat-co/sneat-go-core/coretypes"
)

type ListItemDateTaskLink struct {
	Happening coretypes.ItemRef `json:"happening"`
	Source    coretypes.ItemRef `json:"source"`
	Purpose   string            `json:"purpose"`
	Revision  int64             `json:"revision"`
}

type SaveListItemDateTaskRequest struct {
	SpaceID              coretypes.SpaceID `json:"spaceID"`
	ListID               string            `json:"listID"`
	ItemID               string            `json:"itemID"`
	OperationID          string            `json:"operationID"`
	ExpectedTaskRevision int64             `json:"expectedTaskRevision"`
	Title                string            `json:"title"`
	DueDate              string            `json:"dueDate,omitempty"`
	State                SourceTodoState   `json:"state"`
}

func (v SaveListItemDateTaskRequest) Validate() error {
	if err := coretypes.ValidateSpaceID(v.SpaceID); err != nil {
		return fmt.Errorf("spaceID: %w", err)
	}
	for name, value := range map[string]string{"listID": v.ListID, "itemID": v.ItemID, "operationID": v.OperationID} {
		if value == "" || value != strings.TrimSpace(value) || len(value) > 100 {
			return fmt.Errorf("%s must be trimmed and 1..100 bytes", name)
		}
	}
	if v.Title == "" || v.Title != strings.TrimSpace(v.Title) || len(v.Title) > 100 {
		return fmt.Errorf("title must be trimmed and 1..100 bytes")
	}
	if v.ExpectedTaskRevision < 0 || v.ExpectedTaskRevision > SourceTodoMaxSafeInteger {
		return fmt.Errorf("expectedTaskRevision is outside the safe integer range")
	}
	if v.DueDate != "" {
		parsed, err := time.Parse("2006-01-02", v.DueDate)
		if err != nil || parsed.Format("2006-01-02") != v.DueDate {
			return fmt.Errorf("invalid dueDate")
		}
	}
	switch v.State {
	case SourceTodoActive:
		if v.DueDate == "" {
			return fmt.Errorf("active date task requires dueDate")
		}
	case SourceTodoCompleted, SourceTodoCanceled:
	default:
		return fmt.Errorf("unsupported state %q", v.State)
	}
	return nil
}

type SaveListItemDateTaskResponse struct {
	ItemID   string               `json:"itemID"`
	DateTask ListItemDateTaskLink `json:"dateTask"`
}
