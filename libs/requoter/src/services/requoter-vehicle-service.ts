import { InjectionToken } from '@angular/core';
import { Observable } from 'rxjs';
import { ISpaceContext } from '@sneat/space-models';
import { IVehicleRef, VehicleOwnership } from '../dto/profile/vehicle';

/** What the Car step collects for a new car before it becomes an Assetus asset. */
export interface IVehicleInput {
  readonly name?: string;
  readonly make?: string;
  readonly model?: string;
  readonly reg?: string;
  readonly year?: number;
  /** `owned` (a car the user has) or `prospective` (a car under consideration). */
  readonly ownership: VehicleOwnership;
}

/**
 * Create/select the car through Assetus, returning an {@link IVehicleRef} (assetId
 * + ownership + display cache) — never a forked copy of the asset. The UI depends
 * on this contract token only; the implementation (in `-internal`) talks to
 * Assetus's `AssetService`, so ReQuoter never writes the asset record directly.
 */
export interface IRequoterVehicleService {
  /** Existing vehicles in the space, as refs, so the user can pick one. */
  watchVehicles(space: ISpaceContext): Observable<IVehicleRef[]>;
  /** Create a new car in Assetus (prospective cars are flagged non-owned). */
  createVehicle(space: ISpaceContext, input: IVehicleInput): Observable<IVehicleRef>;
  /** Reference an existing Assetus asset by id — no duplicate is created. */
  selectVehicle(space: ISpaceContext, assetId: string): Observable<IVehicleRef>;
}

export const REQUOTER_VEHICLE_SERVICE =
  new InjectionToken<IRequoterVehicleService>('RequoterVehicleService');
