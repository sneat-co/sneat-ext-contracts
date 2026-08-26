import { InjectionToken } from '@angular/core';
import { Observable } from 'rxjs';
import {
  ModuleSpaceItemService,
  SpaceModuleService,
} from '@sneat/space-services';
import { ITrackerBrief, ITrackerDbo } from '../dbo/i-tracker-dbo';
import { ITrackusSpaceDbo } from '../dbo/i-trackus-space-dbo';
import {
  IAddTrackerPointRequest,
  IAddTrackerPointResponse,
  ICreateTrackerRequest,
  ICreateTrackerResponse,
  IDeleteTrackerPointsRequest,
  ITrackerRequest,
} from './trackus-api.types';

// The trackus pages/components depend only on these contract tokens. The
// concrete services live in -internal and are bound to the tokens by
// provideTrackusInternal() at app bootstrap (the app is the composition root).

// TrackusSpaceService adds no methods over SpaceModuleService, so the contract
// surface IS the base class.
export type ITrackusSpaceService = SpaceModuleService<ITrackusSpaceDbo>;
export const TRACKUS_SPACE_SERVICE = new InjectionToken<ITrackusSpaceService>(
  'TrackusSpaceService',
);

// TrackersService adds no methods over ModuleSpaceItemService.
export type ITrackersService = ModuleSpaceItemService<
  ITrackerBrief,
  ITrackerDbo
>;
export const TRACKERS_SERVICE = new InjectionToken<ITrackersService>(
  'TrackersService',
);

export interface ITrackusApiService {
  addTracker(
    request: ICreateTrackerRequest,
  ): Observable<ICreateTrackerResponse>;
  createTracker(
    request: ICreateTrackerRequest,
  ): Observable<ICreateTrackerResponse>;
  archiveTracker(request: ITrackerRequest): Observable<void>;
  addTrackerPoint(
    request: IAddTrackerPointRequest,
  ): Observable<IAddTrackerPointResponse>;
  deleteTrackerPoints(request: IDeleteTrackerPointsRequest): Observable<void>;
}
export const TRACKUS_API_SERVICE = new InjectionToken<ITrackusApiService>(
  'TrackusApiService',
);
