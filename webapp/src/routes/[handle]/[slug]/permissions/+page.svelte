<script lang="ts">
  import { base } from '$app/paths';
  import DataTable from "$lib/components/permissions/data-table.svelte";
  import { createColumns } from "$lib/components/permissions/columns.js";
  import AddEntityDialog from "$lib/components/permissions/add-entity-dialog.svelte";
  import type { PageData } from './$types';
  import type { Permission } from '$lib/components/permissions/columns.js';
  import { toast } from "svelte-sonner";

  let { data }: { data: PageData } = $props();
  let permissions = $state<Permission[]>(
    (data.permissions ?? []).map(p => ({
      ...p,
      type: p.type === "user" || p.type === "group" ? p.type : undefined
    }))
  );

  async function refreshPermissions() {
    try {
      const permissionsRes = await fetch(`/api/forms/${data.form.id}/permissions`, { credentials: 'include' });
      if (!permissionsRes.ok) return;
      const permissionsData = await permissionsRes.json();
      
      const userPermissions = permissionsData.filter((p: any) => p.user);
      const groupPermissions = permissionsData.filter((p: any) => p.group);
      
      const userIds = Array.from(new Set(userPermissions.map((p: any) => p.user)));
      const groupIds = Array.from(new Set(groupPermissions.map((p: any) => p.group)));
      
      let userMap: Record<string, { name: string; email: string }> = {};
      if (userIds.length > 0) {
        const results = await Promise.all(userIds.map(id => fetch(`/api/users/${id}`, { credentials: 'include' })));
        for (let i = 0; i < results.length; i++) {
          if (results[i].ok) {
            const user = await results[i].json();
            userMap[String(userIds[i])] = { name: user.name, email: user.email };
          }
        }
      }

      let groupMap: Record<string, { name: string; type: string; domain?: string }> = {};
      if (groupIds.length > 0) {
        const results = await Promise.all(groupIds.map(id => fetch(`/api/groups/${id}`, { credentials: 'include' })));
        for (let i = 0; i < results.length; i++) {
          if (results[i].ok) {
            const group = await results[i].json();
            groupMap[String(groupIds[i])] = { 
              name: group.name, 
              type: group.type,
              domain: group.domain 
            };
          }
        }
      }

      const groupedPermissions = [
        ...userIds.map(userId => {
          const userPerms = userPermissions.filter((p: any) => p.user === userId);
          const user = userMap[String(userId)];
          return {
            userId: String(userId),
            userName: user?.name,
            userEmail: user?.email,
            roles: userPerms.map((p: any) => p.role),
            permissions: userPerms,
            type: 'user'
          };
        }),
        ...groupIds.map(groupId => {
          const groupPerms = groupPermissions.filter((p: any) => p.group === groupId);
          const group = groupMap[String(groupId)];
          const displayEmail = group?.type === 'domain' && group?.domain 
            ? group.domain 
            : group?.type === 'list' 
            ? 'List' 
            : 'Group';
          
          return {
            userId: String(groupId),
            userName: group?.name,
            userEmail: displayEmail,
            roles: groupPerms.map((p: any) => p.role),
            permissions: groupPerms,
            type: 'group'
          };
        })
      ];
      
      permissions = groupedPermissions.map(p => ({
        ...p,
        type: p.type === "user" || p.type === "group" ? p.type : undefined
      }));
    } catch (e) {
      console.error('Failed to refresh permissions:', e);
    }
  }

  async function handlePermissionChange(userId: string, action: 'add' | 'remove', role: string) {
    try {
      if (action === 'add') {
        const entry = permissions.find(p => p.userId === userId);
        if (!entry) {
          toast.error('Entry not found');
          return;
        }
        
        const payload = entry.type === 'group' 
          ? { role, group: userId }
          : { role, user: entry.userEmail };
          
        const res = await fetch(`${base}/api/forms/${data.form.id}/permissions`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          credentials: 'include',
          body: JSON.stringify(payload)
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
          const res = await fetch(`${base}/api/forms/${data.form.id}/permissions/${toRemove.id}`, {
            method: 'DELETE',
            credentials: 'include'
          });
          if (res.ok) {
            const entityType = userEntry?.type === 'group' ? 'group' : 'user';
            toast.success(`Removed ${role} permission from ${entityType}`);
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
      const entry = permissions.find(p => p.userId === userId);
      if (entry?.permissions) {
        const results = await Promise.all(
<<<<<<< HEAD
          entry.permissions.map(p =>
            fetch(`/api/forms/${data.form.id}/permissions/${p.id}`, { method: 'DELETE', credentials: 'include' })
=======
          userEntry.permissions.map(p =>
            fetch(`${base}/api/forms/${data.form.id}/permissions/${p.id}`, { method: 'DELETE', credentials: 'include' })
>>>>>>> a4a5648 (fix: support base url everywhere)
          )
        );
        if (results.every(r => r.ok)) {
          const entityType = entry?.type === 'group' ? 'Group' : 'User';
          toast.success(`${entityType} removed`);
          permissions = permissions.filter(p => p.userId !== userId);
        } else {
          toast.error('Failed to remove entity');
        }
      }
    } catch {
      toast.error('Error removing entity');
    }
  }

  async function handleAddEntity(event: CustomEvent<{ type: 'user' | 'group'; email?: string; groupId?: string }>) {
    const { type, email, groupId } = event.detail;
    try {
      if (type === 'user' && email) {
        toast.success('User added. Assign permissions using the menu.');
      } else if (type === 'group' && groupId) {
        const res = await fetch(`${base}/api/forms/${data.form.id}/permissions`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          credentials: 'include',
          body: JSON.stringify({ group: groupId, role: 'view' })
        });
        if (res.ok) {
          toast.success('Group added with view permission.');
          await refreshPermissions();
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

<div class="container mx-auto px-8 py-8 mt-20 min-h-screen">
  <div class="flex items-center mb-3">
    <h2 class="text-xl font-semibold tracking-tight flex-1">Permissions</h2>
    <AddEntityDialog on:add={handleAddEntity} formId={data.form.id} existingPermissions={permissions} />
  </div>
  <DataTable data={permissions} {columns} />
</div>
