// The Assetus asset-visibility values, mirrored here as wire values so the
// listing-visibility ceiling (see `listing-visibility.ts`) can be encoded as
// shared data with zero runtime dependencies.
//
// Source of truth for the VALUES is the pinned assetus-mvp contract
// (`AssetVisibility` in the `@sneat/extension-assetus-contract` lib):
// spec display names `Private` / `Family` / `Friends` / `Friends of Friends` /
// `Specific Space` / `Public` are persisted as the lowercase snake_case wire
// values below. Yardius only READS asset visibility (REQ
// assetus-write-boundary) — this type exists so both the Yardius frontend and
// backend consume the ceiling mapping from this single package.
export const ASSETUS_VISIBILITIES = [
	'private',
	'family',
	'friends',
	'friends_of_friends',
	'specific_space',
	'public',
] as const;

export type AssetusVisibility = (typeof ASSETUS_VISIBILITIES)[number];
