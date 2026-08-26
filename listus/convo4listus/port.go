// Package convo4listus declares the Listus extension's conversational
// capability: the typed actions it understands, the vocabulary it is listening
// for, the deterministic rules that interpret it, and the port a host calls to
// execute them.
//
// Everything here is persistence-free. The port is satisfied in the Listus
// implementation repo (listus/backend/botapi) over facade4listus; a host binds
// that implementation to these declarations. Nothing in this package may import
// DALgo, a facade, or a transport.
package convo4listus

import "context"

// Standard list IDs, in "{type}!{sub}" form.
const (
	GroceriesListID = "buy!groceries"
	TasksListID     = "do!tasks"
	MoviesListID    = "watch!movies"
	BooksListID     = "read!books"
)

// ItemInput is an item to add.
//
// Title is the BARE noun ("milk"), never "milk 2 liters": quantity lives in its
// own field so that "bought milk" still matches by title and two differently
// sized entries of the same product collide as the same item. That is the whole
// reason quantity is structured rather than folded into the title.
type ItemInput struct {
	Title string
	// Quantity is optional; zero means unspecified.
	Quantity float64
	// Unit is optional and only meaningful with a quantity. The allowed
	// vocabulary is enforced by the action specification's Enum, not here — the
	// domain stores what it is given.
	Unit string
}

// AddItemsRequest adds items to a list, creating the list on first use.
type AddItemsRequest struct {
	SpaceID string
	ListID  string
	Items   []ItemInput
}

// Item is the presentation-safe view of a list item.
type Item struct {
	ID       string
	Title    string
	Quantity float64
	Unit     string
	IsDone   bool
}

// List is the presentation-safe view of a list and its items.
type List struct {
	ID    string
	Title string
	Items []Item
}

// SetDoneRequest marks items done or not done.
type SetDoneRequest struct {
	SpaceID string
	ListID  string
	ItemIDs []string
	IsDone  bool
}

// RemoveItemsRequest removes items from a list.
type RemoveItemsRequest struct {
	SpaceID string
	ListID  string
	ItemIDs []string
}

// ChangedItems reports what a mutation affected, by title, for the reply text.
type ChangedItems struct {
	Titles []string
}

// Port is the Listus capability a host binds to the declared actions.
type Port interface {
	// AddItems adds items, creating the list if it does not exist yet.
	AddItems(ctx context.Context, userID string, request AddItemsRequest) (ChangedItems, error)

	// GetList returns a list and its items. It returns an error for which
	// IsNotFound reports true when the list does not exist.
	GetList(ctx context.Context, userID, spaceID, listID string) (List, error)

	// SetDone marks the given items done or not done.
	SetDone(ctx context.Context, userID string, request SetDoneRequest) (ChangedItems, error)

	// RemoveItems removes the given items.
	RemoveItems(ctx context.Context, userID string, request RemoveItemsRequest) (ChangedItems, error)

	// IsNotFound reports whether err means the list or item does not exist, so
	// the runtime can clarify rather than surface a raw error.
	IsNotFound(err error) bool
}
