<script lang="ts">
	import { base } from '$app/paths';
	import { Button } from '$lib/components/ui/button';
	import { IconMessage, IconTrash } from '@tabler/icons-svelte';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { page } from '$app/state';
	import { formatRelativeTime } from '$lib/utils/date';

	type Comment = {
		id: string;
		commenter: string;
		body: string;
		modified: string;
	};

	let comments = $state<Comment[]>([]);
	let newComment = $state('');
	let loading = $state(false);
	let error = $state('');
	let userNames = $state<Record<string, string>>({});

	const formId = $derived(page.params.formId || page.data.form?.id || '');

	let showDeleteDialog = $state(false);
	let commentToDelete = $state<string | null>(null);

	function openDeleteDialog(commentId: string) {
		commentToDelete = commentId;
		showDeleteDialog = true;
	}

	async function fetchUserName(userId: string): Promise<string> {
		if (!userId || userNames[userId]) return userNames[userId] || 'Unknown';
		try {
			const res = await fetch(`${base}/api/users/${userId}`);
			if (!res.ok) throw new Error();
			const data = await res.json();
			const name = data.name || data.handle || 'Unknown';
			userNames[userId] = name;
			return name;
		} catch {
			userNames[userId] = 'Unknown';
			return 'Unknown';
		}
	}

	async function fetchComments() {
		if (!formId) return;
		loading = true;
		error = '';
		try {
			const res = await fetch(`${base}/api/forms/${formId}/comments`);
			if (!res.ok) throw new Error('Failed to fetch comments');
			const data = await res.json();
			const commentList: Comment[] = data.data || [];
			comments = commentList;

			const uniqueIds = Array.from(new Set(commentList.map((c) => c.commenter).filter(Boolean)));
			await Promise.all(uniqueIds.map(fetchUserName));
		} catch (e: any) {
			error = e.message || 'Unknown error';
		} finally {
			loading = false;
		}
	}

	$effect(() => {
		if (formId) {
			fetchComments();
		}
	});

	async function addComment(event: SubmitEvent) {
		event.preventDefault();
		if (!newComment.trim() || !formId) return;
		loading = true;
		error = '';
		try {
			const res = await fetch(`${base}/api/forms/${formId}/comments`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ element: 'none', body: newComment })
			});
			if (!res.ok) throw new Error('Failed to add comment');
			newComment = '';
			await fetchComments();
		} catch (e: any) {
			error = e.message || 'Unknown error';
		} finally {
			loading = false;
		}
	}

	async function confirmDeleteComment() {
		if (!formId || !commentToDelete) return;
		loading = true;
		error = '';
		try {
			const res = await fetch(`${base}/api/forms/${formId}/comments/${commentToDelete}`, {
				method: 'DELETE',
				headers: { 'Content-Type': 'application/json' }
			});
			if (!res.ok) throw new Error('Failed to delete comment');
			showDeleteDialog = false;
			commentToDelete = null;
			await fetchComments();
		} catch (e: any) {
			error = e.message || 'Unknown error';
		} finally {
			loading = false;
		}
	}
</script>

<div class="flex flex-col h-full overflow-hidden">
	{#if error}
		<div class="px-8 text-red-500 text-sm mb-2">{error}</div>
	{/if}

	<div class="flex-1 overflow-y-auto px-8 pt-4 pb-4 space-y-3">
		{#if loading && comments.length === 0}
			<div class="text-muted-foreground text-sm">Loading...</div>
		{:else if comments.length === 0}
			<div class="text-muted-foreground text-sm min-h-40 flex justify-center items-center">
				No comments yet.
			</div>
		{:else}
			{#each comments as c (c.id)}
				<div class="group border rounded-md p-3 bg-muted/30 text-left relative">
					<div class="flex items-center justify-between mb-1">
						<span class="text-xs text-muted-foreground"
							>{userNames[c.commenter] || '...'} &middot;
							{formatRelativeTime(c.modified)}</span
						>
						<button
							type="button"
							class="opacity-0 group-hover:opacity-100 transition-opacity ml-2 text-muted-foreground hover:text-destructive p-1 rounded"
							aria-label="Delete comment"
							onclick={(e) => {
								e.stopPropagation();
								openDeleteDialog(c.id);
							}}
							tabindex="-1"
							disabled={loading}
						>
							<IconTrash class="w-4 h-4" />
						</button>
					</div>
					<div class="text-sm">{c.body}</div>
				</div>
			{/each}
		{/if}
	</div>

	<form
		class="sticky bottom-0 left-0 right-0 bg-background px-8 py-4 flex gap-2 border-t z-10"
		onsubmit={addComment}
	>
		<input
			class="h-11 flex-1 border rounded-md px-3 py-2 text-sm bg-background focus:outline-none"
			placeholder="Add a comment..."
			bind:value={newComment}
			autocomplete="off"
		/>
	</form>
</div>

<Dialog.Root bind:open={showDeleteDialog}>
	<Dialog.Content class="sm:max-w-[400px]">
		<Dialog.Header>
			<Dialog.Title>Delete comment</Dialog.Title>
			<Dialog.Description
				>This action cannot be undone. Are you sure you want to delete this comment?</Dialog.Description
			>
		</Dialog.Header>
		<Dialog.Footer>
			<Button
				type="button"
				variant="destructive"
				class="px-4 py-2 rounded-md"
				onclick={(e) => {
					e.stopPropagation();
					confirmDeleteComment();
				}}
				disabled={loading}
			>
				Delete
			</Button>
			<Button
				type="button"
				variant="ghost"
				class="px-4 py-2 rounded-md ml-2"
				onclick={(e) => {
					e.stopPropagation();
					showDeleteDialog = false;
					commentToDelete = null;
				}}
			>
				Cancel
			</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
