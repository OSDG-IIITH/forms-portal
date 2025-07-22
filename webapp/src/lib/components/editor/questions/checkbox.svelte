<script lang="ts">
  import { Input } from '$lib/components/ui/input';
  import { Label } from '$lib/components/ui/label';
  import { Checkbox } from '$lib/components/ui/checkbox';
  import { Button } from '$lib/components/ui/button';
  import { IconPlus, IconGripVertical, IconTrash } from '@tabler/icons-svelte';
  import { ulid } from 'ulid';
  import { getContext } from 'svelte';
  import type { FormStore, Option } from '../form-store.svelte';

  let { questionId }: { questionId: string } = $props();

  const store: FormStore = getContext('form-store');
  const question = $derived(store.questions.find((q) => q.id === questionId));

  let draggedIndex = $state<number | null>(null);
  let dragOverIndex = $state<number | null>(null);

  function addOption(): void {
    if (!question?.options) return;
    const newId = ulid();
    const newOptions = [...question.options, { id: newId, value: '', label: '' }];
    store.updateQuestion(questionId, { options: newOptions });
  }

  function removeOption(id: string): void {
    if (!question?.options) return;
    const newOptions = question.options.filter((option) => option.id !== id);
    store.updateQuestion(questionId, { options: newOptions });
  }

  function handleDragStart(event: DragEvent, index: number): void {
    if (event.dataTransfer) {
      event.dataTransfer.effectAllowed = 'move';
      draggedIndex = index;
    }
  }

  function handleDragOver(event: DragEvent, index: number): void {
    event.preventDefault();
    if (event.dataTransfer) {
      event.dataTransfer.dropEffect = 'move';
    }
    dragOverIndex = index;
  }

  function handleDragLeave(): void {
    dragOverIndex = null;
  }

  function handleDrop(event: DragEvent, dropIndex: number): void {
    event.preventDefault();

    if (draggedIndex !== null && draggedIndex !== dropIndex && question?.options) {
      const newOptions = [...question.options];
      const [draggedOption] = newOptions.splice(draggedIndex, 1);
      newOptions.splice(dropIndex, 0, draggedOption);
      store.updateQuestion(questionId, { options: newOptions });
    }

    draggedIndex = null;
    dragOverIndex = null;
  }

  function handleDragEnd(): void {
    draggedIndex = null;
    dragOverIndex = null;
  }

  function updateOptionValue(optionId: string, newLabel: string): void {
    if (!question?.options) return;
    const newOptions = question.options.map(opt => 
      opt.id === optionId ? { ...opt, label: newLabel, value: newLabel } : opt
    );
    store.updateQuestion(questionId, { options: newOptions });
  }
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

  <div class="space-y-3">
    <Label class="text-sm font-medium">Options</Label>
    <div class="space-y-2">
      {#if question.options}
        {#each question.options as option, i (option.id)}
          <div
            role="button"
            tabindex="0"
            class="flex items-center gap-3 p-3 border rounded-md transition-colors {dragOverIndex === i
              ? 'border-primary bg-primary/5'
              : ''} {draggedIndex === i ? 'opacity-50' : ''}"
            draggable="true"
            ondragstart={(e) => handleDragStart(e, i)}
            ondragover={(e) => handleDragOver(e, i)}
            ondragleave={handleDragLeave}
            ondrop={(e) => handleDrop(e, i)}
            ondragend={handleDragEnd}
          >
            <IconGripVertical class="cursor-grab" />
            <Input 
              placeholder="Option" 
              value={option.label}
              oninput={(e) => updateOptionValue(option.id, e.currentTarget.value)}
            />
            <Button
              variant="ghost"
              size="icon"
              class="h-8 w-8 hover:bg-destructive/10 hover:text-destructive"
              onclick={() => removeOption(option.id)}
            >
              <IconTrash class="h-4 w-4" />
            </Button>
          </div>
        {/each}
      {/if}
    </div>
  </div>

  <Button variant="outline" class="w-full" onclick={addOption}>
    <IconPlus class="mr-2" />
    Add Option
  </Button>

  <div class="flex items-center gap-2">
    <Checkbox 
      id="required-{question.id}" 
      checked={question.required}
      onchange={() => store.updateQuestion(questionId, { required: !question.required })}
    />
    <Label for="required-{question.id}" class="text-sm">
      Required
    </Label>
  </div>
</div>
{/if}
