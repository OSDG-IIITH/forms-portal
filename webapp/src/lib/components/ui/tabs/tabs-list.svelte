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
			const update = () => {
				requestAnimationFrame(() => {
					const activeButton = tabsContainer.querySelector(`[data-value="${$activeTab}"]`) as HTMLElement;
					if (activeButton) {
						const { left } = activeButton.getBoundingClientRect();
						const { left: containerLeft } = tabsContainer.getBoundingClientRect();
						const offsetLeft = left - containerLeft;
						const offsetWidth = activeButton.offsetWidth;
						activeIndicator.style.transform = `translateX(${offsetLeft}px)`;
						activeIndicator.style.width = `${offsetWidth}px`;
					}
				});
			};
			update();
			window.addEventListener('resize', update);
			const ro = new ResizeObserver(update);
			ro.observe(tabsContainer);
			return () => {
				window.removeEventListener('resize', update);
				ro.disconnect();
			};
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
