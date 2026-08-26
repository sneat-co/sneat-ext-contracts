import { defineConfig } from 'vitest/config';

// contactus-contract ships one real unit spec (contact-requests DTO
// validators). This config exists so the workspace-wide @nx/vite plugin
// (registered in nx.json, testTargetName: "test") infers a "test" target for
// this project — the same gap gameboard-contract hit first (see its
// vite.config.ts): only @nx/vitest actually infers a target, named
// "vite:test", so this project's own project.json also declares an explicit
// "test" target (nx:run-commands, vitest run).
export default defineConfig({
  test: {
    environment: 'node',
    include: ['src/**/*.{test,spec}.ts'],
  },
});
