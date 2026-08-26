# @sneat/extension-contactus-contract

The published TS contract for the contactus extension — contact/person/member
DTOs, contact-group and contact-role types, the API request/response DTOs, the
Angular service tokens (`contact-service`, `contactus-space-service`,
`contact-nav-service`, `contactus-nav-service`, `contact-group-service`,
`contact-role-service`, `invite-service`), and the cross-extension reuse
surfaces (`contact-ref`, `contacts-selector.contract`,
`contacts-list.contract`). This is the fleet's largest cross-family
dependency target (81 of 148 measured cross-family contract imports point at
contactus per the 2026-08-26 fleet audit in `docs/boundaries.md`) — no
contract currently imports another (verified: zero `@sneat/extension-*`
imports in `src/`), only foundational/core plus `@sneat/auth-models`.

## Provenance

**Source: `sneat-co/ext-contactus`, `frontend/libs/extensions/contactus/contract`,
commit `1d703c36cd9a76bb07c54055e546edb9a44974dc`** (`origin/main`, read
2026-08-26). This is the sole viable source: `sneat-co/ext-contactus` owns the
armed contract publisher (`.github/workflows/publish.yml`, `nx release --yes`
on push to `main`); `sneat-co/contactus` (the product/app repo) was checked and
contains **no contract library at all** — its `frontend/package.json` and its
`ui`/`runtime` libs only *consume* `@sneat/extension-contactus-contract` (all
three pin `0.12.2` exact), matching the brief's note that fleet consumers pin
0.12.2 exact. There is no drift to adjudicate between two sources because only
one source exists; the product repo was a false lead.

**The "stuck 0.12.3" tag, explained (git-tag ≠ package.json ≠ npm, and that's
by design):** `ext-contactus`'s tags run `v0.12.2` → `v0.12.3`, but the
contract's checked-in `package.json` has read `"version": "0.12.1"` since the
one commit that ever touched it (`eabf67f`, bootstrap) — it has never been
bumped in git. This is not drift or breakage: the project's `project.json`
sets `"release": {"version": {"currentVersionResolver": "git-tag"}}`, so
`nx release` computes the published version from the next git tag at publish
time and never writes it back to disk. `git diff v0.12.2 v0.12.3` touches only
`backend/go.mod`/`go.sum`/`renovate.json` — **zero contract-source changes**
between the two tags — so `v0.12.3` was cut for the Go side only; npm's actual
latest for the *contract* package is `0.12.2` (confirmed via
`npm view … versions`: `0.12.0`, `0.12.1`, `0.12.2`, no `0.12.3` ever
published — nx release's `nx-release-publish` target only runs for projects
whose conventional-commit-scanned changes actually touch them, and none did).
Disk version here is seeded at `0.12.2` to match npm reality
(REQ `version-continuity-from-npm`, same rule the gameboard migration
applied); the paired version plan (`.nx/version-plans/`) requests the patch
bump to `0.12.3` (now free — the npm number, unlike the git tag, was never
spent on this package).

**`.d.ts` + `.mjs` parity gate:** confirmed `git diff v0.12.2 HEAD --
frontend/libs/extensions/contactus/contract/` is empty — the contract source
has not changed since the tag. Built it fresh at `HEAD` with this repo's own
toolchain (`pnpm install --no-frozen-lockfile --ignore-scripts` +
`nx build ext-contactus-contract`) and diffed the output against
`npm pack @sneat/extension-contactus-contract@0.12.2`'s unpacked tarball:
`types/*.d.ts` byte-identical, `fesm2022/*.mjs` byte-identical. Only
`package.json`'s `version` field differs (`0.12.1` on disk vs `0.12.2`
published — expected, see above). **const-enum grep:** zero `const enum`
usages anywhere in `src/` — no runtime-erasure risk for this migration mode.

**Go half — backend, highest tag `backend/v0.1.8`:** `HEAD`
(`1d703c3`, "build: use Go 1.26 with Go 1.27 toolchain") sits two untagged
commits ahead of `backend/v0.1.8` (`bc2350f`) — both Go-toolchain-only bumps
(`build: bump Go to 1.27.0`, then this Go-1.26+toolchain-1.27 commit), no
`.go` file changed between the tag and `HEAD`
(`git diff backend/v0.1.8 HEAD -- backend/` touches only `go.mod`). Consumers
(togethered, sneat-work, logistus) currently require `v0.1.7`; `v0.1.8` is one
untagged-content-identical release ahead of that. `GOWORK=off go build/vet/test
./...` all green at `HEAD` before the copy.

## Extraction changes

**Flattened `src/lib/` → `src/`** (Nx's default `nx g lib` scaffold nests
everything one level under `lib/`; this repo's convention, established by
`libs/taxus` and every migrated family since, keeps library source directly
under `src/`). Zero import-path risk: every relative import in the source
tree goes at most one level up (`../dto`, `../contexts`), entirely within the
moved subtree, so only the top `src/index.ts` needed its `./lib/X` prefixes
rewritten to `./X` — verified before the move (`grep` for `../../` inside
`src/lib/` found none) and after (full byte-diff above proves the built
output is unaffected — the public API surface is defined by `index.ts`'s
export graph, not by internal file layout).

**`@sneat/auth-models` — a new foundational peer for this repo.** Six files
(`contact-types.ts`, `member.ts`, `contact-requests.ts`, `person.ts`,
`contact-base.ts`, `apidto/requests.ts`) import from `@sneat/auth-models`,
which is not yet in this repo's root `package.json` devDependencies (every
prior family only needed `@sneat/core`/`data`/`dto`/`space-models`). It stays
a declared `peerDependency` here, matching the source's own (and npm's
published) `package.json` exactly — root `package.json` is a forbidden file
for this branch, and per the workspace's `auto-install-peers` (on by default,
documented in `.pnpmfile.cjs`), `pnpm install` at the repo root will pull it
in on its own once this lib exists; **the merger needs to run `pnpm install`
after combining batch-3 branches so `pnpm-lock.yaml` picks up
`@sneat/auth-models`** (this branch does not commit a lockfile change, per the
`pnpm-lock.yaml commits` forbidden-files rule).

**Test target:** one spec, `dto/contact-requests.spec.ts`, already imports
`describe`/`it`/`expect` explicitly from `vitest` (this workspace has no
globals) — ported verbatim, no changes needed. Added `vite.config.ts`
(`environment: 'node'`, no Angular plugin — nothing under test touches
Angular DI/TestBed) and an explicit `test` target in `project.json`
(`nx:run-commands`, `vitest run`), the same fix gameboard's migration
introduced for the same `nx.json` `@nx/vite`/`@nx/vitest` naming gap.
`src/test-setup.ts` (self-contained already in the source — it does not
depend on `@sneat/core/testing`, unlike calendarius's) was copied over but is
currently unreferenced (no `setupFiles` entry): the one spec needs no
DOM/Ionic shims. Kept in the tree for provenance and in case a future spec
does.

## Cross-family edges

None outgoing. Contactus imports only `@angular/core`, `@sneat/auth-models`,
`@sneat/core`, `@sneat/dto`, `@sneat/space-models`, `rxjs` — all
foundational/core. (It is a large *incoming* cross-family target for other
families migrating later — see `docs/boundaries.md`.)
