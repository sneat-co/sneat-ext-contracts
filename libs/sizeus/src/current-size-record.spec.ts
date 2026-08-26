import { describe, expect, it } from 'vitest';
import { currentSizeRecord } from './current-size-record';
import type { ISizeRecord } from './size-record';

const record = (partial: Partial<ISizeRecord>): ISizeRecord => ({
  contactID: 'contact-1',
  sizeTypeID: 'shoe',
  value: '9',
  system: 'US',
  effectiveDate: '2026-01-01',
  createdAt: '2026-01-01T00:00:00.000Z',
  ...partial,
});

describe('currentSizeRecord', () => {
  it('selects the latest record effective today or earlier', () => {
    const expected = record({ effectiveDate: '2026-06-01', createdAt: '2026-06-01T09:00:00.000Z' });
    const future = record({ effectiveDate: '2026-12-01' });
    expect(currentSizeRecord([record({ effectiveDate: '2026-05-01' }), expected, future], '2026-07-02')).toBe(expected);
  });

  it('uses creation time to resolve an equal effective date without mutating input', () => {
    const first = record({ createdAt: '2026-06-01T08:00:00.000Z' });
    const expected = record({ createdAt: '2026-06-01T09:00:00.000Z' });
    const records = [first, expected];
    expect(currentSizeRecord(records, '2026-07-02')).toBe(expected);
    expect(records).toEqual([first, expected]);
  });
});
