<script lang="ts">
  import { Button } from '$lib/components/ui/button';
  import { IconUpload } from '@tabler/icons-svelte';

  // TODO: this is just decorative atp

  interface Props {
    question: {
      uid: string;
      text: string;
      required: boolean;
    };
    value?: string;
  }

  let { question, value = $bindable('') }: Props = $props();
  let fileInput: HTMLInputElement;

  function handleFileSelect(event: Event) {
    const target = event.target as HTMLInputElement;
    const files = target.files;
    
    if (files?.length) {
      const fileName = files.length === 1 ? files[0].name : `${files.length} files`;
      value = `${fileName}`;
      target.value = '';
    } 
  }

  function triggerFileInput() {
    fileInput?.click();
  }
</script>

<div class="space-y-4">
  <input
    bind:this={fileInput}
    type="file"
    multiple
    class="hidden"
    onchange={handleFileSelect}
    accept="*/*"
  />

  <div class="border-2 border-dashed border-muted-foreground/25 rounded-lg p-6 text-center">
    <div class="space-y-4">
      <div class="w-12 h-12 bg-muted rounded-md flex items-center justify-center mx-auto">
        <IconUpload class="h-6 w-6 text-muted-foreground" />
      </div>
      <div>
        <h3 class="font-medium mb-1">Upload Files</h3>
        <p class="text-sm text-muted-foreground mb-4">
          Choose files or drag and drop them here.
          {#if question.required}
            <span class="text-red-500">*</span>
          {/if}
        </p>
        <Button type="button" variant="outline" onclick={triggerFileInput}>
          Choose Files
        </Button>
      </div>
    </div>
  </div>

  {#if value && value !== '[]'}
    <div class="p-3 bg-muted/50 rounded-md border">
      <p class="text-sm text-muted-foreground">{value}</p>
    </div>
  {/if}
</div>