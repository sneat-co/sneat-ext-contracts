// Self-contained Vitest setup for the contactus extension libs.
//
// In sneat-libs these specs run via `@nx/vitest:test` against a shared
// `@sneat/core/testing` harness. That subpath is not part of the published
// `@sneat/core` package, so on extraction we inline the parts the specs need:
//  1. initialise the Angular testing environment (TestBed) via Analog's
//     setup-testbed (zone-based, matching the libs' zone.js runtime), and
//  2. install the jsdom shims that stop Ionic/Stencil from throwing on asset
//     loading (`URL`, `matchMedia`, `CSSStyleSheet`, `document.baseURI`).
// The specs supply their own service mocks via `TestBed.configureTestingModule`,
// so no global Sneat provider mocks are required here.
/* eslint-disable @typescript-eslint/no-explicit-any */
import '@analogjs/vitest-angular/setup-zone';
import { setupTestBed } from '@analogjs/vitest-angular/setup-testbed';

setupTestBed({ zoneless: false });

if (typeof window !== 'undefined') {
  const OriginalURL = (window as any).URL;
  (window as any).URL = class extends OriginalURL {
    constructor(url: string, base?: string | URL) {
      try {
        super(url, base);
      } catch {
        try {
          super(url, 'http://localhost/');
        } catch {
          super(
            url.startsWith('/') ? `http://localhost${url}` : 'http://localhost/',
          );
        }
      }
    }
  };

  if ((window as any).document) {
    Object.defineProperty((window as any).document, 'baseURI', {
      get: () => 'http://localhost/',
      configurable: true,
    });
    if ((window as any).document.dir === undefined) {
      (window as any).document.dir = '';
    }
  }

  if (!window.matchMedia) {
    Object.defineProperty(window, 'matchMedia', {
      writable: true,
      value: (query: string) => ({
        matches: false,
        media: query,
        onchange: null,
        addListener: () => undefined,
        removeListener: () => undefined,
        addEventListener: () => undefined,
        removeEventListener: () => undefined,
        dispatchEvent: () => false,
      }),
    });
  }

  if (!(window as any).CSSStyleSheet) {
    (window as any).CSSStyleSheet = class {
      replaceSync() {
        /* ignore */
      }
      replace() {
        return Promise.resolve();
      }
    };
  } else {
    (window as any).CSSStyleSheet.prototype.replaceSync = function () {
      /* ignore */
    };
    (window as any).CSSStyleSheet.prototype.replace = function () {
      return Promise.resolve();
    };
  }

  if (!(window as any).CSS) {
    (window as any).CSS = { supports: () => false };
  }
}
