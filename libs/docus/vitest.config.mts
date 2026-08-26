import { defineConfig } from 'vitest/config';

export default defineConfig({
  test: {
    globals: true,
    environment: 'node',
    // doc-type-presentation.ts imports @sneat/extension-assetus-contract,
    // which re-exports @sneat/space-models@0.22.1's Angular Package Format
    // bundle (esm2022/sneat-space-models.js: `export * from './index'`) —
    // an extensionless relative specifier Node's native ESM loader refuses
    // to resolve outside a bundler. Vite/Vitest's own module graph handles
    // this correctly when the package is pre-bundled, so inlining it here
    // (same pattern formius's vitest.config.ts already uses for @ionic/*)
    // routes it through Vite instead of Node's raw loader.
    server: { deps: { inline: [/@sneat\/space-models/, /@sneat\/extension-assetus-contract/] } },
  },
});
