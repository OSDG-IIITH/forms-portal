<script lang="ts">
  export let value: string = '';
</script>

<div class="text-sm text-muted-foreground">
  {value ? new Date(value).toLocaleDateString() : '-'}
</div>