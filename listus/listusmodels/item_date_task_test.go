package listusmodels

import "testing"

func TestSaveListItemDateTaskRequestValidation(t *testing.T) {
	valid := SaveListItemDateTaskRequest{SpaceID: "space1", ListID: "do!tasks", ItemID: "item1", OperationID: "op1", DueDate: "2026-09-30", State: SourceTodoActive}
	if err := valid.Validate(); err != nil {
		t.Fatalf("valid request: %v", err)
	}
	for name, mutate := range map[string]func(*SaveListItemDateTaskRequest){
		"active without date": func(v *SaveListItemDateTaskRequest) { v.DueDate = "" },
		"invalid date":        func(v *SaveListItemDateTaskRequest) { v.DueDate = "2026-02-30" },
		"negative revision":   func(v *SaveListItemDateTaskRequest) { v.ExpectedTaskRevision = -1 },
	} {
		t.Run(name, func(t *testing.T) {
			v := valid
			mutate(&v)
			if err := v.Validate(); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
	canceled := valid
	canceled.State = SourceTodoCanceled
	canceled.DueDate = ""
	if err := canceled.Validate(); err != nil {
		t.Fatalf("valid canceled request: %v", err)
	}
}
