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
	whole := ResolveFinancialEnrollmentRequest{OwnerSpaceID: "space1", HappeningID: "happening1", OperationID: "operation2", Scope: FinancialEnrollmentScope{Kind: FinancialEnrollmentScopeWholeHappening, Subjects: []FinancialEnrollmentSubjectRef{}}}
	if err := whole.Validate(); err != nil {
		t.Fatalf("whole-happening request rejected: %v", err)
	}
	if canonicalWhole := whole.Scope.Canonical(); canonicalWhole.Subjects == nil || canonicalWhole.Canonical().Subjects == nil {
		t.Fatal("canonicalization changed an explicit empty subjects array to nil")
	}
	for _, scope := range []FinancialEnrollmentScope{{Kind: FinancialEnrollmentScopeWholeHappening}, {Kind: FinancialEnrollmentScopeSubjects, Subjects: []FinancialEnrollmentSubjectRef{}}, {Kind: FinancialEnrollmentScopeWholeHappening, Subjects: []FinancialEnrollmentSubjectRef{{ExtensionID: FinancialEnrollmentSubjectContactus, EntityID: "child"}}}, {Kind: FinancialEnrollmentScopeSubjects, Subjects: []FinancialEnrollmentSubjectRef{{ExtensionID: FinancialEnrollmentSubjectContactus, EntityID: "child"}, {ExtensionID: FinancialEnrollmentSubjectContactus, EntityID: "child"}}}} {
		if scope.Validate() == nil {
			t.Fatalf("accepted scope %+v", scope)
		}
	}
}
