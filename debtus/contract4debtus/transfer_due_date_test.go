package contract4debtus

import "testing"

func TestSetTransferDueDateRequestRequiresExplicitSpaceAndActor(t *testing.T) {
	valid := SetTransferDueDateRequest{ContractVersion: 1, SpaceID: "house", TransferID: "transfer1", Title: "Repay Alex", DueDate: "2026-10-01", OperationKey: "operation1", ActorUserID: "member1"}
	if err := valid.Validate(); err != nil {
		t.Fatalf("valid request: %v", err)
	}
	for name, mutate := range map[string]func(*SetTransferDueDateRequest){
		"space": func(r *SetTransferDueDateRequest) { r.SpaceID = "" },
		"actor": func(r *SetTransferDueDateRequest) { r.ActorUserID = "" },
		"date":  func(r *SetTransferDueDateRequest) { r.DueDate = "01/10/2026" },
	} {
		t.Run(name, func(t *testing.T) {
			r := valid
			mutate(&r)
			if err := r.Validate(); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}
