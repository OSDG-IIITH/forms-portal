<script lang="ts">
  import { Badge } from "$lib/components/ui/badge/index.js";
  
  interface $$Props {
    icon?: string;
    label: string;
  }
  
  let { icon, label }: $$Props = $props();
</script>

<Badge variant="secondary" class="flex w-fit items-center gap-1">
  {#if icon}
    <span>{icon}</span>
  {/if}
  <span>{label}</span>
</Badge>
