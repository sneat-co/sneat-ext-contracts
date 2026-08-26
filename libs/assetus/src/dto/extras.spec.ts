import {
  assetExtraTypeDigitalAsset,
  assetExtraTypeDocument,
  assetExtraTypeDwelling,
  assetExtraTypeSportsEquipment,
  assetExtraTypeValuable,
  assetExtraTypeVehicle,
  bicycleKinds,
  docTypeSchema,
  standardDocTypesByID,
  type IAssetDigitalAssetExtra,
  type IAssetDocumentExtra,
  type IAssetDwellingExtra,
  type IAssetSportsEquipmentExtra,
  type IAssetValuableExtra,
  type IAssetVehicleExtra,
  type IVehicleRecord,
} from './extras';

// AC: vehicle-extra-frontend-no-field-dropped
describe('IAssetVehicleExtra', () => {
  it('round-trips every field incl. engineSerialNumber with no field dropped', () => {
    const vehicle: IAssetVehicleExtra = {
      make: 'Toyota',
      model: 'Corolla',
      regNumber: '12-D-3456',
      vin: 'JT2BG22K1W0123456',
      engineType: 'combustion',
      engineFuel: 'petrol',
      engineCC: 1800,
      engineKW: 103,
      engineNM: 173,
      engineSerialNumber: 'ENG-SERIAL-001',
      nctExpires: '2026-01-01',
      taxExpires: '2026-02-01',
      nextServiceDue: '2026-03-01',
    };

    const roundTripped: IAssetVehicleExtra = JSON.parse(
      JSON.stringify(vehicle),
    );
    expect(roundTripped).toEqual(vehicle);
    // The engine serial number is preserved.
    expect(roundTripped.engineSerialNumber).toBe('ENG-SERIAL-001');
    // Exact backend camelCase json names are present.
    const keys = Object.keys(roundTripped);
    expect(keys).toContain('engineCC');
    expect(keys).toContain('engineKW');
    expect(keys).toContain('engineNM');
    expect(keys).toContain('engineFuel');
    expect(keys).toContain('regNumber');
    expect(keys).toContain('vin');
  });
});

// AC: vehicle-extra-frontend-no-field-dropped (vehicle record mileage + fuel)
describe('IVehicleRecord', () => {
  it('round-trips a record with mileage AND full fuel payload', () => {
    const record: IVehicleRecord = {
      mileage: { value: 123456, unit: 'km' },
      fuel: {
        volume: 42.5,
        unit: 'l',
        amount: { currency: 'EUR', value: 70 },
        fuelCost: 1.65,
        currency: 'EUR',
      },
      createdAt: '2026-01-01T00:00:00Z',
      createdBy: 'user1',
    };

    const roundTripped: IVehicleRecord = JSON.parse(JSON.stringify(record));
    expect(roundTripped).toEqual(record);
    // mileage {value,unit}
    expect(roundTripped.mileage).toEqual({ value: 123456, unit: 'km' });
    // fuel {volume,unit,amount,fuelCost,currency} — no field dropped
    const fuelKeys = Object.keys(roundTripped.fuel ?? {});
    expect(fuelKeys).toEqual([
      'volume',
      'unit',
      'amount',
      'fuelCost',
      'currency',
    ]);
    expect(roundTripped.fuel?.amount).toEqual({ currency: 'EUR', value: 70 });
  });
});

// AC: document-extra-frontend-full-shape
describe('IAssetDocumentExtra', () => {
  it('round-trips a passport with all 8 fields', () => {
    const passport: IAssetDocumentExtra = {
      docType: 'passport',
      number: 'P1234567',
      batchNumber: 'B-001',
      countryID: 'IE',
      issuedBy: 'DFA',
      issuedOn: '2020-01-01',
      effectiveFrom: '2020-01-01',
      expiresOn: '2030-01-01',
    };

    const roundTripped: IAssetDocumentExtra = JSON.parse(
      JSON.stringify(passport),
    );
    expect(roundTripped).toEqual(passport);
    const keys = Object.keys(roundTripped);
    expect(keys).toEqual([
      'docType',
      'number',
      'batchNumber',
      'countryID',
      'issuedBy',
      'issuedOn',
      'effectiveFrom',
      'expiresOn',
    ]);
  });

  it('round-trips an insurance policy linked to the insured vehicle', () => {
    const policy: IAssetDocumentExtra = {
      docType: 'insurance',
      number: 'POL-9',
      issuedBy: 'Acme Insurance',
      expiresOn: '2027-01-01',
      insuredAssetID: 'veh1',
    };
    const roundTripped: IAssetDocumentExtra = JSON.parse(
      JSON.stringify(policy),
    );
    expect(roundTripped).toEqual(policy);
    expect(roundTripped.insuredAssetID).toBe('veh1');
  });

  it('passport schema requires number + validity (validTill)', () => {
    const schema = docTypeSchema('passport');
    expect(schema?.fields?.number?.required).toBe(true);
    expect(schema?.fields?.validTill?.required).toBe(true);
    expect(schema?.fields?.members?.max).toBe(1);
  });

  it('marriage_cert requires number + issuedOn, excludes validTill, max 2 members', () => {
    const schema = standardDocTypesByID['marriage_cert'];
    expect(schema.fields?.number?.required).toBe(true);
    expect(schema.fields?.issuedOn?.required).toBe(true);
    expect(schema.fields?.validTill?.exclude).toBe(true);
    expect(schema.fields?.members?.max).toBe(2);
  });

  it('insurance schema requires number + validity (validTill)', () => {
    const schema = docTypeSchema('insurance');
    expect(schema?.fields?.number?.required).toBe(true);
    expect(schema?.fields?.validTill?.required).toBe(true);
    expect(schema?.fields?.members?.max).toBe(1);
  });

  it('warranty requires validity, home_insurance requires number + validity, life_insurance requires only number', () => {
    expect(docTypeSchema('warranty')?.fields?.validTill?.required).toBe(true);
    expect(docTypeSchema('home_insurance')?.fields?.number?.required).toBe(true);
    expect(docTypeSchema('home_insurance')?.fields?.validTill?.required).toBe(
      true,
    );
    expect(docTypeSchema('life_insurance')?.fields?.number?.required).toBe(true);
    expect(
      docTypeSchema('life_insurance')?.fields?.validTill?.required,
    ).toBeFalsy();
  });

  it('covers the backend AssetDocumentType taxonomy', () => {
    expect(Object.keys(standardDocTypesByID).sort()).toEqual([
      'birth_cert',
      'driving_license',
      'home_insurance',
      'insurance',
      'life_insurance',
      'marriage_cert',
      'other',
      'passport',
      'warranty',
    ]);
  });

  it('returns undefined for a docType with no standard schema', () => {
    expect(docTypeSchema(undefined)).toBeUndefined();
    expect(docTypeSchema('id_card')).toBeUndefined();
  });
});

// AC: dwelling extra full shape
describe('IAssetDwellingExtra', () => {
  it('round-trips address, rent_price, bedrooms and area', () => {
    const dwelling: IAssetDwellingExtra = {
      address: {
        countryID: 'IE',
        zipCode: 'D02 XY45',
        state: 'Leinster',
        city: 'Dublin',
        lines: '1 Main St',
      },
      rent_price: { value: 2000, currency: 'EUR' },
      numberOfBedrooms: 3,
      areaSqM: 95,
    };
    const roundTripped: IAssetDwellingExtra = JSON.parse(
      JSON.stringify(dwelling),
    );
    expect(roundTripped).toEqual(dwelling);
    expect(roundTripped.rent_price).toEqual({ value: 2000, currency: 'EUR' });
  });
});

describe('extraType discriminator', () => {
  it('mirrors the backend extraType values', () => {
    expect(assetExtraTypeVehicle).toBe('vehicle');
    expect(assetExtraTypeDwelling).toBe('dwelling');
    expect(assetExtraTypeDocument).toBe('document');
    expect(assetExtraTypeValuable).toBe('valuable');
    expect(assetExtraTypeDigitalAsset).toBe('digital_asset');
    expect(assetExtraTypeSportsEquipment).toBe('sports_equipment');
  });
});

// AC: boat/aircraft vehicle fields (additive to IAssetVehicleExtra)
describe('IAssetVehicleExtra boat/aircraft fields', () => {
  it('round-trips lengthM and hullOrTailNumber alongside make/model', () => {
    const boat: IAssetVehicleExtra = {
      make: 'Beneteau',
      model: 'Oceanis',
      regNumber: 'IE-1234',
      lengthM: 12.5,
      hullOrTailNumber: 'HIN-12345',
    };
    const roundTripped: IAssetVehicleExtra = JSON.parse(JSON.stringify(boat));
    expect(roundTripped).toEqual(boat);
    expect(roundTripped.lengthM).toBe(12.5);
    expect(roundTripped.hullOrTailNumber).toBe('HIN-12345');
  });
});

// AC: valuable extra (jewelry/watches/art/collectibles)
describe('IAssetValuableExtra', () => {
  it('round-trips material, appraisedValue/appraisalDate, certificateNumber', () => {
    const valuable: IAssetValuableExtra = {
      material: 'gold',
      appraisedValue: { currency: 'USD', value: 5000 },
      appraisalDate: '2026-01-15',
      certificateNumber: 'CERT-1',
    };
    const roundTripped: IAssetValuableExtra = JSON.parse(
      JSON.stringify(valuable),
    );
    expect(roundTripped).toEqual(valuable);
    expect(roundTripped.appraisedValue).toEqual({
      currency: 'USD',
      value: 5000,
    });
  });
});

// AC: digital asset extra (domain/software licence/digital collectible)
describe('IAssetDigitalAssetExtra', () => {
  it('round-trips registrar/identifier/renewalDate/autoRenew/renewalCost', () => {
    const digitalAsset: IAssetDigitalAssetExtra = {
      registrar: 'Namecheap',
      identifier: 'example.com',
      renewalDate: '2027-05-01',
      autoRenew: true,
      renewalCost: { currency: 'EUR', value: 15 },
    };
    const roundTripped: IAssetDigitalAssetExtra = JSON.parse(
      JSON.stringify(digitalAsset),
    );
    expect(roundTripped).toEqual(digitalAsset);
  });
});

// AC: sports-equipment extra (bicycle make/model/frame number + bicycleKind)
describe('IAssetSportsEquipmentExtra', () => {
  it('round-trips make/model/frameNumber/bicycleKind', () => {
    const bike: IAssetSportsEquipmentExtra = {
      make: 'Trek',
      model: 'Domane',
      frameNumber: 'FR-123',
      bicycleKind: 'road',
    };
    const roundTripped: IAssetSportsEquipmentExtra = JSON.parse(
      JSON.stringify(bike),
    );
    expect(roundTripped).toEqual(bike);
  });

  it('bicycleKinds covers road/mountain/e_bike/kids/hybrid/gravel/other', () => {
    expect(bicycleKinds).toEqual([
      'road',
      'mountain',
      'e_bike',
      'kids',
      'hybrid',
      'gravel',
      'other',
    ]);
  });
});
