import { IRecord } from '@sneat/data';
import { IListDbo, IListInfo } from './dto/list';

export function getListShortUrlId(communeId: string, shortId?: string, id?: string): string | undefined {
  if (shortId) {
    return `${communeId}-${shortId}`;
  }
  if (id) {
    return id;
  }
  return undefined;
}

export function isListInfoMatchesListDto(i: IListInfo, l: IRecord<IListDbo>): boolean {
  return (
    (!!i.id && i.id === l.id) ||
    (i.type === l.dbo?.type && !!i.shortId && i.shortId === l.dbo?.shortId)
  );
}

export function createListInfoFromDto(dto: IListDbo, shortId?: string): IListInfo {
  if (!dto.title) {
    throw new Error('!title');
  }
  const listInfo: IListInfo = {
    type: dto.type,
    title: dto.title,
  };
  if (shortId) {
    listInfo.shortId = shortId;
  }
  if (dto.items && dto.items.length) {
    listInfo.itemsCount = dto.items.length;
  }
  if (dto.emoji) {
    listInfo.emoji = dto.emoji;
  }
  if (dto.restrictions) {
    listInfo.restrictions = dto.restrictions;
  }
  if (dto.commune) {
    listInfo.space = dto.commune;
  }
  return listInfo;
}

export function monthISOOf(dateISO: string): string {
  return dateISO.slice(0, 7);
}

export function maskSurpriseLineItems(
  items: any[],
  options?: {
    reveal?: boolean;
  }
): any[] {
  if (options?.reveal) {
    return items;
  }
  return items.map((item) =>
    item.isSurprise
      ? { ...item, title: '🎁 Hidden surprise', sourceRef: undefined }
      : item
  );
}
