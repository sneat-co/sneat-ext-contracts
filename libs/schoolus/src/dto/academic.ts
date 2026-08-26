import { IWithCreated, IWithSpaceIDs } from '@sneat/dto';

export interface ISubject { readonly id: string; title: string; code?: string; }
export interface ICourse { readonly id: string; title: string; subjectID?: string; description?: string; }
export interface ITerm { readonly id: string; title: string; readonly startDate: string; readonly endDate: string; }
export interface IAcademicYear {
  readonly id: string;
  title: string;
  readonly startDate?: string;
  readonly endDate?: string;
  terms?: ITerm[];
  termIDs?: string[];
}
export type SchoolClassStatus = 'active' | 'archived';
export interface ISchoolClassDbo extends IWithSpaceIDs, IWithCreated {
  title: string;
  subjectID?: string;
  courseID?: string;
  academicYearID: string;
  termID?: string;
  teacherContactIDs: string[];
  classroomID?: string;
  homeroom?: boolean;
  status?: SchoolClassStatus;
}
export interface IGroup { readonly id: string; title: string; enrolmentIDs: string[]; }
