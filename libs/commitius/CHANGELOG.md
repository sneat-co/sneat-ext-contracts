## 0.1.4 (2026-11-01)

### 🩹 Fixes

- Trackus adds canonical tracker metadata, public Firebase timestamp types, and typed stable-entry requests for personal progress tracking. Its Angular support moves to the verified Angular 22 line; the former Angular 21 peer graph cannot resolve with the supported shared packages. Existing published versions remain immutable. ([33b055c](https://github.com/sneat-co/sneat-ext-contracts/commit/33b055c))

  Commitius consumes the compatible published template contract without changing its public exports. Assetus widens its Angular peer range to include Angular 22. The qualified consumer uses Angular 22; this does not assert whole-graph Angular 21 support.

  Docus follows the existing native dependent-release policy for Assetus updates. Its patch republishes the unchanged public API with the generated Assetus dependency update; this is an explicit release-plan entry, not a new application feature.

### ❤️ Thank You

- Alexander Trakhimenok

## 0.1.3 (2026-08-31)

### 🩹 Fixes

- chore(contracts): release peers targeting Sneat Libraries 0.27.6 ([a5f58ec](https://github.com/sneat-co/sneat-ext-contracts/commit/a5f58ec))

### ❤️ Thank You

- Alexander Trakhimenok

## 0.1.2 (2026-08-30)

### 🩹 Fixes

- fix(contracts): accept the Angular 22 and current Sneat Libraries consumers ([98c85ad](https://github.com/sneat-co/sneat-ext-contracts/commit/98c85ad))

  The migrated contract packages still advertised their pre-migration peer
  ranges. Align the metadata with the Angular 22 product workspaces and permit
  the existing and current Sneat Libraries contract-compatible release lines.

### ❤️ Thank You

- Alexander Trakhimenok

## 0.1.1 (2026-08-26)

### 🩹 Fixes

- feat(commitius-contract): migrate commitius contract into sneat-ext-contracts (first-publish) ([9aee482](https://github.com/sneat-co/sneat-ext-contracts/commit/9aee482))

  First publish of @sneat/extension-commitius-contract into sneat-ext-contracts monorepo.
  Source: sneat-co/commitius@03ee7029ad38553729aaf78c2feba22e219fdb38 (libs/extensions/commitius/contract)

  Resulting version: 0.1.1 (first-publish patch bump from declared 0.1.0)

### ❤️ Thank You

- Alexander Trakhimenok
- Claude Fable 5