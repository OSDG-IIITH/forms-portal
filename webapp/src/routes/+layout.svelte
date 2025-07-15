<script lang="ts">
    import Header from '$lib/components/header.svelte';
	import Footer from '$lib/components/footer.svelte';
	import '../app.css';
    import { Toaster } from '$lib/components/ui/sonner';
	
	let { children } = $props();
</script>

<Toaster 
  position="bottom-right"
/>

<Header />
{@render children()}
<Footer />