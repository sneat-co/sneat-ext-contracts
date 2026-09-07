package contract4debtus

import (
	"testing"
	"time"
)

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

func TestRecordTransferRepaymentRequestUsesExactMinorUnitsAndExplicitSpace(t *testing.T) {
	valid := RecordTransferRepaymentRequest{ContractVersion: 1, SpaceID: "house", TransferID: "transfer1", Currency: "EUR", AmountMinor: "123", RepaidAt: time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC), OperationKey: "repay1", ActorUserID: "member1"}
	if err := valid.Validate(); err != nil {
		t.Fatal(err)
	}
	for name, mutate := range map[string]func(*RecordTransferRepaymentRequest){
		"space":             func(r *RecordTransferRepaymentRequest) { r.SpaceID = "" },
		"numeric ambiguity": func(r *RecordTransferRepaymentRequest) { r.AmountMinor = "01" },
		"lower currency":    func(r *RecordTransferRepaymentRequest) { r.Currency = "eur" },
		"missing actor":     func(r *RecordTransferRepaymentRequest) { r.ActorUserID = "" },
	} {
		request := valid
		mutate(&request)
		if request.Validate() == nil {
			t.Fatalf("%s accepted", name)
		}
	}
}

func TestTransferDueDateResultTitleIsOptionalAndBounded(t *testing.T) {
	result := TransferDueDateResult{
		ContractVersion: TransferDueDateContractVersion, SpaceID: "house", TransferID: "transfer1",
		Revision: 1, DueDate: "2026-10-01", HappeningID: "task1", State: SourceObligationDueDateActive,
		UpdatedAt: time.Now().UTC(), UpdatedBy: "member1", Currency: "EUR", OutstandingMinor: "123",
	}
	if err := result.Validate(); err != nil {
		t.Fatalf("legacy result without title: %v", err)
	}
	result.Title = "Repay Alex"
	if err := result.Validate(); err != nil {
		t.Fatalf("result with title: %v", err)
	}
	result.Title = " untrimmed"
	if err := result.Validate(); err == nil {
		t.Fatal("accepted untrimmed title")
	}
}

func TestRecordTransferRepaymentResultRequiresMatchingDueTaskIdentity(t *testing.T) {
	due := TransferDueDateResult{
		ContractVersion: TransferDueDateContractVersion, SpaceID: "house", TransferID: "transfer1",
		Revision: 2, HappeningID: "task1", State: SourceObligationDueDateCompleted,
		UpdatedAt: time.Now().UTC(), UpdatedBy: "member1", Currency: "EUR", OutstandingMinor: "0",
	}
	valid := RecordTransferRepaymentResult{
		ContractVersion: TransferDueDateContractVersion, SpaceID: "house", TransferID: "transfer1",
		RepaymentID: "return1", OutstandingMinor: "0", FullyRepaid: true, DueDateTask: due,
	}
	if err := valid.Validate(); err != nil {
		t.Fatalf("valid receipt: %v", err)
	}
	for name, mutate := range map[string]func(*RecordTransferRepaymentResult){
		"other transfer in task": func(r *RecordTransferRepaymentResult) { r.DueDateTask.TransferID = "transfer2" },
		"other space in task":    func(r *RecordTransferRepaymentResult) { r.DueDateTask.SpaceID = "flat" },
		"active task without due": func(r *RecordTransferRepaymentResult) {
			r.DueDateTask.State = SourceObligationDueDateActive
		},
		"missing repayment ID": func(r *RecordTransferRepaymentResult) { r.RepaymentID = "" },
		"non-canonical amount": func(r *RecordTransferRepaymentResult) { r.OutstandingMinor = "01" },
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
