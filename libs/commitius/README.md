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
Publication requires strict packed-consumer qualification of the complete
Angular 22 peer graph; package build success alone does not establish that
compatibility or Angular 21 support.

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

Schema2 adds explicitly reviewed activity/role context while schema1 remains a
separate supported input/read type. Standard push-up/pull-up names reuse the
public Trackus vocabulary; context is meaning, not a tracker/source binding.
Do not infer or upgrade it from a goal title. The server remains authoritative
for safe integers, bounded text, calendar coherence and duplicate axes.

The bounded owner-page types describe fixed identity/date/compatible-target
orders, current saved-example selection, ID-only unavailable rows and explicit
live coverage. Cursor metadata is non-confidential and grants no authority.
Target values are promises, not activity, verified outcomes or achievements.
Derived caches, raw Space/snapshot/source IDs and reconciliation operations are
not public list fields. These types do not mount a route or enable a release.

The current contextual-goal/page candidate passes local lint, tests, package
build and a separate TypeScript gate. Its exact positive goal examples also
pass the owning Go validator. Strict packed public-entry compilation passes
in ibetican's cached app peer context: template0.0.29, shared/core-public0.27.46
and Angular22.2.4, with skipLibCheck=false and unchanged dependencies.

The contracts workspace's older frozen peer context still fails that strict
check: it resolves core-public0.26.6 and lacks required core-public/telemetry
declarations. The successful app-context check does not qualify that older
graph, a fresh install, every consumer, app workflows or deployment. No check
was suppressed, package version changed or artifact published. Protected
release/registry verification and actual app integration remain required.
