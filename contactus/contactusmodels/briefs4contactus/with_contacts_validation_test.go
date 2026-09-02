package briefs4contactus

import (
	"strings"
	"testing"

	"github.com/sneat-co/sneat-go-core/models/dbmodels"
)

// TestWithMultiSpaceContactsValidateDoesNotSwallowErrors is a straight bug fix.
//
// WithMultiSpaceContacts.Validate() read:
//
//	if err := v.WithMultiSpaceContactIDs.Validate(); err != nil {
//		return nil
//	}
//
// It discarded the verdict of a real format validator -- one that checks each
// contactID is non-blank, has no surrounding whitespace, and is in the
// "spaceID:contactID" shape. WithMultiSpaceContacts is embedded in contactus's
// ContactDbo, so that validation has never actually run on a contact record:
// every malformed contactIDs entry validated clean.
func TestWithMultiSpaceContactsValidateDoesNotSwallowErrors(t *testing.T) {
	// contactIDs[0] is a required "*" sentinel (WithContactIDs.Validate), which
	// is why the per-id loop starts at [1:]. A malformed id therefore has to
	// sit at index 1 or later.
	const sentinel = "*"
	valid := string(dbmodels.NewSpaceItemID("space1", "contact1"))

	for _, tt := range []struct {
		name string
		bad  string
	}{
		{"missing separator", "nosuchseparator"},
		{"empty spaceID", ":contact2"},
		{"empty contactID", "space2:"},
		{"leading whitespace", " space2:contact2"},
		{"blank", ""},
	} {
		t.Run(tt.name, func(t *testing.T) {
			v := &WithMultiSpaceContacts[*ContactBrief]{}
			v.ContactIDs = []string{sentinel, valid, tt.bad}

			err := v.Validate()
			if err == nil {
				t.Fatalf("BUG: a malformed contactID (%q) validated clean --"+
					" WithMultiSpaceContacts.Validate() is discarding the verdict of"+
					" WithMultiSpaceContactIDs.Validate()", tt.bad)
			}
			// ...and it must fail for the ID-FORMAT reason, not incidentally on
			// the ids-vs-briefs check that runs afterwards.
			if !strings.Contains(err.Error(), "contactIDs[") {
				t.Fatalf("expected a contactIDs format error, got: %v", err)
			}
		})
	}

	t.Run("a well-formed set still validates", func(t *testing.T) {
		// Sentinel only: every id must have a matching brief in Contacts
		// (ValidateWithIdsAndBriefs), which is a separate rule from the id
		// format this test is about.
		v := &WithMultiSpaceContacts[*ContactBrief]{}
		v.ContactIDs = []string{sentinel}
		if err := v.Validate(); err != nil {
			t.Fatalf("a valid contactIDs set must still pass: %v", err)
		}
	})
}

// TestWithMultiSpaceContactsValidateEnforcesTheSentinel covers the other half
// the swallow hid: WithContactIDs.Validate() requires contactIDs to be
// non-empty and to start with the "*" sentinel. Both verdicts were discarded
// too, so a contact record with a completely absent or unseeded contactIDs
// array also validated clean.
//
// Note WithMultiSpaceContactIDs.AddSpaceContactID does NOT seed that sentinel,
// unlike its single-space sibling AddContactID -- so anything built purely
// through it will now fail validation. That asymmetry is pre-existing and is
// left alone here; this test simply makes the requirement visible.
func TestWithMultiSpaceContactsValidateEnforcesTheSentinel(t *testing.T) {
	t.Run("empty contactIDs is refused", func(t *testing.T) {
		v := &WithMultiSpaceContacts[*ContactBrief]{}
		if err := v.Validate(); err == nil {
			t.Fatal("BUG: an empty contactIDs validated clean")
		}
	})

	t.Run("a missing sentinel is refused", func(t *testing.T) {
		v := &WithMultiSpaceContacts[*ContactBrief]{}
		v.ContactIDs = []string{string(dbmodels.NewSpaceItemID("space1", "contact1"))}
		if err := v.Validate(); err == nil {
			t.Fatal("BUG: contactIDs without the leading \"*\" sentinel validated clean")
		}
	})
}
