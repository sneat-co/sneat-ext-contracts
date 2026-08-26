import { IProfileBase, validateProfile } from './profile';
import { IDriver } from './driver';
import { IVehicleRef } from './vehicle';
import { IAddress } from './address';

const address: IAddress = { line1: '1 Main St', town: 'Dublin', county: 'Dublin' };
const vehicle: IVehicleRef = {
  assetId: 'asset-1',
  ownership: 'owned',
  displayCache: { reg: '12-D-1', make: 'Toyota', model: 'Yaris', year: 2019 },
};

const driver = (id: string, isPrimary: boolean): IDriver => ({
  id,
  firstName: 'A',
  lastName: 'B',
  dateOfBirth: '1990-01-01',
  licence: { type: 'full', number: 'X', yearsHeld: 3, countryOfIssue: 'IE' },
  occupation: { title: 'Engineer', employmentStatus: 'employed' },
  ncbYears: 2,
  isPrimary,
});

const profile = (id: string, drivers: IDriver[]): IProfileBase => ({
  id,
  proposer: { email: 'a@b.com', phone: '+353870000000', driverId: drivers[0].id },
  drivers,
  vehicle,
  address,
  coverStartDate: '2026-07-01',
});

describe('validateProfile', () => {
  it('passes when exactly one driver is primary (verifies REQ:driver-type)', () => {
    expect(() => validateProfile(profile('p1', [driver('d1', true), driver('d2', false)]))).not.toThrow();
  });

  it('throws when zero drivers are primary', () => {
    expect(() => validateProfile(profile('p1', [driver('d1', false), driver('d2', false)]))).toThrow();
  });

  it('throws when more than one driver is primary', () => {
    expect(() => validateProfile(profile('p1', [driver('d1', true), driver('d2', true)]))).toThrow();
  });
});

describe('multiple profiles per Space', () => {
  it('two profiles coexist, each addressable by its own id (verifies REQ:profile-aggregate)', () => {
    const profiles = [profile('own-car', [driver('d1', true)]), profile('prospective', [driver('d2', true)])];
    expect(profiles.map((p) => p.id)).toEqual(['own-car', 'prospective']);
    expect(new Set(profiles.map((p) => p.id)).size).toBe(2);
  });
});
