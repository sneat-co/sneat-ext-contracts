package contract4chessraiders

import (
	"context"
	"errors"
)

// MatchLobbyApplication is the stable social/bot-shell boundary for
// creating, joining, and starting a match. Every method resolves authority
// for callerUserID independently against the implementation's own trusted
// source; a caller cannot act as another player by supplying a different
// UserID inside Player. Join's side comes from joiner.Side — there is no
// separate side parameter, so a caller cannot express two disagreeing
// answers to "which side" in one call.
//
// Start is deliberately NOT creator-gated: any joined, ready player may
// trigger it once every joined player is ready and the teams are viable.
// SetCaptureMode remains creator-only (ErrNotCreator).
//
// Gameplay itself — moves, board state, legal targets — is out of scope.
// Once a match leaves LifecycleLobby, this interface only reports status;
// the private web/Mini App gameplay surface is the sole place a game is
// played.
type MatchLobbyApplication interface {
	CreateLobby(ctx context.Context, creator Player, mode CaptureMode) (LobbyView, error)
	Join(ctx context.Context, matchID string, joiner Player) (LobbyView, error)
	Leave(ctx context.Context, matchID, callerUserID string) (LobbyView, error)
	SetCaptureMode(ctx context.Context, matchID, callerUserID string, mode CaptureMode) (LobbyView, error)
	VoteReady(ctx context.Context, matchID, callerUserID string, ready bool) (LobbyView, error)
	Start(ctx context.Context, matchID, callerUserID string) (LobbyView, error)
	Get(ctx context.Context, matchID, callerUserID string) (LobbyView, error)
}

var (
	ErrMatchNotFound = errors.New("chess raiders contract: match not found")
	ErrNotInLobby    = errors.New("chess raiders contract: match is not in the lobby stage")
	// ErrNotCreator is returned for a creator-only action — SetCaptureMode
	// today — attempted by a non-creator. It is never returned by Start,
	// which any joined, ready player may trigger.
	ErrNotCreator         = errors.New("chess raiders contract: only the lobby creator may do that")
	ErrPlayerNotFound     = errors.New("chess raiders contract: player not found in this match")
	ErrAlreadyJoined      = errors.New("chess raiders contract: player already joined a side")
	ErrSideFull           = errors.New("chess raiders contract: that side is full")
	ErrInvalidSide        = errors.New("chess raiders contract: side must be white or black")
	ErrInvalidCaptureMode = errors.New("chess raiders contract: invalid capture mode")
	// ErrNotAllReady is returned by Start when at least one seated player has
	// not voted ready.
	ErrNotAllReady = errors.New("chess raiders contract: not every player has voted ready")
	// ErrTeamsNotViable is returned by Start when the seated roster cannot
	// form a match — an empty side with another side seated, and the match is
	// not a deliberate solo/bot match.
	ErrTeamsNotViable = errors.New("chess raiders contract: teams are not viable — need a player on both sides, or a deliberate solo match")
	ErrMatchFinished  = errors.New("chess raiders contract: match has already finished")
)
