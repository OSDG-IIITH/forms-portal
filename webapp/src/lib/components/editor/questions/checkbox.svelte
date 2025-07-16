<script lang="ts">
  import { Input } from '$lib/components/ui/input';
  import { Label } from '$lib/components/ui/label';
  import { Checkbox } from '$lib/components/ui/checkbox';
  import { Button } from '$lib/components/ui/button';
  import { IconPlus, IconGripVertical, IconTrash } from '@tabler/icons-svelte';
  import { ulid } from 'ulid';

  interface Option {
    id: string;
    value: string;
    label: string;
  }

  interface Props {
    question: {
      id: string;
      title: string;
      required: boolean;
      options?: Option[];
    };
  }

  let { question = $bindable() }: Props = $props();
  let draggedIndex = $state<number | null>(null);
  let dragOverIndex = $state<number | null>(null);

  function addOption(): void {
    if (!question.options) question.options = [];
    const newId = ulid();
    question.options = [...question.options, { id: newId, value: '', label: '' }];
  }

  function removeOption(id: string): void {
    if (!question.options) return;
    question.options = question.options.filter((option) => option.id !== id);
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

    if (draggedIndex !== null && draggedIndex !== dropIndex && question.options) {
      const newOptions = [...question.options];
      const [draggedOption] = newOptions.splice(draggedIndex, 1);
      newOptions.splice(dropIndex, 0, draggedOption);
      question.options = newOptions;
    }

    draggedIndex = null;
    dragOverIndex = null;
  }

  function handleDragEnd(): void {
    draggedIndex = null;
    dragOverIndex = null;
  }

  function updateOptionValue(option: Option, newLabel: string): void {
    option.label = newLabel;
    option.value = newLabel;
  }

  if (!question.options || question.options.length === 0) {
    addOption();
    addOption();
  }
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
              oninput={(e: Event) => {
                const target = e.target as HTMLInputElement;
                updateOptionValue(option, target.value);
              }}
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
    <Checkbox id="required-{question.id}" bind:checked={question.required} />
    <Label for="required-{question.id}" class="text-sm">
      Required
    </Label>
  </div>
</div>
