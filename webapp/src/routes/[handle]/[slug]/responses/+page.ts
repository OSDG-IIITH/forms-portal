import type { PageLoad } from './$types';

function readInt(value: string | null, fallback: number, min = 1, max = 100) {
  const parsed = Number.parseInt(value ?? '', 10);
  if (!Number.isFinite(parsed)) return fallback;
  return Math.min(max, Math.max(min, parsed));
}

export const load: PageLoad = async ({ params, fetch, url }) => {
  const { handle, slug } = params;
  const page = readInt(url.searchParams.get('page'), 1, 1, Number.MAX_SAFE_INTEGER);
  const limit = readInt(url.searchParams.get('limit'), 20);
  const offset = (page - 1) * limit;

  const res = await fetch(`/api/forms/${handle}/${slug}`);
  if (!res.ok) return { status: res.status, error: new Error('Form not found') };
  const form = await res.json();

  const responsesRes = await fetch(`/api/forms/${form.id}/responses?limit=${limit}&offset=${offset}`, { credentials: 'include' });
  if (!responsesRes.ok) return { status: responsesRes.status, error: new Error('Failed to fetch responses') };
  const responsesData = await responsesRes.json();
  const responses = responsesData.data ?? [];

  // Fetch respondent names for each response
  const respondentIds = Array.from(new Set(responses.map((r: any) => r.respondent).filter(Boolean)));
  let respondentMap: Record<string, { name: string }> = {};
  if (respondentIds.length > 0) {
    const results = await Promise.all(respondentIds.map(id => fetch(`/api/users/${id}`, { credentials: 'include' })));
    for (let i = 0; i < results.length; i++) {
      if (results[i].ok) {
        const user = await results[i].json();
        respondentMap[String(respondentIds[i])] = { name: user.name };
      }
    }
  }

  return { form, responses, pagination: responsesData.pagination, respondentMap, page, limit };
};
