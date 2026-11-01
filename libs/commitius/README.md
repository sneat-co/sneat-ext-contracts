# @sneat/extension-commitius-contract

Public Commitius DTOs, service interfaces, and dependency-injection tokens.
Commitment tracking, ledger, and UI implementations remain in the private
Commitius repository.

## Provenance

Migrated from `sneat-co/commitius`
(`libs/extensions/commitius/contract`), commit
`03ee7029ad38553729aaf78c2feba22e219fdb38` (`origin/main`, 2026-08-26 read).

## Cross-family dependency

The public contract reexports `@sneat/extension-template-contract` and adds
Commitius-owned type aliases and `COMMITIUS_SERVICE`. Its peer and build-time
reference use the existing published template 0.0.29, whose Angular 22 graph
uses shared Sneat 0.27.46. Keep the reexports and token identity intact; this
compatibility correction does not add or replace service methods.

Source package versions stay unchanged until the protected release process.
Strict packed-consumer validation qualifies the complete Angular 22 graph;
it does not establish Angular 21 support for this graph.
