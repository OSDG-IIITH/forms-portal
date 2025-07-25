<script lang="ts">
  import { onMount } from 'svelte';
  import * as kdljs from 'kdljs';
  import { ulid } from 'ulid';
  import InputView from '$lib/components/viewer/questions/input.svelte';
  import TextareaView from '$lib/components/viewer/questions/textarea.svelte';
  import RadioView from '$lib/components/viewer/questions/radio.svelte';
  import CheckboxView from '$lib/components/viewer/questions/checkbox.svelte';
  import FileUploadView from '$lib/components/viewer/questions/file-upload.svelte';
  import SelectView from '$lib/components/viewer/questions/select.svelte';
  import DateView from '$lib/components/viewer/questions/date.svelte';
  import { Button } from '$lib/components/ui/button';
  import { Card, CardContent } from '$lib/components/ui/card';
  import { toast } from 'svelte-sonner';

  export let form: any;

  interface Option {
    id: string;
    value: string;
    label: string;
  }

  type QuestionType =
    | 'input'
    | 'textarea'
    | 'radio'
    | 'checkbox'
    | 'file'
    | 'select'
    | 'date';

  interface Question {
    id: string;
    type: QuestionType;
    title: string;
    required: boolean;
    options?: Option[];
    placeholder?: string;
    validations?: {
      'max-chars'?: number;
      'min-chars'?: number;
      regex?: string;
    };
    'max-file-size'?: number;
    'max-files'?: number;
    'allowed-types'?: string[];
  }

  interface FormConfig {
    title: string;
    description: string;
    visibility: string;
  }

  let formConfig: FormConfig = {
    title: '',
    description: '',
    visibility: 'public'
  };
  let questions: Question[] = [];
  let responses: Record<string, string> = {};
  let isLoading = true;
  let isSubmitting = false;
  let error: string | null = null;

  function parseKdlValue(node: any): string {
    if (node === null || node === undefined) return '';
    if (typeof node === 'string') return node;
    if (typeof node.value === 'string') return node.value;
    if (typeof node.value === 'number' || typeof node.value === 'boolean') return String(node.value);
    return String(node);
  }

  function safeString(val: any): string {
    if (val === null || val === undefined) return '';
    if (typeof val === 'string') return val;
    if (typeof val === 'number' || typeof val === 'boolean') return String(val);
    return '';
  }

  function parseKdlForm(kdl: string) {
    const ast = kdljs.parse(kdl);
    if (!ast || !ast.output || !ast.output.length) throw new Error('Invalid KDL');
    const formNode = ast.output.find((n: any) => n.name === 'form');
    if (!formNode) throw new Error('No form node');
    const config: FormConfig = {
      title: '',
      description: '',
      visibility: 'public'
    };
    const qs: Question[] = [];
    for (const child of formNode.children) {
      if (child.name === 'title') config.title = parseKdlValue(child.values[0]);
      else if (child.name === 'description') config.description = parseKdlValue(child.values[0]);
      else if (child.name === 'visibility') config.visibility = parseKdlValue(child.values[0]);
      else if (child.name === 'question') {
        const q: Question = {
          id: '',
          type: 'input',
          title: '',
          required: false
        };
        if (child.properties && typeof child.properties === 'object') {
          for (const key in child.properties) {
            const value = child.properties[key];
            if (key === 'id') q.id = safeString(value);
            else if (key === 'type') {
              let t = safeString(value);
              if (t === 'multiple_choice') t = 'radio';
              if (t === 'text') t = 'input';
              if (t === 'textarea') t = 'textarea';
              q.type = t as QuestionType;
            }
          }
        }
        if (Array.isArray(child.values) && child.values.includes('required')) {
          q.required = true;
        }
        if (child.children) {
          for (const c of child.children) {
            if (c.name === 'title') q.title = parseKdlValue(c.values[0]);
            else if (c.name === 'placeholder') q.placeholder = parseKdlValue(c.values[0]);
            else if (c.name === 'max-file-size') {
              q['max-file-size'] = Number(parseKdlValue(c.values[0]));
            }
            else if (c.name === 'max-files') {
              q['max-files'] = Number(parseKdlValue(c.values[0]));
            }
            else if (c.name === 'allowed-types') {
              q['allowed-types'] = c.values?.map((v: any) => parseKdlValue(v)) || [];
            }
            else if (c.name === 'option') {
              if (!q.options) q.options = [];
              let id = '', value = '', label = '';
              if (c.properties && typeof c.properties === 'object') {
                for (const key in c.properties) {
                  const v = c.properties[key];
                  if (key === 'value') value = safeString(v);
                  else if (key === 'label') label = safeString(v);
                  else if (key === 'id') id = safeString(v);
                }
              }
              q.options.push({ id: id || ulid(), value, label });
            }
            else if (c.name === 'validations') {
              q.validations = {};
              for (const v of c.children || []) {
                if (v.name === 'regex') q.validations.regex = parseKdlValue(v.values[0]);
                else if (v.name === 'min-chars') q.validations['min-chars'] = Number(parseKdlValue(v.values[0]));
                else if (v.name === 'max-chars') q.validations['max-chars'] = Number(parseKdlValue(v.values[0]));
              }
            }
          }
        }
        qs.push(q);
      }
    }
    return { config, qs };
  }

  async function loadFormKdl() {
    try {
      isLoading = true;
      error = null;
      const kdl = form.structure;
      const { config, qs } = parseKdlForm(kdl);
      formConfig = config;
      questions = qs;
    } catch (e) {
      error = 'Failed to load form.';
    } finally {
      isLoading = false;
    }
  }

  onMount(() => {
    loadFormKdl();
  });

  function getInitialResponses() {
    const initialResponses: Record<string, string> = {};
    for (const q of questions) {
      if (q.type === 'checkbox') initialResponses[q.id] = '[]';
      else initialResponses[q.id] = '';
    }
    return initialResponses;
  }

  $: if (!isLoading && !error) {
    responses = getInitialResponses();
  }

  function handleReset() {
    responses = getInitialResponses();
  }

  async function handleSubmit() {
    isSubmitting = true;
    error = null;
    try {
      const startRes = await fetch(`/api/forms/${form.id}/responses`, { method: 'POST' });
      if (!startRes.ok) throw new Error('Could not start response');
      const responseObj = await startRes.json();
      const responseId = responseObj.id;
      for (const q of questions) {
        let valueToSubmit: any;
        const rawValue = responses[q.id];

        if (q.type === 'radio') {
          const selectedOption = q.options?.find((opt) => opt.id === rawValue);
          valueToSubmit = selectedOption ? selectedOption.value : '';
        } else if (q.type === 'checkbox') {
          const selectedIds = JSON.parse(rawValue);
          valueToSubmit =
            q.options
              ?.filter((opt) => selectedIds.includes(opt.id))
              .map((opt) => opt.value) || [];
        } else {
          valueToSubmit = rawValue;
        }

        const upsertRes = await fetch(`/api/forms/${form.id}/responses/${responseId}/answers`, {
          method: 'PUT',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ question: q.id, value: JSON.stringify(valueToSubmit) })
        });
        if (!upsertRes.ok) throw new Error('Could not save answer');
      }
      const submitRes = await fetch(`/api/forms/${form.id}/responses/${responseId}/submit`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ "save": true })
      });
      if (!submitRes.ok) throw new Error('Could not submit response');
      toast.success('Form submitted successfully!');
      handleReset();
    } catch (e: any) {
      toast.error(e.message || 'Submission failed');
    } finally {
      isSubmitting = false;
    }
    
  }
</script>

<svelte:head>
  <title>{formConfig.title || 'Form'} - Form Viewer</title>
  <meta name="description" content={formConfig.description || 'Form viewer'} />
</svelte:head>

<div class="min-h-screen bg-background">
  <main class="container mx-auto max-w-4xl px-6 py-8 pt-24">
    {#if isLoading}
      <div class="flex items-center justify-center py-12">
        <div class="text-center space-y-4">
          <div class="w-8 h-8 border-4 border-primary border-t-transparent rounded-full animate-spin mx-auto"></div>
          <p class="text-muted-foreground">Loading form...</p>
        </div>
      </div>
    {:else if error}
      <div class="text-center py-12">
        <div class="w-12 h-12 bg-destructive/10 rounded-full flex items-center justify-center mx-auto mb-4">
          <svg class="w-6 h-6 text-destructive" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z"></path>
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
        <form on:submit|preventDefault={handleSubmit}>
          <div class="space-y-6">
            {#each questions as question (question.id)}
              <div class="p-6 border rounded-lg bg-card shadow-xs space-y-2">
                <div class="font-medium text-base mb-4">
                  {question.title}{#if question.required}<span class="text-destructive ml-1">*</span>{/if}
                </div>
                {#if question.type === 'input'}
                  <InputView {question} bind:value={responses[question.id]} />
                {:else if question.type === 'textarea'}
                  <TextareaView {question} bind:value={responses[question.id]} />
                {:else if question.type === 'radio'}
                  <RadioView
                    question={{
                      id: question.id,
                      title: question.title,
                      required: question.required,
                      options: question.options ?? []
                    }}
                    bind:value={responses[question.id]}
                  />
                {:else if question.type === 'checkbox'}
                  <CheckboxView
                    question={{
                      id: question.id,
                      title: question.title,
                      required: question.required,
                      options: question.options ?? []
                    }}
                    bind:value={responses[question.id]}
                  />
                {:else if question.type === 'file'}
                  <FileUploadView {question} bind:value={responses[question.id]} />
                {:else if question.type === 'select'}
                  <SelectView {question} bind:value={responses[question.id]} />
                {:else if question.type === 'date'}
                  <DateView {question} bind:value={responses[question.id]} />
                {:else}
                  {#if question.options && question.options.length > 0}
                    <div class="mt-2">
                      <div class="font-semibold text-xs mb-1">Options:</div>
                      <ul class="list-disc ml-6">
                        {#each question.options as opt}
                          <li>{opt.label}</li>
                        {/each}
                      </ul>
                    </div>
                  {/if}
                  {#if question.placeholder}
                    <div class="mt-2 text-xs text-muted-foreground">Placeholder: {question.placeholder}</div>
                  {/if}
                  {#if question.validations}
                    <div class="mt-2 text-xs text-muted-foreground">Validations: {JSON.stringify(question.validations)}</div>
                  {/if}
                {/if}
              </div>
            {/each}
          </div>
          {#if questions.length > 0}
            <Card class="py-0 mt-8">
              <CardContent class="py-4">
                <div class="flex justify-between items-center">
                  <div class="flex items-center gap-2">
                    <span class="text-sm text-muted-foreground">
                      {questions.length} question{questions.length !== 1 ? 's' : ''}
                    </span>
                  </div>
                  <div class="flex gap-4">
                    <Button 
                      type="button" 
                      variant="outline" 
                      size="sm"
                      onclick={handleReset}
                      disabled={isSubmitting}
                    >
                      Reset Form
                    </Button>
                    <Button 
                      type="submit" 
                      size="sm"
                      disabled={isSubmitting}
                    >
                      {#if isSubmitting}
                        <div class="w-4 h-4 border-2 border-current border-t-transparent rounded-full animate-spin mr-2"></div>
                        Submitting...
                      {:else}
                        Submit Form
                      {/if}
                    </Button>
                  </div>
                </div>
              </CardContent>
            </Card>
          {/if}
        </form>
      </div>
    {/if}
  </main>
</div>