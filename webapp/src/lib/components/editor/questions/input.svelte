<script lang="ts">
  import Text from './text.svelte';

  let { questionId }: { questionId: string } = $props();
</script>

<Text {questionId} />