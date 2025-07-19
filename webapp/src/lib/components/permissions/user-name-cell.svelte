<script lang="ts">
interface Props {
  userName?: string;
  userEmail?: string;
  userId: string;
}

let { userName, userEmail, userId }: Props = $props();
</script>

<div class="font-medium">
  {userName || userEmail || userId}
</div>
