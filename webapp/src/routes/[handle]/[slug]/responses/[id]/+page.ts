import { base } from '$app/paths';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ params, fetch }) => {
  const { handle, slug, id } = params;
  const formRes = await fetch(`${base}/api/forms/${handle}/${slug}`);
  if (!formRes.ok) return { status: formRes.status, error: new Error('Form not found') };
  const form = await formRes.json();

  const responseRes = await fetch(`${base}/api/forms/${form.id}/responses/${id}`, { credentials: 'include' });
  if (!responseRes.ok) return { status: responseRes.status, error: new Error('Failed to fetch response') };
  const response = await responseRes.json();

  const answersRes = await fetch(`${base}/api/forms/${form.id}/responses/${id}/answers`, {
    credentials: 'include'
  });
  if (!answersRes.ok) return { status: answersRes.status, error: new Error('Failed to fetch answers') };
  const answersData = await answersRes.json();
  const answers = answersData.data;

  return { form, response, answers };
};
