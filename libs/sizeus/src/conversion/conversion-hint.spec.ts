import { describe, expect, it } from 'vitest';
import { conversionHint } from './conversion-hint';
import { conversionTableFor } from './all';

describe('conversionHint', () => {
  it('returns a reliable published conversion', () => {
    expect(conversionHint('footwear-adult', 'EU', '42', 'UK')).toBe('8');
  });

	it('does not invent a conversion for unreliable, unknown or identical pairs', () => {
    expect(conversionHint('footwear-kids', 'EU', '28', 'US')).toBeNull();
    expect(conversionHint('ski-boots', 'EU', '42', 'UK')).toBeNull();
    expect(conversionHint('footwear-adult', 'EU', '42', 'EU')).toBeNull();
	});

	it('exposes only the published conversion tables', () => {
		expect(conversionTableFor('footwear-adult')?.systems).toEqual([
			'EU',
			'UK',
			'US',
		]);
		expect(conversionTableFor('ski-boots')).toBeUndefined();
	});
});
