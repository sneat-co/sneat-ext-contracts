/** The first stable Debtus source-obligation repayment browser contract. */
export const DEBTUS_SOURCE_REPAYMENT_CONTRACT_VERSION = 1 as const;

export const MAX_DEBTUS_SOURCE_OBLIGATION_IDS = 256;

export type DebtusSourceRepaymentContractVersion =
  typeof DEBTUS_SOURCE_REPAYMENT_CONTRACT_VERSION;

/**
 * A canonical non-negative integer number of minor units encoded as a string.
 * Use `parseDebtusExactMinorAmountString()` at every untrusted boundary.
 */
export type DebtusExactMinorAmountString = string;

/** A UTC RFC 3339 timestamp with exactly millisecond precision. */
export type DebtusUtcTimestampString = string;

export interface IDebtusSourceRefV1 {
  readonly namespace: string;
  readonly spaceID: string;
  readonly recordID: string;
}

export interface IDebtusSourceContactRefV1 {
  readonly spaceID: string;
  readonly contactID: string;
}

export type DebtusSourceSettlementStatus =
  | 'unsettled'
  | 'part_settled'
  | 'settled';

export type DebtusSourceRepaymentUnavailableReason =
  | 'readOnlyRole'
  | 'notParty'
  | 'notOutstanding'
  | 'ledgerPartyUnresolved'
  | 'linkedEventLimitReached';

/**
 * Server-authoritative render-time state. It is never an authorization grant;
 * the server reauthorizes the actor and target for every mutation.
 */
export type IDebtusSourceRepaymentCapabilityV1 =
  | {
      readonly canRecordRepayment: true;
      readonly repaymentUnavailableReason?: never;
    }
  | {
      readonly canRecordRepayment: false;
      readonly repaymentUnavailableReason: DebtusSourceRepaymentUnavailableReason;
    };

/** Public exact-string read projection reused by repayment results. */
export interface IDebtusSourceObligationV1 {
  readonly lineID: string;
  readonly obligationIDs: readonly string[];
  readonly debtor: IDebtusSourceContactRefV1;
  readonly creditor: IDebtusSourceContactRefV1;
  readonly currency: string;
  readonly principalMinor: DebtusExactMinorAmountString;
  readonly outstandingMinor: DebtusExactMinorAmountString;
  readonly repaidMinor: DebtusExactMinorAmountString;
  readonly creditMinor: DebtusExactMinorAmountString;
  readonly status: DebtusSourceSettlementStatus;
  readonly repaymentCapability: IDebtusSourceRepaymentCapabilityV1;
}

export type DebtusSourceObligationActivityKind =
  | 'obligation'
  | 'repayment'
  | 'adjustment'
  | 'cancellation'
  | 'credit';

/** Public activity projection; trusted actor identity is intentionally absent. */
export interface IDebtusSourceObligationActivityV1 {
  readonly activityID: string;
  readonly rootActivityID: string;
  readonly lineIDs: readonly string[];
  readonly kind: DebtusSourceObligationActivityKind;
  readonly from: IDebtusSourceContactRefV1;
  readonly to: IDebtusSourceContactRefV1;
  readonly currency: string;
  readonly amountMinor: DebtusExactMinorAmountString;
  readonly occurredAt: DebtusUtcTimestampString;
}

export interface IDebtusSourceRepaymentActivityV1
  extends IDebtusSourceObligationActivityV1 {
  readonly kind: 'repayment';
}

export interface IRecordDebtusSourceRepaymentV1Request {
  readonly contractVersion: DebtusSourceRepaymentContractVersion;
  readonly source: IDebtusSourceRefV1;
  readonly lineID: string;
  readonly obligationID: string;
  readonly currency: string;
  readonly amountMinor: DebtusExactMinorAmountString;
  readonly repaidAt: DebtusUtcTimestampString;
  readonly operationKey: string;
}

export interface IRecordDebtusSourceRepaymentV1Response {
  readonly contractVersion: DebtusSourceRepaymentContractVersion;
  readonly source: IDebtusSourceRefV1;
  readonly operationKey: string;
  readonly obligation: IDebtusSourceObligationV1;
  readonly repayment: IDebtusSourceRepaymentActivityV1;
}

const maxExactMinorUnits = 9_223_372_036_854_775_807n;
const exactMinorPattern = /^(?:0|[1-9][0-9]*)$/;
const utcTimestampPattern =
  /^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}Z$/;

export function parseDebtusExactMinorAmountString(
  value: unknown,
): DebtusExactMinorAmountString {
  if (
    typeof value !== 'string' ||
    value.length > 19 ||
    !exactMinorPattern.test(value)
  ) {
    throw new TypeError(
      'amountMinor must be a canonical non-negative integer string',
    );
  }
  if (BigInt(value) > maxExactMinorUnits) {
    throw new RangeError('amountMinor exceeds signed 64-bit range');
  }
  return value;
}

export function parseRecordDebtusSourceRepaymentV1Request(
  value: unknown,
): IRecordDebtusSourceRepaymentV1Request {
  const input = exactRecord(
    value,
    'source repayment request',
    [
      'contractVersion',
      'source',
      'lineID',
      'obligationID',
      'currency',
      'amountMinor',
      'repaidAt',
      'operationKey',
    ] as const,
  );
  if (input['contractVersion'] !== DEBTUS_SOURCE_REPAYMENT_CONTRACT_VERSION) {
    throw new TypeError('unsupported Debtus source repayment contract version');
  }
  const amountMinor = parseDebtusExactMinorAmountString(input['amountMinor']);
  if (amountMinor === '0') {
    throw new RangeError('amountMinor must be positive');
  }
  return {
    contractVersion: DEBTUS_SOURCE_REPAYMENT_CONTRACT_VERSION,
    source: sourceRef(input['source']),
    lineID: storageID(input['lineID'], 'lineID'),
    obligationID: storageID(input['obligationID'], 'obligationID'),
    currency: currency(input['currency']),
    amountMinor,
    repaidAt: utcTimestamp(input['repaidAt'], 'repaidAt'),
    operationKey: storageID(input['operationKey'], 'operationKey'),
  };
}

export function parseRecordDebtusSourceRepaymentV1Response(
  value: unknown,
  request: IRecordDebtusSourceRepaymentV1Request,
): IRecordDebtusSourceRepaymentV1Response {
  const expected = parseRecordDebtusSourceRepaymentV1Request(request);
  const input = exactRecord(
    value,
    'source repayment response',
    ['contractVersion', 'source', 'operationKey', 'obligation', 'repayment'] as const,
  );
  if (input['contractVersion'] !== DEBTUS_SOURCE_REPAYMENT_CONTRACT_VERSION) {
    throw new TypeError('unsupported Debtus source repayment contract version');
  }
  const source = sourceRef(input['source']);
  assertSource(source, expected.source);
  const operationKey = storageID(input['operationKey'], 'operationKey');
  if (operationKey !== expected.operationKey) {
    throw new TypeError('repayment response operationKey does not match request');
  }
  const obligation = sourceObligation(input['obligation'], source.spaceID);
  if (
    obligation.lineID !== expected.lineID ||
    obligation.currency !== expected.currency ||
    !obligation.obligationIDs.includes(expected.obligationID)
  ) {
    throw new TypeError(
      'repayment response obligation does not match requested source line',
    );
  }
  const repayment = sourceRepaymentActivity(
    input['repayment'],
    source.spaceID,
  );
  if (
    repayment.rootActivityID !== expected.obligationID ||
    repayment.lineIDs.length !== 1 ||
    repayment.lineIDs[0] !== expected.lineID ||
    repayment.currency !== expected.currency ||
    repayment.amountMinor !== expected.amountMinor ||
    repayment.occurredAt !== expected.repaidAt ||
    repayment.from.spaceID !== obligation.debtor.spaceID ||
    repayment.from.contactID !== obligation.debtor.contactID ||
    repayment.to.spaceID !== obligation.creditor.spaceID ||
    repayment.to.contactID !== obligation.creditor.contactID
  ) {
    throw new TypeError(
      'repayment activity does not match the accepted source repayment',
    );
  }
  return {
    contractVersion: DEBTUS_SOURCE_REPAYMENT_CONTRACT_VERSION,
    source,
    operationKey,
    obligation,
    repayment,
  };
}

function sourceRef(value: unknown): IDebtusSourceRefV1 {
  const input = exactRecord(value, 'source', [
    'namespace',
    'spaceID',
    'recordID',
  ] as const);
  return {
    namespace: token(input['namespace'], 'source.namespace'),
    spaceID: storageID(input['spaceID'], 'source.spaceID'),
    recordID: storageID(input['recordID'], 'source.recordID'),
  };
}

function sourceContact(
  value: unknown,
  name: string,
  sourceSpaceID: string,
): IDebtusSourceContactRefV1 {
  const input = exactRecord(value, name, ['spaceID', 'contactID'] as const);
  const spaceID = storageID(input['spaceID'], `${name}.spaceID`);
  if (spaceID !== sourceSpaceID) {
    throw new TypeError(`${name}.spaceID does not match the source Space`);
  }
  return {
    spaceID,
    contactID: storageID(input['contactID'], `${name}.contactID`),
  };
}

function sourceObligation(
  value: unknown,
  sourceSpaceID: string,
): IDebtusSourceObligationV1 {
  const input = exactRecord(
    value,
    'obligation',
    [
      'lineID',
      'obligationIDs',
      'debtor',
      'creditor',
      'currency',
      'principalMinor',
      'outstandingMinor',
      'repaidMinor',
      'creditMinor',
      'status',
      'repaymentCapability',
    ] as const,
  );
  const debtor = sourceContact(input['debtor'], 'obligation.debtor', sourceSpaceID);
  const creditor = sourceContact(
    input['creditor'],
    'obligation.creditor',
    sourceSpaceID,
  );
  if (debtor.contactID === creditor.contactID) {
    throw new TypeError('obligation debtor and creditor must differ');
  }
  return {
    lineID: storageID(input['lineID'], 'obligation.lineID'),
    obligationIDs: identifierArray(
      input['obligationIDs'],
      'obligation.obligationIDs',
      MAX_DEBTUS_SOURCE_OBLIGATION_IDS,
    ),
    debtor,
    creditor,
    currency: currency(input['currency']),
    principalMinor: parseDebtusExactMinorAmountString(input['principalMinor']),
    outstandingMinor: parseDebtusExactMinorAmountString(
      input['outstandingMinor'],
    ),
    repaidMinor: parseDebtusExactMinorAmountString(input['repaidMinor']),
    creditMinor: parseDebtusExactMinorAmountString(input['creditMinor']),
    status: enumValue(
      input['status'],
      ['unsettled', 'part_settled', 'settled'] as const,
      'obligation.status',
    ),
    repaymentCapability: repaymentCapability(input['repaymentCapability']),
  };
}

function sourceRepaymentActivity(
  value: unknown,
  sourceSpaceID: string,
): IDebtusSourceRepaymentActivityV1 {
  const input = exactRecord(
    value,
    'repayment',
    [
      'activityID',
      'rootActivityID',
      'lineIDs',
      'kind',
      'from',
      'to',
      'currency',
      'amountMinor',
      'occurredAt',
    ] as const,
  );
  if (input['kind'] !== 'repayment') {
    throw new TypeError('repayment.kind must be repayment');
  }
  const from = sourceContact(input['from'], 'repayment.from', sourceSpaceID);
  const to = sourceContact(input['to'], 'repayment.to', sourceSpaceID);
  if (from.contactID === to.contactID) {
    throw new TypeError('repayment from and to contacts must differ');
  }
  return {
    activityID: storageID(input['activityID'], 'repayment.activityID'),
    rootActivityID: storageID(
      input['rootActivityID'],
      'repayment.rootActivityID',
    ),
    lineIDs: identifierArray(
      input['lineIDs'],
      'repayment.lineIDs',
      MAX_DEBTUS_SOURCE_OBLIGATION_IDS,
    ),
    kind: 'repayment',
    from,
    to,
    currency: currency(input['currency']),
    amountMinor: positiveMinorAmount(input['amountMinor']),
    occurredAt: utcTimestamp(input['occurredAt'], 'repayment.occurredAt'),
  };
}

function repaymentCapability(
  value: unknown,
): IDebtusSourceRepaymentCapabilityV1 {
  const input = record(value, 'repaymentCapability');
  if (input['canRecordRepayment'] === true) {
    assertKeys(input, 'repaymentCapability', ['canRecordRepayment']);
    return { canRecordRepayment: true };
  }
  if (input['canRecordRepayment'] !== false) {
    throw new TypeError('repaymentCapability.canRecordRepayment must be boolean');
  }
  assertKeys(input, 'repaymentCapability', [
    'canRecordRepayment',
    'repaymentUnavailableReason',
  ]);
  return {
    canRecordRepayment: false,
    repaymentUnavailableReason: enumValue(
      input['repaymentUnavailableReason'],
      [
        'readOnlyRole',
        'notParty',
        'notOutstanding',
        'ledgerPartyUnresolved',
        'linkedEventLimitReached',
      ] as const,
      'repaymentCapability.repaymentUnavailableReason',
    ),
  };
}

function positiveMinorAmount(value: unknown): DebtusExactMinorAmountString {
  const amount = parseDebtusExactMinorAmountString(value);
  if (amount === '0') {
    throw new RangeError('amountMinor must be positive');
  }
  return amount;
}

function storageID(value: unknown, name: string): string {
  if (
    typeof value !== 'string' ||
    value.length === 0 ||
    value.trim() !== value ||
    new TextEncoder().encode(value).length > 512 ||
    Array.from(value).some((character) => {
      const code = character.codePointAt(0) ?? 0;
      return code < 0x20 || code === 0x7f;
    })
  ) {
    throw new TypeError(`${name} is not a valid bounded identifier`);
  }
  return value;
}

function token(value: unknown, name: string): string {
  const result = storageID(value, name);
  if (!/^[a-z0-9_-]+$/.test(result)) {
    throw new TypeError(
      `${name} must contain lowercase ASCII letters, digits, hyphens, or underscores`,
    );
  }
  return result;
}

function currency(value: unknown): string {
  if (typeof value !== 'string' || !/^[A-Z]{3}$/.test(value)) {
    throw new TypeError('currency must be three uppercase ASCII letters');
  }
  return value;
}

function utcTimestamp(value: unknown, name: string): DebtusUtcTimestampString {
  if (
    typeof value !== 'string' ||
    !utcTimestampPattern.test(value) ||
    Number.isNaN(Date.parse(value))
  ) {
    throw new TypeError(
      `${name} must be a valid UTC RFC 3339 timestamp with millisecond precision`,
    );
  }
  const instant = new Date(value);
  const fields = value.match(
    /^(\d{4})-(\d{2})-(\d{2})T(\d{2}):(\d{2}):(\d{2})/,
  );
  if (
    fields === null ||
    instant.getUTCFullYear() !== Number(fields[1]) ||
    instant.getUTCMonth() + 1 !== Number(fields[2]) ||
    instant.getUTCDate() !== Number(fields[3]) ||
    instant.getUTCHours() !== Number(fields[4]) ||
    instant.getUTCMinutes() !== Number(fields[5]) ||
    instant.getUTCSeconds() !== Number(fields[6])
  ) {
    throw new TypeError(
      `${name} must be a valid UTC RFC 3339 timestamp with millisecond precision`,
    );
  }
  return value;
}

function identifierArray(
  value: unknown,
  name: string,
  maximum: number,
): readonly string[] {
  if (!Array.isArray(value) || value.length < 1 || value.length > maximum) {
    throw new RangeError(`${name} must contain 1 to ${maximum} identifiers`);
  }
  const result = value.map((item, index) =>
    storageID(item, `${name}[${index}]`),
  );
  if (new Set(result).size !== result.length) {
    throw new TypeError(`${name} contains duplicate identifiers`);
  }
  return result;
}

function assertSource(
  actual: IDebtusSourceRefV1,
  expected: IDebtusSourceRefV1,
): void {
  if (
    actual.namespace !== expected.namespace ||
    actual.spaceID !== expected.spaceID ||
    actual.recordID !== expected.recordID
  ) {
    throw new TypeError('repayment response source does not match request');
  }
}

function enumValue<const T extends readonly string[]>(
  value: unknown,
  allowed: T,
  name: string,
): T[number] {
  if (typeof value !== 'string' || !allowed.includes(value as T[number])) {
    throw new TypeError(`${name} has an unsupported value`);
  }
  return value;
}

function exactRecord<const K extends readonly string[]>(
  value: unknown,
  name: string,
  keys: K,
): Record<K[number], unknown> {
  const input = record(value, name);
  assertKeys(input, name, keys);
  return input as Record<K[number], unknown>;
}

function record(value: unknown, name: string): Record<string, unknown> {
  if (value === null || typeof value !== 'object' || Array.isArray(value)) {
    throw new TypeError(`${name} must be an object`);
  }
  return value as Record<string, unknown>;
}

function assertKeys(
  value: Record<string, unknown>,
  name: string,
  keys: readonly string[],
): void {
  const allowed = new Set(keys);
  const unexpected = Object.keys(value).find((key) => !allowed.has(key));
  if (unexpected !== undefined) {
    throw new TypeError(`${name} contains unsupported field ${unexpected}`);
  }
}
