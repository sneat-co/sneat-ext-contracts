import { describe, it, expect } from 'vitest';
import { IListingRequestDto, LISTING_REQUEST_STATUSES } from './listing-request.js';

describe('LISTING_REQUEST_STATUSES', () => {
	it('is exactly {pending, approved, declined, withdrawn}', () => {
		expect(LISTING_REQUEST_STATUSES).toEqual([
			'pending',
			'approved',
			'declined',
			'withdrawn',
		]);
	});
});

describe('IListingRequestDto round-trip', () => {
	it('preserves every field, including optional display labels, through JSON serialize/parse', () => {
		const request: IListingRequestDto = {
			id: 'request-1',
			requesterSpaceID: 'space-y',
			requesterMemberID: 'member-bob',
			requesterTitle: 'Bob',
			requesterSpaceTitle: "Emma's Family",
			status: 'pending',
			createdAt: '2026-07-02T09:00:00.000Z',
		};

		const roundTripped: IListingRequestDto = JSON.parse(JSON.stringify(request));

		expect(roundTripped).toEqual(request);
	});

	it('preserves the shape with optional fields absent (not merely undefined)', () => {
		const request: IListingRequestDto = {
			id: 'request-2',
			requesterSpaceID: 'space-y',
			requesterMemberID: 'member-carol',
			status: 'withdrawn',
			createdAt: '2026-07-01T18:30:00.000Z',
		};

		const roundTripped: IListingRequestDto = JSON.parse(JSON.stringify(request));

		expect(roundTripped).toEqual(request);
		expect('requesterTitle' in roundTripped).toBe(false);
		expect('requesterSpaceTitle' in roundTripped).toBe(false);
	});
});
