import { describe, expect, it } from 'vitest';
import {
  defaultHappeningParticipantRole,
  happeningParticipantRoles,
  happeningTagError,
  happeningTagLimits,
  type IHappeningContactRef,
  isKnownHappeningParticipantRole,
  normalizeHappeningTags,
} from './happening';

describe('happening participant roles', () => {
  it('pins the exact wire values', () => {
    // These are persisted as `rolesOfItem` map KEYS in every space's happening
    // documents and must stay identical to the Go
    // dbo4calendarius.HappeningParticipantWellKnownRoles. Renaming one is a
    // data migration, not a refactor.
    expect(happeningParticipantRoles).toEqual([
      'participant',
      'teacher',
      'organizer',
      'assistant',
    ]);
  });

  it('defaults to participant, the value every pre-existing caller produced', () => {
    expect(defaultHappeningParticipantRole).toBe('participant');
    expect(happeningParticipantRoles).toContain(
      defaultHappeningParticipantRole,
    );
  });

  it('recognises every known role', () => {
    for (const role of happeningParticipantRoles) {
      expect(isKnownHappeningParticipantRole(role)).toBe(true);
    }
  });

  it('rejects anything outside the closed vocabulary', () => {
    // Casing and padding are NOT normalised for roles (unlike tags): the
    // backend refuses them with a 400, so the client must too.
    for (const role of [
      '',
      'Teacher',
      'TEACHER',
      ' teacher',
      'student',
      'owner',
      'event-attendee',
    ]) {
      expect(isKnownHappeningParticipantRole(role)).toBe(false);
    }
  });
});

describe('IHappeningContactRef', () => {
  it('stays byte-identical to the pre-role body when no role is given', () => {
    // The whole no-migration claim rests on this: a ref built the old way must
    // serialise with no `role` key at all, not with `role: undefined` or an
    // explicit "participant" the server would have defaulted anyway.
    const ref: IHappeningContactRef = { id: 'contact1', spaceID: 'space1' };
    expect(JSON.stringify(ref)).toBe('{"id":"contact1","spaceID":"space1"}');
  });

  it('adds role as one extra key beside the existing ones', () => {
    const ref: IHappeningContactRef = { id: 'contact1', role: 'teacher' };
    expect(JSON.parse(JSON.stringify(ref))).toEqual({
      id: 'contact1',
      role: 'teacher',
    });
  });

  it('only carries roles the server will accept', () => {
    for (const role of happeningParticipantRoles) {
      const ref: IHappeningContactRef = { id: 'c', role };
      expect(isKnownHappeningParticipantRole(ref.role as string)).toBe(true);
    }
  });
});

describe('normalizeHappeningTags', () => {
  it('returns an empty list for nothing', () => {
    expect(normalizeHappeningTags(undefined)).toEqual([]);
    expect(normalizeHappeningTags([])).toEqual([]);
    expect(normalizeHappeningTags(['', '   '])).toEqual([]);
  });

  it('trims and lower-cases', () => {
    expect(normalizeHappeningTags([' Guitar ', 'PIANO'])).toEqual([
      'guitar',
      'piano',
    ]);
  });

  it('collapses case-insensitive duplicates keeping first appearance order', () => {
    expect(normalizeHappeningTags(['piano', 'Guitar', 'PIANO'])).toEqual([
      'piano',
      'guitar',
    ]);
  });

  it('keeps inner spaces and hyphens', () => {
    expect(normalizeHappeningTags(['Beginner-Group', 'adult class'])).toEqual([
      'beginner-group',
      'adult class',
    ]);
  });

  it('never emits a tag its own validator rejects', () => {
    // The property that makes normalise-then-validate safe: a normaliser that
    // could emit an invalid value would turn user input into a server error.
    const inputs = [
      [' Guitar ', 'GUITAR', '', '  ', 'Beginner-Group'],
      ['ГИТАРА', 'гитара'],
      ['a'],
    ];
    for (const input of inputs) {
      for (const tag of normalizeHappeningTags(input)) {
        expect(happeningTagError(tag)).toBeUndefined();
      }
    }
  });
});

describe('happeningTagError', () => {
  it('accepts an ordinary tag', () => {
    expect(happeningTagError('guitar')).toBeUndefined();
    expect(happeningTagError('beginner-group')).toBeUndefined();
  });

  it('rejects an empty tag', () => {
    expect(happeningTagError('')).toBeDefined();
  });

  it('counts characters, not bytes', () => {
    // 32 Cyrillic characters are 64 UTF-8 bytes and must be accepted; a
    // byte-based limit would silently give non-Latin schools half the length.
    expect(
      happeningTagError('я'.repeat(happeningTagLimits.maxRunes)),
    ).toBeUndefined();
    expect(
      happeningTagError('я'.repeat(happeningTagLimits.maxRunes + 1)),
    ).toBeDefined();
    expect(
      happeningTagError('a'.repeat(happeningTagLimits.maxRunes)),
    ).toBeUndefined();
    expect(
      happeningTagError('a'.repeat(happeningTagLimits.maxRunes + 1)),
    ).toBeDefined();
  });

  it('rejects control characters', () => {
    // A newline turns one tag into two in every rendering of it.
    expect(happeningTagError('guitar\npiano')).toBeDefined();
    expect(happeningTagError('guitar\tpiano')).toBeDefined();
  });

  it('rejects the C1 range too, because the server does', () => {
    // The server validates with Go's unicode.IsControl, which covers
    // U+0080\u2013U+009F. A client guard that stopped at U+007F let `gui<NEL>tar`
    // through and turned it into a 400 nobody could have predicted from the
    // published rules \u2014 which falsifies the only promise this mirror makes.
    for (const codePoint of [0x0080, 0x0085, 0x009f]) {
      const tag = `gui${String.fromCharCode(codePoint)}tar`;
      expect(happeningTagError(tag)).toBeDefined();
    }
    // U+00A0 (NBSP) is the first NON-control code point above that range and
    // must stay acceptable, so the widening does not over-reach.
    expect(happeningTagError('gui\u00a0tar')).toBeUndefined();
  });

  it('publishes the same limits the server enforces', () => {
    expect(happeningTagLimits.maxCount).toBe(10);
    expect(happeningTagLimits.maxRunes).toBe(32);
  });
});
