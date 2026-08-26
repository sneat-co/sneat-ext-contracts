// Car make/model reference data — ported faithfully from the legacy
// legacy mod-assetus-core (core/src/lib/data: car-makes-with-models.ts +
// vehicles.ts). The new lib's components (make-model-card, vehicle-card) consume
// `carMakes` and the public `IMake`/`IModel` shapes; nothing is dropped relative
// to the legacy export surface.

// Engine fuel literal — kept self-contained (the legacy `engines` data used the
// FuelTypes enum, which includes 'electricity' that the new lib's FuelType union
// does not). Preserving the legacy values verbatim avoids losing data.
export type CarEngineFuel =
  | ''
  | 'other'
  | 'petrol'
  | 'diesel'
  | 'hydrogen'
  | 'electricity';

// Engine mirrors the legacy data/vehicles.ts `Engine`.
export interface Engine {
  title: string;
  engineType?: string;
  engineFuel: CarEngineFuel;
  engineCC?: number;
  engineKW?: number;
  engineNM?: number;
}

// engines mirrors the legacy default engine list (data/vehicles.ts).
export const engines: Engine[] = [
  { title: 'Petrol 1L', engineCC: 999, engineFuel: 'petrol' },
  { title: 'Petrol 1.6L', engineCC: 1599, engineFuel: 'petrol' },
  { title: 'Diesel 2L', engineCC: 1999, engineFuel: 'diesel' },
];

// IMake / IModel are the public ergonomic shapes the components bind to (legacy
// data/vehicles.ts: `{ id, title }`).
export interface IMake {
  id: string;
  title: string;
}

export interface IModel {
  id: string;
  title: string;
}

// ICarModel / ICarMake are the internal shapes of the `carMakes` data table
// (legacy data/car-makes-with-models.ts), kept distinct from the public
// IMake/IModel to preserve both surfaces without collision.
export interface ICarModel {
  id: string;
  yearMin?: number;
  yearMax?: string;
  engines: Engine[];
}

export interface ICarMake {
  yearMin?: number;
  yearMax?: string;
  models: ICarModel[];
}

// carMakes mirrors the legacy data/car-makes-with-models.ts table verbatim.
export const carMakes: Record<string, ICarMake> = {
  Audi: {
    models: [
      { id: 'A1', engines },
      { id: 'A2', engines },
      { id: 'A3', engines },
      { id: 'A4', engines },
      { id: 'A5', engines },
      { id: 'A6', engines },
      { id: 'A7', engines },
      { id: 'A8', engines },
      { id: 'S4', engines },
      { id: 'S6', engines },
      { id: 'S8', engines },
      { id: 'Q3', engines },
      { id: 'Q5', engines },
      { id: 'Q7', engines },
    ],
  },
  BMW: {
    models: [
      { id: '1 Series', engines },
      { id: '2 Series', engines },
      { id: '3 Series Saloon', engines },
      { id: '3 Series Touring', engines },
      { id: '3 Series Gran Tourismo', engines },
      { id: '4 Series', engines },
      { id: '5 Series', engines },
      { id: '6 Series', engines },
      { id: '7 Series', engines },
      { id: '8 Series', engines },
      { id: 'i3', engines },
      { id: 'i8 Coupe', engines },
      { id: 'i8 Roadster', engines },
      { id: 'X1', engines },
      { id: 'X2', engines },
      { id: 'X3', engines },
      { id: 'X4', engines },
      { id: 'X5', engines },
      { id: 'X6', engines },
      { id: 'X7', engines },
      { id: 'z3', engines },
      { id: 'z4', engines },
      { id: 'z5', engines },
    ],
  },
  Toyota: {
    models: [
      { id: 'Avensis', engines },
      { id: 'Auris', engines },
      { id: 'Aygo', engines },
      { id: 'GR', engines },
      { id: 'GT86', engines },
      { id: 'C-HR', engines },
      { id: 'Camry', engines },
      { id: 'Corolla', engines },
      { id: 'Hilux', engines },
      { id: 'Prius', engines },
      { id: 'Prius+', engines },
      { id: 'Proace', engines },
      { id: 'Proace Verso', engines },
      { id: 'RAV4', engines },
      { id: 'Yaris', engines },
    ],
  },
  Tesla: {
    yearMin: 2001,
    models: [
      {
        id: 'Model 3',
        yearMin: 2017,
        engines: [{ title: '90KW', engineFuel: 'electricity' }],
      },
      {
        id: 'Model S',
        engines: [{ title: '90KW', engineFuel: 'electricity' }],
      },
      {
        id: 'Model X',
        engines: [{ title: '90KW', engineFuel: 'electricity' }],
      },
    ],
  },
  Hyundai: {
    yearMin: 1900,
    models: [
      { id: 'Kona', engines: [] },
      { id: 'Tuscan', engines: [] },
    ],
  },
};
