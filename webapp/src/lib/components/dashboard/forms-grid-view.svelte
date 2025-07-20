<script lang="ts">
	import FormItem from './form-item.svelte';
	import { createEventDispatcher } from 'svelte';

	const dispatch = createEventDispatcher();

	export let forms: any[];
	export let loading: boolean;
	export let userHandle: string;
	export let animationDelay: number = 75;

	function handleDelete(e: CustomEvent<{ id: string | number }>) {
		dispatch('delete', e.detail);
	}
</script>

{#if loading}
	<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5 gap-4">
		{#each Array(10) as _}
			<div class="h-36 bg-muted rounded-lg animate-pulse"></div>
		{/each}
	</div>
{:else}
	<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5 gap-4">
		{#each forms as form, i}
			<div class="animate-in fade-in-0 slide-in-from-bottom-4 duration-500 ease-out" style="animation-delay: {i * animationDelay}ms; animation-fill-mode: both;">
				<FormItem {form} variant="created" viewMode="grid" {userHandle} on:delete={handleDelete} />
			</div>
		{/each}
	</div>
{/if}
