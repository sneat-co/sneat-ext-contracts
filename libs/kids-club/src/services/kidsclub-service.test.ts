import { describe, expect, it } from 'vitest';
import { KIDSCLUB_SERVICE } from './kidsclub-service.js';

describe('KIDSCLUB_SERVICE', () => {
  it('keeps the public dependency-injection token stable', () => {
    expect(KIDSCLUB_SERVICE.toString()).toBe('InjectionToken KIDSCLUB_SERVICE');
  });
});
