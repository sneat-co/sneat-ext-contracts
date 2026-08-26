import { InjectionToken } from '@angular/core';
import { Observable } from 'rxjs';
import { ICreateSplitRequest, ICreateSplitResponse, ISplit, ISplitListItem } from './dto';

export interface ISplitusService {
  /** REAL: POST /api4splitus/create-split. Payer is the authenticated user. */
  createSplit(request: ICreateSplitRequest): Observable<ICreateSplitResponse>;
  /** REAL: GET /api4splitus/split?spaceID=&id= */
  getSplit(spaceID: string, id: string): Observable<ISplit>;
  /** REAL: GET /api4splitus/splits?spaceID= */
  getSplits(spaceID: string): Observable<ISplitListItem[]>;
}

export const SPLITUS_SERVICE = new InjectionToken<ISplitusService>('SplitusService');
