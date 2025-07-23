<script lang="ts">
	import {
		Dialog,
		DialogContent,
		DialogFooter,
		DialogHeader,
		DialogTitle,
		DialogTrigger
	} from '$lib/components/ui/dialog';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';

	let { onCreated = () => {} }: { onCreated?: () => void } = $props();

	let open = $state(false);
	let name = $state('');
	let description = $state('');
	let type = $state<'domain' | 'list'>('list');
	let domain = $state('');
	let members = $state('');
	let loading = $state(false);
	let error = $state<string | null>(null);

	function resetForm() {
		name = '';
		description = '';
		domain = '';
		members = '';
		type = 'list';
		error = null;
	}

	async function handleSubmit() {
		loading = true;
		error = null;
		const payload: any = { name, type };
		if (description) payload.description = description;
		if (type === 'domain' && domain) payload.domain = domain;
		if (type === 'list' && members) payload.members = members.split(/[,\s]+/).filter(Boolean);

		try {
			const res = await fetch('/api/groups', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				credentials: 'include',
				body: JSON.stringify(payload)
			});

			if (!res.ok) {
				const errText = await res.text();
				throw new Error(errText || 'Failed to create group');
			}

			open = false;
			onCreated();
			resetForm();
		} catch (e: any) {
			error = e.message || 'Unknown error';
		} finally {
			loading = false;
		}
	}
</script>

<Dialog bind:open>
	<DialogTrigger>
		<Button class="ml-auto gap-2" variant="default" aria-label="Add group"> Add Group </Button>
	</DialogTrigger>
	<DialogContent>
		<DialogHeader>
			<DialogTitle>Add New Group</DialogTitle>
		</DialogHeader>
		<form onsubmit={handleSubmit} class="space-y-4">
			<div>
				<label class="block text-sm font-medium mb-1" for="name">Name</label>
				<Input id="name" bind:value={name} required placeholder="Group name" />
			</div>
			<div>
				<label class="block text-sm font-medium mb-1" for="description">Description</label>
				<Input id="description" bind:value={description} placeholder="Description (optional)" />
			</div>
			<div>
				<label class="block text-sm font-medium mb-1" for="type">Type</label>
				<select
					id="type"
					bind:value={type}
					class="w-full border rounded px-2 py-1 bg-background text-foreground h-10"
				>
					<option value="list">List</option>
					<option value="domain">Domain</option>
				</select>
			</div>

			{#if type === 'domain'}
				<div>
					<label class="block text-sm font-medium mb-1" for="domain">Domain</label>
					<Input id="domain" bind:value={domain} required placeholder="example.com" />
				</div>
			{/if}
			{#if type === 'list'}
				<div>
					<label class="block text-sm font-medium mb-1" for="members"
						>Members (comma or space separated emails)</label
					>
					<Input
						id="members"
						bind:value={members}
						placeholder="user1@example.com, user2@example.com"
					/>
				</div>
			{/if}

			{#if error}
				<div class="text-red-500 text-sm">{error}</div>
			{/if}

			<DialogFooter>
				<Button type="submit" disabled={loading}>{loading ? 'Creating...' : 'Create Group'}</Button>
			</DialogFooter>
		</form>
	</DialogContent>
</Dialog>