import { listusHappeningFields } from './list-template';

describe('listusHappeningFields', () => {
  it('creates typed extension data and same-Space related references', () => {
    expect(
      listusHappeningFields({
        sourceListID: 'buy!regular',
        destinationListID: 'buy!groceries',
      }),
    ).toEqual({
      ext: {
        listus: {
          listTemplate: {
            sourceListID: 'buy!regular',
            destinationListID: 'buy!groceries',
          },
        },
      },
      related: {
        listus: {
          lists: {
            'buy!regular': {},
            'buy!groceries': {},
          },
        },
      },
    });
  });
});
