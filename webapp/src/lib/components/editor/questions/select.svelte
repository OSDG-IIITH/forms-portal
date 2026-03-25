<script lang="ts">
	import Options from './options.svelte';

	let { questionId }: { questionId: string } = $props();
</script>

<Options
	{questionId}
	optionsLabel="Options (dropdown)"
	requiredContainerClass="flex items-center gap-2 pt-4 border-t"
/>