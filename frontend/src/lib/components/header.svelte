<script>
    import ThemeSwitcher from "./theme-switcher.svelte";

</script>
  <header class="border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60 fixed top-0 left-0 right-0 z-50">
    <div class="container mx-auto px-6 py-4">
      <div class="flex items-center justify-between">
        <h1 class="text-3xl tracking-tight">
          forms <span class="font-semibold text-accent-foreground">iiit</span>
        </h1>
        <ThemeSwitcher />
      </div>
    </div>
  </header>