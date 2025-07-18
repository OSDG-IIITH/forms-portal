<script lang="ts">
  interface $$Props {
    date?: string;
  }
  
  let { date }: $$Props = $props();
  
  const formatDateTime = (dateString: string) => {
    const dt = new Date(dateString);
    return dt.toLocaleString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
      hour: 'numeric',
      minute: '2-digit',
      hour12: true
    });
  };
</script>

<div class="text-sm text-muted-foreground">
  {date ? formatDateTime(date) : "-"}
</div>
