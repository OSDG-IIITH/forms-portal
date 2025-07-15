<script lang="ts">
	import { setContext } from 'svelte';
	import { writable } from 'svelte/store';

	interface Props {
		value?: string;
		onValueChange?: (value: string) => void;
		children: any;
		class?: string;
	}

	let { value = $bindable(''), onValueChange, children, class: className = '', ...restProps }: Props = $props();

	const activeTab = writable(value);

	setContext('tabs', {
		activeTab,
		setActiveTab: (newValue: string) => {
			activeTab.set(newValue);
			value = newValue;
			onValueChange?.(newValue);
		}
	});

	$effect(() => {
		activeTab.set(value);
	});
</script>

<div class={className} {...restProps}>
	{@render children()}
</div>
