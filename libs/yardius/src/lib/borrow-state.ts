import { canTransition, StateTransitionTable } from './state-transitions.js';

// Borrow-listing lifecycle states, per REQ borrow-lifecycle (spec display
// names `Available → Requested → Approved → Borrowed → Returned → Closed`,
// persisted as the wire values below):
//   - 'available' — open, no pending requests;
//   - 'requested' — ≥1 request pending; the Listing stays visible and accepts
//     further requests;
//   - 'approved'  — a member of the owning Space approved exactly one request;
//     all other pending requests are automatically declined and no further
//     requests are accepted;
//   - 'borrowed'  — a member of the owning Space confirmed handover;
//   - 'returned'  — a member of the owning Space confirmed the item is back;
//   - 'closed'    — terminal; immediately follows 'returned' (or an
//     owner-cancel — see REQ owner-cancel).
export const BORROW_STATES = [
	'available',
	'requested',
	'approved',
	'borrowed',
	'returned',
	'closed',
] as const;

export type BorrowState = (typeof BORROW_STATES)[number];

// The legal borrow transitions, per REQ borrow-lifecycle + REQ owner-cancel.
// Any transition not listed here MUST be rejected.
//
//   available → requested   first request arrives
//   available → closed      owner cancels (before 'borrowed' — legal)
//   requested → approved    owner approves exactly one request
//   requested → available   the only pending request is withdrawn
//   requested → closed      owner cancels; all pending requests declined
//   approved  → borrowed    owner confirms handover
//   approved  → closed      owner cancels (still before 'borrowed')
//   borrowed  → returned    owner confirms the item is back; a 'borrowed'
//                           Listing MUST NOT be cancellable — 'returned' is
//                           its only exit
//   returned  → closed      immediate, automatic
//   closed    → (terminal)
export const BORROW_STATE_TRANSITIONS: StateTransitionTable<BorrowState> = {
	available: ['requested', 'closed'],
	requested: ['approved', 'available', 'closed'],
	approved: ['borrowed', 'closed'],
	borrowed: ['returned'],
	returned: ['closed'],
	closed: [],
};

// True when `from → to` is a legal borrow-lifecycle transition.
export function canTransitionBorrow(from: BorrowState, to: BorrowState): boolean {
	return canTransition(BORROW_STATE_TRANSITIONS, from, to);
}
