<script lang="ts">
  import { Input } from '$lib/components/ui/input';
  import { Label } from '$lib/components/ui/label';
  import { Checkbox } from '$lib/components/ui/checkbox';
  import * as Select from '$lib/components/ui/select/index.js';

  interface Props {
    question: {
      uid: string;
      text: string;
      required: boolean;
      options: any[];
    };
  }

  let { question = $bindable() }: Props = $props();
  
  let maxFileSize = $state('10');
  let allowedFileTypes = $state('any');
  let maxFiles = $state('1');

  const fileSizeOptions = [
    { value: '1', label: '1 MB' },
    { value: '5', label: '5 MB' },
    { value: '10', label: '10 MB' },
    { value: '25', label: '25 MB' },
    { value: '50', label: '50 MB' },
    { value: '100', label: '100 MB' }
  ];

  const fileTypeOptions = [
    { value: 'any', label: 'Any file type' },
    { value: 'images', label: 'Images only (JPG, PNG, GIF)' },
    { value: 'documents', label: 'Documents (PDF, DOC, DOCX)' },
    { value: 'spreadsheets', label: 'Spreadsheets (XLS, XLSX, CSV)' },
    { value: 'custom', label: 'Custom file types' }
  ];

  const maxFilesOptions = [
    { value: '1', label: '1 file' },
    { value: '3', label: '3 files' },
    { value: '5', label: '5 files' },
    { value: '10', label: '10 files' },
    { value: 'unlimited', label: 'Unlimited' }
  ];
</script>

<div class="space-y-4">
  <div class="space-y-2">
    <Label for="question-{question.uid}">Question Text</Label>
    <Input
      id="question-{question.uid}"
      placeholder="Enter your question here"
      bind:value={question.text}
    />
  </div>

  <div class="space-y-4 p-4 border rounded-md bg-muted/20">
    <Label class="text-sm font-medium">File Upload Settings</Label>
    
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
      <div class="space-y-2">
        <Label class="text-xs font-medium text-muted-foreground">Max File Size</Label>
        <Select.Root type="single" bind:value={maxFileSize}>
          <Select.Trigger class="w-full" size="sm">
            {fileSizeOptions.find(opt => opt.value === maxFileSize)?.label ?? "Select"}
          </Select.Trigger>
          <Select.Content>
            {#each fileSizeOptions as option}
              <Select.Item value={option.value}>{option.label}</Select.Item>
            {/each}
          </Select.Content>
        </Select.Root>
      </div>

      <div class="space-y-2">
        <Label class="text-xs font-medium text-muted-foreground">File Types</Label>
        <Select.Root type="single" bind:value={allowedFileTypes}>
          <Select.Trigger class="w-full" size="sm">
            {fileTypeOptions.find(opt => opt.value === allowedFileTypes)?.label ?? "Select"}
          </Select.Trigger>
          <Select.Content>
            {#each fileTypeOptions as option}
              <Select.Item value={option.value}>{option.label}</Select.Item>
            {/each}
          </Select.Content>
        </Select.Root>
      </div>

      <div class="space-y-2">
        <Label class="text-xs font-medium text-muted-foreground">Max Files</Label>
        <Select.Root type="single" bind:value={maxFiles}>
          <Select.Trigger class="w-full" size="sm">
            {maxFilesOptions.find(opt => opt.value === maxFiles)?.label ?? "Select"}
          </Select.Trigger>
          <Select.Content>
            {#each maxFilesOptions as option}
              <Select.Item value={option.value}>{option.label}</Select.Item>
            {/each}
          </Select.Content>
        </Select.Root>
      </div>
    </div>

    {#if allowedFileTypes === 'custom'}
      <div class="space-y-2">
        <Label class="text-xs font-medium text-muted-foreground">Custom File Extensions</Label>
        <Input
          placeholder="e.g., .pdf, .jpg, .png, .docx"
          class="h-9 text-sm"
        />
      </div>
    {/if}
  </div>

  <div class="flex items-center space-x-2">
    <Checkbox 
      id="required-{question.uid}" 
      bind:checked={question.required}
    />
    <Label 
      for="required-{question.uid}" 
      class="text-sm font-normal cursor-pointer"
    >
      Required field
    </Label>
  </div>
</div>
