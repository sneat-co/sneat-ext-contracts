import { InjectionToken } from '@angular/core';
import { ISpaceContext } from '@sneat/space-models';
import { Observable } from 'rxjs';
import { IListGroup } from '../dto';

/**
 * Reads the lists that belong to one Space for overview surfaces.
 *
 * The base Listus service intentionally remains source-compatible: ordinary
 * Firebase-backed implementations do not need to implement this capability.
 * A storage adapter can provide it when a Space's list overview is stored
 * outside the Space DBO (for example, an authenticated OVDB demo session).
 */
export interface IListusListGroupsReader {
  watchListGroups(space: ISpaceContext): Observable<readonly IListGroup[]>;
}

/** Optional capability for Listus overview/menu consumers. */
export const LISTUS_LIST_GROUPS_READER =
  new InjectionToken<IListusListGroupsReader>('ListusListGroupsReader');
