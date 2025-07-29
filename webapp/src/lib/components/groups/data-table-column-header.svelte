<script lang="ts">
  import type { Column } from "@tanstack/table-core";
  import { cn } from "$lib/utils.js";
  import { Button } from "$lib/components/ui/button/index.js";
  import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
  } from "$lib/components/ui/dropdown-menu/index.js";
  import { ArrowDown, ArrowUp, ChevronsUpDown, EyeOff } from "@lucide/svelte";

  interface $$Props {
    column: Column<any, unknown>;
    title: string;
    class?: string;
  }

  let { column, title, class: className, ...restProps }: $$Props = $props();

  const canSort = $derived(column.getCanSort());
  const sorted = $derived(column.getIsSorted());
</script>

<div class={cn("flex items-center space-x-2", className)} {...restProps}>
  {#if canSort}
    <DropdownMenu>
      <DropdownMenuTrigger>
        <Button variant="ghost" size="sm" class="ml-3 h-8 data-[state=open]:bg-accent">
          <span>{title}</span>
          {#if sorted === "desc"}
            <ArrowDown class="ml-2 h-4 w-4" />
          {:else if sorted === "asc"}
            <ArrowUp class="ml-2 h-4 w-4" />
          {:else}
            <ChevronsUpDown class="ml-2 h-4 w-4" />
          {/if}
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="start">
        <DropdownMenuItem onclick={() => column.toggleSorting(false)}>
          <ArrowUp class="mr-2 h-3.5 w-3.5 text-muted-foreground/70" />
          Asc
        </DropdownMenuItem>
        <DropdownMenuItem onclick={() => column.toggleSorting(true)}>
          <ArrowDown class="mr-2 h-3.5 w-3.5 text-muted-foreground/70" />
          Desc
        </DropdownMenuItem>
        <DropdownMenuSeparator />
        <DropdownMenuItem onclick={() => column.toggleVisibility(false)}>
          <EyeOff class="mr-2 h-3.5 w-3.5 text-muted-foreground/70" />
          Hide
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  {:else}
    <span>{title}</span>
  {/if}
</div>
