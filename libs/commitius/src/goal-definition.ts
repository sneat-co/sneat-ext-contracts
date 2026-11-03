/** Additive wire values. These definitions do not establish a result or move money. */
export const COMMITMENT_GOAL_SCHEMA_VERSION = 1 as const;
export const COMMITMENT_GOAL_CONTEXT_SCHEMA_VERSION = 2 as const;
export const GOAL_CALENDAR_NORMALIZATION_VERSION = "goal-calendar-v1" as const;

export type GoalEvaluationKind =
  "by-deadline" | "single-event" | "single-day" | "single-set" | "every-period";

export type GoalComparator = "gte" | "gt" | "lte" | "lt" | "eq";

interface GoalClauseIdentity {
  /** Stable within the definition; the server rejects duplicate or invalid IDs. */
  readonly id: string;
}

interface NumericGoalTarget extends GoalClauseIdentity {
  readonly kind: "numeric";
  readonly comparator: GoalComparator;
  /** Positive safe integer in the canonical unit; validated by the server. */
  readonly target: number;
}

/** Canonical units use the existing Trackus vocabulary, without a runtime dependency. */
type GoalNumericMetric =
  (
    | { readonly metric: "distance"; readonly unitOfMeasure: "meter" }
    | { readonly metric: "duration"; readonly unitOfMeasure: "second" }
    | { readonly metric: "repetitions"; readonly unitOfMeasure: "count" }
    | {
        readonly metric: "distinct-items";
        readonly unitOfMeasure: "count";
        readonly itemKind: string;
      }
    | {
        readonly metric: "occurrences";
        readonly unitOfMeasure: "count";
        readonly actionLabel: string;
      }
  );

/** Existing schema1 remains readable without inventing missing activity. */
export type INumericGoalClause = NumericGoalTarget & GoalNumericMetric & {
  readonly context?: never;
};

/** Immutable goal meaning, not a tracker binding, observation or authority. */
export interface IDeclaredGoalNumericContext {
  readonly version: 1;
  readonly activity: { readonly kind: "declared"; readonly value: string };
  readonly role: string;
}

export interface IStandardGoalNumericContext {
  readonly version: 1;
  readonly activity: {
    readonly kind: "trackus-standard";
    /** Public Trackus standard activity vocabulary; no runtime dependency. */
    readonly value: "_push_ups" | "_pull_ups";
  };
  readonly role: string;
}

/** Schema2 makes ambiguous numeric activity explicit. The server validates
 * trimmed bounded text, safe targets, duplicate axes and all saved semantics. */
export type IContextualNumericGoalClause = NumericGoalTarget & (
  | (Extract<GoalNumericMetric, { metric: "distance" | "duration" }> & {
      readonly context: IDeclaredGoalNumericContext;
    })
  | (Extract<GoalNumericMetric, { metric: "repetitions" }> & {
      readonly context: IDeclaredGoalNumericContext | IStandardGoalNumericContext;
    })
  | (Extract<GoalNumericMetric, { metric: "distinct-items" | "occurrences" }> & {
      readonly context?: IDeclaredGoalNumericContext;
    })
);

export interface IAssessmentGoalClause extends GoalClauseIdentity {
  readonly context?: never;
  readonly kind: "assessment";
  readonly scheme: "cefr";
  readonly subject: string;
  readonly minimumLevel: "A1" | "A2" | "B1" | "B2" | "C1" | "C2";
  readonly rubric: string;
}

export interface IEventGoalClause extends GoalClauseIdentity {
  readonly context?: never;
  readonly kind: "event";
  /** The certified21.0975km event, not an inferred21,100m GPS threshold. */
  readonly eventType: "certified-half-marathon";
  readonly rubric: string;
}

export interface IRubricGoalClause extends GoalClauseIdentity {
  readonly context?: never;
  readonly kind: "rubric";
  readonly description: string;
}

export type IGoalClause =
  | INumericGoalClause
  | IAssessmentGoalClause
  | IEventGoalClause
  | IRubricGoalClause;

interface GoalEvaluationStart {
  /** Reviewed ISO8601 instant. CreatedAt is independently server-owned. */
  readonly startAt: string;
  /** Validated, frozen IANA name; never a device-local fallback. */
  readonly timeZone: string;
}

export type GoalPeriodPolicy =
  | {
      readonly periodUnit: "day";
      readonly anchor: "local-start-date";
      readonly partialPeriods: "include" | "exclude";
      readonly requiredPeriodCount?: number;
    }
  | {
      readonly periodUnit: "week";
      readonly anchor: "local-start-date" | "iso-monday";
      readonly partialPeriods: "include" | "exclude";
      readonly requiredPeriodCount?: number;
    };

export type IGoalEvaluation = GoalEvaluationStart &
  (
    | { readonly kind: Exclude<GoalEvaluationKind, "every-period"> }
    | ({ readonly kind: "every-period" } & GoalPeriodPolicy)
  );

export type IContextualGoalClause =
  | IContextualNumericGoalClause
  | IAssessmentGoalClause
  | IEventGoalClause
  | IRubricGoalClause;

interface CommitmentGoalDefinitionContent {
  readonly templateID?: string;
  readonly conjunction: "all";
  readonly evaluation: IGoalEvaluation;
  readonly verificationRubric?: string;
}

export interface ILegacyCommitmentGoalDefinition extends CommitmentGoalDefinitionContent {
  readonly schemaVersion: typeof COMMITMENT_GOAL_SCHEMA_VERSION;
  readonly clauses: readonly IGoalClause[];
}

export interface IContextualCommitmentGoalDefinition extends CommitmentGoalDefinitionContent {
  readonly schemaVersion: typeof COMMITMENT_GOAL_CONTEXT_SCHEMA_VERSION;
  readonly clauses: readonly IContextualGoalClause[];
}

/** The enclosing create/change request supplies the authoritative deadline.
 * Explicit discriminants preserve schema1 without silently upgrading context. */
export type ICommitmentGoalDefinition =
  | ILegacyCommitmentGoalDefinition
  | IContextualCommitmentGoalDefinition;

export interface IGoalLocalDateTime {
  /** Zone-local ISO date/time without a UTC suffix or offset. */
  readonly localDateTime: string;
  /** Required to select one occurrence when the local time is ambiguous. */
  readonly offsetSeconds?: number;
}

export interface IGoalCalendarHorizon {
  readonly unit: "day" | "week" | "month" | "year";
  readonly count: number;
}

/** Pure calendar inputs; no user, Space, contact, commitment or payment identity. */
export type IGoalCalendarPreviewRequest = {
  readonly normalizationVersion: typeof GOAL_CALENDAR_NORMALIZATION_VERSION;
  readonly timeZone: string;
  /** Only the first preview may omit this; retain its returned concrete start. */
  readonly start?: IGoalLocalDateTime;
  readonly periodPage?: { readonly offset: number; readonly limit: number };
} & (
  | { readonly kind: Exclude<GoalEvaluationKind, "every-period"> }
  | ({ readonly kind: "every-period" } & GoalPeriodPolicy)
) &
  (
    | { readonly cutoff: IGoalLocalDateTime; readonly horizon?: never }
    | { readonly horizon: IGoalCalendarHorizon; readonly cutoff?: never }
  );

export interface INormalizedGoalPeriod {
  readonly id: string;
  readonly localLabel: string;
  /** Envelope bounds for ordering/display; membership uses intervals below. */
  readonly startAt: string;
  readonly endAt: string;
  readonly partial: boolean;
  /** Exact frozen membership; a civil date can have disjoint UTC intervals. */
  readonly intervals: readonly INormalizedGoalInterval[];
}

export interface INormalizedGoalInterval {
  /** Start-inclusive/end-exclusive absolute instants. */
  readonly startAt: string;
  readonly endAt: string;
}

export type IGoalCalendarPreview = {
  readonly normalizationVersion: typeof GOAL_CALENDAR_NORMALIZATION_VERSION;
  readonly normalizedStartAt: string;
  readonly normalizedDeadline: string;
  readonly timeZone: string;
  readonly localStartLabel: string;
  readonly localDeadlineLabel: string;
  /** SHA256 of complete normalized content, never an authorization capability. */
  readonly normalizationDigest: string;
  readonly totalPeriodCount: number;
  readonly periodPage: {
    readonly offset: number;
    readonly periods: readonly INormalizedGoalPeriod[];
    readonly hasMore: boolean;
  };
} & (
  | { readonly kind: Exclude<GoalEvaluationKind, "every-period"> }
  | ({ readonly kind: "every-period" } & GoalPeriodPolicy)
);

/** Structured submit must carry the calendar content actually reviewed. */
export interface IReviewedCommitmentGoal {
  readonly goalDefinition: ICommitmentGoalDefinition;
  readonly expectedCalendarDigest: string;
}

export interface IGoalCalendarOffsetChoice {
  readonly instant: string;
  readonly offsetSeconds: number;
}

export interface IGoalCalendarOffsetChoiceRequired {
  readonly code: "calendar_offset_choice_required";
  readonly message: string;
  readonly field: "start" | "cutoff";
  readonly localDateTime: string;
  readonly candidates: readonly IGoalCalendarOffsetChoice[];
}

export interface IGoalCalendarInputError {
  readonly code:
    | "invalid_calendar"
    | "invalid_json"
    | "request_too_large"
    | "unsupported_media_type"
    | "method_not_allowed"
    | "calendar_unavailable";
  readonly message: string;
  readonly field?: string;
}

export interface IGoalCalendarErrorResponse {
  readonly error: IGoalCalendarOffsetChoiceRequired | IGoalCalendarInputError;
}
