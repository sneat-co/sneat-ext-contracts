import assert from 'node:assert/strict';
import test from 'node:test';

import { createDependencyPlan, isRegistryNotFound } from './check.mjs';

test('uses latest for published contracts and a local package only for an unpublished family', () => {
  const { dependencies, bootstrapped } = createDependencyPlan(
    ['contactus', 'media'],
    (name) => !name.includes('media'),
    (family) => `file:/tmp/${family}.tgz`,
  );

  assert.deepEqual(dependencies, {
    '@sneat/extension-contactus-contract': 'latest',
    '@sneat/extension-media-contract': 'file:/tmp/media.tgz',
  });
  assert.deepEqual(bootstrapped, ['media']);
});

test('preserves custom npm family names in the synthetic dependency plan', () => {
  const { dependencies } = createDependencyPlan(
    [{ dir: 'kids-club', npmName: 'kidsclub' }],
    () => true,
    () => assert.fail('published package must not use a local fallback'),
  );

  assert.deepEqual(dependencies, {
    '@sneat/extension-kidsclub-contract': 'latest',
  });
});

test('recognizes only explicit registry not-found failures as bootstrap candidates', () => {
  assert.equal(isRegistryNotFound({ stderr: 'npm error code E404' }), true);
  assert.equal(isRegistryNotFound({ stderr: '404 Not Found' }), true);
  assert.equal(isRegistryNotFound({ stderr: 'ECONNRESET registry unavailable' }), false);
  assert.equal(isRegistryNotFound({ stderr: 'E401 authentication required' }), false);
});
