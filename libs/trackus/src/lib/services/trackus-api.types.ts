import { ITrackerBrief } from '../dbo/i-tracker-dbo';

export interface ITrackerRequest {
  spaceID: string;
  trackerID: string;
}

export interface ICreateTrackerRequest extends ITrackerBrief {
  readonly spaceID: string;
}

export interface IAddTrackerPointRequest {
  readonly spaceID: string;
  readonly trackerID: string;
  readonly trackByKind: 'space' | 'contact' | 'asset';
  readonly trackByID: string;
  /** Retain across retries of the same point; the backend upserts this ID. */
  readonly entryID?: string;
  /** ISO/RFC3339 event instant, not a Firestore Timestamp. */
  readonly ts?: string;
  readonly n?: string;
  readonly i?: number;
  readonly f?: number;
  readonly s?: string;
  readonly b?: boolean;
  readonly m?: number;
  readonly I?: readonly number[];
  readonly F?: readonly number[];
}

export interface IDeleteTrackerPointsRequest {
  readonly spaceID: string;
  readonly trackerID: string;
  readonly entityRef?: string;
  readonly date?: string;
  readonly pointIDs?: string[];
}

export interface IAddTrackerPointResponse {
  readonly entryID: string;
}

export interface ICreateTrackerResponse {
  trackerID: string;
}
