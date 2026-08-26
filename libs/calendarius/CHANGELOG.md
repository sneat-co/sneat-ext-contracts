## 0.27.3 (2026-08-26)

### 🩹 Fixes

- feat(calendarius-contract): migrate calendarius contract (npm + Go) from sneat-co/sneat-libs and sneat-co/ext-calendarius ([1e1cd5e](https://github.com/sneat-co/sneat-ext-contracts/commit/1e1cd5e))

  Migrate @sneat/extension-calendarius-contract into sneat-ext-contracts monorepo, API-identical to npm's published 0.27.2 (fresh build byte-diffed against the npm tarball: .d.ts and .mjs both identical; 0.27.1 and 0.27.2 are themselves content-identical on npm). Provenance: sneat-co/sneat-libs, libs/extensions/calendarius/contract @ 344831c (sole current publisher; not modified by this branch — removing calendarius from its release set is owed separately). Go backend (github.com/sneat-co/sneat-ext-contracts/calendarius) migrated from sneat-co/ext-calendarius backend/ @ tag backend/v0.0.6 (the weekly/fortnightly/monthly/yearly recurrence-validation fix), GOWORK=off build/vet/test green.

### ❤️ Thank You

- Alexander Trakhimenok
- Claude Fable 5