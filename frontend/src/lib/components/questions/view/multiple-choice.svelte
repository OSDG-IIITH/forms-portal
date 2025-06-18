<script lang="ts">
  import * as RadioGroup from '$lib/components/ui/radio-group';
  import { Label } from '$lib/components/ui/label';

  interface Option {
    uid: string;
    text: string;
    type: 'text';
  }

  interface Props {
    question: {
      uid: string;
      text: string;
      required: boolean;
      options: Option[];
    };
    value?: string;
  }

  let { question, value = $bindable('') }: Props = $props();
</script>

{#if question.options.length > 0}
  <RadioGroup.Root bind:value class="space-y-3">
    {#each question.options as option (option.uid)}
      <div class="flex items-center space-x-3">
        <RadioGroup.Item value={option.uid} id="option-{option.uid}" />
        <Label 
          for="option-{option.uid}" 
          class="text-sm font-normal cursor-pointer flex-1"
        >
          {option.text}
        </Label>
      </div>
    {/each}
  </RadioGroup.Root>
{:else}
  <div class="text-sm text-muted-foreground bg-muted/50 p-4 rounded-md border-dashed border">
    No options available for this question.
  </div>
{/if}
