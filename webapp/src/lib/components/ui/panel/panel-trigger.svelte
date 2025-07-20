<script lang="ts">
	import { getContext, onMount } from 'svelte';
	
	let { children, ...restProps }: { children?: any; [key: string]: any } = $props();
	let buttonElement: HTMLElement;
	
	const panel = getContext<any>('panel');
	
	onMount(() => {
		if (panel && buttonElement) {
			panel.setTriggerElement(buttonElement);
		}
	});
</script>

<button bind:this={buttonElement} {...restProps} onclick={panel?.toggle} type="button">
	{@render children?.()}
</button>
