// Status of a borrow request or give-away claim:
//   - 'pending'   — awaiting the owning Space's decision;
//   - 'approved'  — the owner approved this request (borrow) or selected this
//     claimant (give-away); all other pending requests/claims on the Listing
//     are automatically declined (REQ borrow-lifecycle /
//     REQ giveaway-lifecycle);
//   - 'declined'  — declined by the owner, auto-declined because another
//     request was approved, cancelled with the Listing (REQ owner-cancel), or
//     auto-declined on friendship removal (REQ friendship-removal);
//   - 'withdrawn' — the requester withdrew their own pending request
//     (REQ borrow-request).
export const LISTING_REQUEST_STATUSES = [
	'pending',
	'approved',
	'declined',
	'withdrawn',
] as const;

export type ListingRequestStatus = (typeof LISTING_REQUEST_STATUSES)[number];

// A request to borrow (on a `borrow` Listing) or a claim (on a `giveaway`
// Listing) — one shared shape; which it is follows from the Listing's `type`.
//
// Per REQ borrow-request, requesters are members of a Space that can see the
// Listing per its visibility, EXCLUDING members of the owning Space, and
// owning-Space members see each pending request with the requester's name and
// Space — hence the requester identity + display fields below.
export interface IListingRequestDto {
	readonly id: string;

	// The requester/claimant: their Space plus their member (contactus
	// contact) ID within that Space. Never a member of the owning Space.
	readonly requesterSpaceID: string;
	readonly requesterMemberID: string;

	// Display labels resolved by the backend (contactus member identity /
	// space title) so owning-Space members can render "who is asking" without
	// a cross-space read. Optional: absent when the caller can resolve them
	// locally.
	readonly requesterTitle?: string;
	readonly requesterSpaceTitle?: string;

	readonly status: ListingRequestStatus;

	// ISO timestamp string, server-assigned at creation.
	readonly createdAt: string;
}
