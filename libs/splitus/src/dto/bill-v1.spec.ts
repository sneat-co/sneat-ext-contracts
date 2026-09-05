import { describe, expect, it } from 'vitest';
import {
  assertCreateSplitusBillV1Recorder,
  MAX_SPLITUS_BILL_LIST_PAGE_SIZE,
  parseCreateSplitusBillV1Request,
  parseCreateSplitusBillV1Response,
  parseExactDecimalString,
  parseGetSplitusBillV1Request,
  parseListSplitusBillsV1Request,
  parseListSplitusBillsV1Response,
  SPLITUS_BILL_CONTRACT_VERSION,
} from './bill-v1';

function electricityBill(): Record<string, unknown> {
  return {
    contractVersion: SPLITUS_BILL_CONTRACT_VERSION,
    spaceID: 'housemates-space',
    billID: 'electricity-2026-08',
    recorderUserID: 'recorder-user',
    title: 'August electricity',
    billKind: 'utility',
    currency: 'EUR',
    actualAmount: '90.00',
    paidAllocations: [
      { allocationID: 'paid-alex', contactID: 'alex-contact', amount: '90.00' },
    ],
    owedAllocations: [
      { allocationID: 'owed-alex', contactID: 'alex-contact', amount: '30.00' },
      { allocationID: 'owed-bea', contactID: 'bea-contact', amount: '30.00' },
      { allocationID: 'owed-cam', contactID: 'cam-contact', amount: '30.00' },
    ],
    utility: {
      utilityKind: 'electricity',
      providerName: 'Grid Energy',
      period: { startDate: '2026-08-01', endDate: '2026-08-31' },
    },
    recurringOccurrence: {
      happeningID: 'monthly-electricity',
      occurrenceID: '2026-08',
      expectedAmount: '80.00',
      standingChargeAmount: '15.00',
      expectedComparison: 'increased',
      previousComparable: {
        billID: 'electricity-2026-07',
        actualAmount: '85.00',
        comparison: 'increased',
      },
    },
  };
}

function billResponse(): Record<string, unknown> {
  const bill = electricityBill();
  const digest = 'a'.repeat(64);
  Object.assign(bill, {
    revision: '1',
    posting: {
      status: 'applied',
      operationKey: 'post-electricity-2026-08',
      inputDigest: digest,
      receipt: {
        receiptID: 'debtus-receipt-1',
        operationKey: 'post-electricity-2026-08',
        inputDigest: digest,
        revision: '1',
        obligationLines: [
          { lineID: 'bea-to-alex', obligationIDs: ['debt-bea-alex'] },
          { lineID: 'cam-to-alex', obligationIDs: ['debt-cam-alex'] },
        ],
      },
    },
    debtus: {
      status: 'unsettled',
      obligations: [
        {
          lineID: 'bea-to-alex',
          obligationIDs: ['debt-bea-alex'],
          debtorContactID: 'bea-contact',
          creditorContactID: 'alex-contact',
          currency: 'EUR',
          principalAmount: '30.00',
          outstandingAmount: '30.00',
          repaidAmount: '0.00',
          creditAmount: '0.00',
          status: 'unsettled',
          settlementTarget: {
            route: 'debtus.source-obligations',
            spaceID: 'housemates-space',
            sourceNamespace: 'splitus',
            sourceRecordID: 'electricity-2026-08',
            lineID: 'bea-to-alex',
          },
        },
      ],
      settlementTarget: {
        route: 'debtus.source-obligations',
        spaceID: 'housemates-space',
        sourceNamespace: 'splitus',
        sourceRecordID: 'electricity-2026-08',
      },
    },
    createdAt: '2026-09-05T12:00:00Z',
    updatedAt: '2026-09-05T12:00:01Z',
  });
  return { contractVersion: SPLITUS_BILL_CONTRACT_VERSION, bill };
}

describe('Splitus bill contract version 1', () => {
  it('accepts an actual EUR 90 bill shared equally by three housemates', () => {
    const request = parseCreateSplitusBillV1Request(electricityBill());

    expect(request.actualAmount).toBe('90.00');
    expect(request.owedAllocations.map(({ amount }) => amount)).toEqual([
      '30.00',
      '30.00',
      '30.00',
    ]);
    expect(request.recurringOccurrence).toMatchObject({
      expectedAmount: '80.00',
      expectedComparison: 'increased',
    });
  });

  it('reconciles explicit custom paid and owed allocations exactly', () => {
    const input = electricityBill();
    input['paidAllocations'] = [
      { allocationID: 'paid-alex', contactID: 'alex-contact', amount: '70.00' },
      { allocationID: 'paid-bea', contactID: 'bea-contact', amount: '20.00' },
    ];
    input['owedAllocations'] = [
      { allocationID: 'owed-alex', contactID: 'alex-contact', amount: '10.00' },
      { allocationID: 'owed-bea', contactID: 'bea-contact', amount: '20.00' },
      { allocationID: 'owed-cam', contactID: 'cam-contact', amount: '60.00' },
    ];

    const request = parseCreateSplitusBillV1Request(input);
    expect(request.paidAllocations).toHaveLength(2);
    expect(request.owedAllocations[2]?.amount).toBe('60.00');
  });

  it('keeps the authenticated recorder distinct from every payment party', () => {
    const request = parseCreateSplitusBillV1Request(electricityBill());

    expect(request.recorderUserID).toBe('recorder-user');
    expect(request.paidAllocations.map(({ contactID }) => contactID)).toEqual([
      'alex-contact',
    ]);
    expect(request.recorderUserID).not.toBe(
      request.paidAllocations[0]?.contactID,
    );
    expect(() =>
      assertCreateSplitusBillV1Recorder(request, 'impersonated-user'),
    ).toThrow('trusted authenticated identity');
    expect(() =>
      assertCreateSplitusBillV1Recorder(request, 'recorder-user'),
    ).not.toThrow();
  });

  it.each([
    ['spaceID', 'spaces/foreign'],
    ['spaceID', '..'],
    ['billID', '__reserved__'],
    ['billID', ' bill-padded'],
    ['recorderUserID', 'user\u0000id'],
  ])('rejects unsafe %s %j', (field, value) => {
    const input = electricityBill();
    input[field] = value;
    expect(() => parseCreateSplitusBillV1Request(input)).toThrow(
      'safe identifier',
    );
  });

  it('rejects numeric JSON amounts instead of rounding them in JavaScript', () => {
    const input = JSON.parse(JSON.stringify(electricityBill())) as Record<
      string,
      unknown
    >;
    input['actualAmount'] = 90;

    expect(() => parseCreateSplitusBillV1Request(input)).toThrow(
      'canonical decimal string',
    );
  });

  it('rejects a canonical-looking amount above signed 64-bit minor units', () => {
    expect(() => parseExactDecimalString('92233720368547758.08')).toThrow(
      'contract limit',
    );
  });

  it('rejects an expected occurrence when the actual paid amount is absent', () => {
    const input = electricityBill();
    delete input['actualAmount'];

    expect(() => parseCreateSplitusBillV1Request(input)).toThrow(
      'canonical decimal string',
    );
  });

  it('rejects allocations that do not reconcile to the actual amount', () => {
    const input = electricityBill();
    input['owedAllocations'] = [
      { allocationID: 'owed-alex', contactID: 'alex-contact', amount: '30.00' },
      { allocationID: 'owed-bea', contactID: 'bea-contact', amount: '30.00' },
      { allocationID: 'owed-cam', contactID: 'cam-contact', amount: '29.99' },
    ];

    expect(() => parseCreateSplitusBillV1Request(input)).toThrow(
      'owed allocations must reconcile',
    );
  });

  it('binds the declared expected comparison to the actual amount', () => {
    const input = electricityBill();
    input['recurringOccurrence'] = {
      happeningID: 'monthly-electricity',
      occurrenceID: '2026-08',
      expectedAmount: '100.00',
      expectedComparison: 'increased',
    };

    expect(() => parseCreateSplitusBillV1Request(input)).toThrow(
      'does not match expected and actual',
    );
  });

  it('binds comparison to the previous comparable actual bill', () => {
    const input = electricityBill();
    const recurring = input['recurringOccurrence'] as Record<string, unknown>;
    recurring['previousComparable'] = {
      billID: 'electricity-2026-07',
      actualAmount: '95.00',
      comparison: 'increased',
    };

    expect(() => parseCreateSplitusBillV1Request(input)).toThrow(
      'previous comparison does not match',
    );
  });
});

describe('Splitus bill response contract version 1', () => {
  it('validates stable identifiers for a detail request', () => {
    expect(
      parseGetSplitusBillV1Request({
        contractVersion: SPLITUS_BILL_CONTRACT_VERSION,
        spaceID: 'housemates-space',
        billID: 'electricity-2026-08',
      }),
    ).toMatchObject({ billID: 'electricity-2026-08' });
    expect(() =>
      parseGetSplitusBillV1Request({
        contractVersion: SPLITUS_BILL_CONTRACT_VERSION,
        spaceID: 'housemates-space',
        billID: '../unsafe',
      }),
    ).toThrow('safe identifier');
  });

  it('parses a bounded applied bill with a matching Debtus receipt', () => {
    const response = parseCreateSplitusBillV1Response(billResponse());

    expect(response.bill.posting.receipt?.receiptID).toBe('debtus-receipt-1');
    expect(response.bill.debtus?.obligations[0]?.outstandingAmount).toBe(
      '30.00',
    );
  });

  it('rejects a numeric amount from an untrusted API response', () => {
    const response = billResponse();
    const bill = response['bill'] as Record<string, unknown>;
    const debtus = bill['debtus'] as Record<string, unknown>;
    const obligations = debtus['obligations'] as Record<string, unknown>[];
    const firstObligation = obligations[0];
    expect(firstObligation).toBeDefined();
    if (firstObligation === undefined) {
      throw new Error('test response has no obligation');
    }
    firstObligation['outstandingAmount'] = 30;

    expect(() => parseCreateSplitusBillV1Response(response)).toThrow(
      'canonical decimal string',
    );
  });

  it('rejects mismatched durable posting proof', () => {
    const response = billResponse();
    const bill = response['bill'] as Record<string, unknown>;
    const posting = bill['posting'] as Record<string, unknown>;
    const receipt = posting['receipt'] as Record<string, unknown>;
    receipt['revision'] = '2';

    expect(() => parseCreateSplitusBillV1Response(response)).toThrow(
      'does not match the accepted bill revision',
    );
  });

  it('rejects invalid timestamps and unbounded server attention messages', () => {
    const response = billResponse();
    const bill = response['bill'] as Record<string, unknown>;
    bill['updatedAt'] = 'not-a-time';
    expect(() => parseCreateSplitusBillV1Response(response)).toThrow(
      'RFC 3339',
    );

    bill['updatedAt'] = '2026-09-05T12:00:01Z';
    bill['posting'] = {
      status: 'attention',
      operationKey: 'post-electricity-2026-08',
      inputDigest: 'a'.repeat(64),
      attentionCode: 'raw server error text',
    };
    expect(() => parseCreateSplitusBillV1Response(response)).toThrow(
      'unsupported value',
    );
  });
});

describe('Splitus bill list contract version 1', () => {
  it('accepts the maximum bounded page size', () => {
    expect(
      parseListSplitusBillsV1Request({
        contractVersion: SPLITUS_BILL_CONTRACT_VERSION,
        spaceID: 'housemates-space',
        pageSize: MAX_SPLITUS_BILL_LIST_PAGE_SIZE,
      }).pageSize,
    ).toBe(MAX_SPLITUS_BILL_LIST_PAGE_SIZE);
  });

  it('rejects an unbounded page request', () => {
    expect(() =>
      parseListSplitusBillsV1Request({
        contractVersion: SPLITUS_BILL_CONTRACT_VERSION,
        spaceID: 'housemates-space',
        pageSize: MAX_SPLITUS_BILL_LIST_PAGE_SIZE + 1,
      }),
    ).toThrow('pageSize');
  });

  it('rejects a response larger than its echoed request bound', () => {
    const item = {
      contractVersion: SPLITUS_BILL_CONTRACT_VERSION,
      spaceID: 'housemates-space',
      billID: 'electricity-2026-08',
      billKind: 'utility',
      utilityKind: 'electricity',
      period: { startDate: '2026-08-01', endDate: '2026-08-31' },
      currency: 'EUR',
      actualAmount: '90.00',
      ownPaidAmount: '0.00',
      ownOwedAmount: '30.00',
      postingStatus: 'applied',
      debtusSettlementStatus: 'unsettled',
      createdAt: '2026-09-05T12:00:00Z',
    };
    expect(() =>
      parseListSplitusBillsV1Response(
        {
          contractVersion: SPLITUS_BILL_CONTRACT_VERSION,
          pageSize: 1,
          items: [item, { ...item, billID: 'electricity-2026-07' }],
        },
        1,
      ),
    ).toThrow('exceed');
  });
});
