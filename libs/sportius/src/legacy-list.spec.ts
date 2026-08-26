import { InjectionToken } from '@angular/core';
import {
  SPORTIUS_SERVICE,
  type ICreateListItemRequest,
  type ICreateListItemsRequest,
  type ICreateListRequest,
  type IDeleteListItemsRequest,
  type IListBrief,
  type IListContext,
  type IListDbo,
  type IListGroup,
  type IListInfo,
  type IListItemBrief,
  type IListItemIDsRequest,
  type IListItemResult,
  type IListItemsCommandParams,
  type IReorderListItemsRequest,
  type ISportiusService,
  type ISportiusSpaceDbo,
  type ISetListItemsIsComplete,
  type ListType,
} from './index';

describe('Sportius list compatibility exports', () => {
  it('retains every legacy symbol imported by the current Sportius UI', () => {
    const token: InjectionToken<ISportiusService> = SPORTIUS_SERVICE;
    expect(token.toString()).toContain('SportiusService');

    type CurrentConsumerExports = [
      ICreateListItemRequest,
      ICreateListItemsRequest,
      ICreateListRequest,
      IDeleteListItemsRequest,
      IListBrief,
      IListContext,
      IListDbo,
      IListGroup,
      IListInfo,
      IListItemBrief,
      IListItemIDsRequest,
      IListItemResult,
      IListItemsCommandParams,
      IReorderListItemsRequest,
      ISportiusService,
      ISportiusSpaceDbo,
      ISetListItemsIsComplete,
      ListType,
    ];

    const exportedTypes: CurrentConsumerExports | undefined = undefined;
    expect(exportedTypes).toBeUndefined();
  });
});
