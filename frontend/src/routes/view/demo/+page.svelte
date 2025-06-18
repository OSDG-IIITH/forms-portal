<script lang="ts">
  import { onMount } from 'svelte';
  import { Card, CardContent } from '$lib/components/ui/card';
  import { Button } from '$lib/components/ui/button';
  import TextQuestionView from '$lib/components/questions/view/text.svelte';
  import MultipleChoiceView from '$lib/components/questions/view/multiple-choice.svelte';
  import CheckboxView from '$lib/components/questions/view/checkbox.svelte';
  import FileUploadView from '$lib/components/questions/view/file-upload.svelte';
  import DatePickerView from '$lib/components/questions/view/date-picker.svelte';
  import SelectView from '$lib/components/questions/view/select.svelte';
  import ThemeToggle from '$lib/components/theme-toggle.svelte';

  type QuestionType = 'text' | 'multiple_choice' | 'checkbox' | 'file_upload' | 'select' | 'date_picker';

  interface Option {
    uid: string;
    text: string;
    type: 'text';
  }

  interface Question {
    uid: string;
    type: QuestionType;
    text: string;
    required: boolean;
    display_order: number;
    options: Option[];
  }

  interface FormConfig {
    title: string;
    description: string;
    visibility: string;
  }

  interface FormData {
    form: FormConfig;
    questions: Question[];
  }

  let formConfig = $state<FormConfig>({
    title: '',
    description: '',
    visibility: 'public'
  });
  
  let questions = $state<Question[]>([]);
  let responses = $state<Record<string, string>>({});
  let isSubmitting = $state(false);
  let isSubmitted = $state(false);
  let isLoading = $state(true);
  let error = $state<string | null>(null);

  async function loadFormData() {
    try {
      isLoading = true;
      error = null;
      
      const { default: formData } = await import('$lib/form.json');
      const typedFormData = formData as FormData;
      
      formConfig = {
        title: typedFormData.form.title || 'Untitled Form',
        description: typedFormData.form.description || '',
        visibility: typedFormData.form.visibility || 'public'
      };
      
      questions = (typedFormData.questions || [])
        .map(question => ({ 
          ...question, 
          type: question.type as QuestionType,
          options: question.options || []
        }))
        .sort((a, b) => a.display_order - b.display_order);
    
      const initialResponses: Record<string, string> = {};
      for (const question of questions) {
        const isArrayType = question.type === 'checkbox' || question.type === 'file_upload';
        initialResponses[question.uid] = isArrayType ? '[]' : '';
      }
      responses = initialResponses;
      
    } catch (err) {
      console.error('Failed to load form data:', err);
      error = 'Failed to load form data. Please try again.';
    } finally {
      isLoading = false;
    }
  }

  async function handleSubmit(event: Event) {
    event.preventDefault();
    
    if (isSubmitting) return;
    
    isSubmitting = true;
    
    try {
      const requiredQuestions = questions.filter(q => q.required);
      const invalidFields: Question[] = [];
      
      for (const question of requiredQuestions) {
        const responseValue = responses[question.uid];
        const isValid = validateQuestionResponse(question, responseValue);
        
        if (!isValid) {
          invalidFields.push(question);
        }
      }
    
      if (invalidFields.length > 0) {
        const fieldNames = invalidFields.map(q => q.text).join(', ');
        alert(`Please fill in the following required fields: ${fieldNames}`);
        return;
      }

      await new Promise(resolve => setTimeout(resolve, 1500));
    
      isSubmitted = true;
      window.scrollTo({ top: 0, behavior: 'smooth' });
    
    } catch (err) {
      console.error('Form submission failed:', err);
      alert('Form submission failed. Please try again.');
    } finally {
      isSubmitting = false;
    }
  }

  function validateQuestionResponse(question: Question, responseValue: string): boolean {
    if (question.type === 'checkbox' || question.type === 'file_upload') {
      try {
        const items = JSON.parse(responseValue || '[]');
        return Array.isArray(items) && items.length > 0;
      } catch {
        return false;
      }
    }
  
    return Boolean(responseValue?.trim());
  }

  function handleReset() {
    const resetResponses: Record<string, string> = {};
  
    for (const question of questions) {
      const isArrayType = question.type === 'checkbox' || question.type === 'file_upload';
      resetResponses[question.uid] = isArrayType ? '[]' : '';
    }
  
    responses = resetResponses;
    isSubmitted = false;
  }

  function getDisplayValue(question: Question, responseValue: string): string {
    if (!responseValue) return 'No response provided';
    
    switch (question.type) {
      case 'multiple_choice':
      case 'select': {
        const selectedOption = question.options.find(opt => opt.uid === responseValue);
        return selectedOption ? selectedOption.text : 'Invalid selection';
      }
      
      case 'checkbox': {
        try {
          const selectedOptionIds = JSON.parse(responseValue);
          if (selectedOptionIds.length === 0) return 'No options selected';
          
          const selectedOptions = question.options.filter(opt => 
            selectedOptionIds.includes(opt.uid)
          );
          return selectedOptions.map(opt => opt.text).join(', ');
        } catch {
          return 'Invalid selection';
        }
      }
      
      case 'file_upload': {
        return responseValue || 'No files selected';
      }
      
      case 'date_picker': {
        if (!responseValue.trim()) return 'No date selected';
        try {
          const date = new Date(responseValue);
          return date.toLocaleDateString('en-US', { 
            year: 'numeric', 
            month: 'long', 
            day: 'numeric' 
          });
        } catch {
          return responseValue;
        }
      }
      
      default:
        return responseValue;
    }
  }

  onMount(() => {
    loadFormData();
  });

  const requiredFieldsCount = $derived(questions.filter(q => q.required).length);
  const supportedQuestions = $derived(
    questions.filter(q => 
      ['text', 'multiple_choice', 'checkbox', 'file_upload', 'date_picker', 'select'].includes(q.type)
    )
  );
</script>

<svelte:head>
  <title>{formConfig.title || 'Form'} - Form Viewer</title>
  <meta name="description" content={formConfig.description || 'Form viewer'} />
</svelte:head>

<div class="min-h-screen bg-background">
  <header class="border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60 fixed top-0 left-0 right-0 z-50">
    <div class="container mx-auto px-6 py-4">
      <div class="flex items-center justify-between">
        <h1 class="text-3xl tracking-tight">
          forms <span class="font-semibold">iiit</span>
        </h1>
        <ThemeToggle />
      </div>
    </div>
  </header>

  <main class="container mx-auto max-w-4xl px-6 py-8 pt-24">
    {#if isLoading}
      <div class="flex items-center justify-center py-12">
        <div class="text-center space-y-4">
          <div class="w-8 h-8 border-4 border-primary border-t-transparent rounded-full animate-spin mx-auto"></div>
          <p class="text-muted-foreground">Loading form...</p>
        </div>
      </div>
    {:else if error}
      <Card>
        <CardContent class="text-center py-12">
          <div class="w-12 h-12 bg-destructive/10 rounded-full flex items-center justify-center mx-auto mb-4">
            <svg class="w-6 h-6 text-destructive" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z"></path>
            </svg>
          </div>
          <h3 class="text-lg font-semibold mb-2">Error Loading Form</h3>
          <p class="text-muted-foreground mb-6">{error}</p>
          <Button onclick={loadFormData}>Try Again</Button>
        </CardContent>
      </Card>
    {:else if !isSubmitted}
      <div class="space-y-8">
        <Card>
          <CardContent class="space-y-6">
            <div class="space-y-3">
              <h1 class="text-3xl font-bold tracking-tight">{formConfig.title}</h1>
              {#if formConfig.description && formConfig.description !== 'Add a description'}
                <p class="text-lg text-muted-foreground leading-relaxed">
                  {formConfig.description}
                </p>
              {/if}
            </div>
          </CardContent>
        </Card>

        <form onsubmit={handleSubmit} class="space-y-4">
          {#each questions as question (question.uid)}
            <Card>
              <CardContent class="space-y-6">
                <div class="flex items-start gap-3">
                  <div class="text-base font-medium leading-relaxed">
                    {question.text}
                    {#if question.required}
                      <span class="text-destructive ml-1" title="Required field">*</span>
                    {/if}
                  </div>
                </div>

                <div>
                  {#if question.type === 'text'}
                    <TextQuestionView 
                      {question} 
                      bind:value={responses[question.uid]}
                    />
                  {:else if question.type === 'multiple_choice'}
                    <MultipleChoiceView 
                      {question} 
                      bind:value={responses[question.uid]}
                    />
                  {:else if question.type === 'checkbox'}
                    <CheckboxView 
                      {question} 
                      bind:value={responses[question.uid]}
                    />
                  {:else if question.type === 'file_upload'}
                    <FileUploadView 
                      {question} 
                      bind:value={responses[question.uid]}
                    />
                  {:else if question.type === 'date_picker'}
                    <DatePickerView 
                      {question} 
                      bind:value={responses[question.uid]}
                    />
                  {:else if question.type === 'select'}
                    <SelectView 
                      {question} 
                      bind:value={responses[question.uid]}
                    />
                  {:else}
                    <div class="p-4 bg-muted/30 border border-dashed border-muted-foreground/30 rounded-lg">
                      <p class="text-sm text-muted-foreground text-center">
                        Question type <span class="font-mono font-medium">"{question.type}"</span> is not yet supported in this demo.
                      </p>
                    </div>
                  {/if}
                </div>
              </CardContent>
            </Card>
          {/each}

          <Card>
            <CardContent>
              <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
                <div class="text-sm text-muted-foreground">
                  {#if requiredFieldsCount > 0}
                    <span class="font-medium">{requiredFieldsCount}</span> required field{requiredFieldsCount !== 1 ? 's' : ''}
                  {:else}
                    All fields are optional
                  {/if}
                </div>
                <div class="flex gap-3 w-full sm:w-auto">
                  <Button 
                    type="button" 
                    variant="outline" 
                    size="sm"
                    onclick={handleReset}
                    disabled={isSubmitting}
                    class="flex-1 sm:flex-none"
                  >
                    Reset Form
                  </Button>
                  <Button 
                    type="submit" 
                    size="sm"
                    disabled={isSubmitting}
                    class="flex-1 sm:flex-none min-w-[140px]"
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
        </form>
      </div>
    {:else}
      <div class="space-y-8">
        <Card>
          <CardContent class="text-center py-16">
            <div class="w-16 h-16 bg-green-100 dark:bg-green-900/20 rounded-full flex items-center justify-center mx-auto mb-6">
              <svg class="w-8 h-8 text-green-600 dark:text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
              </svg>
            </div>
            <h2 class="text-2xl font-bold mb-3">Thank you for your submission!</h2>
            <p class="text-muted-foreground mb-8 max-w-md mx-auto">
              Your responses have been recorded successfully. We appreciate you taking the time to complete this form.
            </p>
            <Button onclick={handleReset} size="lg">
              Submit Another Response
            </Button>
          </CardContent>
        </Card>

        {#if supportedQuestions.length > 0}
          <div class="space-y-6">
            <h3 class="text-xl font-semibold">Your Responses</h3>
            {#each supportedQuestions as question}
              <Card>
                <CardContent class="space-y-6">
                  <div class="flex items-start gap-3">
                    <div class="text-base font-medium leading-relaxed">
                      {question.text}
                    </div>
                  </div>
                  <div>
                    <div class="bg-muted/50 p-4 rounded-lg border">
                      <p class="text-sm whitespace-pre-wrap">
                        {getDisplayValue(question, responses[question.uid])}
                      </p>
                    </div>
                  </div>
                </CardContent>
              </Card>
            {/each}
          </div>
        {/if}
      </div>
    {/if}
  </main>
</div>
