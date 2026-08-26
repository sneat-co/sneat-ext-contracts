import { ISizeTypeCatalog } from './size-type';

// The MVP size-type catalog (version 1). Pure data — per AC
// catalog-extends-without-migration, adding a new size type here (e.g. "ski
// boots") never requires a schema change: `ISizeRecord.sizeTypeID` is a plain
// string, so a new entry just needs to be appended to `sizeTypes` (and to a
// `sportKits` grouping if applicable) and it is immediately recordable and
// displayable like any other size type.
//
// Coverage below is deliberately exactly the brief's MVP scope (per AC
// mvp-catalog-coverage): footwear; body measurements (height, weight, chest,
// waist, hips, inseam); tops; bottoms; accessories (hat, gloves, ring,
// helmet); and the football/basketball/cycling/swimming kit types.
export const MVP_SIZE_TYPE_CATALOG: ISizeTypeCatalog = {
	version: 1,
	sizeTypes: [
		// Footwear. Adult and kids are kept as distinct entries (rather than one
		// "footwear" entry) because their EU/UK/US conversion tables do not
		// share a single reliable mapping — see `conversion/footwear.ts`.
		{
			id: 'footwear-adult',
			category: 'footwear',
			title: 'Footwear (adult)',
			valueKind: 'numeric',
			applicableSystems: ['EU', 'UK', 'US', 'JP', 'Mondopoint'],
		},
		{
			id: 'footwear-kids',
			category: 'footwear',
			title: 'Footwear (kids)',
			valueKind: 'numeric',
			applicableSystems: ['EU', 'UK', 'US', 'JP', 'Mondopoint'],
		},

		// Body measurements.
		{
			id: 'height',
			category: 'body-measurement',
			title: 'Height',
			valueKind: 'measurement',
			applicableSystems: ['cm', 'in'],
		},
		{
			id: 'weight',
			category: 'body-measurement',
			title: 'Weight',
			valueKind: 'measurement',
			// `SizeSystem` is an open union — 'kg'/'lb' are not first-party
			// literals but are valid declared systems for this type.
			applicableSystems: ['kg', 'lb'],
		},
		{
			id: 'chest',
			category: 'body-measurement',
			title: 'Chest',
			valueKind: 'measurement',
			applicableSystems: ['cm', 'in'],
		},
		{
			id: 'waist',
			category: 'body-measurement',
			title: 'Waist',
			valueKind: 'measurement',
			applicableSystems: ['cm', 'in'],
		},
		{
			id: 'hips',
			category: 'body-measurement',
			title: 'Hips',
			valueKind: 'measurement',
			applicableSystems: ['cm', 'in'],
		},
		{
			id: 'inseam',
			category: 'body-measurement',
			title: 'Inseam',
			valueKind: 'measurement',
			applicableSystems: ['cm', 'in'],
		},

		// Tops / bottoms.
		{
			id: 'tops',
			category: 'tops',
			title: 'Tops',
			valueKind: 'alpha',
			applicableSystems: ['alpha', 'EU', 'UK', 'US'],
		},
		{
			id: 'bottoms',
			category: 'bottoms',
			title: 'Bottoms',
			valueKind: 'numeric',
			applicableSystems: ['EU', 'UK', 'US', 'in'],
		},

		// Accessories.
		{
			id: 'hat',
			category: 'accessories',
			title: 'Hat',
			valueKind: 'alpha',
			applicableSystems: ['alpha', 'cm'],
		},
		{
			id: 'gloves',
			category: 'accessories',
			title: 'Gloves',
			valueKind: 'alpha',
			applicableSystems: ['alpha'],
		},
		{
			id: 'ring',
			category: 'accessories',
			title: 'Ring',
			valueKind: 'numeric',
			applicableSystems: ['US', 'UK', 'EU'],
		},
		{
			id: 'helmet',
			category: 'accessories',
			title: 'Helmet',
			valueKind: 'alpha',
			applicableSystems: ['alpha', 'cm'],
		},

		// Football kit.
		{
			id: 'football-jersey',
			category: 'sport-kit',
			title: 'Football jersey',
			valueKind: 'alpha',
			applicableSystems: ['alpha'],
		},
		{
			id: 'football-shorts',
			category: 'sport-kit',
			title: 'Football shorts',
			valueKind: 'alpha',
			applicableSystems: ['alpha'],
		},
		{
			id: 'football-socks',
			category: 'sport-kit',
			title: 'Football socks',
			valueKind: 'alpha',
			applicableSystems: ['alpha'],
		},
		{
			id: 'football-boots',
			category: 'sport-kit',
			title: 'Football boots',
			valueKind: 'numeric',
			applicableSystems: ['EU', 'UK', 'US'],
		},
		{
			id: 'football-goalkeeper-gloves',
			category: 'sport-kit',
			title: 'Football goalkeeper gloves',
			valueKind: 'alpha',
			applicableSystems: ['alpha'],
		},
		{
			id: 'football-shin-guards',
			category: 'sport-kit',
			title: 'Football shin guards',
			valueKind: 'alpha',
			applicableSystems: ['alpha'],
		},

		// Basketball kit.
		{
			id: 'basketball-jersey',
			category: 'sport-kit',
			title: 'Basketball jersey',
			valueKind: 'alpha',
			applicableSystems: ['alpha'],
		},
		{
			id: 'basketball-shorts',
			category: 'sport-kit',
			title: 'Basketball shorts',
			valueKind: 'alpha',
			applicableSystems: ['alpha'],
		},
		{
			id: 'basketball-shoes',
			category: 'sport-kit',
			title: 'Basketball shoes',
			valueKind: 'numeric',
			applicableSystems: ['EU', 'UK', 'US'],
		},

		// Cycling kit.
		{
			id: 'cycling-helmet',
			category: 'sport-kit',
			title: 'Cycling helmet',
			valueKind: 'alpha',
			applicableSystems: ['alpha', 'cm'],
		},
		{
			id: 'cycling-gloves',
			category: 'sport-kit',
			title: 'Cycling gloves',
			valueKind: 'alpha',
			applicableSystems: ['alpha'],
		},
		{
			id: 'cycling-jersey',
			category: 'sport-kit',
			title: 'Cycling jersey',
			valueKind: 'alpha',
			applicableSystems: ['alpha'],
		},

		// Swimming kit.
		{
			id: 'swimming-goggles',
			category: 'sport-kit',
			title: 'Swimming goggles',
			valueKind: 'alpha',
			applicableSystems: ['alpha'],
		},
		{
			id: 'swimming-fins',
			category: 'sport-kit',
			title: 'Swimming fins',
			valueKind: 'numeric',
			applicableSystems: ['EU', 'UK', 'US'],
		},
	],
	sportKits: [
		{
			id: 'football',
			title: 'Football kit',
			sizeTypeIDs: [
				'football-jersey',
				'football-shorts',
				'football-socks',
				'football-boots',
				'football-goalkeeper-gloves',
				'football-shin-guards',
			],
		},
		{
			id: 'basketball',
			title: 'Basketball kit',
			sizeTypeIDs: ['basketball-jersey', 'basketball-shorts', 'basketball-shoes'],
		},
		{
			id: 'cycling',
			title: 'Cycling kit',
			sizeTypeIDs: ['cycling-helmet', 'cycling-gloves', 'cycling-jersey'],
		},
		{
			id: 'swimming',
			title: 'Swimming kit',
			sizeTypeIDs: ['swimming-goggles', 'swimming-fins'],
		},
	],
};
