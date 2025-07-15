<script lang="ts">
  import CreateFormCard from '$lib/components/dashboard/create-form-card.svelte';
  import FormTemplatesCard from '$lib/components/dashboard/form-templates-card.svelte';
  import FormsGrid from '$lib/components/dashboard/forms-grid.svelte';
</script>

<div class="min-h-screen flex flex-col bg-background">
  <main class="flex-1 flex flex-col pt-16">
    <section class="bg-background">
      <div class="container mx-auto max-w-7xl px-6 py-6">
        <div class="grid grid-cols-1 lg:grid-cols-20 gap-6">
          <div class="lg:col-span-5">
            <CreateFormCard />
          </div>
          <div class="hidden md:block lg:col-span-15">
            <FormTemplatesCard />
          </div>
        </div>
      </div>
    </section>
    <section class="bg-card border-t flex-1 flex flex-col min-h-0">
      <div class="container mx-auto max-w-7xl px-6 py-6 flex-1 flex flex-col min-h-0">
        <FormsGrid />
      </div>
    </section>
  </main>
</div>
