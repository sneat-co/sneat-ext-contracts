import { defineConfig } from 'vitest/config';

// gameboard-contract is the only family lib in this repo (so far) that ships
// real runtime logic (the event-timeline fold reducer) with real tests — the
// Go<->TS parity oracle and the legacy list-DTO unit spec. This config exists
// so the workspace-wide @nx/vite plugin (registered in nx.json,
// testTargetName: "test") infers a "test" target for this project; taxus has
// none of this because it currently ships no tests.
export default defineConfig({
  test: {
    environment: 'node',
    include: ['src/**/*.{test,spec}.ts'],
  },
});
