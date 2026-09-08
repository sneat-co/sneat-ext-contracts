export type MediaAccess = 'public' | 'private';
export type MediaAssetStatus = 'uploading' | 'ready' | 'orphaned' | 'purging' | 'deleted' | 'failed';
export type MediaLinkStatus = 'active' | 'deleted';
export type MediaLinkRetention = 'follows_source' | 'retained';
export type MediaTargetScope = 'root' | 'space';
export type MediaVariant = 'avatar-xs' | 'avatar-sm' | 'avatar' | 'avatar-lg' | 'thumbnail' | 'preview' | 'large';

export interface IMediaCrop {
  readonly centerX: number;
  readonly centerY: number;
  readonly zoom: number;
}

export interface IMediaRef {
  readonly mediaID: string;
  readonly crop?: IMediaCrop;
}

export interface IMediaTarget {
  readonly scope: MediaTargetScope;
  readonly spaceID?: string;
  readonly type: 'user' | 'contact' | 'asset' | 'list_item';
  readonly id: string;
  readonly parentID?: string;
}

export interface IBeginMediaUploadRequest {
  readonly requestID: string;
  readonly contentType: 'image/jpeg' | 'image/png';
  readonly size: number;
  readonly originalFilename?: string;
  readonly access: MediaAccess;
}

export interface IBeginMediaUploadResponse {
  readonly mediaID: string;
  readonly uploadURL: string;
  readonly uploadMethod: 'POST';
  readonly uploadHeaders: Readonly<Record<string, string>>;
  readonly expiresAt: string;
}

export interface IFinalizeMediaUploadRequest {
  readonly mediaID: string;
  readonly sha256: string;
}

export interface ILinkMediaRequest {
  readonly requestID: string;
  readonly mediaID: string;
  readonly target: IMediaTarget;
  readonly role: 'avatar' | 'photo' | 'preview' | 'attachment';
  readonly crop?: IMediaCrop;
  readonly retention?: MediaLinkRetention;
}

export interface IMediaLinkResult {
  readonly mediaID: string;
  readonly linkID: string;
  readonly refCount: number;
  readonly durableRefCount: number;
  readonly status: MediaLinkStatus;
  readonly purgeAfter?: string;
}

export interface IMediaAccessRequest {
  readonly mediaID: string;
  readonly target: IMediaTarget;
  readonly role: ILinkMediaRequest['role'];
}

export interface IMediaAccessResponse {
  readonly token: string;
  readonly expiresAt: string;
}

export interface IMediaPresentation {
  readonly ref: IMediaRef;
  readonly access: MediaAccess;
  readonly token?: string;
  readonly tokenExpiresAt?: string;
}

export const mediaURL = (mediaID: string, variant: MediaVariant, token?: string): string => {
  const path = `https://media.sneat.co/m/${encodeURIComponent(mediaID)}/${variant}`;
  return token ? `${path}?token=${encodeURIComponent(token)}` : path;
};
