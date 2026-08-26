/** Legacy game-record contract retained alongside the event-timeline model. */
export interface GameRecord {
  id: string;
  home: { name: string; colour: string; spaceID?: string | null };
  away: { name: string; colour: string; spaceID?: string | null };
  status: string;
  scheduledMs?: number;
}

export type CreatorRole =
  | 'coach'
  | 'player'
  | 'score-keeper'
  | 'timekeeper'
  | 'judge'
  | 'spectator';

export const CREATOR_ROLES: readonly CreatorRole[] = [
  'coach',
  'player',
  'score-keeper',
  'timekeeper',
  'judge',
  'spectator',
];
