export interface IChild { name: string; birthYear: number; gender?: 'f' | 'm'; allergies?: string; }
export interface IConsentAcceptance { version: string; acceptedAt: string; byUID: string; }
export interface IRegistrationRequest { children: IChild[]; emergencyContactName: string; emergencyContactPhone: string; pickups?: string[]; acceptedConsentVersion: string; }
export interface IRegistrationResponse { id: string; }
export interface IRosterEntry { id: string; children: IChild[]; emergencyContactName: string; emergencyContactPhone: string; pickups?: string[]; consent: IConsentAcceptance; createdAt: string; }
