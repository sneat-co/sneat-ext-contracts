import { ISizeTypeConversionTable } from './types';
import { FOOTWEAR_ADULT_CONVERSION, FOOTWEAR_KIDS_CONVERSION } from './footwear';

// Every conversion table the contract ships, in the order they should be
// presented. Consumers (e.g. the sizechart.app landing that generates one
// SEO page per chart) enumerate this instead of hard-coding the table names,
// so adding a table is pure data here — mirroring the
// catalog-extends-without-migration principle.
export const ALL_CONVERSION_TABLES: readonly ISizeTypeConversionTable[] = [
	FOOTWEAR_ADULT_CONVERSION,
	FOOTWEAR_KIDS_CONVERSION,
];

// Look up a single table by its catalog size-type id, or undefined if none
// ships for that type yet.
export function conversionTableFor(
	sizeTypeID: string,
): ISizeTypeConversionTable | undefined {
	return ALL_CONVERSION_TABLES.find((t) => t.sizeTypeID === sizeTypeID);
}
