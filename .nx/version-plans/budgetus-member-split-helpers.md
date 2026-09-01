---
budgetus-contract: patch
---

feat(budgetus): share the per-member split rule

Adds splitMinorUnitsAcrossMembers() and memberShareOfLine() so the
projection and any UI showing one member's costs use ONE implementation
of the split. A per-member breakdown must always add up to the currency
total above it; two implementations were two chances to disagree.

Also aligns the @angular/core peer with the fleet convention
(">=22.0.0 <23.0.0").
