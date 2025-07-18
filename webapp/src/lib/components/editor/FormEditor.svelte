<script lang="ts">
  import FormTitleDescription from './FormTitleDescription.svelte';
  import FormConfig from './FormConfig.svelte';
  import QuestionsList from './QuestionsList.svelte';
  import AddQuestionPanel from './AddQuestionPanel.svelte';
  import FormSaveFooter from './FormActionFooter.svelte';
  import { ulid } from 'ulid';
  import * as kdljs from 'kdljs';
  import { toast } from "svelte-sonner";
  import {
    IconFileText,
    IconListCheck,
    IconSquareCheck,
    IconUpload,
    IconCalendar
  } from '@tabler/icons-svelte';
  import { createEventDispatcher, getContext } from 'svelte';
  import { Time } from '@internationalized/date';
  import type {
    QuestionType,
    Option,
    Question,
    FormData
  } from './FormEditor.ts';
  import {
    parseKdlValue,
    safeString,
    escapeKdlString,
    questionToKdl,
    generateFormKdl
  } from './FormEditor';

  const { form, mode = 'edit' }: { form: any; mode?: 'create' | 'edit' } = $props();

  let formData = $state<FormData>({
    title: 'Untitled Form',
    description: 'Add a description',
    visibility: 'private',
    opens: undefined,
    closes: undefined,
    opensTime: undefined,
    closesTime: undefined,
    anonymous: false,
    max_responses: null,
    individual_limit: 1,
    editable_responses: false
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

  let isPanelOpen = $state(false);
  let windowWidth = $state(0);
  let dialogOpen = $state(false);

  function handlePanelChange(event: CustomEvent<{ open: boolean }>) {
    isPanelOpen = event.detail.open;
  }

  function handleDialogOpenChange(e: CustomEvent<boolean>) {
    dialogOpen = e.detail;
    if (windowWidth > 0 && windowWidth < 1380 && !e.detail) {
      isPanelOpen = false;
    }
  }

  const getDisplacement = (width: number, isOpen: boolean) => {
    if (!isOpen) return '0px';
    
    if (width >= 1536) return '-240px';
    if (width >= 1280) return '-200px';
    if (width >= 1024) return '-160px';
    if (width >= 768) return '-120px';
    return '-80px';
  };

  $effect(() => {
    if (typeof window !== 'undefined') {
      windowWidth = window.innerWidth;
      
      const handleResize = () => {
        windowWidth = window.innerWidth;
      };
      
      window.addEventListener('resize', handleResize);
      return () => window.removeEventListener('resize', handleResize);
    }
  });

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
        
        if (form.opens) {
          const [dateStr, timeStr] = form.opens.split('T');
          formData.opens = dateStr;
          if (timeStr) {
            const [hours, minutes] = timeStr.split(':').map(Number);
            formData.opensTime = new Time(hours || 0, minutes || 0);
          }
        } else {
          formData.opens = undefined;
          formData.opensTime = undefined;
        }
        
        if (form.closes) {
          const [dateStr, timeStr] = form.closes.split('T');
          formData.closes = dateStr;
          if (timeStr) {
            const [hours, minutes] = timeStr.split(':').map(Number);
            formData.closesTime = new Time(hours || 0, minutes || 0);
          }
        } else {
          formData.closes = undefined;
          formData.closesTime = undefined;
        }
        
        formData.anonymous = form.anonymous || false;
        formData.max_responses = form.max_responses ?? null;
        formData.individual_limit = form.individual_limit ?? 1;
        formData.editable_responses = form.editable_responses || false;
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

  async function saveForm() {
    isSaving = true;
    try {
      const kdl = generateFormKdl(questions);

      kdljs.parse(kdl);

      if (mode === 'create') {
        dispatch('create', { formData, questions, kdl });
        isSaving = false;
        return;
      }

      const formatDateTime = (dateStr: string | undefined, time: Time | undefined): string | null => {
        if (!dateStr || dateStr.trim() === '') return null;
        try {
          const hours = time?.hour ?? 0;
          const minutes = time?.minute ?? 0;
          const seconds = 0;
          
          const date = new Date(dateStr);
          if (isNaN(date.getTime())) return null;
          
          date.setHours(hours, minutes, seconds, 0);
          return date.toISOString();
        } catch {
          return null;
        }
      };

      const payload = {
        title: formData.title || 'Untitled Form',
        description: formData.description || null,
        structure: kdl,
        opens: formatDateTime(formData.opens, formData.opensTime),
        closes: formatDateTime(formData.closes, formData.closesTime),
        anonymous: Boolean(formData.anonymous),
        max_responses: formData.max_responses !== null && formData.max_responses !== undefined ? Number(formData.max_responses) : null,
        individual_limit: Number(formData.individual_limit) || 1,
        editable_responses: Boolean(formData.editable_responses)
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

<div class="min-h-screen bg-background relative pt-24">
  <div class="absolute top-24 right-8 z-10">
    <FormConfig bind:formData bind:dialogOpen on:panelchange={handlePanelChange} on:dialogopenchange={handleDialogOpenChange} />
  </div>
  <div 
    id="MainContent" 
    class="mx-auto transition-transform duration-300 ease-in-out w-full max-w-4xl px-4" 
    style="transform: translateX({getDisplacement(windowWidth, isPanelOpen)});"
  >
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
      <div class="space-y-0 w-full mb-10">
        <FormTitleDescription bind:title={formData.title} bind:description={formData.description} />
        <section class="space-y-6 w-full">
          {#if questions.length > 0}
            <QuestionsList
              {questions}
              {questionTypeLabels}
              {moveQuestionUp}
              {moveQuestionDown}
              {removeQuestion}
            />
          {/if}
          <AddQuestionPanel
            {questionTypeButtons}
            {addQuestion}
            questionsLength={questions.length}
          />
          {#if questions.length > 0}
            <FormSaveFooter
              questionsLength={questions.length}
              {isSaving}
              {saveForm}
            />
          {/if}
        </section>
      </div>
    {/if}
  </div>
</div>