export interface ICreateClubRequest { title: string; venue?: string; consentText: string; }
export interface ICreateClubResponse { id: string; publicPath: string; }
export interface IClubPublic { id: string; title: string; venue?: string; consentVersion: string; consentText: string; }
