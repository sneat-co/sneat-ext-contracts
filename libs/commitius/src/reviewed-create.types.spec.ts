import { describe, expect, expectTypeOf, it } from 'vitest';
import type {
  IReviewedCommitmentCreateRequest,
  IReviewedCommitmentCreateResponse,
  ReviewedCommitmentCreateErrorCode,
} from './index';

const request = {
  operationID: 'qualified-reviewed-create-operation',
  commitment: {
    homeSpaceID: 'personal1',
    title: 'Launch my project',
    deadline: '2027-04-02T00:00:00Z',
    goalDefinition: {
      schemaVersion: 2, conjunction: 'all',
      evaluation: { kind: 'by-deadline', timeZone: 'UTC', startAt: '2027-01-02T00:00:00Z' },
      clauses: [{ id: 'launch', kind: 'rubric', description: 'Publish a working project sign-up page at its public URL.' }],
    },
    expectedCalendarDigest: '0'.repeat(64),
    stakeAmountCents: 2500,
    stakeCurrency: 'EUR',
    refereeHandles: ['friend@example.test'],
  },
} as const satisfies IReviewedCommitmentCreateRequest;

describe('reviewed web creation public contract', () => {
  it('keeps immutable goal input nested under its independent retry identity', () => {
    expectTypeOf(request).toExtend<IReviewedCommitmentCreateRequest>();
    expect(JSON.parse(JSON.stringify(request))).toEqual(request);
    expectTypeOf({ id: 'saved-original-id' }).toExtend<IReviewedCommitmentCreateResponse>();
    expectTypeOf('operation_unconfirmed' as const).toExtend<ReviewedCommitmentCreateErrorCode>();
  });

  it('requires replay identity, actual goalDefinition wire name, and paid backing', () => {
    expectTypeOf(request.commitment).not.toExtend<IReviewedCommitmentCreateRequest>();
    const { operationID: _operationID, ...withoutOperation } = request;
    expectTypeOf(withoutOperation).not.toExtend<IReviewedCommitmentCreateRequest>();
    const { goalDefinition, ...withoutDefinition } = request.commitment;
    expectTypeOf({ operationID: request.operationID, commitment: { ...withoutDefinition, definition: goalDefinition } }).not.toExtend<IReviewedCommitmentCreateRequest>();
    const { stakeAmountCents: _stake, ...withoutStake } = request.commitment;
    expectTypeOf({ operationID: request.operationID, commitment: withoutStake }).not.toExtend<IReviewedCommitmentCreateRequest>();
    // Keep deliberately removed values observable to satisfy the standard lint policy.
    expect(_operationID).toBe(request.operationID);
    expect(_stake).toBe(2500);
  });
});
