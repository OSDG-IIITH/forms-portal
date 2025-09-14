import { base } from '$app/paths';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ params, fetch }) => {
  const { handle, slug } = params;
  
  const res = await fetch(`${base}/api/forms/${handle}/${slug}`);
  if (!res.ok) return { status: res.status, error: new Error('Form not found') };
  const form = await res.json();
  
  const permissionsRes = await fetch(`${base}/api/forms/${form.id}/permissions`, { credentials: 'include' });
  if (!permissionsRes.ok) return { status: permissionsRes.status, error: new Error('Failed to fetch permissions') };
  const permissions = await permissionsRes.json();
  
  const userPermissions = permissions.filter((p: any) => p.user);
  const groupPermissions = permissions.filter((p: any) => p.group);
  
  const userIds = Array.from(new Set(userPermissions.map((p: any) => p.user)));
  const groupIds = Array.from(new Set(groupPermissions.map((p: any) => p.group)));
  
  let userMap: Record<string, { name: string; email: string }> = {};
  if (userIds.length > 0) {
    const results = await Promise.all(userIds.map(id => fetch(`${base}/api/users/${id}`, { credentials: 'include' })));
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
  
  return { form, permissions: groupedPermissions, userMap, groupMap };
};
