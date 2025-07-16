<script lang="ts">
  import { Card, CardContent } from '$lib/components/ui/card';
  import { Button } from '$lib/components/ui/button';
  import { IconChevronUp, IconChevronDown, IconTrash } from '@tabler/icons-svelte';
  import InputEditor from '$lib/components/editor/questions/input.svelte';
  import TextareaEditor from '$lib/components/editor/questions/textarea.svelte';
  import RadioEditor from '$lib/components/editor/questions/radio.svelte';
  import CheckboxEditor from '$lib/components/editor/questions/checkbox.svelte';
  import FileUploadEditor from '$lib/components/editor/questions/file-upload.svelte';
  import SelectEditor from '$lib/components/editor/questions/select.svelte';
  import DatePickerEditor from '$lib/components/editor/questions/date-picker.svelte';
  import type { Question, QuestionType } from './FormEditor';

  export let question: Question;
  export let i: number;
  export let questionTypeLabels: Record<QuestionType, string>;
  export let moveQuestionUp: (index: number) => void;
  export let moveQuestionDown: (index: number) => void;
  export let removeQuestion: (id: string) => void;
  export let questionsLength: number;
</script>

<Card>
  <CardContent class="space-y-6">
    <!-- header -->
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-3">
        <div class="flex items-center justify-center w-8 h-8 bg-muted text-sm font-medium rounded-md">
          {i + 1}
        </div>
        <span class="font-medium">{questionTypeLabels[question.type]}</span>
      </div>
      <!-- controls -->
      <div class="flex items-center gap-1">
        <Button
          variant="ghost"
          size="icon"
          class="h-8 w-8"
          onclick={() => moveQuestionUp(i)}
          disabled={i === 0}
        >
          <IconChevronUp class="h-4 w-4" />
        </Button>
        <Button
          variant="ghost"
          size="icon"
          class="h-8 w-8"
          onclick={() => moveQuestionDown(i)}
          disabled={i === questionsLength - 1}
        >
          <IconChevronDown class="h-4 w-4" />
        </Button>
        <Button
          variant="ghost"
          size="icon"
          class="h-8 w-8 hover:bg-destructive/10 hover:text-destructive"
          onclick={() => removeQuestion(question.id)}
        >
          <IconTrash class="h-4 w-4" />
        </Button>
      </div>
    </div>
    <!-- editor container -->
    <div class="pl-4">
      {#if question.type === 'input'}
        <InputEditor bind:question />
      {:else if question.type === 'textarea'}
        <TextareaEditor bind:question />
      {:else if question.type === 'radio'}
        <RadioEditor bind:question />
      {:else if question.type === 'checkbox'}
        <CheckboxEditor bind:question />
      {:else if question.type === 'file'}
        <FileUploadEditor bind:question />
      {:else if question.type === 'select'}
        <SelectEditor bind:question />
      {:else if question.type === 'date'}
        <DatePickerEditor bind:question />
      {/if}
    </div>
  </CardContent>
</Card>
