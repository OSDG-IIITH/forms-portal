<script lang="ts">
  import Basics from './basics.svelte';
  import { getContext } from 'svelte';
  import type { FormStore } from '../form-store.svelte';

  let { questionId }: { questionId: string } = $props();

  const store: FormStore = getContext('form-store');
  const question = $derived(store.questions.find((q) => q.id === questionId));
</script>

{#if question}
<div class="space-y-4">
  <Basics
    questionId={question.id}
    title={question.title}
    error={question.error}
    required={question.required}
    requiredLabel="Required field"
    requiredClass="text-sm font-normal cursor-pointer"
    requiredContainerClass="flex items-center space-x-2"
    onTitleChange={(value) => store.updateQuestion(questionId, { title: value })}
    onRequiredChange={(value) => store.updateQuestion(questionId, { required: value })}
  />
</div>
{/if}