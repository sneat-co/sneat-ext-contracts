import { describe, expect, expectTypeOf, it } from 'vitest';
import { writeFileSync } from 'node:fs';
import {
  COMMITMENT_GOAL_CONTEXT_SCHEMA_VERSION,
  COMMITMENT_GOAL_SCHEMA_VERSION,
  type ICommitmentGoalDefinition,
  type ICommitmentListItem,
  type ICommitmentListPage,
  type ICommitmentListView,
  type IContextualNumericGoalClause,
  type IGoalTermsPage,
  type IListCommitmentPageRequest,
} from './index';

const evaluation = {
  kind: 'by-deadline', startAt: '2026-11-03T12:00:00Z', timeZone: 'Europe/Dublin',
} as const;
const running = {
  version: 1, activity: { kind: 'declared', value: 'Running' }, role: 'event finish time',
} as const;
const halfMarathon = {
  schemaVersion: 2, conjunction: 'all', evaluation: {...evaluation,kind:'single-event'},
  clauses: [
    { id: 'event', kind: 'event', eventType: 'certified-half-marathon', rubric: 'Official event result.' },
    { id: 'time', kind: 'numeric', metric: 'duration', comparator: 'lt', target: 10800, unitOfMeasure: 'second', context: running },
  ],
} as const satisfies ICommitmentGoalDefinition;
const pushUps = {
  schemaVersion: 2, conjunction: 'all', evaluation: { ...evaluation, kind: 'single-day' },
  clauses: [{
    id: 'push', kind: 'numeric', metric: 'repetitions', comparator: 'gte', target: 100, unitOfMeasure: 'count',
    context: { version: 1, activity: { kind: 'trackus-standard', value: '_push_ups' }, role: 'daily total' },
  }],
} as const satisfies ICommitmentGoalDefinition;
const pullUps = {
  ...pushUps, evaluation: { ...evaluation, kind: 'single-set' },
  clauses: [{ ...pushUps.clauses[0], id: 'pull', target: 20, context: {
    version: 1, activity: { kind: 'trackus-standard', value: '_pull_ups' }, role: 'consecutive repetitions',
  } }],
} as const satisfies ICommitmentGoalDefinition;
const spanish = {
  schemaVersion: 2, conjunction: 'all', evaluation,
  clauses: [{ id: 'assessment', kind: 'assessment', scheme: 'cefr', subject: 'Spanish', minimumLevel: 'B1', rubric: 'Agreed independent B1 assessment.' }],
} as const satisfies ICommitmentGoalDefinition;
const reading = {
  schemaVersion: 2, conjunction: 'all', evaluation,
  clauses: [{ id: 'books', kind: 'numeric', metric: 'distinct-items', comparator: 'gte', target: 12, unitOfMeasure: 'count', itemKind: 'book' }],
} as const satisfies ICommitmentGoalDefinition;
const meditation = {
  schemaVersion: 2, conjunction: 'all', evaluation: {
    ...evaluation, kind: 'every-period', periodUnit: 'day', anchor: 'local-start-date', partialPeriods: 'include',
  },
  clauses: [{ id: 'duration', kind: 'numeric', metric: 'duration', comparator: 'gte', target: 600, unitOfMeasure: 'second', context: {
    version: 1, activity: { kind: 'declared', value: 'Meditation' }, role: 'daily practice',
  } }],
} as const satisfies ICommitmentGoalDefinition;

// These compile-time assertions are intentionally exercised by a separate tsc
// gate: Vitest transpilation alone is not a type-safety verdict.
describe('public goal and bounded page wire types', () => {
  it('retains schema1 and describes the six original goal families', () => {
    const legacy = {
      schemaVersion: 1, conjunction: 'all', evaluation: {...evaluation,kind:'single-event'},
      clauses: [{ id: 'distance', kind: 'numeric', metric: 'distance', comparator: 'gte', target: 5000, unitOfMeasure: 'meter' }],
    } as const satisfies ICommitmentGoalDefinition;
    expect(COMMITMENT_GOAL_SCHEMA_VERSION).toBe(1);
    expect(COMMITMENT_GOAL_CONTEXT_SCHEMA_VERSION).toBe(2);
    expectTypeOf(legacy).toExtend<ICommitmentGoalDefinition>();
    // Optional release qualification emits these exact checked values for the
    // owning Go validator, not a second hand-authored set of lookalike examples.
    const proofPath=process.env['IBETICAN_GOAL_FIXTURE_PROOF'];
    if(proofPath){
      const goals={legacy,halfMarathon,pushUps,pullUps,spanish,reading,meditation};
      writeFileSync(proofPath,JSON.stringify(Object.entries(goals).map(([name,definition])=>({name,definition,deadline:'2027-01-02T12:00:00Z'}))));
    }
    for (const goal of [halfMarathon, pushUps, pullUps, spanish, reading, meditation]) {
      expectTypeOf(goal).toExtend<ICommitmentGoalDefinition>();
    }
  });

  it('rejects silent context upgrades and incompatible standard activities', () => {
    const missingContext = { ...halfMarathon, clauses: [{ id: 'distance', kind: 'numeric', metric: 'distance', comparator: 'gte', target: 5000, unitOfMeasure: 'meter' }] } as const;
    expectTypeOf(missingContext).not.toExtend<ICommitmentGoalDefinition>();
    const legacyWithContext = { ...halfMarathon, schemaVersion: 1 } as const;
    expectTypeOf(legacyWithContext).not.toExtend<ICommitmentGoalDefinition>();
    const standardDistance = {
      ...pushUps.clauses[0], metric: 'distance', unitOfMeasure: 'meter',
    } as const;
    expectTypeOf(standardDistance).not.toExtend<IContextualNumericGoalClause>();
    const arbitraryStandard = {
      ...pushUps.clauses[0], context: { ...pushUps.clauses[0].context, activity: { kind: 'trackus-standard', value: 'my-private-tracker-id' } },
    } as const;
    expectTypeOf(arbitraryStandard).not.toExtend<IContextualNumericGoalClause>();
    const eventWithContext={...halfMarathon.clauses[0],context:running};
    const injectedContext={...halfMarathon,clauses:[eventWithContext]};
    expectTypeOf(injectedContext).not.toExtend<ICommitmentGoalDefinition>();
  });

  it('requires a saved example only for target order', () => {
    const request = { order: 'target-asc', limit: 20, target: { commitmentID: 'saved-goal', termsVersion: 2, clauseID: 'time' } } as const satisfies IListCommitmentPageRequest;
    expectTypeOf(request).toExtend<IListCommitmentPageRequest>();
    expectTypeOf<{ order: 'target-desc' }>().not.toExtend<IListCommitmentPageRequest>();
    expectTypeOf<{ order: 'deadline-asc'; target: typeof request.target }>().not.toExtend<IListCommitmentPageRequest>();
    expectTypeOf<'upcoming'>().not.toExtend<NonNullable<IListCommitmentPageRequest['status']>>();
  });

  it('keeps unavailable rows and unsupported history free of invented meaning', () => {
    expectTypeOf<{ id: string; kind: 'unavailable'; targetValue: number }>().not.toExtend<ICommitmentListItem>();
    expectTypeOf<Extract<'homeSpaceID' | 'snapshotID' | 'trackerID' | 'targetSort',keyof ICommitmentListView>>().toEqualTypeOf<never>();
    expectTypeOf<{ supported: false; version: number; effectiveVersion: number; schemaVersion: number; totalPeriodCount: 0; definition: typeof halfMarathon }>().not.toExtend<IGoalTermsPage>();
    const page = {
      items: [{ id: 'unsupported-goal', kind: 'unavailable' }], nextCursor: '', limit: 20, order: 'identity-asc', consistency: 'live', coverage: 'owner', evaluatedAt: '2026-11-03T12:00:00.123456Z',
    } as const satisfies ICommitmentListPage;
    expectTypeOf(page).toExtend<ICommitmentListPage>();
    expect(page.evaluatedAt).toContain('.123456Z');
  });
});
