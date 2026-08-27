// Package contract4chessraiders is the public contract for Chess Raiders —
// the Sneat platform's real-time multiplayer chess game. It holds types,
// interfaces and errors only: what a host, bot shell, or distribution portal
// needs in order to talk to Chess Raiders without depending on the private
// engine (sneat-co/chessraiders).
//
// It deliberately does not import server-go/chess, facade4chess, or any
// other private chessraiders package — every type below is independently
// declared. It carries no board state, move, or rules-evaluation concept:
// gameplay stays exclusively on the private web/Mini App surface. What
// crosses here is the lobby/social-shell boundary (creating, joining,
// starting, and watching the status of a match), external portal identity,
// and distribution-platform session binding.
package contract4chessraiders

// ExtensionID is Chess Raiders' stable global namespace — the same role
// contract4competios.ExtensionID plays for Competios. Hosts, adapters, and
// sibling extensions use it to key Chess Raiders data or capabilities
// without hardcoding the string "chess-raiders" in more than one place (it
// was previously duplicated as a bare literal, e.g. in
// contract4competios.DefaultChessRaidersDraftContent).
const ExtensionID = "chess-raiders"

// Side identifies a team in a match. It is independent of the private
// engine's own colour representation.
type Side string

const (
	SideNone  Side = ""
	SideWhite Side = "w"
	SideBlack Side = "b"
)

// CaptureMode is the lobby-level capture rule a creator selects before a
// match starts. It says nothing about how an individual capture resolves on
// the board — that stays inside the private engine.
type CaptureMode string

const (
	CaptureModeKillOnly      CaptureMode = "kill-only"
	CaptureModeCaptureOnly   CaptureMode = "capture-only"
	CaptureModeKillOrCapture CaptureMode = "kill-or-capture"
)

// Lifecycle is the coarse match phase visible to a social/bot shell. It has
// no substates: tactical phase detail (charge cooldowns, wall repair
// progress, ...) stays on the private web/Mini App gameplay surface.
type Lifecycle string

const (
	LifecycleLobby    Lifecycle = "lobby"
	LifecycleStarting Lifecycle = "starting"
	LifecyclePlaying  Lifecycle = "playing"
	LifecycleEnded    Lifecycle = "ended"
)

// Player identifies one human participant on a social/bot shell. It carries
// no Telegram, portal, or authentication detail — a host resolves those
// independently and constructs Player from its own trusted source; a caller
// can never assert another player's identity through this type.
type Player struct {
	UserID      string `json:"userID"`
	DisplayName string `json:"displayName,omitempty"`
	Side        Side   `json:"side,omitempty"`
}

// LobbySeat is one seat's read-only projection inside LobbyView.
type LobbySeat struct {
	UserID      string `json:"userID"`
	DisplayName string `json:"displayName,omitempty"`
	Ready       bool   `json:"ready,omitempty"`
}

// LobbyView is the presentation-neutral snapshot a social/bot shell renders.
// It carries no board state, legal moves, fog, morale, or other tactical
// detail — those belong exclusively to the private web/Mini App gameplay
// surface. ViewerSide and ViewerIsReady are resolved by the implementation
// for the caller's own identity; they are never client-supplied.
type LobbyView struct {
	MatchID       string      `json:"matchID"`
	Lifecycle     Lifecycle   `json:"lifecycle"`
	CaptureMode   CaptureMode `json:"captureMode"`
	CreatedBy     string      `json:"createdBy"`
	Revision      int64       `json:"revision"`
	White         []LobbySeat `json:"white,omitempty"`
	Black         []LobbySeat `json:"black,omitempty"`
	Winner        Side        `json:"winner,omitempty"`
	StartsAtMs    int64       `json:"startsAtMs,omitempty"`
	ViewerSide    Side        `json:"viewerSide,omitempty"`
	ViewerIsReady bool        `json:"viewerIsReady,omitempty"`
}
