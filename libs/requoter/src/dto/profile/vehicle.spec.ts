import { IVehicleRef } from './vehicle';

const vehicle: IVehicleRef = {
  assetId: 'asset-42',
  ownership: 'owned',
  displayCache: { reg: '12-D-3456', make: 'Toyota', model: 'Corolla', year: 2018 },
};

describe('IVehicleRef', () => {
  it('references the Assetus asset by id with a display cache + ownership marker, not a forked copy (verifies REQ:vehicle-reference)', () => {
    expect(vehicle.assetId).toBe('asset-42');
    expect(vehicle.ownership).toBe('owned');
    expect(vehicle.displayCache).toEqual({ reg: '12-D-3456', make: 'Toyota', model: 'Corolla', year: 2018 });
    // reference, not fork: no embedded full Assetus asset record
    expect('asset' in vehicle).toBe(false);
  });

  it('supports a prospective ownership marker and insurance-usage fields', () => {
    const prospective: IVehicleRef = {
      ...vehicle,
      ownership: 'prospective',
      annualMileage: 8000,
      classOfUse: 'social',
      overnightParking: 'garage',
      modifications: ['alloy wheels'],
    };
    expect(prospective.ownership).toBe('prospective');
    expect(prospective.annualMileage).toBe(8000);
    expect(prospective.modifications).toEqual(['alloy wheels']);
  });
});
