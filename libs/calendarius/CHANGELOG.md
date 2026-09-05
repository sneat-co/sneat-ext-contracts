## 0.27.5 (2026-09-05)

### 🩹 Fixes

- Add generic scheduled responsibility contracts for recurring Calendarius ([c0fe406](https://github.com/sneat-co/sneat-ext-contracts/commit/c0fe406))
  happenings, ordered or fixed assignment, stable occurrence references, and
  transactional completion history. The API is additive and enables verticals to
  present chores or duties without adding another recurrence engine.

### ❤️ Thank You

- Alexander Trakhimenok

## 0.27.4 (2026-09-02)

### 🩹 Fixes

- participant role and happening tags on the calendarius contract ([31ec8fe](https://github.com/sneat-co/sneat-ext-contracts/commit/31ec8fe))

  PATCH, not minor, and the distinction is load-bearing rather than modest. Every
  change here is purely ADDITIVE — new exports plus one narrowing of a type with
  zero references — so no consumer has to change to keep working. Under 0.x caret
  semantics `^0.27.0` means `>=0.27.0 <0.28.0`, so cutting 0.28.0 would fall
  OUTSIDE the range three published manifests already declare
  (`@sneat/extension-contactus-ui@0.14.0`, and calendarius' own `-ui` and
  `-runtime` peers), forcing an upstream republish of contactus-ui and a
  consumer-bump wave for a release that breaks nobody. 0.27.4 satisfies all three
  as they stand.

  Adds `IHappeningBase.tags`, the closed `HappeningParticipantRole` vocabulary
  with its list/default/guard, `IHappeningContactRef` (the typed
  add/remove-participant ref carrying the optional role), and the tag limits plus
  the normaliser that mirrors what the server stores. `IHappeningParticipant` was
  adopted rather than deleted: its untyped `roles?: string[]` becomes
  `readonly HappeningParticipantRole[]`, which narrows a type nothing referenced.

  Also widens `happeningTagError`'s control-character class to the C1 range
  (U+0080–U+009F) so it matches Go's `unicode.IsControl` exactly, closing a gap
  where a tag passed the client guard and was then refused by the server.

### ❤️ Thank You

- Alexander Trakhimenok @trakhimenok
- Claude Fable 5.1

## 0.27.3 (2026-08-26)

### 🩹 Fixes

- feat(calendarius-contract): migrate calendarius contract (npm + Go) from sneat-co/sneat-libs and sneat-co/ext-calendarius ([1e1cd5e](https://github.com/sneat-co/sneat-ext-contracts/commit/1e1cd5e))

  Migrate @sneat/extension-calendarius-contract into sneat-ext-contracts monorepo, API-identical to npm's published 0.27.2 (fresh build byte-diffed against the npm tarball: .d.ts and .mjs both identical; 0.27.1 and 0.27.2 are themselves content-identical on npm). Provenance: sneat-co/sneat-libs, libs/extensions/calendarius/contract @ 344831c (sole current publisher; not modified by this branch — removing calendarius from its release set is owed separately). Go backend (github.com/sneat-co/sneat-ext-contracts/calendarius) migrated from sneat-co/ext-calendarius backend/ @ tag backend/v0.0.6 (the weekly/fortnightly/monthly/yearly recurrence-validation fix), GOWORK=off build/vet/test green.

### ❤️ Thank You

- Alexander Trakhimenok
- Claude Fable 5