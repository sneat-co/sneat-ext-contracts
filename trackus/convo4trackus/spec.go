package convo4trackus

import (
	"regexp"

	"github.com/sneat-co/sneat-go-core/convospec"
)

// CatalogID is the Trackus catalog identifier, by convention the extension ID.
const CatalogID = "trackus"

// Standard tracker IDs a conversation can reach. The upstream facade
// materialises these on first entry, which is why recording against one needs
// no setup step. They are listed in the action description rather than as a
// closed Enum, because a tracker created through trackers.create has a
// generated ID that must also be recordable.
const (
	PushUpsTrackerID      = "_push_ups"
	PullUpsTrackerID      = "_pull_ups"
	SquatsTrackerID       = "_squats"
	JumpingJacksTrackerID = "_jumping_jacks"
	WeightTrackerID       = "_weight"
	MileageTrackerID      = "_mileage"
	SpendingTrackerID     = "_spending"
)

var trackerIDArg = convospec.ArgDef{
	Name: "trackerID", Type: convospec.ArgTypeString, Required: true,
	Description: `Target tracker ID. Standard IDs: "_push_ups", "_pull_ups", "_squats", ` +
		`"_jumping_jacks", "_weight", "_mileage", "_spending". Use the exact ID of an ` +
		`existing custom tracker when the measurement is not one of these. If no ` +
		`tracker fits, call trackers.create instead of inventing an ID.`,
}

var addEntryDef = convospec.ActionDef{
	ID: "trackers.add_entry", Extension: CatalogID,
	Summary:     "Record a tracker measurement",
	Description: "Records one measurement against a tracker, e.g. a number of press-ups, a body weight, a distance.",
	Args: []convospec.ArgDef{
		trackerIDArg,
		{Name: "value", Type: convospec.ArgTypeFloat, Required: true, Description: "The measured number"},
		{Name: "date", Type: convospec.ArgTypeString, Description: "ISO date (YYYY-MM-DD) for a backdated entry; omit for now"},
	},
	Examples: []convospec.Example{
		{UserText: "push-ups 20", Args: map[string]any{"trackerID": PushUpsTrackerID, "value": 20.0}},
		{UserText: "20 push-ups", Args: map[string]any{"trackerID": PushUpsTrackerID, "value": 20.0}},
		{UserText: "10 pull-ups", Args: map[string]any{"trackerID": PullUpsTrackerID, "value": 10.0}},
		{UserText: "my weight is 80.5", Args: map[string]any{"trackerID": WeightTrackerID, "value": 80.5}},
	},
	Result: "tracker title, unit and running total",
}

var createTrackerDef = convospec.ActionDef{
	ID: "trackers.create", Extension: CatalogID, Confirm: true,
	Summary: "Create a new tracker and record the first measurement",
	Description: "Use ONLY when the measurement does not fit any existing tracker — for example " +
		"running distance, which has no standard tracker. Requires confirmation, so a " +
		"misheard message never permanently adds a tracker.",
	Args: []convospec.ArgDef{
		{Name: "title", Type: convospec.ArgTypeString, Required: true, Description: `Human name for the tracker, e.g. "Running"`},
		{Name: "unit", Type: convospec.ArgTypeString, Description: `Unit of measure, e.g. "km"`},
		{Name: "value", Type: convospec.ArgTypeFloat, Required: true, Description: "The first measurement to record"},
		{Name: "numberKind", Type: convospec.ArgTypeString, Enum: []string{"summable", "absolute"},
			Description: `"summable" when measurements add up over time (distance, repetitions); "absolute" when each replaces the last (body weight). Defaults to summable.`},
		{Name: "trackerID", Type: convospec.ArgTypeString, Description: "Filled in by resolution; do not set"},
	},
	Examples: []convospec.Example{
		{UserText: "ran 10 km", Args: map[string]any{"title": "Running", "unit": "km", "value": 10.0, "numberKind": "summable"}},
	},
	Result: "created tracker ID and the recorded measurement",
}

var viewDef = convospec.ActionDef{
	ID: "trackers.view", Extension: CatalogID,
	Summary:  "Show recent measurements of a tracker",
	Args:     []convospec.ArgDef{trackerIDArg},
	Examples: []convospec.Example{{UserText: "how many push-ups have I done?", Args: map[string]any{"trackerID": PushUpsTrackerID}}},
	Result:   "tracker title and recent entries",
}

// Declaration is the Trackus conversational capability as plain data: what it
// can do, how to interpret language into it, and what vocabulary it is
// listening for.
func Declaration() convospec.Catalog {
	return convospec.Catalog{
		ID:    CatalogID,
		Title: "Trackers (exercise, weight, distance, spending)",
		PromptHints: []string{
			`"push-ups 20", "20 push-ups" and "did 20 pushups" all mean trackers.add_entry on "_push_ups" with value 20 — word order carries no meaning.`,
			`A bare exercise name with a number is a measurement, never a shopping-list item.`,
			`Only use trackers.create when no standard tracker fits (e.g. running distance); prefer an existing tracker.`,
			`Never invent a tracker ID. If unsure which tracker is meant, call clarify.`,
		},
		// Distinctive vocabulary only. These narrow the action set when they
		// match exactly one catalog, and a false positive is the only way this
		// can misroute — so a generic word like "did" or "add" must NOT be here.
		Triggers: []string{
			"push-ups", "pushups", "push ups",
			"pull-ups", "pullups", "pull ups",
			"squats", "jumping jacks",
			"km", "kg", "lbs", "miles",
			"ran", "jogged", "walked", "weighed",
		},
		Actions: []convospec.ActionDef{addEntryDef, createTrackerDef, viewDef},
	}
}

// Rules is the extension's own deterministic interpretation, used by the mock
// LLM client so tests never depend on a cloud model. Both word orders are
// declared explicitly rather than handled by one clever pattern, because the
// acceptance criteria are written per phrasing and a failure should name the
// phrasing that broke.
func Rules() []convospec.Rule {
	entry := func(pattern, trackerID string, priority int) convospec.Rule {
		return convospec.Rule{
			Pattern:  regexp.MustCompile(pattern),
			ActionID: addEntryDef.ID,
			Args:     map[string]any{"trackerID": trackerID, "value": "$1"},
			Priority: priority,
		}
	}
	reversed := func(pattern, trackerID string, priority int) convospec.Rule {
		return convospec.Rule{
			Pattern:  regexp.MustCompile(pattern),
			ActionID: addEntryDef.ID,
			Args:     map[string]any{"trackerID": trackerID, "value": "$1"},
			Priority: priority,
		}
	}
	return []convospec.Rule{
		// "ran 10 km" must beat any generic distance handling, hence priority 20.
		{
			Pattern:  regexp.MustCompile(`^(?:ran|jogged) ([\d.]+) (?:km|kilometers|kilometres)$`),
			ActionID: createTrackerDef.ID,
			Args: map[string]any{
				"title": "Running", "unit": "km", "value": "$1", "numberKind": "summable",
			},
			Priority: 20,
		},
		entry(`^push-?\s?ups? ([\d.]+)$`, PushUpsTrackerID, 10),
		reversed(`^([\d.]+) push-?\s?ups?$`, PushUpsTrackerID, 10),
		entry(`^pull-?\s?ups? ([\d.]+)$`, PullUpsTrackerID, 10),
		reversed(`^([\d.]+) pull-?\s?ups?$`, PullUpsTrackerID, 10),
		entry(`^squats ([\d.]+)$`, SquatsTrackerID, 10),
		reversed(`^([\d.]+) squats$`, SquatsTrackerID, 10),
		{
			Pattern:  regexp.MustCompile(`^(?:my )?weigh(?:t|ed) (?:is )?([\d.]+)(?: kg)?$`),
			ActionID: addEntryDef.ID,
			Args:     map[string]any{"trackerID": WeightTrackerID, "value": "$1"},
			Priority: 10,
		},
		{
			Pattern:  regexp.MustCompile(`^how many push-?\s?ups?.*$`),
			ActionID: viewDef.ID,
			Args:     map[string]any{"trackerID": PushUpsTrackerID},
			Priority: 30,
		},
	}
}
