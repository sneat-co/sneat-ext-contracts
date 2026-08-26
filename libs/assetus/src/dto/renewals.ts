// Mirrors backend dto4assetus.Renewals* — the unified upcoming-renewals feed
// (GET /v0/assetus/renewals). A Calendarius sync turns each item into a
// renewal happening (recurring/yearly when period === 'yearly').

// RenewalPeriod: how often a renewal recurs.
export type RenewalPeriod = 'once' | 'yearly';

// IRenewalItem is one upcoming renewal/expiry derived from an asset.
export interface IRenewalItem {
  assetID: string;
  assetName?: string;
  category?: string;
  // insurance/home_insurance/life_insurance/warranty/nct/tax/service/utility/
  // contract/passport/driving_license/…
  kind: string;
  dueOn: string; // ISO 'YYYY-MM-DD'
  period: RenewalPeriod;
}

// IRenewalsResponse — GET /v0/assetus/renewals?spaceID=&withinDays=.
export interface IRenewalsResponse {
  renewals: IRenewalItem[];
}
