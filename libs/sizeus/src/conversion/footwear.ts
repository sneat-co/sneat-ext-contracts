import { ISizeTypeConversionTable } from './types';

// Adult footwear — EU/UK/US. A widely used indicative reference chart (the
// same rounded values many retailer size converters publish); it is NOT
// brand-accurate and is intentionally coarse (no fabricated half-size
// precision beyond what is commonly published). EU<->UK<->US are all
// reasonably stable relative to one another across this range, so every pair
// is flagged reliable.
//
// Deliberately covers only EU/UK/US per the task's MVP scope — JP and
// Mondopoint are declared `applicableSystems` on the catalog entry (so the UI
// can accept a value in those systems) but have no conversion rows yet;
// `conversionHint` correctly returns null for them rather than guessing.
export const FOOTWEAR_ADULT_CONVERSION: ISizeTypeConversionTable = {
	sizeTypeID: 'footwear-adult',
	systems: ['EU', 'UK', 'US'],
	pairs: [
		{ from: 'EU', to: 'UK', reliable: true },
		{ from: 'EU', to: 'US', reliable: true },
		{ from: 'UK', to: 'US', reliable: true },
	],
	rows: [
		{ EU: '39', UK: '6', US: '7' },
		{ EU: '40', UK: '6.5', US: '7.5' },
		{ EU: '41', UK: '7', US: '8' },
		{ EU: '42', UK: '8', US: '9' },
		{ EU: '43', UK: '9', US: '10' },
		{ EU: '44', UK: '10', US: '11' },
		{ EU: '45', UK: '11', US: '12' },
	],
};

// Kids footwear — EU/UK/US. EU<->UK is a single stable industry-standard
// scale across the whole kids range, so that pair is kept reliable. Kids' US
// sizing is not: it splits across toddler / little-kid / big-kid scales that
// do not track the EU/UK scale uniformly across the range, so both US pairs
// are deliberately flagged unreliable. The US column values are kept in the
// row data (for reference / possible future use) but `conversionHint` will
// never surface them, per the "no reliable mapping -> no hint" rule.
export const FOOTWEAR_KIDS_CONVERSION: ISizeTypeConversionTable = {
	sizeTypeID: 'footwear-kids',
	systems: ['EU', 'UK', 'US'],
	pairs: [
		{ from: 'EU', to: 'UK', reliable: true },
		{ from: 'EU', to: 'US', reliable: false },
		{ from: 'UK', to: 'US', reliable: false },
	],
	rows: [
		{ EU: '24', UK: '7', US: '8' },
		{ EU: '26', UK: '8.5', US: '9.5' },
		{ EU: '28', UK: '10', US: '11' },
		{ EU: '30', UK: '11.5', US: '12.5' },
		{ EU: '32', UK: '13', US: '1' },
		{ EU: '34', UK: '1.5', US: '2.5' },
	],
};
