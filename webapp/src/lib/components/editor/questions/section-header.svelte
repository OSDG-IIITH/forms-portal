<script lang="ts">
  import { Input } from '$lib/components/ui/input';
  import { Label } from '$lib/components/ui/label';
  import { Textarea } from '$lib/components/ui/textarea';
  import { getContext } from 'svelte';
  import type { FormStore } from '../form-store.svelte';

  let { questionId }: { questionId: string } = $props();

  const store: FormStore = getContext('form-store');
  const question = $derived(store.questions.find((q) => q.id === questionId));
</script>

{#if question}
<div class="space-y-4">
  <div class="space-y-2">
    <Label for="section-title-{question.id}">Section Title</Label>
    <Input
      id="section-title-{question.id}"
      placeholder="Enter section title"
      value={question.title}
      oninput={(e) => store.updateQuestion(questionId, { title: e.currentTarget.value })}
    />
    {#if question.error}
      <p class="text-destructive text-sm">{question.error}</p>
    {/if}
  </div>

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
