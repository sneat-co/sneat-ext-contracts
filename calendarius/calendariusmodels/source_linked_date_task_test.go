package calendariusmodels

import (
	"reflect"
	"strings"
	"testing"

	"github.com/sneat-co/sneat-go-core/coretypes"
)

func TestSourceLinkedDateTaskPersistenceTagsMatchWireNames(t *testing.T) {
	for _, value := range []any{
		SourceLinkedDateTaskRef{},
		SourceLinkedDateTaskRelatedRef{},
		SourceLinkedDateTask{},
		SourceLinkedDateTaskMutation{},
	} {
		typeOf := reflect.TypeOf(value)
		for i := 0; i < typeOf.NumField(); i++ {
			field := typeOf.Field(i)
			jsonName := strings.Split(field.Tag.Get("json"), ",")[0]
			firestoreName := strings.Split(field.Tag.Get("firestore"), ",")[0]
			if jsonName == "" || firestoreName != jsonName {
				t.Fatalf("%s.%s firestore tag %q does not match JSON tag %q", typeOf.Name(), field.Name, firestoreName, jsonName)
			}
		}
	}
}

func TestLinkedDateTaskPreservesCompositeFinancialLineIdentity(t *testing.T) {
	v := validLinkedTaskRequest()
	v.Source.LineID = "shared/bill/" + strings.Repeat("a", 500)
	if err := v.Validate(); err != nil {
		t.Fatalf("valid opaque financial line rejected: %v", err)
	}
	for _, invalid := range []string{strings.Repeat("a", 513), "line\n1", " line", string([]byte{0xff})} {
		v.Source.LineID = invalid
		if err := v.Validate(); err == nil {
			t.Fatalf("invalid source line accepted: %q", invalid)
		}
	}
}

func validLinkedTaskRequest() MutateSourceLinkedDateTaskRequest {
	return MutateSourceLinkedDateTaskRequest{OperationID: "op1", Source: SourceLinkedDateTaskRef{Namespace: "debtus", OwnerSpaceID: "family1", RecordID: "debt1", LineID: "payment1"}, Title: "Pay electricity bill", DueDate: "2026-09-30", State: SourceLinkedDateTaskActive, ActionID: "open-source", ActionDisposition: SourceLinkedDateTaskNavigate, Related: []SourceLinkedDateTaskRelatedRef{{ItemRef: coretypes.ItemRef{ExtID: "debtus", Collection: "obligations", ItemID: "debt1"}, Role: "source"}}}
}

func TestMutateSourceLinkedDateTaskRequestValidation(t *testing.T) {
	v := validLinkedTaskRequest()
	if err := v.Validate(); err != nil {
		t.Fatal(err)
	}
	for _, invalid := range []string{"2026-9-01", "2026-02-30", "2026-09-01T00:00:00Z"} {
		v := validLinkedTaskRequest()
		v.DueDate = invalid
		if err := v.Validate(); err == nil {
			t.Fatalf("accepted invalid date %q", invalid)
		}
	}
	v = validLinkedTaskRequest()
	v.ActionDisposition = "execute"
	if err := v.Validate(); err == nil {
		t.Fatal("accepted unimplemented execute action")
	}
	v = validLinkedTaskRequest()
	v.Related = append(v.Related, v.Related[0])
	if err := v.Validate(); err == nil {
		t.Fatal("accepted duplicate Linkage reference")
	}
}

func TestSourceLinkedDateTaskHappeningItemRef(t *testing.T) {
	v := validLinkedTaskRequest()
	task := SourceLinkedDateTask{HappeningID: "due-task", Source: v.Source}
	ref := task.HappeningItemRef()
	if ref.ExtID != CalendariusExtensionID || ref.Collection != CalendariusHappeningsCollection || ref.ItemID != "due-task@family1" {
		t.Fatalf("ref=%+v", ref)
	}
}

func TestSourceLinkedDateTask_AdditionalValidation(t *testing.T) {
	// SourceLinkedDateTaskRelatedRef: invalid ItemRef
	r1 := SourceLinkedDateTaskRelatedRef{ItemRef: coretypes.ItemRef{}, Role: "source"}
	if err := r1.Validate(); err == nil {
		t.Fatal("expected error for invalid ItemRef")
	}
	// SourceLinkedDateTaskRelatedRef: invalid Role
	r2 := SourceLinkedDateTaskRelatedRef{ItemRef: coretypes.ItemRef{ExtID: "ext", Collection: "col", ItemID: "id"}, Role: ""}
	if err := r2.Validate(); err == nil {
		t.Fatal("expected error for empty Role")
	}

	// SourceLinkedDateTaskRef: invalid SpaceID
	ref := SourceLinkedDateTaskRef{OwnerSpaceID: "invalid/space", Namespace: "debtus", RecordID: "r1", LineID: "l1"}
	if err := ref.Validate(); err == nil {
		t.Fatal("expected error for invalid spaceID")
	}
	// SourceLinkedDateTaskRef: invalid Namespace
	ref = SourceLinkedDateTaskRef{OwnerSpaceID: "space1", Namespace: "", RecordID: "r1", LineID: "l1"}
	if err := ref.Validate(); err == nil {
		t.Fatal("expected error for empty namespace")
	}
	// SourceLinkedDateTaskRef: invalid RecordID (contains tab)
	ref = SourceLinkedDateTaskRef{OwnerSpaceID: "space1", Namespace: "debtus", RecordID: "r\t1", LineID: "l1"}
	if err := ref.Validate(); err == nil {
		t.Fatal("expected error for invalid recordID")
	}

	// SourceLinkedDateTask.Validate
	validRef := SourceLinkedDateTaskRef{OwnerSpaceID: "space1", Namespace: "debtus", RecordID: "r1", LineID: "l1"}
	// invalid happeningID or revision
	t1 := SourceLinkedDateTask{HappeningID: "", Revision: 1, Source: validRef, Title: "task", State: SourceLinkedDateTaskCompleted}
	if err := t1.Validate(); err == nil {
		t.Fatal("expected error for empty happeningID")
	}
	// invalid title
	t2 := SourceLinkedDateTask{HappeningID: "h1", Revision: 1, Source: validRef, Title: "", State: SourceLinkedDateTaskCompleted}
	if err := t2.Validate(); err == nil {
		t.Fatal("expected error for empty title")
	}
	// unsupported state
	t3 := SourceLinkedDateTask{HappeningID: "h1", Revision: 1, Source: validRef, Title: "task", State: "unknown"}
	if err := t3.Validate(); err == nil {
		t.Fatal("expected error for unknown state")
	}
	// active state with empty dueDate
	t4 := SourceLinkedDateTask{HappeningID: "h1", Revision: 1, Source: validRef, Title: "task", State: SourceLinkedDateTaskActive, DueDate: ""}
	if err := t4.Validate(); err == nil {
		t.Fatal("expected error for active task with empty dueDate")
	}
	// actionID without actionDisposition
	t5 := SourceLinkedDateTask{HappeningID: "h1", Revision: 1, Source: validRef, Title: "task", State: SourceLinkedDateTaskCompleted, ActionID: "act"}
	if err := t5.Validate(); err == nil {
		t.Fatal("expected error for actionID without actionDisposition")
	}
	// actionDisposition without actionID
	t6 := SourceLinkedDateTask{HappeningID: "h1", Revision: 1, Source: validRef, Title: "task", State: SourceLinkedDateTaskCompleted, ActionDisposition: SourceLinkedDateTaskNavigate}
	if err := t6.Validate(); err == nil {
		t.Fatal("expected error for actionDisposition without actionID")
	}
	// actionID untrimmed
	t7 := SourceLinkedDateTask{HappeningID: "h1", Revision: 1, Source: validRef, Title: "task", State: SourceLinkedDateTaskCompleted, ActionID: " act ", ActionDisposition: SourceLinkedDateTaskNavigate}
	if err := t7.Validate(); err == nil {
		t.Fatal("expected error for untrimmed actionID")
	}
	// actionDisposition RequiresInput
	t8 := SourceLinkedDateTask{HappeningID: "h1", Revision: 1, Source: validRef, Title: "task", State: SourceLinkedDateTaskCompleted, ActionID: "act", ActionDisposition: SourceLinkedDateTaskRequiresInput}
	if err := t8.Validate(); err != nil {
		t.Fatalf("unexpected error for RequiresInput: %v", err)
	}
	// empty actionDisposition and empty actionID -> return nil
	t9 := SourceLinkedDateTask{HappeningID: "h1", Revision: 1, Source: validRef, Title: "task", State: SourceLinkedDateTaskCompleted}
	if err := t9.Validate(); err != nil {
		t.Fatalf("unexpected error for empty action: %v", err)
	}

	// MutateSourceLinkedDateTaskRequest
	req := validLinkedTaskRequest()
	req.OperationID = ""
	if err := req.Validate(); err == nil {
		t.Fatal("expected error for empty operationID")
	}
	req = validLinkedTaskRequest()
	req.ExpectedRevision = -1
	if err := req.Validate(); err == nil {
		t.Fatal("expected error for negative expectedRevision")
	}
	req = validLinkedTaskRequest()
	req.Related = nil
	if err := req.Validate(); err == nil {
		t.Fatal("expected error for nil related")
	}
	req = validLinkedTaskRequest()
	req.Related = []SourceLinkedDateTaskRelatedRef{}
	if err := req.Validate(); err == nil {
		t.Fatal("expected error for empty related")
	}
	req = validLinkedTaskRequest()
	req.Related = []SourceLinkedDateTaskRelatedRef{{ItemRef: coretypes.ItemRef{}, Role: "source"}}
	if err := req.Validate(); err == nil {
		t.Fatal("expected error for invalid related itemRef")
	}
}
