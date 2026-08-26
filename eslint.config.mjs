import nx from '@nx/eslint-plugin';

export default [
  ...nx.configs['flat/base'],
  ...nx.configs['flat/typescript'],
  ...nx.configs['flat/javascript'],
  {
    ignores: [
      '**/dist',
      '**/vite.config.*.timestamp*',
      '**/vitest.config.*.timestamp*',
    ],
  },
  {
    files: ['**/*.ts', '**/*.tsx', '**/*.js', '**/*.jsx'],
    rules: {
      // Module-boundary lint (REQ explicit-cross-contract-boundaries): default
      // is NO contract may import another contract. Every family lib carries
      // a unique `family:<name>` tag plus `layer:contract`; a depConstraints
      // entry below is the ONLY way to open an edge between two families, and
      // it must name both tags explicitly (`sourceTag` + `onlyDependOnLibsWithTags`).
      //
      // Seeding rule for whoever adds the next family (see docs/boundaries.md
      // for the authoritative table and the ownership test this enforces):
      //   1. Tag the new lib ['family:<name>', 'layer:contract'] in its project.json.
      //   2. Add ONE entry here: { sourceTag: 'family:<name>', onlyDependOnLibsWithTags: ['family:<name>'] }
      //      — this is the default (self-only, no cross-family import).
      //   3. Only widen onlyDependOnLibsWithTags with another family's tag when
      //      that specific edge is an approved platform-like dependency
      //      (contactus and assetus are the only ones named in the spec so
      //      far, Phase 2). Record the edge in docs/boundaries.md in the same
      //      change.
      '@nx/enforce-module-boundaries': [
        'error',
        {
          enforceBuildableLibDependency: true,
          allow: ['^.*/eslint(\\.base)?\\.config\\.[cm]?[jt]s$'],
          depConstraints: [],
        },
      ],
    },
  },
];
