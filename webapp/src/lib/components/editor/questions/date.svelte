<script lang="ts">
  import { Input } from '$lib/components/ui/input';
  import { Label } from '$lib/components/ui/label';
  import { Checkbox } from '$lib/components/ui/checkbox';
  import { getContext } from 'svelte';
  import type { FormStore } from '../form-store.svelte';

  let { questionId }: { questionId: string } = $props();

  const store: FormStore = getContext('form-store');
  const question = $derived(store.questions.find((q) => q.id === questionId));
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

  <div class="flex items-center space-x-2">
    <Checkbox 
      id="required-{question.id}" 
      checked={question.required} 
      onchange={() => store.updateQuestion(questionId, { required: !question.required })}
    />
    <Label for="required-{question.id}" class="text-sm font-normal cursor-pointer">
      Required field
    </Label>
  </div>
</div>
{/if}