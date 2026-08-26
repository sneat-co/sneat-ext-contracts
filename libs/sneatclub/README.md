# @sneat/extension-sneatclub-contract

Public Sneat Club DTOs, contexts, service interfaces, and
dependency-injection tokens. Runtime implementation and UI stay in the
private Sneat Club repository (`sneat-co/sneat-club`,
`libs/extensions/sneatclub/runtime`).

## Provenance

Authored from the vendored placeholder's declared `.d.ts` surface per founder
decision 2026-08-26. `sneat-co/sneat-club` commit
`898b30af9e5b6136296aec492265ca7954a935ab` ("Scaffold the Sneat Club app;
registration moves to /app/register", 2026-08-24) vendored the published
`@sneat/extension-template-contract@0.1.0` package build output (fesm2022 +
types only, no src) at `libs/extensions/sneatclub/contract` with
`Template`/`template` tokens renamed to `Sneat Club`/`sneatclub`, consumed via
`workspace:*`, "until an ext-sneatclub contract repo publishes the real one."

This package is that real contract: real TypeScript source (`src/`) authored
here to reproduce the vendored `.d.ts` export surface exactly (interfaces,
tokens, consts) — not copied dist. Reference read:
`sneat-co/sneat-club@09c4d0cfc4651df02bc2060114ee7ad047d07c45`
(`libs/extensions/sneatclub/contract`, `origin/main`, 2026-08-26 read).

This is a first-publish — no npm baseline exists for
`@sneat/extension-sneatclub-contract`.

## Cross-family dependencies

None.
