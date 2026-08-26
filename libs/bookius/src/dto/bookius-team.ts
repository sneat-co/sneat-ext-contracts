import { BookiusBookingBrief, BookiusBookingTypeBrief } from './booking';

export interface IBookiusSpaceDbo {
  bookingTypes?: BookiusBookingTypeBrief[];
  bookings?: BookiusBookingBrief[];
}
