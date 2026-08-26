import {
  IAssetApplianceExtra,
  IAssetDocumentExtra,
  IAssetElectronicsExtra,
  IAssetUtilityExtra,
  IAssetVehicleExtra,
  IRenewalsResponse,
} from './index';

describe('new asset-type extras', () => {
  it('appliance extra round-trips service/warranty dates', () => {
    const a: IAssetApplianceExtra = {
      make: 'Worcester',
      model: 'Greenstar',
      serialNumber: 'WB-1',
      installedOn: '2022-03-01',
      warrantyExpiresOn: '2030-03-01',
      nextServiceDue: '2027-03-01',
      energyRating: 'A',
    };
    expect(JSON.parse(JSON.stringify(a))).toEqual(a);
  });

  it('electronics extra carries serial + imei + warranty', () => {
    const e: IAssetElectronicsExtra = {
      make: 'Apple',
      model: 'MacBook Pro',
      serialNumber: 'C02X',
      imei: '35-2099',
      warrantyExpiresOn: '2028-01-01',
    };
    expect(JSON.parse(JSON.stringify(e))).toEqual(e);
  });

  it('utility extra carries provider + renewalDate', () => {
    const u: IAssetUtilityExtra = {
      provider: 'Electric Ireland',
      accountNumber: 'ACC-1',
      tariff: 'Standard',
      renewalDate: '2026-12-01',
      contractEndsOn: '2026-12-01',
    };
    expect(JSON.parse(JSON.stringify(u))).toEqual(u);
    expect(u.renewalDate).toBe('2026-12-01');
  });

  it('vehicle extra carries warrantyExpiresOn', () => {
    const v: IAssetVehicleExtra = { make: 'Toyota', warrantyExpiresOn: '2027-01-01' };
    expect(v.warrantyExpiresOn).toBe('2027-01-01');
  });

  it('document extra carries the generalized coversAssetID', () => {
    const d: IAssetDocumentExtra = {
      docType: 'warranty',
      expiresOn: '2030-01-01',
      coversAssetID: 'laptop1',
    };
    expect(d.coversAssetID).toBe('laptop1');
  });
});

describe('renewals feed', () => {
  it('IRenewalsResponse round-trips items with period', () => {
    const resp: IRenewalsResponse = {
      renewals: [
        {
          assetID: 'a1',
          assetName: 'Car',
          category: 'vehicles',
          kind: 'nct',
          dueOn: '2026-09-01',
          period: 'yearly',
        },
        {
          assetID: 'a2',
          kind: 'warranty',
          dueOn: '2027-01-01',
          period: 'once',
        },
      ],
    };
    const rt: IRenewalsResponse = JSON.parse(JSON.stringify(resp));
    expect(rt).toEqual(resp);
    expect(rt.renewals[0].period).toBe('yearly');
    expect(rt.renewals[1].period).toBe('once');
  });
});
