<script lang="ts">
  import DataTable from "$lib/components/permissions/data-table.svelte";
  import { createColumns } from "$lib/components/permissions/columns.js";
  import AddEntityDialog from "$lib/components/permissions/add-entity-dialog.svelte";
  import type { PageData } from './$types';
  import type { Permission } from '$lib/components/permissions/columns.js';
  import { toast } from "svelte-sonner";

  let { data }: { data: PageData } = $props();
  let permissions = $state<Permission[]>(data.permissions ?? []);

  async function handlePermissionChange(userId: string, action: 'add' | 'remove', role: string) {
    try {
      if (action === 'add') {
        const userEntry = permissions.find(p => p.userId === userId);
        const userEmail = userEntry?.userEmail;
        if (!userEmail) {
          toast.error('User email not found');
          return;
        }
        const res = await fetch(`/api/forms/${data.form.id}/permissions`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          credentials: 'include',
          body: JSON.stringify({ role, user: userEmail })
        });
        if (res.ok) {
          const newPerm = await res.json();
          toast.success(`Added ${role} permission`);
          permissions = permissions.map(p =>
            p.userId === userId
              ? {
                  ...p,
                  roles: [...p.roles, role as any],
                  permissions: [...(p.permissions ?? []), newPerm]
                }
              : p
          );
        } else {
          const err = await res.text();
          toast.error(`Failed to add permission: ${err}`);
        }
      } else {
        const userEntry = permissions.find(p => p.userId === userId);
        const toRemove = userEntry?.permissions?.find(p => p.role === role);
        if (toRemove) {
          const res = await fetch(`/api/forms/${data.form.id}/permissions/${toRemove.id}`, {
            method: 'DELETE',
            credentials: 'include'
          });
          if (res.ok) {
            toast.success(`Removed ${role} permission`);
            permissions = permissions
              .map(p =>
                p.userId === userId
                  ? {
                      ...p,
                      roles: p.roles.filter(r => r !== role),
                      permissions: (p.permissions ?? []).filter(pr => pr.role !== role)
                    }
                  : p
              )
              .filter(p => p.roles.length > 0);
          } else {
            toast.error('Failed to remove permission');
          }
        }
      }
    } catch (e) {
      toast.error('Error managing permission');
    }
  }

  async function handleRemoveUser(userId: string) {
    try {
      const userEntry = permissions.find(p => p.userId === userId);
      if (userEntry?.permissions) {
        const results = await Promise.all(
          userEntry.permissions.map(p =>
            fetch(`/api/forms/${data.form.id}/permissions/${p.id}`, { method: 'DELETE', credentials: 'include' })
          )
        );
        if (results.every(r => r.ok)) {
          toast.success('User removed');
          permissions = permissions.filter(p => p.userId !== userId);
        } else {
          toast.error('Failed to remove user');
        }
      }
    } catch {
      toast.error('Error removing user');
    }
  }

  async function handleAddEntity(event: CustomEvent<{ type: 'user' | 'group'; email?: string; groupId?: string }>) {
    const { type, email, groupId } = event.detail;
    try {
      if (type === 'user' && email) {
        toast.success('User added. Assign permissions using the menu.');
      } else if (type === 'group' && groupId) {
        const res = await fetch(`/api/forms/${data.form.id}/permissions`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          credentials: 'include',
          body: JSON.stringify({ group: groupId, role: 'view' })
        });
        if (res.ok) {
          toast.success('Group added with view permission.');
        } else {
          const err = await res.text();
          toast.error(`Failed to add group: ${err}`);
        }
      }
    } catch (e) {
      toast.error('Failed to add entity.');
    }
  }

  const columns = createColumns(handlePermissionChange, handleRemoveUser);
</script>

<div class="container mx-auto py-8 mt-20 min-h-screen">
  <div class="flex items-center mb-3">
    <h2 class="text-xl font-semibold tracking-tight flex-1">Permissions</h2>
    <AddEntityDialog on:add={handleAddEntity} />
  </div>
  <DataTable data={permissions} {columns} />
</div>
