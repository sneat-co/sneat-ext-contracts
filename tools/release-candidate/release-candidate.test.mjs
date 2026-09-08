import assert from 'node:assert/strict';
import { execFileSync } from 'node:child_process';
import { mkdtempSync, mkdirSync, writeFileSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { join } from 'node:path';
import { afterEach, test } from 'node:test';
import {
  bumpVersion,
  classifyPublication,
  createReleaseMetadata,
  metadataPath,
  validateCiRevision,
  validateFinalCandidate,
  validateGeneratedTree,
} from './release-candidate.mjs';

const originalCwd = process.cwd();
afterEach(() => process.chdir(originalCwd));

test('binds one planned package and sibling Go tag to a two-parent checked merge', () => {
  const repo = newFixture();
  process.chdir(repo.path);
  const final = validateFinalCandidate({ checked: repo.checked });

  assert.equal(final.reviewedSource, repo.source);
  assert.equal(final.generatedCandidate, repo.candidate);
  assert.equal(final.checkedMain, repo.checked);
  assert.deepEqual(final.releases, [
    {
      family: 'contactus',
      project: 'contactus-contract',
      package: '@sneat/extension-contactus-contract',
      manifestPath: 'libs/contactus/package.json',
      changelogPath: 'libs/contactus/CHANGELOG.md',
      version: '0.12.9',
      npmTag: 'contactus-contract-v0.12.9',
      siblingGoTag: 'contactus/v0.12.9',
    },
  ]);
});

test('validates the actual generated pull-request head and checked main revision', () => {
  const repo = newFixture();
  process.chdir(repo.path);

  assert.equal(
    validateCiRevision({ event: 'pull_request', revision: repo.candidate, base: repo.source }).mode,
    'generated-candidate',
  );
  assert.equal(validateCiRevision({ event: 'push', revision: repo.checked }).mode, 'checked-main');

  writeFileSync('ordinary.txt', 'later source work\n');
  git(['add', 'ordinary.txt']);
  git(['commit', '-m', 'ordinary source work']);
  assert.equal(validateCiRevision({ event: 'push' }).mode, 'ordinary');
});

test('rejects silent plan consumption and package version changes after release metadata persists', () => {
  const removedPlan = newFixture();
  process.chdir(removedPlan.path);
  mkdirSync('.nx/version-plans', { recursive: true });
  writeFileSync('.nx/version-plans/next.md', '---\ncontactus-contract: patch\n---\n');
  git(['add', '.nx/version-plans/next.md']);
  git(['commit', '-m', 'plan next release']);
  git(['rm', '.nx/version-plans/next.md']);
  git(['commit', '-m', 'silently consume plan']);
  assert.throws(
    () => validateCiRevision({ event: 'push' }),
    /only be consumed by a bound release candidate/,
  );

  const changedVersion = newFixture();
  process.chdir(changedVersion.path);
  writeFileSync(
    'libs/contactus/package.json',
    '{"name":"@sneat/extension-contactus-contract","version":"0.12.10"}\n',
  );
  git(['add', 'libs/contactus/package.json']);
  git(['commit', '-m', 'silently change package version']);
  assert.throws(
    () => validateCiRevision({ event: 'push' }),
    /version changed outside a bound release candidate/,
  );
});

test('rejects package scope that was not present in the reviewed plan', () => {
  const repo = newFixture();
  process.chdir(repo.path);
  const metadata = structuredClone(repo.metadata);
  metadata.releases.push({
    ...metadata.releases[0],
    project: 'forged-contract',
    npmTag: 'forged-contract-v0.12.9',
  });

  assert.throws(
    () => validateGeneratedTree({ source: repo.source, candidate: repo.candidate, metadata }),
    /package scope differs/,
  );
});

test('rejects forged reviewed plan content', () => {
  const repo = newFixture();
  process.chdir(repo.path);
  const metadata = structuredClone(repo.metadata);
  metadata.plans[0].sha256 = '0'.repeat(64);

  assert.throws(
    () => validateGeneratedTree({ source: repo.source, candidate: repo.candidate, metadata }),
    /hash mismatch/,
  );
});

test('rejects a candidate presented against a different reviewed source', () => {
  const repo = newFixture();
  process.chdir(repo.path);

  assert.throws(
    () => validateGeneratedTree({ source: repo.checked, candidate: repo.candidate, metadata: repo.metadata }),
    /reviewed source does not match/,
  );
});

test('rejects a squash because it loses generated candidate ancestry', () => {
  const repo = newFixture({ merge: false });
  process.chdir(repo.path);
  git(['checkout', 'main']);
  git(['merge', '--squash', 'release']);
  git(['commit', '-m', 'squashed release']);

  assert.throws(() => validateFinalCandidate(), /two-parent merge/);
});

test('rejects protected release files amended after candidate generation', () => {
  const repo = newFixture({ merge: false });
  process.chdir(repo.path);
  git(['checkout', 'release']);
  writeFileSync('libs/contactus/CHANGELOG.md', '# changed after review\n');
  git(['add', '.']);
  git(['commit', '-m', 'drift release output']);
  git(['checkout', 'main']);
  git(['merge', '--no-ff', 'release', '-m', 'merge drifted release']);

  assert.throws(() => validateFinalCandidate(), /second parent is the generated candidate/);
});

test('rejects target drift before the protected candidate merge', () => {
  const repo = newFixture({ merge: false });
  process.chdir(repo.path);
  git(['checkout', 'main']);
  writeFileSync('target-drift.txt', 'drift\n');
  git(['add', 'target-drift.txt']);
  git(['commit', '-m', 'target drift']);
  git(['merge', '--no-ff', 'release', '-m', 'merge after target drift']);

  assert.throws(() => validateFinalCandidate(), /two-parent merge/);
});

test('rejects merge resolution that changes the generated candidate tree', () => {
  const repo = newFixture({ merge: false });
  process.chdir(repo.path);
  git(['checkout', 'main']);
  git(['merge', '--no-commit', '--no-ff', 'release']);
  writeFileSync('pnpm-lock.yaml', 'changed during merge\n');
  git(['add', 'pnpm-lock.yaml']);
  git(['commit', '-m', 'merge with changed generated output']);

  assert.throws(() => validateFinalCandidate(), /tree differs/);
});

test('rejects an amended generated pull-request head through the live CI path', () => {
  const repo = newFixture({ merge: false });
  process.chdir(repo.path);
  git(['checkout', 'release']);
  writeFileSync('late-change.txt', 'not generated\n');
  git(['add', 'late-change.txt']);
  git(['commit', '-m', 'late candidate amendment']);

  assert.throws(
    () => validateCiRevision({ event: 'pull_request', base: repo.source }),
    /direct child/,
  );
});

test('rejects duplicate candidate metadata for the same reviewed source', () => {
  const repo = newFixture({ merge: false });
  process.chdir(repo.path);
  git(['checkout', 'release']);
  const metadata = JSON.parse(git(['show', `HEAD:${metadataPath}`], false));
  writeFileSync(metadataPath, `${JSON.stringify(metadata)}\n`);
  git(['add', metadataPath]);
  git(['commit', '--allow-empty', '-m', 'duplicate candidate']);
  git(['checkout', 'main']);
  git(['merge', '--no-ff', 'release', '-m', 'merge duplicate candidate']);

  assert.throws(() => validateFinalCandidate(), /exactly one generated release candidate/);
});

test('preserves partial npm retry and refuses a wrong-SHA tag', () => {
  const release = {
    project: 'contactus-contract',
    package: '@sneat/extension-contactus-contract',
    version: '0.12.9',
    npmTag: 'contactus-contract-v0.12.9',
  };
  assert.deepEqual(
    classifyPublication({
      releases: [release],
      checkedMain: 'a'.repeat(40),
      remoteTags: {
        'contactus-contract-v0.12.9': {
          rawObject: 'c'.repeat(40),
          peeledCommit: 'a'.repeat(40),
        },
      },
      npmIntegrities: {},
      localIntegrities: {
        '@sneat/extension-contactus-contract@0.12.9': 'sha512-local',
      },
    }),
    { publishProjects: ['contactus-contract'], tagsToPush: [] },
  );
  assert.deepEqual(
    classifyPublication({
      releases: [release],
      checkedMain: 'a'.repeat(40),
      remoteTags: { 'unrelated-contract-v9.9.9': null },
      npmIntegrities: { '@sneat/unrelated@9.9.9': null },
      localIntegrities: {
        '@sneat/extension-contactus-contract@0.12.9': 'sha512-local',
      },
    }),
    { publishProjects: ['contactus-contract'], tagsToPush: ['contactus-contract-v0.12.9'] },
  );
  assert.deepEqual(
    classifyPublication({
      releases: [release],
      checkedMain: 'a'.repeat(40),
      remoteTags: {},
      npmIntegrities: {
        '@sneat/extension-contactus-contract@0.12.9': 'sha512-receipt',
      },
      localIntegrities: {
        '@sneat/extension-contactus-contract@0.12.9': 'sha512-receipt',
      },
    }),
    { publishProjects: [], tagsToPush: ['contactus-contract-v0.12.9'] },
  );
  assert.throws(
    () => classifyPublication({
      releases: [release],
      checkedMain: 'a'.repeat(40),
      remoteTags: {
        'contactus-contract-v0.12.9': {
          rawObject: 'c'.repeat(40),
          peeledCommit: 'b'.repeat(40),
        },
      },
      npmIntegrities: {},
      localIntegrities: {
        '@sneat/extension-contactus-contract@0.12.9': 'sha512-local',
      },
    }),
    /not checked main/,
  );
  assert.throws(
    () => classifyPublication({
      releases: [release],
      checkedMain: 'a'.repeat(40),
      remoteTags: {
        'contactus-contract-v0.12.9': {
          rawObject: 'a'.repeat(40),
          peeledCommit: null,
        },
      },
      npmIntegrities: {},
      localIntegrities: {
        '@sneat/extension-contactus-contract@0.12.9': 'sha512-local',
      },
    }),
    /must be an annotated tag/,
  );
  assert.throws(
    () => classifyPublication({
      releases: [release],
      checkedMain: 'a'.repeat(40),
      remoteTags: {},
      npmIntegrities: {
        '@sneat/extension-contactus-contract@0.12.9': 'sha512-registry',
      },
      localIntegrities: {
        '@sneat/extension-contactus-contract@0.12.9': 'sha512-local',
      },
    }),
    /differs from the checked-main artifact/,
  );
});

function newFixture({ merge = true } = {}) {
  const path = mkdtempSync(join(tmpdir(), 'contract-release-'));
  process.chdir(path);
  git(['init', '-b', 'main']);
  git(['config', 'user.email', 'test@example.test']);
  git(['config', 'user.name', 'Release Test']);
  mkdirSync('.nx/version-plans', { recursive: true });
  mkdirSync('libs/contactus', { recursive: true });
  mkdirSync('contactus', { recursive: true });
  writeFileSync('.nx/version-plans/contactus.md', '---\ncontactus-contract: patch\n---\n\nRelease helper.\n');
  writeFileSync('libs/contactus/project.json', '{"name":"contactus-contract"}\n');
  writeFileSync('libs/contactus/package.json', '{"name":"@sneat/extension-contactus-contract","version":"0.12.8"}\n');
  writeFileSync('libs/contactus/CHANGELOG.md', '# Changelog\n');
  writeFileSync('contactus/go.mod', 'module example.test/contactus\n');
  git(['add', '.']);
  git(['commit', '-m', 'reviewed source']);
  const source = git(['rev-parse', 'HEAD']);
  git(['checkout', '-b', 'release']);
  writeFileSync('libs/contactus/package.json', '{"name":"@sneat/extension-contactus-contract","version":"0.12.9"}\n');
  writeFileSync('libs/contactus/CHANGELOG.md', '## 0.12.9\n\nRelease helper.\n');
  git(['rm', '.nx/version-plans/contactus.md']);
  git(['add', '.']);
  git(['commit', '-m', 'generated candidate']);
  git(['tag', 'contactus-contract-v0.12.9']);
  const generated = git(['rev-parse', 'HEAD']);
  const metadata = createReleaseMetadata({ source, candidate: generated });
  mkdirSync('.nx', { recursive: true });
  writeFileSync(metadataPath, `${JSON.stringify(metadata, null, 2)}\n`);
  git(['add', metadataPath]);
  git(['commit', '--amend', '--no-edit']);
  const candidate = git(['rev-parse', 'HEAD']);
  git(['tag', '--force', 'contactus-contract-v0.12.9', candidate]);
  validateGeneratedTree({ source, candidate, metadata });
  git(['checkout', 'main']);
  let checked = source;
  if (merge) {
    git(['merge', '--no-ff', 'release', '-m', 'merge release candidate']);
    checked = git(['rev-parse', 'HEAD']);
  }
  return { path, source, candidate, checked, metadata };
}

function git(args, trim = true) {
  const value = execFileSync('git', args, { encoding: 'utf8' });
  return trim ? value.trim() : value;
}

test('mirrors nx major-zero semantics when computing the expected plan version', () => {
  assert.equal(bumpVersion('0.2.6', 'major'), '0.3.0');
  assert.equal(bumpVersion('0.2.6', 'minor'), '0.2.7');
  assert.equal(bumpVersion('0.2.6', 'patch'), '0.2.7');
  assert.equal(bumpVersion('1.4.2', 'major'), '2.0.0');
  assert.equal(bumpVersion('1.4.2', 'minor'), '1.5.0');
  assert.equal(bumpVersion('1.4.2', 'patch'), '1.4.3');
});
