import { IListItemBrief, ListType } from './list';
import { IRelatedModules } from '@sneat/dto';

export interface IApplyListTemplateRequest {
  readonly spaceID: string;
  readonly sourceListID: string;
  readonly destinationListID: string;
  readonly requestID: string;
}

export interface IApplyListTemplateResult {
  readonly sourceListID: string;
  readonly destinationListID: string;
  readonly listType: Extract<ListType, 'buy' | 'do'>;
  readonly added: readonly IListItemBrief[];
  readonly restored: readonly IListItemBrief[];
  readonly unchanged: readonly IListItemBrief[];
  readonly disposition: 'applied' | 'reused';
}

/** Stored in a Calendarius happening's `ext.listus.listTemplate` payload. */
export interface IListTemplateHappeningLink {
  readonly sourceListID: string;
  readonly destinationListID: string;
}

export interface IListusHappeningExtension {
  readonly listTemplate?: IListTemplateHappeningLink;
}

export interface IListTemplateHappeningFields {
  readonly ext: Readonly<Record<'listus', IListusHappeningExtension>>;
  readonly related: IRelatedModules;
}

export const listusHappeningFields = (
  link: IListTemplateHappeningLink,
): IListTemplateHappeningFields => ({
  ext: { listus: { listTemplate: link } },
  related: {
    listus: {
      lists: {
        [link.sourceListID]: {},
        [link.destinationListID]: {},
      },
    },
  },
});
