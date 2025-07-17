<script lang="ts">
	import { cn } from '$lib/utils.js';
	
	let { children, class: className = '', ...restProps }: { children?: any; class?: string; [key: string]: any } = $props();
</script>

<h2 class={cn("text-lg font-semibold text-foreground", className)} {...restProps}>
	{@render children?.()}
</h2>
