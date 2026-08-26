import { defineConfig } from 'vitest/config';

// A service spec imports @sneat/api, which transitively pulls @ionic/angular.
// @ionic/angular does a bare directory import of '@ionic/core/components' that
// Node's ESM resolver rejects; redirect it to the package's index.js. Inlining
// the @ionic/@sneat packages keeps the rest of the chain bundled by Vitest.
export default defineConfig({
  resolve: {
    alias: [
      {
        find: /^@ionic\/core\/components$/,
        replacement: '@ionic/core/components/index.js',
      },
    ],
  },
  test: {
    // asset.service.spec mocks @angular/fire/firestore (getHistory reads
    // Firestore directly). Isolated forks give every test file its own module
    // registry so that mock can never leak into sibling specs that use the real
    // Firestore — a shared worker registry made the suite load-order flaky.
    pool: 'forks',
    isolate: true,
    // jsdom + the analog TestBed setup (injected by @nx/angular:unit-test) so
    // the ported standalone Angular components/pages can be created and rendered
    // in specs. Pure-logic specs (DTOs, services) run fine under jsdom too.
    environment: 'jsdom',
    environmentOptions: { jsdom: { url: 'http://localhost/' } },
    // Ionic's icon lazy-loader tries to fetch SVGs in jsdom and throws benign
    // async "Invalid URL" / fetch errors when components with icons render.
    // All assertions still pass; don't let these unhandled async errors fail.
    dangerouslyIgnoreUnhandledErrors: true,
    server: { deps: { inline: [/@ionic/, /ionicons/, /@sneat/] } },
    deps: {
      optimizer: {
        web: { include: ['@ionic/angular', '@ionic/core', 'ionicons'] },
      },
    },
  },
});
