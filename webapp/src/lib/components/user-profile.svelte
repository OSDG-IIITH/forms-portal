<script lang="ts">
  import { page } from "$app/stores";
  import { onMount } from "svelte";
  import { Button } from "$lib/components/ui/button";
  import {
    Popover,
    PopoverTrigger,
    PopoverContent,
  } from "$lib/components/ui/popover";
  import { IconLogout } from "@tabler/icons-svelte";

  interface UserInfo {
    id?: string;
    handle?: string;
    email?: string;
    name?: string;
  }

  let user: UserInfo | null = null;

  $: if ($page.data && $page.data.user) {
    user = $page.data.user;
  }

  onMount(async () => {
    if (!user) {
      try {
        const res = await fetch("/api/auth/info", { credentials: "include" });
        if (res.ok) {
          user = await res.json();
        }
      } catch (e) {
        // ignore
      }
    }
  });

  function getInitials(name?: string, handle?: string): string {
    if (name && name.trim()) {
      const parts = name.trim().split(/\s+/);
      if (parts.length >= 2) {
        return (parts[0][0] + parts[parts.length - 1][0]).toUpperCase();
      }
      return name.slice(0, 2).toUpperCase();
    }
    if (handle && handle.trim()) {
      return handle.slice(0, 2).toUpperCase();
    }
    return "U";
  }

  function handleLogout() {
    window.location.href = "/api/auth/logout";
  }
</script>

{#if user}
  <Popover>
    <PopoverTrigger>
      <Button
        variant="outline"
        size="sm"
        class="h-9 w-9 p-0 rounded-lg border hover:bg-accent/50 hover:border-primary/30 transition-all duration-200 hover:scale-105 flex items-center justify-center font-semibold text-xs text-foreground"
        aria-label="User profile menu"
      >
        {getInitials(user.name, user.handle)}
      </Button>
    </PopoverTrigger>

    <PopoverContent
      class="w-56 rounded-lg shadow-xl border border-border p-1 text-xs"
      align="end"
    >
      <div class="px-3 py-2.5 border-b border-border/50">
        <div class="font-semibold text-sm text-foreground truncate">
          {user.name || user.handle}
        </div>
        {#if user.handle}
          <div class="text-xs text-muted-foreground truncate mt-0.5">
            @{user.handle}
          </div>
        {/if}
        {#if user.email}
          <div class="text-xs text-muted-foreground truncate mt-0.5">
            {user.email}
          </div>
        {/if}
      </div>

      <div class="p-1">
        <button
          class="w-full flex items-center gap-2 text-xs rounded-md px-3 py-2 hover:bg-accent hover:text-accent-foreground transition-colors cursor-pointer text-left"
          onclick={handleLogout}
        >
          <IconLogout class="w-4 h-4 text-muted-foreground" />
          <span class="font-medium">Log out</span>
        </button>
      </div>
    </PopoverContent>
  </Popover>
{/if}
