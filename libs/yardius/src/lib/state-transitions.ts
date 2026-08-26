// A lifecycle's legal transitions as data: for each state, the exact set of
// states it may move to. Both lifecycles (`borrow-state.ts`,
// `giveaway-state.ts`) publish their table in this shape so frontend and
// backend validate transitions from the same source of truth. Any transition
// not present in the table MUST be rejected.
export type StateTransitionTable<TState extends string> = Readonly<
	Record<TState, readonly TState[]>
>;

// True when `from → to` is a legal transition per the given table. Pure and
// dependency-free by design; a state never legally "transitions" to itself
// (staying put — e.g. a give-away remaining `claimed` after a failed Assetus
// transfer — is the absence of a transition, not a transition).
export function canTransition<TState extends string>(
	table: StateTransitionTable<TState>,
	from: TState,
	to: TState,
): boolean {
	return table[from].includes(to);
}
