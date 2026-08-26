// Proves the ported typed asset-context wrappers preserve the legacy
// legacy mod-assetus-core surface: each carries its per-category typed extra on
// the dbo while keeping the base IAssetContext shape (incl. owner + brief).
import { describe, expect, it } from 'vitest';
import {
  type IAssetContext,
  type IAssetDocumentContext,
  type IAssetDwellingContext,
  type IAssetVehicleContext,
  type IAssetGroupContext,
} from './asset';

describe('typed asset contexts', () => {
  it('IAssetVehicleContext carries a vehicle extra', () => {
    const ctx: IAssetVehicleContext = {
      id: 'a1',
      space: { id: 's1' },
      brief: {
        id: 'a1',
        name: 'My car',
        category: 'vehicles',
        condition: 'good',
        status: 'active',
        visibility: 'private',
      },
      dbo: {
        id: 'a1',
        name: 'My car',
        category: 'vehicles',
        condition: 'good',
        status: 'active',
        visibility: 'private',
        extra: { make: 'Audi', model: 'A4', regNumber: '12-D-345' },
      },
    };
    expect(ctx.dbo?.extra?.make).toBe('Audi');
    expect(ctx.brief?.name).toBe('My car');
  });

  it('IAssetDwellingContext carries a dwelling extra', () => {
    const ctx: IAssetDwellingContext = {
      id: 'a2',
      space: { id: 's1' },
      dbo: {
        id: 'a2',
        name: 'Flat',
        category: 'dwelling',
        condition: 'good',
        status: 'active',
        visibility: 'private',
        extra: { numberOfBedrooms: 2, areaSqM: 75 },
      },
    };
    expect(ctx.dbo?.extra?.numberOfBedrooms).toBe(2);
  });

  it('IAssetDocumentContext carries a document extra', () => {
    const ctx: IAssetDocumentContext = {
      id: 'a3',
      space: { id: 's1' },
      dbo: {
        id: 'a3',
        name: 'Passport',
        category: 'document',
        condition: 'good',
        status: 'active',
        visibility: 'private',
        extra: { docType: 'passport', number: 'X123' },
      },
    };
    expect(ctx.dbo?.extra?.docType).toBe('passport');
  });

  it('base IAssetContext still carries owner', () => {
    const ctx: IAssetContext = {
      id: 'a4',
      space: { id: 's1' },
      owner: { spaceID: 's1', spaceType: 'family', ownerType: 'family' },
    };
    expect(ctx.owner?.ownerType).toBe('family');
  });

  it('IAssetGroupContext is a nav context over the group DTO', () => {
    const ctx: IAssetGroupContext = {
      id: 'g1',
      dbo: { id: 'g1', order: 0, title: 'Vehicles' },
    };
    expect(ctx.dbo?.title).toBe('Vehicles');
  });
});
