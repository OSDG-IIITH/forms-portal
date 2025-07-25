<script lang="ts">
  import { Checkbox } from '$lib/components/ui/checkbox';
  import { Label } from '$lib/components/ui/label';

  let {
    question,
    value = $bindable('[]'),
    disabled = false
  }: {
    question: {
      id: string;
      title: string;
      required: boolean;
      options: { id: string; value: string; label: string }[];
    };
    value: string;
    disabled?: boolean;
  } = $props();

  let selectedIds = $state<string[]>([]);
  let isUpdatingFromValue = false;

  $effect(() => {
    if (!isUpdatingFromValue) {
      try {
        const parsed = JSON.parse(value);
        selectedIds = parsed;
      } catch {
        selectedIds = [];
      }
    }
  });

  function handleCheckboxChange(optionId: string, checked: boolean | 'indeterminate') {
    const isChecked = checked === true;
    if (isChecked) {
      if (!selectedIds.includes(optionId)) {
        selectedIds = [...selectedIds, optionId];
      }
    } else {
      selectedIds = selectedIds.filter(id => id !== optionId);
    }
    
    isUpdatingFromValue = true;
    value = JSON.stringify(selectedIds);
    isUpdatingFromValue = false;
  }
</script>

{#if question.options.length > 0}
  <div class="space-y-3">
    {#each question.options as option (option.id)}
      <div class="flex items-center space-x-3">
        <Checkbox
          checked={selectedIds.includes(option.id)}
          onCheckedChange={(checked) => handleCheckboxChange(option.id, checked)}
          {disabled}
          id={`checkbox-${option.id}`}
        />
        <Label 
          for={`checkbox-${option.id}`}
          class="text-sm font-normal cursor-pointer flex-1"
        >
          {option.label}
        </Label>
      </div>
    {/each}
  </div>
{:else}
  <div class="text-sm text-muted-foreground bg-muted/50 p-4 rounded-md border-dashed border">
    No options available for this question.
  </div>
{/if}