// Proves the legacy legacy mod-assetus-core companion symbols (engine/fuel
// consts + enums, possession list, currency/unit lists, AssetRealEstateType,
// the generic IAssetDboBase) ported into dto/asset.ts preserve the exact
// surface the migrated components rely on.
import { describe, expect, it } from 'vitest';
import {
  AssetPossessionLeasing,
  AssetPossessionOwning,
  AssetPossessionRenting,
  AssetPossessionUndisclosed,
  AssetPossessions,
  CurrencyEUR,
  CurrencyList,
  CurrencyUSD,
  EngineTypeCombustion,
  EngineTypeElectric,
  EngineTypeHybrid,
  EngineTypeOther,
  EngineTypePHEV,
  EngineTypeSteam,
  EngineTypeUnknown,
  EngineTypes,
  FuelTypeDiesel,
  FuelTypeElectricity,
  FuelTypeHydrogen,
  FuelTypeOther,
  FuelTypePetrol,
  FuelTypeUnknown,
  FuelTypes,
  FuelVolumeUnitTypes,
  MileageUnitTypes,
  type AssetRealEstateType,
  type CurrencyCode,
  type EngineType,
  type FuelType,
  type IAssetDboBase,
  type IAssetExtra,
} from './asset';
import { type IAssetVehicleExtra } from './extras';

describe('legacy engine/fuel companions', () => {
  it('engine consts match the EngineType union', () => {
    const engineTypes: EngineType[] = [
      EngineTypeUnknown,
      EngineTypeOther,
      EngineTypePHEV,
      EngineTypeCombustion,
      EngineTypeElectric,
      EngineTypeHybrid,
      EngineTypeSteam,
    ];
    expect(engineTypes).toEqual([
      '',
      'other',
      'phev',
      'combustion',
      'electric',
      'hybrid',
      'steam',
    ]);
  });

  it('EngineTypes enum members carry the legacy string values', () => {
    expect(EngineTypes.unknown).toBe('');
    expect(EngineTypes.other).toBe('other');
    expect(EngineTypes.phev).toBe('phev');
    expect(EngineTypes.combustion).toBe('combustion');
    expect(EngineTypes.electric).toBe('electric');
    expect(EngineTypes.hybrid).toBe('hybrid');
    expect(EngineTypes.steam).toBe('steam');
  });

  it('fuel consts match the FuelType union (electricity stays an extra)', () => {
    const fuelTypes: FuelType[] = [
      FuelTypeUnknown,
      FuelTypeOther,
      FuelTypePetrol,
      FuelTypeDiesel,
      FuelTypeHydrogen,
    ];
    expect(fuelTypes).toEqual(['', 'other', 'petrol', 'diesel', 'hydrogen']);
    expect(FuelTypeElectricity).toBe('electricity');
  });

  it('FuelTypes enum members carry the legacy string values', () => {
    expect(FuelTypes.unknown).toBe('');
    expect(FuelTypes.other).toBe('other');
    expect(FuelTypes.petrol).toBe('petrol');
    expect(FuelTypes.diesel).toBe('diesel');
    expect(FuelTypes.hydrogen).toBe('hydrogen');
    expect(FuelTypes.electricity).toBe('electricity');
  });
});

describe('legacy possession companions', () => {
  it('AssetPossessions is the ordered option list', () => {
    expect(AssetPossessions).toEqual([
      AssetPossessionOwning,
      AssetPossessionRenting,
      AssetPossessionLeasing,
      AssetPossessionUndisclosed,
    ]);
    expect(AssetPossessions).toEqual([
      'owning',
      'renting',
      'leasing',
      'undisclosed',
    ]);
  });
});

describe('legacy currency + unit lists', () => {
  it('CurrencyList carries USD then EUR', () => {
    const codes: CurrencyCode[] = CurrencyList;
    expect(codes).toEqual([CurrencyUSD, CurrencyEUR]);
    expect(codes).toEqual(['USD', 'EUR']);
  });

  it('fuel-volume and mileage unit lists match the legacy values', () => {
    expect(FuelVolumeUnitTypes).toEqual(['l', 'g']);
    expect(MileageUnitTypes).toEqual(['km', 'mile']);
  });
});

describe('AssetRealEstateType', () => {
  it('admits the legacy dwelling subtypes', () => {
    const types: AssetRealEstateType[] = ['house', 'apartment', 'land'];
    expect(types).toHaveLength(3);
  });
});

describe('IAssetDboBase', () => {
  it('carries the legacy generic extra + audit + title fields', () => {
    const dbo: IAssetDboBase<'vehicle', IAssetVehicleExtra> = {
      id: 'a1',
      category: 'vehicles',
      status: 'draft',
      title: 'My car',
      extraType: 'vehicle',
      extra: { make: 'Toyota', model: 'Corolla', regNumber: '12-D-345' },
      createdAt: '2026-01-01T00:00:00.000Z',
      createdBy: '-',
      updatedAt: '2026-01-01T00:00:00.000Z',
      updatedBy: '-',
      possession: 'owning',
      type: 'car',
    };
    expect(dbo.extra?.make).toBe('Toyota');
    expect(dbo.extraType).toBe('vehicle');
    expect(dbo.title).toBe('My car');
    expect(dbo.createdBy).toBe('-');
    // name/condition/visibility are optional on the in-progress base.
    expect(dbo.name).toBeUndefined();
  });

  it('accepts the open IAssetExtra map by default', () => {
    const dbo: IAssetDboBase = {
      id: 'a2',
      category: 'dwelling',
      status: 'draft',
    };
    const extra: IAssetExtra = { anything: 1 };
    dbo.extra = extra;
    expect(dbo.extra['anything']).toBe(1);
  });
});
