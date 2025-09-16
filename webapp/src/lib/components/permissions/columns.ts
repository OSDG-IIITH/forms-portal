import type { ColumnDef } from "@tanstack/table-core";
import { renderComponent } from "$lib/components/ui/data-table/index.js";
import DataTableColumnHeader from "$lib/components/responses/data-table-column-header.svelte";
import UserNameCell from "./user-name-cell.svelte";
import UserEmailCell from "./user-email-cell.svelte";
import PermissionsCell from "./permissions-cell.svelte";
import PermissionActions from "./permission-actions.svelte";

export type Permission = {
  userId: string;
  userName?: string;
  userEmail?: string;
  roles: Array<"view" | "respond" | "comment" | "analyze" | "edit" | "manage">;
  permissions?: Array<{ id: string; role: string; user?: string; group?: string }>;
  type?: 'user' | 'group';
};

export function createColumns(
  onPermissionChange: (userId: string, action: 'add' | 'remove', role: string) => void,
  onRemoveUser: (userId: string) => void
): ColumnDef<Permission>[] {
  return [
    {
      accessorKey: "userName",
      header: ({ column }) =>
        renderComponent(DataTableColumnHeader, { column, title: "User/Group" }),
      cell: ({ row }) =>
        renderComponent(UserNameCell, { 
          userName: row.original.userName,
          userEmail: row.original.userEmail,
          userId: row.original.userId
        }),
    },
    {
      accessorKey: "userEmail",
      header: ({ column }) =>
        renderComponent(DataTableColumnHeader, { column, title: "Email/Domain" }),
      cell: ({ row }) =>
        renderComponent(UserEmailCell, { userEmail: row.original.userEmail }),
    },
    {
      accessorKey: "roles", 
      header: ({ column }) =>
        renderComponent(DataTableColumnHeader, { column, title: "Permissions" }),
      cell: ({ row }) =>
        renderComponent(PermissionsCell, { roles: row.getValue("roles") as Permission["roles"] }),
    },
    {
      id: "actions",
      enableHiding: false,
      cell: ({ row }) =>
        renderComponent(PermissionActions, {
          userId: row.original.userId,
          userName: row.original.userName,
          userEmail: row.original.userEmail,
          roles: row.original.roles,
          onPermissionChange,
          onRemoveUser,
        }),
    },
  ];
}
