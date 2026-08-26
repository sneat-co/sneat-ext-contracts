// Package eventtimeline is the frozen contract for the gameboard event-timeline:
// the append-only, ordered, timestamped event log that IS a basketball game's
// official record, plus the deterministic fold that projects current state.
//
// This package is hand-implemented to match the frozen wire contract
// (api4gameboard.tsp, kept in sneat-co/ext-gameboard) and is mirrored by the
// TS reducer in @sneat/extension-gameboard-contract (../../libs/gameboard);
// the two MUST fold the same log to identical state (parity is the contract).
package eventtimeline

// TeamSide identifies one of the two sides of a game.
type TeamSide string

const (
	// SideHome is the home side.
	SideHome TeamSide = "home"
	// SideAway is the away side.
	SideAway TeamSide = "away"
)

// Valid reports whether s is a known side.
func (s TeamSide) Valid() bool { return s == SideHome || s == SideAway }

// Side is the inline team descriptor stored per side on a game.
// SpaceID is nil for an ad-hoc name (first-use-backprop fills it later); set
// when the side is a real sneat team space.
type Side struct {
	Name    string  `json:"name"`
	Colour  string  `json:"colour"`
	SpaceID *string `json:"spaceID,omitempty"`
}

// EventType is the official event vocabulary.
type EventType string

const (
	EventStatus       EventType = "status"
	EventPeriod       EventType = "period"
	EventClock        EventType = "clock"
	EventScore        EventType = "score"
	EventTeamFoul     EventType = "team-foul"
	EventTimeout      EventType = "timeout"
	EventSubstitution EventType = "substitution"
	EventPossession   EventType = "possession"
	EventJudgeRuling  EventType = "judge-ruling"
	EventCorrection   EventType = "correction"
)

// GameStatus is the lifecycle status of a game.
type GameStatus string

const (
	StatusScheduled GameStatus = "scheduled"
	StatusLive      GameStatus = "live"
	StatusHalftime  GameStatus = "halftime"
	StatusOvertime  GameStatus = "overtime"
	StatusFinal     GameStatus = "final"
	StatusCancelled GameStatus = "cancelled"
)

// ClockAction is the kind of a clock event.
type ClockAction string

const (
	ClockStart  ClockAction = "start"
	ClockStop   ClockAction = "stop"
	ClockAdjust ClockAction = "adjust"
)

// Source is the authorized appender of an event (a per-game role or consensus).
type Source string

const (
	SourceScorekeeper Source = "scorekeeper" // appends score, team-foul, substitution
	SourceTimekeeper  Source = "timekeeper"  // appends clock, period, possession, timeout, status
	SourceJudge       Source = "judge"       // appends judge-ruling, correction
	SourceConsensus   Source = "consensus"   // appends when no official crew
)

// Event is one immutable entry of the log.
//
// EventID is a client-generated random dashless id used as the Firestore doc
// key and as the idempotency key (a replayed append with the same id is a
// no-op). Total order is by WallClockMs ascending, ties broken by EventID.
type Event struct {
	EventID    string    `json:"eventID"`
	Type       EventType `json:"type"`
	Source     Source    `json:"source"`
	WallClockMs int64    `json:"wallClockMs"`
	Period     int       `json:"period"`
	GameClockMs int64    `json:"gameClockMs"` // remaining time on the game clock

	// Payload (by type; only the fields relevant to Type are set):
	Status      GameStatus  `json:"status,omitempty"`
	ClockAction ClockAction `json:"clockAction,omitempty"`
	Side        TeamSide    `json:"side,omitempty"`
	Points      int         `json:"points,omitempty"`   // score: 1, 2 or 3
	ScorerID    string      `json:"scorerID,omitempty"`  // optional
	AssistID    string      `json:"assistID,omitempty"`  // optional
	PlayerOn    string      `json:"playerOn,omitempty"`  // substitution
	PlayerOff   string      `json:"playerOff,omitempty"` // substitution

	// CorrectionOf references the EventID this event voids/amends (corrections only).
	CorrectionOf string `json:"correctionOf,omitempty"`
}
