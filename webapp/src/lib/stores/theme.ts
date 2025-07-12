import { browser } from '$app/environment';
import { writable } from 'svelte/store';

export type Theme = 'milk' | 'lilac' | 'mint' | 'coal' | 'amethyst' | 'moss';

const THEME_KEY = 'selected-theme';

function getInitialTheme(): Theme {
	if (browser) {
		const savedTheme = localStorage.getItem(THEME_KEY) as Theme | null;
		if (savedTheme) return savedTheme;

		const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
		return prefersDark ? 'coal' : 'milk';
	}

	return 'milk';
}

function createThemeStore() {
	const { subscribe, set } = writable<Theme>(getInitialTheme());

	function applyTheme(theme: Theme) {
		if (browser) {
			document.documentElement.classList.remove('dark', 'theme-lilac', 'theme-mint', 'theme-amethyst', 'theme-moss');

			if (theme === 'coal') {
				document.documentElement.classList.add('dark');
			} else if (theme === 'amethyst') {
				document.documentElement.classList.add('dark');
				document.documentElement.classList.add('theme-lilac');
			} else if (theme === 'moss') {
				document.documentElement.classList.add('dark');
				document.documentElement.classList.add('theme-mint');
			} else if (theme === 'lilac') {
				document.documentElement.classList.add('theme-lilac');
			} else if (theme === 'mint') {
				document.documentElement.classList.add('theme-mint');
			}

			localStorage.setItem(THEME_KEY, theme);
		}
	}

	function setTheme(theme: Theme) {
		applyTheme(theme);
		set(theme);
	}

	if (browser) {
		const theme = getInitialTheme();
		applyTheme(theme);

		window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', (event) => {
			const currentTheme = localStorage.getItem(THEME_KEY);
			if (!currentTheme) {
				const newTheme = event.matches ? 'coal' : 'milk';
				setTheme(newTheme);
			}
		});
	}

	return {
		subscribe,
		set: setTheme,
	};
}

export const theme = createThemeStore();
