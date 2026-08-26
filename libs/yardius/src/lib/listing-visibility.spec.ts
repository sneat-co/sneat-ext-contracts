import { describe, it, expect } from 'vitest';
import { ASSETUS_VISIBILITIES, AssetusVisibility } from './assetus-visibility.js';
import {
	allowedListingVisibilities,
	DEFAULT_LISTING_VISIBILITY,
	isListingVisibilityAllowed,
	LISTING_VISIBILITIES,
	LISTING_VISIBILITY_CEILING,
	ListingVisibility,
} from './listing-visibility.js';

describe('listing visibility', () => {
	it('has exactly the three visibilities from REQ listing-visibility-ceiling', () => {
		expect(LISTING_VISIBILITIES).toEqual(['private', 'members', 'friends']);
	});

	it("defaults to 'members'", () => {
		expect(DEFAULT_LISTING_VISIBILITY).toBe('members');
	});
});

describe('LISTING_VISIBILITY_CEILING', () => {
	// The exact mapping table from REQ listing-visibility-ceiling.
	const expectedCeiling: ReadonlyArray<
		[AssetusVisibility, readonly ListingVisibility[]]
	> = [
		['private', ['private']],
		['family', ['private', 'members']],
		['specific_space', ['private']], // MVP simplification.
		['friends', ['private', 'members', 'friends']],
		['friends_of_friends', ['private', 'members', 'friends']],
		['public', ['private', 'members', 'friends']],
	];

	it.each(expectedCeiling)(
		'asset visibility %s permits exactly %j',
		(assetVisibility, permitted) => {
			expect(LISTING_VISIBILITY_CEILING[assetVisibility]).toEqual(permitted);
			expect(allowedListingVisibilities(assetVisibility)).toEqual(permitted);
		},
	);

	it('covers every Assetus visibility exactly once', () => {
		expect(Object.keys(LISTING_VISIBILITY_CEILING).sort()).toEqual(
			[...ASSETUS_VISIBILITIES].sort(),
		);
		expect(expectedCeiling.map(([asset]) => asset).sort()).toEqual(
			[...ASSETUS_VISIBILITIES].sort(),
		);
	});

	it("always permits 'private' — every asset can be listed to the owning Space only", () => {
		for (const assetVisibility of ASSETUS_VISIBILITIES) {
			expect(isListingVisibilityAllowed(assetVisibility, 'private')).toBe(true);
		}
	});
});

describe('isListingVisibilityAllowed', () => {
	it("rejects a 'friends' listing on a 'private' asset (AC visibility-ceiling-enforced)", () => {
		expect(isListingVisibilityAllowed('private', 'friends')).toBe(false);
	});

	it("allows a 'private' listing on a 'private' asset (AC private-listing-on-private-asset)", () => {
		expect(isListingVisibilityAllowed('private', 'private')).toBe(true);
	});

	it("rejects listing visibilities above the 'family' ceiling", () => {
		expect(isListingVisibilityAllowed('family', 'members')).toBe(true);
		expect(isListingVisibilityAllowed('family', 'friends')).toBe(false);
	});

	it("rejects everything above 'private' for 'specific_space' assets (MVP simplification)", () => {
		expect(isListingVisibilityAllowed('specific_space', 'members')).toBe(false);
		expect(isListingVisibilityAllowed('specific_space', 'friends')).toBe(false);
	});

	it('agrees with the ceiling table for every combination', () => {
		for (const assetVisibility of ASSETUS_VISIBILITIES) {
			for (const listingVisibility of LISTING_VISIBILITIES) {
				expect(isListingVisibilityAllowed(assetVisibility, listingVisibility)).toBe(
					LISTING_VISIBILITY_CEILING[assetVisibility].includes(listingVisibility),
				);
			}
		}
	});
});
