import { IRecord } from '@sneat/data';
import { IWithSpaceIDs, SneatRecordStatus, IWithCreated, IWithRestrictions } from '@sneat/dto';
import { ListType } from './list-types';

export type ListStatus = SneatRecordStatus;

export interface IQuantity {
  value: number;
  unit: string;
}

export interface IListItemCommon {
  subListId?: string;
  subListType?: ListType;
  quantity?: IQuantity;
  category?: string;
  title: string;
  img?: string;
  emoji?: string;
  isDone?: boolean;
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

export interface IListItemDbo extends IWithSpaceIDs, IListItemCommon {
  listId?: string;
  score?: number;
  subListItems?: IListItemBrief[];
}
