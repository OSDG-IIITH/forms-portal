<script lang="ts">
  import DataTable from "$lib/components/responses/data-table.svelte";
  import { columns } from "$lib/components/responses/columns.js";
  import type { PageData } from './$types';
  import CalendarChart from "$lib/components/responses/calendar-chart.svelte";
  import SummaryCard from "$lib/components/responses/summary-card.svelte";
  import Card from "$lib/components/ui/card/card.svelte";
  import { ChartNoAxesColumn, Settings } from "@lucide/svelte";
  
  export let data: PageData;
  
  const responsesData = data.responses.map((r: any) => ({
    ...r,
    respondentName: data.respondentMap?.[r.respondent]?.name
  }));

  const uniqueRespondentIds = new Set(data.responses.map((r: any) => r.respondent));

  const latestResponse = data.responses
    .filter((r: any) => r.submitted || r.started)
    .sort((a: any, b: any) => {
      const aDate = new Date(a.submitted || a.started).getTime();
      const bDate = new Date(b.submitted || b.started).getTime();
      return bDate - aDate;
    })[0];
  const lastResponseTime = latestResponse ? (latestResponse.submitted || latestResponse.started) : "-";
</script>

<div class="container mx-auto pt-24 px-8 py-8 min-h-screen space-y-8">
  <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
    <div class="flex flex-col gap-4 h-full">
      <div class="flex-1">
        <SummaryCard 
          totalResponses={responsesData.length} 
          uniqueRespondents={uniqueRespondentIds.size}
          lastResponse={lastResponseTime}
          opens={data.form?.opens}
          closes={data.form?.closes}
        />
      </div>
      <div class="flex flex-col sm:flex-row gap-5">
        <a href="" class="no-underline flex-1">
          <Card class="flex flex-row items-center gap-6 px-7 py-8 min-h-[160px]">
          </Card>
        </a>
        <a href="" class="no-underline flex-1">
          <Card class="flex flex-row items-center gap-6 px-7 py-8 min-h-[160px]">
          </Card>
        </a>
      </div>
    </div>
    <div class="flex flex-col h-full">
      <div class="flex-1">
        <CalendarChart responses={responsesData} />
      </div>
    </div>
  </div>
  <DataTable data={responsesData} {columns} />
</div>
