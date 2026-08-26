import { SizeSystem } from './size-system';

// A single size record for one contactus contact within a Space.
//
// Per REQ size-record: a size record belongs to exactly one contact, holds a
// size-type (catalog id), a value with its declared size system (stored
// verbatim — see REQ size-systems / `SizeSystem`), an effective date, and a
// server-assigned creation timestamp used to deterministically break
// same-date ties (see REQ size-history / `currentSizeRecord`).
//
// Persistence (backend concern, not part of this DTO's shape): records live
// under `/spaces/{spaceID}/ext/sizeus/...`, keyed so all records for one
// contact are retrievable in a single query.
export interface ISizeRecord {
	readonly contactID: string;
	readonly sizeTypeID: string;

	// Stored verbatim in the declared `system` — never normalized.
	readonly value: string;
	readonly system: SizeSystem;

	// ISO date string (YYYY-MM-DD). Defaults to "today" at creation time in
	// the app layer; the DTO itself always carries an explicit value.
	readonly effectiveDate: string;

	// ISO timestamp string, server-assigned at creation. Exists solely to
	// make the same-effective-date tie-break deterministic — see
	// REQ size-history and the `currentSizeRecord` helper.
	readonly createdAt: string;

	readonly note?: string;
	readonly preferredBrandNote?: string;
}
