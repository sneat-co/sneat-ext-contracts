// Package gameboard is the root of the gameboard contract's Go module
// (github.com/sneat-co/sneat-ext-contracts/gameboard).
//
// This module holds the public contract surface of the gameboard extension —
// the event-timeline model/const shapes, the inline team-side DTO, and the
// deterministic fold reducer (in the eventtimeline sub-package) — and
// depends only on foundational/core packages, never on another extension.
//
// The Go types here are hand-implemented to match the frozen wire contract
// (api4gameboard.tsp, kept in sneat-co/ext-gameboard) per the house
// no-emitter convention. The deterministic fold reducer lives in the
// contract so the backend and the frontend reducer fold the same log to
// identical state.
//
// Provenance: migrated from sneat-co/ext-gameboard's backend/ Go module
// (module github.com/sneat-co/ext-gameboard/backend, package "backend") —
// the root package is renamed backend -> gameboard here since the module
// directory is no longer named backend/. See ../libs/gameboard/README.md
// for the full provenance and parity note.
package gameboard
