// Vehicle reference for the profile-data-model contract schema.
//
// The profile references the car by Assetus `assetId` plus a small display cache
// and an ownership marker — it does NOT embed a forked copy of the asset record.
// Insurance-usage fields live here on the profile side, not in Assetus.

export type VehicleOwnership = 'owned' | 'prospective';

export type ClassOfUse = 'social' | 'social-commuting' | 'business';

export type OvernightParking =
  | 'drive'
  | 'garage'
  | 'street'
  | 'car-park'
  | 'public-car-park';

export interface IVehicleDisplayCache {
  reg: string;
  make: string;
  model: string;
  year: number;
}

export interface IVehicleRef {
  assetId: string;
  ownership: VehicleOwnership;
  displayCache: IVehicleDisplayCache;
  annualMileage?: number;
  classOfUse?: ClassOfUse;
  overnightParking?: OvernightParking;
  modifications?: string[];
}
