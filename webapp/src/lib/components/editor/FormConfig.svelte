<script lang="ts">
  import { Input } from '$lib/components/ui/input';
  import { Label } from '$lib/components/ui/label';
  import { DatePicker } from '$lib/components/ui/date-picker';
  import { Checkbox } from '$lib/components/ui/checkbox';
  import * as Panel from "$lib/components/ui/panel";
  import * as Dialog from "$lib/components/ui/dialog";
  import { buttonVariants } from "$lib/components/ui/button";
  import { TimeField } from 'bits-ui';
  import { IconSettings } from '@tabler/icons-svelte';
  import { createEventDispatcher } from 'svelte';

  let { formData = $bindable(), dialogOpen = $bindable(false) } = $props();
  
  const dispatch = createEventDispatcher<{
    panelchange: { open: boolean };
    dialogopenchange: boolean;
  }>();

  let windowWidth = $state(0);

  const useDialog = $derived(windowWidth > 0 && windowWidth < 1380);

  function handlePanelChange(event: CustomEvent<{ open: boolean }>) {
    dispatch('panelchange', event.detail);
  }

  function handleDialogOpenChange(open: boolean) {
    dialogOpen = open;
    dispatch('dialogopenchange', open);
    if (useDialog && !open) {
      dispatch('panelchange', { open: false });
    }
  }

  $effect(() => {
    if (typeof window !== 'undefined') {
      windowWidth = window.innerWidth;
      
      const handleResize = () => {
        windowWidth = window.innerWidth;
        if (windowWidth >= 1024 && dialogOpen) {
          handleDialogOpenChange(false);
        }
      };
      
      window.addEventListener('resize', handleResize);
      return () => window.removeEventListener('resize', handleResize);
    }
  });
</script>

{#snippet settingsContent()}
  <div class="grid flex-1 auto-rows-min gap-6 px-8 py-2">
    <div class="grid grid-cols-1 gap-6 pt-2">
      <div class="space-y-2">
        <Label for="form-opens">Opens</Label>
        <div class="flex gap-2">
          <DatePicker 
            bind:value={formData.opens} 
            placeholder="Select open date"
            class="flex-1" 
          />
          <div class="w-22">
            <TimeField.Root bind:value={formData.opensTime} hourCycle={24}>
              <TimeField.Input class="flex h-9 w-full rounded-md border border-input bg-muted/50 px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50">
                {#snippet children({ segments })}
                  {#each segments as { part, value }, i (part + value + i)}
                    <div class="inline-block select-none">
                      {#if part === "literal"}
                        <TimeField.Segment {part} class="text-muted-foreground p-1">
                          {value}
                        </TimeField.Segment>
                      {:else}
                        <TimeField.Segment
                          {part}
                          class="rounded px-1 py-1 hover:bg-muted focus:bg-muted focus:text-foreground aria-[valuetext=Empty]:text-muted-foreground focus-visible:ring-0 focus-visible:ring-offset-0"
                        >
                          {value}
                        </TimeField.Segment>
                      {/if}
                    </div>
                  {/each}
                {/snippet}
              </TimeField.Input>
            </TimeField.Root>
          </div>
        </div>
      </div>
      <div class="space-y-2">
        <Label for="form-closes">Closes</Label>
        <div class="flex gap-2">
          <DatePicker 
            bind:value={formData.closes} 
            placeholder="Select close date"
            class="flex-1" 
          />
          <div class="w-22">
            <TimeField.Root bind:value={formData.closesTime} hourCycle={24}>
              <TimeField.Input class="flex h-9 w-full rounded-md border border-input bg-muted/50 px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50">
                {#snippet children({ segments })}
                  {#each segments as { part, value }, i (part + value + i)}
                    <div class="inline-block select-none">
                      {#if part === "literal"}
                        <TimeField.Segment {part} class="text-muted-foreground p-1">
                          {value}
                        </TimeField.Segment>
                      {:else}
                        <TimeField.Segment
                          {part}
                          class="rounded px-1 py-1 hover:bg-muted focus:bg-muted focus:text-foreground aria-[valuetext=Empty]:text-muted-foreground focus-visible:ring-0 focus-visible:ring-offset-0"
                        >
                          {value}
                        </TimeField.Segment>
                      {/if}
                    </div>
                  {/each}
                {/snippet}
              </TimeField.Input>
            </TimeField.Root>
          </div>
        </div>
      </div>
      <div class="space-y-3">
        <Label for="form-max-responses">Max Responses</Label>
        <Input
          id="form-max-responses"
          type="number"
          min="1"
          placeholder="Unlimited"
          bind:value={formData.max_responses}
        />
      </div>
      <div class="space-y-3">
        <Label for="form-individual-limit">Individual Limit</Label>
        <Input
          id="form-individual-limit"
          type="number"
          min="1"
          bind:value={formData.individual_limit}
        />
      </div>
      <div onclick={e => e.stopPropagation()} tabindex="-1" aria-hidden="true" class="flex items-start gap-3 p-3 border rounded-md bg-muted/20">
        <Checkbox id="form-anonymous" bind:checked={formData.anonymous} class="mt-0.5" />
        <div>
          <Label for="form-anonymous" class="text-sm font-medium">Anonymous Responses</Label>
          <p class="text-xs text-muted-foreground">Allow users to submit responses without identification.</p>
        </div>
      </div>
      <div onclick={e => e.stopPropagation()} tabindex="-1" aria-hidden="true" class="flex items-start gap-3 p-3 border rounded-md bg-muted/20">
        <Checkbox id="form-editable-responses" bind:checked={formData.editable_responses} class="mt-0.5" />
        <div>
          <Label for="form-editable-responses" class="text-sm font-medium">Editable Responses</Label>
          <p class="text-xs text-muted-foreground">Allow users to edit their submitted responses.</p>
        </div>
      </div>
    </div>
  </div>
{/snippet}

<div class="fixed top-20 right-8 z-50">
  {#if useDialog}
    <!-- dialog on non-wide screens -->
    <Dialog.Root bind:open={dialogOpen} onOpenChange={handleDialogOpenChange}>
      <Dialog.Trigger class={buttonVariants({ variant: "outline" }) + " flex items-center gap-2"} aria-label="Form settings">
        <IconSettings class="size-5" />
      </Dialog.Trigger>
      <Dialog.Content class="max-w-md max-h-[80vh] overflow-y-auto">
        <Dialog.Header>
          <Dialog.Title>Form Settings</Dialog.Title>
          <Dialog.Description>Advanced configuration for this form.</Dialog.Description>
        </Dialog.Header>
        {@render settingsContent()}
      </Dialog.Content>
    </Dialog.Root>
  {:else}
    <!-- panel on wide screens -->
    <Panel.Root on:openchange={handlePanelChange}>
      <Panel.Trigger class={buttonVariants({ variant: "outline" }) + " flex items-center gap-2"} aria-label="Form settings">
        <IconSettings class="size-5" />
      </Panel.Trigger>
      <Panel.Content>
        <Panel.Header class="px-8 pt-8">
          <Panel.Title>Form Settings</Panel.Title>
          <Panel.Description>Advanced configuration for this form.</Panel.Description>
        </Panel.Header>
        {@render settingsContent()}
        <Panel.Footer>
        </Panel.Footer>
      </Panel.Content>
    </Panel.Root>
  {/if}
</div>
