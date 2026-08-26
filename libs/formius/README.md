# @sneat/extension-formius-contract

Public TypeScript and Angular contract for the Sneat Formius extension.

It publishes `@sneat/extension-formius-contract`, including form/list DTOs,
contexts, service interfaces, and injection tokens. Private Formius runtime
and UI implementation stays in the `formius` repository.

## Provenance

Migrated from `sneat-co/ext-formius`
(`frontend/libs/extensions/formius/contract`), commit
`50e947781c84d1c2f7a74c6dca866fdeff116904` (`origin/main`, 2026-08-26 read).
The source tree is identical to the `v0.1.0`
git tag (`git diff v0.1.0 HEAD` over the contract path is empty) and to npm's
actual published latest, `0.1.0` (published 2026-07-14). No post-tag drift:
git tag, npm latest, and this seed are all API-identical.

`ext-formius`'s `backend/` and `typespec/` directories hold only template
placeholders (no Go contract code, no TypeSpec definitions) — this family has
no Go leg to migrate. This seed is npm-only.
