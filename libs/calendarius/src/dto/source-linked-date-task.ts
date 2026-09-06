export type SourceLinkedDateTaskState = 'active' | 'completed' | 'canceled';
export type SourceLinkedDateTaskActionDisposition =
  | 'navigate'
  | 'requires_input'
  | 'execute';

export interface ISourceLinkedDateTaskRef {
  readonly namespace: string;
  readonly ownerSpaceID: string;
  readonly recordID: string;
  readonly lineID: string;
}

export interface ISourceLinkedDateTaskRelatedRef {
  readonly itemRef: {
    readonly module: string;
    readonly collection: string;
    readonly itemID: string;
    readonly subPath?: string;
  };
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
