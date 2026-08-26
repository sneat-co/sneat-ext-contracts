import { InjectionToken } from '@angular/core';
import { Observable } from 'rxjs';
import {
  IListDbo,
  IListBrief,
  IListKey,
  IBudgetRollup,
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
  deleteList(space: any, listId: string): Observable<void>;
  reorderListItems(request: IReorderListItemsRequest): Observable<void>;
  createListItems(params: IListItemsCommandParams): Observable<IListItemResult>;
  setListItemsIsCompleted(request: ISetListItemsIsComplete): Observable<void>;
  deleteListItems(request: IDeleteListItemsRequest): Observable<void>;
  getListById(space: any, listType: ListType, listID: string): Observable<IListContext>;
  watchBudget(spaceID: string): Observable<IBudgetRollup>;
  setOverride(spaceID: string, lineItemId: string, patch: IBudgetOverridePatch): Promise<void>;
}

export const BUDGETUS_SERVICE = new InjectionToken<IBudgetusService>('BudgetusService');
