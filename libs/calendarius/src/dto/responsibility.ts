import { IRelatedModules } from '@sneat/dto';

export type ResponsibilityAssignmentMode = 'fixed' | 'rotating';

export interface IResponsibilityAssignmentPolicy {
  readonly mode: ResponsibilityAssignmentMode;
  readonly rosterContactIDs: readonly string[];
}

export interface IScheduledResponsibilitySpec {
  readonly title: string;
  readonly description?: string;
  readonly timeZone: string;
  readonly firstDate: string;
  readonly weekday: 'mo' | 'tu' | 'we' | 'th' | 'fr' | 'sa' | 'su';
  readonly startTime: string;
  readonly durationMinutes?: number;
  readonly assignment: IResponsibilityAssignmentPolicy;
}

export interface ICreateScheduledResponsibilityRequest {
  readonly spaceID: string;
  readonly requestID: string;
  readonly spec: IScheduledResponsibilitySpec;
  /** Extension-owned payload and related refs copied to the created happening. */
  readonly happeningFields?: IResponsibilityHappeningFields;
}

export interface IResponsibilityHappeningFields {
  readonly ext?: Readonly<Record<string, unknown>>;
  readonly related?: IRelatedModules;
}

export interface IResponsibilityOccurrenceRef {
  readonly happeningID: string;
  readonly slotID: string;
  readonly date: string;
  readonly start: string;
  readonly end: string;
}

export interface IResponsibilityCompletion {
  readonly occurrenceKey: string;
  readonly happeningID: string;
  readonly slotID: string;
  readonly date: string;
  readonly assignedContactID: string;
  readonly completedBy: string;
  readonly completedAt: string;
}

export interface IResponsibilityOccurrence {
  readonly ref: IResponsibilityOccurrenceRef;
  readonly assignedContactID?: string;
  readonly needsReassignment?: boolean;
  readonly completion?: IResponsibilityCompletion;
}

export interface IScheduledResponsibility {
  readonly id: string;
  readonly spec: IScheduledResponsibilitySpec;
}

export interface ICompleteResponsibilityOccurrenceRequest {
  readonly spaceID: string;
  readonly requestID: string;
  readonly ref: IResponsibilityOccurrenceRef;
}

export type ResponsibilityMutationDisposition =
  | 'created'
  | 'completed'
  | 'unchanged'
  | 'reused';

export interface IScheduledResponsibilityMutation {
  readonly responsibility: IScheduledResponsibility;
  readonly disposition: ResponsibilityMutationDisposition;
}

export interface IResponsibilityCompletionMutation {
  readonly completion: IResponsibilityCompletion;
  readonly disposition: ResponsibilityMutationDisposition;
}
