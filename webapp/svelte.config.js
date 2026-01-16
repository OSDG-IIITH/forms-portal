import adapter from '@sveltejs/adapter-static';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

const config = {
	preprocess: vitePreprocess(),
	kit: {
		adapter: adapter({
			fallback: 'index.html',
			strict: false
		}),
		prerender: {
			entries: ['/'],
			handleUnseenRoutes: 'ignore'
		}
	}
};

export default config;
