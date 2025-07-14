<script lang="ts">
  import { Input } from '$lib/components/ui/input';
  import { Label } from '$lib/components/ui/label';
  import { Checkbox } from '$lib/components/ui/checkbox';
  import { Toggle } from '$lib/components/ui/toggle';
  import {
    IconSquareLetterA,
    IconViewportTall,
    IconRegex,
    IconAt
  } from '@tabler/icons-svelte';
  import { slide } from 'svelte/transition';
  import { cubicOut } from 'svelte/easing';
  const transitionOpts = { duration: 480, easing: cubicOut };

  interface Props {
    question: {
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
  }

  export let question: Props['question'];

  if (!question.validations) {
    question.validations = {};
  }
  if (typeof question.validations.email !== 'boolean') {
    question.validations.email = false;
  }

  let showPlaceholder = !!question.placeholder;
  let showCharLimit =
    question.validations?.['min-chars'] !== undefined ||
    question.validations?.['max-chars'] !== undefined;
  let showRegex = !!question.validations?.regex;

  $: if (!showPlaceholder) {
    question.placeholder = undefined;
  }

  $: if (!showCharLimit) {
    if (question.validations) {
      question.validations['min-chars'] = undefined;
      question.validations['max-chars'] = undefined;
    }
  }
  $: if (!showRegex) {
    if (question.validations) {
      question.validations.regex = undefined;
    }
  }

  $: if (question.validations?.['min-chars']) {
    const strValue = String(question.validations['min-chars']);
    const sanitized = strValue.replace(/\D/g, '');
    if (strValue !== sanitized) {
      question.validations['min-chars'] = sanitized ? parseInt(sanitized, 10) : undefined;
    }
  }

  $: if (question.validations?.['max-chars']) {
    const strValue = String(question.validations['max-chars']);
    const sanitized = strValue.replace(/\D/g, '');
    if (strValue !== sanitized) {
      question.validations['max-chars'] = sanitized ? parseInt(sanitized, 10) : undefined;
    }
  }
</script>

<div class="space-y-2">
  <div class="space-y-2">
    <Label for="question-{question.id}">Question Text</Label>
    <Input
      id="question-{question.id}"
      placeholder="Enter your question here"
      bind:value={question.title}
    />
  </div>

  <div class="flex items-center gap-2 flex-wrap my-4">
    <Toggle bind:pressed={showPlaceholder} aria-label="Toggle placeholder" variant="outline">
      <IconSquareLetterA class="size-4" />
      <span class="text-sm font-medium">Placeholder</span>
    </Toggle>
    <Toggle bind:pressed={showCharLimit} aria-label="Toggle character limit" variant="outline">
      <IconViewportTall class="size-4" />
      <span class="text-sm font-medium">Character Limit</span>
    </Toggle>
    <Toggle bind:pressed={showRegex} aria-label="Toggle regex" variant="outline">
      <IconRegex class="size-4" />
      <span class="text-sm font-medium">Regex</span>
    </Toggle>
    {#if question.validations}
      <Toggle
        bind:pressed={question.validations.email}
        aria-label="Toggle email validation"
        variant="outline"
      >
        <IconAt class="size-4" />
        <span class="text-sm font-medium">Email</span>
      </Toggle>
    {/if}
  </div>

  {#if showPlaceholder}
    <div class="space-y-2 mb-4" transition:slide={transitionOpts}>
      <Label for="placeholder-{question.id}">Placeholder</Label>
      <Input
        id="placeholder-{question.id}"
        placeholder="Enter placeholder text"
        bind:value={question.placeholder}
      />
    </div>
  {/if}

  {#if question.validations}
    <div class="grid grid-cols-2 gap-4 mb-4">
      {#if showCharLimit}
        <div class="space-y-2" transition:slide={transitionOpts}>
          <Label for="min-chars-{question.id}">Min Characters</Label>
          <Input
            id="min-chars-{question.id}"
            placeholder="Minimum Character Limit"
            bind:value={question.validations['min-chars']}
          />
        </div>
        <div class="space-y-2" transition:slide={transitionOpts}>
          <Label for="max-chars-{question.id}">Max Characters</Label>
          <Input
            id="max-chars-{question.id}"
            placeholder="Maximum Character Limit"
            bind:value={question.validations['max-chars']}
          />
        </div>
      {/if}
      {#if showRegex}
        <div class="space-y-2 col-span-2" transition:slide={transitionOpts}>
          <Label for="regex-{question.id}">Regex Pattern</Label>
          <Input
            id="regex-{question.id}"
            placeholder="Regex Pattern"
            bind:value={question.validations.regex}
          />
        </div>
      {/if}
    </div>
  {/if}

  <div class="flex items-center gap-2">
    <Checkbox id="required-{question.id}" bind:checked={question.required} />
    <Label for="required-{question.id}" class="text-sm"> Required </Label>
  </div>
</div>
