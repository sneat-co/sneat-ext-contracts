// The THIN, actually-persisted requoter profile — the Formius reference
// composition written by the Go backend (`facade4requoter.Onboard`) at
// `/spaces/{spaceID}/ext/requoter/profiles/{profileID}`. It holds only reference
// IDs to the canonical records (Assetus car + insurance document, Contactus
// drivers, Calendarius renewal) plus the orphan cover/usage fields no module owns.
//
// This is deliberately NOT the rich `IProfileBase` aggregate (proposer, full
// drivers, address) — that richer shape is the aspirational profile-edit/autofill
// target, not what is stored today. The frontend read path deserializes THIS type
// and composes a display view by re-reading the referenced records.

export interface IStoredProfileBrief {
  /** Assetus vehicle asset id (the car). */
  assetID?: string;
  /** Contactus member ids (the drivers/family). */
  contactIDs?: readonly string[];
  /** Assetus document asset id (the current insurance policy). */
  insuranceDocID?: string;
  /** Calendarius recurring-yearly happening id (the renewal reminder). */
  happeningID?: string;
  /** Orphan fields owned by no canonical module. */
  coverStartDate?: string;
  annualMileage?: number;
  classOfUse?: string;
}

// The stored doc carries created/modified/space/user bases on the backend; the
// read model only needs the reference fields, so Brief and Dbo are the same shape.
export type IStoredProfileDbo = IStoredProfileBrief;
