import { SizeSystem } from '../size-system';

// One row of a conversion table: the same physical size expressed in every
// system the table covers, keyed by system name (matching one of the values
// in `ISizeTypeConversionTable.systems`). Values are strings — the same
// stored-verbatim convention as `ISizeRecord.value` — so a row can hold
// half-sizes ("9.5") without a numeric-precision concern.
export type IConversionRow = Readonly<Record<string, string>>;

// Declares whether converting between `from` and `to` (in either direction)
// is reliable enough to surface as a hint downstream. Per the spec,
// conversions are INDICATIVE ONLY; a table may still hold row values for an
// unreliable pair (e.g. for reference/debugging) while `conversionHint`
// refuses to surface them — see the kids footwear US columns in
// `footwear.ts` for a concrete example of a pair with data but no reliable
// mapping.
export interface IConversionPairReliability {
	readonly from: SizeSystem;
	readonly to: SizeSystem;
	readonly reliable: boolean;
}

// A versioned-by-catalog-id conversion table for one size type. Adding a new
// table (for a new size type, or a new system pair on an existing one) is
// pure data — no schema change — mirroring AC catalog-extends-without-migration
// for the catalog itself.
export interface ISizeTypeConversionTable {
	readonly sizeTypeID: string;
	readonly systems: readonly SizeSystem[];
	readonly pairs: readonly IConversionPairReliability[];
	readonly rows: readonly IConversionRow[];
}
