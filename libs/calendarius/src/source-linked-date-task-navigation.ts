import { InjectionToken } from '@angular/core';
import type { IRelatedModules } from '@sneat/dto';
import type { ISpaceContext } from '@sneat/space-models';
import type { ISourceLinkedDateTaskMetadata } from './dto/source-linked-date-task';

export interface ISourceLinkedDateTaskNavigator {
  navigateToSourceLinkedDateTask(
    space: ISpaceContext,
    task: ISourceLinkedDateTaskMetadata,
    related: IRelatedModules,
  ): Promise<boolean>;
}

export const SOURCE_LINKED_DATE_TASK_NAVIGATOR =
  new InjectionToken<ISourceLinkedDateTaskNavigator>(
    'SourceLinkedDateTaskNavigator',
  );
