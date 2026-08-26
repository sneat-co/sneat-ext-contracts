# @sneat/extension-gameboard-contract

The published TS contract for the gameboard extension — the deterministic
event-timeline fold reducer, legacy list DTOs/contexts/service tokens, and the
legacy game-record shape. Hand-implemented against the frozen wire contract
(`api4gameboard.tsp`, kept in `sneat-co/ext-gameboard`) to mirror the Go
reducer in `gameboard/eventtimeline` (this repo). It depends only on
foundational/core (`@angular/core`, `rxjs`, `@sneat/{core,data,dto,space-models}`)
— never on another `@sneat/extension-*` contract.

## Provenance

Migrated from `sneat-co/ext-gameboard` (`frontend/`), commit
`32ddd0433f2e2f0e593ea55ed0207390f8526770` (`origin/main`, 2026-08-26 read).

**npm reality wins (REQ `version-continuity-from-npm`):** npm's actual latest
is `0.3.0` (verified via `npm view @sneat/extension-gameboard-contract
dist-tags`). Disk version here is seeded at `0.3.0` to match; the paired
version plan (`.nx/version-plans/`) requests the patch bump to `0.3.1`.

**`.d.ts` parity gate:** built `frontend/` at commit `9c64da9` (content-identical
to `32ddd04` for `backend/`+`frontend/` — the only diff between the two is an
unrelated `.github/renovate.json` edit) in an isolated worktree, ran
`tsc -p tsconfig.build.json`, and diffed the resulting `dist/**/*.d.ts` against
`npm pack @sneat/extension-gameboard-contract@0.3.0`'s unpacked `dist/**/*.d.ts`
— byte-identical (`diff -r` exit 0, no output). `npm pack` of `0.2.0` vs `0.3.0`
also confirmed those two published versions are code-identical (only
`package.json`'s `version` field differs) — npm minted two versions in the same
CI run without any git-tracked content change between them.

That byte-identical `.d.ts` set proves the SOURCE content is unchanged, but
this repo builds every contract lib with `@nx/angular:package` (ng-packagr),
not `tsc` — matching taxus and every future family, not `ext-gameboard`'s own
`tsc -p tsconfig.build.json`. ng-packagr rolls every per-file `.d.ts` into one
bundle (`dist/libs/gameboard/types/sneat-extension-gameboard-contract.d.ts`)
with a different declaration style (`declare function` instead of `export
function`, one trailing `export { ... }` list instead of per-file `export`).
So the build-format change was checked separately, symbol-for-symbol: every
exported value and type name in npm 0.3.0's unpacked `dist/**/*.d.ts` (57
names, extracted by regex) has a same-named, same-shape counterpart in the
`nx build gameboard-contract` output's rolled-up `.d.ts` — none missing, none
added, none renamed. The two checks together (byte-identical source content +
symbol-for-symbol bundled output) are the full parity gate for this lib.

**Excluded: an untagged, unpublished post-0.3.0 commit.** `origin/main` HEAD
(`7af1daf`) contains one further commit beyond `32ddd04` that touches this
contract: `8b19343` "feat: add linked competition contract" (merged via PR #1,
dated 2026-08-15 — a month after the 0.3.0 publish), which adds
`LinkedCompetition`/`SubmissionSync`/`ExternalControlAuthorizer` to both
`backend/linkage.go` and `frontend/src/linked-competition.ts`. Per the
migration rule (only carry a post-latest change forward if an unpublished
**tag** backs it, as the taxus migration did for its phantom `v0.0.3`), this
was checked and excluded: `git tag`/`gh api .../tags`/`gh release list` on
`sneat-co/ext-gameboard` show zero tags or releases for the npm/Go contract
side at all (only `backend/v0.1.0`, `v0.2.0`, `v0.2.1`, none of which reach
`8b19343`) — there is no tag evidence this commit was ever cut as an intended
release, unlike taxus's `v0.0.3`. It is therefore left out of this migration.
Also excluded for the same reason: `7af1daf`'s unrelated `backend/go.mod`
toolchain bump (`go 1.25.0` → `go 1.27.0`), which sits on top of the same
untagged commit range. **Flagging for the migration owner:** whoever owns
`sneat-co/ext-gameboard` should decide whether the linked-competition contract
gets its own follow-up migration (with its own version plan) once it has
release evidence, or whether it needs one retroactively cut in the source repo
first.

**Legacy source dead-code trim:** two files (`legacy/dto/list.ts`,
`legacy/services/interfaces.ts`) each dropped one block of already-commented-out
dead code during the copy (stale alternate `IListService`/`IGameboardService`
abstract-class sketches). No `.d.ts`-visible symbol changed; this is a
non-functional cleanup, not a content migration decision.

**Test target — and a repo-wide naming gap found along the way:** unlike
`taxus-contract` (no tests yet), this family ships real runtime logic and
tests (the fold reducer + its Go↔TS parity oracle, plus a legacy DTO unit
spec) — `frontend/src/eventtimeline.parity.test.ts` and
`frontend/src/legacy/dto/list.spec.ts` in the source repo, both ported
verbatim (the parity test's explicit `vitest` imports were already present;
`list.spec.ts` gained explicit `describe`/`it`/`expect` imports from `vitest`
since this workspace does not enable Vitest globals). Added `vite.config.ts`
so this project has a vitest config at all. **Finding:** `nx.json` registers
*two* plugins for test inference — `@nx/vite` with `testTargetName: "test"`
and `@nx/vitest` with `testTargetName: "vite:test"` — but in the installed Nx
version (23.1.1) `@nx/vite`'s executor set no longer includes a test executor
(only `dev-server`/`build`/`preview-server`; confirmed via
`node_modules/@nx/vite/executors.json`), so only `@nx/vitest` actually infers
anything, and only under the name **`vite:test`**, never `test`. Root
`package.json`'s `"test": "nx run-many -t test --all"` script (and CI's
`nx-ci.yml` call with `targets: "lint test build"`) both target literally
`test` — so without a fix, gameboard-contract's tests would be silently
skipped by both. Since `nx.json` is shared across every family and out of
scope for this branch, the fix lives entirely in this project's own
`project.json`: an explicit `test` target (`nx:run-commands`, `vitest run`)
alongside the harmless, redundant inferred `vite:test` one. **This is worth
the merger/repo-owner fixing at the `nx.json` layer** (either drop the dead
`@nx/vite` plugin entry or rename `@nx/vitest`'s `testTargetName` to `"test"`)
so future family libs with tests don't need the same per-project workaround.

**Shared parity fixture relocated one level:** the Go↔TS oracle
(`parity/parity.json` in the source repo, one level above both `backend/` and
`frontend/`) now lives at `gameboard/parity/parity.json` in this repo, since
the Go module directory (`gameboard/`, a sibling of `libs/`) plays the role the
old repo root played. Both language-side references were updated accordingly:
`gameboard/eventtimeline/parity_test.go`'s `parityPath` (two `..` segments →
one) and this lib's `eventtimeline.parity.test.ts` fixture URL (two `..`
segments → three, since `libs/gameboard/src/` is one level deeper than the old
`frontend/src/`). Fixture content is byte-identical to the source; the Go test
(`TestParityFixture`) re-derives it from `Fold()` and asserts equality, so this
was cross-checked at migration time, not just copied on faith.

**Old CI boundary check (`scripts/check-no-extension-deps.sh`) — intent
carried forward, script not ported:** the source repo enforced "gameboard
depends only on foundational/core, never another extension" via a hand-rolled
bash script (`go list -deps` for the Go side, a `package.json` grep for the TS
side) run in its own CI. This monorepo replaces that per-repo script with a
single mechanical boundary lint shared by every family:
`@nx/enforce-module-boundaries` in the workspace root `eslint.config.mjs`,
driven by this project's `family:gameboard` / `layer:contract` tags (see
`docs/boundaries.md`). The same invariant — no cross-contract import, contract
never depends on runtime/app code — is what that lint enforces; no
family-specific script is needed or added here.
