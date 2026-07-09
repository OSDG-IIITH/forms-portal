<script lang="ts">
  import { Button } from "$lib/components/ui/button";
</script>

<div class="min-h-[calc(100vh-4rem)] flex items-center justify-center px-4 pt-16">
  <div class="max-w-md w-full text-center space-y-6">
    <div>
      <h1 class="text-2xl font-semibold tracking-tight">You have been logged out</h1>
    </div>
    <div class="pt-2">
      <Button href="/api/auth/login" class="w-full sm:w-auto px-8">
        Log in again
      </Button>
    </div>
  </div>
</div>
