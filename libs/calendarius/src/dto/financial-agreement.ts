export type FinancialAgreementState =
  | 'offered'
  | 'recorded_external'
  | 'confirmed'
  | 'ended'
  | 'cancelled';

export type FinancialAgreementSide = 'payer' | 'receiver';
export type FinancialPartyKind = 'space' | 'contact';
export type FinancialAgreementTermUnit =
  | 'single'
  | 'second'
  | 'minute'
  | 'hour'
  | 'day'
  | 'week'
  | 'month'
  | 'quarter'
  | 'year';

export interface IFinancialPartyRef {
  readonly kind: FinancialPartyKind;
  readonly spaceID: string;
  readonly contactID?: string;
}

export interface IFinancialPartyShare {
  readonly party: IFinancialPartyRef;
  readonly amountMinor: number;
}

export interface IFinancialPartyWeight {
  readonly party: IFinancialPartyRef;
  readonly shares: number;
}

export interface IFinancialAttribution {
  readonly id: string;
  readonly amountMinor: number;
}

export interface IFinancialAttributionWeight {
  readonly id: string;
  readonly shares: number;
}

export interface IFinancialAgreementTerm {
  readonly unit: FinancialAgreementTermUnit;
  readonly length: number;
}

export interface IAcceptedPriceTerms {
  readonly priceID: string;
  /** Zero is the initial catalog snapshot before the first price edit. */
  readonly priceRevision: number;
  readonly amountMinor: number;
  readonly currency: 'EUR' | 'GBP' | 'USD';
  /** Quantity selected for this agreement, never the whole happening. */
  readonly quantity: number;
  readonly term: IFinancialAgreementTerm;
}

export interface IFinancialConfirmation {
  readonly party: IFinancialPartyRef;
  readonly termsRevision: number;
  readonly actorUserID: string;
  readonly confirmedAt: string;
  readonly revokedAt?: string;
}

/**
 * An owner-issued reporting view of one selected deal. `ownerSpaceID` is part
 * of source identity. `reportingSpaceID` is the authorized financial view.
 */
export interface IFinancialAgreementFact {
  readonly ownerSpaceID: string;
  readonly reportingSpaceID: string;
  readonly agreementID: string;
  readonly enrollmentID: string;
  readonly happeningID: string;
  readonly revision: number;
  readonly termsRevision: number;
  readonly state: FinancialAgreementState;
  readonly reportingSpaceSide: FinancialAgreementSide;
  readonly reportingAmountMinor: number;
  readonly payers: readonly IFinancialPartyShare[];
  readonly receivers: readonly IFinancialPartyShare[];
  readonly contactAttributions: readonly IFinancialAttribution[];
  readonly assetAttributions: readonly IFinancialAttribution[];
  readonly acceptedPrice: IAcceptedPriceTerms;
  readonly effectiveFromISO: string;
  readonly effectiveToISO?: string;
  readonly externallyRecorded?: boolean;
  readonly confirmations: readonly IFinancialConfirmation[];
}

export interface IFinancialAgreementQuery {
  readonly reportingSpaceID: string;
  readonly fromMonthISO: string;
  readonly months: number;
  readonly pageSize: number;
  readonly cursor?: string;
}

export interface IFinancialAgreementListRequest {
  readonly ownerSpaceID: string;
  readonly happeningID: string;
  readonly pageSize: number;
  readonly cursor?: string;
}

export interface IFinancialAgreementReadRequest {
  readonly ownerSpaceID: string;
  readonly agreementID: string;
  readonly historyLimit: number;
  readonly beforeRevision?: number;
}

export interface IFinancialAgreementRevisionView {
  readonly revision: number;
  readonly before?: IFinancialAgreementFact;
  readonly after?: IFinancialAgreementFact;
  readonly actorUserID: string;
  readonly at: string;
  readonly reason?: string;
  readonly operationID: string;
}

export interface IFinancialAgreementReadResponse {
  readonly agreement: IFinancialAgreementFact;
  readonly revisions: readonly IFinancialAgreementRevisionView[];
  readonly historyComplete: boolean;
  readonly nextBeforeRevision?: number;
}

export interface IFinancialAgreementResult {
  readonly agreements: readonly IFinancialAgreementFact[];
  readonly hasMore: boolean;
  readonly nextCursor?: string;
  readonly snapshotConsistent: boolean;
  readonly incompleteReason?:
    | 'query_limit'
    | 'source_collection_mutable_between_pages';
}

export interface ISaveFinancialAgreementRequest {
  readonly ownerSpaceID: string;
  readonly reportingSpaceID: string;
  readonly agreementID: string;
  readonly enrollmentID: string;
  readonly happeningID: string;
  readonly priceID: string;
  readonly expectedPriceRevision: number;
  readonly quantity: number;
  readonly operationID: string;
  readonly expectedRevision: number;
  readonly state: 'offered' | 'recorded_external';
  readonly reportingSpaceSide: FinancialAgreementSide;
  readonly payers: readonly IFinancialPartyWeight[];
  readonly receivers: readonly IFinancialPartyWeight[];
  readonly contactAttributions: readonly IFinancialAttributionWeight[];
  readonly assetAttributions: readonly IFinancialAttributionWeight[];
  readonly effectiveFromISO: string;
  readonly effectiveToISO?: string;
  readonly reason?: string;
}

export interface ISaveFinancialAgreementResponse {
  readonly agreementID: string;
  readonly revision: number;
  readonly updatedAt: string;
}

export interface IConfirmFinancialAgreementRequest {
  readonly reportingSpaceID: string;
  readonly ownerSpaceID: string;
  readonly agreementID: string;
  readonly expectedRevision: number;
  readonly termsRevision: number;
  readonly operationID: string;
  readonly party: IFinancialPartyRef;
}

export interface IRevokeFinancialAgreementGrantRequest
  extends IConfirmFinancialAgreementRequest {
  readonly reason?: string;
}
