import { readFileSync } from 'node:fs';
import { fileURLToPath } from 'node:url';
import { describe, expect, it } from 'vitest';
import { fold, inBonus, order, type GameEvent, type GameState } from './eventtimeline.js';

interface ParityCase {
  name: string;
  events: GameEvent[];
  expected: GameState;
}

// Provenance note: in the source repo (sneat-co/ext-gameboard) this fixture
// lived at repo-root parity/parity.json, referenced here as
// '../../parity/parity.json' (up from frontend/src/). In this monorepo the
// Go module directory (gameboard/, a sibling of libs/) plays the role the old
// repo root played, so the fixture now lives at gameboard/parity/parity.json
// and the relative path below crosses one extra directory level
// ('../../../gameboard/parity/parity.json', up from libs/gameboard/src/).
// Content is otherwise byte-identical to the source (see
// gameboard/eventtimeline/parity_test.go for the Go side of this oracle).
const fixtureUrl = new URL('../../../gameboard/parity/parity.json', import.meta.url);
const cases: ParityCase[] = JSON.parse(readFileSync(fileURLToPath(fixtureUrl), 'utf8'));

describe('cross-language reducer parity (Go ↔ TS via parity.json)', () => {
  it('fixture is non-empty', () => {
    expect(cases.length).toBeGreaterThan(0);
  });

  for (const c of cases) {
    it(`folds "${c.name}" to the same state as the Go reducer`, () => {
      expect(fold(c.events)).toEqual(c.expected);
    });
  }
});

describe('reducer unit behaviour', () => {
  it('orders by wallClock then eventID and dedupes by eventID', () => {
    const ev = (eventID: string, wallClockMs: number): GameEvent => ({
      eventID, wallClockMs, type: 'score', source: 'scorekeeper', period: 0, gameClockMs: 0,
    });
    const ordered = order([ev('b', 10), ev('a', 10), ev('c', 5), ev('a', 10)]);
    expect(ordered.map((e) => e.eventID)).toEqual(['c', 'a', 'b']);
  });

  it('inBonus flips when opponent reaches the foul limit', () => {
    const s = fold([
      { eventID: 'f1', type: 'team-foul', source: 'scorekeeper', wallClockMs: 1, period: 0, gameClockMs: 0, side: 'home' },
      { eventID: 'f2', type: 'team-foul', source: 'scorekeeper', wallClockMs: 2, period: 0, gameClockMs: 0, side: 'home' },
    ]);
    expect(inBonus(s, 'away', 2)).toBe(true);
    expect(inBonus(s, 'home', 2)).toBe(false);
  });
});
