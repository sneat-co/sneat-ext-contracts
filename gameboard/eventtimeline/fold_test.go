package eventtimeline

import (
	"reflect"
	"testing"
)

func score(id string, ms int64, side TeamSide, pts int) Event {
	return Event{EventID: id, Type: EventScore, Source: SourceScorekeeper, WallClockMs: ms, Side: side, Points: pts}
}

// AC: append-immutable-ordered / live-stream-in-order — total order is by
// wall-clock, ties broken by EventID, regardless of arrival order.
func TestOrder_ByWallClockThenEventID(t *testing.T) {
	in := []Event{
		{EventID: "b", WallClockMs: 10},
		{EventID: "a", WallClockMs: 10},
		{EventID: "c", WallClockMs: 5},
	}
	got := Order(in)
	want := []string{"c", "a", "b"}
	for i, id := range want {
		if got[i].EventID != id {
			t.Fatalf("position %d: got %s want %s (%v)", i, got[i].EventID, id, got)
		}
	}
	// input not mutated
	if in[0].EventID != "b" {
		t.Fatalf("input slice was mutated: %v", in)
	}
}

// AC: append-immutable-ordered — idempotency: a replayed append with the same
// EventID is collapsed to a single contribution.
func TestOrder_IdempotentDedupe(t *testing.T) {
	in := []Event{score("x", 1, SideHome, 2), score("x", 1, SideHome, 2)}
	if got := Order(in); len(got) != 1 {
		t.Fatalf("expected dedupe to 1 event, got %d", len(got))
	}
	if s := Fold(in); s.Scores[SideHome] != 2 {
		t.Fatalf("idempotent replay double-counted: %d", s.Scores[SideHome])
	}
}

// AC: state-is-deterministic-fold — two independent folds of the same log
// (shuffled) yield identical state, and score == sum of score events.
func TestFold_DeterministicAndScoreIsSum(t *testing.T) {
	a := []Event{
		score("s1", 1, SideHome, 2),
		score("s2", 2, SideAway, 3),
		score("s3", 3, SideHome, 1),
	}
	b := []Event{a[2], a[0], a[1]} // shuffled arrival
	sa, sb := Fold(a), Fold(b)
	if !reflect.DeepEqual(sa, sb) {
		t.Fatalf("folds differ:\n%#v\n%#v", sa, sb)
	}
	if sa.Scores[SideHome] != 3 || sa.Scores[SideAway] != 3 {
		t.Fatalf("score mismatch: %+v", sa.Scores)
	}
}

// AC: vocabulary-covers-official-events — every event type projects.
func TestFold_FullVocabulary(t *testing.T) {
	events := []Event{
		{EventID: "e1", Type: EventStatus, WallClockMs: 1, Status: StatusLive},
		{EventID: "e2", Type: EventPeriod, WallClockMs: 2, Period: 1},
		{EventID: "e3", Type: EventClock, WallClockMs: 3, ClockAction: ClockStart, GameClockMs: 600000},
		{EventID: "e4", Type: EventSubstitution, WallClockMs: 4, Side: SideHome, PlayerOn: "p1"},
		{EventID: "e5", Type: EventScore, WallClockMs: 5, Side: SideHome, Points: 2, ScorerID: "p1", AssistID: "p2"},
		{EventID: "e6", Type: EventTeamFoul, WallClockMs: 6, Side: SideAway},
		{EventID: "e7", Type: EventTimeout, WallClockMs: 7, Side: SideHome},
		{EventID: "e8", Type: EventPossession, WallClockMs: 8, Side: SideAway},
		{EventID: "e9", Type: EventClock, WallClockMs: 9, ClockAction: ClockStop, GameClockMs: 540000},
		{EventID: "e10", Type: EventJudgeRuling, WallClockMs: 10},
		{EventID: "e11", Type: EventClock, WallClockMs: 11, ClockAction: ClockAdjust, GameClockMs: 530000},
	}
	s := Fold(events)
	if s.Status != StatusLive || s.Period != 1 {
		t.Fatalf("status/period: %+v", s)
	}
	if s.GameClockMs != 530000 || s.ClockRunning {
		t.Fatalf("clock: running=%v ms=%d", s.ClockRunning, s.GameClockMs)
	}
	if s.Scores[SideHome] != 2 {
		t.Fatalf("score: %+v", s.Scores)
	}
	if s.TeamFouls[SideAway] != 1 || s.TimeoutsUsed[SideHome] != 1 {
		t.Fatalf("fouls/timeouts: %+v %+v", s.TeamFouls, s.TimeoutsUsed)
	}
	if s.Possession != SideAway {
		t.Fatalf("possession: %v", s.Possession)
	}
	if !reflect.DeepEqual(s.OnCourt[SideHome], []string{"p1"}) {
		t.Fatalf("oncourt: %+v", s.OnCourt)
	}
}

// AC: correction-is-appended-not-edited — a void correction removes the
// original's contribution; history is retained (the event still in the log).
func TestFold_CorrectionVoids(t *testing.T) {
	events := []Event{
		score("good", 1, SideHome, 2),
		score("oops", 2, SideHome, 3),
		{EventID: "fix", Type: EventCorrection, Source: SourceJudge, WallClockMs: 3, CorrectionOf: "oops"},
	}
	s := Fold(events)
	if s.Scores[SideHome] != 2 {
		t.Fatalf("void correction not applied: %d", s.Scores[SideHome])
	}
	// original retained in the ordered log (audit)
	found := false
	for _, e := range Order(events) {
		if e.EventID == "oops" {
			found = true
		}
	}
	if !found {
		t.Fatal("voided event must remain in the log for audit")
	}
}

// AC: correction-is-appended-not-edited — an amend correction counts the new value.
func TestFold_CorrectionAmends(t *testing.T) {
	events := []Event{
		score("oops", 1, SideHome, 3),
		{EventID: "fix", Type: EventCorrection, Source: SourceJudge, WallClockMs: 2, CorrectionOf: "oops", Side: SideHome, Points: 2},
	}
	if s := Fold(events); s.Scores[SideHome] != 2 {
		t.Fatalf("amend should count 2, got %d", s.Scores[SideHome])
	}
}

func TestFold_SubstitutionOnOff(t *testing.T) {
	events := []Event{
		{EventID: "a", Type: EventSubstitution, WallClockMs: 1, Side: SideHome, PlayerOn: "p1"},
		{EventID: "b", Type: EventSubstitution, WallClockMs: 2, Side: SideHome, PlayerOn: "p2"},
		{EventID: "c", Type: EventSubstitution, WallClockMs: 3, Side: SideHome, PlayerOn: "p2"}, // dup on -> no double
		{EventID: "d", Type: EventSubstitution, WallClockMs: 4, Side: SideHome, PlayerOff: "p1", PlayerOn: "p3"},
	}
	s := Fold(events)
	if !reflect.DeepEqual(s.OnCourt[SideHome], []string{"p2", "p3"}) {
		t.Fatalf("oncourt: %+v", s.OnCourt[SideHome])
	}
}

// Box score: scorer points and assister assists accrue per player (post-game-recap).
func TestFold_BoxScore(t *testing.T) {
	events := []Event{
		{EventID: "a", Type: EventScore, WallClockMs: 1, Side: SideHome, Points: 2, ScorerID: "p1", AssistID: "p2"},
		{EventID: "b", Type: EventScore, WallClockMs: 2, Side: SideHome, Points: 3, ScorerID: "p1"},
		{EventID: "c", Type: EventScore, WallClockMs: 3, Side: SideAway, Points: 1}, // no scorer attributed
	}
	s := Fold(events)
	if s.PlayerPoints["p1"] != 5 {
		t.Fatalf("p1 points = %d, want 5", s.PlayerPoints["p1"])
	}
	if s.PlayerAssists["p2"] != 1 {
		t.Fatalf("p2 assists = %d, want 1", s.PlayerAssists["p2"])
	}
	if len(s.PlayerPoints) != 1 { // the unattributed away point adds no player line
		t.Fatalf("unexpected player points: %+v", s.PlayerPoints)
	}
}

func TestInBonus(t *testing.T) {
	s := newState()
	s.TeamFouls[SideHome] = 5
	if !s.InBonus(SideAway, 5) {
		t.Fatal("away should be in bonus when home has 5 fouls")
	}
	if s.InBonus(SideHome, 5) {
		t.Fatal("home should not be in bonus")
	}
	if !s.InBonus(SideHome, 0) {
		t.Fatal("limit 0 -> always bonus")
	}
}

func TestTeamSideValid(t *testing.T) {
	if !SideHome.Valid() || !SideAway.Valid() || TeamSide("x").Valid() {
		t.Fatal("Valid() wrong")
	}
}

// Ignored/no-op branches: invalid side and zero points must not panic or count.
func TestFold_IgnoresInvalidPayloads(t *testing.T) {
	events := []Event{
		{EventID: "a", Type: EventScore, WallClockMs: 1, Side: "bogus", Points: 2},
		{EventID: "b", Type: EventScore, WallClockMs: 2, Side: SideHome, Points: 0},
		{EventID: "c", Type: EventTeamFoul, WallClockMs: 3, Side: "bogus"},
		{EventID: "d", Type: EventClock, WallClockMs: 4, ClockAction: ClockStart}, // no ms -> keep 0
	}
	s := Fold(events)
	if s.Scores[SideHome] != 0 || s.TeamFouls[SideHome] != 0 || s.TeamFouls[SideAway] != 0 {
		t.Fatalf("invalid payloads counted: %+v", s)
	}
	if !s.ClockRunning || s.GameClockMs != 0 {
		t.Fatalf("clock start without ms: running=%v ms=%d", s.ClockRunning, s.GameClockMs)
	}
}
