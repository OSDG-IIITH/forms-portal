import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig, loadEnv } from 'vite';

export default defineConfig(({ mode }) => {
	const env = loadEnv(mode, '.', '');
	return {
		plugins: [tailwindcss(), sveltekit()],
		define: {
			'import.meta.env.VITE_BACKEND_URL': JSON.stringify(env.VITE_BACKEND_URL)
		}
	};
});
