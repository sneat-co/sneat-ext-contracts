package facade4calendarius

import (
	"context"

	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-ext-contracts/calendarius/calendariusmodels"
)

// PreparedSourceLinkedDateTaskMutation is an in-memory write plan bound to the
// transaction used to prepare it. It deliberately exposes no Calendar DBOs.
// Every participating extension prepares first; only then may callers apply
// all plans and queue their own writes.
type PreparedSourceLinkedDateTaskMutation interface {
	Result() calendariusmodels.SourceLinkedDateTaskMutation
	Apply(ctx context.Context, tx dal.ReadwriteTransaction) error
}

// SourceLinkedDateTaskProvider prepares a native Calendarius task mutation
// inside the caller's transaction without writing or committing it.
type SourceLinkedDateTaskProvider interface {
	PlanSourceLinkedDateTask(ctx context.Context, tx dal.ReadwriteTransaction, actorUserID string, request calendariusmodels.MutateSourceLinkedDateTaskRequest) (PreparedSourceLinkedDateTaskMutation, error)
}
