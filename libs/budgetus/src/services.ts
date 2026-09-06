import { InjectionToken } from '@angular/core';
import { Observable } from 'rxjs';
import { ISpaceContext } from '@sneat/space-models';
import {
  IListDbo,
  IListBrief,
  IListKey,
  IBudgetOverridePatch,
  ListType,
  IListItemBrief,
  IListContext,
  ICreateListRequest,
  IListRequest,
  ICreateListItemsRequest,
  IListItemIDsRequest,
  IReorderListItemsRequest,
  IDeleteListItemsRequest,
  ISetListItemsIsComplete,
  IListItemsCommandParams,
  IListItemResult,
} from './dto';

export interface GetOrCreateCommuneItemIds {
  id?: string;
  shortId?: string;
  communeShortId?: string;
}

export interface IProgress {
  current: number;
  total: number;
  state?: string;
}

export type ReorderListItemsWorker = (listDto: IListDbo) => void;

export interface IBudgetusService {
  createList(request: ICreateListRequest): Observable<IListContext>;
  deleteList(space: ISpaceContext, listId: string): Observable<void>;
  reorderListItems(request: IReorderListItemsRequest): Observable<void>;
  createListItems(params: IListItemsCommandParams): Observable<IListItemResult>;
  setListItemsIsCompleted(request: ISetListItemsIsComplete): Observable<void>;
  deleteListItems(request: IDeleteListItemsRequest): Observable<void>;
  getListById(
    space: ISpaceContext,
    listType: ListType,
    listID: string,
  ): Observable<IListContext>;
  /**
   * Watches local courtesy-masking preferences for financial source effects.
   *
   * The record is keyed by a Budgetus override id. A `true` value asks the
   * current device to mask that effect in ordinary presentation. These flags
   * carry no financial amount, permission, or server-side business state.
   */
  watchLocalCourtesyFlags?(
    spaceID: string,
  ): Observable<Readonly<Record<string, boolean>>>;
  setOverride(spaceID: string, lineItemId: string, patch: IBudgetOverridePatch): Promise<void>;
}

export const BUDGETUS_SERVICE = new InjectionToken<IBudgetusService>('BudgetusService');
