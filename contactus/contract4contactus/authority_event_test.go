package contract4contactus

import (
	"github.com/sneat-co/sneat-go-core/coretypes"
	"strings"
	"testing"
	"time"
)

func validAuthorityEvent() ContactAuthorityEvent {
	return ContactAuthorityEvent{V: 1, SpaceID: "space", ContactID: "contact", UserID: "user", MembershipGeneration: strings.Repeat("a", 32), BindingGeneration: strings.Repeat("b", 32), Status: "active", Member: true, ObservedAt: time.Now().UTC()}
}
func TestAuthorityEventValidationAndIdentity(t *testing.T) {
	e := validAuthorityEvent()
	if err := e.Validate(); err != nil {
		t.Fatal(err)
	}
	id := e.ID()
	if len(id) != 64 {
		t.Fatal("event identity is not SHA256")
	}
	retry := e
	retry.ObservedAt = retry.ObservedAt.Add(time.Minute)
	if retry.ID() != id {
		t.Fatal("retry observation changed event identity")
	}
	for _, mutate := range []func(*ContactAuthorityEvent){func(e *ContactAuthorityEvent) { e.UserID = "another" }, func(e *ContactAuthorityEvent) { e.MembershipGeneration = strings.Repeat("c", 32) }, func(e *ContactAuthorityEvent) { e.BindingGeneration = strings.Repeat("c", 32) }, func(e *ContactAuthorityEvent) { e.Member = false }, func(e *ContactAuthorityEvent) { e.SpaceID = "other" }, func(e *ContactAuthorityEvent) { e.ContactID = "other" }, func(e *ContactAuthorityEvent) { e.Status = "archived" }} {
		changed := e
		mutate(&changed)
		if changed.ID() == id {
			t.Fatal("distinct native authority tuple reused event identity")
		}
	}
	for _, mutate := range []func(*ContactAuthorityEvent){func(e *ContactAuthorityEvent) { e.V = 0 }, func(e *ContactAuthorityEvent) { e.SpaceID = "" }, func(e *ContactAuthorityEvent) { e.SpaceID = "foreign/path" }, func(e *ContactAuthorityEvent) { e.ContactID = "" }, func(e *ContactAuthorityEvent) { e.ObservedAt = time.Time{} }, func(e *ContactAuthorityEvent) { e.MembershipGeneration = "not-hex" }, func(e *ContactAuthorityEvent) { e.BindingGeneration = "" }, func(e *ContactAuthorityEvent) { e.BindingGeneration = strings.Repeat("B", 32) }, func(e *ContactAuthorityEvent) { e.Status = "unknown" }} {
		invalid := e
		mutate(&invalid)
		if invalid.Validate() == nil {
			t.Fatalf("invalid event accepted: %+v", invalid)
		}
	}
	for _, status := range []string{"active", "archived", "deleted"} {
		e.Status = status
		if err := e.Validate(); err != nil {
			t.Fatal(err)
		}
	}
}

func TestAuthorityEventAcceptsNativeSpotSpaceID(t *testing.T) {
	e := validAuthorityEvent()
	e.SpaceID = coretypes.SpotSpaceIDPrefix + "native_spot"
	if err := e.Validate(); err != nil {
		t.Fatal(err)
	}
}
