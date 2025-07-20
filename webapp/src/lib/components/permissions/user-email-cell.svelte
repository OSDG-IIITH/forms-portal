<script lang="ts">
interface Props {
  userEmail?: string;
}
let { userEmail }: Props = $props();
</script>

<div class="text-muted-foreground">
  {userEmail}
</div>
