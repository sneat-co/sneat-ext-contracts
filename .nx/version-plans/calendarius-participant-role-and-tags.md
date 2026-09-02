---
calendarius-contract: minor
---

participant role and happening tags on the calendarius contract

Additive only: `IHappeningBase.tags`, the closed `HappeningParticipantRole`
vocabulary with its list/default/guard, `IHappeningContactRef` (the typed
add/remove-participant ref carrying the optional role), and the tag limits plus
the normaliser that mirrors what the server stores. `IHappeningParticipant` was
adopted rather than deleted: its untyped `roles?: string[]` becomes
`readonly HappeningParticipantRole[]`, which narrows a type nothing referenced.
