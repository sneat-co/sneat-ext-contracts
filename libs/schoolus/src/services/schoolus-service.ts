import { InjectionToken } from '@angular/core';
import { IRecord } from '@sneat/data';
import { ISpaceContext } from '@sneat/space-models';
import { Observable } from 'rxjs';
import { ISchoolClassContext } from '../contexts';
import { IAttendanceRecordDbo, IEnrolmentDbo, IGradeDbo, ITimetableEntry } from '../dto';
import { ICreateClassRequest, IEnrolStudentRequest, IMarkAttendanceRequest, IRecordGradeRequest } from './interfaces';

/** Public, implementation-free API consumed by Schoolus UI and hosts. */
export interface ISchoolusService {
  createClass(request: ICreateClassRequest): Observable<ISchoolClassContext>;
  enrolStudent(request: IEnrolStudentRequest): Observable<IRecord<IEnrolmentDbo>>;
  markAttendance(request: IMarkAttendanceRequest): Observable<IRecord<IAttendanceRecordDbo>>;
  getTimetable(space: ISpaceContext, classID: string): Observable<ITimetableEntry[]>;
  recordGrade(request: IRecordGradeRequest): Observable<IRecord<IGradeDbo>>;
}

export const SCHOOLUS_SERVICE = new InjectionToken<ISchoolusService>('SchoolusService');
