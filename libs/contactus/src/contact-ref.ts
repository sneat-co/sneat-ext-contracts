/**
 * Minimal, foundational contact identity used by the cross-extension contactus
 * contract surfaces below. It carries only what a consumer needs to identify and
 * label a picked contact — it does NOT drag the full contactus contact model.
 *
 * (When the full `@sneat/extension-contactus-contract` API is migrated into this
 * repo, this can be reconciled with the richer `IContactBrief` / `IContactWithBrief`
 * types; for now it is deliberately small and self-contained so this lib depends
 * only on foundational packages — never on another extension.)
 */
export interface IContactPersonNames {
	readonly firstName?: string;
	readonly lastName?: string;
}

/** A contact as seen across an extension boundary: id + display + optional names. */
export interface IContactRef {
	readonly id: string;
	/** A human label to show (full name / title). */
	readonly title: string;
	readonly names?: IContactPersonNames;
	/** Contact roles (e.g. the `driver` role), if the provider supplies them. */
	readonly roles?: readonly string[];
}
