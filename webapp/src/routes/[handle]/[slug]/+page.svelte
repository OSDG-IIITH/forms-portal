<script lang="ts">
  import FormViewer from '$lib/components/FormViewer.svelte';
  export let data: { form: any };
</script>

<FormViewer form={data.form} />
