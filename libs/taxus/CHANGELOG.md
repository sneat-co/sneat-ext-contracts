## 0.0.4 (2026-08-26)

### 🩹 Fixes

- Fix: add `isolatedModules: true` to libs/taxus/tsconfig.json so the ([8ad9853](https://github.com/sneat-co/sneat-ext-contracts/commit/8ad9853))
  `const enum ListPage` re-exported from src/constants.ts is preserved as a
  real runtime object in the built `.mjs` instead of being erased — the
  live npm defect this batch's Batch-1 merge repairs (0.0.3 -> 0.0.4).

### ❤️ Thank You

- Alexander Trakhimenok

## 0.0.3 (2026-08-26)

### 🩹 Fixes

- Seed migration: taxus contract moves from sneat-co/ext-taxus into sneat-ext-contracts, API-identical to npm's published 0.0.2. ([f138c5b](https://github.com/sneat-co/sneat-ext-contracts/commit/f138c5b))

### ❤️ Thank You

- Alexander Trakhimenok