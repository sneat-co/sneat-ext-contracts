import type { ICommitmentAccountabilityView } from './goal-accountability';
import type { IGoalTermsSummary } from './goal-terms';

/** Financial lifecycle and the separately projected execution window differ. */
export type CommitmentLifecycleStatus =
  | 'draft' | 'active' | 'awaiting-verification' | 'kept'
  | 'broken' | 'cancelled' | 'void' | 'not-reported';
export type CommitmentBaseListOrder = 'identity-asc' | 'created-desc' | 'deadline-asc';
export type CommitmentTargetListOrder = 'target-asc' | 'target-desc';
export type CommitmentListOrder = CommitmentBaseListOrder | CommitmentTargetListOrder;

/** Select saved current meaning; clients cannot provide groups, paths or owners. */
export interface ICommitmentTargetExample {
  readonly commitmentID: string;
  readonly termsVersion: number;
  readonly clauseID: string;
}

interface CommitmentListPosition {
  /** Default20, maximum100; zero is accepted as the server default. */
  readonly limit?: number;
  /** Opaque non-confidential query/position metadata, not an access grant.
   * Never decode, edit, log or share. Max decoded64KiB; replay input on retry. */
  readonly cursor?: string;
  readonly status?: CommitmentLifecycleStatus;
}

export type IListCommitmentPageRequest = CommitmentListPosition & (
  | { readonly order?: CommitmentBaseListOrder; readonly target?: never }
  | { readonly order: CommitmentTargetListOrder; readonly target: ICommitmentTargetExample }
);

/** These capitalized keys intentionally retain the existing Go wire shape. */
export interface ICommitmentRefereeStageSummary {
  readonly Invited: number;
  readonly Accepted: number;
  readonly Declined: number;
  readonly Expired: number;
  readonly Voted: number;
}

/** Explicit owning response allowlist, not a raw stored commitment. */
export interface ICommitmentListView {
  readonly title: string;
  readonly status: CommitmentLifecycleStatus;
  readonly accountability: ICommitmentAccountabilityView;
  readonly stakeAmountCents: number;
  readonly stakeCurrency: string;
  readonly stakeHeld: boolean;
  readonly fundingPending: boolean;
  readonly referees: ICommitmentRefereeStageSummary;
  readonly refundDeliveryStatus: string;
  readonly refundDeliveryAttempts: number;
  readonly certificateReady: boolean;
  readonly goalTerms?: IGoalTermsSummary;
}

export type ICommitmentListItem =
  | {
      readonly id: string;
      readonly kind: 'commitment';
      readonly commitment: ICommitmentListView;
      /** Preserve wire timestamp precision; absence is not an invented date. */
      readonly createdAt?: string;
      readonly deadline?: string;
      /** Saved target in its canonical unit; never progress or success. */
      readonly targetValue?: number;
    }
  | {
      readonly id: string;
      readonly kind: 'unavailable';
      readonly commitment?: never;
      readonly createdAt?: never;
      readonly deadline?: never;
      readonly targetValue?: never;
    };

/** A live window, never a total, immutable snapshot or lifetime achievement.
 * Target coverage can omit missing/stale unmatched projections; identity order
 * remains the owning enumeration path. Refresh when selection/filter changes. */
export interface ICommitmentListPage {
  readonly items: readonly ICommitmentListItem[];
  /** Empty means this live query has no next page at evaluation time. */
  readonly nextCursor: string;
  readonly limit: number;
  readonly order: CommitmentListOrder;
  readonly status?: CommitmentLifecycleStatus;
  readonly consistency: 'live';
  readonly coverage: 'owner' | 'ordering-field-present' | 'compatible-target-projection';
  readonly evaluatedAt: string;
}

export type ICommitmentListErrorResponse = {
  readonly error:
    | { readonly code: 'list-order-unavailable'; readonly restartOrder: 'identity-asc' }
    | {
        readonly code: 'method_not_allowed' | 'dependency_unavailable' | 'unauthenticated'
          | 'rate_limit' | 'unsupported_media_type' | 'invalid_request'
          | 'request_too_large' | 'target-selection-unavailable' | 'list_unavailable';
        readonly restartOrder?: never;
      };
};

/** Authenticated owner point read; visibility/referee tokens confer no access. */
export interface IReadOwnedCommitmentRequest {
  readonly commitmentID: string;
}

/** Same filtered row as an owner page, without searching a list or implying
 * history completeness. A cold URL can read current held/refund metadata. */
export interface IOwnedCommitmentRead {
  readonly item: ICommitmentListItem;
  readonly consistency: 'live';
  readonly evaluatedAt: string;
}

export interface IOwnedCommitmentReadErrorResponse {
  readonly error: {
    readonly code: 'method_not_allowed' | 'dependency_unavailable' | 'unauthenticated'
      | 'rate_limit' | 'unsupported_media_type' | 'invalid_request'
      | 'request_too_large' | 'not_found' | 'read_unavailable';
  };
}
