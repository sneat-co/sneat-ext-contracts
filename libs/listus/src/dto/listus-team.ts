import { IListGroup } from './list-group';
import { IListBrief } from './list';

export interface IListusSpaceDbo {
  /** Canonical persisted list briefs keyed by full `${type}!${shortID}` ID. */
  lists?: Record<string, IListBrief>;
  /** Legacy grouped projection retained for backwards compatibility. */
  listGroups?: IListGroup[];
}
