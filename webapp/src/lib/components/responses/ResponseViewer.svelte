<script lang="ts">
  import { onMount } from 'svelte';
  import InputView from '$lib/components/viewer/questions/input.svelte';
  import TextareaView from '$lib/components/viewer/questions/textarea.svelte';
  import RadioView from './questions/radio.svelte';
  import CheckboxView from './questions/checkbox.svelte';
  import FileUploadView from './questions/file.svelte';
  import SelectView from '$lib/components/viewer/questions/select.svelte';
  import DateView from './questions/date.svelte';
  import { parseKdlForm } from '$lib/utils/kdl';
  import type { Question, FormConfig } from '$lib/types/form';

  let { form, answers }: { form: any; answers: any[] } = $props();

  let formConfig = $state<FormConfig>({
    title: '',
    description: '',
    visibility: 'public'
  });
  let questions = $state<Question[]>([]);
  let responses = $state<Record<string, string>>({});
  let isLoading = $state(true);
  let error = $state<string | null>(null);

  function populateResponses() {
    const populatedResponses: Record<string, string> = {};
    if (!answers || !Array.isArray(answers)) {
      error = 'Answers data is not available or is not an array.';
      return;
    }

    for (const q of questions) {
      const answer = answers.find((a) => a.question === q.id);
      if (answer && answer.value !== undefined && answer.value !== null) {
        try {
          const decodedJsonString = decodeURIComponent(
            atob(answer.value)
              .split('')
              .map((c) => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2))
              .join('')
          );

          const actualValue = JSON.parse(decodedJsonString);

          if (q.type === 'checkbox' || q.type === 'file' || typeof actualValue === 'object') {
            populatedResponses[q.id] = JSON.stringify(actualValue);
          } else {
            populatedResponses[q.id] = String(actualValue);
          }
        } catch (e) {
          console.error('Failed to decode or parse answer value:', answer.value, e);
          populatedResponses[q.id] = q.type === 'checkbox' || q.type === 'file' ? '[]' : '';
        }
      } else {
        populatedResponses[q.id] = q.type === 'checkbox' || q.type === 'file' ? '[]' : '';
      }
    }
    responses = populatedResponses;
  }

  onMount(() => {
    try {
      const kdl = form.structure;
      const result = parseKdlForm(kdl);
      formConfig = result.config;
      questions = result.questions;
      populateResponses();
    } catch (e) {
      error = 'Failed to load form and answers.';
    } finally {
      isLoading = false;
    }
  });
</script>

<svelte:head>
  <title>{formConfig.title || 'Form'} - Response Viewer</title>
  <meta name="description" content={formConfig.description || 'Response viewer'} />
</svelte:head>

<div class="min-h-screen bg-background">
  <main class="container mx-auto max-w-4xl px-6 py-8 pt-24">
    {#if isLoading}
      <div class="flex items-center justify-center py-12">
        <div class="text-center space-y-4">
          <div
            class="w-8 h-8 border-4 border-primary border-t-transparent rounded-full animate-spin mx-auto"
          ></div>
          <p class="text-muted-foreground">Loading form...</p>
        </div>
      </div>
    {:else if error}
      <div class="text-center py-12">
        <div
          class="w-12 h-12 bg-destructive/10 rounded-full flex items-center justify-center mx-auto mb-4"
        >
          <svg class="w-6 h-6 text-destructive" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z"
            ></path>
          </svg>
        </div>
        <h3 class="text-lg font-semibold mb-2">Error Loading Form</h3>
        <p class="text-muted-foreground mb-6">{error}</p>
      </div>
    {:else}
      <div class="space-y-8">
        <div class="space-y-3">
          <h1 class="text-3xl font-bold tracking-tight">{form.title}</h1>
          {#if form.description && form.description !== 'Add a description'}
            <p class="text-lg text-muted-foreground leading-relaxed">
              {form.description}
            </p>
          {/if}
        </div>
        <div class="space-y-6">
          {#each questions as question (question.id)}
            <div class="p-6 border rounded-lg bg-card shadow-xs space-y-2">
              <div class="font-medium text-base mb-4">
                {question.title}{#if question.required}<span class="text-destructive ml-1"
                    >*</span
                  >{/if}
              </div>
              {#if question.type === 'input'}
                <InputView {question} value={responses[question.id]} disabled={true} />
              {:else if question.type === 'textarea'}
                <TextareaView {question} value={responses[question.id]} disabled={true} />
              {:else if question.type === 'radio'}
                <RadioView
                  question={{
                    id: question.id,
                    title: question.title,
                    required: question.required,
                    options: question.options ?? []
                  }}
                  value={responses[question.id]}
                  disabled={true}
                  
                />
              {:else if question.type === 'checkbox'}
                <CheckboxView
                  question={{
                    id: question.id,
                    title: question.title,
                    required: question.required,
                    options: question.options ?? []
                  }}
                  value={responses[question.id]}
                  disabled={true}
                />
              {:else if question.type === 'file'}
                <FileUploadView value={responses[question.id]} />
              {:else if question.type === 'select'}
                <SelectView {question} value={responses[question.id]} disabled={true} />
              {:else if question.type === 'date'}
                <DateView value={responses[question.id]} />
              {:else}
                <p class="text-sm text-muted-foreground">
                  No view available for question type: {question.type}
                </p>
              {/if}
            </div>
          {/each}
        </div>
      </div>
    {/if}
  </main>
</div>