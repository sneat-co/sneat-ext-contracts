import { expectTypeOf } from 'vitest';
import type { Timestamp } from 'firebase/firestore';
import type {
  ITrackerDbo,
  NumberKind,
  TrackerValueType,
} from '../dbo/i-tracker-dbo';
import type {
  IAddTrackerPointRequest,
  IAddTrackerPointResponse,
} from './trackus-api.types';

// Compile these fixtures with tsconfig.spec.json as well as running Vitest.
// The negative assertions detect accidental widening of the public contracts.
const legacyPoint: IAddTrackerPointRequest = {
  spaceID: 's1',
  trackerID: '_push_ups',
  trackByKind: 'contact',
  trackByID: 'c1',
  i: 20,
};
const floatingPoint: IAddTrackerPointRequest = {
  ...legacyPoint,
  entryID: 'stable-id',
  ts: '2026-11-01T08:30:00Z',
  n: 'Morning',
  f: 72.5,
};
const legacyResponse: IAddTrackerPointResponse = { entryID: 'stable-id' };
const legacyRecord: ITrackerDbo = {
  trackBy: ['contact'],
  valueType: 'int',
  categories: ['fitness'],
  title: 'Reps',
  entries: {},
};
const sharedRecord: ITrackerDbo = {
  ...legacyRecord,
  sharedTo: { 'contact:c1': { s2: { byUserID: 'u1', permissions: ['view'] } } },
  sharedFrom: {
    s3: { 'contact:c2': { byUserID: 'u2', permissions: ['view'] } },
  },
};

it('preserves legacy requests and required normalized entryID responses', () => {
  expectTypeOf(legacyPoint).toMatchTypeOf<IAddTrackerPointRequest>();
  expectTypeOf(floatingPoint).toMatchTypeOf<IAddTrackerPointRequest>();
  expectTypeOf(legacyResponse.entryID).toEqualTypeOf<string>();
  // @ts-expect-error Raw backend responses must be normalized by the runtime.
  const rawResponse: IAddTrackerPointResponse = { ID: 'stable-id' };
  void rawResponse;
});

it('accepts canonical and deprecated value vocabularies', () => {
  expectTypeOf<'boolean' | 'bool'>().toExtend<TrackerValueType>();
  expectTypeOf<'absolute' | 'summable' | 'cumulative'>().toExtend<NumberKind>();
});

it('separates ISO HTTP timestamps from Firebase read timestamps', () => {
  expectTypeOf<IAddTrackerPointRequest['ts']>().toEqualTypeOf<
    string | undefined
  >();
  const invalidTimestamp: IAddTrackerPointRequest = {
    ...legacyPoint,
    // @ts-expect-error HTTP requests use ISO strings, never Firebase Timestamps.
    ts: {} as Timestamp,
  };
  void invalidTimestamp;
});

it('exposes readonly backend sharing metadata without requiring it on legacy records', () => {
  expectTypeOf(legacyRecord).toMatchTypeOf<ITrackerDbo>();
  expectTypeOf(sharedRecord).toMatchTypeOf<ITrackerDbo>();
  const invalidShare: ITrackerDbo = {
    ...legacyRecord,
    // @ts-expect-error Backend shares are nested maps containing ShareDetail, not booleans.
    sharedTo: { 'contact:c1': { s2: true } },
  };
  const invalidPermissions: ITrackerDbo = {
    ...legacyRecord,
    sharedFrom: {
      s2: {
        'contact:c1': {
          byUserID: 'u1',
          // @ts-expect-error Permissions are a readonly string array.
          permissions: 'view',
        },
      },
    },
  };
  void invalidShare;
  void invalidPermissions;
});
