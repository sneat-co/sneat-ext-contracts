package contract4debtus

import "testing"

func TestSetSourceObligationDueDateRequestValidation(t *testing.T) {
	r := SetSourceObligationDueDateRequest{ContractVersion: 1, Source: SourceRef{Namespace: "splitus", SpaceID: "family1", RecordID: "bill1"}, LineID: "guest-share", Title: "Pay electricity bill", DueDate: "2026-09-30", OperationKey: "due1", ActorUserID: "owner"}
	if err := r.Validate(); err != nil {
		t.Fatal(err)
	}
	r.DueDate = "2026-9-30"
	if err := r.Validate(); err == nil {
		t.Fatal("accepted non-canonical due date")
	}
	r.DueDate = "2026-09-30"
	r.ActorUserID = ""
	if err := r.Validate(); err == nil {
		t.Fatal("accepted missing authenticated actor")
	}
}
