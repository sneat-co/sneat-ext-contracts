import { InjectionToken } from '@angular/core';
import { Observable } from 'rxjs';
import type {
  IRecordDebtusSourceRepaymentV1Request,
  IRecordDebtusSourceRepaymentV1Response,
} from '../models/source-obligation-models';

/** Browser service for the Debtus source repayment contract version 1. */
export interface IDebtusSourceRepaymentServiceV1 {
  recordSourceObligationRepayment(
    request: IRecordDebtusSourceRepaymentV1Request,
  ): Observable<IRecordDebtusSourceRepaymentV1Response>;
}

export const DEBTUS_SOURCE_REPAYMENT_SERVICE_V1 =
  new InjectionToken<IDebtusSourceRepaymentServiceV1>(
    'DebtusSourceRepaymentServiceV1',
  );
