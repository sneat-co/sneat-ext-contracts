---
budgetus-contract: patch
---

feat(budgetus): normalizeMemberID belongs with the member rules

memberIDs on the rollup are documented as the SHORT contact id form, so
normalizing the long `contactID@spaceID` form is a rule about the
contract's own shape. It lived in the runtime, where UI code — which by
fleet convention depends on the contract and never on an extension's
runtime — could not reach it, forcing either a duplicate implementation
or a layering violation.
