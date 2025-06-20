<script lang="ts">
  import CreateFormCard from '$lib/components/dashboard/create-form-card.svelte';
  import FormTemplatesCard from '$lib/components/dashboard/form-templates-card.svelte';
  import RecentlyCreatedCard from '$lib/components/dashboard/recently-created-card.svelte';
  import RecentResponsesCard from '$lib/components/dashboard/recent-responses-card.svelte';
</script>

<div class="min-h-screen bg-background">
  <main class="container mx-auto max-w-7xl px-6 py-8 pt-24">
    <div class="space-y-6">
      <div class="grid grid-cols-1 lg:grid-cols-20 gap-6">
        <div class="lg:col-span-5">
          <CreateFormCard />
        </div>
        <div class="hidden md:block lg:col-span-15">
          <FormTemplatesCard />
        </div>
        <div class="lg:col-span-8">
          <RecentlyCreatedCard />
        </div>
        <div class="lg:col-span-12">
          <RecentResponsesCard />
        </div>
      </div>
    </div>
  </main>
</div>
