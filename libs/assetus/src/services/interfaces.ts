import { CountryId, ITitledRecord } from '@sneat/dto';
import {
  AssetCategory,
  AssetCondition,
  AssetExtraType,
  AssetPossession,
  AssetStatus,
  AssetType,
  AssetVisibility,
  HistoryEventType,
  IAddVehicleRecordRequest as IAddVehicleRecord,
  IAssetDbo,
  IAssetGroupInfo,
  IAssetLiabilityInfo,
  IGeoPoint,
  IHistoryEvent,
  IMoney,
  IOwner,
  ISubAssetInfo,
} from '../dto';

// Optional metadata accepted by both create and update.
export interface IAssetMetadata {
  description?: string;
  acquisitionDate?: string;
  purchasePrice?: IMoney;
  estimatedValue?: IMoney;
  location?: string;
  notes?: string;
  tags?: string[];
}

// IAssetRichFields is the optional unified (superset) set the backend
// create/update requests accept on top of the flat MVP metadata. Field names
// mirror the backend json names exactly (camelCase). Every field is optional so
// a flat-only request still type-checks. Reuses the already-ported unified DTO
// types from dto/asset.ts + dto/extras.ts — nothing is redefined here.
//
// The typed extra mirrors the backend extra.WithExtraField: an `extraType`
// discriminator plus an open `extra` object/map carrying the typed extra's
// fields (vehicle/dwelling/document).
export interface IAssetRichFields {
  type?: AssetType;
  possession?: AssetPossession;
  countryID?: CountryId;
  parentCategoryID?: AssetCategory;
  yearOfBuild?: number;
  isRequest?: boolean;
  geo?: IGeoPoint;

  // Embedded AssetDates (flattened, like the backend).
  dateOfBuild?: string;
  dateOfPurchase?: string;
  dateInsuredTill?: string;
  dateCertifiedTill?: string;

  // Custom fields (backend dbmodels.WithCustomFields json names).
  fieldsStr?: Record<string, string>;
  fieldsInt?: Record<string, number>;
  fieldsDate?: Record<string, string>;
  fieldsAmount?: Record<string, IMoney>;

  // Financial.
  totals?: IMoney[];
  canHaveIncome?: boolean;
  canHaveExpense?: boolean;
  financialDirection?: string;
  liabilities?: IAssetLiabilityInfo[];
  notUsedServiceTypes?: string[];

  // Relationships (backend WithAssetRelationships json names).
  groupID?: string;
  group?: IAssetGroupInfo;
  parentAssetID?: string;
  subAssets?: ISubAssetInfo[];
  sameAssetID?: string;
  relatedAs?: string;
  memberIDs?: string[];
  membersInfo?: ITitledRecord[];

  // Polymorphic typed extra (extra.WithExtraField): extraType + extra map.
  extraType?: AssetExtraType;
  extra?: Record<string, unknown>;
}

export interface ICreateAssetRequest extends IAssetMetadata, IAssetRichFields {
  spaceID: string;
  name: string;
  category: AssetCategory;
  condition: AssetCondition;
  // Optional override; when omitted the backend inherits the space default.
  visibility?: AssetVisibility;
  // Optional on create: when omitted defaults to active; may be supplied as
  // 'draft'. NOT editable via update.
  status?: AssetStatus;
}

export interface IAssetResponse {
  id: string;
  asset: IAssetDbo;
}

export interface IUpdateAssetRequest extends IAssetMetadata, IAssetRichFields {
  spaceID: string;
  assetID: string;
  name: string;
  category: AssetCategory;
  condition: AssetCondition;
  visibility: AssetVisibility;
  // NOTE: no `status` — ownership status is not editable via update.
}

export interface IRemoveAssetRequest {
  spaceID: string;
  assetID: string;
  // Soft-archive by default; hardDelete removes the record entirely.
  hardDelete?: boolean;
}

export interface ITransferAssetRequest {
  spaceID: string;
  assetID: string;
  toSpaceID: string;
}

export interface ITransferAssetResponse {
  id: string;
  owner: IOwner;
}

export interface IRecordHistoryEventRequest {
  spaceID: string;
  assetID: string;
  type: HistoryEventType;
  occurredAt?: string;
  note?: string;
}

export interface IAssetHistoryResponse {
  assetID: string;
  events: IHistoryEvent[];
}

// Appends a vehicle record (mileage and/or fuel reading) to a vehicle asset.
// The flat record shape (IAddVehicleRecord, from the dto) plus the spaceID
// required by the backend's embedded SpaceRequest.
export interface ICreateVehicleRecordRequest extends IAddVehicleRecord {
  spaceID: string;
}

export interface ICreateVehicleRecordResponse {
  id: string;
}
