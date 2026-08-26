import { ICommuneDbo } from '@sneat/dto';
import { IListBrief, IListDbo, IListItemBase, IListItemBrief, IListItemDbo } from '../dto';
import { IListContext } from '../contexts';
import { ISpaceContext, ISpaceRequest } from '@sneat/space-models';

export interface GetOrCreateCommuneItemIds { id?: string; shortId?: string; communeShortId?: string; }
export interface IProgress { current: number; total: number; state?: string; }
export interface IListItemResult { message?: string; changed?: boolean; success: boolean; listDto: IListDbo; communeDto?: ICommuneDbo; listItemDto?: IListItemDbo; }
export interface IListItemsCommandParams { space: ISpaceContext; list: IListContext; items: IListItemBrief[]; }
export type ReorderListItemsWorker = (listDto: IListDbo) => void;
export interface ICreateListRequest extends ISpaceRequest, IListBrief {}
export interface IListRequest extends ISpaceRequest { readonly listID: string; }
export interface ICreateListItemRequest extends IListItemBase { id: string; }
export interface ICreateListItemsRequest extends IListRequest { items: ICreateListItemRequest[]; }
export interface IListItemRequest extends IListRequest { itemID: string; }
export interface IListItemIDsRequest extends IListRequest { readonly itemIDs: string[]; }
export interface IReorderListItemsRequest extends IListItemIDsRequest { toIndex: number; }
export type IDeleteListItemsRequest = IListItemIDsRequest;
export interface ISetListItemsIsComplete extends IListItemIDsRequest { isDone: boolean; }
