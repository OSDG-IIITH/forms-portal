<script lang="ts">
	import { getContext } from 'svelte';
	import { cn } from '$lib/utils';

	interface Props {
		value: string;
		children: any;
		class?: string;
	}

	let { value, children, class: className = '', ...restProps }: Props = $props();

	const { activeTab } = getContext('tabs') as {
		activeTab: any;
	};

	const isActive = $derived($activeTab === value);
</script>

{#if isActive}
	<div
		class={cn(
			"mt-2 ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2",
			className
		)}
		{...restProps}
	>
		{@render children()}
	</div>
{/if}
