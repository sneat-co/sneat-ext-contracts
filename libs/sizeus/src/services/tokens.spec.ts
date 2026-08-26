import { describe, expect, it } from 'vitest';
import { SIZEUS_API_SERVICE } from './sizeus-api';
import { SIZEUS_SERVICE } from './sizeus-service';

describe('Sizeus service tokens', () => {
	it('keep distinct, stable consumer injection identities', () => {
		expect(SIZEUS_API_SERVICE.toString()).toBe('InjectionToken SizeusApiService');
		expect(SIZEUS_SERVICE.toString()).toBe('InjectionToken SizeusService');
		expect(SIZEUS_API_SERVICE).not.toBe(SIZEUS_SERVICE);
	});
});
