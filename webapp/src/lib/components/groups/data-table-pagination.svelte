<script lang="ts">
  import type { Table } from "@tanstack/table-core";
  import { Button } from "$lib/components/ui/button/index.js";
  import * as Select from "$lib/components/ui/select/index.js";
  import {
    ChevronLeft,
    ChevronRight,
    ChevronsLeft,
    ChevronsRight,
  } from "@lucide/svelte";
  import type { Group } from "./columns.js";

  interface $$Props {
    table: Table<Group>;
  }

  let { table }: $$Props = $props();
</script>


<div class="flex items-center justify-between px-2">
  <div class="flex-1 text-sm text-muted-foreground hidden sm:block whitespace-nowrap">
    {table.getFilteredSelectedRowModel().rows.length} of {table.getFilteredRowModel().rows.length} row(s) selected.
  </div>
  <div class="flex items-center space-x-6 lg:space-x-8 w-full justify-end">
    <div class="flex items-center space-x-2">
      <p class="text-sm font-medium sm:inline hidden">Rows per page</p>
      <p class="text-sm font-medium sm:hidden">Rows</p>
      <Select.Root
        allowDeselect={false}
        type="single"
        value={`${table.getState().pagination.pageSize}`}
        onValueChange={(value: string | undefined) => {
          if (value) table.setPageSize(Number(value));
        }}
      >
        <Select.Trigger class="h-8 w-[70px]">
          {String(table.getState().pagination.pageSize)}
        </Select.Trigger>
        <Select.Content side="top">
          {#each [10, 20, 30, 40, 50] as pageSize}
            <Select.Item value={`${pageSize}`}>{pageSize}</Select.Item>
          {/each}
        </Select.Content>
      </Select.Root>
    </div>
    <div class="flex w-[100px] items-center justify-center text-sm font-medium">
      <span class="sm:inline hidden">Page {table.getState().pagination.pageIndex + 1} of {table.getPageCount()}</span>
      <span class="sm:hidden">{table.getState().pagination.pageIndex + 1}/{table.getPageCount()}</span>
    </div>
    <div class="flex items-center space-x-2">
      <Button
        variant="outline"
        class="hidden h-8 w-8 p-0 lg:flex"
        onclick={() => table.setPageIndex(0)}
        disabled={!table.getCanPreviousPage()}
      >
        <span class="sr-only">Go to first page</span>
        <ChevronsLeft class="h-4 w-4" />
      </Button>
      <Button
        variant="outline"
        class="h-8 w-8 p-0"
        onclick={() => table.previousPage()}
        disabled={!table.getCanPreviousPage()}
      >
        <span class="sr-only">Go to previous page</span>
        <ChevronLeft class="h-4 w-4" />
      </Button>
      <Button
        variant="outline"
        class="h-8 w-8 p-0"
        onclick={() => table.nextPage()}
        disabled={!table.getCanNextPage()}
      >
        <span class="sr-only">Go to next page</span>
        <ChevronRight class="h-4 w-4" />
      </Button>
      <Button
        variant="outline"
        class="hidden h-8 w-8 p-0 lg:flex"
        onclick={() => table.setPageIndex(table.getPageCount() - 1)}
        disabled={!table.getCanNextPage()}
      >
        <span class="sr-only">Go to last page</span>
        <ChevronsRight class="h-4 w-4" />
      </Button>
    </div>
  </div>
</div>
