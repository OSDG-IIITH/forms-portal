import type { LayoutLoad } from './$types';
import { redirect } from '@sveltejs/kit';

const PUBLIC_PATHS = ['/api/auth/login', '/api/auth/login/callback', '/api/auth/logout'];

export const load: LayoutLoad = async ({ fetch, url }) => {
	if (PUBLIC_PATHS.some((p) => url.pathname.startsWith(p))) {
		return {};
	}
	const res = await fetch('/api/auth/info');
	if (res.status !== 200) {
		throw redirect(302, '/api/auth/login');
	}
	return {};
};
