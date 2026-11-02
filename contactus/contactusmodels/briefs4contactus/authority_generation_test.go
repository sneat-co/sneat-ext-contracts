package briefs4contactus

import (
	"encoding/json"
	"testing"
)

func TestAuthorityGenerationsRotateWithoutABA(t *testing.T) {
	var c ContactBase
	if c.MembershipGeneration != "" || c.BindingGeneration != "" {
		t.Fatal("legacy epochs must be absent")
	}
	c.RotateAuthorityGenerations()
	firstMembership, firstBinding := c.MembershipGeneration, c.BindingGeneration
	if len(firstMembership) != 32 || len(firstBinding) != 32 || firstMembership == firstBinding {
		t.Fatal("expected distinct random 128-bit epochs")
	}
	c.RotateAuthorityGenerations()
	if c.MembershipGeneration == firstMembership || c.BindingGeneration == firstBinding {
		t.Fatal("authority epoch reused")
	}
	bytes, err := json.Marshal(c)
	if err != nil {
		t.Fatal(err)
	}
	var restored ContactBase
	if err = json.Unmarshal(bytes, &restored); err != nil {
		t.Fatal(err)
	}
	if restored.MembershipGeneration != c.MembershipGeneration || restored.BindingGeneration != c.BindingGeneration {
		t.Fatal("epochs not persisted")
	}
}
