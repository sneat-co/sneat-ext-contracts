import { SizeSystem } from '../size-system';

// The shape a size type's recorded value takes. Per REQ catalog-extensibility
// a new size type is declared purely by adding a data entry — `valueKind`
// tells the UI/validation layer how to render an input for the type without
// any code change:
//   - 'numeric'     — a single number (e.g. shoe size "9", ring size "7").
//   - 'alpha'       — a letter size (e.g. "XS" … "XL").
//   - 'measurement' — a physical measurement (e.g. height in cm, weight in kg).
export type SizeValueKind = 'numeric' | 'alpha' | 'measurement';

// The MVP catalog's top-level groupings, per AC mvp-catalog-coverage.
// 'sport-kit' covers the individual garment/equipment entries that make up a
// football/basketball/cycling/swimming kit (see `ISportKit` below for how
// those entries are grouped per sport).
export type SizeTypeCategory =
	| 'footwear'
	| 'body-measurement'
	| 'tops'
	| 'bottoms'
	| 'accessories'
	| 'sport-kit';

// One entry in the size-type catalog. Adding a new entry (e.g. "ski boots")
// is a pure data change — see AC catalog-extends-without-migration — because
// `ISizeRecord.sizeTypeID` (in size-record.ts) is a plain string, not a union
// tied to this catalog.
export interface ISizeTypeEntry {
	readonly id: string;
	readonly category: SizeTypeCategory;
	readonly title: string;
	readonly valueKind: SizeValueKind;
	// Size systems a value for this type may declare (see `SizeSystem`).
	// Non-exhaustive by design of the open `SizeSystem` union — this list is
	// the set of systems the catalog/UI suggests for this type, not a hard
	// validation constraint.
	readonly applicableSystems: readonly SizeSystem[];
}

// A named grouping of size-type ids that make up one sport's kit (per the
// brief's football/basketball/cycling/swimming kit types). Purely a grouping
// — every id in `sizeTypeIDs` must also appear in `ISizeTypeCatalog.sizeTypes`
// as its own catalog entry.
export interface ISportKit {
	readonly id: string;
	readonly title: string;
	readonly sizeTypeIDs: readonly string[];
}

// The versioned size-type catalog. Shipped as build-time data (see
// `mvp-catalog.ts`) — bumping `version` is how a future breaking reshape of
// the catalog itself would be signalled; additive entries (new size types)
// never require a version bump per AC catalog-extends-without-migration.
export interface ISizeTypeCatalog {
	readonly version: number;
	readonly sizeTypes: readonly ISizeTypeEntry[];
	readonly sportKits: readonly ISportKit[];
}
