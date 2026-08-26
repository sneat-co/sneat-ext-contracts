#!/usr/bin/env node
// Tier-coherence check (REQ tier-coherence-check, spec/features/
// ext-contracts-monorepo/README.md on sneat-co/sneat-libs): proves the latest
// published version of every @sneat/extension-<family>-contract package this
// repo owns installs and type-checks together in one synthetic consumer. This
// is the guardrail that keeps independent per-project versioning skew-safe —
// it is the check that would have caught the five desynced pipelines the
// 2026-08-26 fleet audit found.
//
// With zero owned families (contracts.json's `families` array is empty — true
// at Phase 0, this bootstrap) this passes trivially: no network call, no
// install, no tsc. A family enters coverage by being added to contracts.json,
// normally in the same change that adds its libs/<family>/ lib.

import { execFileSync } from 'node:child_process';
import { mkdirSync, writeFileSync, rmSync, readFileSync } from 'node:fs';
import { fileURLToPath } from 'node:url';
import path from 'node:path';

const here = path.dirname(fileURLToPath(import.meta.url));
const repoRoot = path.resolve(here, '..', '..');
const manifestPath = path.join(repoRoot, 'contracts.json');
const workDir = path.join(here, '.tmp', 'consumer');

function readTypescriptVersion() {
  const rootPkg = JSON.parse(
    readFileSync(path.join(repoRoot, 'package.json'), 'utf8'),
  );
  return rootPkg.devDependencies?.typescript ?? 'latest';
}

function main() {
  const manifest = JSON.parse(readFileSync(manifestPath, 'utf8'));
  const families = manifest.families ?? [];

  if (families.length === 0) {
    console.log(
      'tier-coherence: contracts.json lists no owned families yet — trivial pass.',
    );
    return;
  }

  rmSync(workDir, { recursive: true, force: true });
  mkdirSync(workDir, { recursive: true });

  // "latest" (not a resolved pin): the whole point is to prove whatever npm
  // currently serves as latest for every owned family installs together.
  const dependencies = {};
  for (const family of families) {
    dependencies[`@sneat/extension-${family}-contract`] = 'latest';
  }

  writeFileSync(
    path.join(workDir, 'package.json'),
    JSON.stringify(
      {
        name: 'tier-coherence-synthetic-consumer',
        version: '0.0.0',
        private: true,
        dependencies,
        devDependencies: {
          typescript: readTypescriptVersion(),
        },
      },
      null,
      2,
    ),
  );

  writeFileSync(
    path.join(workDir, 'tsconfig.json'),
    JSON.stringify(
      {
        compilerOptions: {
          target: 'es2022',
          module: 'esnext',
          moduleResolution: 'bundler',
          strict: true,
          skipLibCheck: true,
          noEmit: true,
          types: [],
        },
        include: ['import-everything.ts'],
      },
      null,
      2,
    ),
  );

  const imports = families
    .map((family) => `import '@sneat/extension-${family}-contract';`)
    .join('\n');
  writeFileSync(path.join(workDir, 'import-everything.ts'), `${imports}\n`);

  console.log(
    `tier-coherence: installing latest of ${families.length} contract(s): ${families.join(', ')}`,
  );
  // --ignore-workspace: this dir is NOT a member of the repo's pnpm workspace
  // (it depends on published npm versions, not workspace:* sources) and it
  // sits nested under one anyway (tools/tier-coherence/.tmp/consumer) — without
  // this flag pnpm walks up and treats it as a workspace member.
  execFileSync(
    'pnpm',
    ['install', '--ignore-workspace', '--no-frozen-lockfile'],
    { cwd: workDir, stdio: 'inherit' },
  );

  console.log('tier-coherence: type-checking the synthetic consumer...');
  execFileSync('pnpm', ['exec', 'tsc', '-p', 'tsconfig.json'], {
    cwd: workDir,
    stdio: 'inherit',
  });

  console.log(
    `tier-coherence: OK — ${families.length} contract(s) install and type-check together.`,
  );
}

main();
