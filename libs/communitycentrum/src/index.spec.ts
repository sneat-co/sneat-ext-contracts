import { InjectionToken } from '@angular/core';
import { COMMUNITYCENTRUM_SERVICE, ICommunitycentrumSection } from './index';

describe('communitycentrum contract', () => {
  it('exposes COMMUNITYCENTRUM_SERVICE as an InjectionToken', () => {
    expect(COMMUNITYCENTRUM_SERVICE).toBeInstanceOf(InjectionToken);
  });

  it('names the token CommunitycentrumService', () => {
    expect(COMMUNITYCENTRUM_SERVICE.toString()).toContain(
      'CommunitycentrumService',
    );
  });

  it('describes a navigation section shape (compile-time smoke)', () => {
    const section: ICommunitycentrumSection = {
      segment: 'rooms',
      title: 'Rooms',
      icon: 'albums',
    };
    expect(section.segment).toBe('rooms');
  });
});
