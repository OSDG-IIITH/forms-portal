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


    <div class="grid grid-cols-1 sm:grid-cols-3 lg:grid-cols-5 gap-4">
      {#each Array($num) as _}
        <div class="bg-card space-y-2 rounded-lg border p-4 hover:bg-muted/30 transition-colors cursor-pointer">
          <Skeleton class="h-38 rounded-lg" />
          <Skeleton class="h-4 w-3/4" />
        </div>
      {/each}
    </div>