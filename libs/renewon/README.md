# @sneat/extension-renewon-contract

Public RenewOn DTOs, contexts, service interfaces, and dependency-injection
tokens (lists and renewal tracking). Runtime services and UI implementations
remain in the private RenewOn repository.

## Provenance

Migrated from `sneat-co/ext-renewon`
(`frontend/libs/extensions/renewon/contract`), commit
`3edf6dd65b3393b866216b629b143fe700c2eb70` (`origin/main`, 2026-08-26 read).

npm reality wins (REQ `version-continuity-from-npm`): npm's actual latest is
`0.2.0`. `ext-renewon` carries **no git tags at all**, and the source repo's
on-disk `package.json` at the above commit still reads `"version": "0.1.0"` —
the git history never shows a version-bump commit past `0.1.0`, even though
npm has published both `0.1.0` and `0.2.0` (8 minutes apart, 2026-07-14).
Disk version here is seeded at `0.2.0` to match npm reality; the paired
version plan (`.nx/version-plans/`) requests the patch bump to `0.2.1`.

`.d.ts` **and** compiled-output parity gate (npm latest verified via
`npm view @sneat/extension-renewon-contract dist-tags`):

- `pnpm nx build renewon-contract` against source repo HEAD, then
  `diff dist/.../types/sneat-extension-renewon-contract.d.ts` against
  `npm pack @sneat/extension-renewon-contract@0.2.0`'s unpacked
  `types/sneat-extension-renewon-contract.d.ts` — byte-identical (diff exit 0).
- Same diff against the unpacked `fesm2022/sneat-extension-renewon-contract.mjs`
  (the compiled runtime output, not just type declarations) — also
  byte-identical. So despite the git-side version label lagging and never
  reaching `0.2.0`, the actual `src/` tree at HEAD is both API- and
  runtime-identical to what npm already serves as `0.2.0`; nothing under
  `src/` migrated here goes beyond published reality.

The one real drift between source-repo HEAD and npm `0.2.0` is
`package.json` metadata only: HEAD's `peerDependencies` float
`@sneat/{core,data,dto,space-models}` at `^0.24.0` (bumped by post-publish
dependency-update commits that were never re-published), while npm `0.2.0`
still declares `^0.22.1`. There is no tag or publish backing the `0.24.0`
floor, so per version-continuity-from-npm this migration keeps npm `0.2.0`'s
declared peer floors (`^0.22.1`) rather than carrying the unpublished bump
forward — matching the fleet's existing convention of contracts keeping
conservative peer floors (taxus, debtus, assetus all do the same).
