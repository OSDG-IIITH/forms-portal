<script lang="ts">
  import type { PageData } from './$types';
  import ResponseViewer from '$lib/components/responses/ResponseViewer.svelte';

  export let data: PageData;
</script>

<div class="max-w-4xl mx-auto py-8">
  <ResponseViewer form={data.form} answers={data.answers} />
</div>