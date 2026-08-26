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
  },
});
