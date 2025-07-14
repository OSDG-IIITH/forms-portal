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

  let form = $state<FormData>({
    title: 'Untitled Form',
    description: 'Add a description',
    visibility: 'private'
  });

  let questions = $state<Question[]>([]);

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

  function slugify(title: string, suffix = 0): string {
    let base = title
      .toLowerCase()
      .replace(/[^a-z0-9]+/g, '-')
      .replace(/(^-|-$)/g, '')
      .slice(0, 50);
    if (suffix > 0) return `${base}-${suffix}`.slice(0, 64);
    return base;
  }

  const BACKEND_URL = import.meta.env.VITE_BACKEND_URL || 'http://localhost:8647';

  async function saveFormAsKdl() {
    let kdl = 'form {\n';
    kdl += `  version 1\n`;
    for (const q of questions) {
      kdl += questionToKdl(q);
    }
    kdl += '}\n';

    try {
      kdljs.parse(kdl);
    } catch (e) {
      toast.error("Failed to generate valid KDL", {
        description: e instanceof Error ? e.message : String(e)
      });
      return;
    }

    /* 
    if (form.title.trim() === '') {
      toast.error("Form title is required", { description: "Please enter a title for your form." });
      return;
    }
    */

    if (form.title.trim() === '') {
      form.title = 'Untitled Form';
    }

    const blob = new Blob([kdl], { type: 'text/plain' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `${form.title || 'form'}.kdl`;
    document.body.appendChild(a);
    a.click();
    setTimeout(() => {
      document.body.removeChild(a);
      URL.revokeObjectURL(url);
    }, 0);

    let suffix = 0;
    for (let attempt = 0; attempt < 100; attempt++) {
      const payload = {
        // title: suffix === 0 ? form.title : `${form.title} (${suffix})`,
        title: form.title || 'Untitled Form',
        slug: slugify(form.title, suffix),
        description: form.description,
        structure: kdl
      };
      let res, text = '';
      try {
        res = await fetch(`${BACKEND_URL}/api/forms`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          credentials: 'include',
          body: JSON.stringify(payload)
        });
        text = await res.text();
      } catch (e) {
        toast.error("Network error", { description: e instanceof Error ? e.message : String(e) });
        return;
      }
      if (res.ok) {
        toast.success("Form created successfully");
        return;
      }
      if (
        text.includes('duplicate key value') ||
        text.includes('forms_owner_slug_key') ||
        res.status === 500
      ) {
        suffix++;
        continue;
      }
      toast.error("Error creating form", { description: text || 'Failed to create form.' });
      return;
    }
    toast.error("Error creating form", { description: 'Could not generate a unique title/slug after many attempts.' });
  }
</script>

<div class="min-h-screen bg-background">
  <main class="container mx-auto max-w-4xl px-6 py-8 pt-24">
    <div class="space-y-8">
      <!-- configuration -->
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
                bind:value={form.title}
              />
            </div>
            <div class="space-y-2">
              <Label for="form-description">Description</Label>
              <Textarea
                id="form-description"
                placeholder="Describe what this form is for"
                bind:value={form.description}
                rows={3}
              />
            </div>
          </CardContent>
        </Card>
      </section>

      <!-- questios -->
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

          <!-- add question -->
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
                    class="flex items-center gap-2"
                    onclick={() => addQuestion(type)}
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
                Get started by adding a question to your form
              </p>
              <div class="flex flex-wrap justify-center gap-2">
                {#each questionTypeButtons as { type, icon: Icon, label }}
                  <Button 
                    variant="outline" 
                    size="sm"
                    onclick={() => addQuestion(type)}
                  >
                    <Icon class="h-4 w-4 mr-2" />
                    {label}
                  </Button>
                {/each}
              </div>
            </CardContent>
          </Card>
        {/if}
      </section>

      <!-- footer/actions -->
      {#if questions.length > 0}
        <Card>
          <CardContent>
            <div class="flex justify-between items-center">
              <div class="flex items-center gap-2">
                <span class="text-sm text-muted-foreground font-medium">
                  {questions.length} question{questions.length !== 1 ? 's' : ''}
                </span>
              </div>
              <div class="flex gap-4">
                <Button variant="outline" size="sm">Preview</Button>
                <Button size="sm" onclick={saveFormAsKdl}>Save Form</Button>
              </div>
            </div>
          </CardContent>
        </Card>
      {/if}
    </div>
  </main>
</div>
