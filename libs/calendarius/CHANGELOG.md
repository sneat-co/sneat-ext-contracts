## 0.29.0 (2026-11-28)

### ⚠️  Breaking Changes

- Add the live-show participant role `speaker` (may be in the call right now) to ([a65d72e](https://github.com/sneat-co/sneat-ext-contracts/commit/a65d72e))
  `HappeningParticipantRole` and `happeningParticipantRoles`, after `caller`
  (existing values, order and default unchanged), matching the calendarius
  backend. The path is spectator -> `caller` -> `speaker` -> removed = revoked;
  `participant` keeps its meaning.

  Breaking for TypeScript for the same reason as `calendarius-show-roles.md`
  (an exhaustive `switch` over the union must add a case). If that plan has not
  been released yet, both plans are consumed together and give one bump,
  0.28.0 -> 0.29.0; if it has, this plan gives 0.29.0 -> 0.30.0.

- Add the live-show participant roles `host`, `moderator` and `caller` to ([356fa6b](https://github.com/sneat-co/sneat-ext-contracts/commit/356fa6b))
  `HappeningParticipantRole` and `happeningParticipantRoles`, after the existing
  four (order and `defaultHappeningParticipantRole` unchanged), matching the
  calendarius backend.

  On the wire this is additive. For TypeScript it is breaking: widening the union
  breaks an exhaustive `switch` over it, such as the calendarius UI's
  `happeningParticipantRoleTitle`, which must add the three cases when it adopts
  this version. `major` on the 0.x train gives 0.29.0, outside the `^0.28.0`
  ranges consumers declare, so nobody receives the wider union on a routine
  patch update.

- Add the live-show participant roles `viewer` (a view-only person in a video ([85ebe13](https://github.com/sneat-co/sneat-ext-contracts/commit/85ebe13))
  show's call; a Space member or a contact, like `caller` and `speaker`) and
  `follow-upper` (a staff assignment for a private follow-up room on the show,
  already written by the calendarius backend) to `HappeningParticipantRole` and
  `happeningParticipantRoles`, after `speaker` (existing values and default
  unchanged), so the list matches the backend's
  `HappeningParticipantWellKnownRoles` exactly, in order.
  `isKnownHappeningParticipantRole` now accepts both.

  Breaking for TypeScript for the same reason as `calendarius-show-roles.md`
  (an exhaustive `switch` over the union must add a case). If the pending
  `calendarius-show-roles.md` and `calendarius-speaker-role.md` plans have not
  been released yet, all three are consumed together and give one bump; if they
  have, this plan gives the next one.

### ❤️ Thank You

- Alexander Trakhimenok @trakhimenok

## 0.28.0 (2026-11-08)

### ⚠️  Breaking Changes

- Go: add `facade4calendarius.RecurringScheduleFacade`: idempotent `CreateRecurringHappening` with weekly weekday slots and an IANA zone, plus `CancelOccurrence`, `RescheduleOccurrence`, `SetParticipants` and `CancelHappening`. Each request ID is recorded per occurrence, so a late retry is a no-op. ([b98f309](https://github.com/sneat-co/sneat-ext-contracts/commit/b98f309))

  Occurrences get a stable identity (`HappeningID`, `SlotID`, `LocalDate`) and `OriginalStart`/`Cancelled`/`Rescheduled`. `RecurringHappening` gets `SpaceID` and `Participants`. Also new: `ResolveLocalTime`, `ValidateOccurrenceStart` and the `RunRecurringScheduleFacadeConformance` suite.

  No existing interface signature changes, but the documented `RecurringHappeningsFacade.ExpandOccurrences` behaviour changes:
  - cancelled occurrences are returned flagged `Cancelled`;
  - rescheduled occurrences are returned at their moved time;
  - each slot's stored zone wins over the caller's `timezone`, which is now only a fallback for legacy slots;
  - an invalid zone is an error instead of UTC.

  Consumers that key occurrences by start date (togethered) must move to `LocalDate` and skip `Cancelled`.

### ❤️ Thank You

- Alexander Trakhimenok

## 0.27.8 (2026-09-07)

### 🚀 Features

- Add server-authored financial agreement, enrollment, and linked-task-source ([5d015f7](https://github.com/sneat-co/sneat-ext-contracts/commit/5d015f7))
  contracts across Calendarius, Listus, Debtus, and Splitus: selected agreement
  charges and coverage, linked date-task navigation/reader seams, source-managed
  todo and due-task contracts, owned transfer lifecycle and repayment, and
  resolved/captured source attribution — all backward-compatible additions.

### ❤️ Thank You

- Alexander Trakhimenok
- Claude Fable 5.1

## 0.27.7 (2026-09-06)

### 🩹 Fixes

- Publish Budgetus source completeness and recurring occurrence provenance together with the Calendarius ID-keyed day and adjusted-slot related-link contract. ([ea90e9c](https://github.com/sneat-co/sneat-ext-contracts/commit/ea90e9c))

### ❤️ Thank You

- Alexander Trakhimenok

## 0.27.6 (2026-09-05)

### 🩹 Fixes

- Allow scheduled responsibilities to atomically carry generic extension data ([ee539d3](https://github.com/sneat-co/sneat-ext-contracts/commit/ee539d3))
  and validated same-Space related references on their Calendarius happening.

### ❤️ Thank You

- Alexander Trakhimenok

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