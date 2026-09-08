import { describe, expect, it } from 'vitest';
import { mediaURL } from './index';

describe('mediaURL', () => {
  it('uses the stable semantic path and escapes token text', () => {
    expect(mediaURL('m-1', 'thumbnail', 'a+b')).toBe(
      'https://media.sneat.co/m/m-1/thumbnail?token=a%2Bb',
    );
  });
});
