import { IWithRelatedOnly, IWithSpaceIDs } from '@sneat/dto';
import { ActivityType, RepeatPeriod, WeekdayCode2 } from './happening-types';
import type { IScheduledResponsibilitySpec } from './responsibility';
import { IWithStringID } from './todo_move_funcs';

export interface ISlotParticipant {
  readonly roles?: string[];
  // readonly type: 'member' | 'contact';
  // readonly title: string;
}

/**
 * The role a contact holds ON a happening — the key stored in the related
 * item's `rolesOfItem` map. It answers "what is this person to this lesson?",
 * not "what may this person do in this Space".
 *
 * The vocabulary is CLOSED and calendarius-owned, mirroring the backend's
 * `dbo4calendarius.HappeningParticipantWellKnownRoles` exactly. Five existing
 * fleet vocabularies were considered first and each is a different AXIS:
 * `SpaceMemberRole` is space-membership authorization (a person can teach one
 * lesson and attend another in the same Space), `IRelationshipRoles`' kinship
 * constants are contact↔contact with opposite roles, and eventius'
 * `event-attendee`/`event-host` is a reciprocal edge that does not contain
 * `participant` — the value calendarius has already written everywhere.
 *
 * The two that DO contain `teacher` were rejected for the same reason plus a
 * layering one: `ContactRole` (contactus) classifies what a contact IS to a
 * space — its `teacher` lives inside `ContactRoleKidRelated`, beside `plumber`,
 * `ship` and `port_to_location`, none of which may name a lesson link; and
 * `SchoolRole` (schoolus) is a person's role WITHIN a school, from a module
 * that sits ABOVE calendarius and that calendarius must not depend on. The
 * `teacher` VALUE is deliberately identical in all three, so a schoolus
 * consumer maps `SchoolRoleTeacher` onto this role by identity.
 */
export type HappeningParticipantRole =
  | 'participant'
  | 'teacher'
  | 'organizer'
  | 'assistant';

/**
 * Every HappeningParticipantRole, in the backend's declared order. Use this to
 * render a role picker rather than hand-listing the union, so a role added
 * server-side shows up without a second edit.
 */
export const happeningParticipantRoles: readonly HappeningParticipantRole[] = [
  'participant',
  'teacher',
  'organizer',
  'assistant',
];

/**
 * The role a caller gets when it does not ask for one. The backend resolves an
 * absent role to exactly this, so a UI defaulting to it produces byte-identical
 * requests to the ones every pre-existing client sends.
 */
export const defaultHappeningParticipantRole: HappeningParticipantRole =
  'participant';

export function isKnownHappeningParticipantRole(
  role: string,
): role is HappeningParticipantRole {
  return (happeningParticipantRoles as readonly string[]).includes(role);
}

/**
 * A contact's roles on one happening.
 *
 * ADOPTED rather than deleted: this type existed with zero references and an
 * untyped `roles?: string[]`. Leaving it would have made a third shape once the
 * role vocabulary landed, so it is now the typed projection consumers read —
 * `Object.keys(relatedItem.rolesOfItem)` narrowed to the closed vocabulary.
 * `roles` is a LIST because the underlying linkage model stores a map of role
 * keys: a contact can legitimately be both the teacher and the organizer of a
 * lesson.
 */
export interface IHappeningParticipant {
  readonly roles?: readonly HappeningParticipantRole[];
  // readonly type: 'member' | 'contact';
  // readonly title: string;
}

/**
 * One contact named in an `add_participants` / `remove_participants` request,
 * optionally saying WHICH role that contact holds on the happening.
 *
 * Mirrors the backend's `dto4calendarius.HappeningContactRef`, which embeds
 * `ShortSpaceModuleItemRef` and adds `role` beside it — so `{id, spaceID}`
 * remains a valid body and `role` is purely additive.
 *
 * Only the REF is published here, not the whole request: the request envelope
 * (`spaceID` + `happeningID`) is already typed in the calendarius frontend
 * against `ISpaceRequest`, and duplicating it would create a second shape of
 * the same thing. The ref is the part the role rides on, and the part a
 * non-Angular consumer cannot infer.
 *
 * `role` omitted means `participant` — see `defaultHappeningParticipantRole`.
 * Sending a value outside the closed vocabulary is answered `400`; it is never
 * stored.
 */
export interface IHappeningContactRef {
  readonly id: string;
  /** Omit for a contact in the happening's own space. */
  readonly spaceID?: string;
  readonly role?: HappeningParticipantRole;
}

/**
 * Canonical Happening price-coverage term. This is independent of the
 * Happening recurrence cadence; `quarter` does not make an event recurring.
 */
export type TermUnit =
  | 'single'
  | 'second'
  | 'minute'
  | 'hour'
  | 'day'
  | 'week'
  | 'month'
  | 'quarter'
  | 'year';

export interface ITerm {
  readonly length: number;
  readonly unit: TermUnit;
}

export type CurrencyCode = 'USD' | 'EUR' | 'RUB' | string;

export interface IAmount {
  readonly currency: CurrencyCode;
  /** Integer minor units, matching decimal.Decimal64p2's canonical wire form. */
  readonly value: number;
}

export interface IHappeningPrice {
  readonly id: string;
  readonly term: ITerm;
  readonly amount: IAmount;
  readonly expenseQuantity?: number;
}

export const happeningPriceLimits = {
  idMaxBytes: 200,
  maxItems: 100,
} as const;

const happeningPriceTermUnits: readonly TermUnit[] = [
  'single',
  'second',
  'minute',
  'hour',
  'day',
  'week',
  'month',
  'quarter',
  'year',
];

const canonicalCurrencyCodes = new Set(
  (
    Intl as typeof Intl & {
      supportedValuesOf(key: 'currency'): readonly string[];
    }
  ).supportedValuesOf('currency'),
);

/**
 * Validates the existing Happening-owned price projection. Price item IDs are
 * the stable reference for consumers; multiple items may intentionally share
 * a term. This helper does not introduce an Event-specific pricing authority.
 */
export function assertValidHappeningPrices(
  prices: readonly IHappeningPrice[] | undefined,
): void {
  if (!prices) return;
  if (prices.length > happeningPriceLimits.maxItems)
    throw new Error(
      `prices exceeds maximum item count ${happeningPriceLimits.maxItems}`,
    );
  const seen = new Set<string>();
  for (const [index, price] of prices.entries()) {
    assertHappeningPriceID(`prices[${index}].id`, price.id);
    if (seen.has(price.id))
      throw new Error(`prices[${index}].id duplicates ${price.id}`);
    seen.add(price.id);
    if (!happeningPriceTermUnits.includes(price.term.unit))
      throw new Error(`prices[${index}].term has unknown unit`);
    if (!Number.isSafeInteger(price.term.length) || price.term.length < 1)
      throw new Error(`prices[${index}].term.length must be positive`);
    if (!canonicalCurrencyCodes.has(price.amount.currency))
      throw new Error(
        `prices[${index}].amount.currency must be a canonical ISO 4217 code`,
      );
    if (!Number.isSafeInteger(price.amount.value) || price.amount.value < 0)
      throw new Error(
        `prices[${index}].amount.value must be nonnegative safe-integer minor units`,
      );
    if (
      price.expenseQuantity !== undefined &&
      (!Number.isSafeInteger(price.expenseQuantity) ||
        price.expenseQuantity < 0)
    )
      throw new Error(
        `prices[${index}].expenseQuantity must be a nonnegative safe integer`,
      );
  }
}

function assertHappeningPriceID(field: string, value: string): void {
  if (!value) throw new Error(`${field} is required`);
  if (value.trim() !== value)
    throw new Error(`${field} must not have leading or trailing whitespace`);
  if (!isWellFormedUnicode(value))
    throw new Error(`${field} must encode as valid UTF-8`);
  if (new TextEncoder().encode(value).byteLength > happeningPriceLimits.idMaxBytes)
    throw new Error(
      `${field} exceeds maximum UTF-8 byte length ${happeningPriceLimits.idMaxBytes}`,
    );
  if (value === '*') throw new Error(`${field} must not be '*'`);
}

function isWellFormedUnicode(value: string): boolean {
  for (let i = 0; i < value.length; i++) {
    const unit = value.charCodeAt(i);
    if (unit >= 0xd800 && unit <= 0xdbff) {
      if (i + 1 >= value.length) return false;
      const next = value.charCodeAt(++i);
      if (next < 0xdc00 || next > 0xdfff) return false;
    } else if (unit >= 0xdc00 && unit <= 0xdfff) {
      return false;
    }
  }
  return true;
}

export interface IHappeningBase extends IWithRelatedOnly {
  readonly type: HappeningType;
  readonly status: HappeningStatus;
  readonly responsibility?: IScheduledResponsibilitySpec;
  readonly kind: HappeningKind;
  readonly activityType?: ActivityType; // TODO: Is it same as HappeningKind?
  readonly title: string;
  readonly summary?: string;
  readonly levels?: Level[];
  // readonly contactIDs?: readonly string[]; // obsolete
  readonly slots?: Readonly<Record<string, IHappeningSlot>>;
  readonly prices?: readonly IHappeningPrice[];
  // readonly participants?: Record<string, Readonly<IHappeningParticipant>>;
  /**
   * Filter facets on the happening — `guitar`, `piano`, `beginner-group`.
   *
   * Mirrors the backend's `with.TagsField`, which HappeningDbo has always
   * embedded; what was missing was every path TO it, this declaration included.
   * The server stores them normalised (trimmed, lower-cased, de-duplicated) and
   * enforces `happeningTagLimits`, so a client can render them verbatim.
   * Absent rather than `[]` when a happening has no tags.
   */
  readonly tags?: readonly string[];
  /**
   * Per-extension data embedded on the happening, keyed by extension id (e.g.
   * `eventus`). Lets a hosting module store its own fields on the happening
   * instead of a separate overlay document; calendarius stays agnostic to the
   * blob shapes. Mirrors the backend `HappeningBase.Ext`.
   */
  readonly ext?: Readonly<Record<string, unknown>>;
}

export type IHappeningBrief = IHappeningBase;

/**
 * The server-enforced tag rules, published so a form can refuse a bad tag
 * before a round trip rather than after a 400. These MIRROR the backend
 * (`dbo4calendarius.HappeningTagsMaxCount` / `HappeningTagMaxRunes`); the
 * server remains the authority.
 */
export const happeningTagLimits = {
  maxCount: 10,
  /** RUNES, not bytes — a Cyrillic or CJK tag gets the same 32 characters. */
  maxRunes: 32,
} as const;

/**
 * Canonicalises tags exactly as the server does before storing them: trimmed,
 * lower-cased, blanks dropped, duplicates collapsed, order of first appearance
 * preserved. A UI that normalises before sending shows the user the value that
 * will actually be stored.
 */
export function normalizeHappeningTags(
  tags: readonly string[] | undefined,
): string[] {
  if (!tags?.length) return [];
  const seen = new Set<string>();
  const normalized: string[] = [];
  for (const raw of tags) {
    const tag = raw.trim().toLowerCase();
    if (!tag || seen.has(tag)) continue;
    seen.add(tag);
    normalized.push(tag);
  }
  return normalized;
}

/**
 * Returns why `tag` is not an acceptable happening tag, or undefined when it
 * is. Takes the ALREADY-NORMALIZED value, so it never complains about casing or
 * padding a caller cannot see.
 */
export function happeningTagError(tag: string): string | undefined {
  if (!tag) return 'a tag must not be empty';
  // Spread, not .length: a surrogate pair is one character to a user.
  if ([...tag].length > happeningTagLimits.maxRunes)
    return `a tag must be at most ${happeningTagLimits.maxRunes} characters`;
  // C0 (U+0000\u2013U+001F), DEL (U+007F) and C1 (U+0080\u2013U+009F). The C1 range
  // is not decoration: the server rejects with Go's unicode.IsControl, which
  // covers it, so a narrower client guard lets `gui\u0085tar` pass here and 400 at
  // the server \u2014 breaking the one promise this mirror makes, that it enforces
  // exactly what the server enforces.
  // eslint-disable-next-line no-control-regex
  if (/[\u0000-\u001f\u007f-\u009f]/u.test(tag))
    return 'a tag must not contain control characters';
  return undefined;
}

export interface IWithDates {
  readonly dates?: string[];
}

export interface IWithSpaceDates extends IWithSpaceIDs, IWithDates {
  readonly spaceDates?: string[]; // ISO date strings prefixed with spaceID e.g. [`abc123:2019-12-01`, `abc123:2019-12-02`]
}

export interface IHappeningDbo extends IHappeningBrief, IWithSpaceDates {
  readonly description?: string;
}

export function validateHappeningDto(dto: IHappeningDbo): void {
  if (!dto.title) {
    throw new Error('happening has no title');
  }
  if (dto.title !== dto.title.trim()) {
    throw new Error(
      'happening title has leading or closing whitespace characters',
    );
  }
  switch (dto.type) {
    case 'single':
      break;
    case 'recurring':
      break;
    default:
      if (!dto.type) {
        throw new Error('happening has no type');
      }
      throw new Error('happening has unknown type: ' + dto.type);
  }
  if (!dto.type) {
    throw new Error('happening has no type');
  }
  const slots = Object.entries(dto.slots || {});
  const isPlannedSingleEvent = dto.type === 'single' && dto.kind === 'event';
  if (!slots.length && !isPlannedSingleEvent) {
    throw new Error('!dto.slots?.length');
  }
  switch (dto.type) {
    case 'single':
      slots.forEach(([slotID, slot]) => {
        if (isPlannedSingleEvent) {
          validatePlannedEventSlot(slotID, slot);
        } else {
          validateSingleHappeningSlot(slotID, slot);
        }
      });
      break;
    case 'recurring':
      slots.forEach(([slotID, slot]) =>
        validateRecurringHappeningSlot(slotID, slot),
      );
      break;
  }
}

export function validateRecurringHappeningSlot(
  slotID: string,
  slot: IHappeningSlot,
): void {
  if (slot.repeats === 'once' || slot.repeats === 'UNKNOWN') {
    throw new Error(
      `slots[${slotID}]: slot.repeats is not valid for recurring happening: ${slot.repeats}`,
    );
  }
  validateHappeningSlot(slotID, slot);
}

export function validateSingleHappeningSlot(
  slotID: string,
  slot: IHappeningSlot,
): void {
  if (slot.repeats != 'once') {
    throw new Error(
      `slots[${slotID}]: slot repeats is not 'once': ${slot.repeats}`,
    );
  }
  validateHappeningSlot(slotID, slot);
}

// validatePlannedEventSlot accepts the independently known parts of a
// one-time event plan. A title-only event has no slot; when a slot exists it
// must contribute at least a date, time, location or duration.
export function validatePlannedEventSlot(
  slotID: string,
  slot: IHappeningSlot,
): void {
  if (slot.repeats !== 'once') {
    throw new Error(
      `slots[${slotID}]: planned event slot repeats is not 'once': ${slot.repeats}`,
    );
  }
  if (
    !slot.start?.date &&
    !slot.start?.time &&
    !slot.end?.date &&
    !slot.end?.time &&
    !slot.durationMinutes &&
    !slot.location?.title &&
    !slot.location?.address
  ) {
    throw new Error(`slots[${slotID}]: planned event slot has no planning data`);
  }
  if ((slot.end?.time || slot.durationMinutes) && !slot.start?.time) {
    throw new Error(
      `slots[${slotID}]: planned event end or duration requires a start time`,
    );
  }
}

function validateHappeningSlot(slotID: string, slot: IHappeningSlot): void {
  if (
    !slot.start?.time &&
    !(slot.repeats.startsWith('monthly') || slot.repeats.startsWith('yearly'))
  ) {
    throw new Error(`slots[${slotID}]: slot has no start time: ${slot}`);
  }
}

export type HappeningType = 'recurring' | 'single';

export type HappeningStatus = 'draft' | 'active' | 'canceled' | 'archived';

export type HappeningKind = 'appointment' | 'activity' | 'task' | 'event';

export interface SlotLocation {
  readonly title?: string;
  readonly address?: string;
}

interface IFortnightly {
  readonly title: string;
}

/*
// tslint:disable-next-line:no-magic-numbers
type MonthlyDay = -5 | -4 | -3 | -2 | -1
// tslint:disable-next-line:no-magic-numbers
	| 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10
// tslint:disable-next-line:no-magic-numbers
	| 11 | 12 | 13 | 14 | 15 | 16 | 17 | 18 | 19
// tslint:disable-next-line:no-magic-numbers
	| 20 | 21 | 22 | 23 | 24 | 25 | 26 | 27 | 28;
*/

export interface IDateTime {
  readonly date?: string;
  readonly time?: string;
}

export interface ITiming {
  readonly start?: IDateTime;
  readonly end?: IDateTime;
  readonly durationMinutes?: number;
}

export interface IHappeningSlotSingleRef {
  readonly repeats: RepeatPeriod;
  readonly weekday?: WeekdayCode2;
  readonly week?: number;
}

export type Month =
  | 'January'
  | 'February'
  | 'March'
  | 'April'
  | 'May'
  | 'June'
  | 'July'
  | 'August'
  | 'September'
  | 'October'
  | 'November'
  | 'December';

export interface IHappeningSlotTiming extends ITiming {
  readonly repeats: RepeatPeriod;
  readonly weekdays?: readonly WeekdayCode2[];
  readonly day?: number;
  readonly month?: Month;
  readonly weeks?: readonly number[];
  readonly fortnightly?: Readonly<{
    readonly odd: IFortnightly;
    readonly even: IFortnightly;
  }>;
}

export type Level = 'beginners' | 'intermediate' | 'advanced';

export interface IHappeningTask {
  readonly serviceProvider?: {
    readonly id: string;
    readonly title: string;
  };
}

export interface IHappeningSlot extends IHappeningSlotTiming {
  readonly location?: SlotLocation;
  readonly groupIds?: string[]; // TODO: What is this?
}

export type IHappeningSlotWithID = IWithStringID<IHappeningSlot>;

export const emptyTiming: ITiming = {
  // durationMinutes: 0,
};

export const emptyHappeningSlot: IHappeningSlotWithID = {
  id: '',
  repeats: 'UNKNOWN',
  ...emptyTiming,
};

export interface ISingleHappeningDbo extends IHappeningDbo {
  readonly dtStarts?: number; // UTC
  readonly dtEnds?: number; // UTC
  readonly weekdays?: WeekdayCode2[];
}

export interface DtoSingleTask extends ISingleHappeningDbo {
  readonly isCompleted: boolean;
  readonly completion?: number; // In percents, max value is 100.
}
