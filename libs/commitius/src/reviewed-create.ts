import type { ITypedCommitmentCreateFields } from './goal-terms';

/** Optional existing Trackus metadata, never evidence acceptance or settlement. */
export interface ICommitmentCreateProgressPlan {
  readonly trackerID: string;
  readonly targetPerDay: number;
  readonly label: string;
  readonly unit?: string;
  readonly timeZone: string;
}

/** Positive-stake creation only. The separate self-report command has its own
 * contract. Amounts, currencies, visibility and all goal meaning remain subject
 * to server validation; this type does not establish payment or success. */
export interface IReviewedPaidCommitmentInput extends ITypedCommitmentCreateFields {
  readonly title: string;
  readonly deadline: string;
  readonly stakeAmountCents: number;
  readonly why?: string;
  readonly category?: string;
  readonly progressPlan?: ICommitmentCreateProgressPlan;
  readonly originCommitmentID?: string;
  readonly referrerCommitmentID?: string;
  readonly stakeCurrency?: string;
  readonly visibility?: string;
  readonly whyVisibility?: string;
  readonly stakeVisibility?: string;
  readonly evidenceVisibility?: string;
  readonly refereeHandles?: readonly string[];
  readonly votingDeadlineDays?: number;
}

/** Retain one actor-scoped operation ID and the exact immutable reviewed input
 * across an unknown reply. Different input under the same operation conflicts.
 * Authentication actor and resulting commitment ID are server-owned. */
export interface IReviewedCommitmentCreateRequest {
  readonly operationID: string;
  readonly commitment: IReviewedPaidCommitmentInput;
}

/** HTTP201 returns the original ID for a successful equal replay too. The
 * command never funds a goal. Replays can occur after later lifecycle changes:
 * this ID response proves neither current funding state nor delivery. */
export interface IReviewedCommitmentCreateResponse {
  readonly id: string;
}

/** Unknown or interrupted command outcomes are never proof of rollback. */
export type ReviewedCommitmentCreateErrorCode =
  | 'method_not_allowed' | 'dependency_unavailable' | 'unauthenticated'
  | 'rate_limit' | 'unsupported_media_type' | 'invalid_request'
  | 'request_too_large' | 'operation_conflict' | 'access_denied'
  | 'operation_unconfirmed';
