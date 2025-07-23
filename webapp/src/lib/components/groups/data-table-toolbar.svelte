<script lang="ts">
  import type { Table } from "@tanstack/table-core";
  import { Input } from "$lib/components/ui/input/index.js";
  import DataTableViewOptions from "./data-table-view-options.svelte";
  import { Button } from "$lib/components/ui/button/index.js";
  import { X } from "@lucide/svelte";
  import type { Group } from "./columns.js";

  interface $$Props {
    table: Table<Group>;
  }

  let { table }: $$Props = $props();

  const isFiltered = $derived(table.getState().columnFilters.length > 0);
</script>

<div class="flex items-center justify-between">
  <div class="flex flex-1 items-center space-x-2">
    <Input
      placeholder="Filter groups by name..."
      value={table.getColumn("name")?.getFilterValue() ?? ""}
      oninput={(e: Event) => table.getColumn("name")?.setFilterValue((e.target as HTMLInputElement).value)}
      class="h-8 w-[150px] lg:w-[250px]"
    />
    {#if isFiltered}
      <Button
        variant="ghost"
        onclick={() => table.resetColumnFilters()}
        class="h-8 px-2 lg:px-3"
      >
        Reset
        <X class="ml-2 h-4 w-4" />
      </Button>
    {/if}
  </div>
  <div class="flex items-center gap-2">
    <DataTableViewOptions {table} />
    <slot name="add-group" />
  </div>
</div>
