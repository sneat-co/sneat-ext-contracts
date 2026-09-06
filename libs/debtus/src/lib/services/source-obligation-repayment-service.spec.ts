import { DEBTUS_SOURCE_REPAYMENT_SERVICE_V1 } from './source-obligation-repayment-service';
import { DEBTUS_SERVICE } from './debtus-service';

describe('DEBTUS_SOURCE_REPAYMENT_SERVICE_V1', () => {
  it('is distinct from the legacy Debtus service token', () => {
    expect(DEBTUS_SOURCE_REPAYMENT_SERVICE_V1).toBeTruthy();
    expect(DEBTUS_SOURCE_REPAYMENT_SERVICE_V1).not.toBe(DEBTUS_SERVICE);
  });
});
