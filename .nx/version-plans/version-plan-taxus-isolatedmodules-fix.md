---
taxus-contract: patch
---

Fix: add `isolatedModules: true` to libs/taxus/tsconfig.json so the
`const enum ListPage` re-exported from src/constants.ts is preserved as a
real runtime object in the built `.mjs` instead of being erased — the
live npm defect this batch's Batch-1 merge repairs (0.0.3 -> 0.0.4).
