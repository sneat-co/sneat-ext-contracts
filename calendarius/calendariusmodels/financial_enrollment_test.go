package calendariusmodels

import "testing"

func TestResolveFinancialEnrollmentRequestCanonicalScope(t *testing.T) {
	request := ResolveFinancialEnrollmentRequest{OwnerSpaceID: "space1", HappeningID: "happening1", OperationID: "operation1", Scope: FinancialEnrollmentScope{Kind: FinancialEnrollmentScopeSubjects, Subjects: []FinancialEnrollmentSubjectRef{{ExtensionID: FinancialEnrollmentSubjectContactus, EntityID: "child-b"}, {ExtensionID: FinancialEnrollmentSubjectContactus, EntityID: "child-a"}}}}
	if err := request.Validate(); err != nil {
		t.Fatal(err)
	}
	canonical := request.Scope.Canonical()
	if canonical.Subjects[0].EntityID != "child-a" || canonical.Subjects[1].EntityID != "child-b" {
		t.Fatalf("canonical=%+v", canonical)
	}
	if err := (FinancialEnrollmentScope{Kind: FinancialEnrollmentScopeWholeHappening, Subjects: []FinancialEnrollmentSubjectRef{}}).Validate(); err != nil {
		t.Fatal(err)
	}
	for _, scope := range []FinancialEnrollmentScope{{Kind: FinancialEnrollmentScopeSubjects, Subjects: []FinancialEnrollmentSubjectRef{}}, {Kind: FinancialEnrollmentScopeWholeHappening, Subjects: []FinancialEnrollmentSubjectRef{{ExtensionID: FinancialEnrollmentSubjectContactus, EntityID: "child"}}}, {Kind: FinancialEnrollmentScopeSubjects, Subjects: []FinancialEnrollmentSubjectRef{{ExtensionID: FinancialEnrollmentSubjectContactus, EntityID: "child"}, {ExtensionID: FinancialEnrollmentSubjectContactus, EntityID: "child"}}}} {
		if scope.Validate() == nil {
			t.Fatalf("accepted scope %+v", scope)
		}
	}
}
