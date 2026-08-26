package convo4listus

import (
	"regexp"
	"strings"
	"time"

	"github.com/sneat-co/sneat-go-core/convospec"
)

// CatalogID is the Listus catalog identifier, by convention the extension ID.
const CatalogID = "listus"

// UnitVocabulary is the closed set of units a quantified item may carry.
//
// It is closed deliberately. An open unit field made a quantified turn
// non-deterministic — "2 liters" could arrive as "liters", "litres", "l" or "L"
// — so the interpreter must emit one of these and anything else is rejected by
// ordinary argument validation. Nothing rewrites a value: this is a vocabulary,
// not a normalisation step, and there is no conversion between members.
// Widening the set is a one-line change here, deliberately reviewable.
var UnitVocabulary = []string{"pcs", "g", "kg", "ml", "L", "pack", "box", "bottle", "can"}

var listIDArg = convospec.ArgDef{
	Name: "listID", Type: convospec.ArgTypeString,
	Description: `Target list ID in "{type}!{sub}" form: "buy!groceries" (groceries/shopping), ` +
		`"do!tasks" (to-do), "watch!movies", "read!books". Defaults to "buy!groceries".`,
}

var addItemsDef = convospec.ActionDef{
	ID: "lists.add_items", Extension: CatalogID,
	Summary: "Add items to a list",
	Description: "Adds one or more items to a list (shopping, to-do, movies, books). " +
		"The list is created automatically on first use.",
	Args: []convospec.ArgDef{
		{Name: "items", Type: convospec.ArgTypeStringSlice, Required: true,
			Description: `Item titles as BARE nouns — "milk", not "2 liters of milk". Put the amount in quantity/unit.`},
		{Name: "quantity", Type: convospec.ArgTypeFloat,
			Description: "How much of the item is needed. Only valid when exactly ONE item is given; " +
				"for several quantified items, emit one lists.add_items call per item."},
		{Name: "unit", Type: convospec.ArgTypeString, Enum: UnitVocabulary,
			Description: `Unit for quantity. Must be one of the listed values — use "L" for liters/litres, "kg" for kilos, "pcs" for a plain count.`},
		listIDArg,
	},
	Examples: []convospec.Example{
		{UserText: "buy milk", Args: map[string]any{"items": []string{"milk"}}},
		{UserText: "milk 2 liters", Args: map[string]any{"items": []string{"milk"}, "quantity": 2.0, "unit": "L"}},
		{UserText: "add bread and eggs to the shopping list",
			Args: map[string]any{"items": []string{"bread", "eggs"}, "listID": GroceriesListID}},
	},
	Result: "titles of created items",
}

var setDoneDef = convospec.ActionDef{
	ID: "lists.set_done", Extension: CatalogID,
	Summary: "Mark list items as done / not done",
	Args: []convospec.ArgDef{
		{Name: "itemTitles", Type: convospec.ArgTypeStringSlice, Required: true, Description: "Titles of the items to mark"},
		{Name: "isDone", Type: convospec.ArgTypeBool, Description: "Defaults to true"},
		{Name: "itemIDs", Type: convospec.ArgTypeStringSlice, Description: "Filled in by resolution; do not set"},
		listIDArg,
	},
	Examples: []convospec.Example{{UserText: "bought milk", Args: map[string]any{"itemTitles": []string{"milk"}}}},
	Result:   "titles of changed items",
}

var removeItemsDef = convospec.ActionDef{
	ID: "lists.remove_items", Extension: CatalogID, Confirm: true,
	Summary: "Remove items from a list",
	Args: []convospec.ArgDef{
		{Name: "itemTitles", Type: convospec.ArgTypeStringSlice, Required: true, Description: "Titles of the items to remove"},
		{Name: "itemIDs", Type: convospec.ArgTypeStringSlice, Description: "Filled in by resolution; do not set"},
		listIDArg,
	},
	Examples: []convospec.Example{{UserText: "remove milk from the shopping list", Args: map[string]any{"itemTitles": []string{"milk"}}}},
	Result:   "titles of removed items",
}

var removeDoneDef = convospec.ActionDef{
	ID: "lists.remove_done", Extension: CatalogID, Confirm: true,
	Summary: "Remove done items from a list",
	Args: []convospec.ArgDef{
		{Name: "itemTitles", Type: convospec.ArgTypeStringSlice, Description: "Filled in by resolution; do not set"},
		{Name: "itemIDs", Type: convospec.ArgTypeStringSlice, Description: "Filled in by resolution; do not set"},
		listIDArg,
	},
	Examples: []convospec.Example{{UserText: "remove bought items from the shopping list", Args: map[string]any{"listID": GroceriesListID}}},
	Result:   "titles of removed done items",
}

var removeAllDef = convospec.ActionDef{
	ID: "lists.remove_all", Extension: CatalogID, Confirm: true,
	Summary: "Remove all items from a list",
	Args: []convospec.ArgDef{
		{Name: "itemTitles", Type: convospec.ArgTypeStringSlice, Description: "Filled in by resolution; do not set"},
		{Name: "itemIDs", Type: convospec.ArgTypeStringSlice, Description: "Filled in by resolution; do not set"},
		listIDArg,
	},
	Examples: []convospec.Example{{UserText: "clear the shopping list", Args: map[string]any{"listID": GroceriesListID}}},
	Result:   "titles of all removed items",
}

var viewDef = convospec.ActionDef{
	ID: "lists.view", Extension: CatalogID,
	Summary:  "Show the items of a list",
	Args:     []convospec.ArgDef{listIDArg},
	Examples: []convospec.Example{{UserText: "what's on my shopping list?", Args: map[string]any{"listID": GroceriesListID}}},
	Result:   "list title and items with done-state",
}

// Declaration is the Listus conversational capability as plain data.
func Declaration() convospec.Catalog {
	return convospec.Catalog{
		ID:    CatalogID,
		Title: "Lists (shopping, to-do, movies to watch, books to read)",
		PromptHints: []string{
			`A bare item name like "milk" most likely means lists.add_items to the groceries list.`,
			`"milk 2 liters" is ONE item titled "milk" with quantity 2 and unit "L" — never fold the amount into the title.`,
			`For several quantified items, emit one lists.add_items call per item; quantity applies to a single item only.`,
			`"bought X" or "X done" means lists.set_done.`,
			`"remove done"/"remove bought" is lists.remove_done; "remove all"/"clear the list" is lists.remove_all. Both need confirmation.`,
			`A number with an exercise name ("20 push-ups") is NOT a list item.`,
		},
		// Distinctive vocabulary only. A false positive on exactly one catalog
		// is the one way the routing prefilter can misroute, so generic verbs
		// like "add" and "get" are deliberately absent — a bare item name is
		// meant to fall through to the full action set.
		Triggers: []string{
			"buy", "bought", "purchased",
			"shopping list", "groceries", "grocery",
			"to-do", "todo", "task list",
			"movies to watch", "books to read",
			"clear the list", "remove done",
		},
		Actions: []convospec.ActionDef{addItemsDef, setDoneDef, removeItemsDef, removeDoneDef, removeAllDef, viewDef},
	}
}

// Rules is the extension's own deterministic interpretation, used by the mock
// LLM client so tests never depend on a cloud model.
//
// This is deliberately the SMALL half of the grammar. Everything here needs
// nothing beyond $n-substitution into one target action, so a declarative Rule
// expresses it exactly and gets validated against the action's ArgDefs for
// free. Every other interpretation — anything that must resolve a list-name
// word ("movies"/"books"/"to-do"/…) to a listID, or split "milk and bread"
// into several item titles — needs a lookup or a loop that Rule's $n-only
// templating cannot express, so it lives in RuleFuncs instead.
//
// Priorities matter here: "bought milk" must resolve to set_done, not to adding
// an item called "bought milk", so the set_done patterns outrank the catch-all
// bare-item rule.
func Rules() []convospec.Rule {
	return []convospec.Rule{
		{ // "milk 2 liters" — quantified single item
			Pattern:  regexp.MustCompile(`^([a-z ]+?) ,? ?([\d.]+) (?:liters?|litres?|l)$`),
			ActionID: addItemsDef.ID,
			Args:     map[string]any{"items": "$1", "quantity": "$2", "unit": "L", "listID": GroceriesListID},
			Priority: 25,
		},
		{ // "milk 2 kg"
			Pattern:  regexp.MustCompile(`^([a-z ]+?) ,? ?([\d.]+) (?:kg|kilos?|kilograms?)$`),
			ActionID: addItemsDef.ID,
			Args:     map[string]any{"items": "$1", "quantity": "$2", "unit": "kg", "listID": GroceriesListID},
			Priority: 25,
		},
		// Everything else — view, remove_done, remove_all, remove_items,
		// set_done and the plain "buy X" add — needs the list-name→listID
		// lookup (reListName/listIDFromWord below) and/or multi-item
		// splitting, so it lives in RuleFuncs.
	}
}

// reListName is the closed vocabulary of list-name WORDS a user's phrase may
// use, as distinct from the listID it resolves to (see listIDFromWord). It
// covers all four standard lists so "clear the movies list" and "clear the
// shopping list" both resolve, not just the groceries default.
var reListName = `(?:groceries|grocery|shopping|to-?do|tasks?|watch|movies|read|books)`

// listIDFromWord maps a captured list-name word to its listID. It is a lookup
// table, not a substitution — the reason every rule that names a list is a
// RuleFunc rather than a declarative Rule, whose $n templating can only copy a
// capture verbatim or type it, never look it up.
func listIDFromWord(word string) string {
	switch word {
	case "groceries", "grocery", "shopping":
		return GroceriesListID
	case "todo", "to-do", "task", "tasks", "do":
		return TasksListID
	case "watch", "movies":
		return MoviesListID
	case "read", "books":
		return BooksListID
	default:
		return ""
	}
}

var (
	// Bare "add X" is deliberately NOT claimed. It is ambiguous across
	// extensions — "add contact Jane", "add asset MacBook" — so it must fall
	// through to the full action set and let the model decide. Listus claims
	// only an explicit buy/need, or an add that names the target list.
	reBuyItems  = regexp.MustCompile(`^(?:buy|need|we need|need to buy) (.+)$`)
	reAddToList = regexp.MustCompile(`^add (.+?) to (?:the |my )?(` + reListName + `)(?: list)?$`)

	// reBought/reMarkDone: "bought milk" / "mark milk as done" / "milk is done".
	reBought   = regexp.MustCompile(`^(?:bought|got|purchased) (.+)$`)
	reMarkDone = regexp.MustCompile(`^(?:mark )?(.+?) (?:is |as )?done$`)

	// Bulk removal: "remove done", "remove all", "clear the <list> list". These
	// three patterns are also used, unconditionally, to keep the plain
	// item-removal rule from misreading one of these phrases as an item
	// literally titled "done items" or "all" — see ruleListsRemoveItems.
	reRemoveDone = regexp.MustCompile(`^(?:remove|delete) (?:the )?(?:done|bought|completed)(?: items?)?(?: from (?:the |my )?(?:(` + reListName + `) )?list)?$`)
	reRemoveAll  = regexp.MustCompile(`^(?:remove|delete) (?:all|everything)(?: items?)?(?: from (?:the |my )?(?:(` + reListName + `) )?list)?$`)
	reClearList  = regexp.MustCompile(`^clear (?:the |my )?(?:(` + reListName + `) )?list$`)

	// Single/selected-item removal: "remove milk from the shopping list".
	reRemoveItem = regexp.MustCompile(`^(?:remove|delete) (.+?) from (?:the |my )?(?:(` + reListName + `) )?list$`)

	// "show/view/what's on [the <list>] list".
	reViewList = regexp.MustCompile(`^(?:show|view|whats on|what is on) (?:the |my )?(?:(` + reListName + `) )?list$`)

	reItemJoiner = regexp.MustCompile(`\s*(?:,|;|\band\b)\s*`)
)

// splitItems breaks a captured phrase into individual titles. "milk and bread"
// is two items: folding them into one title creates an entry nobody can tick
// off, and "bought milk" would then never match it.
func splitItems(phrase string) []string {
	items := make([]string, 0, 2)
	for _, field := range reItemJoiner.Split(phrase, -1) {
		if field = strings.TrimSpace(field); field != "" {
			items = append(items, field)
		}
	}
	return items
}

// RuleFuncs covers every interpretation that needs the list-name lookup
// (listIDFromWord) or multi-item splitting (splitItems) — neither of which a
// declarative Rule's $n-substitution-only templating can express.
//
// Order matters and mirrors the pre-migration central grammar exactly: the
// bulk-removal rules run before the plain item-removal rule so
// "remove done items from the list" cannot be misread as removing an item
// literally titled "done items" — see ruleListsRemoveItems's own unconditional
// exclusion check, which is what actually enforces this (RuleFunc order alone
// is not enough once an action drops out of scope; see its comment).
func RuleFuncs() []convospec.RuleFunc {
	return []convospec.RuleFunc{
		ruleListsRemoveDone,
		ruleListsRemoveAll,
		ruleListsSetDone,
		ruleListsRemoveItems,
		ruleListsView,
		ruleListsAdd,
	}
}

// ruleListsRemoveDone: "remove done"/"remove bought items from the movies
// list". The list name is optional; when named, it must resolve through
// listIDFromWord (a function, not a template) to the right listID.
func ruleListsRemoveDone(text string, _ time.Time) ([]convospec.ActionCall, bool) {
	m := reRemoveDone.FindStringSubmatch(text)
	if m == nil {
		return nil, false
	}
	args := map[string]any{}
	if listID := listIDFromWord(m[1]); listID != "" {
		args["listID"] = listID
	}
	return []convospec.ActionCall{{ActionID: removeDoneDef.ID, Args: args}}, true
}

// ruleListsRemoveAll: "remove all"/"remove everything"/"clear the books list".
// Two independent phrasings ("remove all ..." and "clear ... list") feed the
// same action, which a single declarative Rule cannot express.
func ruleListsRemoveAll(text string, _ time.Time) ([]convospec.ActionCall, bool) {
	var listWord string
	if m := reRemoveAll.FindStringSubmatch(text); m != nil {
		listWord = m[1]
	} else if m := reClearList.FindStringSubmatch(text); m != nil {
		listWord = m[1]
	} else {
		return nil, false
	}
	args := map[string]any{}
	if listID := listIDFromWord(listWord); listID != "" {
		args["listID"] = listID
	}
	return []convospec.ActionCall{{ActionID: removeAllDef.ID, Args: args}}, true
}

// ruleListsSetDone: "bought milk and eggs" / "mark milk as done" / "milk is
// done". Two independent phrasings feed one action, and the captured tail must
// be SPLIT into several item titles, not typed as one — a declarative Rule's
// ArgTypeStringSlice substitution always wraps a capture as a single-element
// slice, so "bought milk and eggs" would otherwise become one bad item titled
// "milk and eggs".
func ruleListsSetDone(text string, _ time.Time) ([]convospec.ActionCall, bool) {
	var rest string
	if m := reBought.FindStringSubmatch(text); m != nil {
		rest = m[1]
	} else if m := reMarkDone.FindStringSubmatch(text); m != nil {
		rest = m[1]
	} else {
		return nil, false
	}
	items := splitItems(rest)
	if len(items) == 0 {
		return nil, false
	}
	return []convospec.ActionCall{{ActionID: setDoneDef.ID, Args: map[string]any{"itemTitles": items, "isDone": true}}}, true
}

// ruleListsRemoveItems: "remove milk from the shopping list". It must
// UNCONDITIONALLY decline text that looks like a bulk-removal phrase — even
// when remove_done/remove_all are not themselves on offer this turn — because
// reRemoveItem's own pattern is structurally broad enough to also match
// "remove done items from the list" (capturing "done items" as if it were an
// item title). Declining unconditionally, rather than only when the bulk
// actions are absent from scope, is what keeps that phrase routed to clarify
// instead of silently misfiring as a one-off item removal once remove_done
// drops out of scope. This is exactly the "conditional target" case RuleFunc
// exists for.
func ruleListsRemoveItems(text string, _ time.Time) ([]convospec.ActionCall, bool) {
	if reRemoveDone.MatchString(text) || reRemoveAll.MatchString(text) || reClearList.MatchString(text) {
		return nil, false
	}
	m := reRemoveItem.FindStringSubmatch(text)
	if m == nil {
		return nil, false
	}
	items := splitItems(m[1])
	if len(items) == 0 {
		return nil, false
	}
	args := map[string]any{"itemTitles": items}
	if listID := listIDFromWord(m[2]); listID != "" {
		args["listID"] = listID
	}
	return []convospec.ActionCall{{ActionID: removeItemsDef.ID, Args: args}}, true
}

// ruleListsView: "show/view/what's on [the <list>] list".
func ruleListsView(text string, _ time.Time) ([]convospec.ActionCall, bool) {
	m := reViewList.FindStringSubmatch(text)
	if m == nil {
		return nil, false
	}
	args := map[string]any{}
	if listID := listIDFromWord(m[1]); listID != "" {
		args["listID"] = listID
	}
	return []convospec.ActionCall{{ActionID: viewDef.ID, Args: args}}, true
}

// ruleListsAdd: "add milk and bread to the movies list" / "buy milk and
// bread" / "need milk". Both phrasings must split a multi-item tail, and the
// first must resolve a NAMED list, so neither fits a declarative Rule.
func ruleListsAdd(text string, _ time.Time) ([]convospec.ActionCall, bool) {
	if m := reAddToList.FindStringSubmatch(text); m != nil {
		items := splitItems(m[1])
		if len(items) == 0 {
			return nil, false
		}
		args := map[string]any{"items": items}
		if listID := listIDFromWord(m[2]); listID != "" {
			args["listID"] = listID
		}
		return []convospec.ActionCall{{ActionID: addItemsDef.ID, Args: args}}, true
	}
	if m := reBuyItems.FindStringSubmatch(text); m != nil {
		items := splitItems(m[1])
		if len(items) == 0 {
			return nil, false
		}
		return []convospec.ActionCall{{ActionID: addItemsDef.ID, Args: map[string]any{"items": items}}}, true
	}
	return nil, false
}

// ScopedRuleFuncs covers the one interpretation that depends on which actions
// are on offer, not on the text alone — see ruleBareItemInListScope.
func ScopedRuleFuncs() []convospec.ScopedRuleFunc {
	return []convospec.ScopedRuleFunc{ruleBareItemInListScope}
}

// ruleBareItemInListScope: when ONLY the Listus catalog is in scope (an
// extension-specific bot like ListusBot), a bare short phrase such as "milk"
// most likely means "add milk to the groceries list". In the full action set
// the same phrase is genuinely ambiguous (a contact name? an asset? a
// tracked measurement?), so this must NOT fire there — the interpretation
// depends on WHICH actions are on offer this turn, which is exactly what
// ScopedRuleFunc exists for; neither a declarative Rule nor a plain RuleFunc
// can see that.
func ruleBareItemInListScope(text string, _ time.Time, availableActionIDs []string) ([]convospec.ActionCall, bool) {
	if !convospec.OnlyCatalogInScope(availableActionIDs, "lists.") {
		return nil, false
	}
	hasAddItems := false
	for _, id := range availableActionIDs {
		if id == addItemsDef.ID {
			hasAddItems = true
			break
		}
	}
	if !hasAddItems {
		return nil, false
	}
	words := strings.Fields(text)
	if len(words) == 0 || len(words) > 3 {
		return nil, false
	}
	// Command-like phrases ("list ...", "show ...") are not list items.
	switch words[0] {
	case "list", "show", "view", "what", "who", "how", "delete", "remove", "cancel", "add", "help", "start", "buy", "mark":
		return nil, false
	}
	items := splitItems(text)
	if len(items) == 0 {
		return nil, false
	}
	return []convospec.ActionCall{{ActionID: addItemsDef.ID, Args: map[string]any{"items": items}}}, true
}
