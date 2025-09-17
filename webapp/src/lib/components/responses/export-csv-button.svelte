<script lang="ts">
  import { Button } from "$lib/components/ui/button";
  import { Checkbox } from "$lib/components/ui/checkbox";
  import * as Dialog from "$lib/components/ui/dialog";
  import Card from "$lib/components/ui/card/card.svelte";
  import { IconFileArrowRight } from "@tabler/icons-svelte";

  export let responses: any[];
  export let questions: any[];
  export let form: any;
  export let respondentMap: any;

  let exportDialogOpen = false;
  let selectedColumns: Record<string, boolean> = {};

  $: {
    const initialColumns = {
      responseId: false,
      respondent: true,
      status: false,
      started: false,
      submitted: true,
      ...Object.fromEntries(questions.map(q => [q.id, true]))
    };
    
    const newSelectedColumns = { ...initialColumns };
    for (const q of questions) {
      if (!(q.id in selectedColumns)) {
        newSelectedColumns[q.id] = true;
      }
    }
    selectedColumns = newSelectedColumns;
  }

  function decodeBase64Answer(encodedValue: string): string {
    try {
      if (!encodedValue || typeof encodedValue !== 'string') return '';
      const decoded = atob(encodedValue);
      try {
        const parsed = JSON.parse(decoded);
        if (Array.isArray(parsed)) {
          return parsed.join(', ');
        } else if (typeof parsed === 'object' && parsed !== null) {
          return JSON.stringify(parsed);
        }
        return String(parsed);
      } catch {
        return decoded;
      }
    } catch (e) {
      console.warn('Failed to decode base64 answer:', e);
      return String(encodedValue);
    }
  }

  function formatTimestamp(timestamp: string): string {
    if (!timestamp) return '';
    try {
      const date = new Date(timestamp);
      if (isNaN(date.getTime())) return timestamp;
      
      const year = date.getFullYear();
      const month = String(date.getMonth() + 1).padStart(2, '0');
      const day = String(date.getDate()).padStart(2, '0');
      const hours = String(date.getHours()).padStart(2, '0');
      const minutes = String(date.getMinutes()).padStart(2, '0');
      const seconds = String(date.getSeconds()).padStart(2, '0');
      
      return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
    } catch (e) {
      console.warn('Failed to format timestamp:', e);
      return timestamp;
    }
  }

  async function fetchAllAnswers() {
    return Promise.all(
      responses.map(async (response: any) => {
        try {
          const res = await fetch(`/api/forms/${form.id}/responses/${response.id}/answers`, {
            credentials: 'include'
          });
          if (res.ok) {
            const data = await res.json();
            return { responseId: response.id, answers: data.data || [] };
          }
        } catch (e) {
          console.warn(`Failed to fetch answers for response ${response.id}:`, e);
        }
        return { responseId: response.id, answers: [] };
      })
    );
  }

  async function exportToCSV() {
    if (!responses.length) return;
    
    try {
      const allAnswers = await fetchAllAnswers();
      const answersMap = new Map(allAnswers.map(item => [item.responseId, item.answers]));
      
      const headers = [];
      if (selectedColumns.responseId) headers.push('Response ID');
      if (selectedColumns.respondent) headers.push('Respondent');
      if (selectedColumns.status) headers.push('Status');
      if (selectedColumns.started) headers.push('Started');
      if (selectedColumns.submitted) headers.push('Submitted');
      
      for (const q of questions) {
        if (selectedColumns[q.id]) {
          headers.push(q.title || q.id);
        }
      }
      
      const csvRows = [headers.join(',')];
      
      for (const response of responses) {
        const row = [];
        const answers = answersMap.get(response.id) || [];
        const answerMap = new Map(answers.map((a: any) => [a.question, a.value]));
        
        if (selectedColumns.responseId) row.push(`"${response.id}"`);
        if (selectedColumns.respondent) row.push(`"${(respondentMap?.[response.respondent]?.name) || response.respondent || 'Anonymous'}"`);
        if (selectedColumns.status) row.push(`"${response.status || 'Unknown'}"`);
        if (selectedColumns.started) row.push(`"${formatTimestamp(response.started || '')}"`);
        if (selectedColumns.submitted) row.push(`"${formatTimestamp(response.submitted || '')}"`);
        
        for (const q of questions) {
          if (selectedColumns[q.id]) {
            const answer = answerMap.get(q.id) || '';
            let displayValue = answer ? decodeBase64Answer(String(answer)) : '';
            row.push(`"${displayValue.replace(/"/g, '""')}"`);
          }
        }
        
        csvRows.push(row.join(','));
      }
      
      const csvContent = csvRows.join('\n');
      const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
      const url = URL.createObjectURL(blob);
      
      const a = document.createElement('a');
      a.href = url;
      a.download = `${form?.title || 'form'}-responses.csv`;
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      URL.revokeObjectURL(url);
      
      exportDialogOpen = false;
    } catch (error) {
      console.error('Export failed:', error);
      alert('Failed to export responses. Please try again.');
    }
  }
</script>

<button on:click={() => exportDialogOpen = true} class="w-full h-full">
    <Card class="flex items-center  p-6 min-h-[160px] hover:bg-muted/50 transition-colors cursor-pointer h-full flex-row">
        <div class="bg-muted p-4 rounded-md">
            <IconFileArrowRight class="h-10 w-10 text-muted-foreground" />
        </div>
        <span class="text-xl font-semibold">Export to CSV</span>
    </Card>
</button>

<Dialog.Root bind:open={exportDialogOpen}>
  <Dialog.Content class="sm:max-w-[500px]">
    <Dialog.Header>
      <Dialog.Title>Export Responses to CSV</Dialog.Title>
    </Dialog.Header>
    
    <div class="space-y-4 max-h-[500px] overflow-y-auto p-1">
      <div class="space-y-3">
        <h4 class="text-sm font-medium">Metadata</h4>
        <div class="space-y-2">
          <div class="flex items-center space-x-2">
            <Checkbox id="responseId" bind:checked={selectedColumns.responseId} />
            <label for="responseId" class="text-sm">Response ID</label>
          </div>
          <div class="flex items-center space-x-2">
            <Checkbox id="respondent" bind:checked={selectedColumns.respondent} />
            <label for="respondent" class="text-sm">Respondent</label>
          </div>
          <div class="flex items-center space-x-2">
            <Checkbox id="status" bind:checked={selectedColumns.status} />
            <label for="status" class="text-sm">Status</label>
          </div>
          <div class="flex items-center space-x-2">
            <Checkbox id="started" bind:checked={selectedColumns.started} />
            <label for="started" class="text-sm">Started</label>
          </div>
          <div class="flex items-center space-x-2">
            <Checkbox id="submitted" bind:checked={selectedColumns.submitted} />
            <label for="submitted" class="text-sm">Submitted</label>
          </div>
        </div>
      </div>
      
      {#if questions.length > 0}
        <div class="space-y-3">
          <h4 class="text-sm font-medium">Questions</h4>
          <div class="space-y-2">
            {#each questions as question}
              <div class="flex items-center space-x-2">
                <Checkbox id={question.id} bind:checked={selectedColumns[question.id]} />
                <label for={question.id} class="text-sm truncate" title={question.title || question.id}>
                  {question.title || question.id}
                </label>
              </div>
            {/each}
          </div>
        </div>
      {/if}
    </div>
    
    <Dialog.Footer>
      <Button variant="outline" onclick={() => exportDialogOpen = false}>
        Cancel
      </Button>
      <Button onclick={exportToCSV}>
        Export CSV
      </Button>
    </Dialog.Footer>
  </Dialog.Content>
</Dialog.Root>
