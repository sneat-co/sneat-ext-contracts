import { InjectionToken } from '@angular/core';
import { Observable } from 'rxjs';
import {
  BookiusBookingDbo,
  BookiusBookingTypeDbo,
  BookiusCreateBookingRequest,
  BookiusCreateBookingTypeRequest,
} from '../dto';

// IBookiusService is the runtime-light contract used by Bookius UI components.
// The implementation lives in the internal lib and is provided via the
// BOOKIUS_SERVICE token below.
export interface IBookiusService {
  createBookingType(
    request: BookiusCreateBookingTypeRequest,
  ): Observable<BookiusBookingTypeDbo>;
  createBooking(request: BookiusCreateBookingRequest): Observable<BookiusBookingDbo>;
}

export const BOOKIUS_SERVICE = new InjectionToken<IBookiusService>(
  'BookiusService',
);
