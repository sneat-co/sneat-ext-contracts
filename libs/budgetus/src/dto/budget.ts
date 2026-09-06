import { IListGroup } from "./list";

export interface IBudgetusSpaceDbo {
  listGroups?: IListGroup[];
}

/**
 * A monetary amount in **integer minor units**.
 *
 * `value` is always a whole number of the currency's smallest unit — €12.34 is
 * `{ currency: 'EUR', value: 1234 }`, never `12.34`. This matches the backend
 * wire form (`decimal.Decimal64p2` / `money.Amount`) and the calendarius
 * contract's `IAmount`, so no unit conversion happens in transit.
 *
 * There is exactly one unit convention on this type. Anything that shows an
 * amount to a human divides by 100 at the edge (render with
 * `value | decimal64p2 | currency: currency`); anything that reads a human's
 * input multiplies by 100 before it becomes an `IMoney`. Do not introduce a
 * second, major-unit money shape alongside this one.
 */
export interface IMoney {
  readonly currency: string;
  /** Integer minor units. €12.34 → 1234. Must be a safe integer. */
  readonly value: number;
}

export type BudgetLineSource = "asset-renewal" | "happening" | "gift";

/**
 * Why a line is shown but NOT added to any total.
 *
 * - `no-amount` — the source carries no amount at all (e.g. an asset renewal
 *   with no premium recorded). Rendered as a "set a target" affordance. Summing
 *   it as 0.00 would be a lie, so it is excluded until a target is set.
 * - `unsupported-term` — the price is quoted per day/hour/minute/second. There
 *   is no honest way to turn that into a monthly figure without knowing the
 *   duration of each occurrence, so the amount is shown with its real term and
 *   left out of the totals.
 */
export type BudgetLineExclusionReason = "no-amount" | "unsupported-term";

/**
 * The member bucket for costs that belong to no particular participant.
 *
 * `*` cannot start a contact id in this fleet (it is the reserved
 * `AnyRelatedID` marker), so this can never collide with a real contact.
 */
export const SHARED_BUDGET_MEMBER_ID = "*shared";

export interface IBudgetLineItem {
  id: string;
  title: string;
  dateISO: string;
  amount: IMoney;
  source: BudgetLineSource;
  sourceRef?: string;
  targetAmount?: IMoney;
  isSurprise?: boolean;

  // --- Line identity (0.2.0) -------------------------------------------------
  // A happening line's id is `happening:{happeningID}:{priceID}:{monthISO}`.
  // These fields carry the same components in parsed form so consumers can link
  // back to the source without re-parsing the id.

  /** Source happening, for happening/gift lines. */
  happeningID?: string;
  /** Backend-stable price id. One line per expense-flagged price. */
  priceID?: string;
  /** `YYYY-MM` of the occurrence this line represents. */
  occurrenceMonthISO?: string;

  // --- Rollup inputs (0.2.0) -------------------------------------------------

  /**
   * Contact ids of the participants this cost is attributed to.
   *
   * Empty or absent means the cost belongs to the {@link SHARED_BUDGET_MEMBER_ID}
   * bucket. Ids are short (`contactID`), never the long `contactID@spaceID` form.
   */
  memberIDs?: string[];

  /** Exact attribution when participants vary between occurrences. */
  memberAllocations?: readonly {
    readonly memberID: string;
    readonly amount: IMoney;
  }[];

  /**
   * True for a recurring commitment (a priced recurring happening) as opposed to
   * a one-off dated cost. Regular lines are what the "Regular monthly expenses"
   * headline totals.
   */
  isRegular?: boolean;

  /**
   * The line's normalized cost for one month, in minor units.
   *
   * This is what rollups sum. It is absent exactly when {@link excludedFromTotals}
   * is set.
   */
  monthlyAmount?: IMoney;

  /** Set when the line is displayed but deliberately not summed. */
  excludedFromTotals?: BudgetLineExclusionReason;

  /**
   * Human-readable pricing term, e.g. `"€15.00 per session"` or `"€9.00/hour"`.
   * Shown for excluded lines so the real number stays visible even when it
   * cannot honestly be turned into a monthly figure.
   */
  termLabel?: string;

  /** How occurrence-based cost was derived from its owning schedule. */
  occurrenceProvenance?: IBudgetOccurrenceProvenance;
  /** Source-owned assets linked to the happening; they do not affect payer math. */
  relatedAssets?: readonly {
    readonly assetID: string;
    readonly spaceID?: string;
  }[];
}

export interface IBudgetOccurrenceProvenance {
  /** Scheduled occurrences included in this line's baseline estimate. */
  readonly scheduledCount: number;
  /** Scheduled occurrences carrying a date-specific cancellation marker. */
  readonly canceledCount: number;
  /** Scheduled occurrences whose slot was changed for this date. */
  readonly adjustedCount: number;
  /** The recorded Calendarius price multiplier; participant count is separate. */
  readonly priceQuantity: number;
  /** Cancellations remain charged because Calendarius records no waiver policy. */
  readonly estimateBasis: "pre-adjustment-schedule";
}

export interface IBudgetMonthGroup {
  monthISO: string;
  total: IMoney;
  items: IBudgetLineItem[];
}

/** Per-participant totals within a single currency. */
export interface IMemberBudgetTotals {
  /** A contact id, or {@link SHARED_BUDGET_MEMBER_ID}. */
  memberID: string;
  /** Display name, when the host could resolve one. */
  title?: string;
  /** Recurring commitments attributable to this member, per month. */
  regularMonthlyTotal: IMoney;
  /** Everything attributable to this member across the window. */
  windowTotal: IMoney;
}

/**
 * A complete rollup for ONE currency.
 *
 * Amounts are never summed across currencies — a space with EUR and USD costs
 * produces two of these, side by side.
 */
export interface IBudgetCurrencyRollup {
  currency: string;
  byMonth: IBudgetMonthGroup[];
  /** Sum of every summable line across the window. */
  windowTotal: IMoney;
  /** Sum of the per-month cost of recurring commitments only. */
  regularMonthlyTotal: IMoney;
  mostExpensiveMonthISO: string;
  /** Per-participant breakdown, biggest first. */
  byMember: IMemberBudgetTotals[];
}

/**
 * The span a rollup covers: `months` calendar months starting at `fromMonthISO`.
 */
export interface IBudgetWindow {
  /** `YYYY-MM` — the first month included. */
  readonly fromMonthISO: string;
  /** How many months the window spans. At least 1. */
  readonly months: number;
}

export const DEFAULT_BUDGET_WINDOW_MONTHS = 12;

export interface IBudgetRollup {
  /** One entry per currency present in the space. */
  byCurrency: IBudgetCurrencyRollup[];
  /** The window actually used (a request may omit it and take the default). */
  window: IBudgetWindow;
  /**
   * Lines shown but not summed — see {@link BudgetLineExclusionReason}. Kept out
   * of `byCurrency` so a total can never accidentally include them.
   */
  excludedItems: IBudgetLineItem[];
  /** Completeness of each independently loaded owning source. */
  sourceStatuses?: readonly IBudgetSourceStatus[];
}

export type BudgetSourceID = "assetus-renewals" | "calendarius-happenings";

export type BudgetSourceState =
  "complete" | "partial" | "unavailable" | "unsupported";

export type BudgetSourceStatusReasonCode =
  | "calendar-adjustments-unavailable"
  | "moved-in-occurrences-unknown"
  | "unsupported-recurrence"
  | "cancellation-charge-unknown"
  | "occurrence-cap-reached"
  | "foreign-related-reference"
  | "source-read-failed";

export interface IBudgetSourceStatusReason {
  readonly code: BudgetSourceStatusReasonCode;
  readonly message: string;
}

export interface IBudgetSourceStatus {
  readonly sourceID: BudgetSourceID;
  readonly state: BudgetSourceState;
  readonly reasons?: readonly IBudgetSourceStatusReason[];
  /** Inclusive source query bounds when the source is date-bounded. */
  readonly queryWindow?: {
    readonly fromDateISO: string;
    readonly toDateISO: string;
  };
  readonly cap?: {
    readonly limit: number;
    readonly observed?: number;
  };
}

export interface IBudgetOverridePatch {
  targetAmount?: IMoney;
  isSurprise?: boolean;
}
