import { describe, expect, it } from 'vitest';
import { execFileSync } from 'node:child_process';
import { existsSync, readFileSync, rmSync, writeFileSync } from 'node:fs';
import path from 'node:path';

// tools/contract-generator/src/generators/contract -> repo root is 5 levels up.
const repoRoot = path.resolve(__dirname, '..', '..', '..', '..', '..');

// Kebab-case on purpose: exercises the same hyphen-sanitization path a real
// future family (e.g. "kids-club") will need for its Go package identifier.
const SCRATCH_FAMILY = 'zz-scratch-gen';
const PROJECT_NAME = `${SCRATCH_FAMILY}-contract`;
const NPM_NAME = `@sneat/extension-${SCRATCH_FAMILY}-contract`;
const GO_PACKAGE_NAME = 'zz_scratch_gen';

const LIB_ROOT = path.join(repoRoot, 'libs', SCRATCH_FAMILY);
const GO_ROOT = path.join(repoRoot, SCRATCH_FAMILY);
const CONTRACTS_JSON = path.join(repoRoot, 'contracts.json');
const ESLINT_CONFIG = path.join(repoRoot, 'eslint.config.mjs');
const BOUNDARIES_DOC = path.join(repoRoot, 'docs', 'boundaries.md');
const GO_WORK = path.join(repoRoot, 'go.work');

function read(filePath: string): string {
  return readFileSync(filePath, 'utf-8');
}

function runGenerator(extraArgs: string[] = []): string {
  return execFileSync(
    'pnpm',
    ['nx', 'g', './tools/contract-generator:contract', SCRATCH_FAMILY, '--no-interactive', ...extraArgs],
    { cwd: repoRoot, encoding: 'utf-8' },
  );
}

function cleanupScratchFamily(snapshots: {
  contractsJson: string;
  eslintConfig: string;
  boundariesDoc: string;
  goWork: string;
}): void {
  rmSync(LIB_ROOT, { recursive: true, force: true });
  rmSync(GO_ROOT, { recursive: true, force: true });
  writeFileSync(CONTRACTS_JSON, snapshots.contractsJson);
  writeFileSync(ESLINT_CONFIG, snapshots.eslintConfig);
  writeFileSync(BOUNDARIES_DOC, snapshots.boundariesDoc);
  writeFileSync(GO_WORK, snapshots.goWork);
}

describe('contract generator (real CLI, real workspace)', () => {
  // Runs a real `pnpm nx g` invocation, a full workspace lint, and a shell
  // command against the actual repo tree — slower than a unit test, but this
  // is exactly the "one command" path a human/agent will run, and the only
  // way to prove the lint and Go-discovery acceptance criteria for real.
  it(
    'scaffolds libs/<family> + <family>/go.mod, registers every shared file, keeps the workspace lint-clean, and stays discoverable by CI\'s Go matrix job',
    () => {
      expect(existsSync(LIB_ROOT)).toBe(false);
      expect(existsSync(GO_ROOT)).toBe(false);

      const snapshots = {
        contractsJson: read(CONTRACTS_JSON),
        eslintConfig: read(ESLINT_CONFIG),
        boundariesDoc: read(BOUNDARIES_DOC),
        goWork: read(GO_WORK),
      };

      try {
        runGenerator(['--go']);

        // --- libs/<family> scaffold ---
        for (const relFile of [
          'project.json',
          'package.json',
          'ng-package.json',
          'tsconfig.json',
          'tsconfig.lib.json',
          'tsconfig.lib.prod.json',
          'tsconfig.spec.json',
          'vitest.config.mts',
          'README.md',
          'src/index.ts',
        ]) {
          expect(existsSync(path.join(LIB_ROOT, relFile)), `missing libs/${SCRATCH_FAMILY}/${relFile}`).toBe(true);
        }

        const projectJson = JSON.parse(read(path.join(LIB_ROOT, 'project.json')));
        expect(projectJson.name).toBe(PROJECT_NAME);
        expect(projectJson.tags).toEqual(
          expect.arrayContaining([`family:${SCRATCH_FAMILY}`, 'layer:contract']),
        );
        expect(projectJson.targets.build.executor).toBe('@nx/angular:package');

        const packageJson = JSON.parse(read(path.join(LIB_ROOT, 'package.json')));
        expect(packageJson.name).toBe(NPM_NAME);
        for (const peer of Object.values(packageJson.peerDependencies) as string[]) {
          expect(peer.startsWith('^0.0.') || peer === '^0.0.0').toBe(false);
        }
        const rootPkg = JSON.parse(read(path.join(repoRoot, 'package.json')));
        const angularMajor = String(rootPkg.devDependencies['@angular/core']).split('.')[0];
        const rxjsMajor = String(rootPkg.devDependencies['rxjs']).split('.')[0];
        expect(packageJson.peerDependencies['@angular/core']).toBe(`^${angularMajor}.0.0`);
        expect(packageJson.peerDependencies['rxjs']).toBe(`^${rxjsMajor}.0.0`);
        expect(packageJson.peerDependencies['@sneat/core']).toBe(rootPkg.devDependencies['@sneat/core']);

        // --- <family>/go.mod scaffold ---
        expect(existsSync(path.join(GO_ROOT, 'go.mod'))).toBe(true);
        expect(existsSync(path.join(GO_ROOT, 'doc.go'))).toBe(true);
        const goMod = read(path.join(GO_ROOT, 'go.mod'));
        expect(goMod).toContain(`module github.com/sneat-co/sneat-ext-contracts/${SCRATCH_FAMILY}`);
        const goWorkGoVersion = snapshots.goWork.match(/^go\s+(\S+)/m)?.[1];
        expect(goMod).toContain(`go ${goWorkGoVersion}`);
        const docGo = read(path.join(GO_ROOT, 'doc.go'));
        expect(docGo).toContain(`package ${GO_PACKAGE_NAME}`);

        // --- go.work registration ---
        const updatedGoWork = read(GO_WORK);
        expect(updatedGoWork).toMatch(new RegExp(`\\./${SCRATCH_FAMILY}\\b`));

        // --- contracts.json registration ---
        const contracts = JSON.parse(read(CONTRACTS_JSON));
        expect(contracts.families).toContain(SCRATCH_FAMILY);

        // --- eslint.config.mjs boundary registration ---
        const eslintConfig = read(ESLINT_CONFIG);
        expect(eslintConfig).toContain(
          `{ sourceTag: 'family:${SCRATCH_FAMILY}', onlyDependOnLibsWithTags: ['family:${SCRATCH_FAMILY}'] },`,
        );

        // --- docs/boundaries.md registration ---
        const boundariesDoc = read(BOUNDARIES_DOC);
        expect(boundariesDoc).toContain('## Registered families');
        expect(boundariesDoc).toContain(`\`${SCRATCH_FAMILY}\``);
        expect(boundariesDoc).toContain(`\`${PROJECT_NAME}\``);
        expect(boundariesDoc).toContain(`\`${NPM_NAME}\``);

        // --- idempotency: re-running for the same family touches none of the
        // three shared files a second time (no duplicate rows/entries) ---
        const beforeRerun = {
          contractsJson: read(CONTRACTS_JSON),
          eslintConfig: read(ESLINT_CONFIG),
          boundariesDoc: read(BOUNDARIES_DOC),
        };
        rmSync(LIB_ROOT, { recursive: true, force: true });
        rmSync(GO_ROOT, { recursive: true, force: true });
        runGenerator(['--go']);
        expect(read(CONTRACTS_JSON)).toBe(beforeRerun.contractsJson);
        expect(read(ESLINT_CONFIG)).toBe(beforeRerun.eslintConfig);
        expect(read(BOUNDARIES_DOC)).toBe(beforeRerun.boundariesDoc);

        // --- acceptance: the whole workspace, scratch family included,
        // still passes `nx run-many -t lint` ---
        execFileSync('pnpm', ['nx', 'run-many', '-t', 'lint', '--all'], {
          cwd: repoRoot,
          stdio: 'inherit',
        });

        // --- acceptance: ci.yml's discover-go matrix command tolerates the
        // new module (replicated verbatim from .github/workflows/ci.yml) ---
        const discovered = execFileSync(
          'bash',
          [
            '-c',
            "find . -mindepth 2 -maxdepth 2 -name go.mod -not -path './node_modules/*' | sed -e 's|^\\./||' -e 's|/go\\.mod$||' | jq -R -s -c 'split(\"\\n\") | map(select(length > 0)) | sort'",
          ],
          { cwd: repoRoot, encoding: 'utf-8' },
        );
        const discoveredDirs: string[] = JSON.parse(discovered);
        expect(discoveredDirs).toContain(SCRATCH_FAMILY);
      } finally {
        cleanupScratchFamily(snapshots);
        expect(existsSync(LIB_ROOT)).toBe(false);
        expect(existsSync(GO_ROOT)).toBe(false);
        expect(read(CONTRACTS_JSON)).toBe(snapshots.contractsJson);
        expect(read(ESLINT_CONFIG)).toBe(snapshots.eslintConfig);
        expect(read(BOUNDARIES_DOC)).toBe(snapshots.boundariesDoc);
        expect(read(GO_WORK)).toBe(snapshots.goWork);
      }
    },
    300_000,
  );

  it('rejects an invalid family name before touching the filesystem', () => {
    const invalidLibRoot = path.join(repoRoot, 'libs', 'Not_Valid');
    expect(existsSync(invalidLibRoot)).toBe(false);
    expect(() =>
      execFileSync(
        'pnpm',
        ['nx', 'g', './tools/contract-generator:contract', 'Not_Valid', '--no-interactive'],
        { cwd: repoRoot, encoding: 'utf-8', stdio: 'pipe' },
      ),
    ).toThrow();
    expect(existsSync(invalidLibRoot)).toBe(false);
  });
});
