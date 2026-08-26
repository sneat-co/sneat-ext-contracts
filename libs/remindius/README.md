# @sneat/extension-remindius-contract

Public Remindius DTOs, contexts, service interfaces, and dependency-injection
tokens. Reminder/list scheduling, runtime services, and UI implementations
remain in the private Remindius repository.

## Provenance

Migrated from `sneat-co/ext-remindius`
(`frontend/libs/extensions/remindius/contract`), commit
`c8f674f48dc5ff78766a55364303b9e0cc4bff65` (`origin/main`, 2026-08-26 read).

npm's actual latest is `0.2.0` (REQ `version-continuity-from-npm`): the
source repo's git history never bumped `package.json` past `0.1.0` and
carries no git tags at all, but `0.2.0` is published on the registry.
`npm pack` proves `0.1.0` and `0.2.0` are byte-identical on `.d.ts` — only
the `version` field in `package.json` changed — and the current git tree's
`src/` matches that same API exactly. Disk here is seeded at `0.2.0` to
match npm reality.

One post-`0.2.0` tree change exists in git and is deliberately **not**
carried forward: commit `f9e13a6` (2026-08-05, an unrelated ng-packagr/CI
publish-pipeline fix, PR #11) raised the `@sneat/{core,data,dto,space-models}`
peerDependencies floor from `^0.22.1` to `^0.24.0`. That change was never
published — npm's actual `0.2.0` `package.json` still declares `^0.22.1` —
and no git tag marks it as an intended release, so this seed keeps
`^0.22.1` to match what npm actually shipped.

The migrated `dto/list.spec.ts` unit test is carried forward with a new,
minimal `vitest.config.ts` (the source repo's original config carried
Ionic/jsdom aliasing needed only by an unrelated service spec that is not
part of this contract's public surface).
