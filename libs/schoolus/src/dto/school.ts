import { IWithCreated, IWithSpaceIDs } from '@sneat/dto';

export type SchoolType =
  | 'language' | 'music' | 'tutoring' | 'primary' | 'secondary'
  | 'international' | 'dance' | 'driving' | 'sports' | 'training';

export interface ISchoolProfileDbo extends IWithSpaceIDs, IWithCreated {
  readonly schoolType: SchoolType;
  title: string;
  timezone: string;
  currentAcademicYearID?: string;
  academicYearIDs?: string[];
  campuses?: ICampus[];
}

export interface ICampus {
  readonly id: string;
  title: string;
  address?: string;
  buildings?: IBuilding[];
}

export interface IBuilding { readonly id: string; title: string; campusID?: string; }

export interface IClassroom {
  readonly id: string;
  title: string;
  buildingID?: string;
  capacity?: number;
  assetSpaceRef?: string;
  assetIDs?: string[];
}
