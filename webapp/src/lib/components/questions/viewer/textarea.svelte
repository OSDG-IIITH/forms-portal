<script lang="ts">
  import Textarea from '../../ui/textarea/textarea.svelte';
  export let question: {
    id: string;
    title: string;
    required: boolean;
    placeholder?: string;
    validations?: {
      'max-chars'?: number;
      'min-chars'?: number;
      regex?: string;
      email?: boolean;
    };
  };
  export let value: string = '';
  export let disabled: boolean = false;

  let error: string = '';

  $: if (question.validations?.['min-chars'] && value.length < question.validations['min-chars']) {
    error = `Minimum ${question.validations['min-chars']} characters required.`;
  } else if (question.validations?.['max-chars'] && value.length > question.validations['max-chars']) {
    error = `Maximum ${question.validations['max-chars']} characters allowed.`;
  } else {
    error = '';
  }
</script>

<Textarea
  class="block w-full rounded-md border border-input bg-background px-3 py-2 text-base focus:outline-none disabled:opacity-50"
  id={question.id}
  name={question.id}
  placeholder={question.placeholder || ''}
  bind:value
  {disabled}
  aria-required={question.required}
  maxlength={question.validations?.['max-chars']}
  minlength={question.validations?.['min-chars']}
  rows={4}
/>
{#if value && error}
  <p class="text-sm text-destructive mt-1">{error}</p>
{/if}
