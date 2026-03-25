<script lang="ts">
  import { Input } from '$lib/components/ui/input';
  import { Label } from '$lib/components/ui/label';
  import { Checkbox } from '$lib/components/ui/checkbox';

  let {
    questionId,
    title,
    error,
    required,
    onTitleChange,
    onRequiredChange = () => {},
    titleLabel = 'Question Text',
    titlePlaceholder = 'Enter your question here',
    showRequired = true,
    requiredLabel = 'Required',
    requiredClass = 'text-sm',
    requiredContainerClass = 'flex items-center gap-2'
  }: {
    questionId: string;
    title: string;
    error?: string;
    required: boolean;
    onTitleChange: (value: string) => void;
    onRequiredChange?: (value: boolean) => void;
    titleLabel?: string;
    titlePlaceholder?: string;
    showRequired?: boolean;
    requiredLabel?: string;
    requiredClass?: string;
    requiredContainerClass?: string;
  } = $props();
</script>

<div class="space-y-2">
  <Label for="question-{questionId}">{titleLabel}</Label>
  <Input
    id="question-{questionId}"
    placeholder={titlePlaceholder}
    value={title}
    oninput={(e) => onTitleChange(e.currentTarget.value)}
  />
  {#if error}
    <p class="text-destructive text-sm">{error}</p>
  {/if}
</div>

{#if showRequired}
  <div class={requiredContainerClass}>
    <Checkbox
      id="required-{questionId}"
      checked={required}
      onCheckedChange={(v) => onRequiredChange(!!v)}
    />
    <Label for="required-{questionId}" class={requiredClass}>{requiredLabel}</Label>
  </div>
{/if}