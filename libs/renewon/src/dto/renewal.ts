import { IWithCreated } from '@sneat/dto';

export type RenewalStatus = 'active' | 'overdue' | 'cancelled' | 'done';

export interface IAmount {
  value: number;
  currency: string;
}

// The display/brief fields of a renewal, shared by brief & dbo.
export interface IRenewalBrief {
  title: string;
  renewalDate: string; // ISO yyyy-mm-dd
  expiryDate?: string;
  amount?: IAmount;
  autoRenew?: boolean;
  status: RenewalStatus;
  paymentRisk?: boolean;
  counterpartyID?: string;
  documentID?: string;
  assetID?: string;
}

// Persisted renewal at spaces/{spaceID}/ext/renewon/renewals/{id}.
export type IRenewalDbo = IRenewalBrief & IWithCreated;

// Request body for POST /v0/api4renewon/spaces/{spaceID}/renewals. The backend
// sets status default, createdAt/createdBy and the id.
export interface ICreateRenewalRequest {
  title: string;
  renewalDate: string;
  expiryDate?: string;
  amount?: IAmount;
  autoRenew?: boolean;
  status?: RenewalStatus;
  paymentRisk?: boolean;
  counterpartyID?: string;
  documentID?: string;
  assetID?: string;
}

// A renewal with its id, for the dashboard.
export interface IRenewalWithID {
  id: string;
  dbo: IRenewalDbo;
}

export type RenewalGroupKey =
  | 'upcoming'
  | 'overdue'
  | 'atRisk'
  | 'other';

export interface IRenewalGroup {
  key: RenewalGroupKey;
  title: string;
  emoji: string;
  renewals: IRenewalWithID[];
}

// groupRenewals buckets renewals for the dashboard, relative to todayISO
// (yyyy-mm-dd). Pure + deterministic so it is unit-testable. Overdue: a
// past renewal/expiry date not cancelled/done. At-risk: paymentRisk flag.
// Upcoming: everything else with a future-or-today date. Other: no usable date.
export function groupRenewals(
  items: readonly IRenewalWithID[],
  todayISO: string,
): IRenewalGroup[] {
  const groups: Record<RenewalGroupKey, IRenewalWithID[]> = {
    upcoming: [],
    overdue: [],
    atRisk: [],
    other: [],
  };
  for (const item of items) {
    const { dbo } = item;
    const date = dbo.expiryDate || dbo.renewalDate;
    if (dbo.paymentRisk && dbo.status !== 'cancelled' && dbo.status !== 'done') {
      groups.atRisk.push(item);
    } else if (!date) {
      groups.other.push(item);
    } else if (date < todayISO && dbo.status !== 'cancelled' && dbo.status !== 'done') {
      groups.overdue.push(item);
    } else {
      groups.upcoming.push(item);
    }
  }
  const byDate = (a: IRenewalWithID, b: IRenewalWithID): number =>
    (a.dbo.expiryDate || a.dbo.renewalDate || '').localeCompare(
      b.dbo.expiryDate || b.dbo.renewalDate || '',
    );
  const meta: Record<RenewalGroupKey, { title: string; emoji: string }> = {
    overdue: { title: 'Overdue / expired', emoji: '⚠️' },
    atRisk: { title: 'Failed / at-risk payments', emoji: '💳' },
    upcoming: { title: 'Upcoming', emoji: '🔔' },
    other: { title: 'No date yet', emoji: '❓' },
  };
  const order: RenewalGroupKey[] = ['overdue', 'atRisk', 'upcoming', 'other'];
  return order
    .filter((k) => groups[k].length > 0)
    .map((k) => ({
      key: k,
      title: meta[k].title,
      emoji: meta[k].emoji,
      renewals: groups[k].sort(byDate),
    }));
}
