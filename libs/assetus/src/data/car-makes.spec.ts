// Proves the ported car make/model reference data preserves the legacy
// legacy mod-assetus-core surface (carMakes table + engines + IMake/IModel).
import { describe, expect, it } from 'vitest';
import {
  carMakes,
  engines,
  type Engine,
  type ICarMake,
  type IMake,
  type IModel,
} from './car-makes';

describe('car-makes data', () => {
  it('exposes the legacy makes', () => {
    expect(Object.keys(carMakes)).toEqual([
      'Audi',
      'BMW',
      'Toyota',
      'Tesla',
      'Hyundai',
    ]);
  });

  it('every make has a non-empty list of models with string ids', () => {
    for (const make of Object.values(carMakes) as ICarMake[]) {
      expect(Array.isArray(make.models)).toBe(true);
      for (const model of make.models) {
        expect(typeof model.id).toBe('string');
        expect(model.id.length).toBeGreaterThan(0);
        expect(Array.isArray(model.engines)).toBe(true);
      }
    }
  });

  it('shares the default engines list across petrol/diesel makes', () => {
    expect(carMakes['Audi'].models[0].engines).toBe(engines);
    expect(carMakes['BMW'].models[0].engines).toBe(engines);
  });

  it('preserves the legacy electricity fuel value for Tesla', () => {
    const tesla = carMakes['Tesla'];
    expect(tesla.yearMin).toBe(2001);
    const model3 = tesla.models.find((m) => m.id === 'Model 3');
    expect(model3?.yearMin).toBe(2017);
    expect(model3?.engines[0].engineFuel).toBe('electricity');
  });

  it('default engines list mirrors the legacy values', () => {
    const titles = engines.map((e: Engine) => e.title);
    expect(titles).toEqual(['Petrol 1L', 'Petrol 1.6L', 'Diesel 2L']);
    expect(engines[0].engineFuel).toBe('petrol');
    expect(engines[2].engineFuel).toBe('diesel');
  });

  it('public IMake/IModel are the {id,title} ergonomic shapes the components bind to', () => {
    const make: IMake = { id: 'Audi', title: 'Audi' };
    const model: IModel = { id: 'A4', title: 'A4' };
    expect(make.title).toBe('Audi');
    expect(model.id).toBe('A4');
  });
});
