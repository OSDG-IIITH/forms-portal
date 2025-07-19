<script lang="ts">
import { Dialog, DialogContent, DialogTitle, DialogTrigger, DialogClose, DialogDescription } from '$lib/components/ui/dialog/index.js';
import { Button } from '$lib/components/ui/button/index.js';
import { Input } from '$lib/components/ui/input/index.js';
import * as Select from '$lib/components/ui/select/index.js';
import { IconUserPlus } from '@tabler/icons-svelte';
import { createEventDispatcher } from 'svelte';

const dispatch = createEventDispatcher();
let open = false;
let mode: 'user' | 'group' = 'user';
let email = '';
let groupId = '';
let isLoading = false;
let error = '';
let availableGroups: { id: string; name: string }[] = [];

function reset() {
  email = '';
  groupId = '';
  error = '';
  mode = 'user';
}

async function handleSubmit() {
  error = '';
  isLoading = true;
  if (mode === 'user') {
    if (!email || !/^[^@\s]+@[^@\s]+\.[^@\s]+$/.test(email)) {
      error = 'Enter a valid email.';
      isLoading = false;
      return;
    }
    dispatch('add', { type: 'user', email });
    reset();
    open = false;
    isLoading = false;
    return;
  }
  if (mode === 'group') {
    if (!groupId) {
      error = 'Select a group.';
      isLoading = false;
      return;
    }
    dispatch('add', { type: 'group', groupId });
    reset();
    open = false;
    isLoading = false;
    return;
  }
  isLoading = false;
}
</script>

<Dialog bind:open>
  <DialogTrigger>
    <Button class="ml-auto gap-2" variant="outline" size="sm" aria-label="Add user or group">
      <IconUserPlus class="w-4 h-4" />
      Add
    </Button>
  </DialogTrigger>
  <DialogContent class="max-w-md w-full p-2 pb-4 overflow-hidden rounded-xl">
    <form class="px-6 pt-6 pb-2 flex flex-col gap-4" on:submit|preventDefault={handleSubmit}>
      <DialogTitle class="text-lg font-semibold">Add User or Group</DialogTitle>
      <DialogDescription class="text-muted-foreground text-sm">Invite a user by email or add a group to this form. You can assign permissions after adding.</DialogDescription>
      <div class="flex gap-2">
        <Button type="button" variant={mode === 'user' ? 'outline' : 'ghost'} size="sm" class="flex-1" onclick={() => mode = 'user'}>User</Button>
        <Button type="button" variant={mode === 'group' ? 'outline' : 'ghost'} size="sm" class="flex-1" onclick={() => mode = 'group'}>Group</Button>
      </div>
      {#if mode === 'user'}
        <Input type="email" placeholder="User email" bind:value={email} required autocomplete="off" class="w-full" />
      {:else}
        <Select.Root
          allowDeselect={false}
          type="single"
          bind:value={groupId}
        >
          <Select.Trigger class="w-full">
            {#if groupId}
              {availableGroups.find(g => g.id === groupId)?.name ?? 'Select group'}
            {:else}
              <span class="text-muted-foreground">Select group</span>
            {/if}
          </Select.Trigger>
          <Select.Content>
            {#each availableGroups as group}
              <Select.Item value={group.id}>{group.name}</Select.Item>
            {/each}
          </Select.Content>
        </Select.Root>
      {/if}
      {#if error}
        <div class="text-destructive text-sm mt-1">{error}</div>
      {/if}
      <Button type="submit" class="mt-2 w-full" disabled={isLoading}>
        {isLoading ? 'Adding...' : 'Add'}
      </Button>
    </form>
  </DialogContent>
</Dialog>
