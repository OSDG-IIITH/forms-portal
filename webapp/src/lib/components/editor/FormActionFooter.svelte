<script lang="ts">
  import { Card, CardContent } from '$lib/components/ui/card';
  import { Button } from '$lib/components/ui/button';
  export let questionsLength;
  export let isSaving;
  export let saveForm;
</script>

<Card class="py-0">
  <CardContent class="py-4">
    <div class="flex justify-between items-center">
      <div class="flex items-center gap-2">
        <span class="text-sm text-muted-foreground">
          {questionsLength} question{questionsLength !== 1 ? 's' : ''}
        </span>
      </div>
      <div class="flex gap-4">
        <Button
          onclick={saveForm}
          disabled={isSaving}
          size="sm"
        >
          {#if isSaving}
            <div class="w-4 h-4 border-2 border-current border-t-transparent rounded-full animate-spin mr-2"></div>
            Saving...
          {:else}
            Save Form
          {/if}
        </Button>
      </div>
    </div>
  </CardContent>
</Card>
