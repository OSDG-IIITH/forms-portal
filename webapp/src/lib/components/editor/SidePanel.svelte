<script lang="ts">
  import * as Panel from "$lib/components/ui/panel";
  import * as Dialog from "$lib/components/ui/dialog";
  import * as Sheet from "$lib/components/ui/sheet";
  import { buttonVariants } from "$lib/components/ui/button";
  import { IconSettings, IconMessage } from '@tabler/icons-svelte';
  import { createEventDispatcher } from 'svelte';
  import FormSettings from './FormSettings.svelte';
  import FormComments from './FormComments.svelte';

  let { dialogOpen = $bindable(false) } = $props();
  
  const dispatch = createEventDispatcher<{
    panelchange: { open: boolean };
    dialogopenchange: boolean;
  }>();

  let windowWidth = $state(0);
  let activeTab = $state<'settings' | 'comments'>('settings');
  let panelOpen = $state(false);

  const useSheet = $derived(windowWidth > 0 && windowWidth < 1380);

  function handlePanelChange(event: CustomEvent<{ open: boolean }>) {
    panelOpen = event.detail.open;
    dispatch('panelchange', event.detail);
  }

  function handleSheetOpenChange(open: boolean) {
    dialogOpen = open;
    dispatch('dialogopenchange', open);
    if (useSheet && !open) {
      dispatch('panelchange', { open: false });
    }
  }

  function switchToTab(tab: 'settings' | 'comments') {
    activeTab = tab;
    if (useSheet && !dialogOpen) {
      dialogOpen = true;
      dispatch('dialogopenchange', true);
    } else if (!useSheet && !panelOpen) {
      panelOpen = true;
      dispatch('panelchange', { open: true });
    }
  }

  function handleCommentsClick() {
    switchToTab('comments');
  }

  function handleSettingsClick() {
    switchToTab('settings');
  }

  function handleCommentsOpenChange(event: CustomEvent<{ open: boolean }>) {
    if (event.detail.open) {
      activeTab = 'comments';
    }
    handlePanelChange(event);
  }

  function handleSettingsOpenChange(event: CustomEvent<{ open: boolean }>) {
    if (event.detail.open) {
      activeTab = 'settings';
    }
    handlePanelChange(event);
  }

  $effect(() => {
    if (typeof window !== 'undefined') {
      windowWidth = window.innerWidth;
      
      const handleResize = () => {
        windowWidth = window.innerWidth;
        if (windowWidth >= 1024 && dialogOpen) {
          handleSheetOpenChange(false);
        }
      };
      
      window.addEventListener('resize', handleResize);
      return () => window.removeEventListener('resize', handleResize);
    }
  });
</script>

{#snippet currentContent()}
  {#if activeTab === 'settings'}
    <FormSettings />
  {:else}
    <FormComments />
  {/if}
{/snippet}

<div class="fixed top-20 right-8 z-50 flex gap-2">
  {#if useSheet}
    <button
      class={buttonVariants({ variant: dialogOpen && activeTab === 'comments' ? 'default' : 'outline' }) + ' flex items-center gap-2'}
      aria-label="Form comments"
      onclick={handleCommentsClick}
    >
      <IconMessage class="size-5" />
    </button>
    
    <Sheet.Root bind:open={dialogOpen} onOpenChange={handleSheetOpenChange}>
    <Sheet.Trigger
      class={buttonVariants({ variant: dialogOpen && activeTab === 'settings' ? 'default' : 'outline' }) + ' flex items-center gap-2'}
      aria-label="Form settings"
      onclick={handleSettingsClick}
    >
        <IconSettings class="size-5" />
      </Sheet.Trigger>
      <Sheet.Content class="max-w-md h-screen max-h-screen flex flex-col overflow-hidden pt-8 px-2">
        <Sheet.Header>
          <Sheet.Title>{activeTab === 'settings' ? 'Form Settings' : 'Form Comments'}</Sheet.Title>
          <Sheet.Description>{activeTab === 'settings' ? 'Advanced configuration for this form.' : 'Comments and feedback for this form.'}</Sheet.Description>
        </Sheet.Header>
        {@render currentContent()}
      </Sheet.Content>
    </Sheet.Root>
  {:else}
    <Panel.Root open={panelOpen && activeTab === 'comments'} on:openchange={handleCommentsOpenChange}>
      <Panel.Trigger
        class={buttonVariants({ variant: panelOpen && activeTab === 'comments' ? 'default' : 'outline' }) + ' flex items-center gap-2'}
        aria-label="Form comments"
      >
        <IconMessage class="size-5" />
      </Panel.Trigger>
      <Panel.Content class="flex flex-col">
        <Panel.Header class="px-8 pt-8">
          <Panel.Title>Form Comments</Panel.Title>
          <Panel.Description>Comments and feedback for this form.</Panel.Description>
        </Panel.Header>
        <FormComments />
      </Panel.Content>
    </Panel.Root>
    
    <Panel.Root open={panelOpen && activeTab === 'settings'} on:openchange={handleSettingsOpenChange}>
      <Panel.Trigger
        class={buttonVariants({ variant: panelOpen && activeTab === 'settings' ? 'default' : 'outline' }) + ' flex items-center gap-2'}
        aria-label="Form settings"
      >
        <IconSettings class="size-5" />
      </Panel.Trigger>
      <Panel.Content class="flex flex-col">
        <Panel.Header class="px-8 pt-8">
          <Panel.Title>Form Settings</Panel.Title>
          <Panel.Description>Advanced configuration for this form.</Panel.Description>
        </Panel.Header>
        <FormSettings />
        <Panel.Footer>
        </Panel.Footer>
      </Panel.Content>
    </Panel.Root>
  {/if}
</div>
