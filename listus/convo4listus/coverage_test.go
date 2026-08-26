package convo4listus

import (
	"testing"
	"time"

	"github.com/sneat-co/sneat-go-core/convospec"
)

// This file covers the branches the behavioural tests do not reach, so the
// contract module keeps its 100% coverage floor. Each case is a real
// degenerate input a message could produce, not a synthetic line-hitter:
// an empty item tail, an unrecognised list name, or a scope that offers no
// list actions at all.

func runRuleFuncs(text string) ([]convospec.ActionCall, bool) {
	for _, ruleFunc := range RuleFuncs() {
		if calls, matched := ruleFunc(convospec.NormalizeText(text), time.Time{}); matched {
			return calls, true
		}
	}
	return nil, false
}

// A phrase whose item tail collapses to nothing must not produce an action with
// an empty item list — the list would gain a blank entry nobody can act on.
func TestRuleFuncsRejectEmptyItemTails(t *testing.T) {
	for _, text := range []string{
		"buy ,",
		"buy and",
		"add , to the shopping list",
		"add and to the movies list",
		"remove and from the shopping list",
		"bought and",
		"mark and as done",
	} {
		if calls, matched := runRuleFuncs(text); matched {
			t.Errorf("%q produced %+v; an empty item tail must not match", text, calls)
		}
	}
}

// An unrecognised list name must fall back to the default list rather than
// emitting a listID the domain would reject.
func TestUnknownListNameFallsBackToDefault(t *testing.T) {
	for _, text := range []string{
		"show the list",
		"clear the list",
		"remove done items from the list",
		"remove all from the list",
	} {
		calls, matched := runRuleFuncs(text)
		if !matched {
			t.Fatalf("%q matched no rule", text)
		}
		if listID, ok := calls[0].Args["listID"]; ok && listID == "" {
			t.Errorf("%q emitted an empty listID", text)
		}
	}
}

// Every recognised list-name word must resolve to a declared standard list ID.
func TestListIDFromWordCoversEveryName(t *testing.T) {
	for word, want := range map[string]string{
		"groceries": GroceriesListID,
		"grocery":   GroceriesListID,
		"shopping":  GroceriesListID,
		"todo":      TasksListID,
		"to-do":     TasksListID,
		"task":      TasksListID,
		"tasks":     TasksListID,
		"watch":     MoviesListID,
		"movies":    MoviesListID,
		"read":      BooksListID,
		"books":     BooksListID,
	} {
		if got := listIDFromWord(word); got != want {
			t.Errorf("listIDFromWord(%q) = %q, want %q", word, got, want)
		}
	}
	if got := listIDFromWord("nonsense"); got != "" {
		t.Errorf(`listIDFromWord("nonsense") = %q, want ""`, got)
	}
}

// A named list must be honoured, not silently replaced by groceries.
func TestNamedListIsHonoured(t *testing.T) {
	calls, matched := runRuleFuncs("add inception to the movies list")
	if !matched {
		t.Fatal(`"add inception to the movies list" matched no rule`)
	}
	if got := calls[0].Args["listID"]; got != MoviesListID {
		t.Errorf("listID = %v, want %s", got, MoviesListID)
	}
}

func TestBareItemScopeGuards(t *testing.T) {
	scoped := ScopedRuleFuncs()
	if len(scoped) != 1 {
		t.Fatalf("want exactly one scoped rule, got %d", len(scoped))
	}
	bareItem := scoped[0]
	listusScope := []string{"lists.add_items", "lists.view", "clarify"}

	t.Run("fires when Listus is alone in scope", func(t *testing.T) {
		calls, matched := bareItem(convospec.NormalizeText("milk"), time.Time{}, listusScope)
		if !matched {
			t.Fatal(`"milk" should be an add when only Listus is listening`)
		}
		if calls[0].ActionID != addItemsDef.ID {
			t.Errorf("action = %s, want %s", calls[0].ActionID, addItemsDef.ID)
		}
	})

	t.Run("does not fire when another extension is in scope", func(t *testing.T) {
		mixed := append([]string{"trackers.add_entry"}, listusScope...)
		if _, matched := bareItem(convospec.NormalizeText("milk"), time.Time{}, mixed); matched {
			t.Error(`"milk" is ambiguous when another extension is listening and must not be claimed`)
		}
	})

	t.Run("does not fire without the add action", func(t *testing.T) {
		viewOnly := []string{"lists.view", "clarify"}
		if _, matched := bareItem(convospec.NormalizeText("milk"), time.Time{}, viewOnly); matched {
			t.Error("must not add when the add action is not on offer")
		}
	})

	t.Run("rejects empty and over-long phrases", func(t *testing.T) {
		for _, text := range []string{"", "one two three four five"} {
			if _, matched := bareItem(convospec.NormalizeText(text), time.Time{}, listusScope); matched {
				t.Errorf("%q must not be read as a bare item", text)
			}
		}
	})

	t.Run("rejects command-like openers", func(t *testing.T) {
		for _, opener := range []string{"list", "show", "view", "what", "who", "how", "delete", "remove", "cancel", "add", "help", "start", "buy", "mark"} {
			if _, matched := bareItem(convospec.NormalizeText(opener+" things"), time.Time{}, listusScope); matched {
				t.Errorf("%q opens a command, not a list item", opener)
			}
		}
	})
}

// An empty scope must not be treated as "Listus is the only one listening" —
// OnlyCatalogInScope requires at least one of the extension's own actions.
func TestOnlyCatalogInScopeRequiresOwnAction(t *testing.T) {
	if convospec.OnlyCatalogInScope(nil, "lists.") {
		t.Error("an empty scope must not count as exclusive")
	}
	if convospec.OnlyCatalogInScope([]string{"clarify"}, "lists.") {
		t.Error("clarify alone must not count as exclusive")
	}
	if !convospec.OnlyCatalogInScope([]string{"lists.view", "clarify"}, "lists.") {
		t.Error("own actions plus clarify is exclusive")
	}
}

// Declarative rules must decline text they do not match, and every declared
// rule must resolve against the declaration.
func TestDeclarativeRulesDeclineForeignText(t *testing.T) {
	declaration := Declaration()
	for _, rule := range Rules() {
		def, ok := declaration.Action(rule.ActionID)
		if !ok {
			t.Fatalf("rule names undeclared action %s", rule.ActionID)
		}
		if _, matched := rule.Match(convospec.NormalizeText("20 push-ups"), def); matched {
			t.Errorf("rule for %s must not claim a tracker measurement", rule.ActionID)
		}
	}
}

// ruleListsRemoveItems must decline a bulk-removal phrase UNCONDITIONALLY, not
// merely defer to the bulk rules. Its own pattern is broad enough to capture
// "done items" as if that were an item title, so if the bulk actions ever drop
// out of scope the phrase must still reach clarify rather than misfire as a
// one-off removal. Reached by calling it directly, because the rule loop stops
// at the first match and a bulk rule always claims these first.
func TestRemoveItemsDeclinesBulkPhrasesUnconditionally(t *testing.T) {
	for _, text := range []string{
		"remove done items from the shopping list",
		"remove all from the shopping list",
		"clear the shopping list",
	} {
		if calls, matched := ruleListsRemoveItems(convospec.NormalizeText(text), time.Time{}); matched {
			t.Errorf("%q must not be read as a one-off item removal, got %+v", text, calls)
		}
	}
}

// "mark milk as done" is the second phrasing of set_done; the rule loop reaches
// it only when the "bought" form does not match.
func TestMarkAsDonePhrasing(t *testing.T) {
	calls, matched := ruleListsSetDone(convospec.NormalizeText("mark milk as done"), time.Time{})
	if !matched {
		t.Fatal(`"mark milk as done" matched no set_done phrasing`)
	}
	if calls[0].ActionID != setDoneDef.ID {
		t.Errorf("action = %s, want %s", calls[0].ActionID, setDoneDef.ID)
	}
	titles, _ := calls[0].Args["itemTitles"].([]string)
	if len(titles) != 1 || titles[0] != "milk" {
		t.Errorf("itemTitles = %#v, want [milk]", calls[0].Args["itemTitles"])
	}
	// And it must split, like the "bought" form.
	multi, matched := ruleListsSetDone(convospec.NormalizeText("mark milk and eggs as done"), time.Time{})
	if !matched {
		t.Fatal("multi-item mark-as-done matched no rule")
	}
	if titles, _ := multi[0].Args["itemTitles"].([]string); len(titles) != 2 {
		t.Errorf("itemTitles = %#v, want two items", multi[0].Args["itemTitles"])
	}
}

// A phrase whose only word is a joiner collapses to no items and must decline
// rather than add a blank entry.
func TestBareItemJoinerOnlyDeclines(t *testing.T) {
	bareItem := ScopedRuleFuncs()[0]
	if _, matched := bareItem(convospec.NormalizeText("and"), time.Time{}, []string{"lists.add_items", "clarify"}); matched {
		t.Error(`"and" collapses to no items and must not be added`)
	}
}
