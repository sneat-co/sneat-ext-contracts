import { describe, expect, it } from 'vitest';
import { MVP_SIZE_TYPE_CATALOG } from './mvp-catalog';

describe('MVP_SIZE_TYPE_CATALOG', () => {
  it('is versioned and has unique size-type IDs', () => {
    expect(MVP_SIZE_TYPE_CATALOG.version).toBe(1);
    const ids = MVP_SIZE_TYPE_CATALOG.sizeTypes.map(({ id }) => id);
    expect(new Set(ids).size).toBe(ids.length);
  });

  it('covers each advertised category and resolves sport-kit members', () => {
    for (const category of [
      'footwear',
      'body-measurement',
      'tops',
      'bottoms',
      'accessories',
    ]) {
      expect(MVP_SIZE_TYPE_CATALOG.sizeTypes.some((entry) => entry.category === category)).toBe(true);
    }
    const ids = new Set(MVP_SIZE_TYPE_CATALOG.sizeTypes.map(({ id }) => id));
    for (const kit of MVP_SIZE_TYPE_CATALOG.sportKits) {
      expect(kit.sizeTypeIDs.every((id) => ids.has(id))).toBe(true);
    }
  });
});
