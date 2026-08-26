# @sneat/extension-bookius-contract

Public Bookius booking DTOs and the `IBookiusService` / `BOOKIUS_SERVICE`
dependency-injection token. Booking-engine runtime, persistence, and UI
implementations remain in the private Bookius repository
(`sneat-co/bookius`).

## Provenance

Migrated from `sneat-co/ext-bookius` (`frontend/`), commit
`cc2ef1fa984721b456d1a2618661f119ca13b29e` (`origin/main`, 2026-08-26 read).

- npm reality wins (REQ `version-continuity-from-npm`): npm's actual latest
  is `0.4.0`. `git diff v0.4.0..origin/main` touches only `backend/go.mod`
  (a Go toolchain bump, 1.24.0 → 1.26.0 + toolchain go1.27.0) — zero files
  under `frontend/` changed since the tag that produced the npm `0.4.0`
  publish (`v0.4.0` tag == commit `b8e0e4a`, the same commit tagged
  `backend/v0.1.0`). Disk version here is seeded at `0.4.0` to match both npm
  reality and the source tree exactly; the paired version plan
  (`.nx/version-plans/`) requests the patch bump to `0.4.1`.
- `.d.ts` parity gate: `pnpm nx build bookius-contract`, then the resulting
  bundled `dist/libs/bookius/types/sneat-extension-bookius-contract.d.ts` was
  diffed (import lines and the trailing `export {}` / `export type {}`
  aggregate stripped from both sides, since npm's tarball is an unbundled
  per-file `tsc` build and this repo's convention rebuilds every contract
  with `@nx/angular:package`/ng-packagr) against the concatenated body of
  `npm pack @sneat/extension-bookius-contract@0.4.0`'s unpacked
  `dist/dto/booking.d.ts` + `dist/dto/bookius-team.d.ts` +
  `dist/services/bookius-service.d.ts`. Result: **byte-identical** (diff exit
  0) — every interface, type alias, property (including every `readonly`/`?`
  modifier and string-literal union), and JSDoc comment matches verbatim;
  only the bundling mechanism differs, not the API surface.
- `frontend/src/dto/booking.spec.ts` (6 assertions, covering booking-type
  targets, anonymous public booking requests, the price-free
  competition-entry reservation request, lifecycle-command durable IDs, the
  free-vs-settled confirmation-evidence distinction, and the locked
  participant-cancellation-without-refund case) is carried over verbatim and
  passes unchanged (`pnpm nx test bookius-contract`).

peerDependencies are unchanged from npm `0.4.0`
(`@angular/core": "^21.0.0"`, `@sneat/dto": "^0.22.0 || ^0.24.0"`,
`rxjs": "^7.0.0"`) — matches the fleet's existing convention (e.g.
`@sneat/extension-{debtus,assetus,taxus}-contract` all still declare
`@angular/core": "^21.0.0"` too, despite newer majors being available in this
workspace).

## Go module

`sneat-co/ext-bookius`'s `backend/` Go module (three packages: `botapp`,
`dto4bookius`, `facade4bookius`; no external dependencies beyond stdlib, no
internal cross-package imports) was migrated alongside this npm lib to
`<repo-root>/bookius/` as `github.com/sneat-co/sneat-ext-contracts/bookius`,
same source commit. It is intentionally **not** wired into this repo's
`go.work` (left for the batch's integration merge — see this branch's task
brief) but builds, vets, and tests green standalone
(`cd bookius && GOWORK=off go build ./... && GOWORK=off go test ./...`).
