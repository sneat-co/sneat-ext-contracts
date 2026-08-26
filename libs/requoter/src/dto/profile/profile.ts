// Profile aggregate for the profile-data-model contract schema.
//
// A Space holds zero or more Profile aggregates, each identified by a stable id
// (e.g. own car vs a "what if I buy X" prospective quote). The DBO adds the
// space-scoping and created-metadata bases used elsewhere in this lib.

import { IWithCreated, IWithSpaceIDs } from '@sneat/dto';
import { IProposer } from './proposer';
import { IDriver } from './driver';
import { IVehicleRef } from './vehicle';
import { IAddress } from './address';

export interface IProfileBase {
  id: string;
  proposer: IProposer;
  drivers: IDriver[];
  vehicle: IVehicleRef;
  address: IAddress;
  /** Policy-level cover start date (ISO). The only cover field v1 commits to. */
  coverStartDate: string;
}

export interface IProfileDbo extends IProfileBase, IWithSpaceIDs, IWithCreated {}

/**
 * Validate the primary-driver invariant: exactly one driver MUST be primary.
 * Throws when zero or more than one driver has `isPrimary === true`.
 * Verifies REQ:driver-type (exactly-one-primary-driver).
 */
export function validateProfile(profile: Pick<IProfileBase, 'drivers'>): void {
  const primaries = profile.drivers.filter((d) => d.isPrimary).length;
  if (primaries !== 1) {
    throw new Error(
      `validateProfile: expected exactly one primary driver, found ${primaries}`,
    );
  }
}
