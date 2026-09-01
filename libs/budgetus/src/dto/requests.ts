import { IListItemBrief, IListItemDbo } from './list-item';
import { IListBrief } from './list';
import { ListType } from './list-types';
import {
  ISpaceContext,
  ISpaceItemNavContext,
  ISpaceRequest,
} from '@sneat/space-models';

export interface IListContext extends ISpaceItemNavContext<IListBrief, any> {
  type: ListType;
}

// Every list request addresses a list inside a space, so it must carry the
// spaceID. `extends ISpaceRequest` was dropped from both of these between 0.1.0
// and 0.1.2, which silently removed spaceID from the whole IListRequest family
// (ICreateListItemsRequest, IListItemIDsRequest, IReorderListItemsRequest, ...)
// even though the service still sends it and the backend still routes on it.
export interface ICreateListRequest extends ISpaceRequest {
  id?: string;
  title: string;
  type: ListType;
  emoji?: string;
}

export interface IListRequest extends ISpaceRequest {
  readonly listID: string;
}

export interface ICreateListItemRequest extends IListItemBrief {
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

export interface IListItemsCommandParams {
  space: ISpaceContext;
  list: IListContext;
  items: IListItemBrief[];
}

export interface IListItemResult {
  message?: string;
  changed?: boolean;
  success: boolean;
  listDto: any;
  communeDto?: any;
  listItemDto?: IListItemDbo;
}
