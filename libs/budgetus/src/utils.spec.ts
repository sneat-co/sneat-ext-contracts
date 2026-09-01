import { getListShortUrlId, happeningBudgetLineID, maskSurpriseLineItems, monthISOOf } from './utils';
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

describe('happeningBudgetLineID', () => {
  it('builds the canonical happening:{happeningID}:{priceID}:{monthISO} id', () => {
    expect(happeningBudgetLineID('hap1', 'week1', '2026-09')).toBe(
      'happening:hap1:week1:2026-09',
    );
  });

  it('gives one distinct id per price of the same happening in a month', () => {
    const dropIn = happeningBudgetLineID('chess', 'single1', '2026-09');
    const pass = happeningBudgetLineID('chess', 'month1', '2026-09');
    expect(dropIn).not.toBe(pass);
  });

  it('gives one distinct id per month for the same price', () => {
    expect(happeningBudgetLineID('chess', 'week1', '2026-09')).not.toBe(
      happeningBudgetLineID('chess', 'week1', '2026-10'),
    );
  });
});

describe('maskSurpriseLineItems typing', () => {
  it('preserves fields the caller passed in that masking does not touch', () => {
    const items: IBudgetLineItem[] = [
      {
        id: 'happening:hap1:single1:2026-10',
        title: 'Surprise trip',
        dateISO: '2026-10-12',
        amount: { currency: 'EUR', value: 15000 },
        source: 'gift',
        sourceRef: 'hap1',
        isSurprise: true,
        happeningID: 'hap1',
        priceID: 'single1',
        occurrenceMonthISO: '2026-10',
        memberIDs: ['contact1'],
      },
    ];
    const [masked] = maskSurpriseLineItems(items);
    expect(masked.title).toBe('🎁 Hidden surprise');
    expect(masked.sourceRef).toBeUndefined();
    // Everything else survives masking.
    expect(masked.happeningID).toBe('hap1');
    expect(masked.priceID).toBe('single1');
    expect(masked.memberIDs).toEqual(['contact1']);
    expect(masked.amount).toEqual({ currency: 'EUR', value: 15000 });
  });

  it('returns a new array even when revealing, so callers cannot mutate the source', () => {
    const items: IBudgetLineItem[] = [];
    expect(maskSurpriseLineItems(items, { reveal: true })).not.toBe(items);
  });
});
