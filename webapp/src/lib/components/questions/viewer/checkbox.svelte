<script lang="ts">
  import { Label } from '$lib/components/ui/label';
  import Checkbox from '../../ui/checkbox/checkbox.svelte';
  import { onMount, createEventDispatcher } from 'svelte';
  export let question: {
    id: string;
    title: string;
    required: boolean;
    options: { id: string; value: string; label: string }[];
  };
  export let value: string = '[]';
  export let disabled: boolean = false;
  let selected: Set<string> = new Set();
  const dispatch = createEventDispatcher();
  onMount(() => {
    try {
      const arr = JSON.parse(value);
      if (Array.isArray(arr)) selected = new Set(arr);
    } catch {}
  });
  function handleCheckboxChange(optId: string, event: Event) {
    const input = event.target as HTMLInputElement;
    if (input.checked) selected.add(optId);
    else selected.delete(optId);
    value = JSON.stringify(Array.from(selected));
    dispatch('input', value);
  }
</script>
{#if question.options.length > 0}
  <div class="space-y-4">
    {#each question.options as opt (opt.id)}
      <div class="flex items-center gap-3 p-2 py-0 rounded-md">
        <Checkbox
          id={`checkbox-${opt.id}`}
          checked={selected.has(opt.id)}
          disabled={disabled}
          onchange={(e) => handleCheckboxChange(opt.id, e)}
          class="data-[state=checked]:bg-primary border-input"
        />
        <Label for={`checkbox-${opt.id}`} class="text-sm font-normal cursor-pointer flex-1 select-none">
          {opt.label}
        </Label>
      </div>
    {/each}
  </div>
{:else}
  <div class="text-sm text-muted-foreground bg-muted/50 p-4 rounded-md border-dashed border">
    No options available for this question.
  </div>
{/if}
