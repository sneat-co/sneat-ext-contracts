export interface IDocTypeListItem {
  readonly id: string;
  readonly title: string;
  readonly emoji?: string;
}

/** Sensible groupings for the new-document type picker. */
export type DocTypeCategory = 'identity' | 'family' | 'vehicle' | 'education' | 'legal';

/** A named contact slot a document type requires/allows (e.g. "spouse 1"). */
export interface IDocContactRoleDef {
  readonly id: string;
  readonly label: string;
  /** Minimum number of contacts that must fill this role. */
  readonly min: number;
  /** Maximum number of contacts that may fill this role. */
  readonly max: number;
}

export type DocFieldType = 'text' | 'date' | 'number';

/** A typed field beyond the standard assetus title/number/issuedOn/validTill. */
export interface IDocFieldDef {
  readonly id: string;
  readonly label: string;
  readonly type: DocFieldType;
  readonly required?: boolean;
}

export interface IDocTypeSchema {
  readonly id: string;
  readonly title: string;
  readonly emoji?: string;
  readonly category: DocTypeCategory;
  readonly contactRoles: readonly IDocContactRoleDef[];
  readonly fields: readonly IDocFieldDef[];
}
