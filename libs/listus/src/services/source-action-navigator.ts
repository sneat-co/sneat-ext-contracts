import { InjectionToken } from '@angular/core';
import type { IListItemBrief } from '../dto';
import type { ISpaceContext } from '@sneat/space-models';

export interface IListItemSourceActionNavigator {
  navigateToListItemSource(
    space: ISpaceContext,
    item: IListItemBrief,
  ): Promise<boolean>;
}

export const LIST_ITEM_SOURCE_ACTION_NAVIGATOR =
  new InjectionToken<IListItemSourceActionNavigator>(
    'ListItemSourceActionNavigator',
  );
