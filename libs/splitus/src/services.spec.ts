import { SPLITUS_BILL_SERVICE_V1, SPLITUS_SERVICE } from './services';

describe('SPLITUS_SERVICE token', () => {
  it('should be defined', () => {
    expect(SPLITUS_SERVICE).toBeTruthy();
  });
});

describe('SPLITUS_BILL_SERVICE_V1 token', () => {
  it('is defined independently from the legacy token', () => {
    expect(SPLITUS_BILL_SERVICE_V1).toBeTruthy();
    expect(SPLITUS_BILL_SERVICE_V1).not.toBe(SPLITUS_SERVICE);
  });
});
