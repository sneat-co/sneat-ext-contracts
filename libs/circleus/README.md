# @sneat/extension-circleus-contract

Public Circleus DTOs, contexts, service interfaces, and dependency-injection
tokens. Runtime implementation and UI stay in the private Circleus repository.

## Provenance

Promoted from `sneat-co/ext-circleus`
(`frontend/libs/extensions/template/contract`), commit
`b951948f483b6f1da80958fe7bbd9b78c050488b` (`origin/main`, 2026-08-26 read), per founder
decision 2026-08-26.

ext-circleus was scaffolded from the Sneat extension template and never
rebranded its contract lib: the package was still named
`@sneat/extension-template-contract` and every exported symbol still carried
`Template`/`template` naming. This promotes that stub into the family's first
real contract, with every `Template`/`template` token renamed to
`Circleus`/`circleus` (types, the `TEMPLATE_SERVICE` injection token, the
service interface, `ITemplateSpaceDbo`, and the `template-service.ts` /
`template-team.ts` source files).

This is a first-publish — no npm baseline exists for
`@sneat/extension-circleus-contract`.

## Cross-family dependencies

None.
