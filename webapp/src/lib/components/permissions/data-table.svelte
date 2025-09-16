<script lang="ts">
import type { ColumnDef, TableOptions, SortingState, ColumnFiltersState, VisibilityState } from "@tanstack/table-core";
import { createSvelteTable, FlexRender } from "$lib/components/ui/data-table/index.js";
import { getCoreRowModel, getFacetedRowModel, getFacetedUniqueValues, getFilteredRowModel, getPaginationRowModel, getSortedRowModel } from "@tanstack/table-core";
import { writable } from "svelte/store";
import * as Table from "$lib/components/ui/table/index.js";
import type { Permission } from "./columns.js";

interface Props {
  data: Permission[];
  columns: ColumnDef<Permission>[];
}

let { data, columns }: Props = $props();

const sorting = writable<SortingState>([]);
const columnFilters = writable<ColumnFiltersState>([]);
const columnVisibility = writable<VisibilityState>({});
const rowSelection = writable<Record<string, boolean>>({});

const table = $derived.by(() => {
  const opts: TableOptions<Permission> = {
    data,
    columns,
    state: {
      sorting: $sorting,
      columnFilters: $columnFilters,
      columnVisibility: $columnVisibility,
      rowSelection: $rowSelection,
    },
    enableRowSelection: false,
    onSortingChange: (updater) => {
      if (updater instanceof Function) {
        sorting.update(updater);
      } else {
        sorting.set(updater);
      }
    },
    onColumnFiltersChange: (updater) => {
      if (updater instanceof Function) {
        columnFilters.update(updater);
      } else {
        columnFilters.set(updater);
      }
    },
    onColumnVisibilityChange: (updater) => {
      if (updater instanceof Function) {
        columnVisibility.update(updater);
      } else {
        columnVisibility.set(updater);
      }
    },
    onRowSelectionChange: (updater) => {
      if (updater instanceof Function) {
        rowSelection.update(updater);
      } else {
        rowSelection.set(updater);
      }
    },
    getCoreRowModel: getCoreRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    getSortedRowModel: getSortedRowModel(),
    getFacetedRowModel: getFacetedRowModel(),
    getFacetedUniqueValues: getFacetedUniqueValues(),
  };
  return createSvelteTable(opts);
});
</script>

<div class="space-y-4">
  <div class="rounded-md border overflow-x-auto bg-card">
    <Table.Root class="min-w-full">
      <Table.Header>
        {#each table.getHeaderGroups() as headerGroup}
          <Table.Row>
            {#each headerGroup.headers as header}
              <Table.Head class="first:pl-3 last:w-[70px] last:text-right">
                {#if !header.isPlaceholder}
                  <FlexRender content={header.column.columnDef.header} context={header.getContext()} />
                {/if}
              </Table.Head>
            {/each}
          </Table.Row>
        {/each}
      </Table.Header>
      <Table.Body>
        {#if table.getRowModel().rows?.length}
          {#each table.getRowModel().rows as row}
            <Table.Row data-state={row.getIsSelected() && "selected"} class="hover:bg-muted/10 border-b transition-colors">
              {#each row.getVisibleCells() as cell}
                <Table.Cell class="first:pl-8 last:pr-4 last:text-right">
                  <FlexRender content={cell.column.columnDef.cell} context={cell.getContext()} />
                </Table.Cell>
              {/each}
            </Table.Row>
          {/each}
        {:else}
          <Table.Row>
            <Table.Cell colspan={columns.length} class="h-24 text-center">
              No permissions found.
            </Table.Cell>
          </Table.Row>
        {/if}
      </Table.Body>
    </Table.Root>
  </div>
</div>
