// Proves the ported asset-group DTO preserves the legacy
// legacy mod-assetus-core surface (IAssetDtoGroup + IAssetDtoGroupCounts).
import { describe, expect, it } from 'vitest';
import { type IAssetDtoGroup, type IAssetDtoGroupCounts } from './group';

describe('IAssetDtoGroup', () => {
  it('carries the legacy group fields', () => {
    const counts: IAssetDtoGroupCounts = { assets: 3 };
    const group: IAssetDtoGroup = {
      id: 'g1',
      title: 'Vehicles',
      order: 1,
      desc: 'All vehicles',
      categoryId: 'vehicles',
      numberOf: counts,
      spaceIDs: ['space1'],
    };
    expect(group.id).toBe('g1');
    expect(group.title).toBe('Vehicles');
    expect(group.order).toBe(1);
    expect(group.categoryId).toBe('vehicles');
    expect(group.numberOf?.assets).toBe(3);
    expect(group.spaceIDs).toEqual(['space1']);
  });

  it('allows the minimal required shape', () => {
    const group: IAssetDtoGroup = { id: 'g2', order: 0 };
    expect(group.id).toBe('g2');
    expect(group.numberOf).toBeUndefined();
  });
});
