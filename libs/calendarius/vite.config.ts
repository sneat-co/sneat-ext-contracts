import { defineConfig } from 'vitest/config';

// calendarius-contract ships five unit specs, all pure DTO/logic tests
// (dto/event-happening, dto/happening, dto/todo_move_funcs, view-models,
// contexts/happening-context) — none touch Angular DI/TestBed or the DOM, so
// `environment: 'node'` needs no Angular plugin and no test-setup.ts. This
// config exists so the workspace-wide @nx/vite plugin (registered in
// nx.json, testTargetName: "test") infers a "test" target for this project —
// the same gap gameboard-contract's migration found first (see its
// vite.config.ts): only @nx/vitest actually infers a target, named
// "vite:test", so this project's own project.json also declares an explicit
// "test" target (nx:run-commands, vitest run).
//
// The source workspace (sneat-libs) runs these specs with `globals: true`
// (see its shared vite.config.base.ts); this workspace enables no globals
// (established by libs/taxus onward), so four of the five spec files
// (all but dto/event-happening.spec.ts, which already imported explicitly)
// gained explicit `describe`/`it`/`expect`(/`beforeEach`) imports from
// 'vitest' during this migration — the same fix gameboard's migration made
// for its own ported spec.
export default defineConfig({
  test: {
    environment: 'node',
    include: ['src/**/*.{test,spec}.ts'],
  },
});
