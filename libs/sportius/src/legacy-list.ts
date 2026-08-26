import { InjectionToken } from '@angular/core';
import type {
  ICommuneDbo,
  IShortSpaceInfo,
  IWithCreated,
  IWithRestrictions,
  IWithSpaceIDs,
  SneatRecordStatus,
} from '@sneat/dto';
import type {
  ISpaceContext,
  ISpaceItemNavContext,
  ISpaceRequest,
} from '@sneat/space-models';
import type { Observable } from 'rxjs';

// Compatibility surface for the existing Sportius list UI. These names remain
// additive: the team and club contract in index.ts is the current API.
export type ListStatus = SneatRecordStatus;

export type ListType =
  | 'buy'
  | 'watch'
  | 'cook'
  | 'do'
  | 'other'
  | 'recipes'
  | 'rsvp';

export interface IQuantity {
  value: number;
  unit: string;
}

export interface IListCommon {
  title: string;
  img?: string;
  emoji?: string;
  isDone?: boolean;
}

export interface IListItemCommon extends IListCommon {
  subListId?: string;
  subListType?: ListType;
  quantity?: IQuantity;
  category?: string;
}

export type IListItemBase = IListItemCommon;
export type ListItemStatus = 'done' | 'active';

export interface IListItemBrief extends IListItemBase {
  id: string;
  readonly created?: string;
  readonly emoji?: string;
  readonly status?: ListItemStatus;
  readonly img?: string;
}

export interface ListCounts {
  active?: number;
  completed?: number;
}

export interface IListBase extends IListCommon, IWithSpaceIDs {
  type: ListType;
  shortId?: string;
  status?: ListStatus;
}

export interface IListDbo extends IListBase, IWithRestrictions, IWithCreated {
  dtClosed?: number;
  note?: string;
  numberOf?: ListCounts;
  items?: IListItemBrief[];
  commune?: IShortSpaceInfo;
}

export interface IListItemDbo extends IListBase, IListItemCommon {
  listId?: string;
  score?: number;
  subListItems?: IListItemBrief[];
}

export interface IListInfo extends IWithRestrictions {
  parentListId?: string;
  parentListType?: ListType;
  type: ListType;
  id?: string;
  shortId?: string;
  title?: string;
  hidden?: boolean;
  space?: IShortSpaceInfo;
  emoji?: string;
  img?: string;
  note?: string;
  itemsCount?: number;
}

export interface IListBrief extends IListBase, IWithCreated {
  emoji?: string;
}

export interface IListGroup {
  id: string;
  title?: string;
  type?: ListType;
  emoji?: string;
  lists?: IListInfo[];
}

export interface ISportiusSpaceDbo {
  listGroups?: IListGroup[];
}

export interface IListContext extends ISpaceItemNavContext<
  IListBrief,
  IListDbo
> {
  type: ListType;
}

export interface IListItemResult {
  message?: string;
  changed?: boolean;
  success: boolean;
  listDto: IListDbo;
  communeDto?: ICommuneDbo;
  listItemDto?: IListItemDbo;
}

export interface IListItemsCommandParams {
  space: ISpaceContext;
  list: IListContext;
  items: IListItemBrief[];
}

export interface ICreateListRequest extends ISpaceRequest, IListBrief {}

export interface IListRequest extends ISpaceRequest {
  readonly listID: string;
}

export interface ICreateListItemRequest extends IListItemBase {
  id: string;
}

export interface ICreateListItemsRequest extends IListRequest {
  items: ICreateListItemRequest[];
}

export interface IListItemRequest extends IListRequest {
  itemID: string;
}

export interface IListItemIDsRequest extends IListRequest {
  readonly itemIDs: string[];
}

export interface IReorderListItemsRequest extends IListItemIDsRequest {
  toIndex: number;
}

export type IDeleteListItemsRequest = IListItemIDsRequest;

export interface ISetListItemsIsComplete extends IListItemIDsRequest {
  isDone: boolean;
}

export interface ISportiusService {
  createList(request: ICreateListRequest): Observable<IListContext>;
  deleteList(space: ISpaceContext, listId: string): Observable<void>;
  reorderListItems(request: IReorderListItemsRequest): Observable<void>;
  createListItems(
    params: IListItemsCommandParams,
  ): Observable<IListItemResult>;
  setListItemsIsCompleted(
    request: ISetListItemsIsComplete,
  ): Observable<void>;
  deleteListItems(request: IDeleteListItemsRequest): Observable<void>;
  getListById(
    space: ISpaceContext,
    listType: ListType,
    listID: string,
  ): Observable<IListContext>;
}

export const SPORTIUS_SERVICE = new InjectionToken<ISportiusService>(
  'SportiusService',
);
