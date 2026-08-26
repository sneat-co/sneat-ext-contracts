import { describe, it, expect } from 'vitest';
import {
	canTransitionGiveaway,
	GIVEAWAY_STATE_TRANSITIONS,
	GIVEAWAY_STATES,
	GiveawayState,
} from './giveaway-state.js';
import { canTransition } from './state-transitions.js';

describe('GIVEAWAY_STATES', () => {
	it('is exactly the REQ giveaway-lifecycle state set, in lifecycle order', () => {
		expect(GIVEAWAY_STATES).toEqual(['available', 'claimed', 'transferred', 'closed']);
	});
});

describe('GIVEAWAY_STATE_TRANSITIONS', () => {
	it('has a row for every give-away state', () => {
		expect(Object.keys(GIVEAWAY_STATE_TRANSITIONS).sort()).toEqual(
			[...GIVEAWAY_STATES].sort(),
		);
	});

	it('walks the full happy path Available → Claimed → Transferred → Closed (AC giveaway-transfers-ownership)', () => {
		const happyPath: readonly GiveawayState[] = [
			'available',
			'claimed',
			'transferred',
			'closed',
		];
		for (let i = 0; i + 1 < happyPath.length; i++) {
			expect(canTransitionGiveaway(happyPath[i], happyPath[i + 1])).toBe(true);
		}
	});

	const legal: ReadonlyArray<[GiveawayState, GiveawayState, string]> = [
		['available', 'claimed', 'owner selects exactly one claimant'],
		['available', 'closed', 'owner cancels before a claimant is selected'],
		['claimed', 'transferred', 'Assetus ownership transfer succeeded'],
		['claimed', 'closed', 'owner cancels; the claim is declined'],
		['transferred', 'closed', 'Closed immediately follows Transferred'],
	];

	it.each(legal)('allows %s → %s (%s)', (from, to) => {
		expect(canTransitionGiveaway(from, to)).toBe(true);
	});

	const illegal: ReadonlyArray<[GiveawayState, GiveawayState, string]> = [
		['available', 'transferred', 'cannot transfer without a selected claimant'],
		['claimed', 'available', 'selection auto-declined the other claims; no way back'],
		['transferred', 'available', 'ownership already moved'],
		['transferred', 'claimed', 'ownership already moved'],
		['closed', 'available', 'closed is terminal'],
		['closed', 'claimed', 'closed is terminal'],
	];

	it.each(illegal)('rejects %s → %s (%s)', (from, to) => {
		expect(canTransitionGiveaway(from, to)).toBe(false);
	});

	it('rejects self-transitions — a failed Assetus transfer keeps the Listing Claimed by NOT transitioning (AC giveaway-transfer-failure-stays-claimed)', () => {
		for (const state of GIVEAWAY_STATES) {
			expect(canTransitionGiveaway(state, state)).toBe(false);
		}
	});

	it("makes 'closed' terminal and 'transferred' exit only to 'closed'", () => {
		expect(GIVEAWAY_STATE_TRANSITIONS.closed).toEqual([]);
		expect(GIVEAWAY_STATE_TRANSITIONS.transferred).toEqual(['closed']);
	});

	it('canTransitionGiveaway agrees with the generic canTransition over the table', () => {
		for (const from of GIVEAWAY_STATES) {
			for (const to of GIVEAWAY_STATES) {
				expect(canTransitionGiveaway(from, to)).toBe(
					canTransition(GIVEAWAY_STATE_TRANSITIONS, from, to),
				);
			}
		}
	});
});
