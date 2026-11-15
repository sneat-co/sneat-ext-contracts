---
calendarius-contract: major
---

Add the live-show participant roles `host`, `moderator` and `caller` to
`HappeningParticipantRole` and `happeningParticipantRoles`, after the existing
four (order and `defaultHappeningParticipantRole` unchanged), matching the
calendarius backend.

On the wire this is additive. For TypeScript it is breaking: widening the union
breaks an exhaustive `switch` over it, such as the calendarius UI's
`happeningParticipantRoleTitle`, which must add the three cases when it adopts
this version. `major` on the 0.x train gives 0.29.0, outside the `^0.28.0`
ranges consumers declare, so nobody receives the wider union on a routine
patch update.
