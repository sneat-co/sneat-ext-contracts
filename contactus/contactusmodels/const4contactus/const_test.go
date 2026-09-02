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

// TestSpaceMemberRoleOwnerIsWellKnown closes finding 6 of the adversarial
// review of sneat-core-modules PRs #187/#188.
//
// `owner` was absent from SpaceMemberWellKnownRoles while every other role in
// this file was present. facade4invitus.grantsSpaceUserID checks membership of
// this list and logs
//
//	"invite grants space member role %q which is not a well-known
//	 const4contactus role"
//
// for anything missing -- so EVERY owner-granting claim raised that warning: a
// permanent false alarm in the exact detector PR #186 added to catch roles
// nobody has reasoned about. A detector that always fires on a legitimate,
// central role is a detector people learn to ignore.
func TestSpaceMemberRoleOwnerIsWellKnown(t *testing.T) {
	if SpaceMemberRoleOwner != "owner" {
		t.Fatalf("expected SpaceMemberRoleOwner to be %q, got %q", "owner", SpaceMemberRoleOwner)
	}
	if !slices.Contains(SpaceMemberWellKnownRoles, SpaceMemberRoleOwner) {
		t.Fatal("SpaceMemberRoleOwner should be in SpaceMemberWellKnownRoles")
	}
}

// TestEverySpaceMemberRoleConstIsWellKnown stops this defect recurring for the
// NEXT role, rather than pinning owner alone: every SpaceMemberRole constant
// declared in this package must appear in SpaceMemberWellKnownRoles.
//
// The list is maintained by hand next to the constants, so adding a constant
// and forgetting the list is the natural mistake -- it is exactly how `owner`
// went missing.
func TestEverySpaceMemberRoleConstIsWellKnown(t *testing.T) {
	for _, role := range []SpaceMemberRole{
		SpaceMemberRoleMember,
		SpaceMemberRoleExMember,
		SpaceMemberRoleAdult,
		SpaceMemberRoleChild,
		SpaceMemberRoleCreator,
		SpaceMemberRoleOwner,
		SpaceMemberRoleAdmin,
		SpaceMemberRoleContributor,
		SpaceMemberRoleSpectator,
		SpaceMemberRoleExcluded,
		SpaceMemberRoleCustomer,
	} {
		if !slices.Contains(SpaceMemberWellKnownRoles, role) {
			t.Errorf("SpaceMemberRole %q is declared as a constant but missing from SpaceMemberWellKnownRoles", role)
		}
	}
}

// TestIsKnownSpaceMemberRoleIsANoOpForNilSpaceRoles DOCUMENTS a pre-existing
// defect adjacent to finding 6; it deliberately does NOT fix it.
//
// IsKnownSpaceMemberRole reads
//
//	return spaceRoles == nil || slices.Contains(...) || slices.Contains(...)
//
// so with a nil spaceRoles it returns true for ANY string. The fleet's one
// caller -- sneat-core-modules userus/dbo4userus/user_space_brief.go -- passes
// nil, which makes that validation a no-op: a UserSpaceBrief carrying a
// garbage role validates clean.
//
// Fixing it is a behaviour change in a published contract with an unknown
// blast radius across the fleet (every brief whose roles were never really
// checked would start failing validation), so it is a separate, deliberate
// decision rather than a rider on a security fix. This test pins the CURRENT
// behaviour so the change, when it comes, is visible and intentional.
func TestIsKnownSpaceMemberRoleIsANoOpForNilSpaceRoles(t *testing.T) {
	if !IsKnownSpaceMemberRole("total-nonsense-not-a-role", nil) {
		t.Fatal("PRE-EXISTING behaviour changed: nil spaceRoles no longer short-circuits to true." +
			" That may well be the right fix -- but it is a contract behaviour change," +
			" so update this test deliberately and check every caller.")
	}
}
