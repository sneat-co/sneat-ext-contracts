// Copyright 2026 Sneat.app

// Package botapp contains the persistence-neutral Listus contract consumed by
// delivery adapters and implemented at the application composition boundary.
package botapp

import "context"

// SpaceRef is the presentation-level space selected by a Listus bot. It does
// not reveal a host data model or database implementation.
type SpaceRef struct {
	ID   string
	Type string
}

// NewSpaceRef constructs a delivery-safe space reference.
func NewSpaceRef(id, spaceType string) SpaceRef {
	return SpaceRef{ID: id, Type: spaceType}
}

// HostServices resolves a bot's selected space. The host remains responsible
// for user/profile reads and any space provisioning.
type HostServices interface {
	ResolveSpace(ctx context.Context, userID string, requested SpaceRef) (SpaceRef, error)
}
