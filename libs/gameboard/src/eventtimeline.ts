// Deterministic event-timeline reducer — the TS mirror of the Go reducer in
// gameboard/eventtimeline/fold.go (this repo). Both are hand-implemented to
// match the frozen wire contract (api4gameboard.tsp, kept in
// sneat-co/ext-gameboard) and MUST fold the same log to identical state;
// ./eventtimeline.parity.test.ts's parity.json oracle proves it.

export type TeamSide = 'home' | 'away';

export interface Side {
  name: string;
  colour: string;
  spaceID?: string | null;
}

export type EventType =
  | 'status' | 'period' | 'clock' | 'score' | 'team-foul'
  | 'timeout' | 'substitution' | 'possession' | 'judge-ruling' | 'correction';

export type GameStatus =
  | 'scheduled' | 'live' | 'halftime' | 'overtime' | 'final' | 'cancelled';

export type ClockAction = 'start' | 'stop' | 'adjust';

export type Source = 'scorekeeper' | 'timekeeper' | 'judge' | 'consensus';

export interface GameEvent {
  eventID: string;
  type: EventType;
  source: Source;
  wallClockMs: number;
  period: number;
  gameClockMs: number;
  status?: GameStatus;
  clockAction?: ClockAction;
  side?: TeamSide;
  points?: number;
  scorerID?: string;
  assistID?: string;
  playerOn?: string;
  playerOff?: string;
  correctionOf?: string;
}

export interface GameState {
  status: GameStatus;
  period: number;
  gameClockMs: number;
  clockRunning: boolean;
  scores: Record<TeamSide, number>;
  teamFouls: Record<TeamSide, number>;
  timeoutsUsed: Record<TeamSide, number>;
  possession: TeamSide | '';
  onCourt: Record<TeamSide, string[]>;
  playerPoints: Record<string, number>;
  playerAssists: Record<string, number>;
}

function validSide(s: unknown): s is TeamSide {
  return s === 'home' || s === 'away';
}

function newState(): GameState {
  return {
    status: 'scheduled',
    period: 0,
    gameClockMs: 0,
    clockRunning: false,
    scores: { home: 0, away: 0 },
    teamFouls: { home: 0, away: 0 },
    timeoutsUsed: { home: 0, away: 0 },
    possession: '',
    onCourt: { home: [], away: [] },
    playerPoints: {},
    playerAssists: {},
  };
}

/** Canonical total order: wallClockMs asc, ties by eventID; idempotent dedupe
 * (first occurrence of an eventID wins). Mirrors Go Order(). Input not mutated. */
export function order(events: GameEvent[]): GameEvent[] {
  const seen = new Set<string>();
  const out: GameEvent[] = [];
  for (const e of events) {
    if (seen.has(e.eventID)) continue;
    seen.add(e.eventID);
    out.push(e);
  }
  return out.slice().sort((a, b) => {
    if (a.wallClockMs !== b.wallClockMs) return a.wallClockMs - b.wallClockMs;
    return a.eventID < b.eventID ? -1 : a.eventID > b.eventID ? 1 : 0;
  });
}

/** Deterministic projection of the log. Mirrors Go Fold(). */
export function fold(events: GameEvent[]): GameState {
  const ordered = order(events);

  const voided = new Set<string>();
  for (const e of ordered) {
    if (e.type === 'correction' && e.correctionOf) voided.add(e.correctionOf);
  }

  const s = newState();
  for (const e of ordered) {
    if (voided.has(e.eventID)) continue;
    applyEvent(s, e);
  }
  return s;
}

function applyEvent(s: GameState, e: GameEvent): void {
  switch (e.type) {
    case 'status':
      if (e.status) s.status = e.status;
      break;
    case 'period':
      s.period = e.period;
      break;
    case 'clock':
      if (e.clockAction === 'start') {
        s.clockRunning = true;
        if (e.gameClockMs > 0) s.gameClockMs = e.gameClockMs;
      } else if (e.clockAction === 'stop') {
        s.clockRunning = false;
        s.gameClockMs = e.gameClockMs;
      } else if (e.clockAction === 'adjust') {
        s.gameClockMs = e.gameClockMs;
      }
      break;
    case 'score':
      if (validSide(e.side) && (e.points ?? 0) > 0) {
        s.scores[e.side] += e.points as number;
        if (e.scorerID) s.playerPoints[e.scorerID] = (s.playerPoints[e.scorerID] ?? 0) + (e.points as number);
        if (e.assistID) s.playerAssists[e.assistID] = (s.playerAssists[e.assistID] ?? 0) + 1;
      }
      break;
    case 'team-foul':
      if (validSide(e.side)) s.teamFouls[e.side]++;
      break;
    case 'timeout':
      if (validSide(e.side)) s.timeoutsUsed[e.side]++;
      break;
    case 'possession':
      if (validSide(e.side)) s.possession = e.side;
      break;
    case 'substitution':
      if (validSide(e.side)) {
        if (e.playerOff) s.onCourt[e.side] = s.onCourt[e.side].filter((p) => p !== e.playerOff);
        if (e.playerOn && !s.onCourt[e.side].includes(e.playerOn)) s.onCourt[e.side].push(e.playerOn);
      }
      break;
    case 'correction':
      if (validSide(e.side) && (e.points ?? 0) > 0) s.scores[e.side] += e.points as number;
      break;
    case 'judge-ruling':
      // audit entry; no projection effect unless it carries a correction
      break;
  }
}

/** InBonus: side shoots bonus once the opponent has >= limit team fouls. */
export function inBonus(s: GameState, side: TeamSide, limit: number): boolean {
  const opp: TeamSide = side === 'away' ? 'home' : 'away';
  return s.teamFouls[opp] >= limit;
}
