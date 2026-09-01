---
budgetus-contract: major
---

feat(budgetus): budget rollup for real happening expenses

Breaking: `IBudgetRollup` is now segregated per currency (`byCurrency`)
instead of carrying a single `byMonth`/`annualTotal` pair that silently
assumed one currency. `IMoney.value` is documented as integer minor units.
`watchBudget` takes an optional rollup window, and `maskSurpriseLineItems`
is generic instead of `any[]`.

Adds the per-price line-identity fields (`happeningID`, `priceID`,
`occurrenceMonthISO`) and the per-member rollup types (`memberIDs`,
`IMemberBudgetTotals`) that a real, non-demo budget needs.

On the 0.x train `major` lands as 0.1.2 → 0.2.0.
