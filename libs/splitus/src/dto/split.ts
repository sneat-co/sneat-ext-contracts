/** @deprecated Use the Splitus bill contract types exported from `bill-v1`. */
export type SplitMode = 'equally' | 'exact-amount' | 'percentage';

/** @deprecated Use `SplitusCurrencyCode` from the Splitus bill contract. */
export type CurrencyCode = 'EUR' | 'USD';

/**
 * One participant's custom share for `exact-amount` / `percentage` split
 * modes. An omitted/empty `contactID` denotes the payer's own share. Ignored
 * (and may be omitted) for `equally`, which the backend computes itself.
 */
/** @deprecated Use explicit paid and owed allocations in `ICreateSplitusBillV1Request`. */
export interface ISplitShare {
  readonly contactID?: string;
  /** Decimal string, e.g. "35.00" — required for `exact-amount`. */
  readonly amount?: string;
  /** Decimal string, e.g. "33.34" — required for `percentage`. */
  readonly percent?: string;
}

/** @deprecated Use `ICreateSplitusBillV1Request`. */
export interface ICreateSplitRequest {
  readonly spaceID: string;
  readonly title?: string;
  readonly currency: CurrencyCode;
  /** Decimal string, e.g. "90.00" — the total expense. */
  readonly amount: string;
  /** Defaults to `equally` server-side when omitted. */
  readonly splitMode?: SplitMode;
  /**
   * contactus contact IDs of the space. The payer (the authenticated user)
   * is always a participant and must NOT be listed here.
   */
  readonly participantContactIDs: string[];
  /** Required for `exact-amount` / `percentage`; ignored for `equally`. */
  readonly shares?: ISplitShare[];
}

/** @deprecated Use `ISplitusDebtusObligationV1`. */
export interface ICreateSplitTransfer {
  readonly id: string;
  readonly contactID: string;
  readonly amount: number;
}

/** @deprecated Use `ICreateSplitusBillV1Response`. */
export interface ICreateSplitResponse {
  readonly id: string;
  /**
   * The Debtus transfers holding the balances — the only who-owes-what
   * records for this split. Splitus itself persists no balance.
   */
  readonly transfers: ICreateSplitTransfer[];
}

/**
 * "settled" or "outstanding", derived server-side by reading the linked
 * Debtus transfers — never computed or cached on the client.
 */
/** @deprecated Use `SplitusDebtusSettlementStatus`. */
export type SplitShareStatus = 'settled' | 'outstanding';

/** @deprecated Use `ISplitusBillAllocationV1`. */
export interface ISplitParticipant {
  readonly contactID?: string;
  readonly userID?: string;
  readonly name: string;
  readonly share: number;
  readonly isPayer?: boolean;
  readonly status: SplitShareStatus;
}

/** @deprecated Use `ISplitusBillV1`. */
export interface ISplit {
  readonly id: string;
  readonly title?: string;
  readonly currency: CurrencyCode;
  readonly amount: number;
  readonly status: string;
  readonly participants: ISplitParticipant[];
}

/** @deprecated Use `ISplitusBillListItemV1`. */
export interface ISplitListItem {
  readonly id: string;
  readonly title?: string;
  readonly amount: number;
  readonly currency: CurrencyCode;
  readonly status: string;
  readonly membersCount: number;
}
