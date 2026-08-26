import { defineConfig } from 'vitest/config';

// Pure-TypeScript contract lib (DTOs, contexts, service interfaces) — no DOM,
// no Angular TestBed, so `environment: 'node'`.
export default defineConfig({
  root: __dirname,
  cacheDir: '../../node_modules/.vite/libs/schoolus',
  test: {
    name: 'schoolus-contract',
    watch: false,
    globals: true,
    environment: 'node',
    include: ['src/**/*.spec.ts'],
    reporters: ['default'],
    // No spec exists upstream either (verified against sneat-co/ext-schoolus@a703308) —
    // a genuinely test-free contract lib on the TS side (the Go side has real
    // dto4schoolus/dto_test.go coverage). Without this, Vitest's own "no test
    // files found" default exits 1 and fails CI.
    passWithNoTests: true,
  },
});
