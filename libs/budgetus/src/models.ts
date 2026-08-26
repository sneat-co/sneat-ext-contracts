import { IListItemBrief } from './dto/list-item';

export class ListItemInfoModel {
  static trackBy = (index: number, item: IListItemBrief | null | undefined): string | number | undefined =>
    !item
      ? index
      : (!!item.id && `id:${item.id}`) ||
        (item.subListId && `subList:${item.subListId}`) ||
        item.title;
}

export class ListItemModel {
  static equalListItems(...items: IListItemBrief[]): boolean {
    const { id, title, subListId, category, subListType } = items[0];
    return !items.some((item) => {
      if (id) {
        return item.id !== id;
      }
      return (
        (!!title && item.title !== title) ||
        (!!subListId && item.subListId !== subListId) ||
        (!!category && item.category !== category) ||
        (!!subListType && item.subListType !== subListType)
      );
    });
  }
}
