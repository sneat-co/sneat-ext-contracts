// Proves the ported AssetGroup uimodel preserves the legacy
// legacy mod-assetus-core behaviour (id, totals, numberOf).
import { describe, expect, it } from 'vitest';
import { AssetGroup } from './asset-group';
import { type IAssetGroupContext } from '../contexts';

describe('AssetGroup', () => {
  it('exposes the context id', () => {
    const context: IAssetGroupContext = {
      id: 'group-1',
      dbo: { id: 'group-1', order: 0 },
    };
    const group = new AssetGroup(context);
    expect(group.id).toBe('group-1');
    expect(group.context).toBe(context);
  });

  it('builds Totals from the dbo totals', () => {
    const context: IAssetGroupContext = {
      id: 'group-2',
      dbo: {
        id: 'group-2',
        order: 0,
        totals: { incomes: { count: 1 }, expenses: { count: 0 } },
      },
    };
    const group = new AssetGroup(context);
    expect(group.totals).toBeDefined();
    expect(group.totals.count).toBe(1);
  });

  it('returns numberOf from the dbo, defaulting to {}', () => {
    const withCounts: IAssetGroupContext = {
      id: 'g3',
      dbo: { id: 'g3', order: 0, numberOf: { assets: 5 } },
    };
    expect(new AssetGroup(withCounts).numberOf).toEqual({ assets: 5 });

    const noDbo: IAssetGroupContext = { id: 'g4' };
    expect(new AssetGroup(noDbo).numberOf).toEqual({});
  });
});
