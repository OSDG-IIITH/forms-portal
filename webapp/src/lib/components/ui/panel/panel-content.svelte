<script lang="ts">
	import { getContext } from 'svelte';
	import { cn } from '$lib/utils.js';
	
	let { 
		children, 
		class: className = '', 
		...restProps 
	}: { 
		children?: any; 
		class?: string; 
		[key: string]: any 
	} = $props();

	let panelElement = $state<HTMLElement>();
	let triggerRect = $state<DOMRect | null>(null);
	let mounted = $state(false);
	let animating = $state(false);
	
	const panel = getContext<any>('panel');
	const open = $derived(panel?.open ?? false);
	const triggerElement = $derived(panel?.triggerElement ?? null);

	$effect(() => {
		if (triggerElement) {
			triggerRect = triggerElement.getBoundingClientRect();
		}
	});

	$effect(() => {
		if (open && !mounted) {
			mounted = true;
			setTimeout(() => {
				animating = true;
			}, 10);
		} else if (!open && mounted) {
			animating = false;
			setTimeout(() => {
				mounted = false;
			}, 300);
		}
	});

	function handleClickOutside(event: MouseEvent) {
		if (open && panelElement && !panelElement.contains(event.target as Node) && 
			triggerElement && !triggerElement.contains(event.target as Node)) {
			// Ignore clicks on checkboxes inside the panel
			if ((event.target as HTMLElement)?.closest('[data-slot="checkbox"]')) {
				return;
			}
			panel?.close();
		}
	}

	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'Escape' && open) {
			panel?.close();
		}
	}

	$effect(() => {
		if (open) {
			document.addEventListener('click', handleClickOutside);
			document.addEventListener('keydown', handleKeydown);
			return () => {
				document.removeEventListener('click', handleClickOutside);
				document.removeEventListener('keydown', handleKeydown);
			};
		}
	});
</script>

{#if mounted}
	<div
		bind:this={panelElement}
		class={cn(
			"fixed z-50 bg-card border border-border rounded-lg shadow-lg p-0 max-h-[80vh] overflow-y-auto",
			"w-90 2xl:w-100",
			"transition-transform duration-300 ease-in-out",
			animating ? "translate-x-0" : "translate-x-full",
			className
		)}
		style:top="{triggerRect ? triggerRect.bottom + 20 : 60}px"
		style:right="28px"
		{...restProps}
	>
		{@render children?.()}
	</div>
{/if}
