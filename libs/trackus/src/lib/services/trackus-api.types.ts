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
  readonly i?: number;
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
