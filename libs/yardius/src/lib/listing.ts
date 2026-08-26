import { InjectionToken } from '@angular/core';
import type { Observable } from 'rxjs';

import type { BorrowState } from './borrow-state.js';
import type { GiveawayState } from './giveaway-state.js';
import type { ListingRequestStatus } from './listing-request.js';
import type { ListingVisibility } from './listing-visibility.js';
import type { IGetListingsResponse } from './listings-endpoint.js';

export const LISTING_TYPES = ['borrow', 'giveaway'] as const;
export type ListingType = (typeof LISTING_TYPES)[number];
export type ListingState = BorrowState | GiveawayState;

export function isOpenListingState(state: ListingState): boolean {
	return state !== 'closed';
}

export function listingTypeLabel(type: ListingType): string {
	return type === 'borrow' ? 'Borrow' : 'Give away';
}

export function listingStateLabel(state: ListingState): string {
	return state.charAt(0).toUpperCase() + state.slice(1);
}

export function listingRequestStatusLabel(status: ListingRequestStatus): string {
	return status.charAt(0).toUpperCase() + status.slice(1);
}

export interface IListingDbo {
	readonly assetID: string;
	readonly type: ListingType;
	readonly state: ListingState;
	readonly visibility: ListingVisibility;
	readonly createdAt?: string;
	readonly createdBy?: string;
}

interface IListingIdentity {
	readonly id: string;
	readonly spaceID: string;
}

export interface IBorrowListingDto extends IListingDbo, IListingIdentity {
	readonly type: 'borrow';
	readonly state: BorrowState;
}

export interface IGiveawayListingDto extends IListingDbo, IListingIdentity {
	readonly type: 'giveaway';
	readonly state: GiveawayState;
}

export type IListingDto = IBorrowListingDto | IGiveawayListingDto;

export interface IPublishListingRequest {
	readonly spaceID: string;
	readonly assetID: string;
	readonly type: ListingType;
	readonly visibility?: ListingVisibility;
}

export interface IListingActionResponse {
	readonly id: string;
	readonly listing: IListingDbo;
}

export interface IListingActionRequest {
	readonly spaceID: string;
	readonly listingID: string;
}

export interface IListingRequestActionRequest extends IListingActionRequest {
	readonly requestID: string;
}

export interface IApproveRequestResponse extends IListingActionResponse {
	readonly declinedRequestIDs?: readonly string[];
}

export interface IRequestToBorrowRequest {
	readonly spaceID: string;
	readonly ownerSpaceID: string;
	readonly listingID: string;
}

export interface IWithdrawRequestRequest extends IRequestToBorrowRequest {
	readonly requestID: string;
}

export type AvailableListingSource = 'member' | 'friend';

export interface IAvailableListingDto {
	readonly spaceID: string;
	readonly spaceTitle?: string;
	readonly listingID: string;
	readonly assetID: string;
	readonly type: ListingType;
	readonly state: ListingState;
	readonly visibility: ListingVisibility;
	readonly source: AvailableListingSource;
}

export interface IMyRequestDto {
	readonly id: string;
	readonly ownerSpaceID: string;
	readonly ownerSpaceTitle?: string;
	readonly listingID: string;
	readonly assetID: string;
	readonly type: ListingType;
	readonly status: ListingRequestStatus;
	readonly createdAt: string;
}

export interface IGetMyRequestsResponse {
	readonly requests: readonly IMyRequestDto[];
}

export interface IGetListingResult {
	readonly listing?: IListingDto;
	readonly requests?: readonly import('./listing-request.js').IListingRequestDto[];
}

export interface IYardiusListingService {
	publishListing(request: IPublishListingRequest): Observable<IListingActionResponse>;
	cancelListing(request: IListingActionRequest): Observable<IListingActionResponse>;
	requestToBorrow(request: IRequestToBorrowRequest): Observable<IListingActionResponse>;
	claimListing(request: IRequestToBorrowRequest): Observable<IListingActionResponse>;
	withdrawRequest(request: IWithdrawRequestRequest): Observable<IListingActionResponse>;
	approveRequest(request: IListingRequestActionRequest): Observable<IApproveRequestResponse>;
	declineRequest(request: IListingRequestActionRequest): Observable<IListingActionResponse>;
	confirmHandover(request: IListingActionRequest): Observable<IListingActionResponse>;
	confirmReturn(request: IListingActionRequest): Observable<IListingActionResponse>;
	getAvailableListings(): Observable<IGetListingsResponse>;
	getOurListings(spaceID: string): Observable<readonly IListingDto[]>;
	getListing(spaceID: string, listingID: string): Observable<IGetListingResult>;
	getMyRequests(): Observable<IGetMyRequestsResponse>;
}

export const YARDIUS_LISTING_SERVICE =
	new InjectionToken<IYardiusListingService>('YardiusListingService');
