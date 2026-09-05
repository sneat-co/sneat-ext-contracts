import { InjectionToken } from '@angular/core';
import { Observable } from 'rxjs';
import {
  ICreateSplitRequest,
  ICreateSplitResponse,
  ICreateSplitusBillV1Request,
  ICreateSplitusBillV1Response,
  IGetSplitusBillV1Request,
  IGetSplitusBillV1Response,
  IListSplitusBillsV1Request,
  IListSplitusBillsV1Response,
  ISplit,
  ISplitListItem,
} from './dto';

/** @deprecated Use `ISplitusBillServiceV1`. */
export interface ISplitusService {
  /** REAL: POST /api4splitus/create-split. Payer is the authenticated user. */
  createSplit(request: ICreateSplitRequest): Observable<ICreateSplitResponse>;
  /** REAL: GET /api4splitus/split?spaceID=&id= */
  getSplit(spaceID: string, id: string): Observable<ISplit>;
  /** REAL: GET /api4splitus/splits?spaceID= */
  getSplits(spaceID: string): Observable<ISplitListItem[]>;
}

/** @deprecated Use `SPLITUS_BILL_SERVICE_V1`. */
export const SPLITUS_SERVICE = new InjectionToken<ISplitusService>('SplitusService');

export interface ISplitusBillServiceV1 {
  createBill(
    request: ICreateSplitusBillV1Request,
  ): Observable<ICreateSplitusBillV1Response>;
  getBill(
    request: IGetSplitusBillV1Request,
  ): Observable<IGetSplitusBillV1Response>;
  listBills(
    request: IListSplitusBillsV1Request,
  ): Observable<IListSplitusBillsV1Response>;
}

export const SPLITUS_BILL_SERVICE_V1 =
  new InjectionToken<ISplitusBillServiceV1>('SplitusBillServiceV1');
