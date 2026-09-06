import { InjectionToken } from '@angular/core';
import { Observable } from 'rxjs';

export type FinancialPricingAvailability =
  | 'available'
  | 'ambiguous'
  | 'unsupported';

export interface IFinancialOccurrenceQuery {
  readonly spaceID: string;
  /** Omit to discover bounded candidates across this Space. */
  readonly happeningID?: string;
  readonly fromDate: string;
  readonly toDate: string;
  readonly timezone?: string;
}

export interface IFinancialOccurrenceCandidate {
  readonly happeningID: string;
  /** Opaque owner-issued identity. Consumers must never construct this value. */
  readonly occurrenceID: string;
  readonly slotID?: string;
  /** Calendar date; it is not a bill due, payment, posting, or invoice date. */
  readonly scheduledDate?: string;
  /** Economic source coverage; never inferred cash or invoice dates. */
  readonly periodStartDate: string;
  readonly periodEndDate: string;
  readonly title: string;
  readonly priceID?: string;
  readonly priceRevision?: number;
  readonly currency?: 'EUR' | 'GBP' | 'USD';
  readonly expectedMinor?: number;
  /** Same-Space Assetus identities linked by the source Happening. */
  readonly assetIDs?: readonly string[];
  readonly pricingAvailability: FinancialPricingAvailability;
  readonly pricingUnavailableReason?: string;
  /** `unknown` means cancellation is known but a fee waiver is not. */
  readonly cancellationFinancialEffect?: 'unknown';
}

export interface IFinancialOccurrenceQueryResponse {
  readonly occurrences: readonly IFinancialOccurrenceCandidate[];
  /** True when the bounded owner query reached its source or result cap. */
  readonly hasMore: boolean;
  readonly incompleteReason?: 'query_limit';
}

export interface IFinancialOccurrenceService {
  listFinancialOccurrences(
    query: IFinancialOccurrenceQuery,
  ): Observable<IFinancialOccurrenceQueryResponse>;
}

export const FINANCIAL_OCCURRENCE_SERVICE =
  new InjectionToken<IFinancialOccurrenceService>(
    'CalendariusFinancialOccurrenceService',
  );

export function parseFinancialOccurrenceQueryResponse(
  value: unknown,
): IFinancialOccurrenceQueryResponse {
  if (!isRecord(value) || !Array.isArray(value['occurrences'])) {
    throw new TypeError('financial occurrence response must contain occurrences');
  }
  const hasMore = value['hasMore'];
  if (typeof hasMore !== 'boolean') {
    throw new TypeError('financial occurrence response must contain hasMore');
  }
  const incompleteReason = optionalString(value, 'incompleteReason');
  if (incompleteReason !== undefined && incompleteReason !== 'query_limit') {
    throw new TypeError('invalid financial occurrence incompleteReason');
  }
  if (hasMore !== (incompleteReason === 'query_limit')) {
    throw new TypeError('financial occurrence completion metadata is inconsistent');
  }
  return {
    occurrences: value['occurrences'].map(parseCandidate),
    hasMore,
    incompleteReason,
  };
}

function parseCandidate(value: unknown): IFinancialOccurrenceCandidate {
  if (!isRecord(value)) {
    throw new TypeError('financial occurrence must be an object');
  }
  const availability = requiredString(value, 'pricingAvailability');
  if (
    availability !== 'available' &&
    availability !== 'ambiguous' &&
    availability !== 'unsupported'
  ) {
    throw new TypeError('invalid financial occurrence pricingAvailability');
  }
  const currency = optionalString(value, 'currency');
  if (
    currency !== undefined &&
    currency !== 'EUR' &&
    currency !== 'GBP' &&
    currency !== 'USD'
  ) {
    throw new TypeError('invalid financial occurrence currency');
  }
  const effect = optionalString(value, 'cancellationFinancialEffect');
  if (effect !== undefined && effect !== 'unknown') {
    throw new TypeError('invalid cancellation financial effect');
  }
  const scheduledDate = optionalDate(value, 'scheduledDate');
  const periodStartDate = requiredDate(value, 'periodStartDate');
  const periodEndDate = requiredDate(value, 'periodEndDate');
  if (periodEndDate < periodStartDate) {
    throw new TypeError('financial occurrence period ends before it starts');
  }
  if (scheduledDate !== undefined && (scheduledDate < periodStartDate || scheduledDate > periodEndDate)) {
    throw new TypeError('financial occurrence scheduledDate is outside its period');
  }
  const priceID = optionalString(value, 'priceID');
  const expectedMinor = optionalSafeInteger(value, 'expectedMinor');
  const unavailableReason = optionalString(value, 'pricingUnavailableReason');
  if (availability === 'available') {
    if (!priceID || !currency || expectedMinor === undefined || expectedMinor <= 0 || unavailableReason !== undefined) {
      throw new TypeError('available financial occurrence has inconsistent pricing');
    }
  } else if (priceID !== undefined || currency !== undefined || expectedMinor !== undefined || unavailableReason === undefined) {
    throw new TypeError('unavailable financial occurrence has inconsistent pricing');
  }
  const assetIDs = optionalStrings(value, 'assetIDs');
  if (assetIDs !== undefined && new Set(assetIDs).size !== assetIDs.length) {
    throw new TypeError('financial occurrence assetIDs must be unique');
  }
  return {
    happeningID: requiredString(value, 'happeningID'),
    occurrenceID: requiredString(value, 'occurrenceID'),
    slotID: optionalString(value, 'slotID'),
    scheduledDate,
    periodStartDate,
    periodEndDate,
    title: requiredString(value, 'title'),
    priceID,
    priceRevision: optionalSafeInteger(value, 'priceRevision'),
    currency,
    expectedMinor,
    assetIDs,
    pricingAvailability: availability,
    pricingUnavailableReason: unavailableReason,
    cancellationFinancialEffect: effect,
  };
}

function requiredDate(value: Record<string, unknown>, key: string): string {
  const result = requiredString(value, key);
  validateDate(result, key);
  return result;
}

function optionalDate(value: Record<string, unknown>, key: string): string | undefined {
  const result = optionalString(value, key);
  if (result !== undefined) {
    validateDate(result, key);
  }
  return result;
}

function validateDate(value: string, key: string): void {
  if (!/^\d{4}-\d{2}-\d{2}$/.test(value)) {
    throw new TypeError(`financial occurrence ${key} must be an ISO date`);
  }
  const parsed = new Date(`${value}T00:00:00Z`);
  if (Number.isNaN(parsed.valueOf()) || parsed.toISOString().slice(0, 10) !== value) {
    throw new TypeError(`financial occurrence ${key} must be a real calendar date`);
  }
}

function optionalStrings(
  value: Record<string, unknown>,
  key: string,
): readonly string[] | undefined {
  const result = value[key];
  if (result === undefined) {
    return undefined;
  }
  if (
    !Array.isArray(result) ||
    result.some((item) => typeof item !== 'string' || !item)
  ) {
    throw new TypeError(`financial occurrence ${key} must be strings`);
  }
  return result;
}

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === 'object' && value !== null && !Array.isArray(value);
}

function requiredString(value: Record<string, unknown>, key: string): string {
  const result = value[key];
  if (typeof result !== 'string' || !result) {
    throw new TypeError(`financial occurrence ${key} must be a string`);
  }
  return result;
}

function optionalString(
  value: Record<string, unknown>,
  key: string,
): string | undefined {
  const result = value[key];
  if (result === undefined) {
    return undefined;
  }
  if (typeof result !== 'string' || !result) {
    throw new TypeError(`financial occurrence ${key} must be a string`);
  }
  return result;
}

function optionalSafeInteger(
  value: Record<string, unknown>,
  key: string,
): number | undefined {
  const result = value[key];
  if (result === undefined) {
    return undefined;
  }
  if (typeof result !== 'number' || !Number.isSafeInteger(result) || result < 0) {
    throw new TypeError(`financial occurrence ${key} must be a safe integer`);
  }
  return result;
}
