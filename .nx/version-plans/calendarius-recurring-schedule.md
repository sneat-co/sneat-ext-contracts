---
calendarius-contract: major
---

Go: add `facade4calendarius.RecurringScheduleFacade`: idempotent `CreateRecurringHappening` with weekly weekday slots and an IANA zone, plus `CancelOccurrence`, `RescheduleOccurrence`, `SetParticipants` and `CancelHappening`. Each request ID is recorded per occurrence, so a late retry is a no-op.

Occurrences get a stable identity (`HappeningID`, `SlotID`, `LocalDate`) and `OriginalStart`/`Cancelled`/`Rescheduled`. `RecurringHappening` gets `SpaceID` and `Participants`. Also new: `ResolveLocalTime`, `ValidateOccurrenceStart` and the `RunRecurringScheduleFacadeConformance` suite.

No existing interface signature changes, but the documented `RecurringHappeningsFacade.ExpandOccurrences` behaviour changes:
- cancelled occurrences are returned flagged `Cancelled`;
- rescheduled occurrences are returned at their moved time;
- each slot's stored zone wins over the caller's `timezone`, which is now only a fallback for legacy slots;
- an invalid zone is an error instead of UTC.

Consumers that key occurrences by start date (togethered) must move to `LocalDate` and skip `Cancelled`.
