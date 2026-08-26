import { InjectionToken } from '@angular/core';
import { ISpaceContext } from '@sneat/space-models';
import { Observable } from 'rxjs';
import { IListContext } from '../contexts';
import { ListType } from '../dto';
import { ICreateListRequest, IDeleteListItemsRequest, IListItemResult, IListItemsCommandParams, IReorderListItemsRequest, ISetListItemsIsComplete } from './interfaces';

/** The public service surface implemented by Taxus's private runtime. */
export interface ITaxusService {
  createList(request: ICreateListRequest): Observable<IListContext>;
  deleteList(space: ISpaceContext, listId: string): Observable<void>;
  reorderListItems(request: IReorderListItemsRequest): Observable<void>;
  createListItems(params: IListItemsCommandParams): Observable<IListItemResult>;
  setListItemsIsCompleted(request: ISetListItemsIsComplete): Observable<void>;
  deleteListItems(request: IDeleteListItemsRequest): Observable<void>;
  getListById(space: ISpaceContext, listType: ListType, listID: string): Observable<IListContext>;
}
export const TAXUS_SERVICE = new InjectionToken<ITaxusService>('TaxusService');
