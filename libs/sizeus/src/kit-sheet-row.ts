import { SizeSystem } from './size-system';

// One cell of a kit-sheet export: the current value for a size type, with
// its declared system (per REQ kit-sheet-export the cell renders the value
// together with its system, e.g. "9 US").
export interface IKitSheetCell {
	readonly value: string;
	readonly system: SizeSystem;
}

// One row of the kit-sheet CSV export: one row per member contact, one
// column per size type. Per REQ kit-sheet-export the cell is the *current*
// value (see `currentSizeRecord`) for that size type, or absent if the
// contact has no size record for it.
export interface IKitSheetRow {
	readonly contactID: string;
	readonly contactTitle: string;
	// Keyed by sizeTypeID.
	readonly cells: Readonly<Record<string, IKitSheetCell>>;
}
