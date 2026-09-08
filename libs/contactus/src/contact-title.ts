import { IPersonNames } from '@sneat/auth-models';
import { IIdAndOptionalBriefAndOptionalDbo } from '@sneat/core';
import { IContactBrief, IContactDbo } from './dto';

export type ContactPersonNameStyle = 'display' | 'member';

export interface IContactTitleOptions {
  /** An explicit caller-provided title that precedes persisted contact data. */
  readonly preferredTitle?: string;
  /** Preserve callers that historically displayed brief.shortTitle. */
  readonly includeBriefShortTitle?: boolean;
  /** Select the established UI or member-list structured-name precedence. */
  readonly personNameStyle?: ContactPersonNameStyle;
  /** Defaults to true. Disable when an opaque contact ID is not useful copy. */
  readonly fallbackToID?: boolean;
  /** Returned when no configured contact value can be displayed. */
  readonly fallback?: string;
}

export function contactTitle(
  contact:
    | IIdAndOptionalBriefAndOptionalDbo<IContactBrief, IContactDbo>
    | undefined,
  options: IContactTitleOptions = {},
): string {
  const brief = contact?.brief;
  return (
    options.preferredTitle ||
    brief?.title ||
    (options.includeBriefShortTitle ? brief?.shortTitle : undefined) ||
    contact?.dbo?.title ||
    contactPersonName(brief?.names, options.personNameStyle ?? 'display') ||
    (options.fallbackToID === false ? undefined : contact?.id) ||
    options.fallback ||
    ''
  );
}

export function contactPersonName(
  names: IPersonNames | undefined,
  style: ContactPersonNameStyle = 'display',
): string | undefined {
  if (!names) {
    return undefined;
  }
  if (style === 'member') {
    return (
      [names.firstName, names.lastName].filter(Boolean).join(' ') ||
      names.fullName
    );
  }
  if (names.fullName) {
    return names.fullName;
  }
  if (
    names.firstName &&
    names.lastName &&
    !names.nickName &&
    !names.middleName
  ) {
    return `${names.firstName} ${names.lastName}`;
  }
  return (
    names.nickName ||
    names.firstName ||
    names.lastName ||
    names.middleName
  );
}
