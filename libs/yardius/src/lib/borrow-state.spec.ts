import { describe, it, expect } from 'vitest';
import {
	BORROW_STATE_TRANSITIONS,
	BORROW_STATES,
	BorrowState,
	canTransitionBorrow,
} from './borrow-state.js';
import { canTransition } from './state-transitions.js';

describe('BORROW_STATES', () => {
	it('is exactly the REQ borrow-lifecycle state set, in lifecycle order', () => {
		expect(BORROW_STATES).toEqual([
			'available',
			'requested',
			'approved',
			'borrowed',
			'returned',
			'closed',
		]);
	});
});

describe('BORROW_STATE_TRANSITIONS', () => {
	it('has a row for every borrow state', () => {
		expect(Object.keys(BORROW_STATE_TRANSITIONS).sort()).toEqual(
			[...BORROW_STATES].sort(),
		);
	});

	it('walks the full happy path Available → Requested → Approved → Borrowed → Returned → Closed (AC borrow-full-loop)', () => {
		const happyPath: readonly BorrowState[] = [
			'available',
			'requested',
			'approved',
			'borrowed',
			'returned',
			'closed',
		];
		for (let i = 0; i + 1 < happyPath.length; i++) {
			expect(canTransitionBorrow(happyPath[i], happyPath[i + 1])).toBe(true);
		}
	});

	const legal: ReadonlyArray<[BorrowState, BorrowState, string]> = [
		['available', 'requested', 'first request arrives'],
		['available', 'closed', 'owner cancels before any request'],
		['requested', 'approved', 'owner approves exactly one request'],
		['requested', 'available', 'the only pending request is withdrawn (AC requester-withdraws)'],
		['requested', 'closed', 'owner cancels; pending requests declined (AC owner-cancels-listing)'],
		['approved', 'borrowed', 'owner confirms handover'],
		['approved', 'closed', 'owner cancels before Borrowed'],
		['borrowed', 'returned', 'owner confirms return'],
		['returned', 'closed', 'Closed immediately follows Returned'],
	];

	it.each(legal)('allows %s → %s (%s)', (from, to) => {
		expect(canTransitionBorrow(from, to)).toBe(true);
	});

	const illegal: ReadonlyArray<[BorrowState, BorrowState, string]> = [
		['available', 'approved', 'cannot approve with no request'],
		['available', 'borrowed', 'cannot skip request/approval'],
		['requested', 'borrowed', 'cannot hand over without approval'],
		['approved', 'available', 'approval declined the other requests; no way back'],
		['approved', 'requested', 'no further requests are accepted after approval'],
		['borrowed', 'closed', 'a Borrowed listing MUST NOT be cancellable (AC no-cancel-while-borrowed)'],
		['borrowed', 'available', 'a borrow completes via Returned, never resets'],
		['returned', 'available', 'Returned only closes'],
		['returned', 'borrowed', 'no re-borrow on the same listing'],
		['closed', 'available', 'closed is terminal'],
		['closed', 'requested', 'closed is terminal'],
	];

	it.each(illegal)('rejects %s → %s (%s)', (from, to) => {
		expect(canTransitionBorrow(from, to)).toBe(false);
	});

	it('rejects self-transitions for every state', () => {
		for (const state of BORROW_STATES) {
			expect(canTransitionBorrow(state, state)).toBe(false);
		}
	});

	it("makes 'closed' terminal and 'returned' exit only to 'closed'", () => {
		expect(BORROW_STATE_TRANSITIONS.closed).toEqual([]);
		expect(BORROW_STATE_TRANSITIONS.returned).toEqual(['closed']);
	});

	it("gives 'borrowed' exactly one exit: 'returned'", () => {
		expect(BORROW_STATE_TRANSITIONS.borrowed).toEqual(['returned']);
	});

	it('canTransitionBorrow agrees with the generic canTransition over the table', () => {
		for (const from of BORROW_STATES) {
			for (const to of BORROW_STATES) {
				expect(canTransitionBorrow(from, to)).toBe(
					canTransition(BORROW_STATE_TRANSITIONS, from, to),
				);
			}
		}
	});
});
