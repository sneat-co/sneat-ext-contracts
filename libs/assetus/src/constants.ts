import { EnumAsUnionOfKeys } from '@sneat/core';

export const enum AssetPage {
  asset = 'asset',
}

export type AssetPages = EnumAsUnionOfKeys<typeof AssetPage>;
