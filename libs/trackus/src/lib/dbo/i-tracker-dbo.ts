import type { Timestamp } from 'firebase/firestore';
import { IIdAndBrief, IIdAndOptionalBriefAndOptionalDbo } from '@sneat/core';
import { ISpaceContext } from '@sneat/space-models';

export type TrackerValueType =
  | 'int'
  | 'float'
  | 'string'
  | 'boolean'
  /** @deprecated Use the canonical backend value 'boolean'. */
  | 'bool'
  | 'money'
  | 'integers'
  | 'floats';

export type NumberKind =
  | 'absolute'
  | 'summable'
  /** @deprecated Use the canonical backend value 'summable'. */
  | 'cumulative';

export type TrackBy = 'space' | 'contact' | 'asset';

export interface IWithEntryType {
  readonly trackBy: TrackBy[];
  readonly valueType: TrackerValueType;
  readonly numberKind?: NumberKind;
}

export type TrackerCategory =
  'fitness' | 'health' | 'home' | 'vehicles' | string;

export interface ITrackerBrief extends IWithEntryType {
  readonly categories: readonly TrackerCategory[];
  readonly title: string;
  readonly emoji?: string;
}

export interface ITrackerPointBase {
  readonly ts?: Timestamp; // Time, if empty generate from id that is a UnixTime
  // readonly i?: number;
  // readonly I?: readonly number[];
  // readonly f?: number;
  // readonly F?: readonly number[];
  // readonly s?: string;
  // readonly b?: boolean;
  // readonly m?: number;
}

export interface IIntTrackerPoint extends ITrackerPointBase {
  readonly i: number;
  readonly I?: never;
  readonly f?: never;
  readonly F?: never;
  readonly s?: never;
  readonly b?: never;
  readonly m?: never;
}

export interface IIntsTrackerPoint extends ITrackerPointBase {
  readonly I: readonly number[];
  readonly i?: never;
  readonly f?: never;
  readonly F?: never;
  readonly s?: never;
  readonly b?: never;
  readonly m?: never;
}

export interface IFloatTrackerPoint extends ITrackerPointBase {
  readonly f: number;
  readonly i?: never;
  readonly I?: never;
  readonly F?: never;
  readonly s?: never;
  readonly b?: never;
  readonly m?: never;
}

export interface IFloatsTrackerPoint extends ITrackerPointBase {
  readonly F: readonly number[];
  readonly i?: never;
  readonly I?: never;
  readonly f?: never;
  readonly s?: never;
  readonly b?: never;
  readonly m?: never;
}

export interface IStringTrackerPoint extends ITrackerPointBase {
  readonly s: string;
  readonly i?: never;
  readonly I?: never;
  readonly f?: never;
  readonly F?: never;
  readonly b?: never;
  readonly m?: never;
}

export interface IBoolTrackerPoint extends ITrackerPointBase {
  readonly b: boolean;
  readonly i?: never;
  readonly I?: never;
  readonly f?: never;
  readonly F?: never;
  readonly s?: never;
  readonly m?: never;
}

export interface IMoneyTrackerPoint extends ITrackerPointBase {
  readonly m: number; // Amount in cents
  readonly i?: never;
  readonly I?: never;
  readonly f?: never;
  readonly F?: never;
  readonly s?: never;
}

export type TrackerPointBrief =
  | IIntTrackerPoint
  | IIntsTrackerPoint
  | IFloatTrackerPoint
  | IFloatsTrackerPoint
  | IStringTrackerPoint
  | IBoolTrackerPoint
  | IMoneyTrackerPoint;

// interface ITrackerPointBrief {
// 	readonly ts?: Timestamp; // Time, if empty generate from id that is a UnixTime
//
// 	// // TODO: This should be defined using TrackerPoint
// 	readonly i?: number;
// 	readonly f?: number;
// 	readonly b?: boolean;
// 	readonly s?: string;
// }

/** Sharing metadata is descriptive; it does not establish read/write authority. */
export interface ITrackerShareDetail {
  readonly byUserID: string;
  readonly permissions: readonly string[];
}

export interface ITrackerDbo extends ITrackerBrief {
  /** Absent on older projections. */
  readonly status?: 'active' | 'archived';
  /** Entity reference -> destination Space ID -> share. */
  readonly sharedTo?: Readonly<
    Record<string, Readonly<Record<string, ITrackerShareDetail>>>
  >;
  /** Source Space ID -> entity reference -> share. */
  readonly sharedFrom?: Readonly<
    Record<string, Readonly<Record<string, ITrackerShareDetail>>>
  >;
  readonly entries: Readonly<
    Record<string, Readonly<Record<string, TrackerPointBrief>>>
  >;
}

export interface ITracker extends IIdAndOptionalBriefAndOptionalDbo<
  ITrackerBrief,
  ITrackerDbo
> {
  space: ISpaceContext;
}

export const standardTrackers: readonly IIdAndBrief<ITrackerBrief>[] = [
  {
    id: '_push_ups',
    brief: {
      trackBy: ['contact'],
      valueType: 'int',
      numberKind: 'summable',
      categories: ['fitness'],
      title: 'Push-ups',
      emoji: '🏋️',
    },
  },
  {
    id: '_pull_ups',
    brief: {
      trackBy: ['contact'],
      valueType: 'int',
      numberKind: 'summable',
      categories: ['fitness'],
      title: 'Pull-ups',
      emoji: '🏋️',
    },
  },
  {
    id: '_squats',
    brief: {
      trackBy: ['contact'],
      valueType: 'int',
      categories: ['fitness'],
      title: 'Squats',
      emoji: '🏋️',
    },
  },
  {
    id: '_weight',
    brief: {
      trackBy: ['contact'],
      valueType: 'float',
      categories: ['fitness', 'health'],
      title: 'Weight',
      emoji: '⚖️',
    },
  },
  {
    id: '_body_temperature',
    brief: {
      trackBy: ['contact'],
      valueType: 'float',
      categories: ['health'],
      title: 'Temperature',
      emoji: '🌡️',
    },
  },
  {
    id: '_mileage',
    brief: {
      trackBy: ['asset'],
      valueType: 'int',
      categories: ['vehicles'],
      title: 'Mileage',
      emoji: '🛣️',
    },
  },
  {
    id: '_fuel',
    brief: {
      trackBy: ['asset'],
      valueType: 'float',
      categories: ['vehicles'],
      title: 'Fuel',
      emoji: '⛽',
    },
  },
  {
    id: '_electricity',
    brief: {
      trackBy: ['asset'],
      valueType: 'int',
      categories: ['home'],
      title: 'Electricity',
      emoji: '💡',
    },
  },
  {
    id: '_lpg',
    brief: {
      trackBy: ['asset'],
      valueType: 'int',
      categories: ['home'],
      title: 'LPG',
      emoji: '🔥',
    },
  },
  {
    id: '_heating',
    brief: {
      trackBy: ['asset'],
      valueType: 'money',
      categories: ['home'],
      title: 'Heating',
      emoji: '🔥',
    },
  },
];

export function isStandardTracker(id: string): boolean {
  return id.startsWith('_');
}

export function getStandardTrackerTitle(id: string): string {
  switch (id) {
    case '_pull_ups':
      return 'Pull-ups';
    case '_push_ups':
      return 'Push-ups';
    default:
      if (id.startsWith('_')) {
        return id.substring(1);
      }
      return id;
  }
}
