import { describe, expect, it } from 'vitest';
import { roleCatalog, sportCatalog } from './index.js';

describe('sportCatalog', () => {
  it('uses unique stable codes and localisation keys', () => {
    const ids = sportCatalog.map((sport) => sport.id);
    expect(new Set(ids).size).toBe(ids.length);
    expect(sportCatalog.every((sport) => sport.labelKey.startsWith('sportius.sport.'))).toBe(true);
    expect(ids).toContain('other');
    expect(ids).toEqual(expect.arrayContaining(['chess', 'table-tennis', 'running', 'multi-sport']));
  });
});

describe('roleCatalog', () => {
  it('uses unique stable codes and localisation keys', () => {
    const ids = roleCatalog.map((role) => role.id);
    expect(new Set(ids).size).toBe(ids.length);
    expect(roleCatalog.every((role) => role.labelKey.startsWith('sportius.role.'))).toBe(true);
  });

  it('keeps the default personal selector compact', () => {
    const defaults = roleCatalog.filter((role) => role.defaultPersonal);
    expect(defaults.length).toBeGreaterThan(0);
    expect(defaults.length).toBeLessThanOrEqual(8);
  });

  it('allows general operational experience on a personal sport profile', () => {
    for (const id of ['assistant-coach', 'team-manager']) {
      expect(roleCatalog.find((role) => role.id === id)?.scopes).toContain('personal');
    }
  });
});
