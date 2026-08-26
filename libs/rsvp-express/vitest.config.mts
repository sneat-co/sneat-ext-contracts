import { defineConfig } from 'vitest/config';

// Pure-TypeScript contract lib (DTOs, contexts, service interfaces) — no DOM,
// no Angular TestBed, so `environment: 'node'`. Named "vitest.config.mts"
// (not "vite.config.mts") on purpose: nx.json registers `@nx/vite` and
// `@nx/vitest` as separate plugins with different inferred test-target names
// (testTargetName "test" vs "vite:test") so a vitest-only project like this
// one — whose "build" target stays the explicit `@nx/angular:package`
// (ng-packagr) executor in project.json — never collides with `@nx/vite`'s
// own inferred "build" target. `@nx/vite` only scans "vite.config.*"; this
// file is picked up solely by `@nx/vitest`, giving this project a "vite:test"
// target with zero risk of shadowing "build".
export default defineConfig({
  root: __dirname,
  cacheDir: '../../node_modules/.vite/libs/rsvp-express',
  test: {
    name: 'rsvp-express-contract',
    watch: false,
    globals: true,
    environment: 'node',
    include: ['src/**/*.test.ts', 'src/**/*.spec.ts'],
    reporters: ['default'],
    // No spec exists upstream either (verified against sneat-co/ext-rsvp-express) —
    // a genuinely test-free contract lib, not a dropped suite. Without this,
    // Vitest's own "no test files found" default exits 1 and fails CI.
    passWithNoTests: true,
  },
});
