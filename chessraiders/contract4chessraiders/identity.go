package contract4chessraiders

import (
	"context"
	"errors"
)

// IdentityProvider is an external platform that can assert a player's
// identity before Chess Raiders lets it act inside a match. It is distinct
// from any Sneat-internal auth provider.
type IdentityProvider string

const (
	IdentityProviderCrazyGames IdentityProvider = "crazygames"
	// IdentityProviderItchIO is deliberately "itch-io", not "itch.io" — it
	// matches the stable registry key already live in the private
	// implementation (sneat-co/chessraiders' api4chess.ExternalIdentity.Provider
	// and sneat-go's chessExternalIdentityProviders allow-list), not the
	// provider's own dotted brand name.
	IdentityProviderItchIO IdentityProvider = "itch-io"
)

// ExternalIdentity is what a portal or external provider vouches for after
// its own credential has been verified. It carries no Sneat-internal user
// ID: resolving or minting one from this value is the host's job, never the
// provider's.
type ExternalIdentity struct {
	Provider       IdentityProvider  `json:"provider"`
	ExternalUserID string            `json:"externalUserID"`
	DisplayName    string            `json:"displayName,omitempty"`
	Username       string            `json:"username,omitempty"`
	AvatarURL      string            `json:"avatarURL,omitempty"`
	Metadata       map[string]string `json:"metadata,omitempty"`
}

// ExternalIdentityVerifier turns one provider's opaque credential into an
// ExternalIdentity. It never resolves or mints a Sneat user itself — that
// composition belongs to IdentityLinkApplication, which consumes a verifier.
type ExternalIdentityVerifier interface {
	Provider() IdentityProvider
	Verify(ctx context.Context, credential string) (ExternalIdentity, error)
}

// IdentityLinkApplication binds an ExternalIdentity to a Sneat user, and
// resolves an existing binding without ever creating a second one for the
// same (provider, externalUserID) pair.
type IdentityLinkApplication interface {
	// LinkExternalIdentity binds identity to sneatUserID. It is idempotent for
	// a repeat of the same pair, and rejects reassigning an already-linked
	// externalUserID to a different sneatUserID with ErrIdentityAlreadyLinked.
	LinkExternalIdentity(ctx context.Context, sneatUserID string, identity ExternalIdentity) error
	// ResolveExternalIdentity looks up the Sneat user already linked to
	// provider/externalUserID, if any. found is false, with a nil error, when
	// no binding exists yet — that is not itself an error condition.
	ResolveExternalIdentity(
		ctx context.Context,
		provider IdentityProvider,
		externalUserID string,
	) (sneatUserID string, found bool, err error)
}

var (
	// ErrIdentityProviderDisabled is returned when a provider is recognised
	// but not enabled for this host (e.g. itch.io identity is opt-in).
	ErrIdentityProviderDisabled = errors.New("chess raiders contract: identity provider is not enabled")
	// ErrIdentityAlreadyLinked is returned when the externalUserID in a
	// LinkExternalIdentity call is already bound to a different Sneat user.
	// It is never silently merged or overwritten.
	ErrIdentityAlreadyLinked = errors.New("chess raiders contract: external identity is already linked to a different user")
)
