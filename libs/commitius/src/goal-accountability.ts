/** Supported explicit combinations. Legacy absence is decoded by the server;
 * it must never be inferred from a zero or missing amount in a new request.
 * Optional non-money referee authority has a separate future contract.
 */
export type CommitmentAccountability =
  | {
      accountabilityMode: 'financial';
      verificationMode: 'referee';
    }
  | {
      accountabilityMode: 'none';
      verificationMode: 'self-report';
    };


/** Filtered response provenance. An owner declaration is never referee
 * verification, a payment receipt or a certificate. ISO instants stay strings. */
export interface ISelfReportedCompletion {
  readonly basis: 'self-report';
  readonly reportedAt: string;
  readonly termsVersion: number;
}

export type ICommitmentAccountabilityView =
  | {
      readonly accountabilityMode: 'financial';
      readonly verificationMode: 'referee';
      readonly completion?: never;
    }
  | {
      readonly accountabilityMode: 'none';
      readonly verificationMode: 'self-report';
      readonly completion?: ISelfReportedCompletion;
    };
