// AC: relationships-frontend-preserved
// Proves the ported relationship DTOs preserve and resolve, with the exact
// backend json names, an asset that: belongs to a group (groupID + group
// sub-entity), has two sub-assets under a parentAssetID, is linked via
// relatedAs and sameAssetID, is associated with multiple spaces, and carries
// memberIDs + membersInfo.
import {
  type IAssetDbo,
  type IAssetGroupInfo,
  type ISubAssetInfo,
  type IWithAssetSpaces,
  type IAssetusSpaceBrief,
} from './asset';
import { type ITitledRecord } from '@sneat/dto';

describe('relationship DTOs (relationships-frontend-preserved)', () => {
  const asset: IAssetDbo = {
    id: 'a1',
    name: 'Family car',
    category: 'vehicles',
    condition: 'good',
    status: 'active',
    visibility: 'family',

    // group as a sub-entity (id/title/order/desc/categoryID/numberOf/totals)
    groupID: 'g1',
    group: {
      id: 'g1',
      title: 'Vehicles',
      order: 2,
      desc: 'Household vehicles',
      categoryID: 'vehicles',
      numberOf: { assets: 5 },
      totals: [{ currency: 'EUR', value: 25000 }],
    },

    // parent/sub-asset nesting: two sub-assets under a parent
    parentAssetID: 'parent1',
    subAssets: [
      {
        id: 'sa1',
        title: 'Winter tyres',
        type: 'other',
        countryID: 'IE',
        subType: 'tyres',
        expires: '2030-01-01',
      },
      {
        id: 'sa2',
        title: 'Insurance doc',
        type: 'document',
        countryID: 'IE',
        subType: 'insurance',
        expires: '2026-12-31',
      },
    ],

    // linking
    sameAssetID: 'same1',
    relatedAs: 'tenant',

    // multi-space association (additive to the single owning spaceIDs)
    spaceIDs: ['s1'],
    spaces: {
      s1: { assets: { a1: { id: 'a1', name: 'Family car', category: 'vehicles', condition: 'good', status: 'active', visibility: 'family' } } },
      s2: { assets: { a1: { id: 'a1', name: 'Family car', category: 'vehicles', condition: 'good', status: 'active', visibility: 'family' } } },
    },

    // member info
    memberIDs: ['m1', 'm2'],
    membersInfo: [
      { id: 'm1', title: 'Alice' },
      { id: 'm2', title: 'Bob' },
    ],
  };

  it('round-trips an asset with every relationship dimension through JSON', () => {
    const roundTripped: IAssetDbo = JSON.parse(JSON.stringify(asset));
    expect(roundTripped).toEqual(asset);
  });

  it('resolves the group sub-entity with backend json names', () => {
    const group = asset.group as IAssetGroupInfo;
    expect(group.id).toBe('g1');
    expect(group.title).toBe('Vehicles');
    expect(group.order).toBe(2);
    expect(group.desc).toBe('Household vehicles');
    expect(group.categoryID).toBe('vehicles');
    expect(group.numberOf?.assets).toBe(5);
    expect(group.totals?.[0]).toEqual({ currency: 'EUR', value: 25000 });
    // groupID linkage matches the group sub-entity id.
    expect(asset.groupID).toBe(group.id);
  });

  it('resolves two sub-assets under a parent with per-sub-asset detail', () => {
    expect(asset.parentAssetID).toBe('parent1');
    const subs = asset.subAssets as ISubAssetInfo[];
    expect(subs).toHaveLength(2);
    expect(subs[0]).toEqual({
      id: 'sa1',
      title: 'Winter tyres',
      type: 'other',
      countryID: 'IE',
      subType: 'tyres',
      expires: '2030-01-01',
    });
    expect(subs[1].type).toBe('document');
    expect(subs[1].expires).toBe('2026-12-31');
  });

  it('resolves the relatedAs and sameAssetID links', () => {
    expect(asset.sameAssetID).toBe('same1');
    expect(asset.relatedAs).toBe('tenant');
  });

  it('resolves the multi-space association keyed by spaceID', () => {
    const withSpaces: IWithAssetSpaces = asset;
    expect(Object.keys(withSpaces.spaces ?? {})).toEqual(['s1', 's2']);
    const s1: IAssetusSpaceBrief = withSpaces.spaces!['s1'];
    expect(s1.assets?.['a1'].name).toBe('Family car');
    // multi-space is additive to the single owning spaceIDs.
    expect(asset.spaceIDs).toEqual(['s1']);
  });

  it('resolves memberIDs and membersInfo as ITitledRecord[]', () => {
    expect(asset.memberIDs).toEqual(['m1', 'm2']);
    const members = asset.membersInfo as ITitledRecord[];
    expect(members.map((m) => m.id)).toEqual(['m1', 'm2']);
    expect(members.map((m) => m.title)).toEqual(['Alice', 'Bob']);
  });

  it('accepts an asset with no relationship fields (all optional)', () => {
    const minimal: IAssetDbo = {
      id: 'a2',
      name: 'Bike',
      category: 'sports_equipment',
      condition: 'new',
      status: 'draft',
      visibility: 'private',
    };
    expect(minimal.group).toBeUndefined();
    expect(minimal.subAssets).toBeUndefined();
    expect(minimal.spaces).toBeUndefined();
    expect(minimal.relatedAs).toBeUndefined();
  });
});
