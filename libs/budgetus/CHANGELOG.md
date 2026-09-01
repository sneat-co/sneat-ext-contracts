## 0.2.0 (2026-09-01)

### ⚠️  Breaking Changes

- feat(budgetus): budget rollup for real happening expenses ([f482e58](https://github.com/sneat-co/sneat-ext-contracts/commit/f482e58))

  Breaking: `IBudgetRollup` is now segregated per currency (`byCurrency`)
  instead of carrying a single `byMonth`/`annualTotal` pair that silently
  assumed one currency. `IMoney.value` is documented as integer minor units.
  `watchBudget` takes an optional rollup window, and `maskSurpriseLineItems`
  is generic instead of `any[]`.

  Adds the per-price line-identity fields (`happeningID`, `priceID`,
  `occurrenceMonthISO`) and the per-member rollup types (`memberIDs`,
  `IMemberBudgetTotals`) that a real, non-demo budget needs.

  On the 0.x train `major` lands as 0.1.2 → 0.2.0.

### ❤️ Thank You

- Alexander Trakhimenok @trakhimenok
- Claude Fable 5

## 0.1.2 (2026-08-31)

### 🩹 Fixes

- chore(contracts): release peers targeting Sneat Libraries 0.27.6 ([a5f58ec](https://github.com/sneat-co/sneat-ext-contracts/commit/a5f58ec))

### ❤️ Thank You

- Alexander Trakhimenok

## 0.1.1 (2026-08-26)

### 🩹 Fixes

- feat(budgetus-contract): migrate budgetus contract from sneat-co/budgetus ([525ae73](https://github.com/sneat-co/sneat-ext-contracts/commit/525ae73))

  Migrate @sneat/extension-budgetus-contract into sneat-ext-contracts monorepo. Provenance: sneat-co/budgetus, npm continuity from @sneat/extension-budgetus-contract@0.1.0.

### ❤️ Thank You

- Alexander Trakhimenok
- Claude Fable 5