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