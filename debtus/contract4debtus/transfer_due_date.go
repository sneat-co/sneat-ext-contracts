package contract4debtus

import (
	"context"
	"fmt"
	"strings"
	"time"
)

const TransferDueDateContractVersion = 1

type SetTransferDueDateRequest struct {
	ContractVersion  int    `json:"contractVersion"`
	SpaceID          string `json:"spaceID"`
	TransferID       string `json:"transferID"`
	Title            string `json:"title"`
	DueDate          string `json:"dueDate,omitempty"`
	ExpectedRevision uint64 `json:"expectedRevision"`
	OperationKey     string `json:"operationKey"`
	TodoListID       string `json:"todoListID,omitempty"`
	ActorUserID      string `json:"-"`
}

func (r SetTransferDueDateRequest) Validate() error {
	if r.ContractVersion != TransferDueDateContractVersion {
		return fmt.Errorf("%w: unsupported transfer due-date contract version", ErrInvalidRequest)
	}
	for name, value := range map[string]string{"spaceID": r.SpaceID, "transferID": r.TransferID, "operation key": r.OperationKey} {
		if err := validateID(name, value); err != nil {
			return err
		}
	}
	if r.Title == "" || r.Title != strings.TrimSpace(r.Title) || len(r.Title) > 100 {
		return fmt.Errorf("%w: title must be trimmed and 1..100 bytes", ErrInvalidRequest)
	}
	if r.DueDate != "" {
		parsed, err := time.Parse("2006-01-02", r.DueDate)
		if err != nil || parsed.Format("2006-01-02") != r.DueDate {
			return fmt.Errorf("%w: dueDate must be YYYY-MM-DD", ErrInvalidRequest)
		}
	}
	if r.TodoListID != "" {
		if err := validateID("todoListID", r.TodoListID); err != nil {
			return err
		}
	}
	if r.ExpectedRevision > 9_007_199_254_740_991 {
		return fmt.Errorf("%w: expectedRevision exceeds browser safe integer", ErrInvalidRequest)
	}
	if r.ActorUserID == "" {
		return fmt.Errorf("%w: authenticated actor is required", ErrInvalidRequest)
	}
	return nil
}

type TransferDueDateResult struct {
	ContractVersion int                          `json:"contractVersion"`
	SpaceID         string                       `json:"spaceID"`
	TransferID      string                       `json:"transferID"`
	Revision        uint64                       `json:"revision"`
	DueDate         string                       `json:"dueDate,omitempty"`
	HappeningID     string                       `json:"happeningID"`
	State           SourceObligationDueDateState `json:"state"`
	TodoListID      string                       `json:"todoListID,omitempty"`
	TodoItemID      string                       `json:"todoItemID,omitempty"`
	UpdatedAt       time.Time                    `json:"updatedAt"`
	UpdatedBy       string                       `json:"updatedBy"`
}

func (r TransferDueDateResult) Validate() error {
	if r.ContractVersion != TransferDueDateContractVersion || r.Revision == 0 || r.Revision > 9_007_199_254_740_991 {
		return fmt.Errorf("%w: invalid transfer due-date result version or revision", ErrInvalidRequest)
	}
	for name, value := range map[string]string{"spaceID": r.SpaceID, "transferID": r.TransferID, "happeningID": r.HappeningID, "updatedBy": r.UpdatedBy} {
		if err := validateID(name, value); err != nil {
			return err
		}
	}
	if r.UpdatedAt.IsZero() {
		return fmt.Errorf("%w: updatedAt is required", ErrInvalidRequest)
	}
	if r.DueDate != "" {
		parsed, err := time.Parse("2006-01-02", r.DueDate)
		if err != nil || parsed.Format("2006-01-02") != r.DueDate {
			return fmt.Errorf("%w: dueDate must be YYYY-MM-DD", ErrInvalidRequest)
		}
	}
	switch r.State {
	case SourceObligationDueDateActive:
		if r.DueDate == "" {
			return fmt.Errorf("%w: active due-date task requires dueDate", ErrInvalidRequest)
		}
	case SourceObligationDueDateCompleted, SourceObligationDueDateCanceled:
	default:
		return fmt.Errorf("%w: unsupported due-date task state", ErrInvalidRequest)
	}
	return nil
}

type TransferDueDates interface {
	SetTransferDueDate(context.Context, SetTransferDueDateRequest) (TransferDueDateResult, error)
	GetTransferDueDate(context.Context, GetTransferDueDateRequest) (TransferDueDateResult, error)
	RecordTransferRepayment(context.Context, RecordTransferRepaymentRequest) (RecordTransferRepaymentResult, error)
}

type GetTransferDueDateRequest struct {
	SpaceID     string `json:"spaceID"`
	TransferID  string `json:"transferID"`
	ActorUserID string `json:"-"`
}

func (r GetTransferDueDateRequest) Validate() error {
	for name, value := range map[string]string{"spaceID": r.SpaceID, "transferID": r.TransferID, "actor": r.ActorUserID} {
		if err := validateID(name, value); err != nil {
			return err
		}
	}
	return nil
}

type RecordTransferRepaymentRequest struct {
	ContractVersion int                    `json:"contractVersion"`
	SpaceID         string                 `json:"spaceID"`
	TransferID      string                 `json:"transferID"`
	Currency        string                 `json:"currency"`
	AmountMinor     ExactMinorAmountString `json:"amountMinor"`
	RepaidAt        time.Time              `json:"repaidAt"`
	OperationKey    string                 `json:"operationKey"`
	ActorUserID     string                 `json:"-"`
}

func (r RecordTransferRepaymentRequest) Validate() error {
	if r.ContractVersion != TransferDueDateContractVersion {
		return fmt.Errorf("%w: unsupported transfer repayment contract version", ErrInvalidRequest)
	}
	for name, value := range map[string]string{"spaceID": r.SpaceID, "transferID": r.TransferID, "operationKey": r.OperationKey, "actor": r.ActorUserID} {
		if err := validateID(name, value); err != nil {
			return err
		}
	}
	if len(r.Currency) != 3 || strings.ToUpper(r.Currency) != r.Currency {
		return fmt.Errorf("%w: currency must be 3 uppercase letters", ErrInvalidRequest)
	}
	minor, err := r.AmountMinor.MinorUnits()
	if err != nil || minor <= 0 {
		return fmt.Errorf("%w: amountMinor must be positive: %v", ErrInvalidRequest, err)
	}
	if r.RepaidAt.IsZero() || r.RepaidAt.Location() != time.UTC || r.RepaidAt.Nanosecond()%int(time.Millisecond) != 0 {
		return fmt.Errorf("%w: repaidAt must be UTC with millisecond precision", ErrInvalidRequest)
	}
	return nil
}

type RecordTransferRepaymentResult struct {
	ContractVersion  int                    `json:"contractVersion"`
	SpaceID          string                 `json:"spaceID"`
	TransferID       string                 `json:"transferID"`
	RepaymentID      string                 `json:"repaymentID"`
	OutstandingMinor ExactMinorAmountString `json:"outstandingMinor"`
	FullyRepaid      bool                   `json:"fullyRepaid"`
	DueDateTask      TransferDueDateResult  `json:"dueDateTask"`
}
