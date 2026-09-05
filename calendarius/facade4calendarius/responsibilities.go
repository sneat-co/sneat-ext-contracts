package facade4calendarius

import (
	"context"
	"errors"
	"github.com/sneat-co/sneat-ext-contracts/calendarius/calendariusmodels"
	"time"
)

var (
	ErrResponsibilityInvalid         = errors.New("invalid scheduled responsibility")
	ErrResponsibilityNotFound        = errors.New("scheduled responsibility not found")
	ErrResponsibilityRequestConflict = errors.New("scheduled responsibility request conflict")
	ErrResponsibilityUnauthorized    = errors.New("scheduled responsibility unauthorized")
)

type ScheduledResponsibilitiesFacade interface {
	CreateScheduledResponsibility(context.Context, string, string, calendariusmodels.CreateScheduledResponsibilityRequest) (calendariusmodels.ScheduledResponsibilityMutation, error)
	ListResponsibilityOccurrences(context.Context, string, string, string, time.Time, time.Time) ([]calendariusmodels.ResponsibilityOccurrence, error)
	CompleteResponsibilityOccurrence(context.Context, string, string, calendariusmodels.CompleteResponsibilityOccurrenceRequest) (calendariusmodels.ResponsibilityCompletionMutation, error)
}
