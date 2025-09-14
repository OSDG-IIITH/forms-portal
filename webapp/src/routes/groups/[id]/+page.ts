import { base } from '$app/paths';
import { error } from '@sveltejs/kit';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, params }) => {
	try {
		const groupResponse = await fetch(`${base}/api/groups/${params.id}`, { credentials: 'include' });
		if (!groupResponse.ok) {
			throw error(groupResponse.status, `Failed to fetch group with ID ${params.id}`);
		}
		const group = await groupResponse.json();

		// Fetch owner and all members in parallel
		const memberIds = group.members ? [group.owner, ...group.members] : [group.owner];
		const uniqueIds = [...new Set(memberIds)];

		const userPromises = uniqueIds.map((id) =>
			fetch(`${base}/api/users/${id}`, { credentials: 'include' }).then((res) => {
				if (!res.ok) {
					console.warn(`Failed to fetch user ${id}`);
					return null;
				}
				return res.json();
			})
		);

		const users = (await Promise.all(userPromises)).filter(Boolean);
		const userNameMap = new Map(users.map((user) => [user.id, user.name]));

		group.owner = userNameMap.get(group.owner) || group.owner;

		if (group.members) {
			group.members = group.members.map((id: string) => userNameMap.get(id) || id);
		}

		return {
			group
		};
	} catch (e: any) {
		if (e.status) {
			throw e;
		}
		throw error(500, e.message || 'Could not load group details');
	}
};
