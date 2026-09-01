import { InjectionToken } from '@angular/core';
import { Observable } from 'rxjs';
import { ISpaceContext } from '@sneat/space-models';
import {
  IListDbo,
  IListBrief,
  IListKey,
  IBudgetRollup,
  IBudgetWindow,
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
   * Watches the space's budget rollup.
   *
   * `window` defaults to {@link DEFAULT_BUDGET_WINDOW_MONTHS} months starting at
   * the current month. The window actually used is echoed back on the rollup.
   */
  watchBudget(spaceID: string, window?: IBudgetWindow): Observable<IBudgetRollup>;

  /**
   * Watches the rollup for a space the caller already holds the full context of.
   *
   * A data source may need more than an id — calendarius's HappeningService
   * takes an ISpaceContext — and a caller inside a space page already has the
   * loaded brief/dbo. Going through `watchBudget(spaceID)` throws that away and
   * rebuilds a bare `{ id }`.
   *
   * OPTIONAL so that adding it does not break existing implementations. Callers
   * should use it when present and fall back to `watchBudget`.
   */
  watchBudgetForSpace?(
    space: ISpaceContext,
    window?: IBudgetWindow,
  ): Observable<IBudgetRollup>;
  setOverride(spaceID: string, lineItemId: string, patch: IBudgetOverridePatch): Promise<void>;
}

export const BUDGETUS_SERVICE = new InjectionToken<IBudgetusService>('BudgetusService');
