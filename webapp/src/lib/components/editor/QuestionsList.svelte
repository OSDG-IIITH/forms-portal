<script lang="ts">
  import QuestionItem from './QuestionItem.svelte';
  import { getContext } from 'svelte';
  import type { FormStore, QuestionType } from './form-store.svelte';

  export let questionTypeLabels: Record<QuestionType, string>;

  const store: FormStore = getContext('form-store');
</script>

<div class="space-y-4">
  {#each store.questions as question, i (question.id)}
    <QuestionItem
      {question}
      {i}
      {questionTypeLabels}
      questionsLength={store.questions.length}
    />
  {/each}
</div>
