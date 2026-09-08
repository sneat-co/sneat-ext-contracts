package facade4listus

import (
	"context"

	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-ext-contracts/listus/listusmodels"
)

// PreparedSourceTodo is an opaque, transaction-bound result of all reads and
// authorization checks. Implementations expose only its immutable view.
type PreparedSourceTodo interface {
	SourceTodoPlanView() listusmodels.SourceTodoPlanView
}

// SourceTodoPort lets another source owner coordinate Listus writes inside its
// ambient DALgo transaction. PlanSourceTodo performs every read; ApplySourceTodo
// accepts only a plan created by that implementation and performs no reads.
type SourceTodoPort interface {
	PlanSourceTodo(ctx context.Context, tx dal.ReadwriteTransaction, actorUserID string, spec listusmodels.SourceTodoSpec) (PreparedSourceTodo, error)
	ApplySourceTodo(ctx context.Context, tx dal.ReadwriteTransaction, prepared PreparedSourceTodo) (listusmodels.SourceTodoPlanView, error)
}
