import {
	createListInfoFromDto,
	getListShortUrlId,
	isListInfoMatchesListDto,
	ListItemInfoModel,
	ListItemModel,
} from './list';
import type { IListDbo, IListItemBrief } from './list';
import { describe, expect, it } from 'vitest';

describe('getListShortUrlId', () => {
  it('combines spaceId and shortId when a shortId is given', () => {
    expect(getListShortUrlId('fam1', 'groceries')).toBe('fam1-groceries');
  });

  it('falls back to the full id when there is no shortId', () => {
    expect(getListShortUrlId('fam1', undefined, 'to-do:abc')).toBe('to-do:abc');
  });

	it('returns undefined when neither shortId nor id is provided', () => {
		expect(getListShortUrlId('fam1')).toBeUndefined();
	});

	it('selects stable list-item identities in priority order', () => {
		expect(ListItemInfoModel.trackBy(3, undefined as unknown as IListItemBrief)).toBe(3);
		expect(ListItemInfoModel.trackBy(0, { id: 'item-1', title: 'Milk' })).toBe(
			'id:item-1',
		);
		expect(
			ListItemInfoModel.trackBy(0, {
				id: '',
				subListId: 'list-2',
				title: 'Dinner',
			}),
		).toBe('subList:list-2');
		expect(ListItemInfoModel.trackBy(0, { id: '', title: 'Milk' })).toBe('Milk');
	});

	it('compares identified and provisional list items by their stable fields', () => {
		expect(
			ListItemModel.equalListItems(
				{ id: 'item-1', title: 'Milk' },
				{ id: 'item-1', title: 'Updated title' },
			),
		).toBe(true);
		expect(
			ListItemModel.equalListItems(
				{ id: 'item-1', title: 'Milk' },
				{ id: 'item-2', title: 'Milk' },
			),
		).toBe(false);
		expect(
			ListItemModel.equalListItems(
				{ id: '', title: 'Milk', category: 'dairy' },
				{ id: '', title: 'Milk', category: 'dairy' },
			),
		).toBe(true);
		expect(
			ListItemModel.equalListItems(
				{ id: '', title: 'Milk', category: 'dairy' },
				{ id: '', title: 'Milk', category: 'other' },
			),
		).toBe(false);
	});

	it('matches list records by id or by type and short id', () => {
		const record = {
			id: 'list-1',
			dbo: { title: 'Groceries', type: 'buy', shortId: 'groceries' },
		} as Parameters<typeof isListInfoMatchesListDto>[1];

		expect(isListInfoMatchesListDto({ type: 'do', id: 'list-1' }, record)).toBe(true);
		expect(
			isListInfoMatchesListDto(
				{ type: 'buy', shortId: 'groceries' },
				record,
			),
		).toBe(true);
		expect(
			isListInfoMatchesListDto({ type: 'buy', shortId: 'other' }, record),
		).toBe(false);
	});

	it('creates the consumer list summary without dropping optional metadata', () => {
		const dto = {
			title: 'Groceries',
			type: 'buy',
			items: [{ id: 'item-1', title: 'Milk' }],
			emoji: '🛒',
			restrictions: { role: 'editor' },
			commune: { id: 'space-1', title: 'Family' },
		} as unknown as IListDbo;

		expect(createListInfoFromDto(dto, 'groceries')).toEqual({
			type: 'buy',
			title: 'Groceries',
			shortId: 'groceries',
			itemsCount: 1,
			emoji: '🛒',
			restrictions: dto.restrictions,
			space: dto.commune,
		});
		expect(() => createListInfoFromDto({ type: 'buy' } as IListDbo)).toThrow(
			'!title',
		);
	});
});
