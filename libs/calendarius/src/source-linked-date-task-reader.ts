import { InjectionToken } from '@angular/core';
import type { Observable } from 'rxjs';
import type { ISourceLinkedDateTask } from './dto/source-linked-date-task';

/** Reads the Calendar-owned canonical date task without copying its due date. */
export interface ISourceLinkedDateTaskReader {
  observeSourceLinkedDateTask(
    spaceID: string,
    happeningID: string,
  ): Observable<ISourceLinkedDateTask | undefined>;
}

export const SOURCE_LINKED_DATE_TASK_READER =
  new InjectionToken<ISourceLinkedDateTaskReader>(
    'SourceLinkedDateTaskReader',
  );
