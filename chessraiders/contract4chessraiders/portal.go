package contract4chessraiders

import (
	"context"
	"errors"
)

// PortalSessionRequest asks Chess Raiders to bind a distribution portal's
// own room/session identifier to a match. CrazyGames' roomId is one such
// identifier — an implementation may reuse it 1:1 as MatchID, but is never
// required to. MatchID is empty when the portal is opening a brand new
// session and expects the implementation to create the match.
type PortalSessionRequest struct {
	Provider IdentityProvider `json:"provider"`
	RoomID   string           `json:"roomID"`
	MatchID  string           `json:"matchID,omitempty"`
}

// PortalSession is the durable binding a portal holds after opening a
// session. RoomID is opaque portal data to every other caller, never
// authority.
type PortalSession struct {
	MatchID  string           `json:"matchID"`
	Provider IdentityProvider `json:"provider"`
	RoomID   string           `json:"roomID"`
}

// PortalSessionApplication is the seam a distribution-platform adapter
// (CrazyGames, itch.io, ...) uses to bind its own room/session concept to a
// Chess Raiders match. It carries no launch, identity-verification, or
// gameplay concern of its own — those belong to IdentityLinkApplication and
// MatchLobbyApplication respectively.
type PortalSessionApplication interface {
	OpenPortalSession(ctx context.Context, request PortalSessionRequest) (PortalSession, error)
	GetPortalSession(ctx context.Context, provider IdentityProvider, roomID string) (PortalSession, error)
}

var (
	// ErrPortalRoomAlreadyBound is returned when a RoomID already resolves to
	// a different MatchID than the one requested. A room is never silently
	// rebound to a second match.
	ErrPortalRoomAlreadyBound = errors.New("chess raiders contract: portal room is already bound to a different match")
	ErrPortalSessionNotFound  = errors.New("chess raiders contract: portal session not found")
)
