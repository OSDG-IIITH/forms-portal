<script lang="ts">
	import {
		Dialog,
		DialogContent,
		DialogFooter,
		DialogHeader,
		DialogTitle,
		DialogTrigger
	} from '$lib/components/ui/dialog';
	import * as Select from '$lib/components/ui/select';
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

	const groupTypes = [
		{ value: 'list', label: 'List' },
		{ value: 'domain', label: 'Domain' }
	];

	const selectedLabel = $derived(groupTypes.find((t) => t.value === type)?.label);

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
		<Button class="ml-auto gap-2 h-8" variant="default" aria-label="Add group"> Add Group </Button>
	</DialogTrigger>
	<DialogContent class="max-w-sm w-full">
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
				<Select.Root type="single" bind:value={type}>
					<Select.Trigger class="w-full">
						{selectedLabel}
					</Select.Trigger>
					<Select.Content>
						{#each groupTypes as groupType}
							<Select.Item value={groupType.value}>{groupType.label}</Select.Item>
						{/each}
					</Select.Content>
				</Select.Root>
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