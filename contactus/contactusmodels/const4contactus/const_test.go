package const4contactus

import (
	"slices"
	"testing"
)

func TestKnownPetKinds(t *testing.T) {
	if !IsKnownPetPetKind(PetKindDog) {
		t.Fatal("dog should be known")
	}
	if IsKnownPetPetKind("dragon") {
		t.Fatal("dragon should not be known")
	}
}

func TestKnownSpaceMemberRole(t *testing.T) {
	if !IsKnownSpaceMemberRole(SpaceMemberRoleMember, nil) {
		t.Fatal("standard role should be known")
	}
	if !IsKnownSpaceMemberRole("custom", []SpaceMemberRole{"custom"}) {
		t.Fatal("configured custom role should be known")
	}
	if IsKnownSpaceMemberRole("custom", []SpaceMemberRole{"other"}) {
		t.Fatal("unconfigured custom role should not be known")
	}
}

// TestSpaceMemberRoleCustomerIsWellKnownAndDistinct asserts the founder's
// 2026-08-19 ruling: player invites grant a "customer" role that is a
// standard/well-known SpaceMemberRole in its own right, distinct in value
// from both SpaceMemberRoleMember (venue-content-management authority) and
// SpaceMemberRoleSpectator (a passive observer role).
func TestSpaceMemberRoleCustomerIsWellKnownAndDistinct(t *testing.T) {
	if SpaceMemberRoleCustomer != "customer" {
		t.Fatalf("expected SpaceMemberRoleCustomer to be %q, got %q", "customer", SpaceMemberRoleCustomer)
	}
	if !slices.Contains(SpaceMemberWellKnownRoles, SpaceMemberRoleCustomer) {
		t.Fatal("SpaceMemberRoleCustomer should be in SpaceMemberWellKnownRoles")
	}
	if !IsKnownSpaceMemberRole(SpaceMemberRoleCustomer, nil) {
		t.Fatal("SpaceMemberRoleCustomer should be a known standard role")
	}
	if SpaceMemberRoleCustomer == SpaceMemberRoleMember {
		t.Fatal("SpaceMemberRoleCustomer must never equal SpaceMemberRoleMember")
	}
	if SpaceMemberRoleCustomer == SpaceMemberRoleSpectator {
		t.Fatal("SpaceMemberRoleCustomer must be distinct from SpaceMemberRoleSpectator")
	}
}
