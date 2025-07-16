<script lang="ts">
  import { FileDropZone, displaySize, MEGABYTE, type FileDropZoneProps } from '$lib/components/ui/file-drop-zone';
  import { sleep } from '$lib/utils/sleep';
  import { onDestroy } from 'svelte';
  import { toast } from 'svelte-sonner';
  import { Button } from '$lib/components/ui/button';
  import { IconTrash, IconFileText } from '@tabler/icons-svelte';

  interface UploadedFile {
    name: string;
    type: string;
    size: number;
    uploadedAt: number;
    url: Promise<string>;
  }


  export let question: {
    id: string;
    title: string;
    required: boolean;
    'max-file-size'?: number;
    'max-files'?: number;
    'allowed-types'?: string[];
  };

  export let value: string = '';

  const maxFileSize = (question['max-file-size'] || 5) * MEGABYTE;
  const maxFiles = question['max-files'] === -1 ? undefined : (question['max-files'] || 3);
  const allowedTypes = question['allowed-types'] || [];


  let files: UploadedFile[] = [];
  let interval: any;
  let now = Date.now();

  const updateNow = () => { now = Date.now(); };

  const onUpload: FileDropZoneProps['onUpload'] = async (selectedFiles) => {
    await Promise.allSettled(selectedFiles.map((file) => uploadFile(file)));
  };

  const onFileRejected: FileDropZoneProps['onFileRejected'] = ({ reason, file }) => {
    toast.error(`${file.name} failed to upload!`, { description: reason });
  };

  const uploadFile = async (file: File) => {
    if (files.find((f) => f.name === file.name)) return;
    
    if (allowedTypes.length > 0 && !allowedTypes.includes(file.type)) {
      toast.error(`${file.name} file type not allowed!`, { 
        description: `Allowed types: ${allowedTypes.join(', ')}` 
      });
      return;
    }
    
    const urlPromise = new Promise<string>((resolve) => {
      sleep(1000).then(() => resolve(URL.createObjectURL(file)));
    });
    files = [
      ...files,
      {
        name: file.name,
        type: file.type,
        size: file.size,
        uploadedAt: Date.now(),
        url: urlPromise
      }
    ];
    await urlPromise;
  };

  onDestroy(async () => {
    for (const file of files) {
      URL.revokeObjectURL(await file.url);
    }
    clearInterval(interval);
  });

  interval = setInterval(updateNow, 100);
</script>

<div class="flex w-full flex-col gap-4 p-0 pt-2 md:p-6 md:pt-2 md:pb-0">
  <FileDropZone
    {onUpload}
    {onFileRejected}
    maxFileSize={maxFileSize}
    maxFiles={maxFiles}
    fileCount={maxFiles !== undefined ? files.length : undefined}
    accept={allowedTypes.length > 0 ? allowedTypes.join(',') : undefined}
  />
  <div class="flex flex-col gap-2">
    {#each files as file, i (file.name)}
      <div class="flex place-items-center justify-between gap-2">
        <div class="flex place-items-center gap-2">
          {#await file.url then src}
            {#if file.type.startsWith('image/')}
              <div class="relative size-9 overflow-clip">
                <img
                  {src}
                  alt={file.name}
                  class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 overflow-clip"
                />
              </div>
            {:else}
              <div class="relative size-9 flex items-center justify-center bg-muted rounded">
                <IconFileText class="h-6 w-6 text-muted-foreground" />
              </div>
            {/if}
          {/await}
          <div class="flex flex-col">
            <span>{file.name}</span>
            <span class="text-muted-foreground text-xs">{displaySize(file.size)}</span>
          </div>
        </div>
        {#await file.url}
          <div class="h-2 w-full grow bg-muted-foreground/10 rounded">
            <div class="h-2 bg-primary rounded" style="width: {Math.min(((now - file.uploadedAt) / 1000) * 100, 100)}%"></div>
          </div>
        {:then url}
          <Button
            variant="outline"
            size="icon"
            class="h-8 w-8 hover:bg-destructive/10 hover:text-destructive"
            type="button"
            onclick={() => {
              URL.revokeObjectURL(url);
              files = [...files.slice(0, i), ...files.slice(i + 1)];
            }}
          >
            <IconTrash class="h-4 w-4" />
          </Button>
        {/await}
      </div>
    {/each}
  </div>
</div>
