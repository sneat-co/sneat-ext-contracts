import { InjectionToken } from '@angular/core';
import { ISpaceContext } from '@sneat/space-models';
import { Observable } from 'rxjs';
import { IListContext } from '../contexts';
import { ListType } from '../dto';
import { ICreateListRequest, IDeleteListItemsRequest, IListItemResult, IListItemsCommandParams, IReorderListItemsRequest, ISetListItemsIsComplete } from './interfaces';

/** The public service surface implemented by RSVP Express's private runtime. */
export interface IRsvpexpressService {
  createList(request: ICreateListRequest): Observable<IListContext>;
  deleteList(space: ISpaceContext, listId: string): Observable<void>;
  reorderListItems(request: IReorderListItemsRequest): Observable<void>;
  createListItems(params: IListItemsCommandParams): Observable<IListItemResult>;
  setListItemsIsCompleted(request: ISetListItemsIsComplete): Observable<void>;
  deleteListItems(request: IDeleteListItemsRequest): Observable<void>;
  getListById(space: ISpaceContext, listType: ListType, listID: string): Observable<IListContext>;
}

export const RSVPEXPRESS_SERVICE = new InjectionToken<IRsvpexpressService>('RsvpexpressService');
