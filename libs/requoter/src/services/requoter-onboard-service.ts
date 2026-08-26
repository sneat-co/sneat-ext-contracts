import { InjectionToken } from '@angular/core';
import { Observable } from 'rxjs';
import { IOnboardRequest, IOnboardResponse } from '../dto/onboard';

/**
 * Contract for the ReQuoter onboarding call. UI depends on this interface + the
 * REQUOTER_ONBOARD_SERVICE token; the implementation (which POSTs to
 * /v0/requoter/onboard) lives in the internal lib and is bound by
 * provideRequoterInternal().
 */
export interface IRequoterOnboardService {
	onboard(request: IOnboardRequest): Observable<IOnboardResponse>;
}

export const REQUOTER_ONBOARD_SERVICE =
	new InjectionToken<IRequoterOnboardService>('RequoterOnboardService');
