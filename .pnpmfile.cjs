'use strict';

/**
 * zone.js is not a dependency of this workspace (see the `zone.js` override
 * in pnpm-workspace.yaml) and this repo runs zoneless end to end (it has no
 * components and no zone-based tests at all). Upstream packages still declare
 * it as an *optional* peer dependency — @angular/core today, and
 * @analogjs/vitest-angular once a family adds Angular-TestBed specs — and
 * pnpm's `auto-install-peers` (default on; this repo never flips it) installs
 * unmet peers even when `peerDependenciesMeta` marks them optional. That is a
 * known pnpm bug: https://github.com/pnpm/pnpm/issues/11155. The documented
 * workaround is a `readPackage` hook that strips the peer declaration before
 * pnpm's peer-resolution phase ever sees it, so there is nothing left to
 * auto-install. Reference: debtus frontend/.pnpmfile.cjs (2026-08-25).
 */
function stripOptionalZoneJsPeer(pkg) {
  if (pkg.peerDependencies && 'zone.js' in pkg.peerDependencies) {
    delete pkg.peerDependencies['zone.js'];
  }
  if (pkg.peerDependenciesMeta && 'zone.js' in pkg.peerDependenciesMeta) {
    delete pkg.peerDependenciesMeta['zone.js'];
  }
  return pkg;
}

module.exports = {
  hooks: {
    readPackage: stripOptionalZoneJsPeer,
  },
};
