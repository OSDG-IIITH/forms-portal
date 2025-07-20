<script lang="ts">
import * as DropdownMenu from "$lib/components/ui/dropdown-menu/index.js";
import { Button, buttonVariants } from "$lib/components/ui/button/index.js";
import * as Dialog from "$lib/components/ui/dialog/index.js";
import { Input } from "$lib/components/ui/input/index.js";
import { Label } from "$lib/components/ui/label/index.js";
import { IconDotsVertical } from '@tabler/icons-svelte';
import { goto } from '$app/navigation';
import { createEventDispatcher } from 'svelte';

export let userHandle: string;
export let slug: string;
export let formId: string;

const dispatch = createEventDispatcher();

function goTo(path: string) {
	goto(path);
}

let showRenameDialog = false;
let showDeleteDialog = false;
let newTitle = '';
let error = '';

async function submitRename() {
	if (!newTitle.trim()) {
		error = 'Title cannot be empty.';
		return;
	}
	try {
		const res = await fetch(`/api/forms/${formId}`, {
			method: 'PATCH',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ title: newTitle })
		});
		if (!res.ok) {
			const data = await res.json().catch(() => ({}));
			throw new Error(data.message || 'Failed to rename form');
		}
		showRenameDialog = false;
		dispatch('renamed', { title: newTitle });
	} catch (e) {
		if (e instanceof Error) {
			error = e.message || 'Failed to rename form';
		} else {
			error = 'Failed to rename form';
		}
	}
}

async function confirmDelete() {
	try {
		const res = await fetch(`/api/forms/${formId}`, { method: 'DELETE' });
		if (res.status === 204) {
			showDeleteDialog = false;
			dispatch('deleted');
		} else {
			const data = await res.json().catch(() => ({}));
			throw new Error(data.message || 'Failed to delete form');
		}
	} catch (e) {
		showDeleteDialog = false;
		if (e instanceof Error) {
			alert(e.message || 'Failed to delete form');
		} else {
			alert('Failed to delete form');
		}
	}
}
</script>

<DropdownMenu.Root>
	<DropdownMenu.Trigger>
		<Button variant="ghost" size="icon" aria-label="Open menu" class="hover:bg-muted text-muted-foreground">
			<IconDotsVertical class="w-4 h-4" />
		</Button>
	</DropdownMenu.Trigger>
	<DropdownMenu.Content class="w-48 rounded-md shadow-xl border border-border text-xs" align="end" style="min-width:12rem;">
		<DropdownMenu.Group>
			<DropdownMenu.Item class="text-xs rounded-md px-3 py-2 hover:bg-accent/80 hover:text-primary transition-colors cursor-pointer" onclick={() => goTo(`/${userHandle}/${slug}`)}>
				Respond to form
			</DropdownMenu.Item>
			<DropdownMenu.Item class="text-xs rounded-md px-3 py-2 hover:bg-accent/80 hover:text-primary transition-colors cursor-pointer" onclick={() => goTo(`/${userHandle}/${slug}/edit`)}>
				Edit form
			</DropdownMenu.Item>
			<DropdownMenu.Item class="text-xs rounded-md px-3 py-2 hover:bg-accent/80 hover:text-primary transition-colors cursor-pointer" onclick={() => goTo(`/${userHandle}/${slug}/responses`)}>
				View responses
			</DropdownMenu.Item>
			<DropdownMenu.Item class="text-xs rounded-md px-3 py-2 hover:bg-accent/80 hover:text-primary transition-colors cursor-pointer" onclick={() => goTo(`/${userHandle}/${slug}/permissions`)}>
				Manage permissions
			</DropdownMenu.Item>
		</DropdownMenu.Group>
		<DropdownMenu.Separator />
		<DropdownMenu.Item class="text-xs rounded-md px-3 py-2 hover:bg-accent/80 hover:text-primary transition-colors cursor-pointer" onclick={() => { newTitle = ''; error = ''; showRenameDialog = true; }}>
			Rename form
		</DropdownMenu.Item>
		<DropdownMenu.Item class="text-xs rounded-md px-3 py-2 hover:bg-destructive/80 hover:text-destructive transition-colors cursor-pointer" onclick={() => { showDeleteDialog = true; }}>
			Delete form
		</DropdownMenu.Item>
	</DropdownMenu.Content>
</DropdownMenu.Root>

<!-- rename dialog -->
<Dialog.Root bind:open={showRenameDialog}>
	<Dialog.Content class="sm:max-w-[425px]">
		<Dialog.Header>
			<Dialog.Title>Rename form</Dialog.Title>
			<Dialog.Description>Enter a new name for this form.</Dialog.Description>
		</Dialog.Header>
		<div class="grid gap-4 py-4">
			<Input bind:value={newTitle} placeholder="New form title" onkeydown={(e) => e.key === 'Enter' && submitRename()} />
			{#if error}
				<span class="text-xs text-destructive">{error}</span>
			{/if}
		</div>
		<Dialog.Footer>
			<Button type="button" onclick={submitRename}>Save</Button>
			<Button type="button" variant="ghost" onclick={() => showRenameDialog = false}>Cancel</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>

<!-- delete dialog -->
<Dialog.Root bind:open={showDeleteDialog}>
	<Dialog.Content class="sm:max-w-[400px]">
		<Dialog.Header>
			<Dialog.Title>Delete form</Dialog.Title>
			<Dialog.Description>This action cannot be undone. Are you sure you want to delete this form?</Dialog.Description>
		</Dialog.Header>
		<Dialog.Footer>
			<Button type="button" variant="destructive" onclick={confirmDelete}>Delete</Button>
			<Button type="button" variant="ghost" onclick={() => showDeleteDialog = false}>Cancel</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
