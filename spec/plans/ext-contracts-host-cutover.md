---
format: https://specscore.md/plan-specification
status: Draft
---
# Plan: Ext-Contracts Host Cutover

**Status:** Draft
**Source:** none
**Date:** 2026-08-27
**Owner:** alex
**Supersedes:** —

## Summary

All 30 `sneat-co/ext-*` repositories were archived on 2026-08-27. Archiving
does not break Go module resolution — nothing is broken this minute — but
those modules are frozen and will never take another fix. Everything that
still names an `ext-*` module path in `go.mod`, `go.sum`, or a `.go` import
has to move to its `github.com/sneat-co/sneat-ext-contracts/<family>`
successor, or be documented as blocked on a specific external repo.

This Plan covers the two large Go host repos only: `sneat-co/sneat-go` and
`sneat-co/sneat-bots`. Product/extension repos (`sneat-co/calendarius`,
`sneat-co/bookius`, `sneat-co/gameboard`, `sneat-co/kids-club`,
`sneat-co/sneat-team`, `sneat-co/listus`, `sneat-co/trackus`,
`sneat-co/togethered`, `sneat-co/eventius`, `sneat-co/sportius`,
`sneat-co/competios`, `sneat-co/chessraiders`, `sneat-co/debtus`,
`sneat-co/assetus`, `sneat-co/logistus`, `sneat-co/ovdb`,
`sneat-co/requoter`, `sneat-co/rosycycle`, `sneat-co/sneat-core-modules`,
`sneat-co/sneat-work`, and a third host, `sneat-co/sneat-cli`) are **out of
this Plan's authority** — every place one of them blocks a host-repo family
is called out as a named Open Question, not silently assumed or routed
around.

**Snapshot discipline.** Inventory below is measured against:

- `sneat-go` at `origin/main` = `f42a0ae7` (2026-08-26 20:21 UTC+1). A
  separate lane is actively landing `sneat-go-storedalgo-hotfix`, which
  migrates roughly ten old-contract sites plus a large dependency-pin set.
  **This branch/worktree was not read or touched to produce this Plan.**
  Task 1 below is the explicit re-audit gate for that fact — every other
  `sneat-go` task depends on it, directly or transitively.
- `sneat-bots` at `origin/main` = `6ce22b4` = tag `v0.30.7` (2026-08-27
  08:56 UTC+1, includes PR #144, "move off ext-contactus/backend to
  sneat-ext-contracts/contactus"). **A separate lane has since been
  dispatched to continue the `sneat-bots` cutover itself** (the
  `convosetup`/calendarius/listus/trackus composition hub described in
  Task 2). This Plan documents that scope and its ordering; it does not
  re-assign work already claimed there. Check that lane's live state before
  starting any `sneat-bots` task.

**What "done" means, mechanically, per family:** `grep -rn
"sneat-co/ext-<family>" go.mod go.sum $(git ls-files '*.go')` returns nothing
in either repo, and the `github.com/sneat-co/ext-<family>/*` line is absent
from both `go.mod` files — not "compiles," not "an agent reports success."
Several families cannot reach that state through host-repo work alone; those
are flagged explicitly rather than declared done when they aren't.

## Approach

### The end-to-end journey

1. A developer (or agent) picks one host-repo family task below. They open a
   worktree, find every import of `github.com/sneat-co/ext-<family>/backend/...`
   (or `/<family>/...` for the chessraiders/competios path shape) in that
   repo's source, and rewrite each to the matching package under
   `github.com/sneat-co/sneat-ext-contracts/<family>/...`, using the exact
   sub-package names already published there (e.g. `contactusmodels`,
   `dal4contactus`, `facade4contactus`, `contract4contactus` for contactus —
   confirmed present and tagged `contactus/v0.12.4`).
   **Observable good result (start):** the repo still builds and every
   existing test still passes with only the import path changed — no
   behavior moved, only where the type/port comes from.
2. They run `go mod tidy`.
   **Observable good result (middle):** if no other module in the graph
   still needs the old `ext-<family>` module, `tidy` drops it from
   `go.mod`/`go.sum` entirely. If `go mod graph` still shows another
   requirer — a tripwire dependency named below — `tidy` correctly *leaves*
   it in `go.mod` as `// indirect`. That is expected, not a bug: it is
   exactly what the per-family tripwire table predicts, and the task's
   completion note should name which requirer is still pulling it in.
3. They push. For `sneat-bots`, a local pre-push hook
   (`scripts/pre-push-sneat-go-e2e.sh`, installed via `.wb/templates/`) clones
   `sneat-go@origin/main` and runs its tests against the candidate build —
   this is the mechanism that made today's contactus migration provably safe
   without hand-tracing every call site.
   **Observable good result:** green run proves the family's new contract
   surface still satisfies both the local package's own tests and (for any
   family `sneat-go` also imports through `sneat-bots`) `sneat-go`'s
   consumption of it — *unless* the family is flagged "coordinate, don't
   assume independent" below, in which case a red run here is the expected,
   correct signal to stop and sequence, not a flake to retry past.
   **Trap:** this hook is a *local* pre-push hook, not a GitHub Actions
   check. It only fires on a machine that has it installed and does not pass
   `--no-verify`. Treat a green push from an unknown machine as unproven; run
   the script by hand as this task's explicit verification step if in doubt.
4. PR merges to the repo's `main`.
   **Observable good result — deliberately incomplete for `sneat-bots`:**
   merging to `sneat-bots` main delivers *nothing* downstream by itself.
   `sneat-bots`' `strongo_workflow` runs `disable-version-bumping: true` — no
   auto-tag. A human (or a task explicitly assigned this, see Task 6) must
   run `git tag vX.Y.Z && git push origin vX.Y.Z` by hand, exactly as
   happened today for `v0.30.7`. Until that tag exists, `sneat-go`'s
   `go.mod` still pins the old `sneat-bots` version and none of the
   `sneat-bots`-side fixes are visible to `sneat-go`.
5. Once a new `sneat-bots` tag exists, a `sneat-go` task (Task 12) bumps
   `github.com/sneat-co/sneat-bots` to it, runs `go mod tidy`, and re-audits
   which `ext-*` modules are now indirect-only survivors purely because of
   the stale pin.
   **Observable good result:** families that were only in `sneat-go`'s graph
   via the `sneat-bots` pin (`listus`, `trackus`) either disappear from
   `go.mod` or shrink to "blocked on `<named product repo>`" — a documented,
   not silent, remainder.
   **Trap, dated and confirmed today:** `sneat-bots@v0.30.7`'s `go.mod`
   declares `go 1.27.0`; `sneat-go`'s currently declares `go 1.26.1` and
   currently pins `sneat-bots@v0.30.6`, whose own `go.mod` is `go 1.26.1`
   (confirmed by reading both tags directly). The moment `sneat-go` bumps
   past `v0.30.6`, Go's module-graph max-`go`-directive rule drags
   `sneat-go`'s effective build toolchain to 1.27 for the first time.
   Per standing policy, a coverage floor is bound to the toolchain that
   measured it (go1.26 → go1.27 has been observed to re-measure ~3.5pp lower
   on identical source elsewhere in this fleet today). **A floor failing in
   an untouched `sneat-go` package right after this bump is a measurement
   artifact, not a regression — and lowering a floor is a founder decision,
   never an implementer's.** Task 12 isolates the bump into its own
   commit/PR specifically so this fallout is diagnosable in isolation from
   any family-migration diff landing the same week.
6. **Terminal state for this Plan's authority:** every host-owned source file
   is migrated, and every remaining `go.mod` entry for an `ext-*` module is
   either gone or attached to a named external product-repo blocker (see
   Open Questions). This Plan cannot, by itself, reach zero references
   fleet-wide — several families are blocked on 10+ other repos this Plan
   has no authority to change.

### Composition-threading — the tripwire class that makes "independent family" the wrong default unit

An empirical finding surfaced today while a parallel lane attempted the
`sneat-cli` cutover (a third host, out of this Plan's scope, mentioned here
only because it proves the mechanism): `sneat-cli`'s
`cmd/sneat/commands/convo_trackus.go` wires calendarius's convo service
*through* contactus's — `convoservice4calendarius.New(botservice4contactus.New())`.
That composition's compile-time type contract cascades into `sneat-bots`,
and produced this exact build failure:

```
cannot use botservice4contactus.New() (... "sneat-ext-contracts/contactus/contract4contactus".ConvoService)
as "github.com/sneat-co/ext-contactus/backend/contract4contactus".ConvoService value
in argument to convoservice4calendarius.New: ... wrong type for method CreateContact
```

Two things this proves, confirmed by direct inspection of the actual files
in this repo's clones (not inferred):

- **`sneat-bots/extensions/anybot/convo/convosetup/setup.go`** is a
  single composition-hub file. It imports
  `github.com/sneat-co/ext-calendarius/backend/convo4calendarius` directly
  for its `Services.Calendarius` field type, and its `NewCatalogRegistry`
  wires `actions4contactus.Catalog()`, `actions4calendarius.Catalog(...)`,
  `actions4listus.Catalog()`, and `actions4trackus.Catalog()` together in one
  place. `sneat-go` imports this exact package
  (`sneat-bots/extensions/anybot/convo/convosetup`) directly, so `sneat-go`'s
  own build has a real compile-time dependency on however this file's public
  surface (`Services`, `NewCatalogRegistry`) is typed — even though, per the
  file's own doc comment, "no host binds [a Calendarius service] today," so
  `sneat-go`'s current *runtime* exposure to the calendarius↔contactus type
  clash is lower than `sneat-cli`'s (which does try to bind one).
- `sneat-bots/extensions/listus/convoactions/catalog.go` and
  `extensions/trackus/convoactions/catalog.go` have the *same shape* of
  coupling one level down: each binds `ext-listus`/`ext-trackus`'s
  `Declaration`/`Rules`/`Port` types against a concrete implementation
  returned by the *product repo's own* `botapi.NewPort()`
  (`sneat-co/listus/backend/botapi`, `sneat-co/trackus/backend/botapi`).
  Migrating the host's import does nothing until that product repo's
  `botapi` return type also satisfies the new `sneat-ext-contracts/<family>`
  interface — the same "product-repo return type must match the interface
  the host imports" mechanism, not an independent per-file swap.

- **`calendarius/backend` bundled its own contract migration atomically at
  `v0.6.2`**: every version `≥ v0.6.2` speaks new-contactus and
  new-calendarius together; every version `≤ v0.6.1` speaks old-for-both.
  **There is no in-between release.** `sneat-bots@v0.30.7` (current, latest
  tag) still imports `ext-calendarius/backend/convo4calendarius` directly in
  `convosetup/setup.go`, so it cannot adopt a migrated `calendarius/backend`
  without also finishing the `convosetup.go` swap in the same change — and
  vice versa.

**Generalizing this, per the coordinator's explicit instruction:** treat
"bundled atomic migration in a dependency's release, no partial path" as its
own tripwire class, alongside the already-confirmed "transitive-forcing"
class (`debtus v0.2.33` and `sneat-core-modules v0.65.14` were *each
independently* sufficient to force the *new* `sneat-ext-contracts/contactus
v0.12.3` transitively — no partial path there either; a distinct claim from
Open Question 1 below, which is about the *old* `ext-contactus/backend`
still being required by 15 other repos as consumed today — both are true
simultaneously, at different points in the same graph). **Calendarius is the second
confirmed instance of the bundled-release class.** Whether `listus/backend`
and `trackus/backend` bundle their own switch the same atomic way as
`calendarius/backend` did at `v0.6.2` is **not yet confirmed** — Task 2 below
is explicitly scoped to establish that before assuming a partial path exists
for either.

A same-repo instance of the identical mechanism exists inside `sneat-go`
itself: `pkg/modules/competios/eventreg/` composes bookius, eventius, *and*
competios contract types together in one directory
(`bookius_ports.go`, `bookius_registrar.go`,
`eventius_attendance_projector.go`, `composition.go`, `participant_ref.go`).
`pkg/modules/competios/production_composition.go` additionally threads
sportius (`TeamRosterAuthority`) and chess/chessraiders
(`ChessExecutionProvider`) through the same composition root. Since
competios itself has no target contract yet (see Open Questions), **Tasks 8,
9, and 10 below explicitly exclude every file under `pkg/modules/competios/`**
— those files stay entangled with competios's unresolved status and are not
this Plan's to migrate in isolation. (Reinforcing evidence:
`sneat-ext-contracts/eventius/facade4eventius/` already contains
`competios_attendance.go` and `competios_attendance_commands.go` — eventius's
*own* published contract references competios by name, so even eventius's
contract-side story is not fully decoupled from competios's missing home.)

### Dependency-graph ordering

`sneat-go` requires `sneat-bots` (`sneat-go/go.mod:110`); `sneat-bots` has no
module dependency on `sneat-go`. So:

- Any `sneat-bots`-side family fix can be authored, tested, and merged to
  `sneat-bots` main independently of `sneat-go`'s schedule.
- But `sneat-go`'s `go.mod` cannot fully resolve an entry that flows through
  `sneat-bots` (`listus`, `trackus`, and `gameboard`'s `sneat-bots`-side
  copy) until `sneat-bots` (a) fixes it, (b) merges, (c) is manually tagged,
  and (d) `sneat-go` bumps its pin. Four sequential steps, not one — this is
  why Tasks 5/6/12 are separated rather than folded into the family tasks.
- For families where `sneat-go` imports a `sneat-bots` extension package
  directly (`bookius` via `extensions/bookius/vendorbot`, `eventius` via
  `extensions/eventius/bot/cmds4eventiusbot`, and `contactus`/`calendarius`
  via the `convosetup` hub above), the two sides' migrations must be
  verified together — via the reverse-integration gate or a manual run of
  it — before either side's change is considered final for that family.
  This Plan encodes that as a `Depends-On` edge rather than asserting blind
  independence for families that merely *look* separable file-by-file.
- Families with **zero** `sneat-go`-side source usage and no `sneat-go`
  import of the relevant `sneat-bots` package (`gameboard`, `kids-club`,
  `sneat-team`) are genuinely independent of this ordering — `sneat-go`'s
  only task for them is the pin-bump/re-audit in Task 13, and two of the
  three (`kids-club`, `sneat-team`) don't touch `sneat-bots` at all.

### Measured inventory (file counts, non-test / test / go.mod+go.sum+other, by family)

**`sneat-go`** (12 families; `contactus` file counts here are *pre*-migration —
Task 11 is that migration):

| Family | src | test | other | Total | Direct in `sneat-go`? | Contract exists? |
|---|---|---|---|---|---|---|
| contactus | 5 | 9 | 2 | 16 | yes | yes |
| competios | 9 | 4 | 4 | 17 | yes | **no** |
| sportius | 7 | 2 | 2 | 11 | yes | yes |
| bookius | 6 | 2 | 2 | 10 | yes | yes |
| chessraiders | 3 | 3 | 2 | 8 | yes | **no** |
| eventius | 4 | 2 | 2 | 8 | yes | yes |
| calendarius | 2 | 0 | 2 | 4 | yes | yes |
| gameboard | 0 | 0 | 2 | 2 | no (indirect only) | yes |
| kids-club | 0 | 0 | 2 | 2 | no (indirect only) | yes — **new finding, absent from the original stale hint entirely** |
| listus | 0 | 0 | 2 | 2 | no (indirect only) | yes |
| sneat-team | 0 | 0 | 2 | 2 | no (indirect only) | yes |
| trackus | 0 | 0 | 2 | 2 | no (indirect only) | yes |

**`sneat-bots`** (7 families, measured against `v0.30.7`/`6ce22b4`, i.e.
*after* PR #144):

| Family | src | test | other | Total | Direct in `sneat-bots`? | Contract exists? |
|---|---|---|---|---|---|---|
| trackus | 4 | 4 | 3 | 11 | yes | yes |
| bookius | 3 | 1 | 3 | 7 | yes | yes |
| calendarius | 2 | 3 | 3 | 8 | yes | yes |
| eventius | 3 | 1 | 3 | 7 | yes | yes |
| listus | 1 | 3 | 2 | 6 | yes | yes |
| gameboard | 0 | 0 | 2 | 2 | no (indirect only) | yes |
| contactus | 0 | 0 | 1 (`go.mod` only, `// indirect`) | 1 | **no — PR #144 fully cleared `sneat-bots`' own source** | yes |

`contactus` in `sneat-bots` collapsed from the original stale hint of 37 to 1
(a single indirect `go.mod` line) because of PR #144 — but that one line is
real, confirmed live via `go mod why -m` chaining through
`extensions/gametable/cmds4gametable → gametable/backend → calendarius/backend/facade4calendarius → ext-contactus/backend/dal4contactus`.
`sneat-bots`' own code is genuinely done for contactus; the residual
indirect entry is entirely `sneat-co/calendarius`'s (the product repo) to
clear.

## Tasks

### Task 1: Re-audit `sneat-go` after `sneat-go-storedalgo-hotfix` merges

**Id:** task-1
**Verifies:** re-run this Plan's inventory commands (`grep -rn
"sneat-co/ext-<family>"` per family, plus the `go mod why -m`/`go mod graph`
chains cited in Approach) against `sneat-go@origin/main` once that branch
merges, and update the file counts above and Tasks 7–11 before starting
them if they moved.
**Depends-On:** —
**Status:** planning

That branch migrates roughly ten old-contract sites plus a large pin set;
this Plan's `sneat-go` inventory predates it by construction (the branch was
never read to write this Plan). Every `sneat-go` task below gates on this
one, directly or transitively, so nobody duplicates work that branch already
did.

### Task 2: `sneat-bots` — migrate the `convosetup` composition hub (calendarius/contactus/listus/trackus)

**Id:** task-2
**Verifies:** `go build ./... && go test ./...` green in `sneat-bots`;
zero remaining `sneat-co/ext-calendarius`, `ext-listus`, `ext-trackus`
imports in `extensions/anybot/convo/convosetup/setup.go`,
`extensions/listus/convoactions/`, `extensions/trackus/convoactions/`; a
concrete finding recorded for each of `listus/backend` and `trackus/backend`
stating the exact release at which each bundles its own contract switch
(mirroring the confirmed `calendarius/backend@v0.6.2` finding), and whether a
partial path exists for either.
**Depends-On:** —
**Status:** planning

**This is already claimed by a separately dispatched lane as of this Plan's
writing — verify its live status before starting any part of this task.**
Scope: bump `calendarius/backend` to `≥ v0.6.2`, swap
`convosetup/setup.go`'s `Services.Calendarius` field type and
`NewCatalogRegistry` wiring to `sneat-ext-contracts/calendarius`, and
determine (don't assume) whether `listus/backend`'s and `trackus/backend`'s
own `botapi` packages need an equivalent bundled bump before
`extensions/listus/convoactions` and `extensions/trackus/convoactions` can
swap their `ext-listus`/`ext-trackus` imports. Contactus's side of this file
is already migrated (PR #144); this task's contactus-adjacent work is
limited to confirming the calendarius bump doesn't reintroduce an
old-contactus type via `calendarius/backend`'s internal
`convoservice4calendarius.New(...)` construction.

### Task 3: `sneat-bots` — migrate `bookius` vendorbot off `ext-bookius/backend`

**Id:** task-3
**Verifies:** `go build ./... && go test ./...` green; zero
`sneat-co/ext-bookius` references outside `go.mod`/`go.sum` in
`extensions/bookius/vendorbot/`.
**Depends-On:** —
**Status:** planning

No evidence of composition-threading through `convosetup` or another
product repo's return type for this family in `sneat-bots` — treat as a
genuinely separable single-package swap unless verification surfaces
otherwise.

### Task 4: `sneat-bots` — migrate `eventius` bot off `ext-eventius/backend`

**Id:** task-4
**Verifies:** `go build ./... && go test ./...` green; zero
`sneat-co/ext-eventius` references outside `go.mod`/`go.sum` in
`extensions/eventius/bot/`.
**Depends-On:** —
**Status:** planning

`go mod graph` shows this family's *only* `sneat-bots`-side requirer is
`sneat-bots` itself — `eventius/backend` (the product repo) is not even a
`sneat-bots` dependency today. This is the cheapest, most independent win in
the whole `sneat-bots` track.

### Task 5: `sneat-bots` — tidy and confirm disposition after Tasks 2–4

**Id:** task-5
**Verifies:** `go mod tidy` produces a diff limited to the modules touched
by Tasks 2–4; `go mod why -m github.com/sneat-co/ext-gameboard/backend` and
`...ext-contactus/backend` re-run and their chains recorded verbatim in the
task's completion note.
**Depends-On:** 2, 3, 4
**Status:** planning

Housekeeping after the three parallel family fixes land: confirms `tidy`
didn't touch anything outside their scope and records the current
`go mod why` chain for the two families (`gameboard`, `contactus`) that are
expected to remain `// indirect` regardless.

### Task 6: `sneat-bots` — cut and push a lightweight tag for the batch

**Id:** task-6
**Verifies:** `git tag vX.Y.Z && git push origin vX.Y.Z`; confirm via `git
log -1 vX.Y.Z` that it points at the post-Task-5 commit. Manual step —
`strongo_workflow` runs `disable-version-bumping: true` in this repo, exactly
as it did for today's `v0.30.7`.
**Depends-On:** 5
**Status:** planning

The only way any of Tasks 2–4's fixes become visible to `sneat-go` at all —
without this, they sit merged-but-inert on `sneat-bots` main indefinitely.

### Task 7: `sneat-go` — migrate `calendarius` host code

**Id:** task-7
**Verifies:** `go build ./... && go test ./...` green; zero
`sneat-co/ext-calendarius` references outside `go.mod`/`go.sum`.
**Depends-On:** 1, 2
**Status:** planning

Depends on Task 2, not just the reverse-gate: `sneat-go` imports
`sneat-bots/extensions/anybot/convo/convosetup` directly, so this task's
build depends on that package's post-migration public surface compiling,
even though `sneat-go` does not currently construct a concrete
`Services.Calendarius` value (lower runtime risk than `sneat-cli`'s, per
Approach).

### Task 8: `sneat-go` — migrate `sportius` host code

**Id:** task-8
**Verifies:** `go build ./... && go test ./...` green; zero
`sneat-co/ext-sportius` references outside `go.mod`/`go.sum` and outside
`pkg/modules/competios/`.
**Depends-On:** 1
**Status:** planning

Explicitly excludes `pkg/modules/competios/production_composition.go`,
which also references sportius but stays entangled with competios's
unresolved contract status (see Open Questions). Independent of
`sneat-bots` — `sportius` does not appear in `sneat-bots`' graph at all.

### Task 9: `sneat-go` — migrate `bookius` host code

**Id:** task-9
**Verifies:** `go build ./... && go test ./...` green; zero
`sneat-co/ext-bookius` references outside `go.mod`/`go.sum` and outside
`pkg/modules/competios/eventreg/`; the reverse-integration check
(`sneat-bots/scripts/pre-push-sneat-go-e2e.sh`, run manually if needed)
passes against Task 3's branch.
**Depends-On:** 1, 3
**Status:** planning

Explicitly excludes `pkg/modules/competios/eventreg/bookius_ports.go` and
`bookius_registrar.go` — same competios-entanglement reasoning as Task 8.
`sneat-go` imports `sneat-bots/extensions/bookius/vendorbot` directly, hence
the dependency on Task 3 rather than treating this as independent.

### Task 10: `sneat-go` — migrate `eventius` host code

**Id:** task-10
**Verifies:** `go build ./... && go test ./...` green; zero
`sneat-co/ext-eventius` references outside `go.mod`/`go.sum` and outside
`pkg/modules/competios/eventreg/`; reverse-integration check passes against
Task 4's branch.
**Depends-On:** 1, 4
**Status:** planning

Explicitly excludes `pkg/modules/competios/eventreg/eventius_attendance_projector.go`.
`sneat-go` imports `sneat-bots/extensions/eventius/bot/cmds4eventiusbot`
directly, hence the dependency on Task 4.

### Task 11: `sneat-go` — migrate `contactus` host code

**Id:** task-11
**Verifies:** `go build ./... && go test ./...` green; zero
`sneat-co/ext-contactus` references outside `go.mod`/`go.sum` in
`pkg/modules/contactus/`, `pkg/modules/ovdb/`, `pkg/modules/sportius/`, and
their test files.
**Depends-On:** 1, 2
**Status:** planning

`sneat-bots`' own contactus work is already merged (PR #144); this task is
`sneat-go`'s own local migration of its 5 source + 9 test files. Depends on
Task 2 (not Task 2's contactus portion specifically, but the same
`convosetup` file `sneat-go` imports) rather than being called independent
just because `sneat-bots`' contactus half looks finished.

### Task 12: `sneat-go` — bump the `sneat-bots` pin and isolate the go1.27 toolchain jump

**Id:** task-12
**Verifies:** `go.mod` line 110 updated to Task 6's tag; `go mod tidy`
clean; this lands as its **own commit/PR**, separate from Tasks 7–11, so any
coverage-floor failures in packages untouched by this migration can be
triaged as the expected go1.26→go1.27 re-measurement artifact (~3.5pp lower,
observed elsewhere in this fleet today) rather than blamed on unrelated
work. No floor is lowered as part of this task — that requires separate
founder sign-off.
**Depends-On:** 6
**Status:** planning

### Task 13: `sneat-go` — post-bump final disposition audit

**Id:** task-13
**Verifies:** for each of `gameboard`, `kids-club`, `listus`, `sneat-team`,
`trackus`: `go mod why -m github.com/sneat-co/ext-<family>/backend`
re-run; record either "removed" or "blocked on `<named repo>`" per family in
the task's completion note.
**Depends-On:** 7, 8, 9, 10, 11, 12
**Status:** planning

The single checkpoint that turns this Plan's `sneat-go` scope from "tasks
executed" into "disposition documented" — the honest terminal state this
Plan can claim, per Summary.

### Task 14: `sneat-bots` — post-tidy final disposition audit

**Id:** task-14
**Verifies:** `go mod why -m github.com/sneat-co/ext-contactus/backend` and
`...ext-gameboard/backend` re-run against the state after Task 5; record
"removed" or "blocked on `<named repo>`" for each.
**Depends-On:** 5
**Status:** planning

`sneat-bots`' mirror of Task 13 — its own scope is smaller (two families,
`contactus` and `gameboard`, both already known-indirect going in) but the
same discipline applies: document, don't assume.

## Open Questions

1. **contactus — fleet-wide transitive-forcing blocker.** `go mod graph`
   shows `ext-contactus/backend` still required (directly or transitively)
   by at least 15 other repos' pinned versions as consumed today: `assetus`,
   `calendarius`, `chessraiders`, `competios`, `contactus` (the product repo
   itself), `debtus`, `eventius`, `gametable`, `logistus`, `ovdb`,
   `requoter`, `rosycycle`, `sneat-core-modules`, `sneat-work`, `togethered`,
   `trackus`. `debtus v0.2.33` and `sneat-core-modules v0.65.14` were each
   independently confirmed today to force it. Neither `sneat-go`'s nor
   `sneat-bots`' `go.mod` can drop `ext-contactus/backend` to zero until
   every one of these repos migrates and releases — a campaign well beyond
   two repos. Who owns coordinating that fleet-wide sweep, and is it
   explicitly out of scope for this Plan (recommendation: yes, but the
   founder should say so)?
2. **calendarius — fleet-wide blocker, same shape.** Requirers today:
   `calendarius` (itself), `communitycentrum`, `gameboard`, `gametable`,
   `requoter`, `togethered`, plus the two host repos. Same question as #1.
3. **gameboard — single external blocker, closes two host repos at once.**
   The *sole* remaining requirer in both `sneat-go` and `sneat-bots` is
   `sneat-co/gameboard`'s own `gameboard/eventtimeline` package (confirmed
   via `go mod why`). Fixing it there and bumping the pin in both hosts
   clears this family everywhere with zero host-repo source changes. Who
   owns filing/running that fix in `sneat-co/gameboard`?
4. **kids-club — single external blocker, `sneat-go`-only.** Sole requirer:
   `sneat-co/kids-club`'s own `kidsclub` package. This family never appeared
   in the original stale-hint inventory at all — it is a genuinely new
   finding from today's audit. Same ownership question as #3, narrower
   blast radius.
5. **sneat-team — single external blocker, `sneat-go`-only.** Sole requirer:
   `sneat-co/sneat-team`'s own `team` package. Same ownership question.
6. **bookius / togethered.** Beyond the host-repo tasks above, `bookius` and
   `togethered` (both product repos) still require `ext-bookius/backend`.
   Who tracks their migration once Tasks 3 and 9 land?
7. **eventius / sportius product repos — who owns tracking them?**
   `sneat-co/eventius` and `sneat-co/sportius` still require their
   respective `ext-*` modules independently of the host-repo tasks above
   (Tasks 4, 8, 10). Who owns filing/tracking their migration once those
   land, mirroring Open Questions 3–6?
8. **competios and chessraiders have no contract home at all.** Neither
   `competios` nor `chessraiders` appears in
   `sneat-ext-contracts/contracts.json`, `go.work`, or the `libs/`/top-level
   family directories — confirmed by direct inspection, not inference. Both
   were archived today alongside the other 28. Before Tasks equivalent to
   7–11 can even be written for these two families (17 + 8 files in
   `sneat-go` respectively), someone must decide: (a) who authors
   `sneat-ext-contracts/competios` and `.../chessraiders`, following the
   same `contract4<family>` shape already used everywhere else, and (b) for
   chessraiders specifically, whether that authoring and the subsequent
   `sneat-go` migration should route through the existing chess-plans
   coordination flow (plan walks only on the chess plans-branch worktree,
   CLI-enforced) instead of this Plan. I am not choosing an answer to either
   half.
9. **The reverse-integration gate is a local hook, not a CI check.**
   `sneat-bots/scripts/pre-push-sneat-go-e2e.sh` only runs on a machine that
   has the `wb`-managed hook installed and does not pass `--no-verify`.
   Should this campaign additionally wire an equivalent check into
   `sneat-bots`' GitHub Actions CI so the cross-repo safety net can't be
   silently skipped by a differently-configured machine or agent? This is a
   policy decision beyond finishing the cutover — surfaced, not decided.
10. **`sneat-cli` is a dependent of Task 2, not a parallel track — but it is
    a third host repo this Plan has no authority over.** Once Task 2 lands,
    `sneat-cli`'s `convo_trackus.go` composition should become a clean 2-file
    + 1-dependency-bump swap against the migrated types. Should a
    corresponding Plan task exist in `sneat-cli`'s own repo, and who owns
    writing it?
11. **contactus's terminal state in `sneat-bots` may already be
    "acceptable."** `sneat-bots`' own source is fully migrated (PR #144);
    the sole remaining `go.mod` line is `// indirect`, blocked entirely on
    `sneat-co/calendarius` (via `gametable`). Is "indirect and externally
    blocked, with the blocker named" an acceptable terminal state for this
    Plan's contactus scope in `sneat-bots`, given Task 14 already documents
    it — or does the founder want an issue opened against
    `sneat-co/calendarius` as part of this effort?

---
*This document follows the https://specscore.md/plan-specification*
