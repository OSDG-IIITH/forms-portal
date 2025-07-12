<script lang="ts">
  import { Input } from '$lib/components/ui/input';
  import { Textarea } from '$lib/components/ui/textarea/index.js';

  interface Props {
    question: {
      uid: string;
      text: string;
      required: boolean;
    };
    value?: string;
  }

  let { question, value = $bindable('') }: Props = $props();
  
  // hack for textarea because i couldnt find anything for it in the schema
  // FIXME: later
  const isLongText = $derived(
    question.text.toLowerCase().includes('tell us') ||
    question.text.toLowerCase().includes('essay') ||
    question.text.toLowerCase().includes('describe') ||
    question.text.toLowerCase().includes('explain') ||
    question.text.toLowerCase().includes('feedback') ||
    question.text.toLowerCase().includes('details') ||
    question.text.toLowerCase().includes('thoughts') ||
    question.text.length > 50
  );

  const inputType = $derived(
    question.text.toLowerCase().includes('email') ? 'email' : 'text'
  );
</script>

{#if isLongText}
  <Textarea
    id="question-{question.uid}"
    bind:value
    placeholder="Type your message here"
    required={question.required}
    rows={4}
  />
{:else}
  <Input
    id="question-{question.uid}"
    bind:value
    placeholder="Enter your response here"
    required={question.required}
    type={inputType}
  />
{/if}
