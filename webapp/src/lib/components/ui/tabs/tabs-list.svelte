<script lang="ts">
	import { getContext } from 'svelte';
	import { cn } from '$lib/utils';

	interface Props {
		children: any;
		class?: string;
	}

	let { children, class: className = '', ...restProps }: Props = $props();
	
	const { activeTab } = getContext('tabs') as { activeTab: any };
	
	let tabsContainer: HTMLDivElement;
	let activeIndicator: HTMLDivElement;
	
	$effect(() => {
		if (tabsContainer && activeIndicator && $activeTab) {
			const activeButton = tabsContainer.querySelector(`[data-value="${$activeTab}"]`) as HTMLElement;
			if (activeButton) {
				const { offsetLeft, offsetWidth } = activeButton;
				activeIndicator.style.transform = `translateX(${offsetLeft}px)`;
				activeIndicator.style.width = `${offsetWidth}px`;
			}
		}
	});
</script>

<div
	bind:this={tabsContainer}
	class={cn(
		"inline-flex h-10 items-center justify-center rounded-md bg-[var(--tab-track)] border border-[var(--tab-border)] p-1 text-muted-foreground relative overflow-hidden",
		className
	)}
	{...restProps}
>
	<div
		bind:this={activeIndicator}
		class="absolute inset-y-1 rounded-sm shadow-sm transition-all duration-300 ease-out"
		style="width: 0; transform: translateX(0); background: var(--tab-active);"
	></div>
	{@render children()}
</div>
