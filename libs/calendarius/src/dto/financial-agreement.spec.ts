import { describe, expect, it } from 'vitest';
import {
  IFinancialAgreementFact,
  ISaveFinancialAgreementRequest,
} from './financial-agreement';

describe('financial agreement wire contract', () => {
  it('keeps owner identity, reporting perspective and accepted terms explicit', () => {
    const fact: IFinancialAgreementFact = {
      ownerSpaceID: 'school', reportingSpaceID: 'family', agreementID: 'deal-1',
      enrollmentID: 'child-a', happeningID: 'music', revision: 2, termsRevision: 1,
      state: 'recorded_external', reportingSpaceSide: 'payer', reportingAmountMinor: 4000,
      payers: [{ party: { kind: 'space', spaceID: 'family' }, amountMinor: 4000 }],
      receivers: [{ party: { kind: 'space', spaceID: 'school' }, amountMinor: 4000 }],
      contactAttributions: [{ id: 'child-a', amountMinor: 4000 }], assetAttributions: [],
      acceptedPrice: { priceID: 'monthly', priceRevision: 0, amountMinor: 4000, currency: 'EUR', quantity: 1, term: { unit: 'month', length: 1 } },
      effectiveFromISO: '2026-09-01', externallyRecorded: true, confirmations: [],
    };
    expect(fact.ownerSpaceID).toBe('school');
    expect(fact.acceptedPrice.amountMinor).toBe(4000);
  });

  it('takes share weights rather than client-computed money', () => {
    const request: ISaveFinancialAgreementRequest = {
      ownerSpaceID: 'school', reportingSpaceID: 'family', agreementID: 'deal-1',
      enrollmentID: 'child-a', happeningID: 'music', priceID: 'monthly',
      expectedPriceRevision: 0, quantity: 1, operationID: 'op-1', expectedRevision: 0,
      state: 'recorded_external', reportingSpaceSide: 'payer',
      payers: [{ party: { kind: 'space', spaceID: 'family' }, shares: 1 }],
      receivers: [{ party: { kind: 'space', spaceID: 'school' }, shares: 1 }],
      contactAttributions: [{ id: 'child-a', shares: 1 }], assetAttributions: [],
      effectiveFromISO: '2026-09-01',
    };
    expect(request.payers[0].shares).toBe(1);
    expect('amountMinor' in request.payers[0]).toBe(false);
  });
});
