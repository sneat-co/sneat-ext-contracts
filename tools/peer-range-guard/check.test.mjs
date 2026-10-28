import assert from 'node:assert/strict';
import { mkdtempSync, mkdirSync, writeFileSync, rmSync } from 'node:fs';
import { tmpdir } from 'node:os';
import path from 'node:path';
import { spawnSync } from 'node:child_process';
import { fileURLToPath } from 'node:url';
import test from 'node:test';
import { inspectPackages } from './check.mjs';

function fixture(t, packages) {
  const root = mkdtempSync(path.join(tmpdir(), 'publishable-peers-'));
  t.after(() => rmSync(root, { recursive: true, force: true }));
  for (const [family, manifest] of Object.entries(packages)) {
    mkdirSync(path.join(root, family));
    if (manifest) writeFileSync(path.join(root, family, 'package.json'), JSON.stringify(manifest));
  }
  return root;
}

test('rejects every workspace dependency protocol in published dependency fields', (t) => {
  for (const field of ['dependencies', 'peerDependencies', 'optionalDependencies']) {
    for (const range of ['workspace:^', 'workspace:~', 'workspace:*', 'workspace:^0.1.1', 'workspace:../media']) {
      const root = fixture(t, { contactus: { [field]: { '@sneat/extension-media-contract': range } } });
      const result = inspectPackages(root);
      assert.equal(result.checked, 1);
      assert.deepEqual(result.violations, [{
        file: path.join(root, 'contactus/package.json'), field,
        name: '@sneat/extension-media-contract', range,
        reason: 'workspace protocol cannot be resolved outside this workspace',
      }]);
    }
  }
});

test('retains numeric peer guard while allowing semantic ranges, 0.0.x pins and workspace dev dependencies', (t) => {
  const root = fixture(t, {
    valid: { peerDependencies: { '@sneat/media': '^0.1.1', '@angular/core': '^21.0.0 || ^22.0.0', '@sneat/tiny': '0.0.4' }, devDependencies: { local: 'workspace:*' } },
    exact: { dependencies: { '@sneat/core': '0.27.6', 'third-party': '1.0.0' } },
    private: { private: true, dependencies: { local: 'workspace:*' } },
    missing: null,
  });
  const { checked, violations } = inspectPackages(root);
  assert.equal(checked, 2);
  assert.equal(violations.length, 1);
  assert.equal(violations[0].name, '@sneat/core');
  assert.equal(violations[0].reason, 'missing range operator');
  assert.deepEqual(inspectPackages(path.join(root, 'absent')), { checked: 0, violations: [] });
});

test('CLI refuses built metadata before npm pack and accepts fixed dist metadata', (t) => {
  const root = fixture(t, { contactus: { peerDependencies: { '@sneat/media': 'workspace:^' } } });
  const command = fileURLToPath(new URL('./check.mjs', import.meta.url));
  const bad = spawnSync(process.execPath, [command, root], { encoding: 'utf8' });
  assert.equal(bad.status, 1);
  assert.match(bad.stderr, /workspace protocol cannot be resolved/);
  writeFileSync(path.join(root, 'contactus/package.json'), JSON.stringify({ peerDependencies: { '@sneat/media': '^0.1.1' } }));
  const good = spawnSync(process.execPath, [command, root], { encoding: 'utf8' });
  assert.equal(good.status, 0);
  assert.match(good.stdout, /checked 1 publishable package/);
});
