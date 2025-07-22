<script lang="ts">
  import { Input } from '$lib/components/ui/input';
  import { Label } from '$lib/components/ui/label';
  import { Checkbox } from '$lib/components/ui/checkbox';
  import * as Select from '$lib/components/ui/select';
  import { getContext } from 'svelte';
  import type { FormStore } from '../form-store.svelte';

  let { questionId }: { questionId: string } = $props();

  const store: FormStore = getContext('form-store');
  const question = $derived(store.questions.find((q) => q.id === questionId));

  const fileSize = $derived.by(() => {
    const sizeInMb = question?.['max-file-size'] ?? 10;
    if (sizeInMb < 1) {
      return { value: Math.round(sizeInMb * 1024), unit: 'kb' };
    }
    return { value: sizeInMb, unit: 'mb' };
  });

  const allowedFileTypes = $derived.by(() => {
    const types = question?.['allowed-types'] ?? [];
    if (types.length === 0) return 'any';
    const typesString = JSON.stringify([...types].sort());
    if (typesString === JSON.stringify(['image/jpeg', 'image/png', 'image/gif'].sort())) return 'images';
    if (typesString === JSON.stringify(['application/pdf', 'application/msword', 'application/vnd.openxmlformats-officedocument.wordprocessingml.document'].sort())) return 'documents';
    if (typesString === JSON.stringify(['application/vnd.ms-excel', 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet', 'text/csv'].sort())) return 'spreadsheets';
    return 'custom';
  });

  function handleFileSizeValueChange(e: Event) {
    const target = e.currentTarget as HTMLInputElement;
    const value = parseFloat(target.value);
    const sizeNum = isNaN(value) ? 10 : value;
    const unit = fileSize.unit;
    const newSizeInMb = unit === 'kb' ? sizeNum / 1024 : sizeNum;
    store.updateQuestion(questionId, { 'max-file-size': Math.round(newSizeInMb * 100) / 100 });
  }

  function handleFileSizeUnitChange(unit: string | undefined) {
    if (!unit || !question) return;
    const value = fileSize.value;
    const newSizeInMb = unit === 'kb' ? value / 1024 : value;
    store.updateQuestion(questionId, { 'max-file-size': newSizeInMb });
  }

  function handleMaxFilesChange(e: Event) {
    const value = (e.currentTarget as HTMLInputElement).value.trim();
    const maxFiles = value === '' ? -1 : parseInt(value, 10);
    store.updateQuestion(questionId, { 'max-files': isNaN(maxFiles) ? -1 : maxFiles });
  }
  
  function handleAllowedTypesChange(type: string | undefined) {
    if (!type) return;
    let newTypes: string[] = [];
    if (type === 'images') {
      newTypes = ['image/jpeg', 'image/png', 'image/gif'];
    } else if (type === 'documents') {
      newTypes = ['application/pdf', 'application/msword', 'application/vnd.openxmlformats-officedocument.wordprocessingml.document'];
    } else if (type === 'spreadsheets') {
      newTypes = ['application/vnd.ms-excel', 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet', 'text/csv'];
    }
    store.updateQuestion(questionId, { 'allowed-types': newTypes });
  }
  
  function handleCustomTypesChange(e: Event) {
    const value = (e.currentTarget as HTMLInputElement).value;
    const newTypes = value.split(',').map(t => t.trim()).filter(t => t);
    store.updateQuestion(questionId, { 'allowed-types': newTypes });
  }

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

{#if question}
<div class="space-y-4">
  <div class="space-y-2">
    <Label for="question-{question.id}">Question Text</Label>
    <Input
      id="question-{question.id}"
      placeholder="Enter your question here"
      value={question.title}
      oninput={(e) => store.updateQuestion(questionId, { title: e.currentTarget.value })}
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
            value={fileSize.value}
            oninput={handleFileSizeValueChange}
            class="flex-1 h-9"
          />
          <Select.Root
            type="single"
            value={fileSize.unit}
            onValueChange={handleFileSizeUnitChange}
          >
            <Select.Trigger class="w-20 h-9" size="sm">
              {fileSize.unit.toUpperCase()}
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
        <Select.Root
          type="single"
          value={allowedFileTypes}
          onValueChange={handleAllowedTypesChange}
        >
          <Select.Trigger class="w-full h-9" size="sm">
             {fileTypeOptions.find(opt => opt.value === allowedFileTypes)?.label ?? "Select a type"}
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
          value={question['max-files'] === -1 ? '' : (question['max-files'] ?? '')}
          oninput={handleMaxFilesChange}
          class="w-full h-9"
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
          value={question['allowed-types']?.join(', ') || ''}
          oninput={handleCustomTypesChange}
          class="text-sm"
        />
        <p class="text-xs text-muted-foreground">
          Enter MIME types separated by commas
        </p>
      </div>
    {/if}
  </div>

  <div class="flex items-center space-x-2">
    <Checkbox
      id="required-{question.id}"
      checked={question.required}
      onchange={() => store.updateQuestion(questionId, { required: !question.required })}
    />
    <Label for="required-{question.id}" class="text-sm font-normal cursor-pointer">
      Required field
    </Label>
  </div>
</div>
{/if}