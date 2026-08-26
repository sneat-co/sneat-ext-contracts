import { InjectionToken } from '@angular/core';
import type {
  ITemplateService,
  ITemplateSpaceDbo,
} from '@sneat/extension-template-contract';

export * from '@sneat/extension-template-contract';

// Commitius currently specializes the maintained list extension. Keeping these
// aliases in its contract gives hosts a Commitius-owned API boundary while the
// shared list model and service surface remain defined in one package.
export type ICommitiusService = ITemplateService;
export type ICommitiusSpaceDbo = ITemplateSpaceDbo;

export const COMMITIUS_SERVICE = new InjectionToken<ICommitiusService>('COMMITIUS_SERVICE');
