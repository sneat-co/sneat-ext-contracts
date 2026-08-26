import { IWithCreated, IWithSpaceIDs } from '@sneat/dto';

export type AttendanceStatus = 'present' | 'absent' | 'late' | 'excused';
export interface IAttendanceRecordDbo extends IWithSpaceIDs, IWithCreated {
  readonly lessonID: string;
  readonly enrolmentID: string;
  status: AttendanceStatus;
  markedByContactID: string;
  readonly ts: number;
  note?: string;
}
