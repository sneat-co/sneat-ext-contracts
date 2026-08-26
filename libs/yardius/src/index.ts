// @sneat/extension-yardius-contract — frozen cross-repo contract surface for the yardius extension.
//
// Shared DTOs (listings, transaction-log entries, request/claim shapes, the
// listings-endpoint contract), lifecycle state enums with their legal
// transition tables, and the listing-visibility ceiling mapping are exported
// from here so that both the sneat-go backend models and the Sneat super-app
// extension libs resolve every shared model from this single package. No
// consumer may re-declare these.

export * from './lib/index.js';
