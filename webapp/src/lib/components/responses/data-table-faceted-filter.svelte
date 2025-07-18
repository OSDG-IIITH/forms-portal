<script lang="ts">
  import type { Column } from "@tanstack/table-core";
  import { cn } from "$lib/utils.js";
  import { Badge } from "$lib/components/ui/badge/index.js";
  import { Button } from "$lib/components/ui/button/index.js";
  import {
    Popover,
    PopoverContent,
    PopoverTrigger,
  } from "$lib/components/ui/popover/index.js";
  import { Separator } from "$lib/components/ui/separator/index.js";
  import { Check, PlusCircle } from "@lucide/svelte";
  import type { Response } from "./columns.js";

  interface Option {
    label: string;
    value: string;
    icon?: string;
  }

  interface $$Props {
    column?: Column<Response, unknown>;
    title?: string;
    options: Option[];
  }

  let { column, title, options }: $$Props = $props();

  const facets = $derived(column?.getFacetedUniqueValues());
  const selectedValues = $derived(new Set(column?.getFilterValue() as string[]));
</script>

<Popover>
  <PopoverTrigger>
    <Button variant="outline" size="sm" class="h-8 border-dashed">
      <PlusCircle class="mr-2 h-4 w-4" />
      {title}
      {#if selectedValues.size > 0}
        <Separator orientation="vertical" class="mx-2 h-4" />
        <Badge variant="secondary" class="rounded-sm px-1 font-normal lg:hidden">
          {selectedValues.size}
        </Badge>
        <div class="hidden space-x-1 lg:flex">
          {#if selectedValues.size > 2}
            <Badge variant="secondary" class="rounded-sm px-1 font-normal">
              {selectedValues.size} selected
            </Badge>
          {:else}
            {#each options.filter((option) => selectedValues.has(option.value)) as option}
              <Badge variant="secondary" class="rounded-sm px-1 font-normal">
                {option.label}
              </Badge>
            {/each}
          {/if}
        </div>
      {/if}
    </Button>
  </PopoverTrigger>
  <PopoverContent class="w-[200px] p-0" align="start">
    <div class="p-2">
      <div class="space-y-1">
        {#each options as option}
          {@const isSelected = selectedValues.has(option.value)}
          <button
            type="button"
            class="flex w-full items-center space-x-2 p-2 hover:bg-muted rounded-sm text-left"
            onclick={() => {
              if (isSelected) {
                selectedValues.delete(option.value);
              } else {
                selectedValues.add(option.value);
              }
              const filterValues = Array.from(selectedValues);
              column?.setFilterValue(
                filterValues.length ? filterValues : undefined
              );
            }}
          >
            <div
              class={cn(
                "flex h-4 w-4 items-center justify-center rounded-sm border border-primary",
                isSelected
                  ? "bg-primary text-primary-foreground"
                  : "opacity-50"
              )}
            >
              {#if isSelected}
                <Check class="h-3 w-3" />
              {/if}
            </div>
            {#if option.icon}
              <span class="mr-2">{option.icon}</span>
            {/if}
            <span class="text-sm">{option.label}</span>
            {#if facets?.get(option.value)}
              <span class="ml-auto text-xs text-muted-foreground">
                {facets.get(option.value)}
              </span>
            {/if}
          </button>
        {/each}
      </div>
      {#if selectedValues.size > 0}
        <div class="border-t pt-2 mt-2">
          <Button
            variant="ghost"
            size="sm"
            class="w-full"
            onclick={() => column?.setFilterValue(undefined)}
          >
            Clear filters
          </Button>
        </div>
      {/if}
    </div>
  </PopoverContent>
</Popover>
