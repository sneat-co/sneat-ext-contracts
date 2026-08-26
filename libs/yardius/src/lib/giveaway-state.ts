import { canTransition, StateTransitionTable } from './state-transitions.js';

// Give-away-listing lifecycle states, per REQ giveaway-lifecycle (spec
// display names `Available → Claimed → Transferred → Closed`, persisted as
// the wire values below):
//   - 'available'   — open; eligible viewers claim (claims arriving do NOT
//     change the state — only the owner's selection does);
//   - 'claimed'     — a member of the owning Space selected exactly one
//     claimant; other claims automatically declined;
//   - 'transferred' — the Assetus ownership transfer to the claimant's Space
//     succeeded (Assetus appends its own `Transferred` history event); if the
//     transfer FAILS the Listing stays 'claimed' — staying put is not a
//     transition;
//   - 'closed'      — terminal; immediately follows 'transferred' (or an
//     owner-cancel — see REQ owner-cancel).
export const GIVEAWAY_STATES = ['available', 'claimed', 'transferred', 'closed'] as const;

export type GiveawayState = (typeof GIVEAWAY_STATES)[number];

// The legal give-away transitions, per REQ giveaway-lifecycle +
// REQ owner-cancel. Any transition not listed here MUST be rejected.
//
//   available   → claimed      owner selects exactly one claimant
//   available   → closed       owner cancels (before 'transferred' — legal)
//   claimed     → transferred  Assetus ownership transfer succeeded on
//                              confirmed handover (on failure the Listing
//                              remains 'claimed' — no transition)
//   claimed     → closed       owner cancels; the claim is declined
//   transferred → closed       immediate, automatic
//   closed      → (terminal)
export const GIVEAWAY_STATE_TRANSITIONS: StateTransitionTable<GiveawayState> = {
	available: ['claimed', 'closed'],
	claimed: ['transferred', 'closed'],
	transferred: ['closed'],
	closed: [],
};

// True when `from → to` is a legal give-away-lifecycle transition.
export function canTransitionGiveaway(from: GiveawayState, to: GiveawayState): boolean {
	return canTransition(GIVEAWAY_STATE_TRANSITIONS, from, to);
}
