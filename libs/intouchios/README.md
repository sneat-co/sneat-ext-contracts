# @sneat/extension-intouchios-contract

Public contract for **InTouchKeeper** (intouchkeeper.app), the `intouchios`
Sneat extension. Its implementation lives in
[`sneat-co/intouchios`](https://github.com/sneat-co/intouchios).

**Provenance:** this was seeded on 2026-11-05 from
`sneat-co/sneat-ext-contract-template` (frontend `libs/extensions/template/contract`),
with `template` renamed to `intouchios`. It currently holds the template's
starter list DTOs, contexts and `INTOUCHIOS_SERVICE`, which the scaffolded
implementation compiles against. Relationship-memory types (contacts,
interactions, memories, follow-ups, monitored events) will replace them as
the spec lands.

Types and injection tokens only. It depends on foundational `@sneat/*`
packages, never on another extension.
