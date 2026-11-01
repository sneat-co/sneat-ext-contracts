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

## Structured goal values

Additive goal-definition and goal-terms exports describe reviewed calendar
windows, version-bound evidence/create fields, filtered semantic pages,
full amendment proposals and owner-only version/change history. They add no
payment behavior, result evaluator or runtime service implementation. Keep ISO
instants as strings to preserve server timestamp precision. Supported and
unsupported semantic pages are distinct; clients must honor server action and
support flags instead of inferring consent from progress or approval counts.

Actor and bot-channel identity come from trusted host context, never body
fields. Space membership and referee capabilities do not grant owner-history
listing. Publication and actual host integration are separate release gates;
the existing service alias/token and package versions remain unchanged here.
