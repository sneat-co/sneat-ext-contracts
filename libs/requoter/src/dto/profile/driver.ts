// Driver type for the profile-data-model contract schema.
//
// A Driver references the contactus person it is (or becomes) via an optional
// `contactId` — mirroring how the vehicle references an Assetus `assetId` — and
// carries the insurance-specific fields (licence, occupation, NCB, convictions,
// isPrimary) profile-side. Personal identity is captured here and serves as the
// display copy when a contactId is linked; the profile never forks the full
// contactus record.

import { IAddress } from './address';

export type MaritalStatus =
  | 'single'
  | 'married'
  | 'civil-partnership'
  | 'separated'
  | 'divorced'
  | 'widowed';

export type LicenceType = 'full' | 'provisional' | 'international';

export type EmploymentStatus =
  | 'employed'
  | 'self-employed'
  | 'unemployed'
  | 'retired'
  | 'student'
  | 'homemaker';

export interface ILicence {
  type: LicenceType;
  number: string;
  yearsHeld: number;
  countryOfIssue: string;
}

export interface IOccupation {
  title: string;
  employmentStatus: EmploymentStatus;
}

export interface IConviction {
  offence: string;
  date: string; // ISO date
  penaltyPoints?: number;
  banMonths?: number;
}

export interface IDriver {
  /** Local identifier of this driver within the profile's drivers list; the
   * proposer's `driverId` references it when the proposer drives. */
  id?: string;
  /** Reference to the contactus person this driver is; set when linked or reused. */
  contactId?: string;
  title?: string;
  firstName: string;
  lastName: string;
  dateOfBirth: string; // ISO date
  maritalStatus?: MaritalStatus;
  licence: ILicence;
  occupation: IOccupation;
  /** Per-driver No-Claims-Bonus years. Lives on the driver, not the profile. */
  ncbYears: number;
  convictions?: IConviction[];
  addressOverride?: IAddress;
  isPrimary: boolean;
}
