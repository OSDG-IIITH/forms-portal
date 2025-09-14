import { base } from '$app/paths';
import type { LayoutLoad } from './$types';
import { redirect } from '@sveltejs/kit';

export const ssr = false

const PUBLIC_PATHS = [`${base}/api/auth/login`, `${base}/api/auth/login/callback`, `${base}/api/auth/logout`];

export const load: LayoutLoad = async ({ fetch, url }) => {
	if (PUBLIC_PATHS.some((p) => url.pathname.startsWith(p))) {
		return {};
	}
	const res = await fetch(`${base}/api/auth/info`);
	if (res.status !== 200) {
		throw redirect(302, `${base}/api/auth/login`);
	}
	return {};
};
