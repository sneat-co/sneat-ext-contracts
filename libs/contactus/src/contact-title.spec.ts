import { describe, expect, it } from 'vitest';
import { IContactBrief, IContactDbo } from './dto';
import { contactPersonName, contactTitle } from './contact-title';

function contact(
  brief?: Partial<IContactBrief>,
  dbo?: Partial<IContactDbo>,
) {
  return {
    id: 'contact-1',
    ...(brief ? { brief: { type: 'person' as const, ...brief } } : {}),
    ...(dbo ? { dbo: { type: 'person' as const, ...dbo } as IContactDbo } : {}),
  };
}

describe('contactTitle', () => {
  it('renders the first-name-only brief returned for a newly invited member', () => {
    expect(
      contactTitle(contact({ names: { firstName: 'Jamie' } })),
    ).toBe('Jamie');
  });

  it('preserves display-name precedence', () => {
    expect(
      contactTitle(
        contact(
          { title: 'Brief title', names: { fullName: 'Full name' } },
          { title: 'DBO title' },
        ),
        { preferredTitle: 'Preferred title' },
      ),
    ).toBe('Preferred title');
    expect(
      contactTitle(
        contact({ title: 'Brief title' }, { title: 'DBO title' }),
      ),
    ).toBe('Brief title');
    expect(contactTitle(contact(undefined, { title: 'DBO title' }))).toBe(
      'DBO title',
    );
  });

  it('includes brief short titles only when the caller requests them', () => {
    const value = contact({
      shortTitle: 'Short title',
      names: { firstName: 'Jamie' },
    });
    expect(contactTitle(value)).toBe('Jamie');
    expect(contactTitle(value, { includeBriefShortTitle: true })).toBe(
      'Short title',
    );
  });

  it('preserves caller-specific ID and text fallbacks', () => {
    expect(contactTitle(contact())).toBe('contact-1');
    expect(
      contactTitle(contact(), {
        fallbackToID: false,
        fallback: 'Unnamed contact',
      }),
    ).toBe('Unnamed contact');
    expect(contactTitle(undefined, { fallback: 'No contact' })).toBe(
      'No contact',
    );
  });
});

describe('contactPersonName', () => {
  it('matches the established display-name precedence', () => {
    expect(
      contactPersonName({
        fullName: 'Jamie Rivera',
        firstName: 'Jamie',
        nickName: 'Jay',
      }),
    ).toBe('Jamie Rivera');
    expect(
      contactPersonName({ firstName: 'Jamie', lastName: 'Rivera' }),
    ).toBe('Jamie Rivera');
    expect(
      contactPersonName({ firstName: 'Jamie', nickName: 'Jay' }),
    ).toBe('Jay');
    expect(contactPersonName({ middleName: 'Quinn' })).toBe('Quinn');
  });

  it('preserves the member-list first-plus-last then full-name behavior', () => {
    expect(
      contactPersonName(
        {
          fullName: 'Jamie Q. Rivera',
          firstName: 'Jamie',
          lastName: 'Rivera',
        },
        'member',
      ),
    ).toBe('Jamie Rivera');
    expect(contactPersonName({ fullName: 'Jamie Rivera' }, 'member')).toBe(
      'Jamie Rivera',
    );
  });
});
