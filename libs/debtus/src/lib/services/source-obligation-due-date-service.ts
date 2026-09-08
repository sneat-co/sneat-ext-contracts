import { InjectionToken } from '@angular/core';
import { Observable } from 'rxjs';
import type {
  IDebtusSourceObligationDueDateV1,
  ISetDebtusSourceObligationDueDateV1Request,
} from '../models/source-obligation-models';

export interface IDebtusSourceObligationDueDateServiceV1 {
  setSourceObligationDueDate(
    request: ISetDebtusSourceObligationDueDateV1Request,
  ): Observable<IDebtusSourceObligationDueDateV1>;
}

export const DEBTUS_SOURCE_OBLIGATION_DUE_DATE_SERVICE_V1 =
  new InjectionToken<IDebtusSourceObligationDueDateServiceV1>(
    'DebtusSourceObligationDueDateServiceV1',
  );
