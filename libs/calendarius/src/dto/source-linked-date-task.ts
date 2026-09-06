export type SourceLinkedDateTaskState = 'active' | 'completed' | 'canceled';
export type SourceLinkedDateTaskActionDisposition =
  | 'navigate'
  | 'requires_input';

export interface ISourceLinkedDateTaskRef {
  readonly namespace: string;
  readonly ownerSpaceID: string;
  readonly recordID: string;
  readonly lineID: string;
}

import type { ISpaceModuleItemRef } from '@sneat/dto';

export interface ISourceLinkedDateTaskRelatedRef {
  readonly itemRef: ISpaceModuleItemRef;
  readonly role: string;
}

export interface ISourceLinkedDateTaskMetadata {
  readonly revision: number;
  readonly source: ISourceLinkedDateTaskRef;
  readonly state: SourceLinkedDateTaskState;
  readonly actionID?: string;
  readonly actionDisposition?: SourceLinkedDateTaskActionDisposition;
}

export interface ISourceLinkedDateTask extends ISourceLinkedDateTaskMetadata {
  readonly happeningID: string;
  readonly title: string;
  readonly dueDate?: string;
}
