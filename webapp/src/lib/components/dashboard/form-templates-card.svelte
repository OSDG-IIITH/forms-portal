<script lang="ts">
  import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '$lib/components/ui/card';
  import { Skeleton } from '$lib/components/ui/skeleton';
  import { onMount } from 'svelte';
  import { writable } from 'svelte/store';
  import { browser } from '$app/environment';

  const num = writable(5);

  function updateTemplateCount() {
    if (browser) {
      const newCount = window.innerWidth >= 1024 ? 5 : 3;
      $num = newCount;
    }
  }

  if (browser) {
    updateTemplateCount();
  }

  onMount(() => {
    updateTemplateCount();
    
    if (browser && window.ResizeObserver) {
      const resizeObserver = new ResizeObserver(() => {
        updateTemplateCount();
      });
      
      resizeObserver.observe(document.documentElement);
      
      return () => {
        resizeObserver.disconnect();
      };
    } else {
      window.addEventListener('resize', updateTemplateCount);
      
      return () => {
        window.removeEventListener('resize', updateTemplateCount);
      };
    }
  });
</script>

<Card class="h-full gap-2 shadow-none">
  <CardHeader>
    <CardTitle>Form Templates</CardTitle>
  </CardHeader>
  <CardContent>
    <div class="grid grid-cols-1 sm:grid-cols-3 lg:grid-cols-5 gap-4">
      {#each Array($num) as _}
        <div class="space-y-2 rounded-lg border p-4 hover:border-primary/25 hover:bg-accent/20 transition-colors cursor-pointer">
          <Skeleton class="h-32 rounded-lg" />
          <Skeleton class="h-4 w-3/4" />
        </div>
      {/each}
    </div>
  </CardContent>
</Card>
