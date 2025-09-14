<script lang="ts">
	import { base } from '$app/paths';
	import DataTable from '$lib/components/groups/data-table.svelte';
	import { columns } from '$lib/components/groups/columns';
	import { invalidate } from '$app/navigation';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	let loading = $state(false);

	async function handleGroupCreated() {
		loading = true;
		try {
			await invalidate(`${base}/api/groups`);
		} catch (error) {
			console.error('Failed to reload groups:', error);
		} finally {
			loading = false;
		}
	}
</script>

<div class="container mx-auto px-8 py-8 min-h-screen">
	<div class="flex items-center justify-between mb-6">
		<h2 class="text-3xl font-bold tracking-tight">Groups</h2>
	</div>

	{#if loading}
		<div class="text-center py-8 text-muted-foreground">Loading...</div>
	{:else if data.groups}
		<DataTable data={data.groups} {columns} onGroupCreated={handleGroupCreated} />
	{:else}
		<div class="text-center py-8 text-muted-foreground">No groups found.</div>
	{/if}
</div>
