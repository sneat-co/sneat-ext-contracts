#!/usr/bin/env node
// Peer-range guard (REQ peer-range-no-exact-pins): fails CI if any
// libs/<family>/package.json declares a bare EXACT version (no ^/~/range
// operator) for a @sneat/* dependency or a workspace protocol in any published
// dependency field. Workspace references must resolve to registry ranges before packing.
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
import { fileURLToPath, pathToFileURL } from 'node:url';
import path from 'node:path';

const here = path.dirname(fileURLToPath(import.meta.url));
const repoRoot = path.resolve(here, '..', '..');
const libsDir = path.join(repoRoot, 'libs');

// A 0.0.x version is deliberately excluded: npm's caret on 0.0.x permits
// nothing but that exact patch (^0.0.4 means 0.0.4, full stop -- same match
// set as the bare "0.0.4"), so adding a caret there is a no-op that also
// trips the fleet's OTHER peer-range guard (sneat-co/cicd's
// check-peer-ranges, which rejects bare ^0.0.x/~0.0.x as its own trap). For
// a 0.0.x peer, staying exact is the fleet's established, correct form; only
// non-0.0.x bare exact pins are this script's target.
const isBareExact = (v) => typeof v === 'string' && /^\d/.test(v) && !/^0\.0\.\d/.test(v);

// ng-packagr and npm pack preserve dependency metadata verbatim. The same
// scanner runs on source manifests in CI and on dist manifests before packing.
export function inspectPackages(libsDir) {
  if (!existsSync(libsDir)) return { checked: 0, violations: [] };
  const families = readdirSync(libsDir, { withFileTypes: true })
    .filter((entry) => entry.isDirectory())
    .map((entry) => entry.name);
  const violations = [];
  let checked = 0;
  for (const family of families) {
    const pkgPath = path.join(libsDir, family, 'package.json');
    if (!existsSync(pkgPath)) continue;
    const pkg = JSON.parse(readFileSync(pkgPath, 'utf8'));
    if (pkg.private === true) continue;
    checked++;
    for (const field of ['peerDependencies', 'dependencies', 'optionalDependencies']) {
      for (const [name, range] of Object.entries(pkg[field] ?? {})) {
        let reason;
        if (typeof range === 'string' && range.startsWith('workspace:')) {
          reason = 'workspace protocol cannot be resolved outside this workspace';
        } else if (name.startsWith('@sneat/') && isBareExact(range)) {
          reason = 'missing range operator';
        }
        if (reason) violations.push({ file: pkgPath, field, name, range, reason });
      }
    }
  }
  return { checked, violations };
}

function main() {
  const directory = process.argv[2] ? path.resolve(process.argv[2]) : libsDir;
  const { checked, violations } = inspectPackages(directory);
  if (violations.length) {
    console.error(`peer-range-guard: ${violations.length} unpublishable dependency range(s):`);
    for (const v of violations) {
      console.error(`  ${path.relative(repoRoot, v.file)}: ${v.field}["${v.name}"] = "${v.range}" (${v.reason})`);
    }
    process.exitCode = 1;
    return;
  }
  console.log(`peer-range-guard: checked ${checked} publishable package(s) — dependency ranges are publishable.`);
}

if (process.argv[1] && import.meta.url === pathToFileURL(process.argv[1]).href) main();
