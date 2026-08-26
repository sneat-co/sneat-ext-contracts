package eventtimeline

import "sort"

// State is the deterministic projection of an event log up to a point in order.
// Two consumers folding the same events MUST obtain an identical State.
type State struct {
	Status       GameStatus       `json:"status"`
	Period       int              `json:"period"`
	GameClockMs  int64            `json:"gameClockMs"`
	ClockRunning bool             `json:"clockRunning"`
	Scores       map[TeamSide]int `json:"scores"`
	TeamFouls    map[TeamSide]int `json:"teamFouls"`
	TimeoutsUsed map[TeamSide]int `json:"timeoutsUsed"`
	Possession   TeamSide         `json:"possession"`
	OnCourt      map[TeamSide][]string `json:"onCourt"`

	// Per-player box score (post-game-recap): points → assists, keyed by playerID.
	PlayerPoints  map[string]int `json:"playerPoints"`
	PlayerAssists map[string]int `json:"playerAssists"`
}

func newState() State {
	return State{
		Status:        StatusScheduled,
		Scores:        map[TeamSide]int{SideHome: 0, SideAway: 0},
		TeamFouls:     map[TeamSide]int{SideHome: 0, SideAway: 0},
		TimeoutsUsed:  map[TeamSide]int{SideHome: 0, SideAway: 0},
		OnCourt:       map[TeamSide][]string{SideHome: {}, SideAway: {}},
		PlayerPoints:  map[string]int{},
		PlayerAssists: map[string]int{},
	}
}

// Order returns events in canonical total order: WallClockMs ascending, ties
// broken by EventID. Idempotent duplicates (same EventID) are collapsed to the
// first occurrence, so a replayed append never double-counts. The input slice
// is not mutated.
func Order(events []Event) []Event {
	seen := make(map[string]struct{}, len(events))
	out := make([]Event, 0, len(events))
	for _, e := range events {
		if _, dup := seen[e.EventID]; dup {
			continue
		}
		seen[e.EventID] = struct{}{}
		out = append(out, e)
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].WallClockMs != out[j].WallClockMs {
			return out[i].WallClockMs < out[j].WallClockMs
		}
		return out[i].EventID < out[j].EventID
	})
	return out
}

// Fold projects current State from the log. It is deterministic and pure:
// equal (ordered, deduped) inputs yield equal output regardless of the input
// order in which duplicates/out-of-order events arrive.
func Fold(events []Event) State {
	ordered := Order(events)

	// Corrections are appended, never destructive: collect the set of voided
	// EventIDs first, then fold skipping any voided event.
	voided := make(map[string]struct{})
	for _, e := range ordered {
		if e.Type == EventCorrection && e.CorrectionOf != "" {
			voided[e.CorrectionOf] = struct{}{}
		}
	}

	s := newState()
	for _, e := range ordered {
		if _, isVoid := voided[e.EventID]; isVoid {
			continue
		}
		applyEvent(&s, e)
	}
	return s
}

func applyEvent(s *State, e Event) {
	switch e.Type {
	case EventStatus:
		s.Status = e.Status
	case EventPeriod:
		s.Period = e.Period
	case EventClock:
		switch e.ClockAction {
		case ClockStart:
			s.ClockRunning = true
			if e.GameClockMs > 0 {
				s.GameClockMs = e.GameClockMs
			}
		case ClockStop:
			s.ClockRunning = false
			s.GameClockMs = e.GameClockMs
		case ClockAdjust:
			s.GameClockMs = e.GameClockMs
		}
	case EventScore:
		if e.Side.Valid() && e.Points > 0 {
			s.Scores[e.Side] += e.Points
			if e.ScorerID != "" {
				s.PlayerPoints[e.ScorerID] += e.Points
			}
			if e.AssistID != "" {
				s.PlayerAssists[e.AssistID]++
			}
		}
	case EventTeamFoul:
		if e.Side.Valid() {
			s.TeamFouls[e.Side]++
		}
	case EventTimeout:
		if e.Side.Valid() {
			s.TimeoutsUsed[e.Side]++
		}
	case EventPossession:
		if e.Side.Valid() {
			s.Possession = e.Side
		}
	case EventSubstitution:
		if e.Side.Valid() {
			if e.PlayerOff != "" {
				s.OnCourt[e.Side] = removePlayer(s.OnCourt[e.Side], e.PlayerOff)
			}
			if e.PlayerOn != "" {
				s.OnCourt[e.Side] = addPlayer(s.OnCourt[e.Side], e.PlayerOn)
			}
		}
	case EventCorrection:
		// A correction that carries a replacement score amends (counts the new
		// value); the voided original was already excluded above.
		if e.Side.Valid() && e.Points > 0 {
			s.Scores[e.Side] += e.Points
		}
	case EventJudgeRuling:
		// Judge rulings that change state do so by carrying a correction; a bare
		// ruling is an audit entry with no projection effect.
	}
}

func addPlayer(list []string, id string) []string {
	for _, p := range list {
		if p == id {
			return list
		}
	}
	return append(list, id)
}

func removePlayer(list []string, id string) []string {
	out := make([]string, 0, len(list))
	for _, p := range list {
		if p != id {
			out = append(out, p)
		}
	}
	return out
}

// InBonus reports whether side is shooting bonus free throws: true once the
// OPPONENT has committed at least limit team fouls in the period.
func (s State) InBonus(side TeamSide, limit int) bool {
	opp := SideAway
	if side == SideAway {
		opp = SideHome
	}
	return s.TeamFouls[opp] >= limit
}
