import { error } from '@sveltejs/kit';
import type { PageLoad } from './$types';
import type { Group } from '$lib/types/group';

function readInt(value: string | null, fallback: number, min = 1, max = 100) {
	const parsed = Number.parseInt(value ?? '', 10);
	if (!Number.isFinite(parsed)) return fallback;
	return Math.min(max, Math.max(min, parsed));
}

export const load: PageLoad = async ({ fetch, url }) => {
	try {
		const page = readInt(url.searchParams.get('page'), 1, 1, Number.MAX_SAFE_INTEGER);
		const limit = readInt(url.searchParams.get('limit'), 20);
		const offset = (page - 1) * limit;
		const groupsResponse = await fetch(`/api/groups?limit=${limit}&offset=${offset}`, { credentials: 'include' });

		if (!groupsResponse.ok) {
			throw error(groupsResponse.status, 'Failed to fetch groups');
		}

		const groupsData = await groupsResponse.json();
		const groups: Group[] = groupsData.data ?? [];

		const ownerIds = [...new Set(groups.map((group) => group.owner))];

		const userPromises = ownerIds.map((id) =>
			fetch(`/api/users/${id}`, { credentials: 'include' }).then((res) => {
				if (!res.ok) {
					console.warn(`Failed to fetch user ${id}`);
					return null;
				}
				return res.json();
			})
		);

		const users = (await Promise.all(userPromises)).filter(Boolean);

		const ownerNameMap = new Map(users.map((user) => [user.id, user.name]));

		const groupsWithNames = groups.map((group) => ({
			...group,
			owner: ownerNameMap.get(group.owner) || group.owner
		}));

		return {
			groups: groupsWithNames,
			pagination: groupsData.pagination,
			page,
			limit
		};
	} catch (e: any) {
		if (e.status) {
			throw e;
		}
		throw error(500, e.message || 'Could not load groups');
	}
};