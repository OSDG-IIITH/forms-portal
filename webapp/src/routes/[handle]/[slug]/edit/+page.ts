import type { PageLoad } from './$types';

export const load: PageLoad = async ({ params, fetch }) => {
  const { handle, slug } = params;
  const res = await fetch(`/api/forms/${handle}/${slug}`);
  if (!res.ok) {
    return { status: res.status, error: new Error('Form not found') };
  }
  const form = await res.json();
  return { form };
};