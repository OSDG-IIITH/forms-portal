<script lang="ts">
  export let title: string;
  export let description: string;
</script>

<div class="space-y-2">
  <input
    type="text"
    class="text-3xl tracking-tight leading-tight font-bold w-full bg-transparent outline-none"
    bind:value={title}
    placeholder="Form Title"
    autocomplete="off"
  />
  <textarea
    class="w-full bg-transparent text-lg outline-none text-muted-foreground resize-none"
    bind:value={description}
    placeholder="Add a description"
    rows="2"
    maxlength="300"
    autocomplete="off"
  ></textarea>
</div>
