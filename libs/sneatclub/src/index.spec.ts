import { SNEATCLUB_SERVICE, ListPage } from './index';

describe('@sneat/extension-sneatclub-contract public entry', () => {
  it('exposes the SNEATCLUB_SERVICE injection token', () => {
    expect(SNEATCLUB_SERVICE).toBeDefined();
    expect(SNEATCLUB_SERVICE.toString()).toContain('SneatclubService');
  });

  it('exposes the ListPage enum', () => {
    expect(ListPage.list).toBe('list');
  });
});
