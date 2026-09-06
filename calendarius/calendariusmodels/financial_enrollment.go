package calendariusmodels

import (
	"fmt"
	"sort"
)

const (
	FinancialEnrollmentScopeWholeHappening = "whole_happening"
	FinancialEnrollmentScopeSubjects       = "subjects"
	FinancialEnrollmentSubjectContactus    = "contactus"
	FinancialEnrollmentSubjectAssetus      = "assetus"
	MaxFinancialEnrollmentSubjects         = 32
)

type FinancialEnrollmentSubjectRef struct {
	ExtensionID string `json:"extensionID"`
	EntityID    string `json:"entityID"`
}

func (r FinancialEnrollmentSubjectRef) Validate() error {
	if r.ExtensionID != FinancialEnrollmentSubjectContactus && r.ExtensionID != FinancialEnrollmentSubjectAssetus {
		return fmt.Errorf("financial enrollment subject extension is unsupported")
	}
	return validateFinancialCommitmentID("subject.entityID", r.EntityID)
}

// FinancialEnrollmentScope is stable domain identity, separate from mutable
// financial attribution and contractual payer/receiver roles.
type FinancialEnrollmentScope struct {
	Kind     string                          `json:"kind"`
	Subjects []FinancialEnrollmentSubjectRef `json:"subjects"`
}

func (s FinancialEnrollmentScope) Validate() error {
	if s.Subjects == nil || len(s.Subjects) > MaxFinancialEnrollmentSubjects {
		return fmt.Errorf("financial enrollment subjects must be a bounded array")
	}
	switch s.Kind {
	case FinancialEnrollmentScopeWholeHappening:
		if len(s.Subjects) != 0 {
			return fmt.Errorf("whole-happening enrollment cannot name subjects")
		}
	case FinancialEnrollmentScopeSubjects:
		if len(s.Subjects) == 0 {
			return fmt.Errorf("subject enrollment requires subjects")
		}
	default:
		return fmt.Errorf("financial enrollment scope kind is unsupported")
	}
	previous := ""
	for i, subject := range s.Subjects {
		if err := subject.Validate(); err != nil {
			return fmt.Errorf("subjects[%d]: %w", i, err)
		}
		key := subject.ExtensionID + "\x00" + subject.EntityID
		if key <= previous {
			return fmt.Errorf("financial enrollment subjects must be unique and canonically sorted")
		}
		previous = key
	}
	return nil
}

func (s FinancialEnrollmentScope) Canonical() FinancialEnrollmentScope {
	result := FinancialEnrollmentScope{Kind: s.Kind, Subjects: append([]FinancialEnrollmentSubjectRef(nil), s.Subjects...)}
	sort.Slice(result.Subjects, func(i, j int) bool {
		if result.Subjects[i].ExtensionID == result.Subjects[j].ExtensionID {
			return result.Subjects[i].EntityID < result.Subjects[j].EntityID
		}
		return result.Subjects[i].ExtensionID < result.Subjects[j].ExtensionID
	})
	return result
}

type ResolveFinancialEnrollmentRequest struct {
	OwnerSpaceID string                   `json:"ownerSpaceID"`
	HappeningID  string                   `json:"happeningID"`
	OperationID  string                   `json:"operationID"`
	Scope        FinancialEnrollmentScope `json:"scope"`
}

func (r ResolveFinancialEnrollmentRequest) Validate() error {
	for name, value := range map[string]string{"ownerSpaceID": r.OwnerSpaceID, "happeningID": r.HappeningID, "operationID": r.OperationID} {
		if err := validateFinancialCommitmentID(name, value); err != nil {
			return err
		}
	}
	return r.Scope.Canonical().Validate()
}

type ResolveFinancialEnrollmentResponse struct {
	EnrollmentID   string                   `json:"enrollmentID"`
	Scope          FinancialEnrollmentScope `json:"scope"`
	IdentityStatus string                   `json:"identityStatus"`
	CreatedAt      string                   `json:"createdAt"`
}
