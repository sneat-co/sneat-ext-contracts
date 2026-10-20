package contract4chessraiderstest_test

import (
	"context"
	"errors"
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

type mockLobbyApplication struct {
	*memoryLobbyApplication
	onCreateLobby func(ctx context.Context, creator contract4chessraiders.Player, mode contract4chessraiders.CaptureMode) (contract4chessraiders.LobbyView, error)
	onJoin        func(ctx context.Context, matchID string, joiner contract4chessraiders.Player) (contract4chessraiders.LobbyView, error)
	onVoteReady   func(ctx context.Context, matchID, callerUserID string, ready bool) (contract4chessraiders.LobbyView, error)
	onStart       func(ctx context.Context, matchID, callerUserID string) (contract4chessraiders.LobbyView, error)
	onGet         func(ctx context.Context, matchID, callerUserID string) (contract4chessraiders.LobbyView, error)
}

func (m *mockLobbyApplication) CreateLobby(ctx context.Context, creator contract4chessraiders.Player, mode contract4chessraiders.CaptureMode) (contract4chessraiders.LobbyView, error) {
	if m.onCreateLobby != nil {
		return m.onCreateLobby(ctx, creator, mode)
	}
	return m.memoryLobbyApplication.CreateLobby(ctx, creator, mode)
}

func (m *mockLobbyApplication) Join(ctx context.Context, matchID string, joiner contract4chessraiders.Player) (contract4chessraiders.LobbyView, error) {
	if m.onJoin != nil {
		return m.onJoin(ctx, matchID, joiner)
	}
	return m.memoryLobbyApplication.Join(ctx, matchID, joiner)
}

func (m *mockLobbyApplication) VoteReady(ctx context.Context, matchID, callerUserID string, ready bool) (contract4chessraiders.LobbyView, error) {
	if m.onVoteReady != nil {
		return m.onVoteReady(ctx, matchID, callerUserID, ready)
	}
	return m.memoryLobbyApplication.VoteReady(ctx, matchID, callerUserID, ready)
}

func (m *mockLobbyApplication) Start(ctx context.Context, matchID, callerUserID string) (contract4chessraiders.LobbyView, error) {
	if m.onStart != nil {
		return m.onStart(ctx, matchID, callerUserID)
	}
	return m.memoryLobbyApplication.Start(ctx, matchID, callerUserID)
}

func (m *mockLobbyApplication) Get(ctx context.Context, matchID, callerUserID string) (contract4chessraiders.LobbyView, error) {
	if m.onGet != nil {
		return m.onGet(ctx, matchID, callerUserID)
	}
	return m.memoryLobbyApplication.Get(ctx, matchID, callerUserID)
}

func TestCheckMatchLobbyApplicationRejectionCases(t *testing.T) {
	tests := []struct {
		name string
		mock func() contract4chessraiders.MatchLobbyApplication
	}{
		{
			name: "create lobby returns error",
			mock: func() contract4chessraiders.MatchLobbyApplication {
				return &mockLobbyApplication{
					memoryLobbyApplication: newMemoryLobbyApplication(),
					onCreateLobby: func(_ context.Context, _ contract4chessraiders.Player, _ contract4chessraiders.CaptureMode) (contract4chessraiders.LobbyView, error) {
						return contract4chessraiders.LobbyView{}, errors.New("cannot create")
					},
				}
			},
		},
		{
			name: "create lobby returns empty match ID",
			mock: func() contract4chessraiders.MatchLobbyApplication {
				return &mockLobbyApplication{
					memoryLobbyApplication: newMemoryLobbyApplication(),
					onCreateLobby: func(ctx context.Context, c contract4chessraiders.Player, m contract4chessraiders.CaptureMode) (contract4chessraiders.LobbyView, error) {
						mem := newMemoryLobbyApplication()
						v, err := mem.CreateLobby(ctx, c, m)
						v.MatchID = ""
						return v, err
					},
				}
			},
		},
		{
			name: "create lobby returns non-lobby lifecycle",
			mock: func() contract4chessraiders.MatchLobbyApplication {
				return &mockLobbyApplication{
					memoryLobbyApplication: newMemoryLobbyApplication(),
					onCreateLobby: func(ctx context.Context, c contract4chessraiders.Player, m contract4chessraiders.CaptureMode) (contract4chessraiders.LobbyView, error) {
						mem := newMemoryLobbyApplication()
						v, err := mem.CreateLobby(ctx, c, m)
						v.Lifecycle = contract4chessraiders.LifecyclePlaying
						return v, err
					},
				}
			},
		},
		{
			name: "create lobby returns wrong created by",
			mock: func() contract4chessraiders.MatchLobbyApplication {
				return &mockLobbyApplication{
					memoryLobbyApplication: newMemoryLobbyApplication(),
					onCreateLobby: func(ctx context.Context, c contract4chessraiders.Player, m contract4chessraiders.CaptureMode) (contract4chessraiders.LobbyView, error) {
						mem := newMemoryLobbyApplication()
						v, err := mem.CreateLobby(ctx, c, m)
						v.CreatedBy = "other"
						return v, err
					},
				}
			},
		},
		{
			name: "join returns invalid black seats",
			mock: func() contract4chessraiders.MatchLobbyApplication {
				mem := newMemoryLobbyApplication()
				return &mockLobbyApplication{
					memoryLobbyApplication: mem,
					onJoin: func(ctx context.Context, matchID string, joiner contract4chessraiders.Player) (contract4chessraiders.LobbyView, error) {
						v, err := mem.Join(ctx, matchID, joiner)
						v.Black = nil
						return v, err
					},
				}
			},
		},
		{
			name: "creator vote ready fails",
			mock: func() contract4chessraiders.MatchLobbyApplication {
				mem := newMemoryLobbyApplication()
				return &mockLobbyApplication{
					memoryLobbyApplication: mem,
					onVoteReady: func(ctx context.Context, matchID, callerUserID string, ready bool) (contract4chessraiders.LobbyView, error) {
						if callerUserID == "creator" {
							return contract4chessraiders.LobbyView{}, errors.New("vote error")
						}
						return mem.VoteReady(ctx, matchID, callerUserID, ready)
					},
				}
			},
		},
		{
			name: "opponent vote ready fails",
			mock: func() contract4chessraiders.MatchLobbyApplication {
				mem := newMemoryLobbyApplication()
				return &mockLobbyApplication{
					memoryLobbyApplication: mem,
					onVoteReady: func(ctx context.Context, matchID, callerUserID string, ready bool) (contract4chessraiders.LobbyView, error) {
						if callerUserID == "opponent" {
							return contract4chessraiders.LobbyView{}, errors.New("vote error")
						}
						return mem.VoteReady(ctx, matchID, callerUserID, ready)
					},
				}
			},
		},
		{
			name: "start fails for non-creator",
			mock: func() contract4chessraiders.MatchLobbyApplication {
				mem := newMemoryLobbyApplication()
				return &mockLobbyApplication{
					memoryLobbyApplication: mem,
					onStart: func(ctx context.Context, matchID, callerUserID string) (contract4chessraiders.LobbyView, error) {
						if callerUserID == "opponent" {
							return contract4chessraiders.LobbyView{}, errors.New("start error")
						}
						return mem.Start(ctx, matchID, callerUserID)
					},
				}
			},
		},
		{
			name: "start leaves lifecycle as lobby",
			mock: func() contract4chessraiders.MatchLobbyApplication {
				mem := newMemoryLobbyApplication()
				return &mockLobbyApplication{
					memoryLobbyApplication: mem,
					onStart: func(ctx context.Context, matchID, callerUserID string) (contract4chessraiders.LobbyView, error) {
						v, err := mem.Start(ctx, matchID, callerUserID)
						v.Lifecycle = contract4chessraiders.LifecycleLobby
						return v, err
					},
				}
			},
		},
		{
			name: "join after start returns non-ErrNotInLobby error",
			mock: func() contract4chessraiders.MatchLobbyApplication {
				mem := newMemoryLobbyApplication()
				return &mockLobbyApplication{
					memoryLobbyApplication: mem,
					onJoin: func(ctx context.Context, matchID string, joiner contract4chessraiders.Player) (contract4chessraiders.LobbyView, error) {
						if joiner.UserID == "stranger" {
							return contract4chessraiders.LobbyView{}, errors.New("unexpected error")
						}
						return mem.Join(ctx, matchID, joiner)
					},
				}
			},
		},
		{
			name: "get unknown match returns non-ErrMatchNotFound error",
			mock: func() contract4chessraiders.MatchLobbyApplication {
				mem := newMemoryLobbyApplication()
				return &mockLobbyApplication{
					memoryLobbyApplication: mem,
					onGet: func(ctx context.Context, matchID, callerUserID string) (contract4chessraiders.LobbyView, error) {
						return contract4chessraiders.LobbyView{}, errors.New("unexpected error")
					},
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			violations := contract4chessraiderstest.CheckMatchLobbyApplication(tc.mock)
			if len(violations) == 0 {
				t.Fatalf("case %q: expected violations, got none", tc.name)
			}
		})
	}
}
