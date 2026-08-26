import { InjectionToken } from '@angular/core';

/**
 * A single community-centre navigation section surfaced by NoticeBoard.cc.
 * The app renders one side-menu item and one dashboard tile per section.
 */
export interface ICommunitycentrumSection {
  /** URL segment mounted under /space/:spaceType/:spaceID/ (e.g. 'rooms'). */
  readonly segment: string;
  /** Human-readable menu label. */
  readonly title: string;
  /** ionicons name for the menu / dashboard icon. */
  readonly icon: string;
}

/**
 * The always-on service contract for the communitycentrum extension.
 *
 * Kept intentionally small: NoticeBoard.cc is a thin composition layer over
 * existing Sneat extensions (Calendarius, Contactus, Assetus, Bookius,
 * ToGethered, Competios), so this surface only exposes the extension id and the
 * community-centre navigation sections. Consumers depend on this interface and
 * the token below, never on the runtime implementation class.
 */
export interface ICommunitycentrumService {
  /** Stable extension id, used for space-scoped API/DAL paths. */
  readonly extensionID: 'communitycentrum';
  /** The community-centre navigation sections, in menu order. */
  sections(): readonly ICommunitycentrumSection[];
}

/**
 * DI token binding {@link ICommunitycentrumService} to its concrete runtime
 * implementation. Provided once at the app composition root via
 * `provideCommunitycentrum()`.
 */
export const COMMUNITYCENTRUM_SERVICE =
  new InjectionToken<ICommunitycentrumService>('CommunitycentrumService');
