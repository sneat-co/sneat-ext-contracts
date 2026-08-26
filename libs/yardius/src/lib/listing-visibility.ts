import { AssetusVisibility } from './assetus-visibility.js';

// A Listing's own visibility, per REQ listing-visibility-ceiling (spec display
// names `Private` / `Members` / `Friends`, persisted as the wire values
// below):
//   - 'private' — visible only to members of the owning Space;
//   - 'members' — visible to members of the owning Space (reserved distinction
//     from 'private' for future member-subset control; in MVP the two resolve
//     to the same audience, and 'members' is the default);
//   - 'friends' — additionally visible to members of befriended Spaces.
export const LISTING_VISIBILITIES = ['private', 'members', 'friends'] as const;

export type ListingVisibility = (typeof LISTING_VISIBILITIES)[number];

// Per REQ listing-visibility-ceiling: `members` is the default listing
// visibility.
export const DEFAULT_LISTING_VISIBILITY: ListingVisibility = 'members';

export function listingVisibilityLabel(visibility: ListingVisibility): string {
	switch (visibility) {
		case 'private':
			return 'Private';
		case 'members':
			return 'Members';
		case 'friends':
			return 'Friends';
	}
}

// The asset's Assetus visibility is a CEILING on the listing's visibility,
// per this exact mapping from REQ listing-visibility-ceiling:
//
//   | Asset visibility (Assetus)               | Permitted listing visibilities |
//   |------------------------------------------|--------------------------------|
//   | Private                                  | Private                        |
//   | Family                                   | Private, Members               |
//   | Specific Space                           | Private (MVP simplification)   |
//   | Friends, Friends of Friends, Public      | Private, Members, Friends      |
//
// Encoded as data so frontend and backend enforce the same source of truth.
// Publishing with a listing visibility not permitted here MUST be rejected
// with an actionable error telling the user to raise the asset's visibility
// in Assetus first.
export const LISTING_VISIBILITY_CEILING: Readonly<
	Record<AssetusVisibility, readonly ListingVisibility[]>
> = {
	private: ['private'],
	family: ['private', 'members'],
	specific_space: ['private'], // MVP simplification.
	friends: ['private', 'members', 'friends'],
	friends_of_friends: ['private', 'members', 'friends'],
	public: ['private', 'members', 'friends'],
};

// The listing visibilities permitted for an asset with the given Assetus
// visibility. Pure lookup over `LISTING_VISIBILITY_CEILING`.
export function allowedListingVisibilities(
	assetVisibility: AssetusVisibility,
): readonly ListingVisibility[] {
	return LISTING_VISIBILITY_CEILING[assetVisibility];
}

// True when a listing with `listingVisibility` may be published for an asset
// whose Assetus visibility is `assetVisibility`.
export function isListingVisibilityAllowed(
	assetVisibility: AssetusVisibility,
	listingVisibility: ListingVisibility,
): boolean {
	return LISTING_VISIBILITY_CEILING[assetVisibility].includes(listingVisibility);
}
