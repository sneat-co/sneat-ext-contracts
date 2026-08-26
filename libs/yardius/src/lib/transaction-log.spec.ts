import { describe, it, expect } from 'vitest';
import { ITransactionLogEntry, YARDIUS_EVENT_TYPES } from './transaction-log.js';

describe('YARDIUS_EVENT_TYPES', () => {
	it('is exactly the REQ transaction-log event set', () => {
		expect(YARDIUS_EVENT_TYPES).toEqual([
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
		]);
	});
});

describe('ITransactionLogEntry round-trip', () => {
	it('preserves event type, timestamp, and acting member through JSON serialize/parse (AC transaction-log-append-only)', () => {
		const entry: ITransactionLogEntry = {
			event: 'approve',
			at: '2026-07-02T10:15:30.000Z',
			actingMemberSpaceID: 'space-x',
			actingMemberID: 'member-anna',
		};

		const roundTripped: ITransactionLogEntry = JSON.parse(JSON.stringify(entry));

		expect(roundTripped).toEqual(entry);
	});

	it('round-trips one entry per event type', () => {
		const entries: readonly ITransactionLogEntry[] = YARDIUS_EVENT_TYPES.map(
			(event, i) => ({
				event,
				at: `2026-07-02T10:00:0${i % 10}.000Z`,
				actingMemberSpaceID: 'space-y',
				actingMemberID: 'member-bob',
			}),
		);

		const roundTripped: readonly ITransactionLogEntry[] = JSON.parse(
			JSON.stringify(entries),
		);

		expect(roundTripped).toEqual(entries);
		expect(roundTripped.map((e) => e.event)).toEqual([...YARDIUS_EVENT_TYPES]);
	});
});
