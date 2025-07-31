<script lang="ts">
  import { Card, CardContent } from '$lib/components/ui/card';
  import { Button } from '$lib/components/ui/button';
  import { IconChevronUp, IconChevronDown, IconTrash } from '@tabler/icons-svelte';
  import InputEditor from '$lib/components/editor/questions/input.svelte';
  import TextareaEditor from '$lib/components/editor/questions/textarea.svelte';
  import RadioEditor from '$lib/components/editor/questions/radio.svelte';
  import CheckboxEditor from '$lib/components/editor/questions/checkbox.svelte';
  import FileUploadEditor from '$lib/components/editor/questions/file.svelte';
  import SelectEditor from '$lib/components/editor/questions/select.svelte';
  import DatePickerEditor from '$lib/components/editor/questions/date.svelte';
  import SectionHeaderEditor from '$lib/components/editor/questions/section-header.svelte';
  import { getContext } from 'svelte';
  import type { FormStore, Question, QuestionType } from './form-store.svelte';

  export let question: Question;
  export let i: number;
  export let questionTypeLabels: Record<QuestionType, string>;
  export let questionsLength: number;

  const store: FormStore = getContext('form-store');
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
          onclick={() => store.moveQuestionUp(i)}
          disabled={i === 0}
        >
          <IconChevronUp class="h-4 w-4" />
        </Button>
        <Button
          variant="ghost"
          size="icon"
          class="h-8 w-8"
          onclick={() => store.moveQuestionDown(i)}
          disabled={i === questionsLength - 1}
        >
          <IconChevronDown class="h-4 w-4" />
        </Button>
        <Button
          variant="ghost"
          size="icon"
          class="h-8 w-8 hover:bg-destructive/10 hover:text-destructive"
          onclick={() => store.removeQuestion(question.id)}
        >
          <IconTrash class="h-4 w-4" />
        </Button>
      </div>
    </div>
    <!-- editor container -->
    <div class="pl-4">
      {#if question.type === 'input'}
        <InputEditor questionId={question.id} />
      {:else if question.type === 'textarea'}
        <TextareaEditor questionId={question.id} />
      {:else if question.type === 'checkbox'}
        <CheckboxEditor questionId={question.id} />
      {:else if question.type === 'radio'}
        <RadioEditor questionId={question.id} />
      {:else if question.type === 'select'}
        <SelectEditor questionId={question.id} />
      {:else if question.type === 'date'}
        <DatePickerEditor questionId={question.id} />
      {:else if question.type === 'file'}
        <FileUploadEditor questionId={question.id} />
      {:else if question.type === 'section-header'}
        <SectionHeaderEditor questionId={question.id} />
      {/if}
    </div>
  </CardContent>
</Card>
