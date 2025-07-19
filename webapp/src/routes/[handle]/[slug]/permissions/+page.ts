import type { PageLoad } from './$types';

export const load: PageLoad = async ({ params, fetch }) => {
  const { handle, slug } = params;
  
  const res = await fetch(`/api/forms/${handle}/${slug}`);
  if (!res.ok) return { status: res.status, error: new Error('Form not found') };
  const form = await res.json();
  
  const permissionsRes = await fetch(`/api/forms/${form.id}/permissions`, { credentials: 'include' });
  if (!permissionsRes.ok) return { status: permissionsRes.status, error: new Error('Failed to fetch permissions') };
  const permissions = await permissionsRes.json();
  
  const userPermissions = permissions.filter((p: any) => p.user);
  const userIds = Array.from(new Set(userPermissions.map((p: any) => p.user)));
  
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

  const groupedPermissions = userIds.map(userId => {
    const userPerms = userPermissions.filter((p: any) => p.user === userId);
    const user = userMap[String(userId)];
    return {
      userId: String(userId),
      userName: user?.name,
      userEmail: user?.email,
      roles: userPerms.map((p: any) => p.role),
      permissions: userPerms
    };
  });
  
  return { form, permissions: groupedPermissions, userMap };
};
