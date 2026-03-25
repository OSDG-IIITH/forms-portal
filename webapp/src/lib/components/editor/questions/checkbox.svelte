<script lang="ts">
  import Options from './options.svelte';

  let { questionId }: { questionId: string } = $props();
</script>

<Options {questionId} optionsLabel="Options" />
