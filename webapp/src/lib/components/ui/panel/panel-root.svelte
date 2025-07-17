<script lang="ts">
	import { setContext, createEventDispatcher } from 'svelte';
	
	let { children }: { children?: any } = $props();
	let open = $state(false);
	let triggerElement = $state<HTMLElement | null>(null);
	
	const dispatch = createEventDispatcher<{
		openchange: { open: boolean };
	}>();
	
	function toggle() {
		open = !open;
		dispatch('openchange', { open });
	}

	function close() {
		open = false;
		dispatch('openchange', { open });
	}

	const panelContext = {
		get open() { return open; },
		get triggerElement() { return triggerElement; },
		toggle,
		close,
		setTriggerElement: (el: HTMLElement) => triggerElement = el
	};

	setContext('panel', panelContext);
</script>

{@render children?.()}
