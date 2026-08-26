import { InjectionToken } from '@angular/core';
import { Observable } from 'rxjs';
import { ISpaceContext } from '@sneat/space-models';
import { IProfileView } from '../dto/profile/profile-view';

/**
 * Read contract for a requoter profile. The UI depends on this interface + the
 * token only; the implementation (in `-internal`) reads the stored profile doc
 * and the referenced canonical records directly from Firestore and composes the
 * {@link IProfileView}. Keeping this a token means the store can be repointed
 * (e.g. to a GET endpoint or a different backend) without touching the UI.
 */
export interface IRequoterProfileService {
  watchProfileView(
    space: ISpaceContext,
    profileID: string,
  ): Observable<IProfileView>;
}

export const REQUOTER_PROFILE_SERVICE =
  new InjectionToken<IRequoterProfileService>('RequoterProfileService');
