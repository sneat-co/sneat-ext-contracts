package listusmodels

import (
	"testing"

	"github.com/sneat-co/sneat-go-core/coretypes"
)

func TestSourceTodoSpecValidation(t *testing.T) {
	v := SourceTodoSpec{
		SpaceID: "family1", ListID: "todo", Purpose: "debt-due", Title: "Pay electricity bill",
		Source:       coretypes.NewFullItemRef("debtus", "sourceObligations", "family1", "bill1"),
		DueHappening: coretypes.NewFullItemRef("calendarius", "happenings", "family1", "due1"),
		State:        SourceTodoActive, DueTaskRevision: 1, CompletionActionID: "open-source", CompletionDisposition: SourceTodoNavigate,
	}
	if err := v.Validate(); err != nil {
		t.Fatal(err)
	}
	v.CompletionDisposition = "execute"
	if err := v.Validate(); err == nil {
		t.Fatal("accepted unimplemented execute disposition")
	}
}
