<script lang="ts">
  import { Button } from "$lib/components/ui/button/index.js";
  import * as Select from "$lib/components/ui/select/index.js";
  import {
    ChevronLeft,
    ChevronRight,
    ChevronsLeft,
    ChevronsRight,
  } from "@lucide/svelte";
  interface $$Props {
    page: number;
    limit: number;
    total: number;
    onPageChange: (page: number) => void;
    onLimitChange: (limit: number) => void;
  }

  let { page, limit, total, onPageChange, onLimitChange }: $$Props = $props();

  const pageCount = $derived(Math.max(1, Math.ceil(total / limit)));
  const canPrevious = $derived(page > 1);
  const canNext = $derived(page < pageCount);
</script>


<div class="flex items-center justify-between px-2">
  <div class="flex-1 text-sm text-muted-foreground hidden sm:block whitespace-nowrap">
    {total} total
  </div>
  <div class="flex items-center space-x-6 lg:space-x-8 w-full justify-end">
    <div class="flex items-center space-x-2">
      <p class="text-sm font-medium sm:inline hidden">Rows per page</p>
      <p class="text-sm font-medium sm:hidden">Rows</p>
      <Select.Root
        allowDeselect={false}
        type="single"
        value={`${limit}`}
        onValueChange={(value: string | undefined) => {
          if (value) onLimitChange(Number(value));
        }}
      >
        <Select.Trigger class="h-8 w-[70px]">
          {String(limit)}
        </Select.Trigger>
        <Select.Content side="top">
          {#each [10, 20, 30, 40, 50] as pageSize (pageSize)}
            <Select.Item value={`${pageSize}`}>{pageSize}</Select.Item>
          {/each}
        </Select.Content>
      </Select.Root>
    </div>
    <div class="flex w-[100px] items-center justify-center text-sm font-medium">
      <span class="sm:inline hidden">Page {page} of {pageCount}</span>
      <span class="sm:hidden">{page}/{pageCount}</span>
    </div>
    <div class="flex items-center space-x-2">
      <Button
        variant="outline"
        class="hidden h-8 w-8 p-0 lg:flex"
        onclick={() => onPageChange(1)}
        disabled={!canPrevious}
      >
        <span class="sr-only">Go to first page</span>
        <ChevronsLeft class="h-4 w-4" />
      </Button>
      <Button
        variant="outline"
        class="h-8 w-8 p-0"
        onclick={() => onPageChange(page - 1)}
        disabled={!canPrevious}
      >
        <span class="sr-only">Go to previous page</span>
        <ChevronLeft class="h-4 w-4" />
      </Button>
      <Button
        variant="outline"
        class="h-8 w-8 p-0"
        onclick={() => onPageChange(page + 1)}
        disabled={!canNext}
      >
        <span class="sr-only">Go to next page</span>
        <ChevronRight class="h-4 w-4" />
      </Button>
      <Button
        variant="outline"
        class="hidden h-8 w-8 p-0 lg:flex"
        onclick={() => onPageChange(pageCount)}
        disabled={!canNext}
      >
        <span class="sr-only">Go to last page</span>
        <ChevronsRight class="h-4 w-4" />
      </Button>
    </div>
  </div>
</div>
