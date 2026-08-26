import { IVehicleRef } from './vehicle';

// The composed, read-only view of a requoter profile shown on the details page.
// It is assembled by re-reading the referenced canonical records (the Assetus car
// via its assetId, the Contactus drivers via their contactIds, the Assetus
// insurance document) and merging in the profile's own orphan cover/usage fields.
// Nothing here is persisted as-is — it is a live projection.

export interface IProfileViewDriver {
  readonly contactID: string;
  readonly name: string;
}

export interface IProfileViewInsurance {
  readonly docID?: string;
  readonly insurer?: string;
  readonly policyNumber?: string;
  readonly renewalDate?: string;
}

export interface IProfileView {
  readonly profileID: string;
  /** The car, as an IVehicleRef (assetId + ownership + display cache). */
  readonly vehicle?: IVehicleRef;
  /** Human label for the car (the Assetus asset title / entered name) — shown
   * when make/model are not set, e.g. a nickname-only onboarded car. */
  readonly vehicleName?: string;
  readonly drivers: readonly IProfileViewDriver[];
  readonly insurance?: IProfileViewInsurance;
  readonly coverStartDate?: string;
  readonly annualMileage?: number;
  readonly classOfUse?: string;
}
