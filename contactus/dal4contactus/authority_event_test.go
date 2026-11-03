package dal4contactus

import (
	"strings"
	"testing"
	"time"

	"github.com/sneat-co/sneat-ext-contracts/contactus/contract4contactus"
)

func TestAuthorityEventNativeRootAndValidation(t *testing.T) {
	e := contract4contactus.ContactAuthorityEvent{V: 1, SpaceID: "space", ContactID: "contact", MembershipGeneration: strings.Repeat("a", 32), BindingGeneration: strings.Repeat("b", 32), Status: "active", ObservedAt: time.Now().UTC()}
	r, err := NewAuthorityEventRecord(e)
	if err != nil {
		t.Fatal(err)
	}
	expected := "spaces/space/ext/contactus/authorityEvents/" + e.ID()
	if r.Key().String() != expected {
		t.Fatalf("event escaped native extension root: %s", r.Key())
	}
	e.V = 0
	if r, err = NewAuthorityEventRecord(e); err == nil || r != nil {
		t.Fatal("invalid event received a persistence envelope")
	}
	e.V = 1
	e.TransitionKind = "unknown_kind"
	if r, err = NewAuthorityEventRecord(e); err == nil || r != nil {
		t.Fatal("invalid transition kind received a persistence envelope")
	}
}
