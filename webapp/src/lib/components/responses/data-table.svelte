<script lang="ts">
  import type { ColumnDef, TableOptions, SortingState, ColumnFiltersState, VisibilityState } from "@tanstack/table-core";
  import { createSvelteTable, FlexRender } from "$lib/components/ui/data-table/index.js";
  import { getCoreRowModel, getFacetedRowModel, getFacetedUniqueValues, getFilteredRowModel, getSortedRowModel } from "@tanstack/table-core";
  import * as Table from "$lib/components/ui/table/index.js";
  import DataTablePagination from "./data-table-pagination.svelte";
  import DataTableToolbar from "./data-table-toolbar.svelte";
  import type { Response } from "./columns.js";
  import { page } from "$app/state";
  import { goto } from "$app/navigation";

  interface $$Props {
    data: Response[];
    columns: ColumnDef<Response>[];
    pageNumber: number;
    pageSize: number;
    totalCount: number;
  }

  let { data, columns, pageNumber, pageSize, totalCount }: $$Props = $props();

  const handleRowClick = (responseId: string) => {
    const currentPath = page.url.pathname;
    const responsePath = `${currentPath}/${responseId}`;
    goto(responsePath);
  };

  const pageCount = $derived(Math.max(1, Math.ceil(totalCount / pageSize)));
  const shownPage = $derived(Math.min(Math.max(pageNumber, 1), pageCount));

  function setPagination(nextPage: number, nextLimit: number) {
    if (nextLimit < 1) return;

    const clampedPage = Math.min(Math.max(nextPage, 1), pageCount);

    const params = new URLSearchParams(page.url.searchParams);
    const currentPage = Number.parseInt(params.get('page') ?? '1', 10);
    const currentLimit = Number.parseInt(params.get('limit') ?? String(pageSize), 10);

    if (currentPage === clampedPage && currentLimit === nextLimit) {
      return;
    }

    params.set('page', String(clampedPage));
    params.set('limit', String(nextLimit));
    goto(`${page.url.pathname}?${params.toString()}`, {
      replaceState: true,
      keepFocus: true,
      noScroll: true
    });
  }

  let sorting = $state<SortingState>([]);
  let columnFilters = $state<ColumnFiltersState>([]);
  let columnVisibility = $state<VisibilityState>({ id: false });

  const options = $derived.by(() => {
    const opts: TableOptions<Response> = {
      data,
      columns,
      state: {
        sorting,
        columnFilters,
		columnVisibility,
      },
      onSortingChange: (updater) => {
        sorting = updater instanceof Function ? updater(sorting) : updater;
      },
      onColumnFiltersChange: (updater) => {
        columnFilters = updater instanceof Function ? updater(columnFilters) : updater;
      },
      onColumnVisibilityChange: (updater) => {
        columnVisibility = updater instanceof Function ? updater(columnVisibility) : updater;
      },
      getCoreRowModel: getCoreRowModel(),
      getFilteredRowModel: getFilteredRowModel(),
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
  <div class="rounded-md border bg-card">
    <Table.Root>
      <Table.Header>
        {#each table.getHeaderGroups() as headerGroup (headerGroup.id)}
          <Table.Row>
            {#each headerGroup.headers as header (header.id)}
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
          {#each table.getRowModel().rows as row (row.id)}
            <Table.Row 
              class="cursor-pointer hover:bg-muted/50"
              onclick={() => handleRowClick(row.original.id)}
            >
              {#each row.getVisibleCells() as cell (cell.id)}
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
  <DataTablePagination
    page={shownPage}
    limit={pageSize}
    total={totalCount}
    onPageChange={(nextPage) => setPagination(nextPage, pageSize)}
    onLimitChange={(nextLimit) => setPagination(1, nextLimit)}
  />
</div>
