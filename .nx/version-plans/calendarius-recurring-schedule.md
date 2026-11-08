---
calendarius-contract: major
---

Go: add `facade4calendarius.RecurringScheduleFacade` (idempotent `CreateRecurringHappening` with weekly weekday slots and IANA zone, `CancelOccurrence`, `RescheduleOccurrence`, `SetParticipants`, `CancelHappening`), stable occurrence identity (`HappeningID`, `SlotID`, `LocalDate`) plus `Cancelled`/`Rescheduled`/`OriginalStart` on `calendariusmodels.Occurrence`, `SpaceID`/`Participants` on `RecurringHappening`, the `ResolveLocalTime` DST policy, and a `RunRecurringScheduleFacadeConformance` suite. Additive; existing interfaces are unchanged.
