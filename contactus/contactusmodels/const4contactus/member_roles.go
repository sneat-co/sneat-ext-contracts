package const4contactus

import (
	"slices"
)

type SpaceMemberRole = string

const (
	SpaceMemberRoleMember   = "member"
	SpaceMemberRoleExMember = "ex-member"

	SpaceMemberRoleAdult = "adult"

	SpaceMemberRoleChild = "child"

	// SpaceMemberRoleCreator role of a creator
	SpaceMemberRoleCreator SpaceMemberRole = "creator"

	// SpaceMemberRoleOwner role of an owner
	SpaceMemberRoleOwner SpaceMemberRole = "owner"

	SpaceMemberRoleAdmin SpaceMemberRole = "admin"

	// SpaceMemberRoleContributor role of a contributor
	SpaceMemberRoleContributor SpaceMemberRole = "contributor"

	// SpaceMemberRoleSpectator role of spectator
	SpaceMemberRoleSpectator SpaceMemberRole = "spectator"

	// SpaceMemberRoleExcluded if space members are excluded
	SpaceMemberRoleExcluded SpaceMemberRole = "excluded"

	// SpaceMemberRoleCustomer is a paying/participating patron of a Space (e.g.
	// a player invited to a venue), as distinct from SpaceMemberRoleMember,
	// which on a venue Space means an employee entitled to manage the venue's
	// content. A customer must never be granted SpaceMemberRoleMember. It is
	// also distinct from SpaceMemberRoleSpectator, which is a passive observer
	// role; a customer is an active/paying participant (founder ruling,
	// 2026-08-19).
	SpaceMemberRoleCustomer SpaceMemberRole = "customer"
)

// SpaceMemberWellKnownRoles defines known roles
var SpaceMemberWellKnownRoles = []SpaceMemberRole{
	SpaceMemberRoleAdmin,
	SpaceMemberRoleContributor,
	SpaceMemberRoleCreator,
	SpaceMemberRoleMember,
	SpaceMemberRoleExMember,
	SpaceMemberRoleChild,
	SpaceMemberRoleAdult,
	SpaceMemberRoleSpectator,
	SpaceMemberRoleExcluded,
	SpaceMemberRoleCustomer,
}

// IsKnownSpaceMemberRole checks if a role has a valid value
func IsKnownSpaceMemberRole(role SpaceMemberRole, spaceRoles []SpaceMemberRole) bool {
	return spaceRoles == nil || slices.Contains(SpaceMemberWellKnownRoles, role) || slices.Contains(spaceRoles, role)
}
