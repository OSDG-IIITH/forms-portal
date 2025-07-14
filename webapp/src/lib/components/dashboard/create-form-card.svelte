<script lang="ts">
  import { Card, CardContent } from '$lib/components/ui/card';
  import { IconPlus } from '@tabler/icons-svelte';
</script>

<a href="/create/demo" class="block h-full group">
  <Card class="h-full hover:border-primary/50 hover:bg-accent/20 transition-all duration-200 shadow-none">
    <CardContent class="flex items-center justify-center h-full p-6">
      <IconPlus class="h-12 w-12 text-muted-foreground group-hover:text-primary transition-colors" />
    </CardContent>
  </Card>
</a>
