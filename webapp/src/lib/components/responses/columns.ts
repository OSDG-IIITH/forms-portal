import type { ColumnDef } from "@tanstack/table-core";
import { renderComponent, renderSnippet } from "$lib/components/ui/data-table/index.js";
import { Checkbox } from "$lib/components/ui/checkbox/index.js";
import DataTableColumnHeader from "./data-table-column-header.svelte";
import DataTableRowActions from "./data-table-row-actions.svelte";
import SubmittedCell from "./submitted-cell.svelte";
import ResponseIdCell from "./response-id-cell.svelte";

export const labels = [
  {
    value: "bug",
    label: "Bug",
  },
  {
    value: "feature",
    label: "Feature",
  },
  {
    value: "documentation",
    label: "Documentation",
  },
];

export type Response = {
  id: string;
  status: "draft" | "completed" | "edited";
  respondent?: string;
  respondentName?: string;
  started: string;
  submitted?: string;
  edited?: string;
  modified?: string;
};

export const columns: ColumnDef<Response>[] = [
  {
    id: "select",
    header: ({ table }) =>
      renderComponent(Checkbox, {
        checked: table.getIsAllPageRowsSelected(),
        indeterminate: table.getIsSomePageRowsSelected() && !table.getIsAllPageRowsSelected(),
        onCheckedChange: (value) => table.toggleAllPageRowsSelected(!!value),
        "aria-label": "Select all",
        class: "translate-y-[2px]",
      }),
    cell: ({ row }) =>
      renderComponent(Checkbox, {
        checked: row.getIsSelected(),
        onCheckedChange: (value) => row.toggleSelected(!!value),
        "aria-label": "Select row",
        class: "translate-y-[2px]",
      }),
    enableSorting: false,
    enableHiding: false,
  },
  {
    accessorKey: "id",
    header: ({ column }) =>
      renderComponent(DataTableColumnHeader, { column, title: "Response ID" }),
    cell: ({ row }) =>
      renderComponent(ResponseIdCell, { value: row.getValue("id") as string }),
    enableSorting: false,
    enableHiding: true,
  },
  {
    accessorKey: "respondentName",
    header: ({ column }) =>
      renderComponent(DataTableColumnHeader, { column, title: "Respondent" }),
    cell: ({ row }) => {
      const name = row.getValue("respondentName") as string | undefined;
      return name || row.original.respondent || "Anonymous";
    },
  },
  {
    accessorKey: "submitted",
    header: ({ column }) =>
      renderComponent(DataTableColumnHeader, { column, title: "Submitted" }),
    cell: ({ row }) => {
      const submitted = row.getValue("submitted") as string | undefined;
      return renderComponent(SubmittedCell, {
        date: submitted,
      });
    },
  },
  {
    id: "actions",
    cell: ({ row }) => renderComponent(DataTableRowActions, { row }),
  },
];
