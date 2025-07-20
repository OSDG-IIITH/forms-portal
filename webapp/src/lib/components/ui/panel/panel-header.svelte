<script lang="ts">
	import { cn } from '$lib/utils.js';
	
	let { children, class: className = '', ...restProps }: { children?: any; class?: string; [key: string]: any } = $props();
</script>

<div class={cn("p-6 pb-0", className)} {...restProps}>
	{@render children?.()}
</div>
