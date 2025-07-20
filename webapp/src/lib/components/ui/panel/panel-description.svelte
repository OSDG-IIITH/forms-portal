<script lang="ts">
	import { cn } from '$lib/utils.js';
	
	let { children, class: className = '', ...restProps }: { children?: any; class?: string; [key: string]: any } = $props();
</script>

<p class={cn("text-sm text-muted-foreground", className)} {...restProps}>
	{@render children?.()}
</p>
