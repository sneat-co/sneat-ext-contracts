import assert from 'node:assert/strict';
import { spawnSync } from 'node:child_process';
import { readFileSync } from 'node:fs';
import { fileURLToPath } from 'node:url';
import test from 'node:test';
import { failedRequiredLegs } from './check.mjs';

const successful = {
  NX_RESULT: 'success',
  TIER_COHERENCE_RESULT: 'success',
  DISCOVER_GO_RESULT: 'success',
  RELEASE_CANDIDATE_POLICY_RESULT: 'success',
  GO_RESULT: 'success',
  GO_DIRS: '["calendarius"]',
};

test('accepts all successful legs', () => {
  assert.deepEqual(failedRequiredLegs(successful), []);
});

test('accepts a skipped Go matrix when discovery found no modules', () => {
  assert.deepEqual(failedRequiredLegs({ ...successful, GO_RESULT: 'skipped', GO_DIRS: '[]' }), []);
});

test('rejects skipped Go when discovery found modules', () => {
  assert.deepEqual(failedRequiredLegs({ ...successful, GO_RESULT: 'skipped' }), ['GO_RESULT']);
});

test('rejects missing, malformed, or non-string Go discovery output', () => {
  for (const value of [undefined, 'not-json', '{}', '[1]']) {
    assert.ok(failedRequiredLegs({ ...successful, GO_DIRS: value }).includes('GO_DIRS'));
  }
});

for (const name of Object.keys(successful)) {
  test(`rejects ${name} when failed, cancelled, skipped unexpectedly, or absent`, () => {
    for (const result of ['failure', 'cancelled', 'skipped', undefined]) {
      if (name === 'GO_RESULT' && result === 'skipped') continue;
      assert.deepEqual(
        failedRequiredLegs({ ...successful, [name]: result }),
        [name],
      );
    }
  });
}

test('CLI entrypoint succeeds and fails with the same environment contract', () => {
  const command = fileURLToPath(new URL('./check.mjs', import.meta.url));
  assert.equal(spawnSync(process.execPath, [command], { env: { ...process.env, ...successful } }).status, 0);
  assert.equal(spawnSync(process.execPath, [command], { env: { ...process.env, ...successful, NX_RESULT: 'failure' } }).status, 1);
});

test('workflow checks out the repository and runs the tested entrypoint', () => {
  const workflow = readFileSync('.github/workflows/ci.yml', 'utf8');
  const aggregate = workflow.slice(workflow.indexOf('  required-checks:'));
  assert.match(aggregate, /uses: actions\/checkout@v6/);
  assert.match(aggregate, /GO_DIRS: \$\{\{ needs\.discover-go\.outputs\.dirs \}\}/);
  assert.match(
    aggregate,
    /RELEASE_CANDIDATE_POLICY_RESULT: \$\{\{ needs\.release-candidate-policy\.result \}\}/,
  );
  assert.match(aggregate, /run: node tools\/required-checks\/check\.mjs/);
});
