package contract4contactus

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/sneat-co/sneat-go-core/coretypes"
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
	for _, mutate := range []func(*ContactAuthorityEvent){func(e *ContactAuthorityEvent) { e.V = 0 }, func(e *ContactAuthorityEvent) { e.SpaceID = "" }, func(e *ContactAuthorityEvent) { e.SpaceID = "foreign/path" }, func(e *ContactAuthorityEvent) { e.ContactID = "" }, func(e *ContactAuthorityEvent) { e.ObservedAt = time.Time{} }, func(e *ContactAuthorityEvent) { e.MembershipGeneration = "not-hex" }, func(e *ContactAuthorityEvent) { e.BindingGeneration = "" }, func(e *ContactAuthorityEvent) { e.BindingGeneration = strings.Repeat("B", 32) }, func(e *ContactAuthorityEvent) { e.Status = "unknown" }, func(e *ContactAuthorityEvent) { e.TransitionKind = "unknown_kind" }} {
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

func TestAuthorityEventKnownIDStableAcrossTransitionKindAndTime(t *testing.T) {
	const knownID = "82f8efd5810c512a47e804ce97bc1a79445e42a29879cf9d3c8344ca849f846c"
	e := validAuthorityEvent()
	if got := e.ID(); got != knownID {
		t.Fatalf("published authority tuple identity changed: %s", got)
	}
	for _, kind := range []ContactAuthorityTransitionKind{ContactAuthorityTransitionMetadataRepair, ContactAuthorityTransitionNative, ContactAuthorityTransitionCreation, ContactAuthorityTransitionUnknown} {
		e.TransitionKind = kind
		e.ObservedAt = e.ObservedAt.Add(time.Hour)
		if got := e.ID(); got != knownID {
			t.Fatalf("transition kind %q or observation time changed event identity: %s", kind, got)
		}
	}
}

func TestAuthorityEventTransitionKindValidation(t *testing.T) {
	for _, kind := range []ContactAuthorityTransitionKind{ContactAuthorityTransitionUnknown, ContactAuthorityTransitionMetadataRepair, ContactAuthorityTransitionNative, ContactAuthorityTransitionCreation} {
		e := validAuthorityEvent()
		e.TransitionKind = kind
		if err := e.Validate(); err != nil {
			t.Fatalf("known transition kind %q rejected: %v", kind, err)
		}
	}
	for _, kind := range []ContactAuthorityTransitionKind{"unknown_kind", "authority-transition", " ", "Creation"} {
		e := validAuthorityEvent()
		e.TransitionKind = kind
		if err := e.Validate(); err == nil {
			t.Fatalf("invalid transition kind %q accepted", kind)
		}
	}
}

func TestAuthorityEventTransitionKindJSONAndFirestoreOptional(t *testing.T) {
	field, ok := reflect.TypeOf(ContactAuthorityEvent{}).FieldByName("TransitionKind")
	if !ok {
		t.Fatal("transition kind field is missing")
	}
	if got := field.Tag.Get("json"); got != "transitionKind,omitempty" {
		t.Fatalf("unexpected JSON tag: %q", got)
	}
	if got := field.Tag.Get("firestore"); got != "transitionKind,omitempty" {
		t.Fatalf("unexpected Firestore tag: %q", got)
	}
	e := validAuthorityEvent()
	raw, err := json.Marshal(e)
	if err != nil {
		t.Fatal(err)
	}
	var legacy map[string]any
	if err = json.Unmarshal(raw, &legacy); err != nil {
		t.Fatal(err)
	}
	if _, present := legacy["transitionKind"]; present {
		t.Fatal("legacy empty transition kind was not omitted")
	}
	e.TransitionKind = ContactAuthorityTransitionNative
	raw, err = json.Marshal(e)
	if err != nil {
		t.Fatal(err)
	}
	if err = json.Unmarshal(raw, &legacy); err != nil {
		t.Fatal(err)
	}
	if got := legacy["transitionKind"]; got != string(ContactAuthorityTransitionNative) {
		t.Fatalf("transition kind not serialized: %v", got)
	}
}

func TestClassifyContactAuthorityTransition(t *testing.T) {
	after := ContactAuthorityState{UserID: "user", Status: "active", Member: true}
	if got := ClassifyContactAuthorityTransition(nil, after); got != ContactAuthorityTransitionCreation {
		t.Fatalf("nil before classified as %q", got)
	}
	unchanged := after
	if got := ClassifyContactAuthorityTransition(&unchanged, after); got != ContactAuthorityTransitionMetadataRepair {
		t.Fatalf("unchanged facts classified as %q", got)
	}
	for name, mutate := range map[string]func(*ContactAuthorityState){
		"user":   func(s *ContactAuthorityState) { s.UserID = "another" },
		"status": func(s *ContactAuthorityState) { s.Status = "archived" },
		"member": func(s *ContactAuthorityState) { s.Member = false },
	} {
		before := after
		mutate(&before)
		if got := ClassifyContactAuthorityTransition(&before, after); got != ContactAuthorityTransitionNative {
			t.Fatalf("%s fact delta classified as %q", name, got)
		}
	}
}

func TestClassifyContactAuthorityTransitionHasNoEpochOnlyCause(t *testing.T) {
	stateType := reflect.TypeOf(ContactAuthorityState{})
	if stateType.NumField() != 3 {
		t.Fatalf("classifier state must expose only scalar membership facts, got %d fields", stateType.NumField())
	}
	expected := map[string]bool{"UserID": true, "Status": true, "Member": true}
	for i := 0; i < stateType.NumField(); i++ {
		if name := stateType.Field(i).Name; !expected[name] {
			t.Fatalf("unexpected classifier input %q", name)
		}
	}
	after := ContactAuthorityState{UserID: "user", Status: "active", Member: true}
	before := after
	// With no epoch, timestamp or boolean cause input, the same scalar facts
	// cannot be inferred to be a genuine authority transition.
	if got := ClassifyContactAuthorityTransition(&before, after); got != ContactAuthorityTransitionMetadataRepair {
		t.Fatalf("same scalar facts classified as %q", got)
	}
}

func TestAuthorityEventAcceptsNativeSpotSpaceID(t *testing.T) {
	e := validAuthorityEvent()
	e.SpaceID = coretypes.SpotSpaceIDPrefix + "native_spot"
	if err := e.Validate(); err != nil {
		t.Fatal(err)
	}
}
