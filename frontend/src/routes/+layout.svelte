<script lang="ts">
    import Header from '$lib/components/header.svelte';
	import '../app.css';
	import { ModeWatcher } from 'mode-watcher';
	
	let { children } = $props();
</script>

<ModeWatcher />
<Header />
{@render children()}
