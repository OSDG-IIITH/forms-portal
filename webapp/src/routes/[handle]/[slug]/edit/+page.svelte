<script lang="ts">
  import FormEditor from '$lib/components/editor/FormEditor.svelte';
  
  const { data }: { data: { form: any } } = $props();
</script>

<FormEditor form={data.form} />