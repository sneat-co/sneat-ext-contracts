// Package contract4chessraiderstest is the conformance harness an
// implementor of contract4chessraiders.MatchLobbyApplication runs against
// its own adapter, mirroring ext-competios's
// backend/contract4competiostest package for the same purpose.
//
// Scope note: only MatchLobbyApplication has a conformance checker here.
// IdentityLinkApplication's and PortalSessionApplication's invariants
// (idempotent link, no silent rebind) are narrow enough that an
// implementation's own unit tests can assert them directly; a shared
// checker for those two was left out to keep this first release scoped to
// the interface with real lifecycle sequencing to get wrong. See the
// ext-chessraiders README for the follow-up note.
package contract4chessraiderstest

import (
	"context"
	"errors"
	"fmt"

	"github.com/sneat-co/sneat-ext-contracts/chessraiders/contract4chessraiders"
)

// MatchLobbyApplicationFactory returns a fresh, empty implementation for one
// check. CheckMatchLobbyApplication calls it more than once and expects each
// call to start from a clean slate unless stated otherwise.
type MatchLobbyApplicationFactory func() contract4chessraiders.MatchLobbyApplication

// CheckMatchLobbyApplication runs the lobby lifecycle an implementation must
// support: create, join both sides, ready up, and start. It returns every
// violation found rather than stopping at the first, so a failing
// implementation gets one report instead of one failure per run.
func CheckMatchLobbyApplication(factory MatchLobbyApplicationFactory) []error {
	ctx := context.Background()
	app := factory()
	var violations []error

	creator := contract4chessraiders.Player{UserID: "creator", DisplayName: "Creator"}
	lobby, err := app.CreateLobby(ctx, creator, contract4chessraiders.CaptureModeKillOnly)
	if err != nil {
		return []error{fmt.Errorf("create lobby: %w", err)}
	}
	if lobby.MatchID == "" {
		return []error{errors.New("create lobby: MatchID is empty")}
	}
	if lobby.Lifecycle != contract4chessraiders.LifecycleLobby {
		violations = append(violations, fmt.Errorf("create lobby: Lifecycle = %v, want LifecycleLobby", lobby.Lifecycle))
	}
	if lobby.CreatedBy != creator.UserID {
		violations = append(violations, fmt.Errorf("create lobby: CreatedBy = %q, want %q", lobby.CreatedBy, creator.UserID))
	}

	// Starting before the opposing side has joined must be rejected, not
	// silently accepted as a solo match — a MatchLobbyApplication that starts
	// with only one side present has not implemented ErrTeamsNotViable.
	if _, err := app.Start(ctx, lobby.MatchID, creator.UserID); err == nil {
		violations = append(violations, errors.New("start with one side seated succeeded, want ErrTeamsNotViable or similar"))
	}

	opponent := contract4chessraiders.Player{UserID: "opponent", DisplayName: "Opponent", Side: contract4chessraiders.SideBlack}
	joined, err := app.Join(ctx, lobby.MatchID, opponent)
	if err != nil {
		violations = append(violations, fmt.Errorf("join: %w", err))
	} else if len(joined.Black) != 1 || joined.Black[0].UserID != opponent.UserID {
		violations = append(violations, fmt.Errorf("join: Black = %+v, want one seat for %q", joined.Black, opponent.UserID))
	}

	// Starting before every seated player has voted ready must be rejected.
	if _, err := app.Start(ctx, lobby.MatchID, creator.UserID); err == nil {
		violations = append(violations, errors.New("start before ready votes succeeded, want ErrNotAllReady or similar"))
	}

	if _, err := app.VoteReady(ctx, lobby.MatchID, creator.UserID, true); err != nil {
		violations = append(violations, fmt.Errorf("creator vote ready: %w", err))
	}
	if _, err := app.VoteReady(ctx, lobby.MatchID, opponent.UserID, true); err != nil {
		violations = append(violations, fmt.Errorf("opponent vote ready: %w", err))
	}

	// Start is deliberately NOT creator-gated: any joined, ready player may
	// trigger it once the teams are viable and everyone is ready — so the
	// non-creator opponent starting here must succeed, not fail.
	started, err := app.Start(ctx, lobby.MatchID, opponent.UserID)
	if err != nil {
		violations = append(violations, fmt.Errorf("start by a non-creator, ready player: %w", err))
	} else if started.Lifecycle == contract4chessraiders.LifecycleLobby {
		violations = append(violations, fmt.Errorf("start: Lifecycle stayed %v", started.Lifecycle))
	}

	// A match no longer in the lobby stage must reject a second join.
	stranger := contract4chessraiders.Player{UserID: "stranger", DisplayName: "Stranger", Side: contract4chessraiders.SideWhite}
	if _, err := app.Join(ctx, lobby.MatchID, stranger); !errors.Is(err, contract4chessraiders.ErrNotInLobby) {
		violations = append(violations, fmt.Errorf("join after start error = %v, want ErrNotInLobby", err))
	}

	if _, err := app.Get(ctx, "unknown-match-id", creator.UserID); !errors.Is(err, contract4chessraiders.ErrMatchNotFound) {
		violations = append(violations, fmt.Errorf("get unknown match error = %v, want ErrMatchNotFound", err))
	}

	return violations
}
