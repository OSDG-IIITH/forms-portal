import { base } from '$app/paths';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ params, fetch }) => {
  const { handle, slug } = params;
  const res = await fetch(`${base}/api/forms/${handle}/${slug}`);
  if (!res.ok) return { status: res.status, error: new Error('Form not found') };
  const form = await res.json();
  const responsesRes = await fetch(`${base}/api/forms/${form.id}/responses?limit=20&offset=0`, { credentials: 'include' });
  if (!responsesRes.ok) return { status: responsesRes.status, error: new Error('Failed to fetch responses') };
  const responsesData = await responsesRes.json();
  const responses = responsesData.data ?? [];

  // Fetch respondent names for each response
  const respondentIds = Array.from(new Set(responses.map((r: any) => r.respondent).filter(Boolean)));
  let respondentMap: Record<string, { name: string }> = {};
  if (respondentIds.length > 0) {
    const results = await Promise.all(respondentIds.map(id => fetch(`${base}/api/users/${id}`, { credentials: 'include' })));
    for (let i = 0; i < results.length; i++) {
      if (results[i].ok) {
        const user = await results[i].json();
        respondentMap[String(respondentIds[i])] = { name: user.name };
      }
    }
  }

  return { form, responses, pagination: responsesData.pagination, respondentMap };
};
