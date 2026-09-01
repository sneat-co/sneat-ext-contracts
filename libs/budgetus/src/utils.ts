import { IRecord } from '@sneat/data';
import { IListDbo, IListInfo } from './dto/list';
import { IBudgetLineItem, SHARED_BUDGET_MEMBER_ID } from './dto/budget';

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

/**
 * The canonical id of a budget line produced by an expense-flagged happening
 * price: `happening:{happeningID}:{priceID}:{occurrenceMonthISO}`.
 *
 * `priceID` is backend-stable, so a happening with several flagged prices yields
 * one stable line per price per month. Both producers and consumers build ids
 * through this function — an id built by hand somewhere else is how the two
 * halves of an override lookup silently stop matching.
 */
export function happeningBudgetLineID(
  happeningID: string,
  priceID: string,
  occurrenceMonthISO: string,
): string {
  return `happening:${happeningID}:${priceID}:${occurrenceMonthISO}`;
}

/** The minimum shape {@link maskSurpriseLineItems} needs to mask an item. */
export interface IMaskableLineItem {
  title: string;
  sourceRef?: string;
  isSurprise?: boolean;
}

/**
 * Hides the title and source of items flagged as a surprise, so a budget can be
 * shown to someone the surprise is for.
 *
 * Returns new objects; the input array and its items are never mutated. The
 * item type is preserved, so callers keep every field they passed in.
 */
export function maskSurpriseLineItems<T extends IMaskableLineItem>(
  items: readonly T[],
  options?: {
    reveal?: boolean;
  },
): T[] {
  if (options?.reveal) {
    return [...items];
  }
  return items.map((item) =>
    item.isSurprise
      ? { ...item, title: '🎁 Hidden surprise', sourceRef: undefined }
      : item,
  );
}

/**
 * Splits a value in minor units across members so the parts sum EXACTLY back to
 * it.
 *
 * A shared cost belongs to its participants equally; the odd minor units go to
 * the first members in the given order, so the split is deterministic rather
 * than arbitrary. A cost with no participants belongs wholly to
 * {@link SHARED_BUDGET_MEMBER_ID}.
 *
 * This lives in the contract because it is an invariant of the rollup shape,
 * not an implementation detail: a per-member breakdown must always add up to the
 * currency total it sits under. Two implementations of this rule would be two
 * chances for those numbers to disagree on screen.
 */
export function splitMinorUnitsAcrossMembers(
  value: number,
  memberIDs: readonly string[],
): Map<string, number> {
  const shares = new Map<string, number>();
  if (memberIDs.length === 0) {
    shares.set(SHARED_BUDGET_MEMBER_ID, value);
    return shares;
  }
  const base = Math.trunc(value / memberIDs.length);
  let remainder = value - base * memberIDs.length;
  for (const memberID of memberIDs) {
    const extra = remainder > 0 ? 1 : 0;
    remainder -= extra;
    shares.set(memberID, base + extra);
  }
  return shares;
}

/**
 * One member's share of a single budget line, in minor units.
 *
 * Returns 0 for a line the member is not attributed to, and for a line that
 * carries no summable amount.
 */
export function memberShareOfLine(
  item: Pick<IBudgetLineItem, 'monthlyAmount' | 'memberIDs'>,
  memberID: string,
): number {
  return (
    splitMinorUnitsAcrossMembers(
      item.monthlyAmount?.value ?? 0,
      item.memberIDs ?? [],
    ).get(memberID) ?? 0
  );
}
