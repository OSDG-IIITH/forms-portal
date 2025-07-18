<script lang="ts">
  import Card from '$lib/components/ui/card/card.svelte';
  import { Button } from '$lib/components/ui/button';
  import { IconPlus } from '@tabler/icons-svelte';
  export let questionTypeButtons;
  export let addQuestion;
  export let questionsLength;
</script>

{#if questionsLength > 0}
  <div class="border-2 border-dashed border-muted-foreground/25 rounded-lg p-6">
    <div class="text-center space-y-4">
      <div>
        <h3 class="font-medium mb-1">Add Question</h3>
        <p class="text-sm text-muted-foreground">Choose a question type</p>
      </div>
      <div class="flex flex-wrap justify-center gap-2">
        {#each questionTypeButtons as { type, icon: Icon, label }}
          <Button
            variant="outline"
            size="sm"
            onclick={() => addQuestion(type)}
            class="flex items-center gap-2"
          >
            <Icon class="h-4 w-4" />
            {label}
          </Button>
        {/each}
      </div>
    </div>
  </div>
{:else}
  <Card class="flex flex-col items-center justify-center py-12 text-center gap-0">
    <div class="w-12 h-12 bg-muted rounded-md flex items-center justify-center mb-4">
      <IconPlus class="h-6 w-6 text-muted-foreground" />
    </div>
    <h3 class="text-lg font-medium mb-2">Create your first question</h3>
    <p class="text-muted-foreground mb-6 max-w-sm">
      Get started by adding questions to your form
    </p>
    <div class="flex flex-wrap justify-center gap-3 max-w-xl">
      {#each questionTypeButtons as { type, icon: Icon, label }}
        <Button
          variant="outline"
          size="sm"
          onclick={() => addQuestion(type)}
          class="flex items-center gap-2"
        >
          <Icon class="h-4 w-4" />
          {label}
        </Button>
      {/each}
    </div>
  </Card>
{/if}
