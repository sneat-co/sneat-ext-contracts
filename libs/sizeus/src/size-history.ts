import { ISizeRecord } from './size-record';
import { SizeSystem } from './size-system';

// Full append-first history for one (contact, size type) pair.
//
// Per REQ size-history: history is retrievable chronologically. `records` is
// ordered oldest-first by (effectiveDate, createdAt) — the same ordering the
// `currentSizeRecord` helper scans to resolve the current value.
export interface ISizeHistory {
	readonly contactID: string;
	readonly sizeTypeID: string;
	readonly records: readonly ISizeRecord[];
}

// The resolved "current value" for one (contact, size type) pair — the
// narrow shape used wherever only the current value is needed (e.g. a
// current-size list row), as opposed to the full `ISizeRecord` this value
// was resolved from via `currentSizeRecord`.
export interface ICurrentSize {
	readonly contactID: string;
	readonly sizeTypeID: string;
	readonly value: string;
	readonly system: SizeSystem;
	readonly effectiveDate: string;
}
