---
contactus-contract: patch
---

fix(contactus): WithMultiSpaceContacts.Validate must not swallow the contactIDs error

It returned nil on failure, discarding the verdict of a real format validator.
Since the mixin is embedded in contactus's ContactDbo, contactIDs were never
actually validated on any contact record.
