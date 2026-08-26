package convo4listus

import (
	"testing"
	"time"

	"github.com/sneat-co/sneat-go-core/convospec"
)

func TestDeclaredExamplesValidate(t *testing.T) {
	for _, action := range Declaration().Actions {
		for i, example := range action.Examples {
			if _, err := action.ValidateArgs(example.Args); err != nil {
				t.Errorf("%s example %d (%q): %v", action.ID, i, example.UserText, err)
			}
		}
	}
}

func TestRulesAgreeWithDeclaration(t *testing.T) {
	if err := Declaration().ValidateRules(Rules()); err != nil {
		t.Fatal(err)
	}
}

// firstMatch mimics the mock engine: highest priority first, declaration order
// breaking ties.
func firstMatch(t *testing.T, text string) (string, map[string]any) {
	t.Helper()
	declaration := Declaration()
	normalized := convospec.NormalizeText(text)

	// Mirror the engine: imperative RuleFuncs run before declarative Rules,
	// because they exist for the cases a Rule cannot express and must not be
	// pre-empted by one.
	for _, ruleFunc := range RuleFuncs() {
		if calls, matched := ruleFunc(normalized, time.Time{}); matched && len(calls) > 0 {
			return calls[0].ActionID, calls[0].Args
		}
	}

	best := -1 << 31
	var bestID string
	var bestArgs map[string]any
	for _, rule := range Rules() {
		def, ok := declaration.Action(rule.ActionID)
		if !ok {
			t.Fatalf("rule names undeclared action %s", rule.ActionID)
		}
		args, matched := rule.Match(normalized, def)
		if !matched || rule.Priority <= best {
			continue
		}
		best, bestID, bestArgs = rule.Priority, rule.ActionID, args
	}
	return bestID, bestArgs
}

// The headline pilot case: the amount is data, and the title stays the bare noun
// so that "bought milk" can still match it.
func TestQuantifiedItemKeepsBareTitle(t *testing.T) {
	actionID, args := firstMatch(t, "milk 2 liters")
	if actionID != addItemsDef.ID {
		t.Fatalf(`"milk 2 liters" -> %s, want %s`, actionID, addItemsDef.ID)
	}
	items, ok := args["items"].([]string)
	if !ok || len(items) != 1 || items[0] != "milk" {
		t.Errorf("items = %#v, want []string{\"milk\"} — the amount must not be in the title", args["items"])
	}
	if args["quantity"] != 2.0 {
		t.Errorf("quantity = %#v, want 2.0", args["quantity"])
	}
	if args["unit"] != "L" {
		t.Errorf("unit = %#v, want canonical \"L\"", args["unit"])
	}
}

// "bought milk" must outrank the add rules, or it would create an item called
// "bought milk" instead of completing the existing one.
func TestBoughtResolvesToSetDone(t *testing.T) {
	actionID, args := firstMatch(t, "bought milk")
	if actionID != setDoneDef.ID {
		t.Fatalf(`"bought milk" -> %s, want %s`, actionID, setDoneDef.ID)
	}
	titles, _ := args["itemTitles"].([]string)
	if len(titles) != 1 || titles[0] != "milk" {
		t.Errorf("itemTitles = %#v, want [milk]", args["itemTitles"])
	}
	if args["isDone"] != true {
		t.Errorf("isDone = %#v, want true", args["isDone"])
	}
}

func TestRulePriorities(t *testing.T) {
	for _, tt := range []struct{ text, wantAction string }{
		{"buy milk", addItemsDef.ID},
		{"milk 2 kg", addItemsDef.ID},
		{"bought milk", setDoneDef.ID},
		{"clear the shopping list", removeAllDef.ID},
		{"remove bought items from the shopping list", removeDoneDef.ID},
		{"what's on my shopping list?", viewDef.ID},
	} {
		if got, _ := firstMatch(t, tt.text); got != tt.wantAction {
			t.Errorf("%q -> %s, want %s", tt.text, got, tt.wantAction)
		}
	}
}

// The list-name word must resolve to the RIGHT listID across all four
// standard lists, not just the groceries default — this is the legacy
// grammar's list-name capture that must not get lost when its rules moved
// here from the shared mock package.
func TestListNameResolvesAcrossLists(t *testing.T) {
	for _, tt := range []struct{ text, wantAction, wantListID string }{
		{"what's on my movies list?", viewDef.ID, MoviesListID},
		{"what's on my books list?", viewDef.ID, BooksListID},
		{"what's on my to-do list?", viewDef.ID, TasksListID},
		{"clear the movies list", removeAllDef.ID, MoviesListID},
		{"remove done items from the books list", removeDoneDef.ID, BooksListID},
	} {
		actionID, args := firstMatch(t, tt.text)
		if actionID != tt.wantAction {
			t.Fatalf("%q -> %s, want %s", tt.text, actionID, tt.wantAction)
		}
		if got := args["listID"]; got != tt.wantListID {
			t.Errorf("%q listID = %#v, want %q", tt.text, got, tt.wantListID)
		}
	}
}

// "remove milk from the shopping list" must resolve to remove_items with the
// item title AND the named list, not just remove_done/remove_all — the plain
// item-removal rule that the earlier, weaker convo4listus rule set was
// missing entirely.
func TestRemoveNamedItemFromList(t *testing.T) {
	actionID, args := firstMatch(t, "remove milk from the shopping list")
	if actionID != removeItemsDef.ID {
		t.Fatalf(`"remove milk from the shopping list" -> %s, want %s`, actionID, removeItemsDef.ID)
	}
	titles, _ := args["itemTitles"].([]string)
	if len(titles) != 1 || titles[0] != "milk" {
		t.Errorf("itemTitles = %#v, want [milk]", args["itemTitles"])
	}
	if got := args["listID"]; got != GroceriesListID {
		t.Errorf("listID = %#v, want %q", got, GroceriesListID)
	}
}

// A bulk-removal phrase must never be misread as removing an item literally
// titled "done items" — ruleListsRemoveItems declines it unconditionally.
func TestRemoveDoneItemsIsNotAnItemTitle(t *testing.T) {
	actionID, _ := firstMatch(t, "remove done items from the list")
	if actionID != removeDoneDef.ID {
		t.Fatalf(`"remove done items from the list" -> %s, want %s`, actionID, removeDoneDef.ID)
	}
}

// "bought milk and eggs" must produce TWO item titles, not one item literally
// titled "milk and eggs" — a declarative Rule's ArgTypeStringSlice wraps a
// capture as a single-element slice, so this needs the RuleFunc's split.
func TestSetDoneSplitsMultipleItems(t *testing.T) {
	actionID, args := firstMatch(t, "bought milk and eggs")
	if actionID != setDoneDef.ID {
		t.Fatalf(`"bought milk and eggs" -> %s, want %s`, actionID, setDoneDef.ID)
	}
	titles, _ := args["itemTitles"].([]string)
	if len(titles) != 2 || titles[0] != "milk" || titles[1] != "eggs" {
		t.Errorf("itemTitles = %#v, want [milk eggs]", args["itemTitles"])
	}
}

// ruleBareItemInListScope: a bare noun means "add to groceries" ONLY when
// Listus is the sole extension in scope this turn.
func TestScopedBareItem(t *testing.T) {
	scoped := ScopedRuleFuncs()
	if len(scoped) != 1 {
		t.Fatalf("ScopedRuleFuncs() = %d funcs, want 1", len(scoped))
	}
	rule := scoped[0]

	listusOnly := []string{addItemsDef.ID, setDoneDef.ID, removeItemsDef.ID, removeDoneDef.ID, removeAllDef.ID, viewDef.ID, "clarify"}
	calls, matched := rule(convospec.NormalizeText("milk"), time.Time{}, listusOnly)
	if !matched || len(calls) != 1 || calls[0].ActionID != addItemsDef.ID {
		t.Fatalf(`"milk" in listus-only scope -> matched=%v calls=%#v, want one lists.add_items call`, matched, calls)
	}
	items, _ := calls[0].Args["items"].([]string)
	if len(items) != 1 || items[0] != "milk" {
		t.Errorf("items = %#v, want [milk]", calls[0].Args["items"])
	}

	fullScope := []string{addItemsDef.ID, "contacts.create", "clarify"}
	if _, matched := rule(convospec.NormalizeText("milk"), time.Time{}, fullScope); matched {
		t.Error(`"milk" must not match when another extension's action is also in scope`)
	}
}

// A destructive action must be staged for confirmation, never executed straight.
func TestDestructiveActionsRequireConfirmation(t *testing.T) {
	declaration := Declaration()
	for _, id := range []string{removeItemsDef.ID, removeDoneDef.ID, removeAllDef.ID} {
		def, ok := declaration.Action(id)
		if !ok {
			t.Fatalf("%s not declared", id)
		}
		if !def.Confirm {
			t.Errorf("%s must require confirmation", id)
		}
	}
	for _, id := range []string{addItemsDef.ID, setDoneDef.ID, viewDef.ID} {
		def, _ := declaration.Action(id)
		if def.Confirm {
			t.Errorf("%s should not require confirmation", id)
		}
	}
}

// An out-of-vocabulary unit must be rejected by ordinary argument validation
// rather than stored as free text.
func TestUnitVocabularyIsClosed(t *testing.T) {
	if _, err := addItemsDef.ValidateArgs(map[string]any{
		"items": []string{"milk"}, "quantity": 2.0, "unit": "liters",
	}); err == nil {
		t.Error(`unit "liters" is outside the vocabulary and must be rejected`)
	}
	if _, err := addItemsDef.ValidateArgs(map[string]any{
		"items": []string{"milk"}, "quantity": 2.0, "unit": "L",
	}); err != nil {
		t.Errorf(`unit "L" must be accepted: %v`, err)
	}
}

// Triggers must not claim a Trackus message, or every measurement would narrow
// to Listus.
func TestTriggersDoNotClaimTrackusMessages(t *testing.T) {
	declaration := Declaration()
	for _, text := range []string{"20 push-ups", "push-ups 20", "10 pull-ups", "my weight is 80.5"} {
		if declaration.MatchesTriggers(convospec.NormalizeText(text)) {
			t.Errorf("%q must not be claimed by Listus triggers", text)
		}
	}
	for _, text := range []string{"buy milk", "bought milk", "clear the shopping list"} {
		if !declaration.MatchesTriggers(convospec.NormalizeText(text)) {
			t.Errorf("%q should be claimed by Listus triggers", text)
		}
	}
}

// "milk 2 liters" must be claimed by NEITHER pilot catalog, so the prefilter
// fails open and the model sees the full action set — this is the acceptance
// criterion prefilter-fails-open-on-no-match depends on.
func TestBareQuantifiedItemIsUnclaimed(t *testing.T) {
	if Declaration().MatchesTriggers(convospec.NormalizeText("milk 2 liters")) {
		t.Error(`"milk 2 liters" must not be claimed by a trigger; the prefilter is meant to fail open here`)
	}
}
