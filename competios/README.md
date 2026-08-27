# Competios contract

`contract4competios` is the public, dependency-free contract for Competios —
the Sneat competitions and tournament engine. Types and interfaces only: what
a game must implement to be launched as a contest, and what it must report
back when the contest ends. The engine itself is private
(`sneat-co/competios`).

`contract4competiostest` provides reusable conformance harnesses an
implementor runs against its own adapter.

This module holds one invariant, enforced by having zero `require` entries:
it depends on nothing but the Go standard library. A dependency arriving here
would reintroduce, one level down, exactly the credential requirement a
public contract exists to remove.

## Provenance

Migrated from `sneat-co/ext-competios@9928d7f035f4c20e5e3e8ad3d9ebeb2f0a630cf8`
(`backend/contract4competios` and `backend/contract4competiostest`; module
renamed from `github.com/sneat-co/ext-competios/backend` to
`github.com/sneat-co/sneat-ext-contracts/competios`, package names unchanged).

`ext-competios`'s third package, `grants` (a separate Go module implementing
`contract4competios`'s `OperationGrantIssuer`/`OperationGrantVerifier` ports
with HMAC-signed tokens and dalgo-backed replay stores), is implementation,
not contract — it depends on `dal-go/dalgo` and `golang-jwt/jwt`, which this
module's zero-dependency invariant forbids. Per founder decision it belongs in
`sneat-co/competios`, not here; moving it there is separate, later work.

## No npm sibling (by design) and independent versioning

Every other family in this repo ships an `@sneat/extension-<family>-contract`
npm package alongside its Go module. Neither `ext-competios` nor
`ext-chessraiders` ever had one, and no npm package
`@sneat/extension-competios-contract` exists on the registry today — verified
via `npm view`, and no repo in the fleet imports one. This module is
deliberately Go-only. If an Angular/TS consumer needs this contract later, add
`libs/competios/` then; nothing about this module's shape requires it.

Consequence: `contracts.json` does NOT list `competios` (that manifest feeds
the npm-side tier-coherence check, and there is no npm package here for it to
cover), and this module's version comes from its own `version.go`, not from
any npm sibling — see `../README.md` "Go-only family versioning" for the full
mechanism and journey. `chessraiders` (`../chessraiders/README.md`) is in the
same shape but versions **independently**: bumping this family's
`version.go` never bumps `chessraiders`'s, and vice versa, even though the two
are wired together at the `sneat-go` composition root (see
`../chessraiders/README.md`'s "Composition-threading note").

First release: `competios/v0.1.0`, cut 2026-08-27 via the `go_module_tags`
manual dispatch escape hatch (the release-pipeline fix that makes this
automatic going forward had not yet landed when this family merged).
