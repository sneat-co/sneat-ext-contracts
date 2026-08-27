// Package competios holds the version identity of this Go module only —
// see ../README.md "Go-only family versioning" for why this file exists
// and how the release pipeline reads it.
package competios

// Version holds the version of this Go module (mirrors the
// <family>/version.go convention already used across the Sneat/dalgo Go
// ecosystem, e.g. dal-go/dalgo/version.go). Bump it, in the same PR that
// changes this family's contract, to the next version this module should
// be tagged at; .github/workflows/publish.yml reads it and tags
// competios/v<Version> the next time it has not already been tagged.
const Version = "0.1.0"
