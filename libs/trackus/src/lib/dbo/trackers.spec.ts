import {
  getStandardTrackerTitle,
  isStandardTracker,
  standardTrackers,
  TrackBy,
  TrackerValueType,
  NumberKind,
} from './i-tracker-dbo';

// The value/target vocabularies the AnyMeter standalone-product Feature pins
// (spec/features/anymeter/standalone-product). Every standard tracker MUST stay
// inside these sets, so a template can never be created with a value type the
// engine can't log or a trackBy target it can't resolve.
const VALUE_TYPES: readonly TrackerValueType[] = [
  'int',
  'float',
  'string',
  'bool',
  'money',
  'integers',
  'floats',
];
const TRACK_BY: readonly TrackBy[] = ['space', 'contact', 'asset'];
const NUMBER_KINDS: readonly NumberKind[] = ['absolute', 'cumulative'];

describe('standardTrackers catalogue', () => {
  it('defines at least one standard tracker', () => {
    expect(standardTrackers.length).toBeGreaterThan(0);
  });

  it('has unique ids', () => {
    const ids = standardTrackers.map((t) => t.id);
    expect(new Set(ids).size).toBe(ids.length);
  });

  it('every id follows the underscore convention and is recognised', () => {
    for (const t of standardTrackers) {
      expect(t.id.startsWith('_')).toBe(true);
      expect(isStandardTracker(t.id)).toBe(true);
    }
  });

  describe('every tracker brief is well-formed', () => {
    for (const t of standardTrackers) {
      describe(`${t.id}`, () => {
        const b = t.brief;

        it('has a non-empty title', () => {
          expect(typeof b.title).toBe('string');
          expect(b.title.trim().length).toBeGreaterThan(0);
        });

        it('has a value type within the engine set', () => {
          expect(VALUE_TYPES).toContain(b.valueType);
        });

        it('has at least one trackBy target, all valid', () => {
          expect(b.trackBy.length).toBeGreaterThan(0);
          for (const target of b.trackBy) {
            expect(TRACK_BY).toContain(target);
          }
        });

        it('has no duplicate trackBy targets', () => {
          expect(new Set(b.trackBy).size).toBe(b.trackBy.length);
        });

        it('has at least one category', () => {
          expect(b.categories.length).toBeGreaterThan(0);
          for (const c of b.categories) {
            expect(typeof c).toBe('string');
            expect(c.trim().length).toBeGreaterThan(0);
          }
        });

        it('has a valid numberKind when present', () => {
          if (b.numberKind !== undefined) {
            expect(NUMBER_KINDS).toContain(b.numberKind);
          }
        });
      });
    }
  });

  // Anchor a few well-known templates the landing + Feature reference by name,
  // so a rename or a wrong value type is caught.
  it.each([
    ['_weight', 'float', 'contact'],
    ['_mileage', 'int', 'asset'],
    ['_fuel', 'float', 'asset'],
    ['_electricity', 'int', 'asset'],
  ])('template %s is %s tracked by %s', (id, valueType, trackBy) => {
    const t = standardTrackers.find((x) => x.id === id);
    expect(t).toBeTruthy();
    expect(t?.brief.valueType).toBe(valueType);
    expect(t?.brief.trackBy).toContain(trackBy);
  });
});

describe('isStandardTracker', () => {
  it('is true only for underscore-prefixed ids', () => {
    expect(isStandardTracker('_weight')).toBe(true);
    expect(isStandardTracker('weight')).toBe(false);
    expect(isStandardTracker('')).toBe(false);
    expect(isStandardTracker('definitely-not-a-standard-tracker')).toBe(false);
  });
});

describe('getStandardTrackerTitle', () => {
  it('returns a non-empty title for every standard tracker', () => {
    for (const t of standardTrackers) {
      expect(getStandardTrackerTitle(t.id).trim().length).toBeGreaterThan(0);
    }
  });

  it('strips the leading underscore for unknown standard ids', () => {
    expect(getStandardTrackerTitle('_custom_thing')).toBe('custom_thing');
  });

  it('returns a non-standard id unchanged', () => {
    expect(getStandardTrackerTitle('user-made-id')).toBe('user-made-id');
  });

  it('maps the special-cased pull/push-up ids to friendly titles', () => {
    expect(getStandardTrackerTitle('_pull_ups')).toBe('Pull-ups');
    expect(getStandardTrackerTitle('_push_ups')).toBe('Push-ups');
  });
});
