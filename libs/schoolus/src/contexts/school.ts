import { ISpaceItemNavContext } from '@sneat/space-models';
import { IEnrolmentDbo, ISchoolClassDbo } from '../dto';

export type ISchoolClassContext = ISpaceItemNavContext<ISchoolClassDbo, ISchoolClassDbo>;
export type IEnrolmentContext = ISpaceItemNavContext<IEnrolmentDbo, IEnrolmentDbo>;
