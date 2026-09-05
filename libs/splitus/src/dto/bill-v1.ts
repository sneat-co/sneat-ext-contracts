/** The first stable Splitus bill browser/host wire contract. */
export const SPLITUS_BILL_CONTRACT_VERSION = 1 as const;

export const MAX_SPLITUS_BILL_PARTICIPANTS = 256;
export const MAX_SPLITUS_BILL_LIST_PAGE_SIZE = 100;
export const MAX_SPLITUS_BILL_OBLIGATIONS = 256;
export const MAX_SPLITUS_OBLIGATION_IDS_PER_LINE = 256;

export type SplitusBillContractVersion =
  typeof SPLITUS_BILL_CONTRACT_VERSION;

/**
 * A canonical, non-negative major-unit amount with exactly two fraction
 * digits. Examples: `0.00`, `30.00`, `90.00`.
 *
 * The alias documents the wire type; use `parseExactDecimalString()` at an
 * untrusted boundary. JavaScript numbers and minor-unit numbers are not part
 * of this contract.
 */
export type ExactDecimalString = string;

/** Currencies whose minor unit is exactly two decimal digits in this contract. */
export type SplitusCurrencyCode = 'EUR' | 'GBP' | 'USD';

export type SplitusBillKind = 'general' | 'utility';
export type SplitusUtilityKind =
  | 'electricity'
  | 'gas'
  | 'water'
  | 'internet'
  | 'other';

export type SplitusBillPostingStatus =
  | 'pending'
  | 'posting'
  | 'applied'
  | 'attention';

export type SplitusDebtusSettlementStatus =
  | 'unsettled'
  | 'part_settled'
  | 'settled';

export type SplitusExpectedActualComparison =
  | 'not_available'
  | 'matches'
  | 'increased'
  | 'decreased';

export type SplitusBillAttentionCode =
  | 'authorization_changed'
  | 'source_conflict'
  | 'provider_rejected'
  | 'invalid_provider_receipt'
  | 'operator_action_required';

export interface ISplitusBillAllocationV1 {
  /** Stable within the bill revision, independent of array ordering. */
  readonly allocationID: string;
  /**
   * Contactus identity in this Space. Hosts resolve its real display data from
   * Contactus; this opaque ID is never a user-facing label.
   */
  readonly contactID: string;
  readonly amount: ExactDecimalString;
}

export interface ISplitusBillingPeriodV1 {
  /** Inclusive ISO calendar date. */
  readonly startDate: string;
  /** Inclusive ISO calendar date; must not precede `startDate`. */
  readonly endDate: string;
}

export interface ISplitusUtilityDetailsV1 {
  readonly utilityKind: SplitusUtilityKind;
  readonly providerName?: string;
  readonly period: ISplitusBillingPeriodV1;
}

/**
 * A reference to the Calendarius occurrence represented by this actual bill.
 * `expectedAmount` and `standingChargeAmount` are context only; neither is a
 * paid expense. `actualAmount` on the enclosing bill remains mandatory.
 */
export interface ISplitusRecurringOccurrenceV1 {
  readonly happeningID: string;
  readonly occurrenceID: string;
  readonly expectedAmount?: ExactDecimalString;
  readonly standingChargeAmount?: ExactDecimalString;
  readonly expectedComparison: SplitusExpectedActualComparison;
  readonly previousComparable?: ISplitusPreviousComparableBillV1;
}

export interface ISplitusPreviousComparableBillV1 {
  readonly billID: string;
  readonly actualAmount: ExactDecimalString;
  readonly comparison: Exclude<
    SplitusExpectedActualComparison,
    'not_available'
  >;
}

export interface ICreateSplitusBillV1Request {
  readonly contractVersion: SplitusBillContractVersion;
  readonly spaceID: string;
  /**
   * Stable across duplicate submission and lost-response retries. Reusing the
   * same ID with changed paid/owed allocations is a provider conflict.
   */
  readonly billID: string;
  /**
   * Audit identity of the authenticated actor. The server must bind this to
   * its trusted authentication context. It never implies a paid or owed
   * allocation.
   */
  readonly recorderUserID: string;
  readonly title?: string;
  readonly billKind: SplitusBillKind;
  readonly currency: SplitusCurrencyCode;
  /** Actual paid amount. An expectation is never accepted in its place. */
  readonly actualAmount: ExactDecimalString;
  /** Explicit sources of payment; one contact may also owe a share. */
  readonly paidAllocations: readonly ISplitusBillAllocationV1[];
  /** Explicit responsibility shares. */
  readonly owedAllocations: readonly ISplitusBillAllocationV1[];
  /** Required exactly when `billKind` is `utility`. */
  readonly utility?: ISplitusUtilityDetailsV1;
  readonly recurringOccurrence?: ISplitusRecurringOccurrenceV1;
}

export interface ISplitusDebtusReceiptLineV1 {
  readonly lineID: string;
  readonly obligationIDs: readonly string[];
}

export interface ISplitusBillPostingReceiptV1 {
  readonly receiptID: string;
  readonly operationKey: string;
  readonly inputDigest: string;
  readonly revision: string;
  readonly obligationLines: readonly ISplitusDebtusReceiptLineV1[];
}

export interface ISplitusBillPostingV1 {
  readonly status: SplitusBillPostingStatus;
  /** Durable identity of the retry-safe provider operation. */
  readonly operationKey: string;
  /** Digest of the exact accepted source revision and allocations. */
  readonly inputDigest: string;
  /** Present exactly when `status` is `applied`. */
  readonly receipt?: ISplitusBillPostingReceiptV1;
  /** Present exactly when `status` is `attention`. */
  readonly attentionCode?: SplitusBillAttentionCode;
}

/**
 * An application-relative target. The host resolves it through its injected
 * Debtus navigation adapter; contracts never hard-code debtus.app or another
 * deployment origin.
 */
export interface ISplitusDebtusSettlementTargetV1 {
  readonly route: 'debtus.source-obligations';
  readonly spaceID: string;
  readonly sourceNamespace: 'splitus';
  readonly sourceRecordID: string;
  readonly lineID?: string;
}

export interface ISplitusDebtusObligationV1 {
  readonly lineID: string;
  readonly obligationIDs: readonly string[];
  readonly debtorContactID: string;
  readonly creditorContactID: string;
  readonly currency: SplitusCurrencyCode;
  readonly principalAmount: ExactDecimalString;
  readonly outstandingAmount: ExactDecimalString;
  readonly repaidAmount: ExactDecimalString;
  readonly creditAmount: ExactDecimalString;
  readonly status: SplitusDebtusSettlementStatus;
  readonly settlementTarget: ISplitusDebtusSettlementTargetV1;
}

export interface ISplitusDebtusStatusV1 {
  /** Current state read from Debtus, never a Splitus-maintained balance. */
  readonly status: SplitusDebtusSettlementStatus;
  readonly obligations: readonly ISplitusDebtusObligationV1[];
  readonly settlementTarget: ISplitusDebtusSettlementTargetV1;
}

export interface ISplitusBillV1 extends ICreateSplitusBillV1Request {
  /** Canonical positive decimal integer encoded as a string. */
  readonly revision: string;
  readonly posting: ISplitusBillPostingV1;
  /** Absent until a Debtus financial projection is available. */
  readonly debtus?: ISplitusDebtusStatusV1;
  readonly createdAt: string;
  readonly updatedAt: string;
}

export interface ICreateSplitusBillV1Response {
  readonly contractVersion: SplitusBillContractVersion;
  readonly bill: ISplitusBillV1;
}

export interface IGetSplitusBillV1Request {
  readonly contractVersion: SplitusBillContractVersion;
  readonly spaceID: string;
  readonly billID: string;
}

export interface IGetSplitusBillV1Response {
  readonly contractVersion: SplitusBillContractVersion;
  readonly bill: ISplitusBillV1;
}

export interface ISplitusBillListItemV1 {
  readonly contractVersion: SplitusBillContractVersion;
  readonly spaceID: string;
  readonly billID: string;
  readonly title?: string;
  readonly billKind: SplitusBillKind;
  readonly utilityKind?: SplitusUtilityKind;
  readonly period?: ISplitusBillingPeriodV1;
  readonly currency: SplitusCurrencyCode;
  readonly actualAmount: ExactDecimalString;
  readonly ownPaidAmount: ExactDecimalString;
  readonly ownOwedAmount: ExactDecimalString;
  readonly postingStatus: SplitusBillPostingStatus;
  /** Debtus-derived and absent before the bill has a financial projection. */
  readonly debtusSettlementStatus?: SplitusDebtusSettlementStatus;
  readonly createdAt: string;
}

export interface IListSplitusBillsV1Request {
  readonly contractVersion: SplitusBillContractVersion;
  readonly spaceID: string;
  readonly pageSize: number;
  readonly cursor?: string;
  readonly utilityKind?: SplitusUtilityKind;
  readonly period?: ISplitusBillingPeriodV1;
}

export interface IListSplitusBillsV1Response {
  readonly contractVersion: SplitusBillContractVersion;
  /** Echoes the accepted request bound so callers can verify the page. */
  readonly pageSize: number;
  /** Contains no more than `pageSize` items and never more than 100. */
  readonly items: readonly ISplitusBillListItemV1[];
  readonly nextCursor?: string;
}

const MAX_EXACT_MINOR_UNITS = 9_223_372_036_854_775_807n;
const exactDecimalPattern = /^(?:0|[1-9][0-9]*)\.[0-9]{2}$/;
const datePattern = /^([0-9]{4})-([0-9]{2})-([0-9]{2})$/;

export function parseExactDecimalString(value: unknown): ExactDecimalString {
  if (typeof value !== 'string' || !exactDecimalPattern.test(value)) {
    throw new TypeError(
      'amount must be a canonical decimal string with exactly two fraction digits',
    );
  }
  const minorUnits = BigInt(value.replace('.', ''));
  if (minorUnits > MAX_EXACT_MINOR_UNITS) {
    throw new RangeError('amount exceeds the Splitus contract limit');
  }
  return value;
}

export function parseCreateSplitusBillV1Request(
  value: unknown,
): ICreateSplitusBillV1Request {
  const input = record(value, 'create bill request');
  if (input['contractVersion'] !== SPLITUS_BILL_CONTRACT_VERSION) {
    throw new TypeError('unsupported Splitus bill contract version');
  }
  const billKind = enumValue(
    input['billKind'],
    ['general', 'utility'] as const,
    'billKind',
  );
  const actualAmount = positiveAmount(input['actualAmount'], 'actualAmount');
  const spaceID = storageID(input['spaceID'], 'spaceID');
  const billID = storageID(input['billID'], 'billID');
  const recorderUserID = storageID(
    input['recorderUserID'],
    'recorderUserID',
  );
  const request: ICreateSplitusBillV1Request = {
    contractVersion: SPLITUS_BILL_CONTRACT_VERSION,
    spaceID,
    billID,
    recorderUserID,
    title: optionalText(input['title'], 'title', 256),
    billKind,
    currency: currency(input['currency']),
    actualAmount,
    paidAllocations: allocations(input['paidAllocations'], 'paidAllocations'),
    owedAllocations: allocations(input['owedAllocations'], 'owedAllocations'),
    utility:
      input['utility'] === undefined ? undefined : utility(input['utility']),
    recurringOccurrence:
      input['recurringOccurrence'] === undefined
        ? undefined
        : recurring(input['recurringOccurrence'], actualAmount, billID),
  };
  if ((billKind === 'utility') !== (request.utility !== undefined)) {
    throw new TypeError('utility details are required only for utility bills');
  }
  validateAllocationTotals(request);
  return request;
}

/**
 * Executable host-boundary check for the request's audit claim. The trusted
 * authenticated identity is supplied by the host, never derived from the
 * request itself.
 */
export function assertCreateSplitusBillV1Recorder(
  request: ICreateSplitusBillV1Request,
  authenticatedUserID: string,
): void {
  const trustedUserID = storageID(
    authenticatedUserID,
    'authenticatedUserID',
  );
  if (request.recorderUserID !== trustedUserID) {
    throw new TypeError(
      'recorderUserID must match the trusted authenticated identity',
    );
  }
}

export function parseListSplitusBillsV1Request(
  value: unknown,
): IListSplitusBillsV1Request {
  const input = record(value, 'list bills request');
  if (input['contractVersion'] !== SPLITUS_BILL_CONTRACT_VERSION) {
    throw new TypeError('unsupported Splitus bill contract version');
  }
  const pageSize = input['pageSize'];
  if (
    typeof pageSize !== 'number' ||
    !Number.isSafeInteger(pageSize) ||
    pageSize < 1 ||
    pageSize > MAX_SPLITUS_BILL_LIST_PAGE_SIZE
  ) {
    throw new RangeError(
      `pageSize must be an integer from 1 to ${MAX_SPLITUS_BILL_LIST_PAGE_SIZE}`,
    );
  }
  return {
    contractVersion: SPLITUS_BILL_CONTRACT_VERSION,
    spaceID: storageID(input['spaceID'], 'spaceID'),
    pageSize,
    cursor: optionalText(input['cursor'], 'cursor', 2048),
    utilityKind:
      input['utilityKind'] === undefined
        ? undefined
        : enumValue(
            input['utilityKind'],
            ['electricity', 'gas', 'water', 'internet', 'other'] as const,
            'utilityKind',
          ),
    period:
      input['period'] === undefined ? undefined : period(input['period']),
  };
}

export function parseGetSplitusBillV1Request(
  value: unknown,
): IGetSplitusBillV1Request {
  const input = versionedRecord(value, 'get bill request');
  return {
    contractVersion: SPLITUS_BILL_CONTRACT_VERSION,
    spaceID: storageID(input['spaceID'], 'spaceID'),
    billID: storageID(input['billID'], 'billID'),
  };
}

export function parseCreateSplitusBillV1Response(
  value: unknown,
): ICreateSplitusBillV1Response {
  const input = versionedRecord(value, 'create bill response');
  return {
    contractVersion: SPLITUS_BILL_CONTRACT_VERSION,
    bill: parseSplitusBillV1(input['bill']),
  };
}

export function parseGetSplitusBillV1Response(
  value: unknown,
): IGetSplitusBillV1Response {
  const input = versionedRecord(value, 'get bill response');
  return {
    contractVersion: SPLITUS_BILL_CONTRACT_VERSION,
    bill: parseSplitusBillV1(input['bill']),
  };
}

export function parseListSplitusBillsV1Response(
  value: unknown,
  requestedPageSize?: number,
): IListSplitusBillsV1Response {
  const input = versionedRecord(value, 'list bills response');
  const pageSize = boundedPageSize(input['pageSize']);
  if (requestedPageSize !== undefined && pageSize !== requestedPageSize) {
    throw new RangeError('response pageSize does not match the accepted request');
  }
  if (!Array.isArray(input['items']) || input['items'].length > pageSize) {
    throw new RangeError('response items exceed the bounded pageSize');
  }
  return {
    contractVersion: SPLITUS_BILL_CONTRACT_VERSION,
    pageSize,
    items: input['items'].map((item) => splitusBillListItem(item)),
    nextCursor: optionalText(input['nextCursor'], 'nextCursor', 2048),
  };
}

export function parseSplitusBillV1(value: unknown): ISplitusBillV1 {
  const input = record(value, 'bill');
  const request = parseCreateSplitusBillV1Request(input);
  const revisionValue = positiveIntegerString(input['revision'], 'revision');
  const createdAt = timestamp(input['createdAt'], 'createdAt');
  const updatedAt = timestamp(input['updatedAt'], 'updatedAt');
  if (Date.parse(updatedAt) < Date.parse(createdAt)) {
    throw new RangeError('updatedAt must not precede createdAt');
  }
  const postingValue = posting(input['posting'], revisionValue);
  const debtusValue =
    input['debtus'] === undefined
      ? undefined
      : debtusStatus(
          input['debtus'],
          request.spaceID,
          request.billID,
          request.currency,
        );
  if (debtusValue !== undefined && postingValue.status !== 'applied') {
    throw new TypeError('Debtus status requires applied posting');
  }
  if (debtusValue !== undefined) {
    const receipt = postingValue.receipt;
    if (receipt === undefined) {
      throw new TypeError('Debtus status requires an applied posting receipt');
    }
    validateDebtusMatchesReceipt(receipt, debtusValue);
  }
  return {
    ...request,
    revision: revisionValue,
    posting: postingValue,
    debtus: debtusValue,
    createdAt,
    updatedAt,
  };
}

function splitusBillListItem(value: unknown): ISplitusBillListItemV1 {
  const input = versionedRecord(value, 'bill list item');
  const billKind = enumValue(
    input['billKind'],
    ['general', 'utility'] as const,
    'billKind',
  );
  const utilityKind =
    input['utilityKind'] === undefined
      ? undefined
      : enumValue(
          input['utilityKind'],
          ['electricity', 'gas', 'water', 'internet', 'other'] as const,
          'utilityKind',
        );
  const periodValue =
    input['period'] === undefined ? undefined : period(input['period']);
  if (
    (billKind === 'utility') !==
    (utilityKind !== undefined && periodValue !== undefined)
  ) {
    throw new TypeError('utility list items require utilityKind and period');
  }
  const postingStatus = enumValue(
    input['postingStatus'],
    ['pending', 'posting', 'applied', 'attention'] as const,
    'postingStatus',
  );
  const debtusSettlementStatus =
    input['debtusSettlementStatus'] === undefined
      ? undefined
      : settlementStatus(
          input['debtusSettlementStatus'],
          'debtusSettlementStatus',
        );
  if (
    debtusSettlementStatus !== undefined &&
    postingStatus !== 'applied'
  ) {
    throw new TypeError('Debtus settlement status requires applied posting');
  }
  return {
    contractVersion: SPLITUS_BILL_CONTRACT_VERSION,
    spaceID: storageID(input['spaceID'], 'spaceID'),
    billID: storageID(input['billID'], 'billID'),
    title: optionalText(input['title'], 'title', 256),
    billKind,
    utilityKind,
    period: periodValue,
    currency: currency(input['currency']),
    actualAmount: positiveAmount(input['actualAmount'], 'actualAmount'),
    ownPaidAmount: nonNegativeAmount(input['ownPaidAmount'], 'ownPaidAmount'),
    ownOwedAmount: nonNegativeAmount(input['ownOwedAmount'], 'ownOwedAmount'),
    postingStatus,
    debtusSettlementStatus,
    createdAt: timestamp(input['createdAt'], 'createdAt'),
  };
}

function posting(value: unknown, billRevision: string): ISplitusBillPostingV1 {
  const input = record(value, 'posting');
  const status = enumValue(
    input['status'],
    ['pending', 'posting', 'applied', 'attention'] as const,
    'posting.status',
  );
  const operationKey = storageID(input['operationKey'], 'posting.operationKey');
  const inputDigest = digest(input['inputDigest'], 'posting.inputDigest');
  const receiptValue =
    input['receipt'] === undefined
      ? undefined
      : postingReceipt(
          input['receipt'],
          billRevision,
          operationKey,
          inputDigest,
        );
  const attentionCode =
    input['attentionCode'] === undefined
      ? undefined
      : enumValue(
          input['attentionCode'],
          [
            'authorization_changed',
            'source_conflict',
            'provider_rejected',
            'invalid_provider_receipt',
            'operator_action_required',
          ] as const,
          'posting.attentionCode',
        );
  if ((status === 'applied') !== (receiptValue !== undefined)) {
    throw new TypeError('applied posting status requires exactly one receipt');
  }
  if ((status === 'attention') !== (attentionCode !== undefined)) {
    throw new TypeError('attention posting status requires exactly one code');
  }
  return {
    status,
    operationKey,
    inputDigest,
    receipt: receiptValue,
    attentionCode,
  };
}

function postingReceipt(
  value: unknown,
  billRevision: string,
  operationKey: string,
  inputDigest: string,
): ISplitusBillPostingReceiptV1 {
  const input = record(value, 'posting.receipt');
  const revisionValue = positiveIntegerString(
    input['revision'],
    'posting.receipt.revision',
  );
  const receiptOperationKey = storageID(
    input['operationKey'],
    'posting.receipt.operationKey',
  );
  const receiptInputDigest = digest(
    input['inputDigest'],
    'posting.receipt.inputDigest',
  );
  if (
    revisionValue !== billRevision ||
    receiptOperationKey !== operationKey ||
    receiptInputDigest !== inputDigest
  ) {
    throw new TypeError('posting receipt does not match the accepted bill revision');
  }
  if (
    !Array.isArray(input['obligationLines']) ||
    input['obligationLines'].length > MAX_SPLITUS_BILL_OBLIGATIONS
  ) {
    throw new RangeError('posting receipt obligationLines are unbounded');
  }
  const lineIDs = new Set<string>();
  const obligationIDs = new Set<string>();
  const obligationLines = input['obligationLines'].map((value, index) => {
    const line = record(value, `posting.receipt.obligationLines[${index}]`);
    const lineID = storageID(
      line['lineID'],
      `posting.receipt.obligationLines[${index}].lineID`,
    );
    if (lineIDs.has(lineID)) {
      throw new TypeError('posting receipt repeats an obligation line');
    }
    lineIDs.add(lineID);
    const ids = identifierArray(
      line['obligationIDs'],
      `posting.receipt.obligationLines[${index}].obligationIDs`,
      MAX_SPLITUS_OBLIGATION_IDS_PER_LINE,
    );
    for (const id of ids) {
      if (obligationIDs.has(id)) {
        throw new TypeError('posting receipt repeats an obligation ID');
      }
      obligationIDs.add(id);
    }
    return { lineID, obligationIDs: ids };
  });
  return {
    receiptID: storageID(input['receiptID'], 'posting.receipt.receiptID'),
    operationKey: receiptOperationKey,
    inputDigest: receiptInputDigest,
    revision: revisionValue,
    obligationLines,
  };
}

function debtusStatus(
  value: unknown,
  spaceID: string,
  billID: string,
  billCurrency: string,
): ISplitusDebtusStatusV1 {
  const input = record(value, 'debtus');
  if (
    !Array.isArray(input['obligations']) ||
    input['obligations'].length > MAX_SPLITUS_BILL_OBLIGATIONS
  ) {
    throw new RangeError('Debtus obligations are unbounded');
  }
  const lineIDs = new Set<string>();
  const obligations = input['obligations'].map((item, index) => {
    const obligation = record(item, `debtus.obligations[${index}]`);
    const lineID = storageID(
      obligation['lineID'],
      `debtus.obligations[${index}].lineID`,
    );
    if (lineIDs.has(lineID)) {
      throw new TypeError('Debtus status repeats an obligation line');
    }
    lineIDs.add(lineID);
    const obligationCurrency = currency(obligation['currency']);
    if (obligationCurrency !== billCurrency) {
      throw new TypeError('Debtus obligation currency differs from the bill');
    }
    return {
      lineID,
      obligationIDs: identifierArray(
        obligation['obligationIDs'],
        `debtus.obligations[${index}].obligationIDs`,
        MAX_SPLITUS_OBLIGATION_IDS_PER_LINE,
      ),
      debtorContactID: storageID(
        obligation['debtorContactID'],
        `debtus.obligations[${index}].debtorContactID`,
      ),
      creditorContactID: storageID(
        obligation['creditorContactID'],
        `debtus.obligations[${index}].creditorContactID`,
      ),
      currency: obligationCurrency,
      principalAmount: positiveAmount(
        obligation['principalAmount'],
        `debtus.obligations[${index}].principalAmount`,
      ),
      outstandingAmount: nonNegativeAmount(
        obligation['outstandingAmount'],
        `debtus.obligations[${index}].outstandingAmount`,
      ),
      repaidAmount: nonNegativeAmount(
        obligation['repaidAmount'],
        `debtus.obligations[${index}].repaidAmount`,
      ),
      creditAmount: nonNegativeAmount(
        obligation['creditAmount'],
        `debtus.obligations[${index}].creditAmount`,
      ),
      status: settlementStatus(
        obligation['status'],
        `debtus.obligations[${index}].status`,
      ),
      settlementTarget: settlementTarget(
        obligation['settlementTarget'],
        spaceID,
        billID,
        lineID,
      ),
    };
  });
  return {
    status: settlementStatus(input['status'], 'debtus.status'),
    obligations,
    settlementTarget: settlementTarget(
      input['settlementTarget'],
      spaceID,
      billID,
    ),
  };
}

function validateDebtusMatchesReceipt(
  receipt: ISplitusBillPostingReceiptV1,
  debtus: ISplitusDebtusStatusV1,
): void {
  if (receipt.obligationLines.length !== debtus.obligations.length) {
    throw new TypeError(
      'Debtus obligations do not exactly match the applied posting receipt',
    );
  }
  const receiptLines = new Map(
    receipt.obligationLines.map((line) => [
      line.lineID,
      new Set(line.obligationIDs),
    ]),
  );
  for (const obligation of debtus.obligations) {
    const receiptIDs = receiptLines.get(obligation.lineID);
    if (
      receiptIDs === undefined ||
      receiptIDs.size !== obligation.obligationIDs.length ||
      obligation.obligationIDs.some((id) => !receiptIDs.has(id))
    ) {
      throw new TypeError(
        'Debtus obligations do not exactly match the applied posting receipt',
      );
    }
  }
}

function settlementTarget(
  value: unknown,
  spaceID: string,
  billID: string,
  expectedLineID?: string,
): ISplitusDebtusSettlementTargetV1 {
  const input = record(value, 'settlementTarget');
  const lineID =
    input['lineID'] === undefined
      ? undefined
      : storageID(input['lineID'], 'settlementTarget.lineID');
  if (
    input['route'] !== 'debtus.source-obligations' ||
    input['sourceNamespace'] !== 'splitus' ||
    input['spaceID'] !== spaceID ||
    input['sourceRecordID'] !== billID ||
    lineID !== expectedLineID
  ) {
    throw new TypeError('settlement target does not match the Splitus bill source');
  }
  return {
    route: 'debtus.source-obligations',
    spaceID,
    sourceNamespace: 'splitus',
    sourceRecordID: billID,
    lineID,
  };
}

function versionedRecord(value: unknown, name: string): Record<string, unknown> {
  const input = record(value, name);
  if (input['contractVersion'] !== SPLITUS_BILL_CONTRACT_VERSION) {
    throw new TypeError('unsupported Splitus bill contract version');
  }
  return input;
}

function boundedPageSize(value: unknown): number {
  if (
    typeof value !== 'number' ||
    !Number.isSafeInteger(value) ||
    value < 1 ||
    value > MAX_SPLITUS_BILL_LIST_PAGE_SIZE
  ) {
    throw new RangeError(
      `pageSize must be an integer from 1 to ${MAX_SPLITUS_BILL_LIST_PAGE_SIZE}`,
    );
  }
  return value;
}

function settlementStatus(
  value: unknown,
  name: string,
): SplitusDebtusSettlementStatus {
  return enumValue(
    value,
    ['unsettled', 'part_settled', 'settled'] as const,
    name,
  );
}

function identifierArray(value: unknown, name: string, maximum: number): string[] {
  if (!Array.isArray(value) || value.length < 1 || value.length > maximum) {
    throw new RangeError(`${name} must contain 1 to ${maximum} identifiers`);
  }
  const seen = new Set<string>();
  return value.map((item, index) => {
    const id = storageID(item, `${name}[${index}]`);
    if (seen.has(id)) {
      throw new TypeError(`${name} contains a duplicate identifier`);
    }
    seen.add(id);
    return id;
  });
}

function digest(value: unknown, name: string): string {
  if (typeof value !== 'string' || !/^[0-9a-f]{64}$/.test(value)) {
    throw new TypeError(`${name} must be a lowercase SHA-256 digest`);
  }
  return value;
}

function positiveIntegerString(value: unknown, name: string): string {
  if (typeof value !== 'string' || !/^[1-9][0-9]*$/.test(value)) {
    throw new TypeError(`${name} must be a canonical positive integer string`);
  }
  if (BigInt(value) > 18_446_744_073_709_551_615n) {
    throw new RangeError(`${name} exceeds unsigned 64-bit range`);
  }
  return value;
}

function timestamp(value: unknown, name: string): string {
  if (
    typeof value !== 'string' ||
    !/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d+)?(?:Z|[+-]\d{2}:\d{2})$/.test(
      value,
    ) ||
    Number.isNaN(Date.parse(value))
  ) {
    throw new TypeError(`${name} must be an RFC 3339 timestamp`);
  }
  date(value.slice(0, 10), name);
  return value;
}

function nonNegativeAmount(
  value: unknown,
  name: string,
): ExactDecimalString {
  try {
    return parseExactDecimalString(value);
  } catch (error) {
    if (error instanceof RangeError) {
      throw new RangeError(`${name}: ${error.message}`);
    }
    if (error instanceof Error) {
      throw new TypeError(`${name}: ${error.message}`);
    }
    throw error;
  }
}

function validateAllocationTotals(request: ICreateSplitusBillV1Request): void {
  const contacts = new Set<string>();
  for (const allocation of [
    ...request.paidAllocations,
    ...request.owedAllocations,
  ]) {
    contacts.add(allocation.contactID);
  }
  if (contacts.size < 2 || contacts.size > MAX_SPLITUS_BILL_PARTICIPANTS) {
    throw new RangeError(
      `bill must have 2 to ${MAX_SPLITUS_BILL_PARTICIPANTS} participants`,
    );
  }
  const actual = amountToMinorUnits(request.actualAmount);
  if (sumAllocations(request.paidAllocations) !== actual) {
    throw new RangeError('paid allocations must reconcile to actualAmount');
  }
  if (sumAllocations(request.owedAllocations) !== actual) {
    throw new RangeError('owed allocations must reconcile to actualAmount');
  }
}

function allocations(value: unknown, name: string): ISplitusBillAllocationV1[] {
  if (!Array.isArray(value) || value.length < 1 || value.length > MAX_SPLITUS_BILL_PARTICIPANTS) {
    throw new RangeError(
      `${name} must contain 1 to ${MAX_SPLITUS_BILL_PARTICIPANTS} allocations`,
    );
  }
  const allocationIDs = new Set<string>();
  const contactIDs = new Set<string>();
  return value.map((item, index) => {
    const input = record(item, `${name}[${index}]`);
    const allocationID = storageID(
      input['allocationID'],
      `${name}[${index}].allocationID`,
    );
    const contactID = storageID(
      input['contactID'],
      `${name}[${index}].contactID`,
    );
    if (allocationIDs.has(allocationID) || contactIDs.has(contactID)) {
      throw new TypeError(`${name} contains a duplicate allocation or contact`);
    }
    allocationIDs.add(allocationID);
    contactIDs.add(contactID);
    return {
      allocationID,
      contactID,
      amount: positiveAmount(input['amount'], `${name}[${index}].amount`),
    };
  });
}

function sumAllocations(allocationsToSum: readonly ISplitusBillAllocationV1[]): bigint {
  let total = 0n;
  for (const allocation of allocationsToSum) {
    total += amountToMinorUnits(allocation.amount);
    if (total > MAX_EXACT_MINOR_UNITS) {
      throw new RangeError('allocation total exceeds the Splitus contract limit');
    }
  }
  return total;
}

function amountToMinorUnits(value: ExactDecimalString): bigint {
  return BigInt(value.replace('.', ''));
}

function positiveAmount(value: unknown, name: string): ExactDecimalString {
  const amount = parseExactDecimalString(value);
  if (amountToMinorUnits(amount) === 0n) {
    throw new RangeError(`${name} must be positive`);
  }
  return amount;
}

function utility(value: unknown): ISplitusUtilityDetailsV1 {
  const input = record(value, 'utility');
  return {
    utilityKind: enumValue(
      input['utilityKind'],
      ['electricity', 'gas', 'water', 'internet', 'other'] as const,
      'utility.utilityKind',
    ),
    providerName: optionalText(
      input['providerName'],
      'utility.providerName',
      256,
    ),
    period: period(input['period']),
  };
}

function recurring(
  value: unknown,
  actualAmount: ExactDecimalString,
  billID: string,
): ISplitusRecurringOccurrenceV1 {
  const input = record(value, 'recurringOccurrence');
  const expectedAmount =
    input['expectedAmount'] === undefined
      ? undefined
      : positiveAmount(
          input['expectedAmount'],
          'recurringOccurrence.expectedAmount',
        );
  const standingChargeAmount =
    input['standingChargeAmount'] === undefined
      ? undefined
      : positiveAmount(
          input['standingChargeAmount'],
          'recurringOccurrence.standingChargeAmount',
        );
  const expectedComparison = enumValue(
    input['expectedComparison'],
    ['not_available', 'matches', 'increased', 'decreased'] as const,
    'recurringOccurrence.expectedComparison',
  );
  if (
    (expectedAmount === undefined) !==
    (expectedComparison === 'not_available')
  ) {
    throw new TypeError(
      'comparison must be not_available exactly when expectedAmount is absent',
    );
  }
  if (expectedAmount !== undefined) {
    const expectedComparisonForAmounts = compareAmounts(
      actualAmount,
      expectedAmount,
    );
    if (expectedComparison !== expectedComparisonForAmounts) {
      throw new TypeError('comparison does not match expected and actual amounts');
    }
  }
  if (
    standingChargeAmount !== undefined &&
    amountToMinorUnits(standingChargeAmount) > amountToMinorUnits(actualAmount)
  ) {
    throw new RangeError('standing charge cannot exceed actualAmount');
  }
  return {
    happeningID: storageID(
      input['happeningID'],
      'recurringOccurrence.happeningID',
    ),
    occurrenceID: storageID(
      input['occurrenceID'],
      'recurringOccurrence.occurrenceID',
    ),
    expectedAmount,
    standingChargeAmount,
    expectedComparison,
    previousComparable:
      input['previousComparable'] === undefined
        ? undefined
        : previousComparable(
            input['previousComparable'],
            actualAmount,
            billID,
          ),
  };
}

function previousComparable(
  value: unknown,
  actualAmount: ExactDecimalString,
  billID: string,
): ISplitusPreviousComparableBillV1 {
  const input = record(value, 'recurringOccurrence.previousComparable');
  const previousBillID = storageID(
    input['billID'],
    'recurringOccurrence.previousComparable.billID',
  );
  if (previousBillID === billID) {
    throw new TypeError('previous comparable bill must have a different billID');
  }
  const previousActualAmount = positiveAmount(
    input['actualAmount'],
    'recurringOccurrence.previousComparable.actualAmount',
  );
  const comparison = enumValue(
    input['comparison'],
    ['matches', 'increased', 'decreased'] as const,
    'recurringOccurrence.previousComparable.comparison',
  );
  const expectedComparison = compareAmounts(
    actualAmount,
    previousActualAmount,
  );
  if (comparison !== expectedComparison) {
    throw new TypeError(
      'previous comparison does not match current and prior actual amounts',
    );
  }
  return {
    billID: previousBillID,
    actualAmount: previousActualAmount,
    comparison,
  };
}

function compareAmounts(
  current: ExactDecimalString,
  baseline: ExactDecimalString,
): Exclude<SplitusExpectedActualComparison, 'not_available'> {
  const currentMinor = amountToMinorUnits(current);
  const baselineMinor = amountToMinorUnits(baseline);
  return currentMinor === baselineMinor
    ? 'matches'
    : currentMinor > baselineMinor
      ? 'increased'
      : 'decreased';
}

function period(value: unknown): ISplitusBillingPeriodV1 {
  const input = record(value, 'period');
  const startDate = date(input['startDate'], 'period.startDate');
  const endDate = date(input['endDate'], 'period.endDate');
  if (endDate < startDate) {
    throw new RangeError('period.endDate must not precede period.startDate');
  }
  return { startDate, endDate };
}

function date(value: unknown, name: string): string {
  if (typeof value !== 'string') {
    throw new TypeError(`${name} must be an ISO calendar date`);
  }
  const match = datePattern.exec(value);
  if (match === null) {
    throw new TypeError(`${name} must be an ISO calendar date`);
  }
  const parsed = new Date(`${value}T00:00:00.000Z`);
  if (Number.isNaN(parsed.valueOf()) || parsed.toISOString().slice(0, 10) !== value) {
    throw new TypeError(`${name} must be a real ISO calendar date`);
  }
  return value;
}

function currency(value: unknown): SplitusCurrencyCode {
  return enumValue(value, ['EUR', 'GBP', 'USD'] as const, 'currency');
}

function storageID(value: unknown, name: string): string {
  if (
    typeof value !== 'string' ||
    value.length === 0 ||
    value.trim() !== value ||
    new TextEncoder().encode(value).length > 512 ||
    value.includes('/') ||
    hasControlCharacters(value) ||
    value === '.' ||
    value === '..' ||
    (/^__/.test(value) && /__$/.test(value))
  ) {
    throw new TypeError(`${name} is not a safe identifier`);
  }
  return value;
}

function optionalText(
  value: unknown,
  name: string,
  maxBytes: number,
): string | undefined {
  if (value === undefined) {
    return undefined;
  }
  if (
    typeof value !== 'string' ||
    value.length === 0 ||
    value.trim() !== value ||
    new TextEncoder().encode(value).length > maxBytes ||
    hasControlCharacters(value)
  ) {
    throw new TypeError(`${name} is empty, padded, too long, or contains controls`);
  }
  return value;
}

function hasControlCharacters(value: string): boolean {
  for (const character of value) {
    const code = character.charCodeAt(0);
    if (code < 0x20 || code === 0x7f) {
      return true;
    }
  }
  return false;
}

function record(value: unknown, name: string): Record<string, unknown> {
  if (typeof value !== 'object' || value === null || Array.isArray(value)) {
    throw new TypeError(`${name} must be an object`);
  }
  return value as Record<string, unknown>;
}

function enumValue<const T extends readonly string[]>(
  value: unknown,
  allowed: T,
  name: string,
): T[number] {
  if (typeof value !== 'string' || !allowed.includes(value)) {
    throw new TypeError(`${name} has an unsupported value`);
  }
  return value as T[number];
}
