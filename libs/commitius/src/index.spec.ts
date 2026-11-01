import { Injector } from '@angular/core';
import { describe, expect, it } from 'vitest';
import * as template from '@sneat/extension-template-contract';
import * as commitius from './index';
import type { ITemplateService, ITemplateSpaceDbo } from '@sneat/extension-template-contract';
import type { ICommitiusService, ICommitiusSpaceDbo } from './index';

// Keep the aliases assignable in both directions, not merely similarly named.
export const asTemplateService = (value: ICommitiusService): ITemplateService => value;
export const asCommitiusService = (value: ITemplateService): ICommitiusService => value;
export const asTemplateSpace = (value: ICommitiusSpaceDbo): ITemplateSpaceDbo => value;
export const asCommitiusSpace = (value: ITemplateSpaceDbo): ICommitiusSpaceDbo => value;

describe('Commitius public boundary', () => {
  it('preserves every runtime template reexport by identity', () => {
    for (const [name, value] of Object.entries(template)) {
      expect((commitius as Record<string, unknown>)[name]).toBe(value);
    }
  });
  it('keeps the public service token injectable independently of its template', () => {
    const service = { marker: 'host-owned Commitius service' };
    const injector = Injector.create({ providers: [{ provide: commitius.COMMITIUS_SERVICE, useValue: service }] });
    expect(injector.get(commitius.COMMITIUS_SERVICE)).toBe(service);
  });
});
