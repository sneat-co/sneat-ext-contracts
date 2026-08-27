# Chess Raiders contract

`contract4chessraiders` is the public, dependency-free contract for Chess
Raiders — the Sneat platform's real-time multiplayer chess game. Types,
interfaces, and errors only: what a host, a Telegram/social bot shell, or a
distribution-portal adapter needs in order to integrate with Chess Raiders,
without depending on the private engine (`sneat-co/chessraiders`). It carries
no board state, move, or rules-evaluation logic.

`contract4chessraiderstest` provides a conformance harness for
`MatchLobbyApplication`. `IdentityLinkApplication` and
`PortalSessionApplication` were left without one in this first release — see
the package doc comment.

This module holds one invariant, enforced by having zero `require` entries:
it depends on nothing but the Go standard library.

## Composition-threading note

The private `sneat-co/chessraiders` engine already implements
`contract4competios.GameLauncher` directly
(`server-go/facade4chess.NewCompetiosAdapter`), wired together by
`sneat-go/pkg/modules/competios/production_composition.go`. That means this
family and `competios` are NOT fully independent at the composition root: a
future migration of `sneat-co/chessraiders`'s `facade4chess` package onto
`sneat-ext-contracts/competios` (instead of the archived `ext-competios`) is
a prerequisite for that adapter to keep compiling, alongside the equivalent
migration of `sneat-go`'s own imports. Neither is done by this change — both
are explicitly out of scope here (no edits to `sneat-go`, `sneat-bots`, or any
game/engine repo).

## Provenance

Migrated from
`sneat-co/ext-chessraiders@34974ac9e0cf4d97a1a7dff5f49268b5af964341`
(`backend/contract4chessraiders` and `backend/contract4chessraiderstest`;
module renamed from `github.com/sneat-co/ext-chessraiders/backend` to
`github.com/sneat-co/sneat-ext-contracts/chessraiders`, package names
unchanged).

## No npm sibling (by design) and independent versioning

See `../competios/README.md`'s "No npm sibling (by design) and independent
versioning" section — identical reasoning applies here: no
`@sneat/extension-chessraiders-contract` npm package ever existed, none is
added by this change, `contracts.json` does not list `chessraiders`, and this
module's version comes from its own `version.go` (see `../README.md` "Go-only
family versioning"). Despite the composition-root coupling noted above, this
family's version line is fully independent of `competios`'s — the founder's
ruling ("each contract get versioned independently") applies per contract,
not per composition graph.

First release: `chessraiders/v0.1.0`, cut 2026-08-27 via the `go_module_tags`
manual dispatch escape hatch (the release-pipeline fix that makes this
automatic going forward had not yet landed when this family merged).
