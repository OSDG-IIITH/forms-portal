<script lang="ts">
  import * as Select from "$lib/components/ui/select/index.js";

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

  const selectOptions = $derived(
    question.options.map(option => ({
      value: option.uid,
      label: option.text
    }))
  );

  const triggerContent = $derived(
    selectOptions.find((option) => option.value === value)?.label ?? "Select an option"
  );
</script>

{#if question.options.length > 0}
  <Select.Root type="single" bind:value>
    <Select.Trigger class="w-full">
      {triggerContent}
    </Select.Trigger>
    <Select.Content>
      <Select.Group>
        {#each selectOptions as option (option.value)}
          <Select.Item
            value={option.value}
            label={option.label}
          >
            {option.label}
          </Select.Item>
        {/each}
      </Select.Group>
    </Select.Content>
  </Select.Root>
{:else}
  <div class="text-sm text-muted-foreground bg-muted/50 p-4 rounded-md border-dashed border">
    No options available for this question.
  </div>
{/if}
