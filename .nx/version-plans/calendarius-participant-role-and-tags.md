---
calendarius-contract: patch
---

participant role and happening tags on the calendarius contract

PATCH, not minor, and the distinction is load-bearing rather than modest. Every
change here is purely ADDITIVE — new exports plus one narrowing of a type with
zero references — so no consumer has to change to keep working. Under 0.x caret
semantics `^0.27.0` means `>=0.27.0 <0.28.0`, so cutting 0.28.0 would fall
OUTSIDE the range three published manifests already declare
(`@sneat/extension-contactus-ui@0.14.0`, and calendarius' own `-ui` and
`-runtime` peers), forcing an upstream republish of contactus-ui and a
consumer-bump wave for a release that breaks nobody. 0.27.4 satisfies all three
as they stand.

Adds `IHappeningBase.tags`, the closed `HappeningParticipantRole` vocabulary
with its list/default/guard, `IHappeningContactRef` (the typed
add/remove-participant ref carrying the optional role), and the tag limits plus
the normaliser that mirrors what the server stores. `IHappeningParticipant` was
adopted rather than deleted: its untyped `roles?: string[]` becomes
`readonly HappeningParticipantRole[]`, which narrows a type nothing referenced.

Also widens `happeningTagError`'s control-character class to the C1 range
(U+0080–U+009F) so it matches Go's `unicode.IsControl` exactly, closing a gap
where a tag passed the client guard and was then refused by the server.
