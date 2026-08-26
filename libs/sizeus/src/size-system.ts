// Declared size system for a size record's `value`. Per REQ size-systems the
// value is stored verbatim in whichever system the caller declares — it is
// never normalized or converted between systems.
//
// The type is intentionally OPEN (mirrors the `SneatApp` idiom in
// sneat-libs/libs/core/src/lib/app.service.ts): the `string & {}` arm accepts
// any system string so a niche/regional system can be recorded without a
// contract release, while the explicit first-party literals below are kept
// for editor autocomplete and documentation. Do not remove the open arm.
export type SizeSystem =
	| 'EU'
	| 'UK'
	| 'US'
	| 'JP'
	| 'Mondopoint'
	| 'cm' // centimetres
	| 'in' // inches
	| 'alpha' // XS / S / M / L / XL / …
	// Open arm: accept any declared system string while keeping the literal
	// suggestions above extensible without a contract release.
	| (string & {});
