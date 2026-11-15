---
calendarius-contract: major
---

Add the live-show participant role `speaker` (may be in the call right now) to
`HappeningParticipantRole` and `happeningParticipantRoles`, after `caller`
(existing values, order and default unchanged), matching the calendarius
backend. The path is spectator -> `caller` -> `speaker` -> removed = revoked;
`participant` keeps its meaning.

Breaking for TypeScript for the same reason as `calendarius-show-roles.md`
(an exhaustive `switch` over the union must add a case). If that plan has not
been released yet, both plans are consumed together and give one bump,
0.28.0 -> 0.29.0; if it has, this plan gives 0.29.0 -> 0.30.0.
