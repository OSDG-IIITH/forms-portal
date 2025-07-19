<script lang="ts">
interface Props {
  roles: Array<"view" | "respond" | "comment" | "analyze" | "edit" | "manage">;
}

let { roles }: Props = $props();

const roleLabels = {
  view: "View",
  respond: "Respond", 
  comment: "Comment",
  analyze: "Analyze",
  edit: "Edit",
  manage: "Manage"
};
</script>

<div class="flex flex-wrap gap-1">
  {#each roles as role}
    <span class="inline-flex items-center rounded-md bg-muted px-2 py-1 text-xs font-medium text-muted-foreground transition-colors cursor-pointer hover:bg-primary hover:text-primary-foreground focus:bg-primary focus:text-primary-foreground">
      {roleLabels[role]}
    </span>
  {/each}
</div>
