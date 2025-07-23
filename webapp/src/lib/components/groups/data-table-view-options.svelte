<script lang="ts">
  import type { Table } from "@tanstack/table-core";
  import { Button } from "$lib/components/ui/button/index.js";
  import {
    DropdownMenu,
    DropdownMenuCheckboxItem,
    DropdownMenuContent,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
  } from "$lib/components/ui/dropdown-menu/index.js";
  import { SlidersHorizontal } from "@lucide/svelte";
  import type { Group } from "./columns.js";

  interface $$Props {
    table: Table<Group>;
  }

  let { table }: $$Props = $props();
</script>

<DropdownMenu>
  <DropdownMenuTrigger>
    <Button variant="outline" size="sm" class="ml-auto hidden h-8 lg:flex">
      <SlidersHorizontal class="mr-2 h-4 w-4" />
      View
    </Button>
  </DropdownMenuTrigger>
  <DropdownMenuContent align="end" class="w-[150px]">
    <DropdownMenuLabel>Toggle columns</DropdownMenuLabel>
    <DropdownMenuSeparator />
    {#each table
      .getAllColumns()
      .filter((column) => typeof column.accessorFn !== "undefined" && column.getCanHide()) as column}
      <DropdownMenuCheckboxItem
        class="capitalize"
        checked={column.getIsVisible()}
        onCheckedChange={(value) => column.toggleVisibility(!!value)}
      >
        {column.id}
      </DropdownMenuCheckboxItem>
    {/each}
  </DropdownMenuContent>
</DropdownMenu>
