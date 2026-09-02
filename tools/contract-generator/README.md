# contract-generator

Local Nx generator collection (not an npm-published package, not a pnpm
workspace member — see `package.json`) that makes the DEFAULT path for a new
extension contract — "it lives in `sneat-ext-contracts`" — also the EASIEST
path: one command scaffolds a fully-conventional `libs/<family>` (plus an
optional Go module) and registers it wherever the repo tracks its family
roster.

## Usage

```sh
pnpm nx g ./tools/contract-generator:contract <family>
pnpm nx g ./tools/contract-generator:contract <family> --go
```

`<family>` is lowercase kebab-case (e.g. `taxus`, `kids-club`). This creates:

- `libs/<family>/` — Nx project `<family>-contract`, npm package
  `@sneat/extension-<family>-contract`: `project.json`, `package.json` (peer
  floors read live from the workspace root — never a hardcoded/bare `^0.0.x`
  range), `ng-package.json`, `tsconfig*.json`, `vitest.config.mts` +
  `tsconfig.spec.json`, `src/index.ts` stub, `README.md` with a
  provenance/purpose TODO.
- With `--go`: `<family>/go.mod` (module
  `github.com/sneat-co/sneat-ext-contracts/<family>`, Go version read live
  from `go.work`'s own `go <version>` line) and a `doc.go` stub, plus
  `./<family>` appended to `go.work`'s `use ( ... )` block. Go package
  identifiers can't contain hyphens, so a kebab-case family (`kids-club`)
  gets Go package `kids_club`; the module *path* keeps the hyphen.

It also registers the family in the repo's shared files (see "Touching
shared files" below):

- `contracts.json` — appended to `families` (sorted, de-duplicated).
- `eslint.config.mjs` — the default self-only `depConstraints` entry
  (`{ sourceTag: 'family:<family>', onlyDependOnLibsWithTags: ['family:<family>'] }`).
- `docs/boundaries.md` — a row in the "Registered families" table (created on
  first use; this table doesn't exist in the doc until the first family is
  scaffolded through this generator — see below).

All three registrations are idempotent: re-running the generator for an
already-registered family is a no-op on those three files (it still fails
fast if `libs/<family>` already exists, since the lib scaffold itself is not
meant to be re-run).

Every invocation prints the next steps: filling in the contract, the version-
plan ritual (`pnpm nx release plan <family>-contract`, including the 0.x
downshift table from the repo README), the commit-scope rule
(`(<family>-contract)` must match the project name or the bump caps at
patch), and — with `--go` — the Go build/CI-discovery step.

## Design notes for whoever touches this next

**Touching shared files.** `docs/boundaries.md`'s "Adding a new family lib"
checklist already told humans to hand-edit `contracts.json` and
`eslint.config.mjs` for every new family; this generator automates exactly
that checklist, plus a new `docs/boundaries.md` table so the family roster is
readable without parsing JSON. That's a deliberate exception to this repo's
general rule that migration lanes never touch shared files: a generator run
is a single human/agent-invoked change (one CLI call, one PR), not a parallel
lane racing four siblings over the same lines — so the collision risk that
rule exists to prevent doesn't apply here. The generator's own template/code
is the only place allowed to make that exception; it must never be used as
precedent for a migration lane to edit these files directly.

**Why `vitest.config.mts`, not `vite.config.mts`.** `nx.json` registers
`@nx/vite` and `@nx/vitest` as separate inferred-target plugins with
*different* target names (`@nx/vite` → `test`/`build`; `@nx/vitest` →
`vite:test`). `@nx/vite`'s plugin only scans `vite.config.*`; `@nx/vitest`
scans both `vite.config.*` and `vitest.config.*`, always creating a test
target for the latter. Naming the scaffolded file `vitest.config.mts` means
only `@nx/vitest` ever touches it, so the project's real "build" target
(`@nx/angular:package`/ng-packagr, defined explicitly in `project.json`)
never collides with `@nx/vite`'s own inferred "build" target. Run a scaffolded
family's tests with `pnpm nx run <family>-contract:vite:test`.

**Peer floors are read live, not templated.** The generator reads the
workspace root `package.json`'s `devDependencies` for `@angular/core` and
`rxjs` and floors them at their current major (`^22.0.0`, `^7.0.0` today);
`@sneat/core`/`data`/`dto`/`space-models` peer ranges are read from the same
root `devDependencies` and given a caret (root pins these exactly, e.g.
`0.27.6`, for its own build — the generator must never copy that exact pin
straight through as a peer, only `^0.27.6`). This is what keeps every new
contract's peer floors honest as the workspace's own majors move, and is how
the generator satisfies `peer-range-strict` (ci.yml) by construction instead
of by author discipline.

**Testing.** `generator.spec.ts` shells out to the real
`pnpm nx g ./tools/contract-generator:contract` CLI path (not an in-memory
`Tree` fixture) against a scratch family name, asserts the scaffolded files
and all three shared-file registrations, runs `pnpm nx run-many -t lint`
across the (now larger) workspace, and replicates `ci.yml`'s `discover-go`
matrix-discovery command to prove a `--go` scratch module is found. It
snapshots and restores every shared file it touches and deletes every path it
created, in a `finally`, so a failed run never leaves scratch state behind.
Run it directly (it is intentionally NOT an Nx project — see `package.json`'s
comment on why — so `nx run-many` never sees it):

```sh
pnpm test:contract-generator
```
