## 0.12.7 (2026-09-02)

### 🩹 Fixes

- fix(contactus): WithMultiSpaceContacts.Validate must not swallow the contactIDs error ([282435c](https://github.com/sneat-co/sneat-ext-contracts/commit/282435c))

  It returned nil on failure, discarding the verdict of a real format validator.
  Since the mixin is embedded in contactus's ContactDbo, contactIDs were never
  actually validated on any contact record.

### ❤️ Thank You

- Alexander Trakhimenok @trakhimenok
- Claude Fable 5

## 0.12.6 (2026-09-02)

### 🩹 Fixes

- fix(contactus): SpaceMemberRoleOwner belongs in SpaceMemberWellKnownRoles ([54aacc7](https://github.com/sneat-co/sneat-ext-contracts/commit/54aacc7))

  `owner` was the one SpaceMemberRole constant missing from the well-known list.
  Consumers read absence from that list as "a role nobody has reasoned about" and
  warn on it, so every owner-granting invite claim raised a permanent false alarm
  in the exact detector added to catch unreviewed roles.

  Go-side change only, but contactus ships both artifacts, so the family version
  moves in lockstep per the repo's versioning rule.

### ❤️ Thank You

- Alexander Trakhimenok @trakhimenok
- Claude Fable 5

## 0.12.5 (2026-08-31)

### 🩹 Fixes

- chore(contracts): release peers targeting Sneat Libraries 0.27.6 ([a5f58ec](https://github.com/sneat-co/sneat-ext-contracts/commit/a5f58ec))

### ❤️ Thank You

- Alexander Trakhimenok

## 0.12.4 (2026-08-27)

### 🩹 Fixes

- fix(contactus-contract): guard SetContactBrief against a nil Contacts map (sneat-co/contactus#60) ([3da261f](https://github.com/sneat-co/sneat-ext-contracts/commit/3da261f))

### ❤️ Thank You

- Alexander Trakhimenok
- Claude Opus 5

## 0.12.3 (2026-08-26)

### 🩹 Fixes

- feat(contactus-contract): migrate contactus contract (npm + Go) from sneat-co/ext-contactus ([fd4997f](https://github.com/sneat-co/sneat-ext-contracts/commit/fd4997f))

  Migrate @sneat/extension-contactus-contract into sneat-ext-contracts monorepo, API-identical to npm's published 0.12.2 (fresh build byte-diffed against the npm tarball: .d.ts and .mjs both identical). Provenance: sneat-co/ext-contactus (the sole publisher — sneat-co/contactus, the product repo, is a consumer only and holds no contract source), frontend/libs/extensions/contactus/contract @ 1d703c3. Go backend (github.com/sneat-co/sneat-ext-contracts/contactus) migrated alongside from the same repo's backend/, highest source tag backend/v0.1.8, GOWORK=off build/vet/test green.

### ❤️ Thank You

- Alexander Trakhimenok
- Claude Fable 5