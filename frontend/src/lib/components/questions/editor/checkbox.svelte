<script lang="ts">
  import { Input } from '$lib/components/ui/input';
  import { Label } from '$lib/components/ui/label';
  import { Checkbox } from '$lib/components/ui/checkbox';
  import { Button } from '$lib/components/ui/button';
  import { IconPlus, IconGripVertical, IconTrash } from '@tabler/icons-svelte';

  interface Option {
    uid: string;
    text: string;
    type: 'text';
  }

  interface Props {
    question: {
      uid: string;
      text: string;
      required: boolean;
      options: Option[];
    };
  }

  let { question = $bindable() }: Props = $props();
  let draggedIndex = $state<number | null>(null);
  let dragOverIndex = $state<number | null>(null);

  function addOption(): void {
    question.options = [
      ...question.options,
      { uid: crypto.randomUUID(), text: '', type: 'text' }
    ];
  }

  function removeOption(uid: string): void {
    question.options = question.options.filter(option => option.uid !== uid);
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
    
    if (draggedIndex !== null && draggedIndex !== dropIndex) {
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

  if (question.options.length === 0) {
    addOption();
    addOption();
  }
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

  <div class="space-y-3">
    <Label class="text-sm font-medium">Options (select multiple)</Label>
    <div class="space-y-2">
      {#each question.options as option, i (option.uid)}
        <div 
          role="button"
          tabindex="0"
          class="flex items-center gap-3 p-3 border rounded-md transition-colors {dragOverIndex === i ? 'border-primary bg-primary/5' : ''} {draggedIndex === i ? 'opacity-50' : ''}"
          draggable="true"
          ondragstart={(e) => handleDragStart(e, i)}
          ondragover={(e) => handleDragOver(e, i)}
          ondragleave={handleDragLeave}
          ondrop={(e) => handleDrop(e, i)}
          ondragend={handleDragEnd}
        >
          <IconGripVertical class="h-4 w-4 text-muted-foreground cursor-grab active:cursor-grabbing" />
          <div class="w-4 h-4 border-2 border-muted-foreground/40 rounded-sm"></div>
          <Input
            placeholder="Enter option text"
            bind:value={option.text}
            class="flex-1 border-0 bg-transparent focus-visible:ring-0 focus-visible:ring-offset-0"
          />
          {#if question.options.length > 1}
            <Button
              variant="ghost"
              size="sm"
              class="h-8 w-8 p-0 text-muted-foreground hover:text-destructive hover:bg-destructive/10 transition-colors"
              onclick={() => removeOption(option.uid)}
              aria-label="Remove option"
            >
              <IconTrash class="h-4 w-4" />
            </Button>
          {/if}
        </div>
      {/each}
    </div>
    
    <Button 
      variant="outline" 
      size="sm" 
      class="w-full"
      onclick={addOption}
    >
      <IconPlus class="h-4 w-4 mr-2" />
      Add Option
    </Button>
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
