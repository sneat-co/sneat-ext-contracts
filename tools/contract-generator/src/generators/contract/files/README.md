# <%= npmName %>

TODO: one-paragraph purpose — what the `<%= family %>` extension is, and what
lives in this contract (DTOs, contexts, service interfaces, injection tokens)
versus what stays in the private `<%= family %>` runtime/UI repo. See
`docs/boundaries.md`'s ownership test if unsure whether a type belongs here.

## Provenance

TODO — keep exactly one of the following and delete the rest:

- Migrated from `sneat-co/ext-<%= family %>` (`<path-in-old-repo>`), commit
  `<sha>` (`origin/main`, `<date>` read). State whether API parity against
  npm's actual latest was verified (REQ `version-continuity-from-npm`,
  `spec/features/ext-contracts-monorepo/README.md` on `sneat-co/sneat-libs`).
- Extracted directly from `sneat-co/sneat-libs`
  (`libs/extensions/<%= family %>/contract`), commit `<sha>`.
- New contract authored directly in this monorepo — no prior repo.
