import { getListShortUrlId, maskSurpriseLineItems, monthISOOf } from './utils';
import { IBudgetLineItem } from './dto/budget';

describe('getListShortUrlId', () => {
  it('combines spaceId and shortId when a shortId is given', () => {
    expect(getListShortUrlId('fam1', 'groceries')).toBe('fam1-groceries');
  });

  it('falls back to the full id when there is no shortId', () => {
    expect(getListShortUrlId('fam1', undefined, 'to-do:abc')).toBe('to-do:abc');
  });

  it('returns undefined when neither shortId nor id is provided', () => {
    expect(getListShortUrlId('fam1')).toBeUndefined();
  });
});

describe('monthISOOf', () => {
  it('extracts the YYYY-MM prefix of an ISO date', () => {
    expect(monthISOOf('2026-09-12')).toBe('2026-09');
  });
});

describe('maskSurpriseLineItems', () => {
  const items: IBudgetLineItem[] = [
    {
      id: 'happening:mum-bday',
      title: "Mum's 60th birthday",
      dateISO: '2026-10-12',
      amount: { currency: 'EUR', value: 0 },
      source: 'gift',
      sourceRef: 'mum-bday',
      targetAmount: { currency: 'EUR', value: 150 },
      isSurprise: true,
    },
    {
      id: 'asset-renewal:car',
      title: 'Car — insurance renewal',
      dateISO: '2026-09-01',
      amount: { currency: 'EUR', value: 620 },
      source: 'asset-renewal',
    },
  ];

  it('masks the title and sourceRef of surprise items by default', () => {
    const masked = maskSurpriseLineItems(items);
    expect(masked[0].title).toBe('🎁 Hidden surprise');
    expect(masked[0].sourceRef).toBeUndefined();
    // Non-surprise items pass through unchanged.
    expect(masked[1]).toEqual(items[1]);
  });

  it('does not mutate the input array/items', () => {
    maskSurpriseLineItems(items);
    expect(items[0].title).toBe("Mum's 60th birthday");
  });

  it('leaves everything unmasked when reveal is true', () => {
    const revealed = maskSurpriseLineItems(items, { reveal: true });
    expect(revealed).toEqual(items);
  });
});
