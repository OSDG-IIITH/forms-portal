<script lang="ts">
  import { Input } from '$lib/components/ui/input';
  import { Label } from '$lib/components/ui/label';
  import { Checkbox } from '$lib/components/ui/checkbox';
  import Chip from '../chips/Chip.svelte';
  import AddChip from '../chips/AddChip.svelte';
  import {
    IconSquareLetterA,
    IconViewportTall,
    IconRegex,
    IconNumber123,
    IconAt
  } from '@tabler/icons-svelte';
  import { fly, slide } from 'svelte/transition';
  import { cubicOut } from 'svelte/easing';
  import { tick } from 'svelte';
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
      };
    };
  }

  export let question: Props['question'];

  if (!question.validations) {
    question.validations = {};
  }

  let chips: Array<{
    key: string;
    label: string;
    icon: any;
    active: boolean;
    clickable: boolean;
  }> = [];

  let activeChip: string | null = null;

  let regexType: 'number' | 'email' | 'regex' | null = null;

  let chipOrder: string[] = [];

  $: {
    const regex = question.validations?.regex;
    const emailPattern = '^[^\\s@]+@[^\\s@]+\\.[^\\s@]{2,}$';
    if (typeof regex === 'string') {
      if (regex.length > 0) {
        if (/^\\d\+$/.test(regex) || /^\^\\d\+\$$/.test(regex)) regexType = 'number';
        else if (regex === emailPattern) regexType = 'email';
        else regexType = 'regex';
      } else {
        regexType = 'regex';
      }
    } else {
      regexType = null;
    }
  }

  $: {
    const activeKeys = chips.filter(v => v.active).map(v => v.key);
    for (const key of activeKeys) {
      if (!chipOrder.includes(key)) chipOrder.push(key);
    }
    chipOrder = chipOrder.filter(key => activeKeys.includes(key));
  }

  function handleChipClick(key: string) {
    activeChip = activeChip === key ? null : key;
  }

  function removeChip(key: string) {
    if (key === 'placeholder') {
      question.placeholder = undefined;
      question = { ...question };
    }
    if (key === 'charlimit' && question.validations) {
      delete question.validations['min-chars'];
      delete question.validations['max-chars'];
    }
    if ((key === 'regex' || key === 'number' || key === 'email') && question.validations) {
      question.validations.regex = undefined;
      regexType = null;
    }
    chipOrder = chipOrder.filter(k => k !== key);
    activeChip = null;
    question = { ...question };
  }

  function addValidation(key: string) {
    if (["number", "email", "regex"].includes(key)) {
      if (question.validations) question.validations.regex = undefined;
      regexType = null;
    }
    if (key === 'placeholder') {
      question.placeholder = '';
      question = { ...question };
    }
    if (key === 'charlimit') {
      if (!question.validations) question.validations = {};
      question.validations['min-chars'] = undefined;
      question.validations['max-chars'] = undefined;
    }
    if (key === 'regex') {
      if (!question.validations) question.validations = {};
      question.validations.regex = '';
      regexType = 'regex';
    }
    if (key === 'number') {
      if (!question.validations) question.validations = {};
      question.validations.regex = "^\\d+$";
      regexType = 'number';
    }
    if (key === 'email') {
      if (!question.validations) question.validations = {};
      question.validations.regex = '^[^\\s@]+@[^\\s@]+\\.[^\\s@]{2,}$';
      regexType = 'email';
      console.log(question.placeholder)
    }
    if (!chipOrder.includes(key)) chipOrder.push(key);
    activeChip = key;
    question = { ...question };
  }

  function updateMinChars(value: string): void {
    if (!question.validations) question.validations = {};
    const sanitized = value.replace(/\D/g, '');
    question.validations['min-chars'] = sanitized ? parseInt(sanitized, 10) : undefined;
  }

  function updateMaxChars(value: string): void {
    if (!question.validations) question.validations = {};
    const sanitized = value.replace(/\D/g, '');
    question.validations['max-chars'] = sanitized ? parseInt(sanitized, 10) : undefined;
  }

  function focusValidationInput(key: string) {
    tick().then(() => {
      if (key === 'placeholder') {
        const el = document.getElementById(`placeholder-${question.id}`);
        el?.focus();
      } else if (key === 'charlimit') {
        const el = document.getElementById(`min-chars-${question.id}`);
        el?.focus();
      } else if (key === 'regex' || key === 'number' || key === 'email') {
        const el = document.getElementById(`regex-${question.id}`);
        el?.focus();
      }
    });
  }

  $: chips = [
    {
      key: 'placeholder',
      label: 'Placeholder',
      icon: IconSquareLetterA,
      active: question.placeholder !== undefined,
      clickable: true
    },
    {
      key: 'charlimit',
      label: 'Character Limit',
      icon: IconViewportTall,
      active: 'min-chars' in (question.validations ?? {}) || 'max-chars' in (question.validations ?? {}),
      clickable: true
    },
    {
      key: 'number',
      label: 'Number',
      icon: IconNumber123,
      active: regexType === 'number',
      clickable: false
    },
    {
      key: 'email',
      label: 'Email',
      icon: IconAt,
      active: regexType === 'email',
      clickable: false
    },
    {
      key: 'regex',
      label: 'Regex',
      icon: IconRegex,
      active: regexType === 'regex',
      clickable: true
    }
  ];

  $: availableToAdd = chips.filter(v => !v.active);

  $: showPlaceholder = !!question.placeholder;
  $: showCharLimit =
    !!(question.validations &&
    (question.validations['min-chars'] !== undefined ||
    question.validations['max-chars'] !== undefined));
  $: showRegex = !!(question.validations && question.validations.regex);

  function handleInputMinChars(event: Event) {
    const target = event.target as HTMLInputElement;
    updateMinChars(target.value);
  }
  function handleInputMaxChars(event: Event) {
    const target = event.target as HTMLInputElement;
    updateMaxChars(target.value);
  }
  function handleInputRegex(event: Event) {
    if (!question.validations) question.validations = {};
    const target = event.target as HTMLInputElement;
    question.validations.regex = target.value;
  }

  function getActiveValidationByKey(key: string) {
    return chips.find(v => v.key === key && v.active);
  }
</script>

<div class="space-y-2" transition:slide={transitionOpts}>
  <div class="space-y-2">
    <Label for="question-{question.id}">Question Text</Label>
    <Input
      id="question-{question.id}"
      placeholder="Enter your question here"
      bind:value={question.title}
    />
  </div>

  <div class="my-4">
    <div class="flex items-center gap-2 flex-wrap">
      {#each chipOrder as key (key)}
        {@const v = getActiveValidationByKey(key)}
        {#if v}
          <div transition:fly={{ y: 10, opacity: 0, duration: 220, easing: cubicOut }}>
            <Chip
              label={v.label}
              icon={v.icon}
              active={v.active && v.clickable ? activeChip === v.key : false}
              clickable={v.clickable}
              on:click={() => v.clickable && handleChipClick(v.key)}
              onRemove={() => removeChip(v.key)}
            />
          </div>
        {/if}
      {/each}
      {#if availableToAdd.length > 0}
        <AddChip available={availableToAdd} onAdd={addValidation} on:focusValidationInput={e => focusValidationInput(e.detail.key)} totalValidations={chips.length} />
      {/if}
    </div>
  </div>

  {#if activeChip === 'placeholder'}
    <div class="space-y-2 mb-4" transition:slide={transitionOpts}>
      <Label for="placeholder-{question.id}">Placeholder</Label>
      <Input
        id="placeholder-{question.id}"
        placeholder="Enter placeholder text"
        bind:value={question.placeholder}
      />
    </div>
  {/if}

  {#if activeChip === 'charlimit'}
    <div class="grid grid-cols-2 gap-4 mb-4" transition:slide={transitionOpts}>
      <div class="space-y-2">
        <Label for="min-chars-{question.id}">Min Characters</Label>
        <Input
          id="min-chars-{question.id}"
          placeholder="Minimum Character Limit"
          value={question.validations && question.validations['min-chars'] || ''}
          oninput={handleInputMinChars}
        />
      </div>
      <div class="space-y-2">
        <Label for="max-chars-{question.id}">Max Characters</Label>
        <Input
          id="max-chars-{question.id}"
          placeholder="Maximum Character Limit"
          value={question.validations && question.validations['max-chars'] || ''}
          oninput={handleInputMaxChars}
        />
      </div>
    </div>
  {/if}

  {#if activeChip === 'regex'}
    <div class="space-y-2 mb-4" transition:slide={transitionOpts}>
      <Label for="regex-{question.id}">Regex Pattern</Label>
      <Input
        id="regex-{question.id}"
        placeholder="Regex Pattern"
        value={question.validations && question.validations.regex || ''}
        oninput={handleInputRegex}
      />
    </div>
  {/if}

  <div class="flex items-center gap-2">
    <Checkbox id="required-{question.id}" bind:checked={question.required} />
    <Label for="required-{question.id}" class="text-sm"> Required </Label>
  </div>
</div>
