<script lang="ts">
  import type { Row } from "@tanstack/table-core";
  import { Button } from "$lib/components/ui/button/index.js";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogFooter,
  DialogTitle,
  DialogDescription,
  DialogClose
} from "$lib/components/ui/dialog/index.js";
  import { createEventDispatcher } from 'svelte';
  import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
  } from "$lib/components/ui/dropdown-menu/index.js";
  import { MoreHorizontal } from "@lucide/svelte";
  import type { Group } from "./columns.js";
  import { page } from "$app/stores";
  import { goto } from "$app/navigation";

  interface $$Props {
    row: Row<Group>;
  }


  let { row }: $$Props = $props();
let showDeleteDialog = $state(false);
let isDeleting = $state(false);
  const dispatch = createEventDispatcher();

  const viewGroup = () => {
    const currentPath = $page.url.pathname;
    const groupPath = `${currentPath}/${row.original.id}`;
    goto(groupPath);
  };

async function deleteGroup() {
  isDeleting = true;
  try {
    const groupId = row.original.id;
    const res = await fetch(`/api/groups/${groupId}`, {
      method: 'DELETE',
      credentials: 'same-origin'
    });
    if (res.status === 204) {
      dispatch('deleted', { id: groupId });
      showDeleteDialog = false;
    } else {
      let msg = 'Failed to delete group. (ERR' + res.status + ')';
      try {
        const data = await res.json();
        if (data && data.message) msg = data.message;
      } catch {}
      alert(msg);
    }
  } catch (e) {
    alert('Failed to delete group.');
  } finally {
    isDeleting = false;
  }
}
</script>


<DropdownMenu>
  <DropdownMenuTrigger>
    <Button
      variant="ghost"
      class="flex h-8 w-8 p-0 data-[state=open]:bg-muted"
    >
      <MoreHorizontal class="h-4 w-4" />
      <span class="sr-only">Open menu</span>
    </Button>
  </DropdownMenuTrigger>
  <DropdownMenuContent align="end" class="w-[160px]">
    <DropdownMenuItem onclick={viewGroup}>
      View Group
    </DropdownMenuItem>
    <DropdownMenuSeparator />
    <DropdownMenuItem class="text-destructive focus:text-destructive" onclick={() => showDeleteDialog = true}>
      Delete
    </DropdownMenuItem>
  </DropdownMenuContent>
</DropdownMenu>

{#if showDeleteDialog}
  <Dialog open={showDeleteDialog} onOpenChange={v => showDeleteDialog = v}>
    <DialogContent>
      <DialogHeader>
        <DialogTitle>Delete group?</DialogTitle>
        <DialogDescription>
          This action cannot be undone. This will permanently delete the group <b>{row.original.name}</b> and remove all its data.
        </DialogDescription>
      </DialogHeader>
      <DialogFooter>
        <DialogClose >
          <Button variant="outline" disabled={isDeleting} onclick={() => showDeleteDialog = false}>Cancel</Button>
        </DialogClose>
        <Button variant="destructive" onclick={deleteGroup} disabled={isDeleting}>
          {isDeleting ? 'Deleting...' : 'Delete'}
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
{/if}
