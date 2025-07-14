<script lang="ts">
  import * as RadioGroup from '$lib/components/ui/radio-group';
  import { Label } from '$lib/components/ui/label';

  export let question: {
    id: string;
    title: string;
    required: boolean;
    options: { id: string; value: string; label: string }[];
  };
  export let value: string = '';
  export let disabled: boolean = false;
</script>

{#if question.options.length > 0}
  <RadioGroup.Root bind:value disabled={disabled}>
    {#each question.options as option (option.id)}
      <div class="flex items-center space-x-3">
        <RadioGroup.Item value={option.id} id={`option-${option.id}`} />
        <Label 
          for={`option-${option.id}`}
          class="text-sm font-normal cursor-pointer flex-1"
        >
          {option.label}
        </Label>
      </div>
    {/each}
  </RadioGroup.Root>
{:else}
  <div class="text-sm text-muted-foreground bg-muted/50 p-4 rounded-md border-dashed border">
    No options available for this question.
  </div>
{/if}
