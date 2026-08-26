// Package backend is the root of the ext-sneat-team backend Go module.
//
// This module holds the public contract surface of the team extension —
// the follow DTO shapes and the Follower facade interface — and depends
// only on foundational/core packages, never on another extension.
//
// TypeSpec (../typespec/api4team.tsp) is the frozen wire contract; the
// Go types here are hand-implemented to match it (no emitters), per the house
// convention.
package backend
