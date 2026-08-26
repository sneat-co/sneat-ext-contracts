# Module boundaries

This repo holds every `@sneat/extension-<family>-contract` npm package as a
sibling Nx library under `libs/<family>/`. That proximity is exactly what
makes an undeclared contract→contract dependency easy to introduce by
accident — this document plus `eslint.config.mjs`'s
`@nx/enforce-module-boundaries` `depConstraints` is how the repo stops that
mechanically instead of by convention alone.

Source: `spec/features/ext-contracts-monorepo/README.md`
(`sneat-co/sneat-libs`), REQs `explicit-cross-contract-boundaries` and
`ownership-test-carries-over`.

## Default: no cross-contract imports

Every family lib is tagged `family:<name>` and `layer:contract` in its
`project.json`. The default posture is that a family may only depend on
itself — no contract imports another contract unless an entry below (and a
matching `depConstraints` entry in `eslint.config.mjs`) says so explicitly.

Contract → runtime and contract → app dependencies are **forbidden
absolutely**, with no allow-list exception — a contract lib never depends on
its extension's implementation or on any product app.

## Ownership test (carries over verbatim from `extension-contract-repo`)

> An interface or type belongs in a contract only if its entire signature is
> expressible in that extension's own types plus foundational/core types. If
> any part of a signature references a **consumer's** types, that interface
> is the consumer's contract, not this extension's — it stays consumer-owned
> and does not move into the contract.

Both cross-extension interaction directions carry over unchanged:

- **facade-call-in** — the facade interface and its DTOs for *calling into* an
  extension live in that extension's contract. The extension's own runtime
  provides the implementation, wired by the app at bootstrap (an Angular
  `InjectionToken` provider on the frontend; registration on the backend). A
  consumer imports only the contract and never the implementation.
- **caller-satisfied-callback** — for the inverse direction (the extension
  needs behaviour or data *from the caller*), the callback interface is also
  declared in the extension's contract, but the **consumer** supplies the
  implementation and passes it in at call time.

## Allowed cross-contract edges

The measured import graph (2026-08-26 fleet audit: 934 contract imports, 84%
intra-family / 16% cross-family, contactus receiving 81 of 148 cross-family
imports) names **contactus** and **assetus** as the only platform-like
targets expected to gain further entries, arriving in Phase 2 of the
migration plan.

When a family needs to depend on another family's contract, add a row here
AND the matching `depConstraints` entry in `eslint.config.mjs` in the same
change:

| Source family | Allowed target | Why | Added in |
| --- | --- | --- | --- |
| commitius | template | commitius specializes the maintained list/template extension; `libs/commitius/src/index.ts` re-exports `@sneat/extension-template-contract` types (`ITemplateService`, `ITemplateSpaceDbo`) and depends on it as a peer. `template` is not yet migrated into this monorepo (still external npm), so this edge is a forward declaration for when it lands. | batch 3 (families: assetus, listus, debtus, sportius, contactus, calendarius, commitius, communitycentrum, togethered) |

## Adding a new family lib

1. Generate the lib at `libs/<family>/` (Nx project name `<family>-contract`,
   npm name `@sneat/extension-<family>-contract`).
2. Tag it `['family:<family>', 'layer:contract']` in `project.json`.
3. Add the default (self-only) `depConstraints` entry in `eslint.config.mjs`.
4. Add the family to `contracts.json` so the tier-coherence check
   (`tools/tier-coherence/`) starts covering it.
5. Only widen the constraint — and this table — when a specific cross-family
   edge is an approved, reviewed exception.
