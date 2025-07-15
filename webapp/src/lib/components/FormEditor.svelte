<script lang="ts">
  import { Card, CardContent } from '$lib/components/ui/card';
  import { Button } from '$lib/components/ui/button';
  import { Input } from '$lib/components/ui/input';
  import { Textarea } from '$lib/components/ui/textarea';
  import { Label } from '$lib/components/ui/label';
  import InputEditor from '$lib/components/questions/editor/input.svelte';
  import TextareaEditor from '$lib/components/questions/editor/textarea.svelte';
  import RadioEditor from '$lib/components/questions/editor/radio.svelte';
  import CheckboxEditor from '$lib/components/questions/editor/checkbox.svelte';
  import FileUploadEditor from '$lib/components/questions/editor/file-upload.svelte';
  import SelectEditor from '$lib/components/questions/editor/select.svelte';
  import DatePickerEditor from '$lib/components/questions/editor/date-picker.svelte';
  import { ulid } from 'ulid';
  import * as kdljs from 'kdljs';
  import { toast } from "svelte-sonner";
  import {
    IconPlus,
    IconFileText,
    IconListCheck,
    IconChevronUp,
    IconChevronDown,
    IconTrash,
    IconSquareCheck,
    IconUpload,
    IconCalendar
  } from '@tabler/icons-svelte';
  import { createEventDispatcher } from 'svelte';

  const { form, mode = 'edit' }: { form: any; mode?: 'create' | 'edit' } = $props();

  type QuestionType = 'input' | 'textarea' | 'radio' | 'checkbox' | 'file' | 'select' | 'date';

  type Option = {
    id: string;
    value: string;
    label: string;
  };

  type Question = {
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
      email?: boolean;
    };
    'max-file-size'?: number;
    'max-files'?: number;
    'allowed-types'?: string[];
  };

  type FormData = {
    title: string;
    description: string;
    visibility: string;
  };

  let formData = $state<FormData>({
    title: 'Untitled Form',
    description: 'Add a description',
    visibility: 'private'
  });

  let questions = $state<Question[]>([]);
  let isLoading = $state(true);
  let isSaving = $state(false);
  let error = $state<string | null>(null);

  const questionTypeLabels: Record<QuestionType, string> = {
    input: 'Input',
    textarea: 'Textarea',
    radio: 'Multiple Choice',
    checkbox: 'Checkbox',
    file: 'File Upload',
    select: 'Select',
    date: 'Date Picker'
  };

  const questionTypeButtons = [
    { type: 'input' as QuestionType, icon: IconFileText, label: 'Text' },
    { type: 'textarea' as QuestionType, icon: IconFileText, label: 'Long Text' },
    { type: 'radio' as QuestionType, icon: IconListCheck, label: 'Multiple Choice' },
    { type: 'checkbox' as QuestionType, icon: IconSquareCheck, label: 'Checkbox' },
    { type: 'file' as QuestionType, icon: IconUpload, label: 'File Upload' },
    { type: 'select' as QuestionType, icon: IconListCheck, label: 'Select' },
    { type: 'date' as QuestionType, icon: IconCalendar, label: 'Date Picker' }
  ];

  const dispatch = createEventDispatcher<{
    create: { formData: FormData; questions: Question[]; kdl: string };
  }>();

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
    const qs: Question[] = [];
    for (const child of formNode.children) {
      if (child.name === 'question') {
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
              const optionId = id || (value && !q.options.some(opt => opt.id === value) ? value : ulid());
              q.options.push({ id: optionId, value, label });
            }
            else if (c.name === 'validations') {
              q.validations = {};
              for (const v of c.children || []) {
                if (v.name === 'regex') q.validations.regex = parseKdlValue(v.values[0]);
                else if (v.name === 'min-chars') q.validations['min-chars'] = Number(parseKdlValue(v.values[0]));
                else if (v.name === 'max-chars') q.validations['max-chars'] = Number(parseKdlValue(v.values[0]));
                else if (v.name === 'email') q.validations.email = true;
              }
            }
          }
        }
        qs.push(q);
      }
    }
    return qs;
  }

  function loadFormKdl() {
    try {
      isLoading = true;
      error = null;
      
      if (form) {
        formData.title = form.title || 'Untitled Form';
        formData.description = form.description || 'Add a description';
      }

      if (!form || !form.structure) {
        questions = [];
        isLoading = false;
        return;
      }

      const kdl = form.structure;
      const qs = parseKdlForm(kdl);
      questions = qs;
      isLoading = false;
    } catch (e) {
      console.error('Error parsing form KDL:', e);
      error = `Failed to load form: ${e instanceof Error ? e.message : 'Unknown error'}`;
      questions = [];
      isLoading = false;
    }
  }

  $effect(() => {
    if (form) {
      loadFormKdl();
    } else {
      questions = [];
      isLoading = false;
    }
  });

  function addQuestion(type: QuestionType): void {
    const newQuestion: Question = {
      id: ulid(),
      type,
      title: '',
      required: false
    };
    if (type === 'radio' || type === 'checkbox' || type === 'select') {
      newQuestion.options = [];
    }
    if (type === 'input' || type === 'textarea') {
      newQuestion.placeholder = '';
      newQuestion.validations = {};
    }
    if (type === 'file') {
      newQuestion['max-file-size'] = 10;
      newQuestion['max-files'] = -1;
      newQuestion['allowed-types'] = [];
    }
    questions = [...questions, newQuestion];
  }

  function removeQuestion(id: string): void {
    questions = questions.filter((q) => q.id !== id);
  }

  function moveQuestionUp(index: number): void {
    if (index > 0) {
      const newQuestions = [...questions];
      [newQuestions[index - 1], newQuestions[index]] = [newQuestions[index], newQuestions[index - 1]];
      questions = newQuestions;
    }
  }

  function moveQuestionDown(index: number): void {
    if (index < questions.length - 1) {
      const newQuestions = [...questions];
      [newQuestions[index], newQuestions[index + 1]] = [newQuestions[index + 1], newQuestions[index]];
      questions = newQuestions;
    }
  }

  function escapeKdlString(str: string): string {
    return '"' + str.replace(/\\/g, '\\\\').replace(/"/g, '\\"') + '"';
  }

  function questionToKdl(q: Question): string {
    let kdl = `  question id=${escapeKdlString(q.id)} type=${escapeKdlString(q.type)}${q.required ? ' required' : ''} {\n`;
    kdl += `    title ${escapeKdlString(q.title)}\n`;
    if (q.placeholder) {
      kdl += `    placeholder ${escapeKdlString(q.placeholder)}\n`;
    }
    if (q.options && q.options.length > 0) {
      for (const opt of q.options) {
        kdl += `    option value=${escapeKdlString(opt.value)} label=${escapeKdlString(opt.label)}\n`;
      }
    }
    if (q['max-file-size'] || q['max-files'] || q['allowed-types']?.length) {
      if (q['max-file-size']) {
        kdl += `    "max-file-size" ${q['max-file-size']} mb\n`;
      }
      if (q['max-files'] && q['max-files'] !== -1) {
        kdl += `    "max-files" ${q['max-files']}\n`;
      }
      if (q['allowed-types'] && q['allowed-types'].length > 0) {
        const typesString = q['allowed-types'].map(t => escapeKdlString(t)).join(' ');
        kdl += `    "allowed-types" ${typesString}\n`;
      }
    }
    if (q.validations && (q.validations['max-chars'] || q.validations['min-chars'] || q.validations.regex || q.validations.email)) {
      kdl += '    validations {\n';
      if (q.validations.regex) {
        kdl += `      regex ${escapeKdlString(q.validations.regex)}\n`;
      }
      if (q.validations['min-chars']) {
        kdl += `      "min-chars" ${q.validations['min-chars']}\n`;
      }
      if (q.validations['max-chars']) {
        kdl += `      "max-chars" ${q.validations['max-chars']}\n`;
      }
      if (q.validations.email) {
        kdl += `      email true\n`;
      }
      kdl += '    }\n';
    }
    kdl += `  }\n`;
    return kdl;
  }

  async function saveForm() {
    isSaving = true;
    try {
      let kdl = 'form {\n';
      kdl += `  version 1\n`;
      for (const q of questions) {
        kdl += questionToKdl(q);
      }
      kdl += '}\n';

      kdljs.parse(kdl);

      if (mode === 'create') {
        dispatch('create', { formData, questions, kdl });
        isSaving = false;
        return;
      }

      const payload = {
        title: formData.title || 'Untitled Form',
        description: formData.description,
        structure: kdl
      };

      const res = await fetch(`/api/forms/${form.id}`, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify(payload)
      });

      if (!res.ok) {
        const text = await res.text();
        throw new Error(text || 'Failed to save form');
      }

      toast.success("Form saved successfully");
    } catch (e: any) {
      toast.error("Error saving form", { description: e.message || 'Failed to save form.' });
    } finally {
      isSaving = false;
    }
  }
</script>

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
        <section class="space-y-6">
          <div>
            <h2 class="text-2xl font-semibold tracking-tight mb-2">Form Configuration</h2>
          </div>
          
          <Card>
            <CardContent class="space-y-6">
              <div class="space-y-2">
                <Label for="form-title">Form Title</Label>
                <Input
                  id="form-title"
                  placeholder="Enter form title"
                  bind:value={formData.title}
                />
              </div>
              <div class="space-y-2">
                <Label for="form-description">Description</Label>
                <Textarea
                  id="form-description"
                  placeholder="Describe what this form is for"
                  bind:value={formData.description}
                  rows={3}
                />
              </div>
            </CardContent>
          </Card>
        </section>

        <section class="space-y-6">
          {#if questions.length > 0}
            <div class="flex items-center justify-between">
              <div>
                <h2 class="text-2xl font-semibold tracking-tight mb-2">Form Questions</h2>
              </div>
              <div class="flex items-center gap-2 px-3 py-1.5 bg-muted rounded-md text-sm">
                <span class="font-medium">{questions.length}</span>
                <span class="text-muted-foreground">question{questions.length !== 1 ? 's' : ''}</span>
              </div>
            </div>
            
            <div class="space-y-4">
              {#each questions as question, i (question.id)}
                <Card>
                  <CardContent class="space-y-6">
                    <!-- header -->
                    <div class="flex items-center justify-between">
                      <div class="flex items-center gap-3">
                        <div class="flex items-center justify-center w-8 h-8 bg-muted text-sm font-medium rounded-md">
                          {i + 1}
                        </div>
                        <span class="font-medium">{questionTypeLabels[question.type]}</span>
                      </div>
                      <!-- controls -->
                      <div class="flex items-center gap-1">
                        <Button
                          variant="ghost"
                          size="icon"
                          class="h-8 w-8"
                          onclick={() => moveQuestionUp(i)}
                          disabled={i === 0}
                        >
                          <IconChevronUp class="h-4 w-4" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="icon"
                          class="h-8 w-8"
                          onclick={() => moveQuestionDown(i)}
                          disabled={i === questions.length - 1}
                        >
                          <IconChevronDown class="h-4 w-4" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="icon"
                          class="h-8 w-8 hover:bg-destructive/10 hover:text-destructive"
                          onclick={() => removeQuestion(question.id)}
                        >
                          <IconTrash class="h-4 w-4" />
                        </Button>
                      </div>
                    </div>
                    <!-- editor container -->
                    <div class="pl-4">
                      {#if question.type === 'input'}
                        <InputEditor bind:question={questions[i]} />
                      {:else if question.type === 'textarea'}
                        <TextareaEditor bind:question={questions[i]} />
                      {:else if question.type === 'radio'}
                        <RadioEditor bind:question={questions[i]} />
                      {:else if question.type === 'checkbox'}
                        <CheckboxEditor bind:question={questions[i]} />
                      {:else if question.type === 'file'}
                        <FileUploadEditor bind:question={questions[i]} />
                      {:else if question.type === 'select'}
                        <SelectEditor bind:question={questions[i]} />
                      {:else if question.type === 'date'}
                        <DatePickerEditor bind:question={questions[i]} />
                      {/if}
                    </div>
                  </CardContent>
                </Card>
              {/each}
            </div>

            <div class="border-2 border-dashed border-muted-foreground/25 rounded-lg p-6">
              <div class="text-center space-y-4">
                <div>
                  <h3 class="font-medium mb-1">Add Question</h3>
                  <p class="text-sm text-muted-foreground">Choose a question type</p>
                </div>
                <div class="flex flex-wrap justify-center gap-2">
                  {#each questionTypeButtons as { type, icon: Icon, label }}
                    <Button
                      variant="outline"
                      size="sm"
                      onclick={() => addQuestion(type)}
                      class="flex items-center gap-2"
                    >
                      <Icon class="h-4 w-4" />
                      {label}
                    </Button>
                  {/each}
                </div>
              </div>
            </div>
          {:else}
            <Card>
              <CardContent class="flex flex-col items-center justify-center py-12 text-center">
                <div class="w-12 h-12 bg-muted rounded-md flex items-center justify-center mb-4">
                  <IconPlus class="h-6 w-6 text-muted-foreground" />
                </div>
                <h3 class="text-lg font-medium mb-2">Create your first question</h3>
                <p class="text-muted-foreground mb-6 max-w-sm">
                  Get started by adding questions to your form
                </p>
                <div class="flex flex-wrap justify-center gap-2">
                  {#each questionTypeButtons as { type, icon: Icon, label }}
                    <Button
                      variant="outline"
                      size="sm"
                      onclick={() => addQuestion(type)}
                      class="flex items-center gap-2"
                    >
                      <Icon class="h-4 w-4" />
                      {label}
                    </Button>
                  {/each}
                </div>
              </CardContent>
            </Card>
          {/if}
        </section>

        {#if questions.length > 0}
          <Card class="py-0">
            <CardContent class="py-4">
              <div class="flex justify-between items-center">
                <div class="flex items-center gap-2">
                  <span class="text-sm text-muted-foreground">
                    {questions.length} question{questions.length !== 1 ? 's' : ''}
                  </span>
                </div>
                <div class="flex gap-4">
                  <Button
                    onclick={saveForm}
                    disabled={isSaving}
                    size="sm"
                  >
                    {#if isSaving}
                      <div class="w-4 h-4 border-2 border-current border-t-transparent rounded-full animate-spin mr-2"></div>
                      Saving...
                    {:else}
                      Save Form
                    {/if}
                  </Button>
                </div>
              </div>
            </CardContent>
          </Card>
        {/if}
      </div>
    {/if}
  </main>
</div>