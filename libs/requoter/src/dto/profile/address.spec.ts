import { IAddress, resolveDriverAddress } from './address';

describe('resolveDriverAddress', () => {
  const profileAddress: IAddress = {
    eircode: 'D01AB12',
    line1: '1 Main St',
    town: 'Dublin',
    county: 'Dublin',
  };
  const override: IAddress = {
    line1: '99 Other Rd',
    town: 'Cork',
    county: 'Cork',
  };

  it('returns the override when the driver has one (verifies REQ:address-type)', () => {
    expect(resolveDriverAddress(profileAddress, override)).toBe(override);
  });

  it('falls back to the profile address when there is no override', () => {
    expect(resolveDriverAddress(profileAddress, undefined)).toBe(profileAddress);
  });
});
