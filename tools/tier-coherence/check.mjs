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
import { fileURLToPath, pathToFileURL } from 'node:url';
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

// A contracts.json entry is either a plain string (npm package follows the
// default @sneat/extension-<entry>-contract pattern) or, when the published
// name doesn't match its dir/family tag, { dir, npmName }. `dir` is what we
// display and log; `npmName` (defaulting to `dir`) is the middle segment of
// the actual npm package name.
function entryDir(entry) {
  return typeof entry === 'string' ? entry : entry.dir;
}
function entryNpmName(entry) {
  return typeof entry === 'string' ? entry : (entry.npmName ?? entry.dir);
}

function packageName(entry) {
  return `@sneat/extension-${entryNpmName(entry)}-contract`;
}

export function isRegistryNotFound(error) {
  const output = `${error?.stdout ?? ''}\n${error?.stderr ?? ''}`;
  return /(?:^|\s)(?:E404|ERR_PNPM_FETCH_404)(?:\s|$)/m.test(output)
    || /\b404\s+Not Found\b/i.test(output);
}

function isPublished(name) {
  try {
    execFileSync('npm', ['view', `${name}@latest`, 'version', '--json'], {
      cwd: repoRoot,
      encoding: 'utf8',
      stdio: ['ignore', 'pipe', 'pipe'],
    });
    return true;
  } catch (error) {
    if (isRegistryNotFound(error)) {
      return false;
    }
    throw error;
  }
}

function packLocalContract(entry) {
  const dir = entryDir(entry);
  const project = JSON.parse(
    readFileSync(path.join(repoRoot, 'libs', dir, 'project.json'), 'utf8'),
  );
  const packDir = path.join(workDir, 'local-packages');
  mkdirSync(packDir, { recursive: true });
  execFileSync('pnpm', ['nx', 'build', project.name, '--skip-nx-cache'], {
    cwd: repoRoot,
    stdio: 'inherit',
  });
  const packed = JSON.parse(execFileSync(
    'npm',
    [
      'pack',
      path.join(repoRoot, 'dist', 'libs', dir),
      '--json',
      '--pack-destination',
      packDir,
    ],
    { cwd: repoRoot, encoding: 'utf8', stdio: ['ignore', 'pipe', 'inherit'] },
  ));
  if (packed.length !== 1 || packed[0].name !== packageName(entry)) {
    throw new Error(`tier-coherence: local package identity mismatch for ${dir}`);
  }
  return `file:${path.join(packDir, packed[0].filename)}`;
}

export function createDependencyPlan(families, published, localPackage) {
  const dependencies = {};
  const bootstrapped = [];
  for (const family of families) {
    const name = packageName(family);
    if (published(name)) {
      dependencies[name] = 'latest';
      continue;
    }
    dependencies[name] = localPackage(family);
    bootstrapped.push(entryDir(family));
  }
  return { dependencies, bootstrapped };
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
  // A family's first pull request is the one unavoidable exception: npm cannot
  // serve a package until that change has passed this check and released. Build
  // and pack only an exact registry-404 family locally so its real distributable
  // is still checked alongside every already-published latest package. Network,
  // authentication, and other registry failures remain hard failures.
  const { dependencies, bootstrapped } = createDependencyPlan(
    families,
    isPublished,
    packLocalContract,
  );
  if (bootstrapped.length > 0) {
    console.log(
      `tier-coherence: bootstrapping unpublished local contract(s): ${bootstrapped.join(', ')}`,
    );
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
    .map((family) => `import '@sneat/extension-${entryNpmName(family)}-contract';`)
    .join('\n');
  writeFileSync(path.join(workDir, 'import-everything.ts'), `${imports}\n`);

  console.log(
    `tier-coherence: installing latest of ${families.length} contract(s): ${families.map(entryDir).join(', ')}`,
  );
  // --ignore-workspace: this dir is NOT a member of the repo's pnpm workspace
  // (it depends on published npm versions, not workspace:* sources) and it
  // sits nested under one anyway (tools/tier-coherence/.tmp/consumer) — without
  // this flag pnpm walks up and treats it as a workspace member.
  //
  // --ignore-scripts: being outside the workspace, this install has none of
  // the root pnpm-workspace.yaml's `onlyBuiltDependencies` allowlist, so pnpm
  // hard-errors (ERR_PNPM_IGNORED_BUILDS) the moment any transitive dep from
  // an owned family's dependency tree ships a postinstall/build script that
  // isn't pre-approved here too (first hit: batch-2 families transitively
  // pulling in @firebase/util and protobufjs). This consumer only needs each
  // package's published .d.ts/.js on disk for `tsc` to type-check against —
  // it never executes any package's code — so skipping build scripts
  // entirely is correct here, not a risk tolerated for convenience.
  execFileSync(
    'pnpm',
    ['install', '--ignore-workspace', '--ignore-scripts', '--no-frozen-lockfile'],
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

if (
  process.argv[1]
  && import.meta.url === pathToFileURL(path.resolve(process.argv[1])).href
) {
  main();
}
