package botapp

import "testing"

func TestSpaceRefCarriesOnlyPresentationFields(t *testing.T) {
	ref := NewSpaceRef("family-1", "family")
	if ref.ID != "family-1" || ref.Type != "family" {
		t.Fatalf("unexpected space reference: %#v", ref)
	}
}
