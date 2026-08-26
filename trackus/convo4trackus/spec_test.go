package convo4trackus

import (
	"testing"

	"github.com/sneat-co/sneat-go-core/convospec"
)

// Every declared example must satisfy its own action's argument schema. A
// drifting example is a drifting prompt: examples are shown to the model.
func TestDeclaredExamplesValidate(t *testing.T) {
	for _, action := range Declaration().Actions {
		for i, example := range action.Examples {
			if _, err := action.ValidateArgs(example.Args); err != nil {
				t.Errorf("%s example %d (%q): %v", action.ID, i, example.UserText, err)
			}
		}
	}
}

// A rule naming an action or argument the catalog does not declare is caught
// here rather than showing up as a message that silently never matches.
func TestRulesAgreeWithDeclaration(t *testing.T) {
	if err := Declaration().ValidateRules(Rules()); err != nil {
		t.Fatal(err)
	}
}

// The pilot phrasings must resolve to the tracker their names claim.
func TestRulesResolvePilotPhrasings(t *testing.T) {
	declaration := Declaration()
	for _, tt := range []struct {
		text      string
		actionID  string
		trackerID string
		value     float64
	}{
		{"push-ups 20", addEntryDef.ID, PushUpsTrackerID, 20},
		{"20 push-ups", addEntryDef.ID, PushUpsTrackerID, 20},
		{"Push-ups: 20", addEntryDef.ID, PushUpsTrackerID, 20},
		{"10 pull-ups", addEntryDef.ID, PullUpsTrackerID, 10},
		{"pull-ups 10", addEntryDef.ID, PullUpsTrackerID, 10},
		{"30 squats", addEntryDef.ID, SquatsTrackerID, 30},
		{"my weight is 80.5", addEntryDef.ID, WeightTrackerID, 80.5},
	} {
		normalized := convospec.NormalizeText(tt.text)
		var matched bool
		for _, rule := range Rules() {
			def, ok := declaration.Action(rule.ActionID)
			if !ok {
				t.Fatalf("rule names undeclared action %s", rule.ActionID)
			}
			args, ok := rule.Match(normalized, def)
			if !ok {
				continue
			}
			matched = true
			if rule.ActionID != tt.actionID {
				t.Errorf("%q -> action %s, want %s", tt.text, rule.ActionID, tt.actionID)
			}
			if args["trackerID"] != tt.trackerID {
				t.Errorf("%q -> trackerID %v, want %s", tt.text, args["trackerID"], tt.trackerID)
			}
			if args["value"] != tt.value {
				t.Errorf("%q -> value %v, want %v", tt.text, args["value"], tt.value)
			}
			break
		}
		if !matched {
			t.Errorf("%q matched no rule", tt.text)
		}
	}
}

// "ran 10 km" has no standard tracker, so it must route to the confirmed
// creation action rather than inventing a tracker ID.
func TestRunningRoutesToConfirmedCreation(t *testing.T) {
	declaration := Declaration()
	normalized := convospec.NormalizeText("ran 10 km")
	for _, rule := range Rules() {
		def, _ := declaration.Action(rule.ActionID)
		args, ok := rule.Match(normalized, def)
		if !ok {
			continue
		}
		if rule.ActionID != createTrackerDef.ID {
			t.Fatalf(`"ran 10 km" -> %s, want %s`, rule.ActionID, createTrackerDef.ID)
		}
		if args["unit"] != "km" || args["value"] != 10.0 || args["title"] != "Running" {
			t.Errorf("args = %#v", args)
		}
		if !def.Confirm {
			t.Error("tracker creation must require confirmation")
		}
		return
	}
	t.Fatal(`"ran 10 km" matched no rule`)
}

// Triggers must be distinctive: a generic verb would narrow the action set on
// messages that belong to another extension, which is the one way a fail-open
// prefilter can misroute.
func TestTriggersAreDistinctive(t *testing.T) {
	generic := map[string]bool{
		"add": true, "did": true, "get": true, "new": true, "set": true,
		"buy": true, "have": true, "do": true, "make": true, "list": true,
	}
	for _, trigger := range Declaration().Triggers {
		if generic[trigger] {
			t.Errorf("trigger %q is too generic to narrow safely", trigger)
		}
	}
}

// The pilot messages must actually be claimed, or the prefilter can never
// narrow to Trackus.
func TestTriggersClaimPilotMessages(t *testing.T) {
	declaration := Declaration()
	for _, text := range []string{"push-ups 20", "20 push-ups", "10 pull-ups", "ran 10 km"} {
		if !declaration.MatchesTriggers(convospec.NormalizeText(text)) {
			t.Errorf("%q is not claimed by any declared trigger", text)
		}
	}
	// A Listus message must NOT be claimed, or every grocery item would narrow
	// to Trackus.
	for _, text := range []string{"milk 2 liters", "bought milk", "add bread to the shopping list"} {
		if declaration.MatchesTriggers(convospec.NormalizeText(text)) {
			t.Errorf("%q must not be claimed by Trackus triggers", text)
		}
	}
}

// The two Value constructors are the only way a caller should build a
// measurement, so they are worth pinning: sending a float where the tracker
// declares an integer is what made a rep count read back wrong.
func TestValueConstructors(t *testing.T) {
	reps := IntValue(20)
	if reps.Int != 20 || reps.IsFloat || reps.Float != 0 {
		t.Errorf("IntValue(20) = %+v, want an integer measurement", reps)
	}
	distance := FloatValue(10.5)
	if !distance.IsFloat || distance.Float != 10.5 || distance.Int != 0 {
		t.Errorf("FloatValue(10.5) = %+v, want a fractional measurement", distance)
	}
	// A whole number expressed as a float stays flagged as fractional: the
	// caller asked for a float tracker, and the flag is what the port reads.
	whole := FloatValue(10)
	if !whole.IsFloat {
		t.Error("FloatValue(10) must remain flagged as fractional")
	}
}
