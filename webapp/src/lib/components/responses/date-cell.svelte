<script lang="ts">
  interface $$Props {
    date?: string;
  }
  
  let { date }: $$Props = $props();
</script>

<div class="text-sm text-muted-foreground">
  {date ? new Date(date).toLocaleDateString() : "-"}
</div>
