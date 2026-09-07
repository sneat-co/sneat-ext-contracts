import { firstValueFrom, of } from 'rxjs';
import { describe, expect, it } from 'vitest';
import {
  IFinancialOccurrenceService,
  parseFinancialOccurrenceQueryResponse,
} from './financial-occurrence-service';

describe('financial occurrence contract', () => {
  it('preserves owner-issued occurrence identity and cancellation uncertainty', async () => {
    const parsed = parseFinancialOccurrenceQueryResponse({
      hasMore: false,
      occurrences: [
        {
          happeningID: 'electricity',
          occurrenceID: 'owner-issued-opaque-id',
          scheduledDate: '2026-09-15',
          periodStartDate: '2026-09-15',
          periodEndDate: '2026-09-15',
          title: 'Electricity',
          priceID: 'electricity-price',
          currency: 'EUR',
          expectedMinor: 10000,
          assetIDs: ['home-1'],
          contactLinks: [{ contactID: 'child-1', roles: ['participant'] }],
          priceRevision: 2,
          pricingAvailability: 'available',
          cancellationFinancialEffect: 'unknown',
        },
      ],
    });
    const service: IFinancialOccurrenceService = {
      listFinancialOccurrences: () => of(parsed),
      listFinancialAgreementCharges: () =>
        of({
          charges: [],
          coverages: [],
          hasMore: false,
          snapshotConsistent: true,
          snapshotDigest: 'empty-source-digest',
        }),
    };

    const result = await firstValueFrom(
      service.listFinancialOccurrences({
        spaceID: 'space1',
        happeningID: 'electricity',
        fromDate: '2026-09-01',
        toDate: '2026-09-30',
      }),
    );

    expect(result.occurrences[0]).toMatchObject({
      occurrenceID: 'owner-issued-opaque-id',
      expectedMinor: 10000,
      cancellationFinancialEffect: 'unknown',
      assetIDs: ['home-1'],
      contactLinks: [{ contactID: 'child-1', roles: ['participant'] }],
    });
  });

  it('rejects unsafe amounts and unknown availability values', () => {
    expect(() =>
      parseFinancialOccurrenceQueryResponse({
        hasMore: false,
        occurrences: [
          {
            happeningID: 'electricity',
            occurrenceID: 'occurrence',
            periodStartDate: '2026-09-01',
            periodEndDate: '2026-09-01',
            title: 'Electricity',
            priceID: 'price',
            currency: 'EUR',
            expectedMinor: Number.MAX_SAFE_INTEGER + 1,
            pricingAvailability: 'available',
          },
        ],
      }),
    ).toThrow(/safe integer/);
    expect(() =>
      parseFinancialOccurrenceQueryResponse({
        hasMore: false,
        occurrences: [
          {
            happeningID: 'electricity',
            occurrenceID: 'occurrence',
            periodStartDate: '2026-09-01',
            periodEndDate: '2026-09-01',
            title: 'Electricity',
            pricingAvailability: 'guessed',
          },
        ],
      }),
    ).toThrow(/pricingAvailability/);
  });

  it('rejects inconsistent completion, dates, and pricing metadata', () => {
    const candidate = {
      happeningID: 'electricity', occurrenceID: 'occurrence', title: 'Electricity',
      periodStartDate: '2026-09-01', periodEndDate: '2026-09-30',
      priceID: 'price', currency: 'EUR', expectedMinor: 12000,
      pricingAvailability: 'available',
    };
    expect(() => parseFinancialOccurrenceQueryResponse({hasMore: true, occurrences: [candidate]})).toThrow(/completion metadata/);
    expect(() => parseFinancialOccurrenceQueryResponse({hasMore: false, occurrences: [{...candidate, periodEndDate: '2026-08-31'}]})).toThrow(/ends before/);
    expect(() => parseFinancialOccurrenceQueryResponse({hasMore: false, occurrences: [{...candidate, periodStartDate: '2026-02-30'}]})).toThrow(/real calendar date/);
    expect(() => parseFinancialOccurrenceQueryResponse({hasMore: false, occurrences: [{...candidate, expectedMinor: undefined}]})).toThrow(/inconsistent pricing/);
    expect(() => parseFinancialOccurrenceQueryResponse({hasMore: false, occurrences: [{...candidate, pricingAvailability: 'unsupported', pricingUnavailableReason: 'unsupported_term'}]})).toThrow(/inconsistent pricing/);
  });
});
