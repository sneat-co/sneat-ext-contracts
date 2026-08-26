// Proposer type for the profile-data-model contract schema.
//
// The proposer is the account-holder identity and the home of contact details.
// It always holds email + phone. Its name (firstName/lastName) is optional and
// authoritative only when the proposer does NOT drive: when the proposer drives,
// it carries a driverId linking to a drivers-list entry and its name is read from
// that Driver (not duplicated on the proposer).

import { IDriver } from './driver';

export interface IProposer {
  email: string;
  phone: string;
  firstName?: string;
  lastName?: string;
  /** Set when the proposer drives — references the IDriver.id in the drivers list. */
  driverId?: string;
}

/**
 * Resolve the proposer's name: read from the linked Driver when `driverId` is
 * set, otherwise from the proposer's own firstName/lastName.
 * Verifies REQ:proposer (driving- and non-driving-proposer ACs).
 */
export function resolveProposerName(
  proposer: IProposer,
  drivers: readonly IDriver[],
): { firstName: string; lastName: string } {
  if (proposer.driverId) {
    const driver = drivers.find((d) => d.id === proposer.driverId);
    if (!driver) {
      throw new Error(
        `resolveProposerName: driverId '${proposer.driverId}' not found in drivers list`,
      );
    }
    return { firstName: driver.firstName, lastName: driver.lastName };
  }
  return {
    firstName: proposer.firstName ?? '',
    lastName: proposer.lastName ?? '',
  };
}
