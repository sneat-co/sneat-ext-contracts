package calendariusmodels

import (
	"strings"
	"testing"

	"github.com/sneat-co/sneat-go-core/coretypes"
)

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
