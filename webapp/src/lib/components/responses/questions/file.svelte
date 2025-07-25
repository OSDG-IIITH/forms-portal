<script lang="ts">
  import { onMount } from 'svelte';
  import { IconFileText } from '@tabler/icons-svelte';
  import { displaySize } from '$lib/components/ui/file-drop-zone';

  export let value: string = '[]';

  interface UploadedFile {
    name: string;
    type: string;
    size: number;
    url: string;
  }

  let files: UploadedFile[] = [];

  onMount(() => {
    try {
      files = JSON.parse(value);
    } catch {
      files = [];
    }
  });
</script>

<div class="flex flex-col gap-2">
  {#each files as file (file.name)}
    <div class="flex place-items-center justify-between gap-2">
      <div class="flex place-items-center gap-2">
        {#if file.type.startsWith('image/')}
          <div class="relative size-9 overflow-clip">
            <img
              src={file.url}
              alt={file.name}
              class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 overflow-clip"
            />
          </div>
        {:else}
          <div class="relative size-9 flex items-center justify-center bg-muted rounded">
            <IconFileText class="h-6 w-6 text-muted-foreground" />
          </div>
        {/if}
        <div class="flex flex-col">
          <span>{file.name}</span>
          <span class="text-muted-foreground text-xs">{displaySize(file.size)}</span>
        </div>
      </div>
    </div>
  {/each}
</div>