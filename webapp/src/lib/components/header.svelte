<script>
  import ThemeSwitcher from "./theme-switcher.svelte";
  import { page } from "$app/stores";
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
</script>

<header class="border-b bg-card fixed top-0 left-0 right-0 z-50">
  <div class="container mx-auto px-6 py-4">
    <div class="flex items-center justify-between">
      <h1 class="text-3xl tracking-tight flex-shrink-0">
        <a href="/"
          >forms <span class="font-semibold text-accent-foreground">iiit</span
          ></a
        >
      </h1>
      {#if showtabs}
        <nav
          class="flex gap-4 items-center bg-transparent flex-1 justify-center"
        >
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
      <ThemeSwitcher />
    </div>
  </div>
</header>
