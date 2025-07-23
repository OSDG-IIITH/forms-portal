import type { ColumnDef } from '@tanstack/table-core';
import { renderComponent } from '$lib/components/ui/data-table/index.js';
import { Checkbox } from '$lib/components/ui/checkbox/index.js';
import DataTableColumnHeader from './data-table-column-header.svelte';
import DataTableRowActions from './data-table-row-actions.svelte';

export type Group = {
	id: string;
	owner: string;
	name: string;
	description?: string;
	type: 'domain' | 'list';
	domain?: string;
	members?: string[];
};

export const columns: ColumnDef<Group>[] = [
	{
		id: 'select',
		header: ({ table }) =>
			renderComponent(Checkbox, {
				checked: table.getIsAllPageRowsSelected(),
				indeterminate: table.getIsSomePageRowsSelected() && !table.getIsAllPageRowsSelected(),
				onCheckedChange: (value) => table.toggleAllPageRowsSelected(!!value),
				'aria-label': 'Select all',
				class: 'translate-y-[2px]'
			}),
		cell: ({ row }) =>
			renderComponent(Checkbox, {
				checked: row.getIsSelected(),
				onCheckedChange: (value) => row.toggleSelected(!!value),
				'aria-label': 'Select row',
				class: 'translate-y-[2px]'
			}),
		enableSorting: false,
		enableHiding: false
	},
	{
		accessorKey: 'id',
		header: ({ column }) => renderComponent(DataTableColumnHeader, { column, title: 'Group ID' }),
		cell: ({ row }) => row.getValue('id'),
		enableSorting: false,
		enableHiding: true
	},
	{
		accessorKey: 'name',
		header: ({ column }) => renderComponent(DataTableColumnHeader, { column, title: 'Name' }),
		cell: ({ row }) => row.getValue('name')
	},
	{
		accessorKey: 'type',
		header: ({ column }) => renderComponent(DataTableColumnHeader, { column, title: 'Type' }),
		cell: ({ row }) => row.getValue('type')
	},
	{
		accessorKey: 'owner',
		header: ({ column }) => renderComponent(DataTableColumnHeader, { column, title: 'Owner' }),
		cell: ({ row }) => row.getValue('owner')
	},
	{
		accessorKey: 'description',
		header: ({ column }) => renderComponent(DataTableColumnHeader, { column, title: 'Description' }),
		cell: ({ row }) => row.getValue('description') ?? ''
	},
	{
		id: 'actions',
		cell: ({ row }) => renderComponent(DataTableRowActions, { row })
	}
];