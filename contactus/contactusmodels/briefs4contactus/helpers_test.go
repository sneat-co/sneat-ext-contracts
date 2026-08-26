package briefs4contactus

import (
	"testing"

	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
)

func TestGetFullContactID(t *testing.T) {
	if got := GetFullContactID("space1", "contact1"); got != "space1:contact1" {
		t.Fatalf("GetFullContactID() = %q", got)
	}
	for _, test := range []struct {
		name, spaceID, contactID string
	}{
		{name: "missing space", contactID: "contact1"},
		{name: "missing contact", spaceID: "space1"},
	} {
		t.Run(test.name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatal("GetFullContactID() did not panic")
				}
			}()
			GetFullContactID(coretypes.SpaceID(test.spaceID), test.contactID)
		})
	}
}

func TestIsUniqueShortTitle(t *testing.T) {
	contacts := map[string]*ContactBrief{
		"member": {ShortTitle: "Alex"},
		"other":  {ShortTitle: "Sam"},
	}
	if IsUniqueShortTitle("Alex", contacts, "") {
		t.Fatal("Alex should not be unique")
	}
	if !IsUniqueShortTitle("Taylor", contacts, "") {
		t.Fatal("Taylor should be unique")
	}
}

func TestWithContactIDs(t *testing.T) {
	for _, test := range []struct {
		name string
		ids  []string
		want bool
	}{
		{name: "missing", want: false},
		{name: "missing sentinel", ids: []string{"contact1"}, want: false},
		{name: "valid sentinel", ids: []string{"*"}, want: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			err := (&WithContactIDs{ContactIDs: test.ids}).Validate()
			if (err == nil) != test.want {
				t.Fatalf("Validate() error = %v, want valid=%v", err, test.want)
			}
		})
	}

	single := &WithSingleSpaceContactIDs{}
	single.AddContactID("contact1")
	if got, want := single.ContactIDs, []string{"*", "contact1"}; len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("AddContactID() = %#v, want %#v", got, want)
	}
	if !single.HasContactID("contact1") || single.HasContactID("missing") {
		t.Fatalf("HasContactID() did not reflect %#v", single.ContactIDs)
	}
	if err := single.Validate(); err != nil {
		t.Fatalf("single Validate() = %v", err)
	}
}

func TestWithMultiSpaceContactIDs(t *testing.T) {
	valid := &WithMultiSpaceContactIDs{WithContactIDs: WithContactIDs{ContactIDs: []string{"*"}}}
	valid.AddSpaceContactID(dbmodels.NewSpaceItemID("space1", "contact1"))
	if err := valid.Validate(); err != nil {
		t.Fatalf("valid Validate() = %v", err)
	}
	if !valid.HasSpaceContactID(dbmodels.NewSpaceItemID("space1", "contact1")) {
		t.Fatal("HasSpaceContactID() = false")
	}
	for _, ids := range [][]string{
		{"*", " "},
		{"*", " space1:contact1"},
		{"*", "not-a-space-item"},
		{"*", ":contact1"},
		{"*", "space1:"},
	} {
		if err := (&WithMultiSpaceContactIDs{WithContactIDs: WithContactIDs{ContactIDs: ids}}).Validate(); err == nil {
			t.Fatalf("Validate() accepted invalid IDs %#v", ids)
		}
	}
}

func TestContactGroupBriefValidate(t *testing.T) {
	if err := (&ContactGroupBrief{}).Validate(); err == nil {
		t.Fatal("empty title should be invalid")
	}
	if err := (&ContactGroupBrief{Title: "Friends"}).Validate(); err != nil {
		t.Fatalf("valid title: %v", err)
	}
}
