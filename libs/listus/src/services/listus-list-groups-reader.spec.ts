import {
  IListusListGroupsReader,
  LISTUS_LIST_GROUPS_READER,
} from './listus-list-groups-reader';

describe('LISTUS_LIST_GROUPS_READER', () => {
  it('is an optional capability token, separate from IListusService', () => {
    const reader = LISTUS_LIST_GROUPS_READER;
    expect(reader.toString()).toContain('ListusListGroupsReader');

    const implementation: IListusListGroupsReader = {
      watchListGroups: () => {
        throw new Error('test-only implementation');
      },
    };
    expect(implementation.watchListGroups).toBeTypeOf('function');
  });
});
