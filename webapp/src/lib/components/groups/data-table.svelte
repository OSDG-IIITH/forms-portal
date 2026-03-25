<script lang="ts">
  import type { ColumnDef, TableOptions, SortingState, ColumnFiltersState, VisibilityState } from "@tanstack/table-core";
	import {
		getCoreRowModel,
		getFacetedRowModel,
		getFacetedUniqueValues,
		getFilteredRowModel,
		getSortedRowModel
	} from '@tanstack/table-core';
    import { createSvelteTable, FlexRender } from "$lib/components/ui/data-table/index.js";
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import * as Table from '$lib/components/ui/table/index.js';
	import type { Group } from '$lib/types/group';
	import DataTablePagination from './data-table-pagination.svelte';
	import DataTableToolbar from './data-table-toolbar.svelte';
	import AddGroupDialog from './add-group-dialog.svelte';

	interface $$Props {
		data: Group[];
		columns: ColumnDef<Group>[];
		onGroupCreated: () => void;
		pageNumber: number;
		pageSize: number;
		totalCount: number;
	}

	let { data, columns, onGroupCreated, pageNumber, pageSize, totalCount }: $$Props = $props();

	const handleRowClick = (groupId: string) => {
		const groupPath = `${page.url.pathname}/${groupId}`;
		goto(groupPath);
	};

	function setPagination(nextPage: number, nextLimit: number) {
		const params = new URLSearchParams(page.url.searchParams);
		params.set('page', String(nextPage));
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
	let rowSelection = $state<Record<string, boolean>>({});

	const options = $derived.by(() => {
		const opts: TableOptions<Group> = {
			data,
			columns,
			state: {
				sorting,
				columnFilters,
				columnVisibility,
				rowSelection
			},
			enableRowSelection: true,
			onSortingChange: (updater) => (sorting = updater instanceof Function ? updater(sorting) : updater),
			onColumnFiltersChange: (updater) =>
				(columnFilters = updater instanceof Function ? updater(columnFilters) : updater),
			onColumnVisibilityChange: (updater) =>
				(columnVisibility = updater instanceof Function ? updater(columnVisibility) : updater),
			onRowSelectionChange: (updater) =>
				(rowSelection = updater instanceof Function ? updater(rowSelection) : updater),
			getCoreRowModel: getCoreRowModel(),
			getFilteredRowModel: getFilteredRowModel(),
			getSortedRowModel: getSortedRowModel(),
			getFacetedRowModel: getFacetedRowModel(),
			getFacetedUniqueValues: getFacetedUniqueValues()
		};
		return opts;
	});

	const table = $derived(createSvelteTable(options));
</script>

<div class="space-y-4">
	<DataTableToolbar {table}>
		<AddGroupDialog slot="add-group" onCreated={onGroupCreated} />
	</DataTableToolbar>
	<div class="rounded-md border bg-card">
		<Table.Root>
			<Table.Header>
				{#each table.getHeaderGroups() as headerGroup (headerGroup.id)}
					<Table.Row>
						{#each headerGroup.headers as header (header.id)}
							<Table.Head>
								{#if !header.isPlaceholder}
									<FlexRender
										content={header.column.columnDef.header}
										context={header.getContext()}
									/>
								{/if}
							</Table.Head>
						{/each}
					</Table.Row>
				{/each}
			</Table.Header>
			<Table.Body>
				{#if data.length > 0}
					{#each table.getRowModel().rows as row (row.id)}
						<Table.Row
							data-state={row.getIsSelected() && 'selected'}
							class="cursor-pointer hover:bg-muted/30"
							onclick={() => handleRowClick(row.original.id)}
						>
							{#each row.getVisibleCells() as cell (cell.id)}
								<Table.Cell
									onclick={(e) => {
										if (cell.column.id === 'select' || cell.column.id === 'actions') {
											e.stopPropagation();
										}
									}}
								>
									<FlexRender content={cell.column.columnDef.cell} context={cell.getContext()} />
								</Table.Cell>
							{/each}
						</Table.Row>
					{/each}
				{:else}
					<Table.Row>
						<Table.Cell colspan={columns.length} class="h-20 text-center">
							No groups found.
						</Table.Cell>
					</Table.Row>
				{/if}
			</Table.Body>
		</Table.Root>
	</div>
	<DataTablePagination
		page={pageNumber}
		limit={pageSize}
		total={totalCount}
		onPageChange={(nextPage) => setPagination(nextPage, pageSize)}
		onLimitChange={(nextLimit) => setPagination(1, nextLimit)}
	/>
</div>