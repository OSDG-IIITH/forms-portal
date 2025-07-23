<script lang="ts">
	import type { PageProps } from './$types';
	import { Badge } from '$lib/components/ui/badge';

	let { data }: PageProps = $props();
	const { group } = data;
</script>

<div class="container mx-auto min-h-screen pt-24">
	<div class="max-w-4xl mx-auto">
		<header class="mb-8">
			<div class="flex items-center justify-between mb-2">
				<h1 class="text-4xl font-bold tracking-tight">{group.name}</h1>
				<Badge variant="outline" class="capitalize">{group.type}</Badge>
			</div>
			{#if group.description}
				<p class="text-lg text-muted-foreground">{group.description}</p>
			{/if}
			<p class="text-sm text-muted-foreground mt-2">
				Owned by: <span class="font-medium">{group.owner}</span>
			</p>
		</header>

		<div class="rounded-md border p-6">
			{#if group.type === 'domain'}
				<div>
					<h3 class="text-lg font-semibold mb-2">Associated Domain</h3>
					<p class="font-mono bg-secondary text-secondary-foreground rounded p-2 inline-block">
						{group.domain}
					</p>
					<p class="text-sm text-muted-foreground mt-2">
						All users with an email address at this domain are automatically members of this group.
					</p>
				</div>
			{:else if group.type === 'list' && group.members && group.members.length > 0}
				<div>
					<h3 class="text-lg font-semibold mb-4">Members ({group.members.length})</h3>
					<ul class="space-y-2">
						{#each group.members as member}
							<li class="font-mono text-sm">{member}</li>
						{/each}
					</ul>
				</div>
			{:else}
				<p class="text-muted-foreground">This group has no members.</p>
			{/if}
		</div>
	</div>
</div>