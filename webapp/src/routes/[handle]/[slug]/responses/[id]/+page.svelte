<script lang="ts">
import { Card, CardContent } from '$lib/components/ui/card';
import type { PageData } from './$types';
export let data: PageData;
let { response } = data;
</script>

<div class="max-w-2xl mx-auto py-8">
  <h2 class="text-2xl font-semibold mb-6">Response Details</h2>
  <Card>
    <CardContent class="overflow-x-auto">
      <pre class="text-sm whitespace-pre-wrap break-words">{JSON.stringify(response, null, 2)}</pre>
    </CardContent>
  </Card>
</div>
