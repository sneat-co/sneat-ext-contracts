package contract4debtus

import (
	"strings"
	"testing"
	"time"
)

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

func TestSourceObligationDueDateResultTitleIsOptionalAndBounded(t *testing.T) {
	result := SourceObligationDueDateResult{
		ContractVersion: SourceDueDateContractVersion,
		Source:          SourceRef{Namespace: "splitus", SpaceID: "family1", RecordID: "bill1"},
		LineID:          "guest-share", Revision: 1, DueDate: "2026-09-30", HappeningID: "task1",
		State: SourceObligationDueDateActive, UpdatedAt: time.Now().UTC(), UpdatedBy: "owner",
	}
	if err := result.Validate(); err != nil {
		t.Fatalf("legacy result without title: %v", err)
	}
	result.Title = "Pay electricity bill"
	if err := result.Validate(); err != nil {
		t.Fatalf("result with title: %v", err)
	}
	result.Title = strings.Repeat("x", 101)
	if err := result.Validate(); err == nil {
		t.Fatal("accepted oversized title")
	}
}
