#!/usr/bin/env node
// Peer-range guard (REQ peer-range-no-exact-pins): fails CI if any
// libs/<family>/package.json declares a bare EXACT version (no ^/~/range
// operator) for a @sneat/* entry in peerDependencies or dependencies.
//
// Why this exists: `@nx/angular:package` (ng-packagr) copies a lib's source
// package.json fields — including peerDependencies — straight into the
// published dist/libs/<family>/package.json with no transformation. So the
// source file IS the built artifact for these fields; scanning it here is
// equivalent to scanning dist output, without needing a full build in CI.
//
// An exact @sneat/* peer pin (e.g. "@sneat/core": "0.27.6" instead of
// "^0.27.6") forces a consumer whose installed version has drifted even one
// patch to get a second copy installed by pnpm to satisfy the unmet peer —
// which is exactly the root cause of the fleet's Angular NG0201
// duplicate-DI-token errors. See tools/contract-generator/generator.ts's
// `caretOf()` for the generator-side half of this fix; this script is the
// guard that catches it if it ever creeps back in (by hand-edit, a future
// generator regression, or a copy-pasted libs/<family>/package.json).

import { readFileSync, readdirSync, existsSync } from 'node:fs';
import { fileURLToPath } from 'node:url';
import path from 'node:path';

const here = path.dirname(fileURLToPath(import.meta.url));
const repoRoot = path.resolve(here, '..', '..');
const libsDir = path.join(repoRoot, 'libs');

const isBareExact = (v) => typeof v === 'string' && /^\d/.test(v);

function main() {
  if (!existsSync(libsDir)) {
    console.log('peer-range-guard: no libs/ directory — trivial pass.');
    return;
  }

  const families = readdirSync(libsDir, { withFileTypes: true })
    .filter((e) => e.isDirectory())
    .map((e) => e.name);

  const violations = [];

  for (const family of families) {
    const pkgPath = path.join(libsDir, family, 'package.json');
    if (!existsSync(pkgPath)) continue;
    const pkg = JSON.parse(readFileSync(pkgPath, 'utf8'));

    for (const field of ['peerDependencies', 'dependencies']) {
      const deps = pkg[field];
      if (!deps) continue;
      for (const [name, range] of Object.entries(deps)) {
        if (name.startsWith('@sneat/') && isBareExact(range)) {
          violations.push({
            file: path.relative(repoRoot, pkgPath),
            field,
            name,
            range,
          });
        }
      }
    }
  }

  if (violations.length > 0) {
    console.error(
      `peer-range-guard: ${violations.length} exact @sneat/* peer/dependency pin(s) found — ` +
        'these must be caret ranges (e.g. "^0.27.6", not "0.27.6"), or a consumer whose ' +
        'installed version has drifted even one patch gets a second copy installed to satisfy ' +
        'the unmet peer (Angular NG0201 duplicate-DI-token errors).\n',
    );
    for (const v of violations) {
      console.error(`  ${v.file}: ${v.field}["${v.name}"] = "${v.range}" (missing ^)`);
    }
    process.exitCode = 1;
    return;
  }

  console.log(
    `peer-range-guard: checked ${families.length} lib(s) — no exact @sneat/* peer/dependency pins.`,
  );
}

main();
