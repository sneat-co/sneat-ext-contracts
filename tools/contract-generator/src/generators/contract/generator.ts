import {
  formatFiles,
  generateFiles,
  logger,
  readJson,
  updateJson,
  type Tree,
} from '@nx/devkit';
import * as path from 'node:path';
import { fileURLToPath } from 'node:url';
import type { ContractGeneratorSchema } from './schema';

// Nx loads this generator's compiled output as ESM in this workspace
// (tsconfig.base.json: "module": "esnext", no CJS override) — `__dirname`
// isn't defined there, so derive it from `import.meta.url` instead.
const __dirname = path.dirname(fileURLToPath(import.meta.url));

// Kebab-case: lowercase letters/digits, hyphen-separated, no leading/
// trailing/double hyphens. Plain single-word families (taxus, bookius) match
// too. Hyphens are fine in the Nx project name and npm package name — Go
// package identifiers can't contain them, so `--go` derives a sanitized
// `goPackageName` substitution (hyphens -> underscores) below.
const FAMILY_PATTERN = /^[a-z][a-z0-9]*(-[a-z0-9]+)*$/;

/**
 * Turns a devDependency version string ("22.1.3", "^0.22.0", "7.8.2") into a
 * caret floor pinned at its current major ("^22.0.0", "^7.0.0"). Reading the
 * live workspace root instead of hardcoding a peer floor is what keeps new
 * contracts honest as the workspace's own Angular/rxjs majors move — and
 * guarantees the generator can never emit a bare "^0.0.x" range (REQ
 * peer-range-strict, ci.yml).
 */
function majorFloor(version: string): string {
  const cleaned = version.replace(/^[\^~]/, '');
  const major = cleaned.split('.')[0] || '0';
  return `^${major}.0.0`;
}

function assertValidFamily(family: string | undefined): asserts family is string {
  if (!family || !FAMILY_PATTERN.test(family)) {
    throw new Error(
      `Invalid family name "${family ?? ''}" — must be lowercase, kebab-case ` +
        `(letters, digits, single hyphens, starting with a letter — e.g. "taxus", ` +
        `"kids-club"). This becomes the Nx project name "<family>-contract" and the ` +
        `npm package "@sneat/extension-<family>-contract".`,
    );
  }
}

/** go.work always opens with a "go <version>" directive — reuse it verbatim
 * for the new module rather than hardcoding a Go version the generator could
 * drift from. */
function readGoVersion(tree: Tree): string {
  const content = tree.read('go.work', 'utf-8');
  const match = content?.match(/^go\s+(\S+)/m);
  if (!match) {
    throw new Error(
      'go.work not found or missing a "go <version>" directive at its start.',
    );
  }
  return match[1];
}

/** Appends "./<family>" to go.work's `use ( ... )` block, keeping entries
 * sorted and de-duplicated. Idempotent: re-running for an already-registered
 * family is a no-op. */
function updateGoWork(tree: Tree, family: string): void {
  const goWorkPath = 'go.work';
  const content = tree.read(goWorkPath, 'utf-8');
  if (content == null) {
    throw new Error(`${goWorkPath} not found.`);
  }
  const useBlockRegex = /use \(([\s\S]*?)\)/;
  const match = content.match(useBlockRegex);
  if (!match) {
    throw new Error(
      `${goWorkPath}: could not find a "use ( ... )" block in the expected shape.`,
    );
  }
  const existing = match[1]
    .split('\n')
    .map((line) => line.trim())
    .filter(Boolean)
    .map((line) => line.replace(/^\.\//, ''));
  if (!existing.includes(family)) {
    existing.push(family);
  }
  existing.sort();
  const newBlock = `use (\n${existing.map((dir) => `\t./${dir}`).join('\n')}\n)`;
  tree.write(goWorkPath, content.replace(useBlockRegex, newBlock));
}

/** Adds the family to contracts.json's `families` array (sorted,
 * de-duplicated) so the tier-coherence check starts covering it — per
 * docs/boundaries.md's "Adding a new family lib" step 4. */
function updateContractsManifest(tree: Tree, family: string): void {
  updateJson(tree, 'contracts.json', (json) => {
    const families: string[] = Array.isArray(json.families) ? json.families : [];
    if (!families.includes(family)) {
      families.push(family);
    }
    json.families = families.sort();
    return json;
  });
}

/** Adds the default self-only `depConstraints` entry for the family to
 * eslint.config.mjs — per docs/boundaries.md's "Adding a new family lib" step
 * 3. Idempotent, and edits only the one line this generator owns rather than
 * reformatting/reserializing the whole file (a hand-authored .mjs file, not
 * JSON — round-tripping it through a JS/JSON parser would reformat every line
 * and blow up the diff). */
function updateEslintBoundaries(tree: Tree, family: string): void {
  const eslintPath = 'eslint.config.mjs';
  const content = tree.read(eslintPath, 'utf-8');
  if (content == null) {
    throw new Error(`${eslintPath} not found.`);
  }
  const entryLine = `{ sourceTag: 'family:${family}', onlyDependOnLibsWithTags: ['family:${family}'] },`;
  if (content.includes(entryLine)) {
    return;
  }
  const marker = 'depConstraints: [';
  const markerIndex = content.indexOf(marker);
  if (markerIndex === -1) {
    throw new Error(
      `${eslintPath}: could not find "${marker}" to register the new family's boundary.`,
    );
  }
  const afterMarker = content.slice(markerIndex + marker.length);
  const closeMatch = afterMarker.match(/\n(\s*)\],/);
  if (!closeMatch || closeMatch.index === undefined) {
    throw new Error(
      `${eslintPath}: could not find the end of the "depConstraints" array.`,
    );
  }
  const entryIndent = `${closeMatch[1]}  `;
  const insertion = `\n${entryIndent}${entryLine}`;
  const closeIndex = markerIndex + marker.length + closeMatch.index;
  tree.write(eslintPath, content.slice(0, closeIndex) + insertion + content.slice(closeIndex));
}

/** Adds/updates a row for the family in docs/boundaries.md's "Registered
 * families" table — a human-readable mirror of contracts.json. Creates the
 * section on first use; idempotent thereafter. This table does not exist yet
 * in the checked-in doc (only the cross-contract-edges table does), so the
 * first family scaffolded via this generator adds the section. */
function updateBoundariesDoc(
  tree: Tree,
  family: string,
  projectName: string,
  npmName: string,
): void {
  const docPath = 'docs/boundaries.md';
  const content = tree.read(docPath, 'utf-8');
  if (content == null) {
    throw new Error(`${docPath} not found.`);
  }
  const sectionHeader = '## Registered families';
  const row = `| \`${family}\` | \`${projectName}\` | \`${npmName}\` |`;
  if (content.includes(row)) {
    return;
  }

  if (content.includes(sectionHeader)) {
    const lines = content.split('\n');
    const headerIdx = lines.findIndex((line) => line.trim() === sectionHeader);
    let separatorIdx = -1;
    for (let i = headerIdx; i < lines.length; i++) {
      if (/^\|\s*-+\s*\|/.test(lines[i])) {
        separatorIdx = i;
        break;
      }
    }
    if (separatorIdx === -1) {
      throw new Error(`${docPath}: found "${sectionHeader}" but no table beneath it.`);
    }
    let end = separatorIdx + 1;
    const rows: string[] = [];
    while (end < lines.length && lines[end].startsWith('|')) {
      rows.push(lines[end]);
      end++;
    }
    rows.push(row);
    rows.sort();
    lines.splice(separatorIdx + 1, end - (separatorIdx + 1), ...rows);
    tree.write(docPath, lines.join('\n'));
  } else {
    const trimmed = content.replace(/\n+$/, '\n');
    const section = `\n${sectionHeader}\n\nEvery family scaffolded via \`tools/contract-generator\` (\`pnpm nx g ./tools/contract-generator:contract <family>\`), for a human-readable view alongside \`contracts.json\`:\n\n| Family | Nx project | npm package |\n| --- | --- | --- |\n${row}\n`;
    tree.write(docPath, `${trimmed}${section}`);
  }
}

export default async function contractGenerator(
  tree: Tree,
  options: ContractGeneratorSchema,
): Promise<void> {
  const family = options.family?.trim();
  assertValidFamily(family);

  const libRoot = `libs/${family}`;
  const projectName = `${family}-contract`;
  const npmName = `@sneat/extension-${family}-contract`;

  if (tree.exists(libRoot)) {
    throw new Error(`${libRoot} already exists — pick a different family name or remove it first.`);
  }
  if (options.go && tree.exists(family)) {
    throw new Error(`${family}/ already exists at the repo root — pick a different family name or remove it first.`);
  }

  const rootPkg = readJson(tree, 'package.json');
  const dev = rootPkg.devDependencies ?? {};

  const substitutions = {
    family,
    projectName,
    npmName,
    angularPeer: majorFloor(dev['@angular/core'] ?? '0.0.0'),
    rxjsPeer: majorFloor(dev['rxjs'] ?? '0.0.0'),
    sneatCorePeer: dev['@sneat/core'] ?? '^0.0.0',
    sneatDataPeer: dev['@sneat/data'] ?? '^0.0.0',
    sneatDtoPeer: dev['@sneat/dto'] ?? '^0.0.0',
    sneatSpaceModelsPeer: dev['@sneat/space-models'] ?? '^0.0.0',
    goVersion: options.go ? readGoVersion(tree) : '',
    // Go package identifiers can't contain hyphens (module paths can) — a
    // kebab-case family like "kids-club" gets Go package `kids_club`.
    goPackageName: family.replace(/-/g, '_'),
    tmpl: '',
  };

  generateFiles(tree, path.join(__dirname, 'files'), libRoot, substitutions);

  if (options.go) {
    generateFiles(tree, path.join(__dirname, 'files-go'), family, substitutions);
    updateGoWork(tree, family);
  }

  updateContractsManifest(tree, family);
  updateEslintBoundaries(tree, family);
  updateBoundariesDoc(tree, family, projectName, npmName);

  await formatFiles(tree);

  logger.info(`
Scaffolded ${libRoot} (${projectName} / ${npmName})${
    options.go
      ? ` and ${family}/go.mod (github.com/sneat-co/sneat-ext-contracts/${family})`
      : ''
  }.

Next steps:
1. Fill in the contract under ${libRoot}/src/ (DTOs, contexts, service
   interfaces), exported from src/index.ts — and this lib's README.md
   provenance/purpose stub.
2. Version-plan ritual (this family starts on the 0.x train):
     pnpm nx release plan ${projectName}
   0.x downshift table: patch -> patch, minor -> patch-level,
   major -> minor-level (no real 1.0.0 via a plan — see README.md "0.x
   semantics" for the full table and why).
3. Commit-scope rule: a conventional-commit scope must be exactly
   "${projectName}" (e.g. "feat(${projectName}): ...") — nx.json's
   release.version.conventionalCommits reads that scope literally; anything
   else caps the bump at patch.
4. Open one PR. CI enforces lint, module boundaries (docs/boundaries.md),
   peer-range strictness, and the tier-coherence check (this generator already
   added "${family}" to contracts.json).${
    options.go
      ? `\n5. Go module: fill in ${family}/doc.go, then \`cd ${family} && go build ./...\`. CI's discover-go job auto-picks up "${family}/go.mod".`
      : ''
  }
`);
}
