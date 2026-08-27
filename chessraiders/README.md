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

## No npm sibling (yet)

See `../competios/README.md`'s "No npm sibling" section — identical
reasoning and identical release-pipeline gap apply here: no
`@sneat/extension-chessraiders-contract` npm package ever existed, none is
added by this change, `contracts.json` does not list `chessraiders`, and
`.github/workflows/publish.yml`'s Go-tag resolution will not fire for this
module without either (a) a `libs/chessraiders/package.json` added later, or
(b) using the workflow's existing manual `go_module_tags` escape hatch to tag
it directly — see this task's report for the exact recommended dispatch.
