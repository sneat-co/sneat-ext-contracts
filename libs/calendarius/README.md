# @sneat/extension-calendarius-contract

The published TS contract for the calendarius extension — happening/event
DTOs (including the `EventHappening` port: I-prefixed DTOs,
`RepeatPeriod`-by-reference recurrence), the happening context, weekday/slot
view-models, and the schedule-nav service token. Depends only on
foundational/core (`@sneat/core`, `@sneat/dto`, `@sneat/space-models`) plus
Angular — zero `@sneat/extension-*` imports anywhere in `src/` (verified by
grep before the copy).

## Provenance

**Source: `sneat-co/sneat-libs`, `libs/extensions/calendarius/contract`,
commit `344831c0fc91b25ea01a21b4d1e276d3a589d513`** ("chore(release): publish
0.27.2", 2026-08-26 — the exact commit that produced the npm version this
migration targets; `origin/main` HEAD `bddf332` is one docs-only merge ahead
and does not touch this path). sneat-libs is the current sole publisher of
this contract (no `sneat-co/ext-calendarius` frontend/contract exists — that
repo ships only the `backend/` Go module, see below). This branch does not
modify sneat-libs; removing `calendarius` contract from its own release set
is a separate step **owed to sneat-libs' owner**, not done here.

**npm reality:** `npm view @sneat/extension-calendarius-contract versions`
shows `latest` moved from `0.27.1` (this migration's original target,
published 2026-08-26T08:23) to `0.27.2` (published 2026-08-26T09:08, mid-task)
— `npm pack` of both and `diff -r` confirms they are **content-identical**
(only `package.json`'s `version` field differs). Disk here is seeded at
`0.27.2` to track true npm latest.

**Full byte parity gate:** built fresh with this repo's own toolchain
(`nx build calendarius-contract`) and diffed the output against
`npm pack @sneat/extension-calendarius-contract@0.27.2`'s unpacked tarball:
`types/*.d.ts` byte-identical, `fesm2022/*.mjs` byte-identical. **const-enum
grep:** zero `const enum` usages in `src/` — no runtime-erasure risk.
`package.json`'s declared surface (`peerDependencies`: only `@angular/common`
and `@angular/core`; `dependencies`: `tslib`) is reproduced exactly as
published — this contract's own `package.json` does not declare
`@sneat/core`/`@sneat/dto`/`@sneat/space-models` as peers even though `src/`
imports all three (a pre-existing gap in the source's own published metadata,
carried forward unchanged rather than "fixed," since real npm metadata is
what byte parity means here; those three packages are already root workspace
devDependencies in this repo, so the build resolves them regardless).

## Extraction changes

**Flattened `src/lib/` → `src/`** (this repo's convention since `libs/taxus`).
Zero import-path risk: every relative import inside `src/lib/` goes at most
one level up (`../dto`, checked before the move), entirely within the moved
subtree, so only the top `src/index.ts`'s `./lib/X` prefixes needed rewriting
to `./X`.

**`test-setup.ts` dropped, not ported.** The source's
`src/test-setup.ts` calls `setupTestEnvironment()` from `@sneat/core/testing`
— a workspace-only tsconfig path, explicitly **not** part of the published
`@sneat/core` npm package (`libs/core/tsconfig.lib.json` excludes
`src/lib/testing/**` in sneat-libs; its own docstring says so). Read that
function's implementation in sneat-libs: it installs a zoneless Angular
`TestBed` environment plus Ionic/Stencil DOM shims. None of this contract's
five spec files touch Angular DI, `TestBed`, or the DOM at all (verified:
no `@angular/core/testing`, no `window`/`document` reference anywhere in
`src/`) — they are pure DTO/logic tests. Porting a stub replacement for an
unused harness would add risk (an import that can silently drift out of sync
with sneat-libs) for zero benefit, so `vite.config.ts` here declares no
`setupFiles` at all and `test-setup.ts` is not copied. `ext-contactus`'s own
already-extracted contract solved the analogous problem by inlining a
self-contained shim (see `libs/contactus/README.md`, same batch); calendarius
doesn't need that shim because nothing in its specs exercises the DOM/DI
surface it would shim.

**Explicit `vitest` imports added to four of five specs.** This workspace
enables no Vitest globals (established by `libs/taxus` onward); sneat-libs'
shared `vite.config.base.ts` sets `globals: true`. `dto/event-happening.spec.ts`
already imported `describe`/`expect`/`it` explicitly and needed no change.
The other four (`contexts/happening-context.spec.ts` — also gained
`beforeEach`, `dto/happening.spec.ts`, `dto/todo_move_funcs.spec.ts`,
`view-models.spec.ts`) gained the same explicit import — the identical fix
gameboard-contract's migration made for its one ported spec. Added
`vite.config.ts` (`environment: 'node'`) and an explicit `test` target in
`project.json` (`nx:run-commands`, `vitest run`) for the same `nx.json`
`@nx/vite`/`@nx/vitest` naming-gap reason as every other family lib with
tests in this repo.

## Go half

Sourced separately from `sneat-co/ext-calendarius` `backend/` — see
`../../calendarius/doc.go` and this repo's commit message for that half's own
provenance (tag `backend/v0.0.6`, the recurrence fix accepting
weekly/fortnightly/monthly/yearly). `sneat-co/ext-calendarius` has no frontend
contract of its own; the TS and Go halves of this family currently live in
two different upstream repos, unified here for the first time.

## Cross-family edges

None. Calendarius imports only `@angular/common`, `@angular/core`,
`@sneat/core`, `@sneat/dto`, `@sneat/space-models` — all foundational/core.

## Follow-up owed

`sneat-co/sneat-libs` should stop publishing
`@sneat/extension-calendarius-contract` from `libs/extensions/calendarius/contract`
once this monorepo's copy is the canonical source (a release-set change in
sneat-libs' own Nx release config) — not done in this branch, which never
touches sneat-libs.
