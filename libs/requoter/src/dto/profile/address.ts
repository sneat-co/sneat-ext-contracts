// Profile address types — the residence address on a Profile and the optional
// per-driver override. Part of the profile-data-model contract schema.

export interface IAddress {
  eircode?: string;
  line1: string;
  line2?: string;
  town: string;
  county: string;
}

/**
 * Resolve a driver's effective address: the per-driver override when present,
 * otherwise the profile's residence address.
 * Verifies REQ:address-type (driver-address-resolves-to-override-then-profile).
 */
export function resolveDriverAddress(
  profileAddress: IAddress,
  driverOverride?: IAddress,
): IAddress {
  return driverOverride ?? profileAddress;
}
