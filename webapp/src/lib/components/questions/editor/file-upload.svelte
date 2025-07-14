<script lang="ts">
  import { Input } from '$lib/components/ui/input';
  import { Label } from '$lib/components/ui/label';
  import { Checkbox } from '$lib/components/ui/checkbox';
  import * as Select from '$lib/components/ui/select/index.js';

  interface Props {
    question: {
      id: string;
      title: string;
      required: boolean;
      'max-file-size'?: number;
      'max-files'?: number;
      'allowed-types'?: string[];
    };
  }

  let { question = $bindable() }: Props = $props();

  if (!question['max-file-size']) question['max-file-size'] = 10;
  if (!question['max-files']) question['max-files'] = -1;
  if (!question['allowed-types']) question['allowed-types'] = [];

  let fileSizeValue = $state(String(question['max-file-size'] || 10));
  let fileSizeUnit = $state('mb');
  let maxFilesValue = $state(question['max-files'] === -1 ? '' : String(question['max-files'] || ''));
  let allowedFileTypes = $state(
    question['allowed-types']?.length 
      ? 'custom' 
      : 'any'
  );
  let customFileTypes = $state(
    question['allowed-types']?.join(', ') || ''
  );

  $effect(() => {
    const sizeNum = parseFloat(fileSizeValue) || 1;
    if (fileSizeUnit === 'kb') {
      question['max-file-size'] = Math.round((sizeNum / 1024) * 100) / 100;
    } else {
      question['max-file-size'] = sizeNum;
    }
    
    const maxFilesNum = parseInt(maxFilesValue);
    question['max-files'] = maxFilesValue.trim() === '' ? -1 : (isNaN(maxFilesNum) ? -1 : maxFilesNum);
    
    if (allowedFileTypes === 'any') {
      question['allowed-types'] = [];
    } else if (allowedFileTypes === 'images') {
      question['allowed-types'] = ['image/jpeg', 'image/png', 'image/gif'];
    } else if (allowedFileTypes === 'documents') {
      question['allowed-types'] = ['application/pdf', 'application/msword', 'application/vnd.openxmlformats-officedocument.wordprocessingml.document'];
    } else if (allowedFileTypes === 'spreadsheets') {
      question['allowed-types'] = ['application/vnd.ms-excel', 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet', 'text/csv'];
    } else if (allowedFileTypes === 'custom' && customFileTypes.trim()) {
      question['allowed-types'] = customFileTypes.split(',').map(t => t.trim()).filter(t => t);
    }
  });

  const fileSizeUnits = [
    { value: 'kb', label: 'KB' },
    { value: 'mb', label: 'MB' }
  ];

  const fileTypeOptions = [
    { value: 'any', label: 'Any file type' },
    { value: 'images', label: 'Images only (JPG, PNG, GIF)' },
    { value: 'documents', label: 'Documents (PDF, DOC, DOCX)' },
    { value: 'spreadsheets', label: 'Spreadsheets (XLS, XLSX, CSV)' },
    { value: 'custom', label: 'Custom file types' }
  ];
</script>

<div class="space-y-4">
  <div class="space-y-2">
    <Label for="question-{question.id}">Question Text</Label>
    <Input
      id="question-{question.id}"
      placeholder="Enter your question here"
      bind:value={question.title}
    />
  </div>

  <div class="space-y-4 p-4 border rounded-md bg-muted/20">
    <Label class="text-sm font-medium">File Upload Settings</Label>
    
    <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
      <div class="space-y-2">
        <Label class="text-xs font-medium text-muted-foreground">Max File Size</Label>
        <div class="flex gap-2">
          <Input
            type="text"
            placeholder="10"
            bind:value={fileSizeValue}
            class="flex-1 h-9"
          />
          <Select.Root type="single" bind:value={fileSizeUnit}>
            <Select.Trigger class="w-20 h-9" size="sm">
              {fileSizeUnits.find(opt => opt.value === fileSizeUnit)?.label ?? "MB"}
            </Select.Trigger>
            <Select.Content>
              {#each fileSizeUnits as unit}
                <Select.Item value={unit.value}>{unit.label}</Select.Item>
              {/each}
            </Select.Content>
          </Select.Root>
        </div>
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
        <Input
          type="text"
          placeholder="Unlimited"
          bind:value={maxFilesValue}
          class="w-full"
        />
        <p class="text-xs text-muted-foreground">
          Leave blank for unlimited files
        </p>
      </div>
    </div>

    {#if allowedFileTypes === 'custom'}
      <div class="space-y-2">
        <Label class="text-xs font-medium text-muted-foreground">Custom File Types</Label>
        <Input
          placeholder="e.g., image/png, application/pdf, text/csv"
          bind:value={customFileTypes}
          class="text-sm"
        />
        <p class="text-xs text-muted-foreground">
          Enter MIME types separated by commas
        </p>
      </div>
    {/if}
  </div>

  <div class="flex items-center space-x-2">
    <Checkbox id="required-{question.id}" bind:checked={question.required} />
    <Label for="required-{question.id}" class="text-sm font-normal cursor-pointer">
      Required field
    </Label>
  </div>
</div>
