// Package contract4contactus defines persistence-free Contactus contracts for
// conversational and other interactive clients.
package contract4contactus

import "context"

// Contact is the compact, presentation-safe view of a contact used by a bot.
// It intentionally contains no DAL, DBO, facade, or framework types.
type Contact struct {
	ID   string
	Name string
}

// CreateContactRequest contains the fields a conversational client may use to
// create a person contact.
type CreateContactRequest struct {
	SpaceID string
	Name    string
	Email   string
	Phone   string
}

// ConvoService is the Contactus application port for conversational clients.
// Implementations own authorization and persistence; callers never do.
type ConvoService interface {
	CreateContact(ctx context.Context, request CreateContactRequest) (Contact, error)
	ListContacts(ctx context.Context, spaceID string) ([]Contact, error)
	DeleteContact(ctx context.Context, spaceID, contactID string) error
}
