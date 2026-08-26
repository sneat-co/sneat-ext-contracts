import { ISizeRecord } from './size-record';

// Implements REQ size-history's "current" rule: the current size record is
// the one with the latest non-future effective date; ties on the same
// effective date are broken by the most-recently-created record (using the
// server-assigned `createdAt` timestamp).
//
// Pure and dependency-free by design — both the extension libs and the
// standalone app (and their tests) share this single implementation instead
// of re-deriving the rule.
export function currentSizeRecord(
	records: readonly ISizeRecord[],
	todayISO: string,
): ISizeRecord | undefined {
	let current: ISizeRecord | undefined;

	for (const record of records) {
		if (record.effectiveDate > todayISO) {
			continue; // Future-dated records are dormant until their day arrives.
		}
		if (!current) {
			current = record;
			continue;
		}
		if (
			record.effectiveDate > current.effectiveDate ||
			(record.effectiveDate === current.effectiveDate &&
				record.createdAt > current.createdAt)
		) {
			current = record;
		}
	}

	return current;
}
