<script>
    import { page } from '$app/stores';

    const notFoundMessage = "Sorry, that page could not be found.";
    const genericErrorMessage = "Something went wrong on our end. Please try again later.";
    const displayMessage = ($page.status === 404) ? notFoundMessage : genericErrorMessage;
</script>

<div class="min-h-screen flex flex-col justify-center items-center text-center">
    <h1 class="text-3xl font-bold mb-4">{$page.status}: {$page.error?.message || displayMessage}</h1>
    
    <p class="mb-6">
        {#if $page.status === 404}
            It looks like the page you are looking for has been moved or does not exist.
        {:else}
            Please try again later. If the problem persists, contact us at osdg@iiit.ac.in
        {/if}
    </p>

    <a href="/" class="hover:underline">Go home.</a>
</div>