<script lang="ts">
  import { getContext } from 'svelte';
  import type { FormStore } from './form-store.svelte';

  const store: FormStore = getContext('form-store');
  const { formData } = store;
</script>

<div class="space-y-2">
  <input
    type="text"
    class="text-3xl tracking-tight leading-tight font-bold w-full bg-transparent outline-none border-2 border-transparent hover:border-border/75 p-2 px-4 rounded-md mb-1"
    placeholder="Form Title"
    autocomplete="off"
    value={formData.title}
    oninput={(e) => store.updateFormData({ title: e.currentTarget.value })}
  />
  <textarea
    class="w-full bg-transparent text-lg outline-none text-muted-foreground resize-none mb-4 border-2 border-transparent hover:border-border/75 p-2 px-4 rounded-md"
    placeholder="Add a description"
    rows="2"
    maxlength="300"
    autocomplete="off"
    value={formData.description}
    oninput={(e) => store.updateFormData({ description: e.currentTarget.value })}
  ></textarea>
</div>