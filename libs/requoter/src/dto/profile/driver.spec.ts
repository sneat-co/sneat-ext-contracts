import { IConviction, IDriver } from './driver';
import { IAddress } from './address';

const baseDriver: IDriver = {
  firstName: 'Jane',
  lastName: 'Doe',
  dateOfBirth: '1990-01-01',
  licence: { type: 'full', number: 'X123', yearsHeld: 5, countryOfIssue: 'IE' },
  occupation: { title: 'Engineer', employmentStatus: 'employed' },
  ncbYears: 4,
  isPrimary: true,
};

describe('IDriver', () => {
  it('carries NCB and convictions on the driver (verifies REQ:driver-type)', () => {
    const conviction: IConviction = { offence: 'SP30', date: '2023-05-01', penaltyPoints: 3 };
    const driver: IDriver = { ...baseDriver, ncbYears: 7, convictions: [conviction] };
    expect(driver.ncbYears).toBe(7);
    expect(driver.convictions).toEqual([conviction]);
  });

  it('references a contactus person by contactId when linked, with no embedded contact (verifies REQ:driver-type)', () => {
    const linked: IDriver = { ...baseDriver, contactId: 'contact-123' };
    expect(linked.contactId).toBe('contact-123');
    expect('contact' in linked).toBe(false);
  });

  it('leaves contactId unset when not yet linked to a contactus person', () => {
    expect(baseDriver.contactId).toBeUndefined();
  });

  it('supports an optional per-driver addressOverride', () => {
    const override: IAddress = { line1: '1 X St', town: 'Y', county: 'Z' };
    const driver: IDriver = { ...baseDriver, addressOverride: override };
    expect(driver.addressOverride).toBe(override);
  });
});
