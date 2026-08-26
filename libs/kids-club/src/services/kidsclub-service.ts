import { InjectionToken } from '@angular/core';
import type { Observable } from 'rxjs';
import type { ICreateClubRequest, ICreateClubResponse, IClubPublic } from '../dto/club.js';
import type { IRegistrationRequest, IRegistrationResponse, IRosterEntry } from '../dto/registration.js';
export interface IKidsclubService { createClub(spaceID: string, request: ICreateClubRequest): Observable<ICreateClubResponse>; getClubPublic(clubID: string): Observable<IClubPublic>; register(clubID: string, request: IRegistrationRequest): Observable<IRegistrationResponse>; getRoster(spaceID: string, clubID: string): Observable<IRosterEntry[]>; }
export const KIDSCLUB_SERVICE = new InjectionToken<IKidsclubService>('KIDSCLUB_SERVICE');
