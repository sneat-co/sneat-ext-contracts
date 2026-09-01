## 0.2.4 (2026-09-01)

### 🩹 Fixes

- feat(budgetus): normalizeMemberID belongs with the member rules ([e709fc8](https://github.com/sneat-co/sneat-ext-contracts/commit/e709fc8))

  memberIDs on the rollup are documented as the SHORT contact id form, so
  normalizing the long `contactID@spaceID` form is a rule about the
  contract's own shape. It lived in the runtime, where UI code — which by
  fleet convention depends on the contract and never on an extension's
  runtime — could not reach it, forcing either a duplicate implementation
  or a layering violation.

### ❤️ Thank You

- Alexander Trakhimenok @trakhimenok
- Claude Fable 5

## 0.2.3 (2026-09-01)

### 🩹 Fixes

- feat(budgetus): share the per-member split rule ([b4469e5](https://github.com/sneat-co/sneat-ext-contracts/commit/b4469e5))

  Adds splitMinorUnitsAcrossMembers() and memberShareOfLine() so the
  projection and any UI showing one member's costs use ONE implementation
  of the split. A per-member breakdown must always add up to the currency
  total above it; two implementations were two chances to disagree.

  Adds an OPTIONAL watchBudgetForSpace(space, window) so a caller that
  already holds the full ISpaceContext can pass it through to the data
  source instead of it being rebuilt as a bare { id }.

  Also aligns the @angular/core peer with the fleet convention
  (">=22.0.0 <23.0.0").

### ❤️ Thank You

- Alexander Trakhimenok @trakhimenok
- Claude Fable 5

## 0.2.2 (2026-09-01)

### 🩹 Fixes

- fix(budgetus): restore types erased to `any` ([edc9fda](https://github.com/sneat-co/sneat-ext-contracts/commit/edc9fda))

  The 0.1.0 -> 0.1.2 migration replaced five published types with `any`:
  IListContext's Dbo generic, IBudgetusSpaceDbo.listGroups, the space
  parameters of IBudgetusService.deleteList/getListById, and
  IListItemResult.listDto. Contract 0.1.0 contained no `any` at all.

  Consumers lost type checking on list DBOs and could not compile against
  anything newer than 0.1.0.

### ❤️ Thank You

- Alexander Trakhimenok @trakhimenok
- Claude Fable 5

## 0.2.1 (2026-09-01)

### 🩹 Fixes

- fix(budgetus): restore spaceID on list requests ([6150f74](https://github.com/sneat-co/sneat-ext-contracts/commit/6150f74))

  IListRequest and ICreateListRequest lost `extends ISpaceRequest` between
  0.1.0 and 0.1.2, dropping spaceID from the whole list-request family even
  though the service still sends it and the backend still routes on it.
  Consumers could not compile against anything newer than 0.1.0.

### ❤️ Thank You

- Alexander Trakhimenok @trakhimenok
- Claude Fable 5

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