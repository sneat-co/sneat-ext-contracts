import { defineConfig } from 'vitest/config';
export default defineConfig({ root: import.meta.dirname, test: { name: 'commitius-contract', environment: 'node', include: ['src/**/*.spec.ts'] } });
