<script lang="ts">
  import FormEditor from '$lib/components/editor/FormEditor.svelte';
  import type { PageData } from './$types';

  const { data }: { data: PageData } = $props();
</script>

<FormEditor form={data.form} />
