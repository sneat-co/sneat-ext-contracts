import { IRecord } from '@sneat/data';
import { IShortSpaceInfo, IWithCreated, IWithRestrictions, IWithSpaceIDs, SneatRecordStatus } from '@sneat/dto';

export type ListStatus = SneatRecordStatus;
export interface IQuantity { value: number; unit: string; }
export interface IListCommon { title: string; img?: string; emoji?: string; isDone?: boolean; }
export type ListType = 'buy' | 'watch' | 'cook' | 'do' | 'other' | 'recipes' | 'rsvp';
export type ListItemStatus = 'done' | 'active';
export interface IListItemCommon extends IListCommon { subListId?: string; subListType?: ListType; quantity?: IQuantity; category?: string; }
export type IListItemBase = IListItemCommon;
export interface IListItemBrief extends IListItemBase { id: string; readonly created?: string; readonly emoji?: string; readonly status?: ListItemStatus; readonly img?: string; }
export interface ListCounts { active?: number; completed?: number; }
export interface IListBase extends IListCommon, IWithSpaceIDs { type: ListType; shortId?: string; status?: ListStatus; }
export interface IListDbo extends IListBase, IWithRestrictions, IWithCreated { dtClosed?: number; note?: string; numberOf?: ListCounts; items?: IListItemBrief[]; commune?: IShortSpaceInfo; }
export class ListItemInfoModel {
  static trackBy: (index: number, item: IListItemBrief) => string | number | undefined = (index, item) => !item ? index : (!!item.id && `id:${item.id}`) || (item.subListId && `subList:${item.subListId}`) || item.title;
}
export class ListItemModel {
  static equalListItems(...items: IListItemBrief[]): boolean {
    const { id, title, subListId, category, subListType } = items[0];
    return !items.some((item) => id ? item.id !== id : ((!!title && item.title !== title) || (!!subListId && item.subListId !== subListId) || (!!category && item.category !== category) || (!!subListType && item.subListType !== subListType)));
  }
}
export interface IListItemDbo extends IListBase, IListItemCommon { listId?: string; score?: number; subListItems?: IListItemBrief[]; }
export function getListShortUrlId(communeId: string, shortId?: string, id?: string): string | undefined { return shortId ? `${communeId}-${shortId}` : id; }
export interface IListInfo extends IWithRestrictions { parentListId?: string; parentListType?: ListType; type: ListType; id?: string; shortId?: string; title?: string; hidden?: boolean; space?: IShortSpaceInfo; emoji?: string; img?: string; note?: string; itemsCount?: number; }
export interface IListBrief extends IListBase, IWithCreated { emoji?: string; }
export function isListInfoMatchesListDto(i: IListInfo, l: IRecord<IListDbo>): boolean { return (!!i.id && i.id === l.id) || (i.type === l.dbo?.type && !!i.shortId && i.shortId === l.dbo?.shortId); }
export function createListInfoFromDto(dto: IListDbo, shortId?: string): IListInfo {
  if (!dto.title) throw new Error('!title');
  const listInfo: IListInfo = { type: dto.type, title: dto.title };
  if (shortId) listInfo.shortId = shortId;
  if (dto.items?.length) listInfo.itemsCount = dto.items.length;
  if (dto.emoji) listInfo.emoji = dto.emoji;
  if (dto.restrictions) listInfo.restrictions = dto.restrictions;
  if (dto.commune) listInfo.space = dto.commune;
  return listInfo;
}
