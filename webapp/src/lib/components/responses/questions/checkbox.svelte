<script lang="ts">
  import { Checkbox } from '$lib/components/ui/checkbox';
  import { Label } from '$lib/components/ui/label';

  export let question: {
    id: string;
    title: string;
    required: boolean;
    options: { id: string; value: string; label: string }[];
  };
  export let value: string = '[]';
  export let disabled: boolean = true;

  let selectedValues: string[] = [];

  $: {
    try {
      selectedValues = JSON.parse(value);
      if (!Array.isArray(selectedValues)) {
        selectedValues = [];
      }
    } catch {
      selectedValues = [];
    }
  }
</script>

{#if question.options.length > 0}
  <div class="space-y-3">
    {#each question.options as option (option.id)}
      <div class="flex items-center space-x-3">
        <Checkbox
          checked={selectedValues.includes(option.value)}
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