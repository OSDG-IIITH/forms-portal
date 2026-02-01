<script lang="ts">
  import { IconFileText } from "@tabler/icons-svelte";
  import { Card, CardContent } from "$lib/components/ui/card";
  import { cn } from "$lib/utils";
  import { onMount } from "svelte";
  import { browser } from "$app/environment";

  export let templates: { slug: string; title: string; description: string }[] =
    [];

  let limit = 5;

  function updateLimit() {
    if (!browser) return;
    const w = window.innerWidth;
    if (w >= 1536) limit = 5;
    else if (w >= 1280) limit = 4;
    else if (w >= 1024) limit = 3;
    else if (w >= 640) limit = 2;
    else limit = 1;
  }

  onMount(() => {
    updateLimit();
    window.addEventListener("resize", updateLimit);
    return () => window.removeEventListener("resize", updateLimit);
  });
</script>

<div>
  <div
    class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5 gap-4"
  >
    {#each templates.slice(0, limit) as template}
      <a
        href="/create?template={template.slug}"
        class="block h-full no-underline"
      >
        <Card
          class={cn(
            "p-0 bg-muted/10 hover:bg-muted/30 rounded-md shadow-xs border transition-all duration-300 ease-[cubic-bezier(.4,0,.2,1)] h-40 cursor-pointer will-change-transform will-change-shadow hover:shadow-lg hover:-translate-y-0 hover:scale-[1.012] group",
          )}
        >
          <CardContent class="p-0 h-full">
            <div class="h-full flex flex-col">
              <div
                class="h-24 w-full bg-muted rounded-t-md flex items-center justify-center relative"
              >
                <IconFileText class="w-8 h-8 text-muted-foreground" />
              </div>
              <div class="flex-1 px-4 py-4 flex flex-col justify-center">
                <h3 class="font-medium text-sm truncate text-foreground">
                  {template.title}
                </h3>
              </div>
            </div>
          </CardContent>
        </Card>
      </a>
    {/each}

    {#if templates.length === 0}
      <div class="col-span-full text-center text-muted-foreground py-8">
        No templates found.
      </div>
    {/if}
  </div>
</div>
