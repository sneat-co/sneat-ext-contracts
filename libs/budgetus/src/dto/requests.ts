import { IListItemBrief, IListItemDbo } from './list-item';
import { IListBrief } from './list';
import { ListType } from './list-types';
import { ISpaceContext, ISpaceItemNavContext } from '@sneat/space-models';

export interface IListContext extends ISpaceItemNavContext<IListBrief, any> {
  type: ListType;
}

export interface ICreateListRequest {
  id?: string;
  title: string;
  type: ListType;
  emoji?: string;
}

export interface IListRequest {
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
