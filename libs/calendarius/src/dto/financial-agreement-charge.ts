import type {
  IAcceptedPriceTerms,
  IFinancialAttribution,
} from './financial-agreement';

export type FinancialChargeDirection = 'expense' | 'income' | 'transfer';
export type FinancialChargeTemporalBasis =
  | 'occurrence_service_cost'
  | 'contract_period_cost'
  | 'normalized_service_cost';
export type FinancialChargeBillingTiming = 'unknown' | 'known_due';

export interface IFinancialAgreementChargeQuery {
  readonly reportingSpaceID: string;
  readonly fromMonthISO: string;
  readonly months: number;
  readonly pageSize: number;
  readonly cursor?: string;
}

export interface IFinancialAgreementChargeFact {
  readonly ownerSpaceID: string;
  readonly reportingSpaceID: string;
  readonly agreementID: string;
  readonly enrollmentID: string;
  readonly happeningID: string;
  readonly agreementRevision: number;
  readonly termsRevision: number;
  readonly chargeID: string;
  readonly occurrenceID?: string;
  readonly direction?: FinancialChargeDirection;
  readonly amountMinor?: number;
  readonly currency: 'EUR' | 'GBP' | 'USD';
  readonly economicPeriod: { readonly startDate: string; readonly endDate: string };
  /** Economic meaning; only exact occurrence/contract periods can match invoices. */
  readonly temporalBasis: FinancialChargeTemporalBasis;
  /** Payment timing remains unknown unless the source owns a real due date. */
  readonly billingTiming: FinancialChargeBillingTiming;
  readonly dueDate?: string;
  readonly ownerTimezone: string;
  readonly invoiceReconciliationEligible: boolean;
  readonly acceptedPrice: IAcceptedPriceTerms;
  readonly contactAttributions: readonly IFinancialAttribution[];
  readonly assetAttributions: readonly IFinancialAttribution[];
  readonly status: 'available' | 'unavailable';
  readonly diagnostics: readonly string[];
}

export interface IFinancialAgreementChargePage {
  readonly charges: readonly IFinancialAgreementChargeFact[];
  readonly hasMore: boolean;
  readonly nextCursor?: string;
  readonly snapshotConsistent: boolean;
  readonly incompleteReason?: string;
}
