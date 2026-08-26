// Asset domain model — mirrors the Go backend
// (github.com/sneat-co/assetus/backend). Assetus is OWNERSHIP-ONLY: there is no
// sharing/availability/borrow/lend field and no ext.yardius.

import { CountryId, ITitledRecord, IWithSpaceIDs } from '@sneat/dto';

// --- Enums (string values must match the backend exactly) ---

// AssetCategory mirrors backend const4assetus.Category exactly: the MVP set
// unioned with the ported legacy categories (legacy sport_gear→sports_equipment,
// vehicle→vehicles, misc→other already reconciled; dwelling/document/debt
// retained).
export type AssetCategory =
  | 'books'
  | 'games'
  | 'toys'
  | 'sports_equipment'
  | 'tools'
  | 'electronics'
  | 'appliances'
  | 'clothing'
  | 'vehicles'
  | 'camping_equipment'
  | 'other'
  | 'dwelling'
  | 'document'
  | 'debt'
  | 'utility'
  | 'valuables'
  | 'digital_asset';

export type AssetCondition =
  | 'new'
  | 'excellent'
  | 'good'
  | 'fair'
  | 'needs_repair'
  | 'broken';

// Ownership lifecycle — read-only-ish (the backend drives transitions).
// Mirrors backend const4assetus.Status: the MVP set unioned with the ported
// legacy 'draft' pre-active state.
export type AssetStatus =
  | 'draft'
  | 'active'
  | 'transferred'
  | 'archived'
  | 'disposed'
  | 'lost';

// AssetPossession mirrors backend const4assetus.Possession.
export type AssetPossession =
  | 'unknown'
  | 'undisclosed'
  | 'owning'
  | 'leasing'
  | 'renting';

// Legacy possession companions — ported from legacy mod-assetus-core
// (core/src/lib/dto/assetus-types.ts). AssetPossessions is the ordered option
// list the possession-card binds (owning/renting/leasing/undisclosed).
export const AssetPossessionUndisclosed: AssetPossession = 'undisclosed';
export const AssetPossessionOwning: AssetPossession = 'owning';
export const AssetPossessionRenting: AssetPossession = 'renting';
export const AssetPossessionLeasing: AssetPossession = 'leasing';
export const AssetPossessions: AssetPossession[] = [
  AssetPossessionOwning,
  AssetPossessionRenting,
  AssetPossessionLeasing,
  AssetPossessionUndisclosed,
];

// AssetType is the optional per-category subtype (backend const4assetus.Type).
// It is an open string so any category's subtype set is admitted; the named
// per-category unions below document the backend sets.
export type AssetVehicleType =
  | 'aircraft'
  | 'boat'
  | 'bus'
  | 'car'
  | 'helicopter'
  | 'motorcycle'
  | 'truck'
  | 'van';
export type AssetDwellingType =
  | 'apartment'
  | 'house'
  | 'office'
  | 'shop'
  | 'land'
  | 'garage'
  | 'warehouse';
// AssetRealEstateType is the legacy legacy mod-assetus-core name for the
// dwelling subtype set the add-dwelling component binds; a subset of
// AssetDwellingType, kept under its original name for the migrated component.
export type AssetRealEstateType = 'house' | 'apartment' | 'land';
export type AssetSportsEquipmentType =
  | 'bicycle'
  | 'kite'
  | 'kite_bar'
  | 'kite_board'
  | 'kite_hydrofoil'
  | 'prone_hydrofoil'
  | 'surf_board'
  | 'wetsuit'
  | 'wing'
  | 'wing_board'
  | 'wing_hydrofoil';
export type AssetDocumentType =
  | 'passport'
  | 'id_card'
  | 'driving_license'
  | 'insurance'
  | 'home_insurance'
  | 'life_insurance'
  | 'warranty'
  | 'marriage_cert'
  | 'birth_cert';
export type AssetElectronicsType =
  | 'laptop'
  | 'desktop'
  | 'phone'
  | 'tablet'
  | 'camera'
  | 'tv'
  | 'monitor'
  | 'console'
  | 'audio'
  | 'wearable'
  | 'printer'
  | 'networking'
  | 'other_electronics';
export type AssetApplianceType =
  | 'boiler'
  | 'washing_machine'
  | 'dishwasher'
  | 'fridge'
  | 'freezer'
  | 'oven'
  | 'cooker'
  | 'hvac'
  | 'water_heater'
  | 'dryer'
  | 'microwave'
  | 'vacuum'
  | 'other_appliance';
export type AssetUtilityType =
  | 'electricity'
  | 'gas'
  | 'water'
  | 'broadband'
  | 'mobile'
  | 'tv'
  | 'streaming'
  | 'waste'
  | 'other_utility';
export type AssetValuableType =
  | 'jewelry'
  | 'watch'
  | 'art'
  | 'collectible'
  | 'antique'
  | 'other_valuable';
export type AssetDigitalAssetType =
  | 'domain'
  | 'software_license'
  | 'digital_collectible'
  | 'other_digital';
export type AssetType =
  | AssetVehicleType
  | AssetDwellingType
  | AssetSportsEquipmentType
  | AssetDocumentType
  | AssetElectronicsType
  | AssetApplianceType
  | AssetUtilityType
  | AssetValuableType
  | AssetDigitalAssetType
  | string;

// EngineType mirrors backend const4assetus.EngineType ('' = unknown).
export type EngineType =
  | ''
  | 'other'
  | 'combustion'
  | 'electric'
  | 'phev'
  | 'hybrid'
  | 'steam';

// FuelType mirrors backend const4assetus.FuelType ('' = unknown).
export type FuelType =
  | ''
  | 'other'
  | 'bio'
  | 'petrol'
  | 'diesel'
  | 'hydrogen';

// Legacy engine/fuel companions — ported faithfully from
// legacy mod-assetus-core (core/src/lib/dto/assetus-types.ts) so the migrated
// vehicle components keep their exact const/enum-member references. Values match
// the EngineType/FuelType unions above.
export const EngineTypeUnknown = '';
export const EngineTypeOther = 'other';
export const EngineTypePHEV = 'phev';
export const EngineTypeCombustion = 'combustion';
export const EngineTypeElectric = 'electric';
export const EngineTypeHybrid = 'hybrid';
export const EngineTypeSteam = 'steam';

export enum EngineTypes {
  unknown = '',
  other = 'other',
  phev = 'phev',
  combustion = 'combustion',
  electric = 'electric',
  hybrid = 'hybrid',
  steam = 'steam',
}

export const FuelTypeUnknown = '';
export const FuelTypeOther = 'other';
export const FuelTypePetrol = 'petrol';
export const FuelTypeDiesel = 'diesel';
export const FuelTypeHydrogen = 'hydrogen';
export const FuelTypeElectricity = 'electricity';

export enum FuelTypes {
  unknown = '',
  other = 'other',
  petrol = 'petrol',
  diesel = 'diesel',
  hydrogen = 'hydrogen',
  electricity = 'electricity',
}

// Fuel-volume unit, mileage unit and currency are OPEN strings on the backend
// (no typed enum). The named unions below are kept for ergonomics only and MUST
// serialize as plain strings.
export type FuelVolumeUnit = 'l' | 'g' | string;
export type MileageUnit = 'km' | 'mile' | string;

// Legacy unit lists — ported from legacy mod-assetus-core
// (core/src/lib/dto/assetus-types.ts). The migrated mileage dialog binds these
// directly as its select-option source.
export const FuelVolumeUnitTypes: FuelVolumeUnit[] = ['l', 'g'];
export const MileageUnitTypes: MileageUnit[] = ['km', 'mile'];

// Legacy currency companions — ported from legacy mod-assetus-core
// (core/src/lib/types.ts). The mileage dialog binds CurrencyList as its
// currency-select source. IMoney.currency stays an open string; these are the
// curated picker options only.
export type CurrencyCode = 'USD' | 'EUR';
export const CurrencyUSD: CurrencyCode = 'USD';
export const CurrencyEUR: CurrencyCode = 'EUR';
export const CurrencyList: CurrencyCode[] = [CurrencyUSD, CurrencyEUR];

export type AssetVisibility =
  | 'private'
  | 'family'
  | 'friends'
  | 'friends_of_friends'
  | 'specific_space'
  | 'public';

// Derived from the owning space, read-only.
export type OwnerType =
  | 'individual'
  | 'family'
  | 'sports_club'
  | 'community'
  | 'school'
  | 'organisation';

export type HistoryEventType =
  | 'purchased'
  | 'repaired'
  | 'transferred'
  | 'sold'
  | 'donated'
  | 'lost';

// --- Value objects ---

// IMoney mirrors the backend MonetaryAmount / money.Amount ({currency,value}).
export interface IMoney {
  currency: string;
  value: number;
}

export interface IOwner {
  spaceID: string;
  spaceType: string;
  ownerType: OwnerType;
}

// IGeoPoint mirrors the backend GeoPoint ({lat,lng}).
export interface IGeoPoint {
  lat: number;
  lng: number;
}

// IAssetDates mirrors the backend embedded AssetDates — optional ISO
// 'YYYY-MM-DD' date strings.
export interface IAssetDates {
  dateOfBuild?: string;
  dateOfPurchase?: string;
  dateInsuredTill?: string;
  dateCertifiedTill?: string;
}

// --- Relationship sub-entities (mirror backend WithAssetRelationships) ---

// ISubAssetInfo mirrors the backend SubAssetInfo ({id,title,type,countryID,
// subType,expires}).
export interface ISubAssetInfo extends ITitledRecord {
  type: AssetCategory;
  countryID?: CountryId;
  subType?: string;
  expires?: string; // ISO 'YYYY-MM-DD'
}

// IAssetGroupCounts mirrors the backend AssetGroupCounts ({assets}).
export interface IAssetGroupCounts {
  assets?: number;
}

// IAssetGroupInfo mirrors the backend AssetGroupInfo
// ({id,title,order,desc,categoryID,numberOf,totals}).
export interface IAssetGroupInfo extends ITitledRecord {
  order?: number;
  desc?: string;
  categoryID?: AssetCategory;
  numberOf?: IAssetGroupCounts;
  totals?: IMoney[];
}

// --- Multi-space association (mirror backend WithAssetSpaces) ---

// IAssetusSpaceBrief mirrors the backend AssetusSpaceBrief ({assets}): the
// per-space projection of an asset's briefs (backend AssetBriefs:
// map[assetID]*AssetBrief).
export interface IAssetusSpaceBrief {
  assets?: Record<string, IAssetBrief>;
}

// IWithAssetSpaces mirrors the backend WithAssetSpaces ({spaces}): the
// multi-space association mapping spaceID -> per-space asset briefs, so a single
// asset record can be associated with multiple spaces. This is ADDITIVE to the
// single owning space carried via IWithSpaceIDs.spaceIDs.
export interface IWithAssetSpaces {
  spaces?: Record<string, IAssetusSpaceBrief>;
}

// --- Financial / liability sub-entities ---

// IAssetLiabilityInfo mirrors the backend AssetLiabilityInfo ({id,serviceTypes}).
export interface IAssetLiabilityInfo {
  id: string;
  serviceTypes?: string[];
}

// --- Asset ---

export interface IAssetBrief {
  id: string;
  name: string;
  category: AssetCategory;
  condition: AssetCondition;
  status: AssetStatus;
  visibility: AssetVisibility;
}

export interface IAssetDbo extends IAssetBrief, IWithSpaceIDs, IWithAssetSpaces {
  description?: string;
  acquisitionDate?: string; // ISO date
  purchasePrice?: IMoney;
  estimatedValue?: IMoney;
  location?: string;
  notes?: string;
  tags?: string[];
  photos?: string[];
  // Backend serializes these as ISO datetime strings.
  createdAt?: string;
  updatedAt?: string;

  // --- Ported legacy optional fields (backend AssetBase json names) ---
  isRequest?: boolean;
  countryID?: CountryId;
  type?: AssetType; // per-category subtype
  possession?: AssetPossession;
  parentCategoryID?: AssetCategory;
  yearOfBuild?: number;
  geo?: IGeoPoint;

  // Polymorphic typed extra (backend embeds extra.WithExtraField on the asset):
  // an `extraType` discriminator plus the typed `extra` map (vehicle/dwelling/
  // document). Optional — an asset with no extra (empty extraType) stays valid.
  extraType?: AssetExtraType;
  extra?: IAssetExtra;

  // Embedded AssetDates (flattened, like the backend).
  dateOfBuild?: string;
  dateOfPurchase?: string;
  dateInsuredTill?: string;
  dateCertifiedTill?: string;

  // --- Ported financial fields (backend AssetBase json names) ---
  totals?: IMoney[];
  canHaveIncome?: boolean;
  canHaveExpense?: boolean;
  financialDirection?: 'income' | 'expense';
  liabilities?: IAssetLiabilityInfo[];
  notUsedServiceTypes?: string[];

  // --- Ported relationship fields (backend WithAssetRelationships json names) ---
  groupID?: string;
  group?: IAssetGroupInfo;
  parentAssetID?: string;
  subAssets?: ISubAssetInfo[];
  sameAssetID?: string;
  // relatedAs mirrors the backend dbmodels.WithOptionalRelatedAs.RelatedAs — a
  // plain optional string naming the relationship role (json 'relatedAs').
  relatedAs?: string;
  memberIDs?: string[];
  membersInfo?: ITitledRecord[];
}

// AssetExtraType is the polymorphic extra discriminator
// (backend extras4assetus.AssetExtraType*). Defined here (the foundational dto
// module) so IAssetDboBase can reference it without importing dto/extras.ts
// (which itself imports from this module); dto/extras.ts re-exports it.
export type AssetExtraType =
  | 'vehicle'
  | 'dwelling'
  | 'document'
  | 'appliance'
  | 'electronics'
  | 'utility'
  | 'valuable'
  | 'digital_asset'
  | 'sports_equipment';

// IAssetExtra is the open per-category extra map (legacy
// legacy mod-assetus-core IAssetExtra). The typed vehicle/dwelling/document
// extras in dto/extras.ts are assignable to it.
export interface IAssetExtra {
  [key: string]: unknown;
}

// IAssetDboBase is the generic, in-progress asset DBO the add-asset components
// build before submit — ported faithfully from legacy legacy mod-assetus-core
// (core/src/lib/dto/dto-asset.ts IAssetDboBase). It extends the unified
// IAssetDbo but relaxes the briefs's required name/condition/visibility (the
// add flow sets `title`/`status`/`category` first) and carries the legacy
// polymorphic `extraType` + typed `extra`, the working `title`, and the
// created/updated audit fields. Nothing here is dropped relative to the legacy
// shape the migrated components relied on.
export interface IAssetDboBase<
  ExtraType extends AssetExtraType = AssetExtraType,
  // The typed vehicle/dwelling/document extras (dto/extras.ts) do not carry an
  // index signature, so the constraint is the structural `object` rather than
  // IAssetExtra (the open map). Either an open IAssetExtra or a typed extra fits.
  Extra extends object = IAssetExtra,
> extends Omit<
    IAssetDbo,
    'name' | 'condition' | 'visibility' | 'type' | 'extraType' | 'extra'
  > {
  name?: string;
  condition?: AssetCondition;
  visibility?: AssetVisibility;
  type?: AssetType;
  title?: string;
  extraType?: ExtraType;
  extra?: Extra;
  createdBy?: string;
  updatedBy?: string;
}

// --- History (append-only) ---

export interface IHistoryEvent {
  id: string;
  type: HistoryEventType;
  occurredAt: string; // ISO datetime
  actorRef: string;
  note?: string;
  fromOwner?: IOwner;
  toOwner?: IOwner;
}

// --- Select option helpers ---

export interface ILabeledOption<T extends string> {
  value: T;
  label: string;
}

const titleize = (s: string): string =>
  s
    .split('_')
    .map((w) => w.charAt(0).toUpperCase() + w.slice(1))
    .join(' ');

function options<T extends string>(values: readonly T[]): ILabeledOption<T>[] {
  return values.map((value) => ({ value, label: titleize(value) }));
}

export const assetCategories: readonly AssetCategory[] = [
  'books',
  'games',
  'toys',
  'sports_equipment',
  'tools',
  'electronics',
  'appliances',
  'clothing',
  'vehicles',
  'camping_equipment',
  'other',
  'dwelling',
  'document',
  'debt',
  'utility',
  'valuables',
  'digital_asset',
];

export const assetConditions: readonly AssetCondition[] = [
  'new',
  'excellent',
  'good',
  'fair',
  'needs_repair',
  'broken',
];

export const assetVisibilities: readonly AssetVisibility[] = [
  'private',
  'family',
  'friends',
  'friends_of_friends',
  'specific_space',
  'public',
];

export const categoryOptions = options(assetCategories);
export const conditionOptions = options(assetConditions);
export const visibilityOptions = options(assetVisibilities);

// --- Derivations matching the backend ---

// The default visibility a new asset inherits from its owning space type:
// private→private, family→family, everything else→specific_space.
export function defaultVisibilityForSpaceType(
  spaceType: string | undefined,
): AssetVisibility {
  switch (spaceType) {
    case 'private':
      return 'private';
    case 'family':
      return 'family';
    default:
      return 'specific_space';
  }
}

// Derives the OwnerType from the owning space type. Mirrors the backend's
// read-only derivation.
export function deriveOwnerType(spaceType: string | undefined): OwnerType {
  switch (spaceType) {
    case 'private':
      return 'individual';
    case 'family':
      return 'family';
    case 'sports_club':
      return 'sports_club';
    case 'community':
      return 'community';
    case 'school':
      return 'school';
    default:
      return 'organisation';
  }
}
