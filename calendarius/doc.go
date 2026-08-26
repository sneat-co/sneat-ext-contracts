// Package calendarius is the root of the calendarius contract's Go module
// (github.com/sneat-co/sneat-ext-contracts/calendarius).
//
// This module holds the public contract surface of the calendarius
// extension — happening/event models (calendariusmodels), the convo-side
// facade (convo4calendarius) and its conformance suite
// (convo4calendariustest), and the recurring-happening facade
// (facade4calendarius) and its conformance suite (facade4calendariustest) —
// and depends only on foundational/core packages, never on another
// extension.
//
// Provenance: migrated from sneat-co/ext-calendarius's backend/ Go module
// (module github.com/sneat-co/ext-calendarius/backend, tag backend/v0.0.6,
// which contains the recurrence fix widening validation to the general
// weekly/fortnightly/monthly/yearly repeats vocabulary, PR #33). That source
// module had no root-level .go file of its own (no package to rename); this
// file exists so the module directory has a documented entry point. See
// ../libs/calendarius/README.md for the TS half's provenance and parity
// note — the two halves currently live in different upstream repos
// (sneat-co/ext-calendarius for Go, sneat-co/sneat-libs for TS) and are
// unified here for the first time.
package calendarius
