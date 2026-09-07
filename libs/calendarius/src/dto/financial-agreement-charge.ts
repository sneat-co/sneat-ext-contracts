import type {
  IAcceptedPriceTerms,
  IFinancialAttribution,
} from './financial-agreement';
import type { IFinancialEnrollmentScope } from './financial-enrollment';

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
  /** Owner-verified coverage; consumers must not infer it from mutable attributions. */
  readonly enrollmentScope?: IFinancialEnrollmentScope;
  readonly happeningID: string;
  /** Source-owned label safe to display in an authorized reporting Space. */
  readonly title: string;
  /** The accepted agreement repeats; this does not assert a payment due date. */
  readonly regular: boolean;
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
  /** Selected-deal authority independent of charge rows in this window. */
  readonly coverages: readonly IFinancialAgreementCoverageFact[];
  readonly hasMore: boolean;
  readonly nextCursor?: string;
  readonly snapshotConsistent: boolean;
  /** Digest of the full bounded result, stable across every page. */
  readonly snapshotDigest?: string;
  readonly incompleteReason?: string;
}

export interface IFinancialAgreementCoverageFact {
  readonly ownerSpaceID: string;
  readonly reportingSpaceID: string;
  readonly agreementID: string;
  readonly enrollmentID: string;
  readonly enrollmentScope: IFinancialEnrollmentScope;
  readonly happeningID: string;
  readonly state: 'recorded_external' | 'confirmed';
  readonly effectiveFromISO: string;
  readonly effectiveToISO?: string;
  readonly verified: true;
}
