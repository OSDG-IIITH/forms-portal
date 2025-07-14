<script lang="ts">
  import { Card, CardContent, CardHeader, CardTitle, CardDescription, CardFooter } from '$lib/components/ui/card';
  import { Button } from '$lib/components/ui/button';
  import { IconFileText, IconChevronLeft, IconChevronRight } from '@tabler/icons-svelte';
  import recentForms from '$lib/data/recent-forms.json';
  import * as Pagination from '$lib/components/ui/pagination/index.js';
  import { MediaQuery } from 'svelte/reactivity';
  
  const isDesktop = new MediaQuery("(min-width: 768px)");
  
  const count = recentForms.length;
  // TODO: remove these if not needed
  const perPage = $derived(isDesktop.current ? 5 : 5);
  const siblingCount = $derived(isDesktop.current ? 0 : 0);
  
  let currentPageForms = $state<typeof recentForms>([]);
  
  function updateVisibleForms(page: number) {
    const startIndex = (page - 1) * perPage;
    const endIndex = Math.min(startIndex + perPage, recentForms.length);
    currentPageForms = recentForms.slice(startIndex, endIndex);
  }
  
  updateVisibleForms(1);
</script>

<Card class="h-full flex flex-col gap-4 shadow-none">
  <CardHeader>
    <CardTitle>Recently Created</CardTitle>
  </CardHeader>
  <CardContent class="flex-1 mt-0 pt-0">
    <div class="space-y-4">
      {#if currentPageForms?.length}
        {#each currentPageForms as form}
          <div class="flex items-center gap-4">
            <div class="bg-muted p-3 rounded-md">
              <IconFileText class="h-5 w-5 text-muted-foreground" />
            </div>
            <div class="flex-1">
              <!--svelte-ignore a11y_invalid_attribute-->
              <a href="" class="font-medium hover:underline">{form.title}</a>
              <p class="text-sm text-muted-foreground">{form.createdAt}</p>
            </div>
            <Button variant="outline" size="sm">View</Button>
          </div>
        {/each}
      {:else}
        <div class="text-center py-4 text-muted-foreground">No forms found</div>
      {/if}
    </div>
  </CardContent>
  <CardFooter>
    <div class="flex justify-center w-full">
      <Pagination.Root {count} {perPage} {siblingCount} onPageChange={updateVisibleForms}>
        {#snippet children({ pages, currentPage })}
          <Pagination.Content>
            <Pagination.Item>
              <Pagination.PrevButton>
                <IconChevronLeft class="h-4 w-4" />
                <span class="hidden sm:block ml-1">Previous</span>
              </Pagination.PrevButton>
            </Pagination.Item>
            {#each pages as page (page.key)}
              {#if page.type === "ellipsis"}
                <Pagination.Item>
                  <Pagination.Ellipsis />
                </Pagination.Item>
              {:else}
                <Pagination.Item>
                  <Pagination.Link {page} isActive={currentPage === page.value}>
                    {page.value}
                  </Pagination.Link>
                </Pagination.Item>
              {/if}
            {/each}
            <Pagination.Item>
              <Pagination.NextButton>
                <span class="hidden sm:block mr-1">Next</span>
                <IconChevronRight class="h-4 w-4" />
              </Pagination.NextButton>
            </Pagination.Item>
          </Pagination.Content>
        {/snippet}
      </Pagination.Root>
    </div>
  </CardFooter>
</Card>
