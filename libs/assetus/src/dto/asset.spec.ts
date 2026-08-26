import {
  assetCategories,
  assetConditions,
  assetVisibilities,
  categoryOptions,
  conditionOptions,
  defaultVisibilityForSpaceType,
  deriveOwnerType,
  visibilityOptions,
  type AssetCategory,
  type AssetStatus,
  type AssetPossession,
  type EngineType,
  type FuelType,
  type FuelVolumeUnit,
  type MileageUnit,
  type IAssetDbo,
} from './asset';

describe('defaultVisibilityForSpaceType', () => {
  it('maps private space to private visibility', () => {
    expect(defaultVisibilityForSpaceType('private')).toBe('private');
  });
  it('maps family space to family visibility', () => {
    expect(defaultVisibilityForSpaceType('family')).toBe('family');
  });
  it('falls back to specific_space for other space types', () => {
    expect(defaultVisibilityForSpaceType('sports_club')).toBe('specific_space');
    expect(defaultVisibilityForSpaceType(undefined)).toBe('specific_space');
  });
});

describe('deriveOwnerType', () => {
  it('maps private to individual', () => {
    expect(deriveOwnerType('private')).toBe('individual');
  });
  it('maps family to family', () => {
    expect(deriveOwnerType('family')).toBe('family');
  });
  it('passes through known org-like space types', () => {
    expect(deriveOwnerType('sports_club')).toBe('sports_club');
    expect(deriveOwnerType('community')).toBe('community');
    expect(deriveOwnerType('school')).toBe('school');
  });
  it('defaults unknown space types to organisation', () => {
    expect(deriveOwnerType('whatever')).toBe('organisation');
  });
});

// AC: dto-superset-roundtrip
describe('IAssetDbo superset', () => {
  it('type-checks and round-trips a full object with every legacy optional', () => {
    const full: IAssetDbo = {
      // required MVP fields
      id: 'a1',
      name: 'My car',
      category: 'vehicles',
      condition: 'good',
      status: 'active',
      visibility: 'private',
      // retained MVP optionals
      description: 'desc',
      acquisitionDate: '2020-01-01',
      purchasePrice: { currency: 'EUR', value: 10000 },
      estimatedValue: { currency: 'EUR', value: 8000 },
      location: 'Dublin',
      notes: 'note',
      tags: ['x'],
      photos: ['p1'],
      createdAt: '2020-01-01T00:00:00Z',
      updatedAt: '2020-02-01T00:00:00Z',
      spaceIDs: ['s1'],
      // ported legacy optionals (backend json names)
      isRequest: true,
      countryID: 'IE',
      type: 'car',
      possession: 'owning',
      parentCategoryID: 'other',
      yearOfBuild: 2019,
      geo: { lat: 53.3, lng: -6.2 },
      dateOfBuild: '2019-01-01',
      dateOfPurchase: '2020-01-01',
      dateInsuredTill: '2025-01-01',
      dateCertifiedTill: '2024-01-01',
      totals: [{ currency: 'EUR', value: 100 }],
      canHaveIncome: true,
      canHaveExpense: false,
      financialDirection: 'income',
      liabilities: [{ id: 'l1', serviceTypes: ['electricity'] }],
      notUsedServiceTypes: ['gas'],
      groupID: 'g1',
      group: {
        id: 'g1',
        title: 'Group',
        order: 1,
        desc: 'd',
        categoryID: 'vehicles',
        numberOf: { assets: 3 },
        totals: [{ currency: 'EUR', value: 50 }],
      },
      parentAssetID: 'p1',
      subAssets: [
        {
          id: 'sa1',
          title: 'Sub',
          type: 'document',
          countryID: 'IE',
          subType: 'passport',
          expires: '2030-01-01',
        },
      ],
      sameAssetID: 'same1',
      relatedAs: 'tenant',
      memberIDs: ['m1', 'm2'],
      membersInfo: [{ id: 'm1', title: 'Member' }],
    };

    const roundTripped: IAssetDbo = JSON.parse(JSON.stringify(full));
    expect(roundTripped).toEqual(full);
  });

  it('accepts a minimal MVP-only object (every legacy field optional)', () => {
    const minimal: IAssetDbo = {
      id: 'a2',
      name: 'Bike',
      category: 'sports_equipment',
      condition: 'new',
      status: 'draft',
      visibility: 'private',
    };
    expect(minimal.id).toBe('a2');
    // No legacy optionals present.
    expect(minimal.possession).toBeUndefined();
    expect(minimal.groupID).toBeUndefined();
  });
});

// AC: enum-union-no-value-dropped
describe('enum unions admit every backend value', () => {
  it('AssetCategory admits document/debt/dwelling and MVP values', () => {
    const cats: AssetCategory[] = [
      'books',
      'games',
      'toys',
      'sports_equipment',
      'tools',
      'electronics',
      'appliances',
      'clothing',
      'vehicles',
      'camping_equipment',
      'other',
      'dwelling',
      'document',
      'debt',
      'utility',
      'valuables',
      'digital_asset',
    ];
    expect(assetCategories).toEqual(cats);
  });

  it('AssetStatus admits draft and disposed/lost/transferred', () => {
    const statuses: AssetStatus[] = [
      'draft',
      'active',
      'transferred',
      'archived',
      'disposed',
      'lost',
    ];
    expect(statuses).toContain('draft');
    expect(statuses).toContain('disposed');
    expect(statuses).toContain('lost');
    expect(statuses).toContain('transferred');
  });

  it('AssetPossession admits owning/leasing/renting/unknown/undisclosed', () => {
    const p: AssetPossession[] = [
      'unknown',
      'undisclosed',
      'owning',
      'leasing',
      'renting',
    ];
    expect(p).toHaveLength(5);
  });

  it('EngineType admits phev/hybrid/steam (and empty unknown)', () => {
    const e: EngineType[] = [
      '',
      'other',
      'combustion',
      'electric',
      'phev',
      'hybrid',
      'steam',
    ];
    expect(e).toContain('phev');
    expect(e).toContain('hybrid');
    expect(e).toContain('steam');
  });

  it('FuelType admits bio/hydrogen (and empty unknown)', () => {
    const f: FuelType[] = [
      '',
      'other',
      'bio',
      'petrol',
      'diesel',
      'hydrogen',
    ];
    expect(f).toContain('bio');
    expect(f).toContain('hydrogen');
  });

  it('fuel-volume unit / mileage unit / currency are open strings', () => {
    const vol: FuelVolumeUnit = 'gallons-or-anything';
    const mileage: MileageUnit = 'nautical-mile';
    const currency = 'XYZ';
    expect(typeof vol).toBe('string');
    expect(typeof mileage).toBe('string');
    expect(typeof currency).toBe('string');
  });
});

describe('select option lists', () => {
  it('produces an option per category/condition/visibility value', () => {
    expect(categoryOptions).toHaveLength(assetCategories.length);
    expect(conditionOptions).toHaveLength(assetConditions.length);
    expect(visibilityOptions).toHaveLength(assetVisibilities.length);
  });
  it('titleizes multi-word values for labels', () => {
    const sports = categoryOptions.find((o) => o.value === 'sports_equipment');
    expect(sports?.label).toBe('Sports Equipment');
    const repair = conditionOptions.find((o) => o.value === 'needs_repair');
    expect(repair?.label).toBe('Needs Repair');
  });
});
