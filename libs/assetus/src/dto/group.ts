// Asset-group DTO — ported faithfully from legacy legacy mod-assetus-core
// (core/src/lib/dto/dto-asset.ts: IAssetDtoGroup + IAssetDtoGroupCounts). The new
// lib already carries an IAssetGroupInfo (the backend AssetGroupInfo projection);
// IAssetDtoGroup is the legacy UI-facing group DTO the asset-group page binds to.
// Kept as a distinct named export so consumers migrate without losing fields.

import { ITitledRecord, ITotalsHolder, IWithSpaceIDs } from '@sneat/dto';
import { AssetCategory } from './asset';

export interface IAssetDtoGroupCounts {
  assets?: number;
}

export interface IAssetDtoGroup
  extends IWithSpaceIDs,
    ITitledRecord,
    ITotalsHolder {
  id: string;
  order: number;
  desc?: string;
  categoryId?: AssetCategory;
  numberOf?: IAssetDtoGroupCounts;
}
