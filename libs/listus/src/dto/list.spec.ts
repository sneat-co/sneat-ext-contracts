import { IListItemBrief, getListShortUrlId } from './list';
import { IMediaRef } from '@sneat/extension-media-contract';

describe('getListShortUrlId', () => {
  it('combines spaceId and shortId when a shortId is given', () => {
    expect(getListShortUrlId('fam1', 'groceries')).toBe('fam1-groceries');
  });

  it('falls back to the full id when there is no shortId', () => {
    expect(getListShortUrlId('fam1', undefined, 'to-do:abc')).toBe('to-do:abc');
  });

  it('returns undefined when neither shortId nor id is provided', () => {
    expect(getListShortUrlId('fam1')).toBeUndefined();
  });
});

describe('IListItemBrief photo forward reference', () => {
  it('uses the shared media contract rather than an extension-specific URL', () => {
    const photo: IMediaRef = { mediaID: 'media-oatmeal' };
    const item: IListItemBrief = { id: 'oatmeal', title: 'Oatmeal', photo };

    expect(item.photo).toEqual(photo);
  });
});
