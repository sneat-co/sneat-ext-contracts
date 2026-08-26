// The exact set of lifecycle event types recorded in a Listing's per-listing
// transaction log, per REQ transaction-log.
export const YARDIUS_EVENT_TYPES = [
	'publish',
	'request',
	'withdraw',
	'approve',
	'decline',
	'handover',
	'return',
	'claim',
	'transfer',
	'cancel',
] as const;

export type YardiusEventType = (typeof YARDIUS_EVENT_TYPES)[number];

// One entry in a Listing's transaction log, per REQ transaction-log: every
// lifecycle event is appended with its event type, timestamp, and acting
// member. The log is APPEND-ONLY — existing entries MUST NOT be mutated or
// removed by any normal operation (which is also why every field here is
// readonly).
//
// Persistence (backend concern, not part of this DTO's shape): the log lives
// with its Listing under `/spaces/{spaceID}/ext/yardius/...`
// (REQ persistence-convention).
export interface ITransactionLogEntry {
	readonly event: YardiusEventType;

	// ISO timestamp string, server-assigned when the event is appended.
	readonly at: string;

	// The acting member: their Space plus their member (contactus contact) ID
	// within that Space. For owner-side events (publish, approve, decline,
	// handover, return, cancel) this is a member of the owning Space; for
	// requester-side events (request, withdraw, claim) it is a member of the
	// requesting/claiming Space.
	readonly actingMemberSpaceID: string;
	readonly actingMemberID: string;
}
