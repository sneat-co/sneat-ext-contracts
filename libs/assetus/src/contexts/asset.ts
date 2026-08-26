import { INavContext } from '@sneat/core';
import { ISpaceItemNavContext } from '@sneat/space-models';
import {
  IAssetBrief,
  IAssetDbo,
  IAssetDocumentExtra,
  IAssetDwellingExtra,
  IAssetVehicleExtra,
  IAssetDtoGroup,
  IOwner,
} from '../dto';

export interface IAssetContext
  extends ISpaceItemNavContext<IAssetBrief, IAssetDbo> {
  owner?: IOwner;
}

// Typed asset-context wrappers — ported from legacy legacy mod-assetus-core
// (core/src/lib/contexts/asset-context.ts). Each carries the per-category typed
// `extra` on its dbo so consumers (e.g. docus) get a strongly-typed extra
// without losing the base IAssetContext surface (incl. `owner`).
export interface IAssetVehicleContext extends IAssetContext {
  readonly dbo?: IAssetDbo & { extra?: IAssetVehicleExtra };
}

export interface IAssetDwellingContext extends IAssetContext {
  readonly dbo?: IAssetDbo & { extra?: IAssetDwellingExtra };
}

export interface IAssetDocumentContext extends IAssetContext {
  readonly dbo?: IAssetDbo & { extra?: IAssetDocumentExtra };
}

// IAssetGroupContext mirrors the legacy group context (a plain INavContext over
// the group DTO). Used by the AssetGroup uimodel.
export type IAssetGroupContext = INavContext<IAssetDtoGroup, IAssetDtoGroup>;
