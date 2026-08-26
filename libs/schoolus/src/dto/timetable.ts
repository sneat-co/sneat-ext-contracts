export interface ILesson {
  readonly id: string;
  calendarEventID: string;
  classID: string;
  teacherContactIDs: string[];
  classroomID?: string;
  topic?: string;
}
export type DayOfWeek = 'mo' | 'tu' | 'we' | 'th' | 'fr' | 'sa' | 'su';
export interface ITimetableEntry {
  readonly dayOfWeek: DayOfWeek;
  readonly startTime: string;
  readonly endTime: string;
  classID: string;
  lessonID?: string;
  classroomID?: string;
}
