## 0.12.3 (2026-08-26)

### 🩹 Fixes

- feat(contactus-contract): migrate contactus contract (npm + Go) from sneat-co/ext-contactus ([fd4997f](https://github.com/sneat-co/sneat-ext-contracts/commit/fd4997f))

  Migrate @sneat/extension-contactus-contract into sneat-ext-contracts monorepo, API-identical to npm's published 0.12.2 (fresh build byte-diffed against the npm tarball: .d.ts and .mjs both identical). Provenance: sneat-co/ext-contactus (the sole publisher — sneat-co/contactus, the product repo, is a consumer only and holds no contract source), frontend/libs/extensions/contactus/contract @ 1d703c3. Go backend (github.com/sneat-co/sneat-ext-contracts/contactus) migrated alongside from the same repo's backend/, highest source tag backend/v0.1.8, GOWORK=off build/vet/test green.

### ❤️ Thank You

- Alexander Trakhimenok
- Claude Fable 5