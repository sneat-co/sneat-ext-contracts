import { ISpaceRequest } from '@sneat/space-models';
import { AttendanceStatus } from '../dto';

export interface ICreateClassRequest extends ISpaceRequest {
  readonly title: string;
  readonly academicYearID: string;
  readonly termID?: string;
  readonly subjectID?: string;
  readonly courseID?: string;
  readonly teacherContactIDs: string[];
  readonly classroomID?: string;
  readonly homeroom?: boolean;
}
export interface IEnrolStudentRequest extends ISpaceRequest {
  readonly studentContactID: string;
  readonly classID: string;
  readonly academicYearID: string;
  readonly termID?: string;
}
export interface IMarkAttendanceRequest extends ISpaceRequest {
  readonly lessonID: string;
  readonly enrolmentID: string;
  readonly status: AttendanceStatus;
  readonly markedByContactID: string;
  readonly note?: string;
}
export interface IRecordGradeRequest extends ISpaceRequest {
  readonly assessmentID: string;
  readonly enrolmentID: string;
  readonly value: string | number;
  readonly scale?: string;
  readonly gradedByContactID: string;
}
