import { defineConfig } from 'vitest/config';

// The contract package is pure TypeScript (an InjectionToken + interfaces) with
// no DOM or Ionic dependencies, so a plain node environment is enough.
export default defineConfig({
  test: {
    pool: 'forks',
    environment: 'node',
  },
});
