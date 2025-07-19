<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { X } from '@lucide/svelte';
  export let label: string;
  export let active: boolean = false;
  export let icon: any;
  export let clickable: boolean = true;
  export let onRemove: (() => void) | undefined;
  const dispatch = createEventDispatcher();
  function handleClick() { if (clickable) dispatch('click'); }
  function handleKeydown(e: KeyboardEvent) {
    if (!clickable) return;
    if (e.key === 'Enter' || e.key === ' ') {
      e.preventDefault();
      handleClick();
    }
  }
</script>

<div
  class="group relative flex items-center gap-2 px-3 py-1.5 rounded-md border text-sm font-medium shadow-xs select-none transition-all duration-200
    {active && clickable
      ? 'bg-primary text-primary-foreground border border-primary cursor-pointer'
      : clickable
        ? 'bg-muted/60 hover:bg-muted/90 text-accent-foreground/85 border cursor-pointer'
        : 'bg-muted/60 hover:bg-muted/90 text-accent-foreground/85 border cursor-default'}"
  role={clickable ? 'button' : undefined}
  {...(clickable ? { tabindex: 0, onkeydown: handleKeydown } : {})}
  on:click={clickable ? handleClick : undefined}
>
  {#if icon}
    <svelte:component this={icon} class="size-4 shrink-0"></svelte:component>
  {/if}
  <span>{label}</span>
  {#if onRemove}
    <button
      type="button"
      class="absolute -top-1 -right-1 rounded-full p-0.5 bg-background border border-muted text-muted-foreground transition-colors
        opacity-0 pointer-events-none
        group-hover:opacity-100 group-hover:pointer-events-auto
        hover:bg-destructive/60 hover:border-destructive hover:text-white"
      on:click|stopPropagation={onRemove}
      aria-label="Remove validation"
      tabindex="-1"
    >
      <X class="size-3" />
    </button>
  {/if}
</div>
