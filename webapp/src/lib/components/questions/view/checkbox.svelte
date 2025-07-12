<script lang="ts">
  import { Checkbox } from '$lib/components/ui/checkbox';
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
  let selectedOptions = $state<string[]>([]);

  $effect(() => {
    if (value) {
      try {
        selectedOptions = JSON.parse(value);
      } catch {
        selectedOptions = [];
      }
    } else {
      selectedOptions = [];
    }
  });

  $effect(() => {
    value = JSON.stringify(selectedOptions);
  });

  function handleOptionChange(optionUid: string, checked: boolean) {
    if (checked) {
      selectedOptions = [...selectedOptions, optionUid];
    } else {
      selectedOptions = selectedOptions.filter(id => id !== optionUid);
    }
  }
</script>

{#if question.options.length > 0}
  <div class="space-y-3">
    {#each question.options as option (option.uid)}
      <div class="flex items-center space-x-3">
        <Checkbox
          id="option-{option.uid}"
          checked={selectedOptions.includes(option.uid)}
          onCheckedChange={(checked) => handleOptionChange(option.uid, !!checked)}
        />
        <Label 
          for="option-{option.uid}" 
          class="text-sm font-normal cursor-pointer flex-1"
        >
          {option.text}
        </Label>
      </div>
    {/each}
  </div>
{:else}
  <div class="text-sm text-muted-foreground bg-muted/50 p-4 rounded-md border-dashed border">
    No options available for this question.
  </div>
{/if}
