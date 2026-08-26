import { IWithCreated, IWithSpaceIDs } from '@sneat/dto';

export interface IHomework { readonly id: string; title: string; classID: string; readonly dueDate: string; assessmentID?: string; description?: string; }
export type IAssignment = IHomework;
export type AssessmentKind = 'test' | 'assignment' | 'observation';
export interface IAssessment { readonly id: string; title: string; classID?: string; subjectID?: string; kind: AssessmentKind; maxScore?: number; scale?: string; }
export interface IGradeDbo extends IWithSpaceIDs, IWithCreated {
  readonly assessmentID: string;
  readonly enrolmentID: string;
  value: string | number;
  scale?: string;
  gradedByContactID: string;
  readonly ts: number;
}
