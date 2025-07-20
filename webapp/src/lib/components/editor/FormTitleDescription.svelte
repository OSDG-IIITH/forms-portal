<script lang="ts">
  export let title: string;
  export let description: string;
</script>

<div class="space-y-2">
  <input
    type="text"
    class="text-3xl tracking-tight leading-tight font-bold w-full bg-transparent outline-none border-2 border-transparent hover:border-border/75 p-2 px-4 rounded-md mb-1"
    bind:value={title}
    placeholder="Form Title"
    autocomplete="off"
  />
  <textarea
    class="w-full bg-transparent text-lg outline-none text-muted-foreground resize-none mb-4 border-2 border-transparent hover:border-border/75 p-2 px-4 rounded-md"
    bind:value={description}
    placeholder="Add a description"
    rows="2"
    maxlength="300"
    autocomplete="off"
  ></textarea>
</div>
