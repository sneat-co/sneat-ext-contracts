import { IWithCreated, IWithSpaceIDs } from '@sneat/dto';

export type BehaviourKind = 'positive' | 'concern';
export type BehaviourVisibility = 'staff' | 'parents' | 'student';
export interface IBehaviourRecordDbo extends IWithSpaceIDs, IWithCreated {
  readonly enrolmentID: string;
  kind: BehaviourKind;
  note: string;
  authorContactID: string;
  readonly ts: number;
  visibility?: BehaviourVisibility;
}
