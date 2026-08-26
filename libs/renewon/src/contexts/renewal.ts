import { ISpaceItemNavContext } from '@sneat/space-models';
import { IRenewalBrief, IRenewalDbo } from '../dto/renewal';

export type IRenewalContext = ISpaceItemNavContext<IRenewalBrief, IRenewalDbo>;
