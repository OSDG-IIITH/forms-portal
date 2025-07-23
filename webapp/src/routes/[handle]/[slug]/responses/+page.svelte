<script lang="ts">
  import DataTable from "$lib/components/responses/data-table.svelte";
  import { columns } from "$lib/components/responses/columns.js";
  import type { PageData } from './$types';
  
  export let data: PageData;
  
  const responsesData = data.responses.map((r: any) => ({
    ...r,
    respondentName: data.respondentMap?.[r.respondent]?.name
  }));
</script>

<div class="container mx-auto px-8 py-8 min-h-screen">
  <h2 class="text-3xl font-bold tracking-tight mb-6">Responses</h2>
  <DataTable data={responsesData} {columns} />
</div>
