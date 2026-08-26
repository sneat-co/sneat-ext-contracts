import { IWithCreated, IWithSpaceIDs } from '@sneat/dto';

export type EnrolmentStatus = 'active' | 'withdrawn' | 'graduated';
export interface IEnrolmentDbo extends IWithSpaceIDs, IWithCreated {
  readonly studentContactID: string;
  readonly classID: string;
  academicYearID: string;
  termID?: string;
  status: EnrolmentStatus;
  role: 'student';
  enrolledDate?: string;
}
