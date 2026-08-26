import { ISpaceContext } from '@sneat/space-models';
import { IContactRef } from './contact-ref';

/**
 * Cross-extension CONTRACT for reusing the contactus contacts-LIST component
 * WITHOUT importing its implementation.
 *
 * The component is authored in contactus `-shared` and registered as an Angular
 * **custom element** (Web Component) by the host app. A consumer extension uses
 * the tag `<sneat-contacts-list>` in its own template (with `CUSTOM_ELEMENTS_SCHEMA`)
 * and passes these contract-typed inputs — so it depends only on this contract,
 * never on `-shared`. This is the UI counterpart of the service-token pattern in
 * `contacts-selector.contract.ts`.
 */

/** The registered custom-element tag for the contactus contacts-list component. */
export const CONTACTUS_CONTACTS_LIST_TAG = 'sneat-contacts-list';

/** Inputs of `<sneat-contacts-list>` as a stable, contract-level shape. */
export interface IContactsListProps {
	readonly space: ISpaceContext;
	readonly contacts: readonly IContactRef[];
	readonly emptyText?: string;
}
