<script lang="ts">
	import { setContext, createEventDispatcher } from 'svelte';
	
	let { children, open: controlledOpen }: { children?: any; open?: boolean } = $props();
	let internalOpen = $state(false);
	let triggerElement = $state<HTMLElement | null>(null);
	
	const dispatch = createEventDispatcher<{
		openchange: { open: boolean };
	}>();
	
	const open = $derived(controlledOpen !== undefined ? controlledOpen : internalOpen);
	
	function toggle() {
		if (controlledOpen !== undefined) {
			dispatch('openchange', { open: !open });
		} else {
			internalOpen = !internalOpen;
			dispatch('openchange', { open: internalOpen });
		}
	}

	function close() {
		if (controlledOpen !== undefined) {
			dispatch('openchange', { open: false });
		} else {
			internalOpen = false;
			dispatch('openchange', { open: false });
		}
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
