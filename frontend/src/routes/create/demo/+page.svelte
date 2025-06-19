<script lang="ts">
  import { Card, CardContent } from '$lib/components/ui/card';
  import { Button } from '$lib/components/ui/button';
  import { Input } from '$lib/components/ui/input';
  import { Textarea } from '$lib/components/ui/textarea';
  import { Label } from '$lib/components/ui/label';
  import TextQuestionEditor from '$lib/components/questions/editor/text.svelte';
  import MultipleChoiceEditor from '$lib/components/questions/editor/multiple-choice.svelte';
  import CheckboxEditor from '$lib/components/questions/editor/checkbox.svelte';
  import FileUploadEditor from '$lib/components/questions/editor/file-upload.svelte';
  import SelectEditor from '$lib/components/questions/editor/select.svelte';
  import DatePickerEditor from '$lib/components/questions/editor/date-picker.svelte';
  import ThemeSwitcher from '$lib/components/theme-switcher.svelte';
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

  type QuestionType = 'text' | 'multiple_choice' | 'checkbox' | 'file_upload' | 'select' | 'date_picker';

  type Option = {
    uid: string;
    text: string;
    type: 'text';
  };

  type Question = {
    uid: string;
    type: QuestionType;
    text: string;
    required: boolean;
    display_order: number;
    options: Option[];
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
    text: 'Text',
    multiple_choice: 'Multiple Choice',
    checkbox: 'Checkbox',
    file_upload: 'File Upload',
    select: 'Select',
    date_picker: 'Date Picker'
  };

  const questionTypeButtons = [
    { type: 'text' as QuestionType, icon: IconFileText, label: 'Text' },
    { type: 'multiple_choice' as QuestionType, icon: IconListCheck, label: 'Multiple Choice' },
    { type: 'checkbox' as QuestionType, icon: IconSquareCheck, label: 'Checkbox' },
    { type: 'file_upload' as QuestionType, icon: IconUpload, label: 'File Upload' },
    { type: 'select' as QuestionType, icon: IconListCheck, label: 'Select' },
    { type: 'date_picker' as QuestionType, icon: IconCalendar, label: 'Date Picker' }
  ];

  function addQuestion(type: QuestionType): void {
    const newQuestion: Question = {
      uid: crypto.randomUUID(),
      type,
      text: '',
      required: false,
      display_order: questions.length,
      options: []
    };
    questions = [...questions, newQuestion];
  }

  function removeQuestion(uid: string): void {
    questions = questions
      .filter(q => q.uid !== uid)
      .map((q, index) => ({ ...q, display_order: index }));
  }

  function moveQuestionUp(index: number): void {
    if (index > 0) {
      const newQuestions = [...questions];
      [newQuestions[index - 1], newQuestions[index]] = [newQuestions[index], newQuestions[index - 1]];
      questions = newQuestions.map((q, i) => ({ ...q, display_order: i }));
    }
  }

  function moveQuestionDown(index: number): void {
    if (index < questions.length - 1) {
      const newQuestions = [...questions];
      [newQuestions[index], newQuestions[index + 1]] = [newQuestions[index + 1], newQuestions[index]];
      questions = newQuestions.map((q, i) => ({ ...q, display_order: i }));
    }
  }

  function saveFormAsJson() {
    const data = {
      form: {
        title: form.title,
        description: form.description,
        visibility: form.visibility
      },
      questions: questions.map(q => ({
        uid: q.uid,
        type: q.type,
        text: q.text,
        required: q.required,
        display_order: q.display_order,
        options: q.options?.map(o => ({
          uid: o.uid,
          text: o.text,
          type: o.type
        })) ?? []
      }))
    };

    const json = JSON.stringify(data, null, 2);
    const blob = new Blob([json], { type: "application/json" });
    const url = URL.createObjectURL(blob);

    const a = document.createElement("a");
    a.href = url;
    a.download = `${form.title || "form"}.json`;
    document.body.appendChild(a);
    a.click();
    setTimeout(() => {
      document.body.removeChild(a);
      URL.revokeObjectURL(url);
    }, 0);
  }
</script>

<div class="min-h-screen bg-background">
  <!-- header bar -->
  <header class="border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60 fixed top-0 left-0 right-0 z-50">
    <div class="container mx-auto px-6 py-4">
      <div class="flex items-center justify-between">
        <h1 class="text-3xl tracking-tight">
          forms <span class="font-semibold text-accent-foreground">iiit</span>
        </h1>
        <ThemeSwitcher />
      </div>
    </div>
  </header>

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
            {#each questions as question, i (question.uid)}
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
                        size="sm"
                        class="h-8 w-8 p-0"
                        onclick={() => moveQuestionUp(i)}
                        disabled={i === 0}
                        aria-label="Move question up"
                      >
                        <IconChevronUp class="h-4 w-4" />
                      </Button>
                      <Button
                        variant="ghost"
                        size="sm"
                        class="h-8 w-8 p-0"
                        onclick={() => moveQuestionDown(i)}
                        disabled={i === questions.length - 1}
                        aria-label="Move question down"
                      >
                        <IconChevronDown class="h-4 w-4" />
                      </Button>
                      <Button
                        variant="ghost"
                        size="sm"
                        class="h-8 w-8 p-0 ml-1 text-muted-foreground hover:text-destructive hover:bg-destructive/10 transition-colors"
                        onclick={() => removeQuestion(question.uid)}
                        aria-label="Delete question"
                      >
                        <IconTrash class="h-4 w-4" />
                      </Button>
                    </div>
                  </div>

                  <!-- editor container -->
                  <div class="pl-4">
                    {#if question.type === 'text'}
                      <TextQuestionEditor bind:question={questions[i]} />
                    {:else if question.type === 'multiple_choice'}
                      <MultipleChoiceEditor bind:question={questions[i]} />
                    {:else if question.type === 'checkbox'}
                      <CheckboxEditor bind:question={questions[i]} />
                    {:else if question.type === 'file_upload'}
                      <FileUploadEditor bind:question={questions[i]} />
                    {:else if question.type === 'select'}
                      <SelectEditor bind:question={questions[i]} />
                    {:else if question.type === 'date_picker'}
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
              <div class="flex flex-wrap gap-2">
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
                <Button size="sm" onclick={saveFormAsJson}>Save Form</Button>
              </div>
            </div>
          </CardContent>
        </Card>
      {/if}
    </div>
  </main>
</div>