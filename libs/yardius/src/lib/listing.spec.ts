import { describe, it, expect } from 'vitest';
import {
	IBorrowListingDto,
	IGiveawayListingDto,
	IListingDto,
	LISTING_TYPES,
} from './listing.js';

describe('LISTING_TYPES', () => {
	it("is exactly {'borrow', 'giveaway'} per REQ publish-listing", () => {
		expect(LISTING_TYPES).toEqual(['borrow', 'giveaway']);
	});
});

describe('IListingDto round-trip', () => {
	it('preserves a borrow Listing through JSON serialize/parse', () => {
		const listing: IBorrowListingDto = {
			id: 'listing-1',
			spaceID: 'space-x',
			assetID: 'asset-basketball',
			type: 'borrow',
			visibility: 'friends',
			state: 'requested',
		};

		const roundTripped: IListingDto = JSON.parse(JSON.stringify(listing));

		expect(roundTripped).toEqual(listing);
	});

	it('preserves a giveaway Listing through JSON serialize/parse', () => {
		const listing: IGiveawayListingDto = {
			id: 'listing-2',
			spaceID: 'space-x',
			assetID: 'asset-pram',
			type: 'giveaway',
			visibility: 'members',
			state: 'claimed',
		};

		const roundTripped: IListingDto = JSON.parse(JSON.stringify(listing));

		expect(roundTripped).toEqual(listing);
	});

	it("narrows `state` to the right lifecycle via the `type` discriminant", () => {
		const listing: IListingDto = {
			id: 'listing-3',
			spaceID: 'space-x',
			assetID: 'asset-drill',
			type: 'borrow',
			visibility: 'private',
			state: 'available',
		};

		// Compile-time check: inside the narrowed arm, `state` is a BorrowState.
		if (listing.type === 'borrow') {
			const state: IBorrowListingDto['state'] = listing.state;
			expect(state).toBe('available');
		}
	});

	it('references the asset by ID only — no duplicated Assetus fields (REQ publish-listing)', () => {
		const listing: IBorrowListingDto = {
			id: 'listing-4',
			spaceID: 'space-x',
			assetID: 'asset-tent',
			type: 'borrow',
			visibility: 'members',
			state: 'available',
		};

		const keys = Object.keys(JSON.parse(JSON.stringify(listing))).sort();
		expect(keys).toEqual(['assetID', 'id', 'spaceID', 'state', 'type', 'visibility']);
	});
});
