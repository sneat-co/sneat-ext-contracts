package convo4listus

import (
	"testing"
	"time"

	"github.com/sneat-co/sneat-go-core/convospec"
)

// This pins the bug that escaped all the way to sneat-go's host smoke test:
// NormalizeText used to replace a comma with a space, so
// "buy milk, bread, eggs and apples" had nothing left to split on but "and" and
// produced 2 items ("milk bread eggs", "apples") instead of 4. Nothing at this
// level or in sneat-bots covered a comma-separated list, which is why it reached
// two repos downstream before failing.
func TestCommaListAndQuantity(t *testing.T) {
	for _, tt := range []struct {
		text      string
		wantItems int
	}{
		{"buy milk, bread, eggs and apples", 4},
		{"buy milk and bread", 2},
		{"buy milk", 1},
	} {
		n := convospec.NormalizeText(tt.text)
		for _, f := range RuleFuncs() {
			if calls, ok := f(n, time.Time{}); ok {
				items, _ := calls[0].Args["items"].([]string)
				if len(items) != tt.wantItems {
					t.Errorf("%q -> %d items %v, want %d (normalized %q)", tt.text, len(items), items, tt.wantItems, n)
				}
				break
			}
		}
	}
	// The quantified form must survive the comma token too.
	for _, text := range []string{"milk 2 liters", "milk, 2 liters"} {
		n := convospec.NormalizeText(text)
		var hit bool
		for _, r := range Rules() {
			d, _ := Declaration().Action(r.ActionID)
			if args, ok := r.Match(n, d); ok {
				hit = true
				if args["quantity"] != 2.0 || args["unit"] != "L" {
					t.Errorf("%q -> %v (normalized %q)", text, args, n)
				}
				break
			}
		}
		if !hit {
			t.Errorf("%q matched no rule (normalized %q)", text, n)
		}
	}
}
