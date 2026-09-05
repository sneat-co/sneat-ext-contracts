import { IListusSpaceDbo } from './listus-team';

describe('IListusSpaceDbo', () => {
  it('models the persisted lists map returned by the Listus module document', () => {
    const dbo: IListusSpaceDbo = {
      lists: {
        'do!cleaning': {
          type: 'do',
          title: 'Saturday cleaning',
          createdAt: { seconds: 1_788_609_600, nanoseconds: 0 },
          createdBy: 'user-1',
        },
      },
    };

    expect(dbo.lists?.['do!cleaning']?.title).toBe('Saturday cleaning');
  });
});
