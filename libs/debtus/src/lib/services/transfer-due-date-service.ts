import { InjectionToken } from '@angular/core';
import { Observable } from 'rxjs';
import {
  IDebtusTransferDueDateV1,
  IDebtusTransferRepaymentV1,
  IRecordDebtusTransferRepaymentV1Request,
  ISetDebtusTransferDueDateV1Request,
} from '../models';

export interface IDebtusTransferDueDateServiceV1 {
  setTransferDueDate(
    request: ISetDebtusTransferDueDateV1Request,
  ): Observable<IDebtusTransferDueDateV1>;
  getTransferDueDate(
    spaceID: string,
    transferID: string,
  ): Observable<IDebtusTransferDueDateV1>;
  recordTransferRepayment(
    request: IRecordDebtusTransferRepaymentV1Request,
  ): Observable<IDebtusTransferRepaymentV1>;
}

export const DEBTUS_TRANSFER_DUE_DATE_SERVICE =
  new InjectionToken<IDebtusTransferDueDateServiceV1>(
    'DebtusTransferDueDateServiceV1',
  );
