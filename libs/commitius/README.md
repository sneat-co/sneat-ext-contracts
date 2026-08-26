# @sneat/extension-commitius-contract

Public Commitius DTOs, service interfaces, and dependency-injection tokens.
Commitment tracking, ledger, and UI implementations remain in the private
Commitius repository.

## Provenance

Migrated from `sneat-co/commitius`
(`libs/extensions/commitius/contract`), commit
`03ee7029ad38553729aaf78c2feba22e219fdb38` (`origin/main`, 2026-08-26 read).

This is a first-publish — no npm baseline exists for
`@sneat/extension-commitius-contract`.

## Cross-family dependency

This contract's entire public surface specializes
`@sneat/extension-template-contract` (`export * from` plus two type aliases
and one `InjectionToken`) — declared as a plain npm `peerDependency`, pinned
`^0.1.0` to match the source repo's own pin, exactly as it was in
`sneat-co/commitius`. `@sneat/extension-template-contract` is a real,
already-published npm package (versions 0.0.2, 0.0.3, 0.1.0 on the public
registry); this is not workspace/runtime code vendored from the product
repo.
