# @sneat/extension-taxus-contract

Public Taxus DTOs, contexts, service interfaces, and dependency-injection
tokens. Tax calculation, runtime services, and UI implementations remain in the
private Taxus repository.

## Provenance

Migrated from `sneat-co/ext-taxus`
(`frontend/libs/extensions/taxus/contract`), commit
`8f2433850169d7e135c22816818c16e94fda3a73` (`origin/main`, 2026-08-26 read).
Source has been API-identical since the lib's first commit
(`cc2741b`, "feat(contract): publish Taxus contract") — the `src/` tree never
changed across ext-taxus's whole history, including the never-published
`v0.0.3` git tag. This migration carries the `v0.0.3` tag's only real content
(this README's current wording, from "fix(contract): publish Taxus package
publicly") forward; the tag added no DTO, context, or service-interface
changes, so npm's actual latest (`0.0.2`) and this seed are API-identical.
