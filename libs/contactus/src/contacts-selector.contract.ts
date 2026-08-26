import { InjectionToken } from '@angular/core';
import { ISpaceContext } from '@sneat/space-models';
import { IContactRef } from './contact-ref';

/**
 * Cross-extension CONTRACT for the contactus "select contacts" capability.
 *
 * A sibling extension (e.g. requoter picking drivers) depends ONLY on this
 * contract — it injects `CONTACTS_SELECTOR` and calls `selectMultipleContacts`,
 * never importing `@sneat/extension-contactus-ui`. The concrete
 * implementation (which opens the contactus selector modal, incl. the "add a new
 * person" tab) lives in the contactus UI package and is bound to this token by
 * the host app's `provideContactus()`.
 */
export interface IContactsSelectorProps {
	readonly space: ISpaceContext;
	/** Restrict to a contact type, e.g. `'person'`. */
	readonly contactType?: string;
	/** Restrict/seed by a contact role, e.g. the `driver` role. */
	readonly contactRole?: string;
	/** Contact ids to exclude from the pick list (e.g. already-added drivers). */
	readonly excludeContactIDs?: readonly string[];
	readonly title?: string;
	readonly okButtonLabel?: string;
}

export interface IContactsSelectorService {
	/** Open the picker; resolves to the chosen contacts, or `undefined` if cancelled. */
	selectMultipleContacts(
		props: IContactsSelectorProps,
	): Promise<readonly IContactRef[] | undefined>;
}

export const CONTACTS_SELECTOR = new InjectionToken<IContactsSelectorService>(
	'ContactsSelector',
);
