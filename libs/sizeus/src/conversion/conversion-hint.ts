import { SizeSystem } from '../size-system';
import { FOOTWEAR_ADULT_CONVERSION, FOOTWEAR_KIDS_CONVERSION } from './footwear';
import { IConversionPairReliability, ISizeTypeConversionTable } from './types';

// Every conversion table shipped in this build. Adding a table for a new
// size type (or a new system pair on an existing one) is pure data — append
// here — mirroring AC catalog-extends-without-migration for the catalog
// itself; no change to `conversionHint` below is required.
const CONVERSION_TABLES: readonly ISizeTypeConversionTable[] = [
	FOOTWEAR_ADULT_CONVERSION,
	FOOTWEAR_KIDS_CONVERSION,
];

function findPair(
	table: ISizeTypeConversionTable,
	from: SizeSystem,
	to: SizeSystem,
): IConversionPairReliability | undefined {
	return table.pairs.find(
		(pair) => (pair.from === from && pair.to === to) || (pair.from === to && pair.to === from),
	);
}

// Looks up an indicative conversion hint for `value` (declared in `system`)
// of size type `sizeTypeID`, expressed in `targetSystem`.
//
// Returns null — meaning "downstream shows no hint" per the spec — when:
//   - there is no conversion table for `sizeTypeID`;
//   - `system` and `targetSystem` are the same;
//   - the table has no reliability entry for the (system, targetSystem)
//     pair, or the pair is explicitly flagged unreliable;
//   - `value` is not present in the table under `system` (e.g. an
//     out-of-range or free-typed value).
//
// Never mutates `value` or any table data — this is a pure lookup.
export function conversionHint(
	sizeTypeID: string,
	system: SizeSystem,
	value: string,
	targetSystem: SizeSystem,
): string | null {
	if (system === targetSystem) {
		return null;
	}

	const table = CONVERSION_TABLES.find((candidate) => candidate.sizeTypeID === sizeTypeID);
	if (!table) {
		return null;
	}

	const pair = findPair(table, system, targetSystem);
	if (!pair || !pair.reliable) {
		return null;
	}

	const row = table.rows.find((candidate) => candidate[system] === value);
	if (!row) {
		return null;
	}

	return row[targetSystem] ?? null;
}
