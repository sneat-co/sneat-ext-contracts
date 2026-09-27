package listusmodels

import (
	"testing"

	"github.com/sneat-co/sneat-go-core/coretypes"
)

func TestSourceTodoSpecValidation(t *testing.T) {
	valid := SourceTodoSpec{
		SpaceID: "family1", ListID: "todo", Purpose: "debt-due", Title: "Pay electricity bill",
		Source:       coretypes.NewFullItemRef("debtus", "sourceObligations", "family1", "bill1"),
		DueHappening: coretypes.NewFullItemRef("calendarius", "happenings", "family1", "due1"),
		State:        SourceTodoActive, DueTaskRevision: 1, CompletionActionID: "open-source", CompletionDisposition: SourceTodoNavigate,
	}
	if err := valid.Validate(); err != nil {
		t.Fatal(err)
	}

	validRequiresInput := valid
	validRequiresInput.CompletionDisposition = SourceTodoRequiresInput
	validRequiresInput.State = SourceTodoCompleted
	if err := validRequiresInput.Validate(); err != nil {
		t.Fatalf("validRequiresInput: %v", err)
	}

	for name, mutate := range map[string]func(*SourceTodoSpec){
		"invalid spaceID":        func(v *SourceTodoSpec) { v.SpaceID = "" },
		"empty listID":           func(v *SourceTodoSpec) { v.ListID = "" },
		"untrimmed listID":       func(v *SourceTodoSpec) { v.ListID = "  todo" },
		"invalid record listID":  func(v *SourceTodoSpec) { v.ListID = "todo list" },
		"empty title":            func(v *SourceTodoSpec) { v.Title = "" },
		"invalid source":         func(v *SourceTodoSpec) { v.Source = coretypes.ItemRef{} },
		"invalid dueHappening":   func(v *SourceTodoSpec) { v.DueHappening = coretypes.ItemRef{} },
		"zero revision":          func(v *SourceTodoSpec) { v.DueTaskRevision = 0 },
		"overflow revision":      func(v *SourceTodoSpec) { v.DueTaskRevision = SourceTodoMaxSafeInteger + 1 },
		"invalid state":          func(v *SourceTodoSpec) { v.State = "unknown" },
		"invalid disposition":    func(v *SourceTodoSpec) { v.CompletionDisposition = "execute" },
	} {
		t.Run(name, func(t *testing.T) {
			v := valid
			mutate(&v)
			if err := v.Validate(); err == nil {
				t.Fatalf("expected validation error for %s", name)
			}
		})
	}
}

func TestSourceTodoPlanViewValidation(t *testing.T) {
	valid := SourceTodoPlanView{
		ListID:       "todo",
		ItemID:       "item1",
		Item:         coretypes.NewFullItemRef("listus", "items", "family1", "item1"),
		DueHappening: coretypes.NewFullItemRef("calendarius", "happenings", "family1", "due1"),
		State:        SourceTodoActive,
	}
	if err := valid.Validate(); err != nil {
		t.Fatal(err)
	}

	validCompleted := valid
	validCompleted.State = SourceTodoCompleted
	if err := validCompleted.Validate(); err != nil {
		t.Fatal(err)
	}

	validCanceled := valid
	validCanceled.State = SourceTodoCanceled
	if err := validCanceled.Validate(); err != nil {
		t.Fatal(err)
	}

	for name, mutate := range map[string]func(*SourceTodoPlanView){
		"empty listID":         func(v *SourceTodoPlanView) { v.ListID = "" },
		"empty itemID":         func(v *SourceTodoPlanView) { v.ItemID = "" },
		"invalid item":         func(v *SourceTodoPlanView) { v.Item = coretypes.ItemRef{} },
		"invalid dueHappening": func(v *SourceTodoPlanView) { v.DueHappening = coretypes.ItemRef{} },
		"invalid state":        func(v *SourceTodoPlanView) { v.State = "invalid" },
	} {
		t.Run(name, func(t *testing.T) {
			v := valid
			mutate(&v)
			if err := v.Validate(); err == nil {
				t.Fatalf("expected validation error for %s", name)
			}
		})
	}
}

