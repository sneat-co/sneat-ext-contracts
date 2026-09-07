// eslint-disable-next-line @nx/enforce-module-boundaries -- this golden file is shared with the public Go module's wire test.
import wholeSecondWireExample from '../../../../../debtus/contract4debtus/testdata/source_repayment_v1.json';
import {
  DEBTUS_SOURCE_DUE_DATE_CONTRACT_VERSION,
  DEBTUS_SOURCE_REPAYMENT_CONTRACT_VERSION,
  parseDebtusExactMinorAmountString,
  parseDebtusSourceObligationDueDateV1,
  parseRecordDebtusSourceRepaymentV1Request,
  parseRecordDebtusSourceRepaymentV1Response,
} from './source-obligation-models';

describe('Debtus source due-date result', () => {
  const result = (updatedAt: string) => ({
    contractVersion: DEBTUS_SOURCE_DUE_DATE_CONTRACT_VERSION,
    source: { namespace: 'splitus', spaceID: 'house-1', recordID: 'bill-1' },
    lineID: 'line-1', revision: 2, dueDate: '2026-09-21', happeningID: 'task-1',
    state: 'active', updatedAt, updatedBy: 'user-1',
  });

  it.each(['2026-09-07T05:56:12Z', '2026-09-07T05:56:12.1Z', '2026-09-07T05:56:12.123Z', '2026-09-07T05:56:12.123456789Z'])(
    'accepts Go RFC3339Nano timestamp %s',
    (updatedAt) => expect(parseDebtusSourceObligationDueDateV1(result(updatedAt)).updatedAt).toBe(updatedAt),
  );

  it.each(['2026-09-07T05:56:12.1234567890Z', '2026-09-07T05:56:12.123+01:00', '2026-02-30T05:56:12Z'])(
    'rejects invalid server timestamp %s',
    (updatedAt) => expect(() => parseDebtusSourceObligationDueDateV1(result(updatedAt))).toThrow(/UTC RFC 3339 timestamp/),
  );
});

const request = () => ({
  contractVersion: DEBTUS_SOURCE_REPAYMENT_CONTRACT_VERSION,
  source: { namespace: 'splitus', spaceID: 'house-1', recordID: 'bill-1' },
  lineID: 'bea-to-alex',
  obligationID: 'obligation-1',
  currency: 'EUR',
  amountMinor: '3000',
  repaidAt: '2026-09-06T08:30:00.000Z',
  operationKey: 'bill-1/repayment-1',
});

const response = () => ({
  contractVersion: DEBTUS_SOURCE_REPAYMENT_CONTRACT_VERSION,
  source: { namespace: 'splitus', spaceID: 'house-1', recordID: 'bill-1' },
  operationKey: 'bill-1/repayment-1',
  obligation: {
    lineID: 'bea-to-alex',
    obligationIDs: ['obligation-1'],
    debtor: { spaceID: 'house-1', contactID: 'bea' },
    creditor: { spaceID: 'house-1', contactID: 'alex' },
    currency: 'EUR',
    principalMinor: '6000',
    outstandingMinor: '3000',
    repaidMinor: '3000',
    creditMinor: '0',
    status: 'part_settled',
    repaymentCapability: { canRecordRepayment: true },
  },
  repayment: {
    activityID: 'repayment-1',
    rootActivityID: 'obligation-1',
    lineIDs: ['bea-to-alex'],
    kind: 'repayment',
    from: { spaceID: 'house-1', contactID: 'bea' },
    to: { spaceID: 'house-1', contactID: 'alex' },
    currency: 'EUR',
    amountMinor: '3000',
    occurredAt: '2026-09-06T08:30:00.000Z',
  },
});

describe('Debtus source repayment request contract', () => {
  it('accepts an exact EUR 30 repayment without browser actor authority', () => {
    expect(parseRecordDebtusSourceRepaymentV1Request(request())).toEqual(
      request(),
    );
    expect('actorUserID' in request()).toBe(false);
  });

  it.each([
    ['numeric JSON amount', 3000],
    ['zero', '0'],
    ['leading zero', '03000'],
    ['negative', '-1'],
    ['fractional', '30.00'],
    ['overflow', '9223372036854775808'],
  ])('rejects %s amountMinor', (_name, amountMinor) => {
    expect(() =>
      parseRecordDebtusSourceRepaymentV1Request({
        ...request(),
        amountMinor,
      }),
    ).toThrow();
  });

  it('rejects malformed identities, non-UTC time, and browser actor claims', () => {
    expect(() =>
      parseRecordDebtusSourceRepaymentV1Request({
        ...request(),
        obligationID: ' obligation-1',
      }),
    ).toThrow(/identifier/);
    expect(() =>
      parseRecordDebtusSourceRepaymentV1Request({
        ...request(),
        repaidAt: '2026-09-06T09:30:00+01:00',
      }),
    ).toThrow(/UTC RFC 3339/);
    expect(() =>
      parseRecordDebtusSourceRepaymentV1Request({
        ...request(),
        repaidAt: '2026-09-06T08:30:00.0001Z',
      }),
    ).toThrow(/millisecond precision/);
    expect(() =>
      parseRecordDebtusSourceRepaymentV1Request({
        ...request(),
        actorUserID: 'bea-user',
      }),
    ).toThrow(/unsupported field actorUserID/);
  });

  it('accepts the maximum signed 64-bit amount and rejects overflow', () => {
    expect(parseDebtusExactMinorAmountString('9223372036854775807')).toBe(
      '9223372036854775807',
    );
    expect(() =>
      parseDebtusExactMinorAmountString('9223372036854775808'),
    ).toThrow(/signed 64-bit/);
    expect(() => parseDebtusExactMinorAmountString('9'.repeat(10_000))).toThrow(
      /canonical non-negative integer string/,
    );
  });
});

describe('Debtus source repayment response contract', () => {
  it('accepts the shared Go/browser whole-second wire example', () => {
    const parsedRequest = parseRecordDebtusSourceRepaymentV1Request(
      wholeSecondWireExample.request,
    );
    expect(
      parseRecordDebtusSourceRepaymentV1Response(
        wholeSecondWireExample.response,
        parsedRequest,
      ).repayment.occurredAt,
    ).toBe('2026-09-06T08:30:00.000Z');
  });

  it('accepts a bound result with exact read amounts and capability', () => {
    expect(
      parseRecordDebtusSourceRepaymentV1Response(response(), request()),
    ).toEqual(response());
  });

  it('accepts a full repayment whose updated capability is unavailable', () => {
    const value = response();
    value.obligation.outstandingMinor = '0';
    value.obligation.repaidMinor = '6000';
    value.obligation.status = 'settled';
    value.obligation.repaymentCapability = {
      canRecordRepayment: false,
      repaymentUnavailableReason: 'notOutstanding',
    } as never;
    expect(
      parseRecordDebtusSourceRepaymentV1Response(value, request()),
    ).toEqual(value);
  });

  it.each([
    ['source', (value: ReturnType<typeof response>) => (value.source.recordID = 'bill-2')],
    ['line', (value: ReturnType<typeof response>) => (value.obligation.lineID = 'line-2')],
    [
      'obligation',
      (value: ReturnType<typeof response>) =>
        (value.obligation.obligationIDs = ['obligation-2']),
    ],
    [
      'activity root',
      (value: ReturnType<typeof response>) =>
        (value.repayment.rootActivityID = 'obligation-2'),
    ],
    [
      'activity amount',
      (value: ReturnType<typeof response>) => (value.repayment.amountMinor = '2999'),
    ],
    [
      'activity timestamp',
      (value: ReturnType<typeof response>) =>
        (value.repayment.occurredAt = '2026-09-06T08:30:00.001Z'),
    ],
    [
      'activity debtor',
      (value: ReturnType<typeof response>) =>
        (value.repayment.from.contactID = 'cam'),
    ],
    [
      'activity creditor',
      (value: ReturnType<typeof response>) =>
        (value.repayment.to.contactID = 'cam'),
    ],
  ])('rejects a valid-looking result for a different %s', (_name, mutate) => {
    const value = response();
    mutate(value);
    expect(() =>
      parseRecordDebtusSourceRepaymentV1Response(value, request()),
    ).toThrow(/does not match/);
  });

  it('rejects numeric read amounts and inconsistent capability shapes', () => {
    const numeric = response();
    numeric.obligation.outstandingMinor = 3000 as never;
    expect(() =>
      parseRecordDebtusSourceRepaymentV1Response(numeric, request()),
    ).toThrow(/canonical non-negative integer string/);

    const unauthorized = response();
    unauthorized.obligation.repaymentCapability = {
      canRecordRepayment: false,
    } as never;
    expect(() =>
      parseRecordDebtusSourceRepaymentV1Response(unauthorized, request()),
    ).toThrow(/unsupported value/);
  });
});
