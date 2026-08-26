/** Mirrors models4splitus.SplitMode (backend/splitus/models4splitus). */
export type SplitMode = 'equally' | 'exact-amount' | 'percentage';

export type CurrencyCode = 'EUR' | 'USD';

/**
 * One participant's custom share for `exact-amount` / `percentage` split
 * modes. An omitted/empty `contactID` denotes the payer's own share. Ignored
 * (and may be omitted) for `equally`, which the backend computes itself.
 */
export interface ISplitShare {
  readonly contactID?: string;
  /** Decimal string, e.g. "35.00" — required for `exact-amount`. */
  readonly amount?: string;
  /** Decimal string, e.g. "33.34" — required for `percentage`. */
  readonly percent?: string;
}

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

export interface ICreateSplitTransfer {
  readonly id: string;
  readonly contactID: string;
  readonly amount: number;
}

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
export type SplitShareStatus = 'settled' | 'outstanding';

export interface ISplitParticipant {
  readonly contactID?: string;
  readonly userID?: string;
  readonly name: string;
  readonly share: number;
  readonly isPayer?: boolean;
  readonly status: SplitShareStatus;
}

export interface ISplit {
  readonly id: string;
  readonly title?: string;
  readonly currency: CurrencyCode;
  readonly amount: number;
  readonly status: string;
  readonly participants: ISplitParticipant[];
}

export interface ISplitListItem {
  readonly id: string;
  readonly title?: string;
  readonly amount: number;
  readonly currency: CurrencyCode;
  readonly status: string;
  readonly membersCount: number;
}
