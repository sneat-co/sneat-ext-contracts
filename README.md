# sneat-ext-contracts

Every `@sneat/extension-<family>-contract` npm package, and every extension's
Go contract module, lives here as one sibling in a single repo — not as ~24
separate single-contract repos. npm package names never change; only sources
and publish pipelines move here, family by family. Runtimes, UIs, and backends
stay in their product repos: this repo holds contracts only, and its name
says so.

Full rationale, requirements, and the phased migration plan:
`spec/features/ext-contracts-monorepo/README.md` on `sneat-co/sneat-libs`.

## Layout

```text
libs/<family>/            @sneat/extension-<family>-contract (Nx project "<family>-contract")
<family>/go.mod           github.com/sneat-co/sneat-ext-contracts/<family>, tagged <family>/v<version>
<family>/version.go       Go-only families ONLY (no libs/<family>/): const Version = "x.y.z" is
                          the version source (see "Go-only family versioning" below)
go.work                   lists every "<family>/go.mod" that exists (empty at Phase 0)
contracts.json            manifest of families this repo owns/publishes (empty at Phase 0)
docs/boundaries.md        the module-boundary map + the contract ownership test
tools/tier-coherence/     CI check: every owned contract's npm-latest installs & type-checks together
.github/workflows/ci.yml       lint/test/build (npm side) + per-module Go build/test (auto-discovered)
.github/workflows/publish.yml  version-plan-driven release, npm first, Go tags to follow (see below)
```

Nothing here is a "frontend workspace nested in the repo" the way the
single-contract `ext-<name>` repos and their `sneat-ext-contract-template`
were shaped — the pnpm/Nx workspace **is** the repo root, because this repo
holds nothing but contracts.

## The contract-author journey

Verbatim from `spec/features/ext-contracts-monorepo/README.md`
(`sneat-co/sneat-libs`), Behavior → End-to-end journeys:

> **Contract author (usually an AI agent):** clones one repo, edits one
> contract lib, adds an Nx version plan naming that project and its bump,
> opens one PR. CI enforces boundary rules, peer ranges, and tier coherence on
> the PR. On merge, the release pipeline versions and publishes ONLY the
> projects named in consumed version plans. *Good result: npm shows the new
> version of exactly that contract, minutes after merge, with correct
> provenance; no other contract republished.*

In practice, that is:

1. Clone this repo. Edit one contract lib under `libs/<family>/`.
2. Add an **Nx version plan** naming that project and its bump
   (`pnpm nx release plan`, or hand-author a file under `.nx/version-plans/`)
   — this is what makes `nx release` touch that project and only that
   project. No version plan, no version bump, no publish, no matter what else
   changed in the PR.
3. Open one PR. CI enforces, on every PR, the module-boundary rules
   (`docs/boundaries.md`), peer-range strictness, the zone.js guard, and the
   tier-coherence check (every owned contract's npm-latest still installs and
   type-checks together).
4. On merge to `main`, the release pipeline (`publish.yml`) runs
   `nx release --yes` in **independent** (per-project) mode: it versions,
   changelogs, and publishes **only** the projects named in the version plans
   that were merged. Every other contract in the repo — however many there
   are — is untouched. No other npm package republishes because your PR
   touched a shared file.

Good result: npm shows the new version of exactly the contract you changed,
minutes after merge, with correct provenance. Nothing else on npm moves.

### Lockstep family versioning (npm + Go)

A family's Go module (`<family>/go.mod`) tracks the **same** version number
as its npm package — one version plan, one bump, applied to both artifacts.
This applies **only when the family has an npm sibling** (`libs/<family>/`):
it keeps "what version is `<family>` at" a single question instead of two
independently-drifting answers for a family that ships both artifacts. `nx
release` publishes the npm side; `publish.yml` then tags the Go module at the
same version (see "Go module tagging" below).

### Go-only family versioning (no npm sibling)

Some families are pure Go — `competios` and `chessraiders` today — with no
Angular/TS consumer and therefore no `libs/<family>/`. Per the founder's
ruling ("I'd prefer if each contract get versioned independently"),
manufacturing an unused npm shell just to source a version number for one of
these would itself be the coupling the ruling rejects — a Go-only family's
version line is its own, and does not ride any npm sibling, any other
family's cadence, or a repo-wide version.

Its version instead comes from a plain `<family>/version.go` at the module
root — the same `const Version = "x.y.z"` convention already used across the
wider Sneat/dalgo Go ecosystem (e.g. `dal-go/dalgo/version.go`,
`bots-go-framework/bots-api-telegram/version.go`), not a new one invented for
this repo. The journey:

1. Clone this repo. Edit the contract under `<family>/`.
2. Bump `<family>/version.go`'s `Version` constant to the next version this
   change should ship as, in the **same PR** — this is the Go-only
   equivalent of adding an Nx version plan; there is no separate plan file
   because there is no Nx project to plan against (Nx only ever discovers
   `libs/*/package.json`; a Go-only family's directory is invisible to it).
3. Open one PR. CI builds/tests/vets the module same as any other Go family
   (`.github/workflows/ci.yml`'s `discover-go`/`go` jobs auto-discover every
   `<family>/go.mod`).
4. On merge to `main`, `publish.yml`'s `release-on-main` job reads
   `<family>/version.go`, and if `<family>/v<Version>` isn't already tagged,
   tags it — no npm release, no version plan, no manual dispatch needed. This
   fires on every CI-success-on-main run, so it also self-heals: if a tag
   attempt ever fails, the very next run tries again as long as the tag is
   still missing.

Good result: pushing a `version.go` bump is the entire release act for a
Go-only family — a new Go tag appears minutes after merge, with nothing on
npm ever touched.

**Silence is impossible by design.** If a `<family>/go.mod` exists with
neither a `libs/<family>/package.json` (npm-sibling path) nor a
`<family>/version.go` (Go-only path), or `version.go` exists but its
`Version` line can't be parsed, `publish.yml` fails the run (`::error::` +
non-zero exit) instead of skipping — a merged, unversionable family must
never pass silently behind a green Publish run. (That silent skip is exactly
how `competios` and `chessraiders` merged and went untagged for a day before
this section and the pipeline fix that backs it were added — recovered via
the `go_module_tags` manual escape hatch documented below.)

### 0.x semantics — read this before picking a version-plan bump

Every contract here starts on the `0.x` train. Nx's conventional-commits-style
version resolver **downshifts** on `0.x`, same as npm/semver convention:

| Version plan says | Effective bump on 0.x | Example |
| --- | --- | --- |
| `patch` | patch | `0.4.2` → `0.4.3` |
| `minor` | **patch-level** | `0.4.2` → `0.4.3` |
| `major` | **minor-level** | `0.4.2` → `0.5.0` |

There is no way to get an actual `1.0.0`-style major bump on the `0.x` train
via a version plan — that only happens once a family is deliberately promoted
past `0.x` (a separate, explicit decision per family, not an accident of
picking "major" in a version plan). If you want your contract's next version
to actually look different from a patch, remember the ceiling: "major" gets
you a minor bump, nothing more, until the family leaves `0.x`.

### Tags stay disambiguated: npm vs. Go

- npm release tags: `<family>-contract-v<version>` (Nx `release.releaseTag.
  pattern: "{projectName}-v{version}"` in `nx.json`, since the Nx project name
  for every npm contract is `<family>-contract`).
- Go module tags: `<family>/v<version>` — **reserved for Go modules only**.
  Never reuse this pattern for an npm release tag, and never expect an
  npm-side tag to satisfy a Go module's version resolution.

## The consumer journey

Nothing changes at adoption time. Package names and the npm registry location
are identical to what they were in the old per-extension `ext-<name>` repo.
`pnpm up @sneat/extension-<family>-contract` works before, during, and after
that family's source migrates into this repo — no consumer-side edit, ever.

## The migration-operator journey (per family)

This repo ships with **no family migrated** (Phase 0 — foundation only; see
the spec's migration plan for phase order and batching). Migrating a family
in later phases means:

1. **Import the source** into `libs/<family>/` (and `<family>/go.mod` if/when
   the Go leg is armed) with a provenance note (where it came from, at what
   commit/version).
2. **Prove API parity** against npm's actual latest for that family — not
   against the old repo's git tags. This is deliberate (REQ
   `version-continuity-from-npm`): several families' git tags and npm latest
   have already drifted apart, and npm reality wins.
3. **Publish the next version from this repo** (a version plan + merge to
   `main`, same as any other contract-author change).
4. **Add the family to `contracts.json`** so the tier-coherence check starts
   covering it, and to `docs/boundaries.md` if it needs any cross-contract
   edge (default is none — see that doc).
5. **Disarm the old pipeline and archive the old `ext-<family>` repo.** Point
   its description at this repo. An org-wide audit must find no other repo
   with an armed workflow that can still publish `@sneat/extension-<family>-
   contract`.

Good result: npm latest for that family is published from this repo, the old
repo is archived and cannot publish, and the org-wide audit finds no second
armed publisher for that npm name.

### Go module tagging

`publish.yml` tags a family's Go module via `sneat-co/cicd`'s reusable
`.github/workflows/go-module-tags.yml` (input contract: `modules: JSON array
of {dir, version}` + optional `ref`), armed since 2026-08-26. Two independent
resolution paths feed it {dir, version} pairs, one per family shape — see
"Lockstep family versioning" and "Go-only family versioning" above:

- npm-sibling family → version from `libs/<family>/package.json`, gated on
  that family's npm release tag landing on `HEAD` this run.
- Go-only family → version from `<family>/version.go`, gated on
  `<family>/v<version>` not already existing as a tag.
- Neither source present for an on-disk `<family>/go.mod` → the run fails
  loudly rather than skipping (see "Go-only family versioning").

A rare manual escape hatch also exists: `workflow_dispatch` with `dry_run:
false` and an explicit `go_module_tags` JSON input forwards {dir, version}
pairs straight to `go-module-tags.yml`, bypassing both resolution paths above
— use it only when the normal resolution can no longer see what it needs
(e.g. a family's release commit is no longer `HEAD` by the time you
re-dispatch), never as the routine way to cut a release.

## Module boundaries

Default: **no contract imports another contract.** See `docs/boundaries.md`
for the full rule, the allow-list mechanism, and the ownership test that
decides what belongs in a contract at all.
