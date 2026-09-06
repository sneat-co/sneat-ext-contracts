import { InjectionToken } from '@angular/core';
import { Observable } from 'rxjs';
import {
  IDebtusTransferDueDateV1,
  ISetDebtusTransferDueDateV1Request,
} from '../models';

export interface IDebtusTransferDueDateServiceV1 {
  setTransferDueDate(
    request: ISetDebtusTransferDueDateV1Request,
  ): Observable<IDebtusTransferDueDateV1>;
}

export const DEBTUS_TRANSFER_DUE_DATE_SERVICE =
  new InjectionToken<IDebtusTransferDueDateServiceV1>(
    'DebtusTransferDueDateServiceV1',
  );
