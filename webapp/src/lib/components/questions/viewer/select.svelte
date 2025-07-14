<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { Select } from '$lib/components/ui/select/index';
  import { Trigger as SelectTrigger, Content as SelectContent, Item as SelectItem } from '$lib/components/ui/select/index';
  export let question: {
    id: string;
    title: string;
    required: boolean;
    options?: { id: string; value: string; label: string }[];
    placeholder?: string;
  };
  export let value: string = '';
  export let disabled: boolean = false;
  export let placeholder: string = 'Select an option';
  const dispatch = createEventDispatcher();
</script>

<Select type="single" allowDeselect bind:value {disabled}>
  <SelectTrigger class="w-full">
    {#if value}
      {question.options?.find(opt => opt.value === value)?.label || placeholder}
    {:else}
      <span class="text-muted-foreground">{placeholder}</span>
    {/if}
  </SelectTrigger>
  <SelectContent>
    {#each question.options ?? [] as opt}
      <SelectItem value={opt.value} label={opt.label}>
        {opt.label}
      </SelectItem>
    {/each}
  </SelectContent>
</Select>
