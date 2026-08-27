package briefs4contactus

import (
	"testing"

	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
)

// TestWithMultiSpaceContacts_SetContactBrief_NilContactsMap reproduces the
// panic reported in sneat-co/contactus#60: a freshly constructed parent
// record has a nil Contacts map (no constructor initializes it, and nothing
// guarantees a value was unmarshalled into it), so linking the very first
// child via SetContactBrief must not panic with "assignment to entry in nil
// map".
func TestWithMultiSpaceContacts_SetContactBrief_NilContactsMap(t *testing.T) {
	v := &WithMultiSpaceContacts[*ContactBrief]{}
	if v.Contacts != nil {
		t.Fatalf("test setup invalid: Contacts should start nil, got %#v", v.Contacts)
	}

	spaceID := coretypes.SpaceID("space1")
	contactID := "contact1"
	brief := &ContactBrief{Type: ContactTypePerson, Title: "Alex"}

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("SetContactBrief() panicked on nil Contacts map: %v", r)
		}
	}()

	updates := v.SetContactBrief(spaceID, contactID, brief)

	id := string(dbmodels.NewSpaceItemID(spaceID, contactID))
	if got := v.Contacts[id]; got != brief {
		t.Fatalf("Contacts[%q] = %#v, want %#v", id, got, brief)
	}
	if len(updates) != 2 {
		t.Fatalf("SetContactBrief() updates = %#v, want 2 entries (contactIDs + contact brief)", updates)
	}
}
