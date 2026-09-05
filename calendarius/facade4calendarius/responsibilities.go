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
	CreateScheduledResponsibility(ctx context.Context, actorUserID, spaceID string, request calendariusmodels.CreateScheduledResponsibilityRequest) (calendariusmodels.ScheduledResponsibilityMutation, error)
	ListResponsibilityOccurrences(ctx context.Context, actorUserID, spaceID, happeningID string, from, to time.Time) ([]calendariusmodels.ResponsibilityOccurrence, error)
	CompleteResponsibilityOccurrence(ctx context.Context, actorUserID, spaceID string, request calendariusmodels.CompleteResponsibilityOccurrenceRequest) (calendariusmodels.ResponsibilityCompletionMutation, error)
}
