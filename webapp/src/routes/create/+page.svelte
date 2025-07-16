<script lang="ts">
  import FormEditor from '$lib/components/FormEditor.svelte';
  import { toast } from 'svelte-sonner';
  import { Time } from '@internationalized/date';

  function slugify(title: string, suffix = 0): string {
    let base = title
      .toLowerCase()
      .replace(/[^a-z0-9]+/g, '-')
      .replace(/(^-|-$)/g, '')
      .slice(0, 50);
    if (suffix > 0) return `${base}-${suffix}`.slice(0, 64);
    return base;
  }

  let form = {
    title: 'Untitled Form',
    description: 'Add a description',
    visibility: 'private',
    structure: ''
  };

  async function handleCreate(event: CustomEvent) {
    const { formData, questions, kdl } = event.detail;
    let suffix = 0;
    
    const formatDateTime = (dateStr: string | undefined, time: Time | undefined): string | null => {
      if (!dateStr || dateStr.trim() === '') return null;
      try {
        const hours = time?.hour ?? 0;
        const minutes = time?.minute ?? 0;
        const seconds = 0;
        
        const date = new Date(dateStr);
        if (isNaN(date.getTime())) return null;
        
        date.setHours(hours, minutes, seconds, 0);
        return date.toISOString();
      } catch {
        return null;
      }
    };

    for (let attempt = 0; attempt < 100; attempt++) {
      const payload = {
        title: formData.title || 'Untitled Form',
        slug: slugify(formData.title, suffix),
        description: formData.description || null,
        structure: kdl,
        opens: formatDateTime(formData.opens, formData.opensTime),
        closes: formatDateTime(formData.closes, formData.closesTime),
        anonymous: Boolean(formData.anonymous),
        max_responses: formData.max_responses !== null && formData.max_responses !== undefined ? Number(formData.max_responses) : null,
        individual_limit: Number(formData.individual_limit) || 1,
        editable_responses: Boolean(formData.editable_responses)
      };
      let res, text = '';
      try {
        res = await fetch(`/api/forms`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          credentials: 'include',
          body: JSON.stringify(payload)
        });
        text = await res.text();
      } catch (e) {
        toast.error('Network error', { description: e instanceof Error ? e.message : String(e) });
        return;
      }
      if (res.ok) {
        let data;
        try {
          data = JSON.parse(text);
        } catch {}
        let slug = data?.slug || payload.slug;
        let handle;
        try {
          const userRes = await fetch('/api/auth/info', { credentials: 'include' });
          if (userRes.ok) {
            const user = await userRes.json();
            handle = user.handle;
          }
        } catch {}
        if (handle && slug) {
          window.location.href = `/${handle}/${slug}/edit`;
        } else {
          toast.success('Form created successfully!');
        }
        return;
      }
      if (
        res.status === 409 ||
        text.includes('duplicate key value') ||
        text.includes('forms_owner_slug_key')
      ) {
        suffix++;
        continue;
      }
      toast.error('Error creating form', { description: text || 'Failed to create form.' });
      return;
    }
    toast.error('Error creating form', { description: 'Could not generate a unique title/slug after many attempts.' });
  }
</script>

<div class="min-h-screen bg-background">
  <main>
    <FormEditor {form} on:create={handleCreate} mode="create" />
  </main>
</div>