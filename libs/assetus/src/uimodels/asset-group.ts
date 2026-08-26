// AssetGroup uimodel — ported faithfully from legacy legacy mod-assetus-core
// (core/src/lib/uimodels/asset-group.ts). Wraps an IAssetGroupContext and exposes
// the group's id, totals and counts. Consumed by budgetus.

import { Totals } from '@sneat/space-models';
import { IAssetGroupContext } from '../contexts';
import { IAssetDtoGroupCounts } from '../dto';

export class AssetGroup {
  public readonly totals: Totals;

  constructor(public readonly context: IAssetGroupContext) {
    this.totals = new Totals(context.dbo?.totals);
  }

  get id(): string {
    return this.context.id;
  }

  public get numberOf(): IAssetDtoGroupCounts {
    return this.context.dbo?.numberOf || {};
  }
}
