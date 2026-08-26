import { InjectionToken } from '@angular/core';
import { Observable } from 'rxjs';
import { ISizeRecord } from '../size-record';
import { SizeSystem } from '../size-system';

// Request to record a size for a contact of a Space. Mirrors the backend
// facade4sizeus.CreateSizeRecordRequest: it intentionally has NO createdAt
// field — the creation timestamp is stamped server-side (or by the in-memory
// implementation standing in for the server).
export interface ICreateSizeRecordRequest {
  readonly spaceID: string;
  readonly contactID: string;
  readonly sizeTypeID: string;

  // Stored verbatim in the declared `system` — never normalized.
  readonly value: string;
  readonly system: SizeSystem;

  // ISO date string (YYYY-MM-DD). The UI defaults it to "today".
  readonly effectiveDate: string;

  readonly note?: string;
  readonly preferredBrandNote?: string;
}

// A retrieved size record with its ID — mirrors facade4sizeus.SizeRecordItem.
export interface ISizeRecordItem {
  readonly id: string;
  readonly record: ISizeRecord;
}

// ISizeusApiService is the runtime-light contract the Sizes UI depends on.
// Its members mirror the backend facade4sizeus endpoints (CreateSizeRecord,
// GetContactSizes, GetContactSizeHistory, GetContactCurrentSizes,
// GetSpaceSizes) so a later task can swap the in-memory implementation
// (internal lib) for the real HTTP/Firestore-backed one without touching any
// consumer. Records are persisted under /spaces/{spaceID}/ext/sizeus/... and
// keyed to the contact ID.
export interface ISizeusApiService {
  createSizeRecord(
    request: ICreateSizeRecordRequest,
  ): Observable<ISizeRecordItem>;

  // All size records of one contact of a Space (all size types).
  getContactSizes(
    spaceID: string,
    contactID: string,
  ): Observable<readonly ISizeRecordItem[]>;

  // Dated history of one contact's records for a single size type.
  getContactSizeHistory(
    spaceID: string,
    contactID: string,
    sizeTypeID: string,
  ): Observable<readonly ISizeRecordItem[]>;

  // The current record per size type for one contact (per the
  // REQ size-history "current" rule — see currentSizeRecord()).
  getContactCurrentSizes(
    spaceID: string,
    contactID: string,
  ): Observable<readonly ISizeRecordItem[]>;

  // All size records of a Space (any contact) — used by the space-level
  // Sizes overview to list the contacts that have recorded sizes.
  getSpaceSizes(spaceID: string): Observable<readonly ISizeRecordItem[]>;
}

export const SIZEUS_API_SERVICE = new InjectionToken<ISizeusApiService>(
  'SizeusApiService',
);
