// Package facade4bookius defines Bookius capabilities exposed to other extensions.
package facade4bookius

import "context"

// VendorBot is the minimal public read model needed to deep-link a customer to
// a vendor's payment bot. Provider tokens and implementation storage stay
// private to Bookius.
type VendorBot struct {
	Code string
}

// VendorBotResolver resolves the enabled vendor bot serving a business space.
// A missing or disabled bot is represented by ok=false, not a provider-private
// sentinel error.
type VendorBotResolver interface {
	ResolveVendorBotBySpace(ctx context.Context, vendorSpaceID string) (bot VendorBot, ok bool, err error)
}
