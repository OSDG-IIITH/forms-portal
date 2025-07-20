<script lang="ts">
  import { onMount } from 'svelte';
  import { Button } from '$lib/components/ui/button';
  import { IconMessage, IconSend2, IconTrash } from '@tabler/icons-svelte';
  import * as Dialog from "$lib/components/ui/dialog/index.js";
  let showDeleteDialog = false;
  let commentToDelete: string | null = null;

  function openDeleteDialog(commentId: string) {
    commentToDelete = commentId;
    showDeleteDialog = true;
  }

  async function confirmDeleteComment() {
    if (!formId || !commentToDelete) return;
    loading = true;
    error = '';
    try {
      const res = await fetch(`/api/forms/${formId}/comments/${commentToDelete}`, {
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
  import { page } from '$app/stores';
  import { derived } from 'svelte/store';
  import { formatRelativeTime } from '$lib/utils/date';

  let comments: Array<{ id: string; commenter: string; body: string; modified: string }> = [];
  let newComment = '';
  let loading = false;
  let error = '';
  let userNames: Record<string, string> = {};

  const formIdStore = derived(page, ($page) => {
    if ($page.params?.formId) return $page.params.formId;
    if ($page.data?.form?.id) return $page.data.form.id;
    return '';
  });
  let formId = '';
  $: formIdStore.subscribe((id) => formId = id);

  async function fetchUserName(userId: string): Promise<string> {
    if (!userId) return 'Unknown';
    if (userNames[userId]) return userNames[userId];
    try {
      const res = await fetch(`/api/users/${userId}`);
      if (!res.ok) throw new Error();
      const data = await res.json();
      userNames[userId] = data.name || data.handle || 'Unknown';
      return userNames[userId];
    } catch {
      userNames[userId] = 'Unknown';
      return 'Unknown';
    }
  }

  async function fetchAllUserNames(commentsList: typeof comments) {
    const uniqueIds = Array.from(new Set(commentsList.map(c => c.commenter).filter(Boolean)));
    await Promise.all(uniqueIds.map(fetchUserName));
  }

  async function fetchComments() {
    if (!formId) return;
    loading = true;
    error = '';
    try {
      const res = await fetch(`/api/forms/${formId}/comments`);
      if (!res.ok) throw new Error('Failed to fetch comments');
      const data = await res.json();
      comments = data.data || [];
      await fetchAllUserNames(comments);
    } catch (e: any) {
      error = e.message || 'Unknown error';
    } finally {
      loading = false;
    }
  }

  async function addComment() {
    if (!newComment.trim() || !formId) return;
    loading = true;
    error = '';
    try {
      const res = await fetch(`/api/forms/${formId}/comments`, {
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

  onMount(fetchComments);
  $: if (formId) fetchComments();
</script>

<div class="flex flex-col h-full overflow-hidden">
  {#if error}
    <div class="px-8 text-red-500 text-sm mb-2">{error}</div>
  {/if}

  <div class="flex-1 overflow-y-auto px-8 pt-4 pb-4 space-y-3">
    {#if loading}
      <div class="text-muted-foreground text-sm">Loading...</div>
    {:else if comments.length === 0}
      <div class="text-muted-foreground text-sm min-h-40 flex justify-center items-center">No comments yet.</div>
    {:else}
      {#each comments as c}
        <div class="group border rounded-md p-3 bg-muted/30 text-left relative">
          <div class="flex items-center justify-between mb-1">
            <span class="text-xs text-muted-foreground">{userNames[c.commenter] || '...'} &middot; {formatRelativeTime(c.modified)}</span>
            <button
              type="button"
              class="opacity-0 group-hover:opacity-100 transition-opacity ml-2 text-muted-foreground hover:text-destructive p-1 rounded"
              aria-label="Delete comment"
              on:click={(e) => { e.stopPropagation(); openDeleteDialog(c.id); }}
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
    on:submit|preventDefault={addComment}
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
      <Dialog.Description>This action cannot be undone. Are you sure you want to delete this comment?</Dialog.Description>
    </Dialog.Header>
    <Dialog.Footer>
      <Button
        type="button"
        variant="destructive"
        class="px-4 py-2 rounded-md"
        onclick={(e) => { e.stopPropagation(); confirmDeleteComment(); }}
        disabled={loading}
      >
        Delete
      </Button>
      <Button
        type="button"
        variant="ghost"
        class="px-4 py-2 rounded-md ml-2"
        onclick={(e) => { e.stopPropagation(); showDeleteDialog = false; commentToDelete = null; }}
      >
        Cancel
    </Button>
    </Dialog.Footer>
  </Dialog.Content>
</Dialog.Root>

