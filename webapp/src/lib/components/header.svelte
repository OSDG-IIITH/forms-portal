<script>
  import ThemeSwitcher from "./theme-switcher.svelte";
  import { page } from "$app/stores";
  import { Button } from '$lib/components/ui/button';
  import { IconMenu2 } from '@tabler/icons-svelte';
  $: currentPath = $page.url.pathname;
  $: match = currentPath.match(/\/(edit|responses|permissions)$/);
  $: showtabs = !!match;
  $: activeTab = match ? match[1] : null;
  const tabs = [
    { name: "Edit", key: "edit" },
    { name: "Responses", key: "responses" },
    { name: "Permissions", key: "permissions" },
  ];
  $: basePath = showtabs
    ? currentPath.replace(/\/(edit|responses|permissions)$/, "")
    : "";

  let showMobileNav = false;
  let prevPath = $page.url.pathname;
  $: if ($page.url.pathname !== prevPath) {
    showMobileNav = false;
    prevPath = $page.url.pathname;
  }
</script>

<header class="border-b bg-card fixed top-0 left-0 right-0 z-50">
  <div class="container mx-auto px-6 py-4">
    <div class="flex items-center w-full">
      <h1 class="text-3xl tracking-tight flex-shrink-0">
        <a href="/"
          >forms <span class="font-semibold text-accent-foreground">iiit</span
          ></a
        >
      </h1>
      {#if showtabs}
        <!-- regular navbar -->
        <nav class="hidden md:flex gap-4 items-center bg-transparent ml-auto mr-4 flex-1 justify-center">
          {#each tabs as tab}
            <a
              href="{basePath}/{tab.key}"
              class="px-3 py-1 rounded transition-colors duration-150 font-medium {activeTab ===
              tab.key
                ? 'text-foreground/90 hover:text-foreground'
                : 'text-muted-foreground/80 hover:text-foreground'}"
            >
              {tab.name}
            </a>
          {/each}
        </nav>
      {/if}
      <div class="flex items-center gap-2 ml-auto">
        {#if showtabs}
          <!-- menu button -->
          <button
            type="button"
            class="md:hidden h-9 px-3 rounded-lg bg-background shadow-xs hover:bg-accent hover:text-accent-foreground dark:bg-input/30 dark:border-input dark:hover:bg-input/50 border"
            aria-label="Open navigation menu"
            on:click={() => (showMobileNav = !showMobileNav)}
          >
            <IconMenu2 class="w-4 h-4" />
          </button>
        {/if}
        <ThemeSwitcher />
      </div>
    </div>
  </div>
</header>
{#if showtabs && showMobileNav}
  <!-- mobile navbar -->
  <div class="md:hidden w-full bg-card border-b px-4 py-2 sticky top-[64px] z-40 flex justify-center">
    <nav class="flex gap-2 items-center">
      {#each tabs as tab}
        <a
          href="{basePath}/{tab.key}"
          class="px-3 py-1 rounded transition-colors duration-150 font-medium whitespace-nowrap {activeTab ===
          tab.key
            ? 'text-foreground/90 hover:text-foreground'
            : 'text-muted-foreground/80 hover:text-foreground'}"
        >
          {tab.name}
        </a>
      {/each}
    </nav>
  </div>
{/if}
