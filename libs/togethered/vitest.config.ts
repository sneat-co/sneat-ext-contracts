import { defineConfig } from 'vitest/config';

// Unlike ext-togethered's original vitest.config.ts (which carried heavy
// Ionic/jsdom aliasing for a service spec importing @sneat/api), the only
// spec migrated here (dto/list.spec.ts) exercises a pure function with no
// Angular/Ionic/DOM dependency, so this config stays minimal.
export default defineConfig({
  test: {
    pool: 'forks',
    isolate: true,
    environment: 'node',
  },
});
