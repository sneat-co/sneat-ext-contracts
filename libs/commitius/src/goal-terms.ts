import type {
  GOAL_CALENDAR_NORMALIZATION_VERSION,
  GoalEvaluationKind,
  GoalPeriodPolicy,
  ICommitmentGoalDefinition,
  INormalizedGoalPeriod,
  IReviewedCommitmentGoal,
} from './goal-definition';

/** Additive structured fields for the existing positive-stake create command.
 * Authenticated actor identity is separate trusted host context. */
export interface ITypedCommitmentCreateFields extends IReviewedCommitmentGoal {
  readonly homeSpaceID: string;
}

/** Required for new structured evidence; legacy omission is a distinct path. */
export interface ITypedGoalVersionGuard {
  readonly expectedTermsVersion: number;
}

export type IGoalCalendarWindow = {
  readonly normalizationVersion: typeof GOAL_CALENDAR_NORMALIZATION_VERSION;
  readonly normalizedStartAt: string;
  readonly normalizedDeadline: string;
  readonly timeZone: string;
  readonly localStartLabel: string;
  readonly localDeadlineLabel: string;
} & (
  | { readonly kind: Exclude<GoalEvaluationKind, 'every-period'> }
  | ({ readonly kind: 'every-period' } & GoalPeriodPolicy)
);

export interface IGoalCalendarPeriodPage {
  readonly offset: number;
  readonly periods: readonly INormalizedGoalPeriod[];
  readonly hasMore: boolean;
}

interface GoalTermsIdentity {
  readonly version: number;
  readonly schemaVersion: number;
}

interface GoalTermsPageIdentity extends GoalTermsIdentity {
  readonly effectiveVersion: number;
}

/** Filtered semantic data, never the raw Space-bound snapshot. Preserve ISO
 * instants as strings: JS Date round-trips would discard server microseconds. */
export type IGoalTermsPage = GoalTermsPageIdentity & (
  | {
      readonly supported: true;
      readonly title: string;
      readonly deadline: string;
      readonly definition: ICommitmentGoalDefinition;
      readonly calendar: IGoalCalendarWindow;
      readonly calendarDigest: string;
      readonly totalPeriodCount: number;
      readonly periodPage: IGoalCalendarPeriodPage;
    }
  | {
      readonly supported: false;
      readonly totalPeriodCount: 0;
      readonly title?: never;
      readonly deadline?: never;
      readonly definition?: never;
      readonly calendar?: never;
      readonly calendarDigest?: never;
      readonly periodPage?: never;
    }
);

export interface IReadGoalTermsRequest {
  readonly commitmentID: string;
  /** Omitted/zero selects current; historical selection remains role-filtered. */
  readonly version?: number;
  readonly offset?: number;
  readonly limit?: number;
  readonly refereeToken?: string;
}

export type GoalExecutionPhase = 'upcoming' | 'open' | 'ended';
export type GoalReadinessBlocker =
  | 'evaluation-not-started' | 'evaluation-ended' | 'funding-unresolved'
  | 'stake-not-funded' | 'referee-not-accepted' | 'acceptance-window-ended'
  | 'verification-in-progress' | 'commitment-closed';

export type IGoalTermsSummary = GoalTermsIdentity & (
  | {
      readonly supported: true;
      readonly evaluationKind: GoalEvaluationKind;
      readonly startAt: string;
      readonly activityEligibleFrom: string;
      readonly timeZone: string;
      readonly periodCount?: number;
      readonly executionPhase: GoalExecutionPhase;
      /** Owner-only, omitted for other permitted audiences. */
      readonly ready?: boolean;
      readonly blockers?: readonly GoalReadinessBlocker[];
    }
  | { readonly supported: false }
);

export interface IGoalTermsChangeInput extends IReviewedCommitmentGoal, ITypedGoalVersionGuard {
  readonly title: string;
  readonly deadline: string;
}

export interface IRequestGoalTermsChange {
  readonly commitmentID: string;
  readonly field: 'goal-terms';
  readonly goalTerms: IGoalTermsChangeInput;
  readonly newValue?: never;
}

export type GoalChangeStatus = 'pending' | 'approved' | 'rejected' | 'withdrawn' | 'unappliable' | 'owner-confirmed';

export interface IGoalTermsChangeOutcome {
  readonly changeID: string;
  readonly baseVersion: number;
  readonly proposedVersion: number;
  /** Paid proposals await consent; self-report amendments close by owner decision. */
  readonly status: 'pending' | 'owner-confirmed';
}

export interface IReadGoalChangeRequest {
  readonly commitmentID: string;
  readonly changeID: string;
  readonly baseOffset?: number;
  readonly proposedOffset?: number;
  readonly limit?: number;
  readonly refereeToken?: string;
}

export interface IGoalTermsChangeView {
  readonly changeID: string;
  readonly baseVersion: number;
  readonly proposedVersion: number;
  readonly status: GoalChangeStatus;
  readonly unappliableReason?: string;
  readonly requestedAt: string;
  readonly closedAt?: string;
  readonly base: IGoalTermsPage;
  readonly proposed: IGoalTermsPage;
  readonly requiredApprovalCount: number;
  readonly approvalCount: number;
  readonly yourDecision?: 'pending' | 'approved' | 'rejected';
  /** Server-authorized actions; clients cannot infer authority from a count. */
  readonly canApprove: boolean;
  readonly canReject: boolean;
  readonly canWithdraw: boolean;
}

export interface IWithdrawGoalChangeRequest extends ITypedGoalVersionGuard {
  readonly commitmentID: string;
  readonly changeID: string;
}

/** Owner-only listing; public visibility/referee links confer no list access. */
export interface IListGoalVersionsRequest {
  readonly commitmentID: string;
  readonly beforeVersion?: number;
  readonly limit?: number;
}

export interface IGoalVersionSummary {
  readonly version: number;
  readonly current: boolean;
  readonly recordedAt: string;
  readonly schemaVersion: number;
  readonly supported: boolean;
}

export interface IGoalVersionHistory {
  readonly effectiveVersion: number;
  readonly items: readonly IGoalVersionSummary[];
  readonly hasMore: boolean;
  readonly nextBeforeVersion?: number;
}

export interface IListGoalChangesRequest {
  readonly commitmentID: string;
  readonly afterChangeID?: string;
  readonly limit?: number;
}

export interface IGoalChangeSummary {
  readonly changeID: string;
  readonly baseVersion: number;
  readonly proposedVersion: number;
  readonly requestedAt: string;
  readonly closedAt?: string;
  readonly status: GoalChangeStatus;
  readonly unappliableReason?: string;
  readonly supported: boolean;
}

export interface IGoalChangeHistory {
  readonly effectiveVersion: number;
  readonly items: readonly IGoalChangeSummary[];
  readonly hasMore: boolean;
  readonly nextAfterChangeID?: string;
}
