<script lang="ts">
	import { getContext } from 'svelte';
	import { cn } from '$lib/utils';

	interface Props {
		value: string;
		children: any;
		class?: string;
	}

	let { value, children, class: className = '', ...restProps }: Props = $props();

	const { activeTab, setActiveTab } = getContext('tabs') as {
		activeTab: any;
		setActiveTab: (value: string) => void;
	};

	const isActive = $derived($activeTab === value);
</script>

<button
	data-value={value}
	class={cn(
		"flex items-center justify-center h-full w-full px-3 py-0 text-sm font-medium leading-normal whitespace-nowrap rounded-sm ring-offset-background transition-all duration-300 ease-out focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 relative z-10",
		isActive ? "text-primary" : "text-muted-foreground hover:text-primary",
		className
	)}
	onclick={() => setActiveTab(value)}
	{...restProps}
>
	{@render children()}
</button>
