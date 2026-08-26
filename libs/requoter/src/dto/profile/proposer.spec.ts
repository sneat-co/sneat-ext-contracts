import { IProposer, resolveProposerName } from './proposer';
import { IDriver } from './driver';

const drivingDriver: IDriver = {
  id: 'drv-1',
  firstName: 'Jane',
  lastName: 'Doe',
  dateOfBirth: '1990-01-01',
  licence: { type: 'full', number: 'X123', yearsHeld: 5, countryOfIssue: 'IE' },
  occupation: { title: 'Engineer', employmentStatus: 'employed' },
  ncbYears: 4,
  isPrimary: true,
};

describe('resolveProposerName', () => {
  it('reads the name from the linked driver when the proposer drives, without duplicating it (verifies REQ:proposer)', () => {
    const proposer: IProposer = { email: 'jane@example.com', phone: '+353871234567', driverId: 'drv-1' };
    expect(proposer.firstName).toBeUndefined();
    expect(proposer.lastName).toBeUndefined();
    expect(resolveProposerName(proposer, [drivingDriver])).toEqual({ firstName: 'Jane', lastName: 'Doe' });
  });

  it("uses the proposer's own name when they do not drive (verifies REQ:proposer)", () => {
    const proposer: IProposer = { email: 'pat@example.com', phone: '+353870000000', firstName: 'Pat', lastName: 'Murphy' };
    expect(proposer.driverId).toBeUndefined();
    expect(resolveProposerName(proposer, [drivingDriver])).toEqual({ firstName: 'Pat', lastName: 'Murphy' });
  });

  it('throws when driverId references a driver not in the list', () => {
    const proposer: IProposer = { email: 'x@example.com', phone: '+353870000001', driverId: 'missing' };
    expect(() => resolveProposerName(proposer, [drivingDriver])).toThrow();
  });
});
