<script lang="ts">
	import FormItem from './form-item.svelte';

	export let forms: any[];
	export let loading: boolean;
	export let userHandle: string;
	export let animationDelay: number = 40;
</script>

{#if loading}
	<div class="space-y-2">
		{#each Array(10) as _}
			<div class="h-16 bg-muted rounded-lg animate-pulse"></div>
		{/each}
	</div>
{:else}
	<div class="space-y-0.5">
		{#each forms as form, i}
			<div class="animate-in fade-in-0 slide-in-from-left-4 duration-400 ease-out" style="animation-delay: {i * animationDelay}ms; animation-fill-mode: both;">
				<FormItem {form} variant="created" viewMode="list" {userHandle} />
			</div>
		{/each}
	</div>
{/if}
