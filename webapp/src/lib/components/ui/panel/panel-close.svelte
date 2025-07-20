<script lang="ts">
	import { getContext } from 'svelte';
	
	let { children, ...restProps }: { children?: any; [key: string]: any } = $props();
	
	const panel = getContext<any>('panel');
</script>

<button {...restProps} onclick={panel?.close} type="button">
	{@render children?.()}
</button>
