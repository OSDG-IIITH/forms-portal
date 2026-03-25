<script lang="ts">
  import { Label } from '$lib/components/ui/label';
  import { Textarea } from '$lib/components/ui/textarea';
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
    required={false}
    titleLabel="Section Title"
    titlePlaceholder="Enter section title"
    showRequired={false}
    onTitleChange={(value) => store.updateQuestion(questionId, { title: value })}
  />

  <div class="space-y-2">
    <Label for="section-description-{question.id}">Section Description</Label>
    <Textarea
      id="section-description-{question.id}"
      placeholder="Enter section description (optional)"
      value={question.description || ''}
      oninput={(e) => store.updateQuestion(questionId, { description: e.currentTarget.value })}
      rows={3}
    />
  </div>
</div>
{/if}
