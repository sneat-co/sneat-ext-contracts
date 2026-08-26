import { IRecord } from '@sneat/data';
import { IWithSpaceIDs, SneatRecordStatus, IWithCreated, IWithRestrictions, IShortSpaceInfo } from '@sneat/dto';
import { ListType } from './list-types';
import { IListItemBrief, ListCounts, ListStatus } from './list-item';

export interface IListCommon {
  title: string;
  img?: string;
  emoji?: string;
  isDone?: boolean;
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

export interface IListKey {
  id: string;
  type: ListType;
}

export interface IListGroup {
  id: string;
  title?: string;
  type?: ListType;
  emoji?: string;
  lists?: IListInfo[];
}
