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
  import { tick, getContext } from 'svelte';
  import type { FormStore } from '../form-store.svelte';

  const transitionOpts = { duration: 480, easing: cubicOut };

  let { questionId }: { questionId: string } = $props();

  const store: FormStore = getContext('form-store');
  const question = $derived(store.questions.find((q) => q.id === questionId));

  let activeChip: string | null = $state(null);
  let chipOrder: string[] = $state([]);

  const regexType = $derived.by(() => {
    if (!question?.validations?.regex) return null;
    const regex = question.validations.regex;
    const emailPattern = '^[^\\s@]+@[^\\s@]+\\.[^\\s@]{2,}$';
    if (/^\\d\+$/.test(regex) || /^\^\\d\+\$$/.test(regex)) return 'number';
    if (regex === emailPattern) return 'email';
    return 'regex';
  });

  const chips = $derived([
    {
      key: 'placeholder',
      label: 'Placeholder',
      icon: IconSquareLetterA,
      active: question?.placeholder !== undefined,
      clickable: true
    },
    {
      key: 'charlimit',
      label: 'Character Limit',
      icon: IconViewportTall,
      active:
        'min-chars' in (question?.validations ?? {}) || 'max-chars' in (question?.validations ?? {}),
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
  ]);

  $effect(() => {
    const activeKeys = chips.filter((v) => v.active).map((v) => v.key);
    const newOrder = chipOrder.filter((key) => activeKeys.includes(key));
    for (const key of activeKeys) {
      if (!newOrder.includes(key)) newOrder.push(key);
    }
    chipOrder = newOrder;
  });

  function handleChipClick(key: string) {
    activeChip = activeChip === key ? null : key;
  }

  function removeChip(key: string) {
    if (!question) return;
    const newValidations = { ...question.validations };
    let update: Partial<typeof question> = {};

    if (key === 'placeholder') {
      update.placeholder = undefined;
    }
    if (key === 'charlimit') {
      delete newValidations['min-chars'];
      delete newValidations['max-chars'];
      update.validations = newValidations;
    }
    if (key === 'regex' || key === 'number' || key === 'email') {
      delete newValidations.regex;
      update.validations = newValidations;
    }

    store.updateQuestion(questionId, update);
    chipOrder = chipOrder.filter((k) => k !== key);
    activeChip = null;
  }

  function addValidation(key: string) {
    if (!question) return;
    const newValidations = { ...question.validations };
    let update: Partial<typeof question> = {};

    if (key === 'placeholder') {
      update.placeholder = '';
    }
    if (key === 'charlimit') {
      newValidations['min-chars'] = undefined;
      newValidations['max-chars'] = undefined;
    }
    if (key === 'regex') {
      newValidations.regex = '';
    }
    if (key === 'number') {
      newValidations.regex = '^\\d+$';
    }
    if (key === 'email') {
      newValidations.regex = '^[^\\s@]+@[^\\s@]+\\.[^\\s@]{2,}$';
    }

    update.validations = newValidations;
    store.updateQuestion(questionId, update);

    if (!chipOrder.includes(key)) chipOrder.push(key);
    activeChip = key;
  }

  function updateValidation(field: 'min-chars' | 'max-chars' | 'regex', value: string) {
    if (!question?.validations) return;
    const newValidations = { ...question.validations };
    if (field === 'regex') {
      newValidations.regex = value;
    } else {
      const sanitized = value.replace(/\D/g, '');
      newValidations[field] = sanitized ? parseInt(sanitized, 10) : undefined;
    }
    store.updateQuestion(questionId, { validations: newValidations });
  }

  function focusValidationInput(key: string) {
    tick().then(() => {
      const el = document.getElementById(`${key}-${questionId}`);
      el?.focus();
    });
  }

  const availableToAdd = $derived(chips.filter((v) => !v.active));

  function getActiveValidationByKey(key: string) {
    return chips.find((v) => v.key === key && v.active);
  }
</script>

{#if question}
  <div class="space-y-2" transition:slide={transitionOpts}>
    <div class="space-y-2">
      <Label for="question-{question.id}">Question Text</Label>
      <Input
        id="question-{question.id}"
        placeholder="Enter your question here"
        value={question.title}
        oninput={(e) => store.updateQuestion(questionId, { title: e.currentTarget.value })}
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
          <AddChip
            available={availableToAdd}
            onAdd={addValidation}
            on:focusValidationInput={(e) => focusValidationInput(e.detail.key)}
            totalValidations={chips.length}
          />
        {/if}
      </div>
    </div>

    {#if activeChip === 'placeholder'}
      <div class="space-y-2 mb-4" transition:slide={transitionOpts}>
        <Label for="placeholder-{question.id}">Placeholder</Label>
        <Input
          id="placeholder-{question.id}"
          placeholder="Enter placeholder text"
          value={question.placeholder}
          oninput={(e) => store.updateQuestion(questionId, { placeholder: e.currentTarget.value })}
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
            value={(question.validations && question.validations['min-chars']) || ''}
            oninput={(e) => updateValidation('min-chars', e.currentTarget.value)}
          />
        </div>
        <div class="space-y-2">
          <Label for="max-chars-{question.id}">Max Characters</Label>
          <Input
            id="max-chars-{question.id}"
            placeholder="Maximum Character Limit"
            value={(question.validations && question.validations['max-chars']) || ''}
            oninput={(e) => updateValidation('max-chars', e.currentTarget.value)}
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
          value={(question.validations && question.validations.regex) || ''}
          oninput={(e) => updateValidation('regex', e.currentTarget.value)}
        />
      </div>
    {/if}

    <div class="flex items-center gap-2">
      <Checkbox
        id="required-{question.id}"
        checked={question.required}
        onchange={() => store.updateQuestion(questionId, { required: !question.required })}
      />
      <Label for="required-{question.id}" class="text-sm"> Required </Label>
    </div>
  </div>
{/if}
