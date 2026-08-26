import { InjectionToken } from '@angular/core';
import type { Observable } from 'rxjs';

export type SpaceFriendshipRole = 'friend' | 'neighbour';

export interface ICreateFriendshipInviteRequest {
	readonly spaceID: string;
	readonly role?: SpaceFriendshipRole;
}

export interface ICreateFriendshipInviteResponse {
	readonly id: string;
	readonly pin: string;
}

export interface IAcceptFriendshipInviteRequest {
	readonly inviteID: string;
	readonly pin: string;
	readonly toSpaceID: string;
}

export interface IAcceptFriendshipInviteResponse {
	readonly invitingSpaceID: string;
	readonly toSpaceID: string;
	readonly role: SpaceFriendshipRole;
}

export interface IFriendSpace {
	readonly id: string;
	readonly roles: readonly SpaceFriendshipRole[];
	readonly title?: string;
}

export interface IListFriendSpacesResponse {
	readonly friendSpaces?: readonly IFriendSpace[];
}

export interface IRemoveFriendshipRequest {
	readonly spaceID: string;
	readonly friendSpaceID: string;
}

export interface IYardiusFriendshipService {
	createFriendshipInvite(
		request: ICreateFriendshipInviteRequest,
	): Observable<ICreateFriendshipInviteResponse>;
	acceptFriendshipInvite(
		request: IAcceptFriendshipInviteRequest,
	): Observable<IAcceptFriendshipInviteResponse>;
	listFriendSpaces(spaceID: string): Observable<IListFriendSpacesResponse>;
	removeFriendship(request: IRemoveFriendshipRequest): Observable<void>;
}

export const YARDIUS_FRIENDSHIP_SERVICE =
	new InjectionToken<IYardiusFriendshipService>('YardiusFriendshipService');

export const BEFRIEND_SPACE_PAGE_PATH = 'befriend';

export function friendshipInvitePath(inviteID: string, pin: string): string {
	return `${BEFRIEND_SPACE_PAGE_PATH}?${new URLSearchParams({ id: inviteID, pin }).toString()}`;
}
