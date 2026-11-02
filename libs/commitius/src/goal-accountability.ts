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
