import assert from 'node:assert/strict';
import test from 'node:test';
import { failedRequiredLegs } from './check.mjs';

const successful = {
  NX_RESULT: 'success',
  TIER_COHERENCE_RESULT: 'success',
  DISCOVER_GO_RESULT: 'success',
  GO_RESULT: 'success',
};

test('accepts all successful legs', () => {
  assert.deepEqual(failedRequiredLegs(successful), []);
});

test('accepts a skipped Go matrix when discovery found no modules', () => {
  assert.deepEqual(failedRequiredLegs({ ...successful, GO_RESULT: 'skipped' }), []);
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
