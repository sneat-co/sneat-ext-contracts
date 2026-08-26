# @sneat/extension-togethered-contract

Public ToGethered DTOs, contexts, service interfaces, and dependency-injection
tokens. List/list-item runtime services and UI implementations remain in the
private ToGethered repository (`libs/extensions/togethered/internal` and
`libs/extensions/togethered/shared`).

## Provenance

Migrated from `sneat-co/togethered`
(`libs/extensions/togethered/contract`), commit
`a4d6c9eb030e12946a959d0a89ee77bac8d0496c` (`origin/main`, 2026-08-26 read).

This is a first-publish — no npm baseline exists for
`@sneat/extension-togethered-contract`.

## Cross-family dependencies

None. The source repo also carries `@sneat/extension-togethered-internal` and
`@sneat/extension-togethered-shared` workspace libs, but the contract imports
neither — only `@angular/core`, `rxjs`, and the foundational `@sneat/core`,
`@sneat/data`, `@sneat/dto`, `@sneat/space-models` packages already peer-depended
on by other families in this repo.
