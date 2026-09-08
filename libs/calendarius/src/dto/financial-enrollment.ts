export type FinancialEnrollmentScopeKind = 'whole_happening' | 'subjects';
export type FinancialEnrollmentSubjectExtensionID = 'contactus' | 'assetus';

export interface IFinancialEnrollmentSubjectRef {
  readonly extensionID: FinancialEnrollmentSubjectExtensionID;
  readonly entityID: string;
}

export interface IFinancialEnrollmentScope {
  readonly kind: FinancialEnrollmentScopeKind;
  readonly subjects: readonly IFinancialEnrollmentSubjectRef[];
}

export interface IResolveFinancialEnrollmentRequest {
  readonly ownerSpaceID: string;
  readonly happeningID: string;
  readonly operationID: string;
  readonly scope: IFinancialEnrollmentScope;
}

export interface IResolveFinancialEnrollmentResponse {
  readonly enrollmentID: string;
  readonly scope: IFinancialEnrollmentScope;
  readonly identityStatus: 'verified';
  readonly createdAt: string;
}
