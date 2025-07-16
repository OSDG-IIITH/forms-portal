<script lang="ts">
  import QuestionItem from './QuestionItem.svelte';
  import type { Question, QuestionType } from './FormEditor';

  export let questions: Question[];
  export let questionTypeLabels: Record<QuestionType, string>;
  export let moveQuestionUp: (index: number) => void;
  export let moveQuestionDown: (index: number) => void;
  export let removeQuestion: (id: string) => void;

</script>

<div class="space-y-4">
  {#each questions as question, i (question.id)}
    <QuestionItem 
      {question} 
      {i} 
      {questionTypeLabels}
      {moveQuestionUp}
      {moveQuestionDown}
      {removeQuestion}
      questionsLength={questions.length}
    />
  {/each}
</div>
