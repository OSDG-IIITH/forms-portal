<script lang="ts">
  import { Card, CardContent, CardHeader, CardTitle, CardDescription, CardFooter } from '$lib/components/ui/card';
  import { Button } from '$lib/components/ui/button';
  import { IconMessageCircle } from '@tabler/icons-svelte';
  import mockResponses from '$lib/data/recent-responses.json';
  import * as Pagination from '$lib/components/ui/pagination/index.js';
  import { IconChevronLeft, IconChevronRight } from '@tabler/icons-svelte';
  import { MediaQuery } from 'svelte/reactivity';
  
  const isDesktop = new MediaQuery("(min-width: 768px)");
  
  const count = mockResponses.length;
  // TODO: remove these if not needed
  const perPage = $derived(isDesktop.current ? 5 : 5);
  const siblingCount = $derived(isDesktop.current ? 0 : 0);
  
  let currentPageResponses = $state<typeof mockResponses>([]);
  
  function updateVisibleResponses(page: number) {
    const startIndex = (page - 1) * perPage;
    const endIndex = Math.min(startIndex + perPage, mockResponses.length);
    currentPageResponses = mockResponses.slice(startIndex, endIndex);
  }
  
  updateVisibleResponses(1);
</script>

<Card class="h-full flex flex-col gap-4">
  <CardHeader>
    <CardTitle>Recent Responses</CardTitle>
  </CardHeader>
  <CardContent class="flex-1">
    <div class="space-y-4">
      {#if currentPageResponses?.length}
        {#each currentPageResponses as response}
          <div class="flex items-center gap-4">
            <div class="bg-muted p-3 rounded-md">
              <IconMessageCircle class="h-5 w-5 text-muted-foreground" />
            </div>
            <div class="flex-1">
              <p class="font-medium">
                <!--svelte-ignore a11y_invalid_attribute-->
                New response for <a href="" class="text-primary hover:underline">{response.formTitle}</a>
              </p>
              <p class="text-sm text-muted-foreground">from {response.respondent} &bull; {response.receivedAt}</p>
            </div>
            <Button variant="outline" size="sm">View</Button>
          </div>
        {/each}
      {:else}
        <div class="text-center py-4 text-muted-foreground">No responses found</div>
      {/if}
    </div>
  </CardContent>
  <CardFooter>
    <div class="flex justify-center w-full">
      <Pagination.Root {count} {perPage} {siblingCount} onPageChange={updateVisibleResponses}>
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