<script lang="ts">
    import Header from '$lib/components/header.svelte';
	import Footer from '$lib/components/footer.svelte';
	import '../app.css';
	
	let { children } = $props();
</script>

<Header />
{@render children()}
<Footer />