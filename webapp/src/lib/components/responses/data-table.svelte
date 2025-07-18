<script lang="ts">
  import type { ColumnDef, TableOptions, SortingState, ColumnFiltersState, VisibilityState } from "@tanstack/table-core";
  import { createSvelteTable, FlexRender } from "$lib/components/ui/data-table/index.js";
  import { getCoreRowModel, getFacetedRowModel, getFacetedUniqueValues, getFilteredRowModel, getPaginationRowModel, getSortedRowModel } from "@tanstack/table-core";
  import { writable } from "svelte/store";
  import * as Table from "$lib/components/ui/table/index.js";
  import DataTablePagination from "./data-table-pagination.svelte";
  import DataTableToolbar from "./data-table-toolbar.svelte";
  import type { Response } from "./columns.js";
  import { page } from "$app/stores";
  import { goto } from "$app/navigation";

  interface $$Props {
    data: Response[];
    columns: ColumnDef<Response>[];
  }

  let { data, columns }: $$Props = $props();

  const handleRowClick = (responseId: string) => {
    const currentPath = $page.url.pathname;
    const responsePath = `${currentPath}/${responseId}`;
    goto(responsePath);
  };

  const sorting = writable<SortingState>([]);
  const columnFilters = writable<ColumnFiltersState>([]);
  const columnVisibility = writable<VisibilityState>({ id: false });
  const rowSelection = writable<Record<string, boolean>>({});

  const options = $derived.by(() => {
    const opts: TableOptions<Response> = {
      data,
      columns,
      state: {
        sorting: $sorting,
        columnFilters: $columnFilters,
        columnVisibility: $columnVisibility,
        rowSelection: $rowSelection,
      },
      enableRowSelection: true,
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
    return opts;
  });

  const table = $derived(createSvelteTable(options));
</script>

<div class="space-y-4">
  <DataTableToolbar {table} />
  <div class="rounded-md border">
    <Table.Root>
      <Table.Header>
        {#each table.getHeaderGroups() as headerGroup}
          <Table.Row>
            {#each headerGroup.headers as header}
              <Table.Head>
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
            <Table.Row 
              data-state={row.getIsSelected() && "selected"}
              class="cursor-pointer hover:bg-muted/50"
              onclick={() => handleRowClick(row.original.id)}
            >
              {#each row.getVisibleCells() as cell}
                <Table.Cell>
                  <FlexRender content={cell.column.columnDef.cell} context={cell.getContext()} />
                </Table.Cell>
              {/each}
            </Table.Row>
          {/each}
        {:else}
          <Table.Row>
            <Table.Cell colspan={columns.length} class="h-24 text-center">
              No results.
            </Table.Cell>
          </Table.Row>
        {/if}
      </Table.Body>
    </Table.Root>
  </div>
  <DataTablePagination {table} />
</div>
