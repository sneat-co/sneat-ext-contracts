package contract4chessraiderstest_test

import (
	"context"
	"sync"
	"testing"

	"github.com/sneat-co/sneat-ext-contracts/chessraiders/contract4chessraiders"
	"github.com/sneat-co/sneat-ext-contracts/chessraiders/contract4chessraiderstest"
)

// memoryLobbyApplication is a minimal, correct-by-construction reference
// implementation used only to prove the conformance checker itself accepts a
// compliant implementation and rejects a non-compliant one. It is not part
// of the public contract.
type memoryLobbyApplication struct {
	mu      sync.Mutex
	matches map[string]*memoryMatch
	nextID  int
}

type memoryMatch struct {
	lifecycle   contract4chessraiders.Lifecycle
	captureMode contract4chessraiders.CaptureMode
	createdBy   string
	white       []contract4chessraiders.LobbySeat
	black       []contract4chessraiders.LobbySeat
	ready       map[string]bool
}

func newMemoryLobbyApplication() *memoryLobbyApplication {
	return &memoryLobbyApplication{matches: map[string]*memoryMatch{}}
}

func (a *memoryLobbyApplication) CreateLobby(
	_ context.Context,
	creator contract4chessraiders.Player,
	mode contract4chessraiders.CaptureMode,
) (contract4chessraiders.LobbyView, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.nextID++
	matchID := "match-" + itoa(a.nextID)
	match := &memoryMatch{
		lifecycle:   contract4chessraiders.LifecycleLobby,
		captureMode: mode,
		createdBy:   creator.UserID,
		white:       []contract4chessraiders.LobbySeat{{UserID: creator.UserID, DisplayName: creator.DisplayName}},
		ready:       map[string]bool{},
	}
	a.matches[matchID] = match
	return a.viewLocked(matchID, match, creator.UserID), nil
}

func (a *memoryLobbyApplication) Join(
	_ context.Context,
	matchID string,
	joiner contract4chessraiders.Player,
) (contract4chessraiders.LobbyView, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	match, ok := a.matches[matchID]
	if !ok {
		return contract4chessraiders.LobbyView{}, contract4chessraiders.ErrMatchNotFound
	}
	if match.lifecycle != contract4chessraiders.LifecycleLobby {
		return contract4chessraiders.LobbyView{}, contract4chessraiders.ErrNotInLobby
	}
	seat := contract4chessraiders.LobbySeat{UserID: joiner.UserID, DisplayName: joiner.DisplayName}
	switch joiner.Side {
	case contract4chessraiders.SideWhite:
		match.white = append(match.white, seat)
	case contract4chessraiders.SideBlack:
		match.black = append(match.black, seat)
	default:
		return contract4chessraiders.LobbyView{}, contract4chessraiders.ErrInvalidSide
	}
	return a.viewLocked(matchID, match, joiner.UserID), nil
}

func (a *memoryLobbyApplication) Leave(_ context.Context, matchID, callerUserID string) (contract4chessraiders.LobbyView, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	match, ok := a.matches[matchID]
	if !ok {
		return contract4chessraiders.LobbyView{}, contract4chessraiders.ErrMatchNotFound
	}
	return a.viewLocked(matchID, match, callerUserID), nil
}

func (a *memoryLobbyApplication) SetCaptureMode(
	_ context.Context,
	matchID, callerUserID string,
	mode contract4chessraiders.CaptureMode,
) (contract4chessraiders.LobbyView, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	match, ok := a.matches[matchID]
	if !ok {
		return contract4chessraiders.LobbyView{}, contract4chessraiders.ErrMatchNotFound
	}
	if match.createdBy != callerUserID {
		return contract4chessraiders.LobbyView{}, contract4chessraiders.ErrNotCreator
	}
	match.captureMode = mode
	return a.viewLocked(matchID, match, callerUserID), nil
}

func (a *memoryLobbyApplication) VoteReady(
	_ context.Context,
	matchID, callerUserID string,
	ready bool,
) (contract4chessraiders.LobbyView, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	match, ok := a.matches[matchID]
	if !ok {
		return contract4chessraiders.LobbyView{}, contract4chessraiders.ErrMatchNotFound
	}
	match.ready[callerUserID] = ready
	return a.viewLocked(matchID, match, callerUserID), nil
}

// Start is deliberately NOT creator-gated — any joined, ready player may
// trigger it, matching the real facade4chess.Service.Start rule this
// reference implementation exists to approximate.
func (a *memoryLobbyApplication) Start(_ context.Context, matchID, callerUserID string) (contract4chessraiders.LobbyView, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	match, ok := a.matches[matchID]
	if !ok {
		return contract4chessraiders.LobbyView{}, contract4chessraiders.ErrMatchNotFound
	}
	if !seatedIn(match.white, callerUserID) && !seatedIn(match.black, callerUserID) {
		return contract4chessraiders.LobbyView{}, contract4chessraiders.ErrPlayerNotFound
	}
	if len(match.white) == 0 || len(match.black) == 0 {
		return contract4chessraiders.LobbyView{}, contract4chessraiders.ErrTeamsNotViable
	}
	for _, seat := range append(append([]contract4chessraiders.LobbySeat{}, match.white...), match.black...) {
		if !match.ready[seat.UserID] {
			return contract4chessraiders.LobbyView{}, contract4chessraiders.ErrNotAllReady
		}
	}
	match.lifecycle = contract4chessraiders.LifecyclePlaying
	return a.viewLocked(matchID, match, callerUserID), nil
}

func (a *memoryLobbyApplication) Get(_ context.Context, matchID, callerUserID string) (contract4chessraiders.LobbyView, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	match, ok := a.matches[matchID]
	if !ok {
		return contract4chessraiders.LobbyView{}, contract4chessraiders.ErrMatchNotFound
	}
	return a.viewLocked(matchID, match, callerUserID), nil
}

func (a *memoryLobbyApplication) viewLocked(matchID string, match *memoryMatch, viewer string) contract4chessraiders.LobbyView {
	return contract4chessraiders.LobbyView{
		MatchID:       matchID,
		Lifecycle:     match.lifecycle,
		CaptureMode:   match.captureMode,
		CreatedBy:     match.createdBy,
		White:         match.white,
		Black:         match.black,
		ViewerIsReady: match.ready[viewer],
	}
}

func seatedIn(seats []contract4chessraiders.LobbySeat, userID string) bool {
	for _, seat := range seats {
		if seat.UserID == userID {
			return true
		}
	}
	return false
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	digits := []byte{}
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}

func TestCheckMatchLobbyApplicationAcceptsCompliantImplementation(t *testing.T) {
	violations := contract4chessraiderstest.CheckMatchLobbyApplication(
		func() contract4chessraiders.MatchLobbyApplication { return newMemoryLobbyApplication() },
	)
	for _, violation := range violations {
		t.Errorf("unexpected violation against a compliant implementation: %v", violation)
	}
}

// brokenLobbyApplication lets Start succeed with only one side seated,
// exercising the checker's ErrTeamsNotViable expectation.
type brokenLobbyApplication struct {
	*memoryLobbyApplication
}

func (b brokenLobbyApplication) Start(ctx context.Context, matchID, callerUserID string) (contract4chessraiders.LobbyView, error) {
	b.mu.Lock()
	match, ok := b.matches[matchID]
	b.mu.Unlock()
	if !ok {
		return contract4chessraiders.LobbyView{}, contract4chessraiders.ErrMatchNotFound
	}
	b.mu.Lock()
	match.lifecycle = contract4chessraiders.LifecyclePlaying
	b.mu.Unlock()
	return b.viewLocked(matchID, match, callerUserID), nil
}

func TestCheckMatchLobbyApplicationRejectsNonCompliantImplementation(t *testing.T) {
	violations := contract4chessraiderstest.CheckMatchLobbyApplication(
		func() contract4chessraiders.MatchLobbyApplication {
			return brokenLobbyApplication{newMemoryLobbyApplication()}
		},
	)
	if len(violations) == 0 {
		t.Fatal("want at least one violation against an implementation that starts with one side seated, got none")
	}
}
