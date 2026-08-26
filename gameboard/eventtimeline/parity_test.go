package eventtimeline

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// parityCase is one shared oracle entry: a log and the State both the Go reducer
// (here) and the TS reducer (@sneat/extension-gameboard-contract) MUST produce.
type parityCase struct {
	Name     string  `json:"name"`
	Events   []Event `json:"events"`
	Expected State   `json:"expected"`
}

// parityScenarios are the canonical fold scenarios. The committed fixture
// ../parity/parity.json is generated from these via Fold (set GB_REGEN=1) and
// is the cross-language contract: the TS reducer test asserts against the same
// file, so the two reducers are proven to agree.
//
// Provenance note: in the source repo (sneat-co/ext-gameboard) this fixture
// lived at repo-root parity/parity.json, two levels up from
// backend/eventtimeline/. Here the Go module directory (gameboard/) plays the
// role the old repo root played, so the fixture now lives one level up, at
// gameboard/parity/parity.json — parityPath below was adjusted from
// filepath.Join("..", "..", "parity", "parity.json") to
// filepath.Join("..", "parity", "parity.json") accordingly. Content is
// otherwise byte-identical to the source.
func parityScenarios() []parityCase {
	mk := func(name string, evs []Event) parityCase {
		return parityCase{Name: name, Events: evs, Expected: Fold(evs)}
	}
	return []parityCase{
		mk("empty", []Event{}),
		mk("basic-scores", []Event{
			{EventID: "a", Type: EventScore, Source: SourceScorekeeper, WallClockMs: 1, Side: SideHome, Points: 2},
			{EventID: "b", Type: EventScore, Source: SourceScorekeeper, WallClockMs: 2, Side: SideAway, Points: 3},
			{EventID: "c", Type: EventScore, Source: SourceScorekeeper, WallClockMs: 3, Side: SideHome, Points: 1},
		}),
		mk("idempotent-and-order", []Event{
			{EventID: "x", Type: EventScore, Source: SourceScorekeeper, WallClockMs: 5, Side: SideHome, Points: 2},
			{EventID: "x", Type: EventScore, Source: SourceScorekeeper, WallClockMs: 5, Side: SideHome, Points: 2},
			{EventID: "w", Type: EventScore, Source: SourceScorekeeper, WallClockMs: 1, Side: SideAway, Points: 3},
		}),
		mk("correction-void", []Event{
			{EventID: "good", Type: EventScore, Source: SourceScorekeeper, WallClockMs: 1, Side: SideHome, Points: 2},
			{EventID: "oops", Type: EventScore, Source: SourceScorekeeper, WallClockMs: 2, Side: SideHome, Points: 3},
			{EventID: "fix", Type: EventCorrection, Source: SourceJudge, WallClockMs: 3, CorrectionOf: "oops"},
		}),
		mk("correction-amend", []Event{
			{EventID: "oops", Type: EventScore, Source: SourceScorekeeper, WallClockMs: 1, Side: SideHome, Points: 3},
			{EventID: "fix", Type: EventCorrection, Source: SourceJudge, WallClockMs: 2, CorrectionOf: "oops", Side: SideHome, Points: 2},
		}),
		mk("full-vocabulary", []Event{
			{EventID: "e1", Type: EventStatus, Source: SourceTimekeeper, WallClockMs: 1, Status: StatusLive},
			{EventID: "e2", Type: EventPeriod, Source: SourceTimekeeper, WallClockMs: 2, Period: 1},
			{EventID: "e3", Type: EventClock, Source: SourceTimekeeper, WallClockMs: 3, ClockAction: ClockStart, GameClockMs: 600000},
			{EventID: "e4", Type: EventSubstitution, Source: SourceScorekeeper, WallClockMs: 4, Side: SideHome, PlayerOn: "p1"},
			{EventID: "e5", Type: EventSubstitution, Source: SourceScorekeeper, WallClockMs: 5, Side: SideHome, PlayerOn: "p2"},
			{EventID: "e6", Type: EventScore, Source: SourceScorekeeper, WallClockMs: 6, Side: SideHome, Points: 2, ScorerID: "p1", AssistID: "p2"},
			{EventID: "e7", Type: EventTeamFoul, Source: SourceScorekeeper, WallClockMs: 7, Side: SideAway},
			{EventID: "e8", Type: EventTimeout, Source: SourceTimekeeper, WallClockMs: 8, Side: SideHome},
			{EventID: "e9", Type: EventPossession, Source: SourceTimekeeper, WallClockMs: 9, Side: SideAway},
			{EventID: "e10", Type: EventSubstitution, Source: SourceScorekeeper, WallClockMs: 10, Side: SideHome, PlayerOff: "p1", PlayerOn: "p3"},
			{EventID: "e11", Type: EventClock, Source: SourceTimekeeper, WallClockMs: 11, ClockAction: ClockStop, GameClockMs: 540000},
		}),
	}
}

func parityPath(t *testing.T) string {
	t.Helper()
	return filepath.Join("..", "parity", "parity.json")
}

func TestParityFixture(t *testing.T) {
	cases := parityScenarios()
	data, err := json.MarshalIndent(cases, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	path := parityPath(t)

	if os.Getenv("GB_REGEN") == "1" {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
			t.Fatal(err)
		}
		t.Logf("regenerated %s", path)
		return
	}

	committed, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read parity fixture (run with GB_REGEN=1 to generate): %v", err)
	}
	var want []parityCase
	if err := json.Unmarshal(committed, &want); err != nil {
		t.Fatal(err)
	}
	// The committed fixture's expected states MUST equal a fresh Go fold.
	for i, c := range want {
		got := Fold(c.Events)
		if !reflect.DeepEqual(got, c.Expected) {
			t.Errorf("case %q (#%d): Go fold drifted from fixture:\n got=%#v\nwant=%#v", c.Name, i, got, c.Expected)
		}
	}
	if len(want) != len(cases) {
		t.Errorf("fixture has %d cases, scenarios define %d — run GB_REGEN=1", len(want), len(cases))
	}
}
