// Polymorphic typed asset extras — mirror the Go backend
// (github.com/sneat-co/assetus/backend/extras4assetus + dbo4assetus).
//
// An asset's typed extra is resolved by its `extraType` discriminator
// (extra.WithExtraField on the backend). No field is dropped relative to the
// backend wire contract; an asset with no extra stays valid.

import { CountryId } from '@sneat/dto';
import {
  AssetDocumentType,
  AssetExtraType,
  EngineType,
  FuelType,
  FuelVolumeUnit,
  IMoney,
  MileageUnit,
} from './asset';

// --- extraType discriminator (backend extras4assetus.AssetExtraType*) ---

export type { AssetExtraType };

export const assetExtraTypeVehicle: AssetExtraType = 'vehicle';
export const assetExtraTypeDwelling: AssetExtraType = 'dwelling';
export const assetExtraTypeDocument: AssetExtraType = 'document';
export const assetExtraTypeAppliance: AssetExtraType = 'appliance';
export const assetExtraTypeElectronics: AssetExtraType = 'electronics';
export const assetExtraTypeUtility: AssetExtraType = 'utility';
export const assetExtraTypeValuable: AssetExtraType = 'valuable';
export const assetExtraTypeDigitalAsset: AssetExtraType = 'digital_asset';
export const assetExtraTypeSportsEquipment: AssetExtraType = 'sports_equipment';

// --- Value objects ---

// IAddress mirrors the backend sneat-go-core dbmodels.Address.
export interface IAddress {
  countryID: CountryId; // ISO 3166-1 alpha-2
  zipCode?: string;
  state?: string;
  city?: string;
  lines?: string;
}

// --- Vehicle extra (extras4assetus.AssetVehicleExtra) ---

// IEngine mirrors the backend extras4assetus.WithEngineData. Cubic-centimetre /
// kilowatt / Newton-metre figures use the backend json names engineCC/KW/NM.
export interface IEngine {
  engineType?: EngineType;
  engineFuel?: FuelType;
  engineCC?: number; // Engine volume in cubic centimetres
  engineKW?: number; // Engine power in kilowatts
  engineNM?: number; // Engine torque in Newton metres
  engineSerialNumber?: string;
}

// IWithMakeAndModel mirrors extras4assetus.WithMakeModelFields.
export interface IWithMakeAndModel {
  make?: string;
  model?: string;
}

// IAssetVehicleExtra mirrors extras4assetus.AssetVehicleExtra: make/model/
// regNumber/vin, the engine data (incl. engineSerialNumber) and the plain
// service/tax/inspection (NCT) due-dates.
export interface IAssetVehicleExtra extends IWithMakeAndModel, IEngine {
  regNumber?: string;
  vin?: string;
  nctExpires?: string; // ISO 'YYYY-MM-DD'
  taxExpires?: string; // ISO 'YYYY-MM-DD'
  nextServiceDue?: string; // ISO 'YYYY-MM-DD'
  warrantyExpiresOn?: string; // ISO 'YYYY-MM-DD'
  // lengthM/hullOrTailNumber are primarily used for boats/aircraft (regNumber/
  // vin alone don't capture size or the hull/tail identifier).
  lengthM?: number;
  hullOrTailNumber?: string;
}

// --- Appliance extra (extras4assetus.AssetApplianceExtra) ---

// IAssetApplianceExtra mirrors extras4assetus.AssetApplianceExtra — identity +
// the service/warranty dates that feed maintenance and renewals.
export interface IAssetApplianceExtra {
  make?: string;
  model?: string;
  serialNumber?: string;
  installedOn?: string; // ISO 'YYYY-MM-DD'
  warrantyExpiresOn?: string; // ISO
  nextServiceDue?: string; // ISO
  energyRating?: string;
}

// --- Electronics extra (extras4assetus.AssetElectronicsExtra) ---

// IAssetElectronicsExtra mirrors extras4assetus.AssetElectronicsExtra.
export interface IAssetElectronicsExtra {
  make?: string;
  model?: string;
  serialNumber?: string;
  imei?: string;
  storageCapacity?: string;
  warrantyExpiresOn?: string; // ISO 'YYYY-MM-DD'
}

// --- Utility extra (extras4assetus.AssetUtilityExtra) ---

// IAssetUtilityExtra mirrors extras4assetus.AssetUtilityExtra — a utility/
// service account/contract; renewalDate drives the (annual) renewal.
export interface IAssetUtilityExtra {
  provider?: string;
  accountNumber?: string;
  tariff?: string;
  monthlyCost?: IMoney;
  renewalDate?: string; // ISO 'YYYY-MM-DD'
  contractEndsOn?: string; // ISO
}

// --- Valuable extra (extras4assetus.AssetValuableExtra) ---

// IAssetValuableExtra mirrors extras4assetus.AssetValuableExtra. The kind
// (jewelry/watch/art/collectible/antique) is carried by the asset's `type`,
// not this extra.
export interface IAssetValuableExtra {
  material?: string;
  appraisedValue?: IMoney;
  appraisalDate?: string; // ISO 'YYYY-MM-DD'
  certificateNumber?: string;
  insuredValue?: IMoney;
}

// --- Digital asset extra (extras4assetus.AssetDigitalAssetExtra) ---

// IAssetDigitalAssetExtra mirrors extras4assetus.AssetDigitalAssetExtra — a
// domain name / software licence / other digital asset; renewalDate drives
// renewals & the ReQuoter switch prompt, like IAssetUtilityExtra.
export interface IAssetDigitalAssetExtra {
  registrar?: string;
  identifier?: string; // domain name / account handle / licence key
  renewalDate?: string; // ISO 'YYYY-MM-DD'
  autoRenew?: boolean;
  renewalCost?: IMoney;
}

// --- Sports-equipment extra (extras4assetus.AssetSportsEquipmentExtra) ---

// BicycleKind mirrors extras4assetus.BicycleKind — only meaningful when the
// asset's `type` is 'bicycle'.
export type BicycleKind =
  | 'road'
  | 'mountain'
  | 'e_bike'
  | 'kids'
  | 'hybrid'
  | 'gravel'
  | 'other';

export const bicycleKinds: readonly BicycleKind[] = [
  'road',
  'mountain',
  'e_bike',
  'kids',
  'hybrid',
  'gravel',
  'other',
];

// IAssetSportsEquipmentExtra mirrors extras4assetus.AssetSportsEquipmentExtra.
export interface IAssetSportsEquipmentExtra {
  make?: string;
  model?: string;
  frameNumber?: string;
  bicycleKind?: BicycleKind;
}

// --- Vehicle records (dbo4assetus.VehicleRecordDbo) ---

// IVehicleMileage mirrors dbo4assetus.VehicleMileage.
export interface IVehicleMileage {
  value: number;
  unit: MileageUnit;
}

// IVehicleFuelRecord mirrors dbo4assetus.VehicleFuelRecord — the persisted fuel
// reading carrying volume/unit/amount and the legacy fuelCost/currency.
export interface IVehicleFuelRecord {
  volume?: number;
  unit?: FuelVolumeUnit;
  amount?: IMoney;
  fuelCost?: number;
  currency?: string;
}

// IVehicleRecord mirrors the persisted dbo4assetus.VehicleRecordDbo (a mileage
// and/or fuel reading in the asset's vehicleRecords child collection).
export interface IVehicleRecord {
  fuel?: IVehicleFuelRecord;
  mileage?: IVehicleMileage;
  // with.CreatedFields (backend) — serialized alongside the record.
  createdAt?: string;
  createdBy?: string;
}

// IAddVehicleRecordRequest mirrors the flat append request
// dto4assetus.AddVehicleRecordRequest (fuelVolume/fuelVolumeUnit/fuelCost/
// currency/mileage/mileageUnit). Wired to the backend create_vehicle_record
// route by the asset service (FE Task 4).
export interface IAddVehicleRecordRequest {
  assetID: string;
  fuelVolume?: number;
  fuelVolumeUnit?: FuelVolumeUnit;
  fuelCost?: number;
  currency?: string;
  mileage?: number;
  mileageUnit?: MileageUnit;
}

// --- Dwelling extra (extras4assetus.AssetDwellingExtra) ---

// IAssetDwellingExtra mirrors extras4assetus.AssetDwellingExtra: address, rent
// price ({value,currency} under the backend `rent_price` json name), bedrooms
// and floor area in square metres.
export interface IAssetDwellingExtra {
  address?: IAddress;
  rent_price?: {
    value?: number;
    currency?: string;
  };
  numberOfBedrooms?: number;
  areaSqM?: number;
}

// --- Document extra (extras4assetus.AssetDocumentExtra + doc_type_schema.go) ---

// IAssetDocumentExtra mirrors extras4assetus.AssetDocumentExtra — the full
// 8-field shape. `regNumber` is the legacy alias for a document `number`.
export interface IAssetDocumentExtra {
  docType?: AssetDocumentType;
  number?: string;
  batchNumber?: string;
  countryID?: CountryId;
  issuedBy?: string;
  issuedOn?: string; // ISO 'YYYY-MM-DD'
  effectiveFrom?: string; // ISO 'YYYY-MM-DD'
  expiresOn?: string; // ISO 'YYYY-MM-DD'
  regNumber?: string; // legacy alias for `number`
  // The asset this document covers — the vehicle/dwelling an insurance policy
  // insures, or the appliance/electronics/vehicle a warranty covers. General
  // link; `insuredAssetID` is the legacy insurance-specific alias.
  coversAssetID?: string;
  insuredAssetID?: string;
}

// --- Per-doc-type validation schema (doc_type_schema.go / standardDocTypesByID) ---

// IDocTypeField is the per-field rule of a document-type validation schema
// (backend DocTypeField).
export interface IDocTypeField {
  type?: 'str' | 'int' | 'date';
  required?: boolean;
  exclude?: boolean;
  max?: number;
  min?: number;
}

// IDocTypeStandardFields is the set of standard fields a document type may
// constrain (backend DocTypeStandardFields).
export interface IDocTypeStandardFields {
  title?: IDocTypeField;
  number?: IDocTypeField;
  issuedBy?: IDocTypeField;
  issuedOn?: IDocTypeField;
  validTill?: IDocTypeField;
  members?: IDocTypeField;
}

// DocTypeDef is a document-type definition with its validation schema (backend
// DocTypeDef). The `id` is an AssetDocumentType, with `other` admitted for the
// generic fallback type.
export interface IDocTypeDef {
  id: AssetDocumentType | 'other';
  fields?: IDocTypeStandardFields;
}

// standardDocTypesByID is the per-doc-type validation schema mirroring the
// backend standardDocTypesByID. Passport and driving_license require
// number + validity (validTill); marriage_cert/birth_cert require
// number + issuedOn and exclude validity (marriage_cert allows up to 2 members).
export const standardDocTypesByID: Record<string, IDocTypeDef> = {
  other: {
    id: 'other',
    fields: {
      title: { required: true },
    },
  },
  passport: {
    id: 'passport',
    fields: {
      number: { required: true },
      validTill: { required: true },
      members: { max: 1 },
    },
  },
  driving_license: {
    id: 'driving_license',
    fields: {
      number: { required: true },
      validTill: { required: true },
      members: { max: 1 },
    },
  },
  // Insurance policy / contract — the Assetus side of the ReQuoter car re-quote:
  // policy number + insurer (issuedBy) + expiry (validTill), which drives the
  // renewal reminder and the re-quote prompt.
  insurance: {
    id: 'insurance',
    fields: {
      number: { required: true },
      issuedBy: {},
      validTill: { required: true },
      members: { max: 1 },
    },
  },
  // Home / contents insurance — covers a dwelling (coversAssetID). Yearly renewal.
  home_insurance: {
    id: 'home_insurance',
    fields: {
      number: { required: true },
      issuedBy: {},
      validTill: { required: true },
    },
  },
  // Life / critical-illness — covers a person; validity optional (whole-life).
  life_insurance: {
    id: 'life_insurance',
    fields: {
      number: { required: true },
      issuedBy: {},
    },
  },
  // Warranty — covers another asset (coversAssetID); the end date is required.
  warranty: {
    id: 'warranty',
    fields: {
      issuedBy: {},
      validTill: { required: true },
    },
  },
  birth_cert: {
    id: 'birth_cert',
    fields: {
      number: { required: true },
      issuedBy: {},
      issuedOn: { required: true },
      validTill: { exclude: true },
      members: { max: 1 },
    },
  },
  marriage_cert: {
    id: 'marriage_cert',
    fields: {
      number: { required: true },
      issuedBy: {},
      issuedOn: { required: true },
      validTill: { exclude: true },
      members: { max: 2 },
    },
  },
};

// docTypeSchema returns the validation schema for a document type, or undefined
// when the document type imposes no standard schema (mirrors the backend
// DocTypeSchema lookup).
export function docTypeSchema(
  docType: AssetDocumentType | 'other' | undefined,
): IDocTypeDef | undefined {
  if (!docType) {
    return undefined;
  }
  return standardDocTypesByID[docType];
}
