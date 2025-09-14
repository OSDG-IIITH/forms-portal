<script lang="ts">
  import { base } from '$app/paths';
  import FormTitleDescription from './FormTitleDescription.svelte';
  import SidePanel from './SidePanel.svelte';
  import QuestionsList from './QuestionsList.svelte';
  import AddQuestionPanel from './AddQuestionPanel.svelte';
  import FormSaveFooter from './FormActionFooter.svelte';
  import * as kdljs from 'kdljs';
  import { toast } from 'svelte-sonner';
  import { IconFileText, IconListCheck, IconSquareCheck, IconUpload, IconCalendar, IconHeading } from '@tabler/icons-svelte';
  import { createEventDispatcher, setContext } from 'svelte';
  import type { Time } from '@internationalized/date';
  import { createFormStore, type FormStore, type QuestionType } from './form-store.svelte';
  import { generateFormKdl } from './FormEditor';

  const { form, mode = 'edit' }: { form: any; mode?: 'create' | 'edit' } = $props();

  const store: FormStore = createFormStore(form);
  setContext('form-store', store);

  let isSaving = $state(false);

  const questionTypeLabels: Record<QuestionType, string> = {
    input: 'Input',
    textarea: 'Textarea',
    radio: 'Multiple Choice',
    checkbox: 'Checkbox',
    file: 'File Upload',
    select: 'Select',
    date: 'Date Picker',
    'section-header': 'Section Header'
  };

  const questionTypeButtons = [
    { type: 'input' as QuestionType, icon: IconFileText, label: 'Text' },
    { type: 'textarea' as QuestionType, icon: IconFileText, label: 'Long Text' },
    { type: 'radio' as QuestionType, icon: IconListCheck, label: 'Multiple Choice' },
    { type: 'checkbox' as QuestionType, icon: IconSquareCheck, label: 'Checkbox' },
    { type: 'file' as QuestionType, icon: IconUpload, label: 'File Upload' },
    { type: 'select' as QuestionType, icon: IconListCheck, label: 'Select' },
    { type: 'date' as QuestionType, icon: IconCalendar, label: 'Date Picker' },
    { type: 'section-header' as QuestionType, icon: IconHeading, label: 'Section Header' }
  ];

  const dispatch = createEventDispatcher<{
    create: { formData: any; questions: any[]; kdl: string };
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

  async function saveForm() {
    store.validateQuestions();
    
    if (!store.isFormValid) {
      toast.error("Please fix validation errors before saving");
      return;
    }
    
    isSaving = true;
    try {
      const kdl = generateFormKdl(store.questions);

      kdljs.parse(kdl);

      if (mode === 'create') {
        dispatch('create', { formData: store.formData, questions: store.questions, kdl });
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
        title: store.formData.title || 'Untitled Form',
        description: store.formData.description || null,
        structure: kdl,
        opens: formatDateTime(store.formData.opens, store.formData.opensTime),
        closes: formatDateTime(store.formData.closes, store.formData.closesTime),
        anonymous: Boolean(store.formData.anonymous),
        max_responses:
          store.formData.max_responses !== null && store.formData.max_responses !== undefined
            ? Number(store.formData.max_responses)
            : null,
        individual_limit: Number(store.formData.individual_limit) || 1,
        editable_responses: Boolean(store.formData.editable_responses)
      };

      const res = await fetch(`${base}/api/forms/${form.id}`, {
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
    <SidePanel
      bind:dialogOpen
      on:panelchange={handlePanelChange}
      on:dialogopenchange={handleDialogOpenChange}
    />
  </div>
  <div
    id="MainContent"
    class="mx-auto transition-transform duration-300 ease-in-out w-full max-w-4xl px-4"
    style="transform: translateX({getDisplacement(windowWidth, isPanelOpen)});"
  >
    {#if false}
      <!-- ... -->
    {:else}
      <div class="space-y-0 w-full mb-10">
        <FormTitleDescription
        />
        <section class="space-y-6 w-full">
          {#if store.questions.length > 0}
            <QuestionsList {questionTypeLabels} />
          {/if}
          <AddQuestionPanel
            {questionTypeButtons}
            addQuestion={store.addQuestion}
            questionsLength={store.questions.length}
          />
          {#if store.questions.length > 0}
            <FormSaveFooter
              questionsLength={store.questions.length}
              {isSaving}
              isFormValid={store.isFormValid}
              {saveForm}
            />
          {/if}
        </section>
      </div>
    {/if}
  </div>
</div>
