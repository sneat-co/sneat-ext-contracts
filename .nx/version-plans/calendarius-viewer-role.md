---
calendarius-contract: major
---

Add the live-show participant roles `viewer` (a view-only person in a video
show's call; a Space member or a contact, like `caller` and `speaker`) and
`follow-upper` (a staff assignment for a private follow-up room on the show,
already written by the calendarius backend) to `HappeningParticipantRole` and
`happeningParticipantRoles`, after `speaker` (existing values and default
unchanged), so the list matches the backend's
`HappeningParticipantWellKnownRoles` exactly, in order.
`isKnownHappeningParticipantRole` now accepts both.

Breaking for TypeScript for the same reason as `calendarius-show-roles.md`
(an exhaustive `switch` over the union must add a case). If the pending
`calendarius-show-roles.md` and `calendarius-speaker-role.md` plans have not
been released yet, all three are consumed together and give one bump; if they
have, this plan gives the next one.
