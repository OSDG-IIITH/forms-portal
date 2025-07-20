<script lang="ts">
	import FormItem from './form-item.svelte';
	import { createEventDispatcher } from 'svelte';

	const dispatch = createEventDispatcher();

	export let forms: any[];
	export let loading: boolean;
	export let userHandle: string;
	export let animationDelay: number = 40;

	function handleDelete(e: CustomEvent<{ id: string | number }>) {
		dispatch('delete', e.detail);
	}
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
				<FormItem {form} variant="created" viewMode="list" {userHandle} class="w-full" on:delete={handleDelete} />
			</div>
		{/each}
	</div>
{/if}
