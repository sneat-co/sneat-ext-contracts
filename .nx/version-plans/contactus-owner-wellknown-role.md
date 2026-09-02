---
contactus-contract: patch
---

fix(contactus): SpaceMemberRoleOwner belongs in SpaceMemberWellKnownRoles

`owner` was the one SpaceMemberRole constant missing from the well-known list.
Consumers read absence from that list as "a role nobody has reasoned about" and
warn on it, so every owner-granting invite claim raised a permanent false alarm
in the exact detector added to catch unreviewed roles.

Go-side change only, but contactus ships both artifacts, so the family version
moves in lockstep per the repo's versioning rule.
