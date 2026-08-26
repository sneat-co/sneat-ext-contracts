# @sneat/extension-communitycentrum-contract

Public NoticeBoard.cc / Community Centrum navigation-section interface,
service contract, and dependency-injection token. NoticeBoard.cc is a thin
composition layer over existing Sneat extensions (Calendarius, Contactus,
Assetus, Bookius, ToGethered, Competios); the runtime implementation and UI
stay in the private communitycentrum repository.

## Provenance

Migrated from `sneat-co/communitycentrum`
(`libs/extensions/communitycentrum/contract`), commit
`4ccc3801a48e8b875019bf2f368689170f30f8db` (`origin/main`, 2026-08-26 read).

This is a first-publish — no npm baseline exists for
`@sneat/extension-communitycentrum-contract`.

## Cross-family dependencies

None. The source mentions other extensions (Calendarius, Contactus, Assetus,
Bookius, ToGethered, Competios) only in a doc comment describing the product's
composition — there is no import of any other `@sneat/extension-*-contract`
package.
