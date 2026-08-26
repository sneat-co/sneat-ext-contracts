import { describe, it, expect } from 'vitest';
import {
	IGetListingsRequest,
	IGetListingsResponse,
	IListingsSpaceGroup,
} from './listings-endpoint.js';

describe('IGetListingsRequest', () => {
	it('carries no parameters — the caller is auth-implied (REQ backend-mediated-query)', () => {
		const request: IGetListingsRequest = {};

		const roundTripped: IGetListingsRequest = JSON.parse(JSON.stringify(request));

		expect(roundTripped).toEqual({});
		expect(Object.keys(roundTripped)).toEqual([]);
	});
});

describe('IGetListingsResponse round-trip', () => {
	it('preserves groups labelled by owning Space through JSON serialize/parse (AC listings-endpoint-scope)', () => {
		// The caller is a member of Space Y, which is befriended with Space X:
		// the response carries Y's own listings plus X's friends-visible one,
		// each group labelled with its owning Space.
		const friendGroup: IListingsSpaceGroup = {
			spaceID: 'space-x',
			spaceTitle: "Emma's Family",
			reason: 'friend',
			listings: [
				{
					id: 'listing-1',
					spaceID: 'space-x',
					assetID: 'asset-basketball',
					type: 'borrow',
					visibility: 'friends',
					state: 'available',
				},
			],
		};
		const memberGroup: IListingsSpaceGroup = {
			spaceID: 'space-y',
			spaceTitle: 'Our Family',
			reason: 'member',
			listings: [
				{
					id: 'listing-2',
					spaceID: 'space-y',
					assetID: 'asset-pram',
					type: 'giveaway',
					visibility: 'private',
					state: 'claimed',
				},
			],
		};
		const response: IGetListingsResponse = {
			groups: [memberGroup, friendGroup],
		};

		const roundTripped: IGetListingsResponse = JSON.parse(JSON.stringify(response));

		expect(roundTripped).toEqual(response);
	});

	it('round-trips an empty response (no member spaces with open listings, no friend listings)', () => {
		const response: IGetListingsResponse = { groups: [] };

		const roundTripped: IGetListingsResponse = JSON.parse(JSON.stringify(response));

		expect(roundTripped).toEqual(response);
	});
});
