import { base } from '$app/paths';
import { error } from '@sveltejs/kit';
import type { PageLoad } from './$types';
import type { Group } from '$lib/components/groups/columns';

export const load: PageLoad = async ({ fetch }) => {
	try {
		const groupsResponse = await fetch(`${base}/api/groups`, { credentials: 'include' });

		if (!groupsResponse.ok) {
			throw error(groupsResponse.status, 'Failed to fetch groups');
		}

		const groupsData = await groupsResponse.json();
		const groups: Group[] = groupsData.data ?? [];

		const ownerIds = [...new Set(groups.map((group) => group.owner))];

		const userPromises = ownerIds.map((id) =>
			fetch(`${base}/api/users/${id}`, { credentials: 'include' }).then((res) => {
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
			pagination: groupsData.pagination
		};
	} catch (e: any) {
		if (e.status) {
			throw e;
		}
		throw error(500, e.message || 'Could not load groups');
	}
};
