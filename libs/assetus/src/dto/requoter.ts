// Typed response shapes for the Assetus → ReQuoter read surface. They mirror the
// Go backend DTOs (dto4assetus) exactly so ReQuoter/Formius can type the
// responses of the assetus read endpoints. See docs/requoter-integration.md.

import { AssetDocumentType, IAssetDbo } from './asset';

// IAssetListItem — an asset with its id (backend dto4assetus.AssetListItem).
export interface IAssetListItem {
  id: string;
  asset: IAssetDbo;
}

// IListAssetsResponse — GET /v0/assetus/assets?spaceID=.
export interface IListAssetsResponse {
  assets: IAssetListItem[];
}

// IReQuoteBundle — GET /v0/assetus/requote_bundle?spaceID=&vehicleID=. The
// one-call car re-quote bundle: the vehicle, the insurance policies covering it
// (latest expiry first), and the Space's driving licences.
export interface IReQuoteBundle {
  vehicle?: IAssetListItem;
  policies: IAssetListItem[];
  licences: IAssetListItem[];
}

// IInsuranceForAssetResponse — GET /v0/assetus/insurance_for_asset?spaceID=&assetID=.
export interface IInsuranceForAssetResponse {
  policies: IAssetListItem[];
}

// IExpiringDocument — a document expiring on/before the requested horizon.
export interface IExpiringDocument {
  id: string;
  docType?: AssetDocumentType;
  expiresOn: string; // ISO 'YYYY-MM-DD'
  asset: IAssetDbo;
}

// IExpiringDocumentsResponse — GET /v0/assetus/expiring_documents?spaceID=&withinDays=.
export interface IExpiringDocumentsResponse {
  documents: IExpiringDocument[];
}
