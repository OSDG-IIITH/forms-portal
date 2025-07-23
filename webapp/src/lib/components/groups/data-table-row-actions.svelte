<script lang="ts">
  import type { Row } from "@tanstack/table-core";
  import { Button } from "$lib/components/ui/button/index.js";
  import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuSeparator,
    DropdownMenuShortcut,
    DropdownMenuTrigger,
  } from "$lib/components/ui/dropdown-menu/index.js";
  import { MoreHorizontal } from "@lucide/svelte";
  import type { Group } from "./columns.js";
  import { page } from "$app/stores";
  import { goto } from "$app/navigation";

  interface $$Props {
    row: Row<Group>;
  }

  let { row }: $$Props = $props();

  const viewGroup = () => {
    const currentPath = $page.url.pathname;
    const groupPath = `${currentPath}/${row.original.id}`;
    goto(groupPath);
  };
</script>

<DropdownMenu>
  <DropdownMenuTrigger>
    <Button
      variant="ghost"
      class="flex h-8 w-8 p-0 data-[state=open]:bg-muted"
    >
      <MoreHorizontal class="h-4 w-4" />
      <span class="sr-only">Open menu</span>
    </Button>
  </DropdownMenuTrigger>
  <DropdownMenuContent align="end" class="w-[160px]">
<DropdownMenuItem onclick={viewGroup}>
  View Group
</DropdownMenuItem>
    <DropdownMenuSeparator />
    <DropdownMenuItem>
      Delete
      <DropdownMenuShortcut>⌘⌫</DropdownMenuShortcut>
    </DropdownMenuItem>
  </DropdownMenuContent>
</DropdownMenu>
