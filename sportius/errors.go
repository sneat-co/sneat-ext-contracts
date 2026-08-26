package sportius

import "fmt"

// ErrorCode is stable across Go, TypeScript and the HTTP wire contract.
type ErrorCode string

const (
	ErrorCodeValidation        ErrorCode = "validation"
	ErrorCodeForbidden         ErrorCode = "forbidden"
	ErrorCodeNotFound          ErrorCode = "not-found"
	ErrorCodeMerged            ErrorCode = "merged"
	ErrorCodeConflict          ErrorCode = "conflict"
	ErrorCodeInvitationExpired ErrorCode = "invitation-expired"
	ErrorCodeInviteRequired    ErrorCode = "invite-required"
	ErrorCodeRetryable         ErrorCode = "retryable"
	ErrorCodeInternal          ErrorCode = "internal"
)

// Error is safe application-level metadata. MessageKey is localised by the
// presentation surface; Cause is never serialised or returned to users.
type Error struct {
	Code       ErrorCode `json:"code"`
	MessageKey string    `json:"messageKey"`
	Field      string    `json:"field,omitempty"`
	Retryable  bool      `json:"retryable"`
	Cause      error     `json:"-"`
}

func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Field == "" {
		return fmt.Sprintf("sportius: %s", e.Code)
	}
	return fmt.Sprintf("sportius: %s (%s)", e.Code, e.Field)
}

func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}
