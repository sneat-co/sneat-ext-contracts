import { defineConfig } from 'vitest/config';
import { fileURLToPath } from 'node:url';
import path from 'node:path';

const root = path.dirname(fileURLToPath(import.meta.url));

// Standalone config for the generator's own test suite, runnable from any
// cwd via `pnpm test:contract-generator` (root package.json) or Nx's own
// inferred "vite:test" target on the "contract-generator" project (it has a
// package.json — see that file's comment — so `@nx/vitest` infers targets
// from this config the same as any other project).
export default defineConfig({
  root,
  test: {
    include: ['src/**/*.spec.ts'],
    environment: 'node',
    globals: true,
    watch: false,
    reporters: ['default'],
    // The real-CLI generator test shells out to `pnpm nx g`, a full
    // workspace lint, and a shell discovery command — sequential by design
    // (they share/mutate the same scratch family on disk).
    fileParallelism: false,
  },
});
