import type { IAvailableListingDto, IListingDto } from './listing.js';

// The listings-endpoint request/response contract, per REQ
// backend-mediated-query: listings available to a user are served ONLY by
// this backend endpoint — web/mobile clients never read another Space's
// listings directly from Firestore.

// The request carries no parameters: the caller is identified by the request
// auth (the backend resolves the caller's member-Spaces and their friend
// edges server-side).
export type IGetListingsRequest = Record<string, never>;

// Why a group of listings is visible to the caller:
//   - 'member' — the caller is a member of the owning Space (all of that
//     Space's open Listings are returned);
//   - 'friend' — the owning Space is befriended with one of the caller's
//     Spaces (only its `friends`-visible open Listings are returned).
// There is deliberately no third arm: the endpoint never returns listings via
// friends-of-friends traversal or from non-friend Spaces.
export type ListingsGroupReason = 'member' | 'friend';

// Open Listings of one owning Space, labelled for relationship-first
// presentation (e.g. "Emma's Family is lending: Basketball").
export interface IListingsSpaceGroup {
	// The owning Space every listing in `listings` belongs to.
	readonly spaceID: string;

	// Display title of the owning Space, resolved by the backend so the
	// client can label the group without a cross-space read.
	readonly spaceTitle: string;

	readonly reason: ListingsGroupReason;

	readonly listings: readonly IListingDto[];
}

// The endpoint's response: for the calling user, (a) open Listings of every
// Space the caller is a member of and (b) `friends`-visible open Listings of
// every Space befriended with any of the caller's Spaces — grouped by the
// owning Space. Nothing else.
export interface IGetListingsResponse {
	/** Grouped response used by the relationship-first listing endpoint. */
	readonly groups?: readonly IListingsSpaceGroup[];
	/** Flat response used by the existing `get_available_listings` facade. */
	readonly listings?: readonly IAvailableListingDto[];
}
