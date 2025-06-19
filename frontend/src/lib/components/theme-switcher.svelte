<script lang="ts">
  import { onMount } from 'svelte';
  import {
    Popover,
    PopoverTrigger,
    PopoverContent
  } from '$lib/components/ui/popover';
  import { Button } from '$lib/components/ui/button';

  import { IconPalette } from '@tabler/icons-svelte';
  import { CheckIcon } from '@lucide/svelte';
  import { theme, type Theme } from '$lib/stores/theme';

  interface ThemeOption {
    value: Theme;
    label: string;
    isDark: boolean;
    accent: string;
  }

  const themeOptions: ThemeOption[] = [
    {
      value: 'milk',
      label: 'Milk',
      isDark: false,
      accent: 'oklch(0.9 0 0)'
    },
    {
      value: 'lilac',
      label: 'Lilac',
      isDark: false,
      accent: 'oklch(0.8 0.15 310)'
    },
    {
      value: 'mint',
      label: 'Mint',
      isDark: false,
      accent: 'oklch(0.65 0.12 140)'
    },
    {
      value: 'moss',
      label: 'Moss',
      isDark: true,
      accent: 'oklch(0.5 0.15 145)'
    },
    {
      value: 'amethyst',
      label: 'Amethyst',
      isDark: true,
      accent: 'oklch(0.4 0.2 285)'
    },
    {
      value: 'coal',
      label: 'Coal',
      isDark: true,
      accent: 'oklch(0.3 0 0)'
    }
  ];

  let selectedTheme: Theme = $state('milk');
  let isOpen = $state(false);

  onMount(() => {
    const unsubscribe = theme.subscribe((value) => {
      selectedTheme = value;
    });

    return unsubscribe;
  });

  function handleThemeSelect(themeValue: Theme) {
    selectedTheme = themeValue;
    theme.set(themeValue);
    isOpen = false;
  }
</script>

<Popover bind:open={isOpen}>
  <PopoverTrigger>
    <Button
      variant="outline"
      size="sm"
      class="h-9 px-3 rounded-lg border hover:bg-accent/50 hover:border-primary/30 transition-all duration-200 hover:scale-105"
      aria-label="Change theme"
    >
      <IconPalette class="w-4 h-4" />
    </Button>
  </PopoverTrigger>

  <PopoverContent class="w-32 p-2" align="end">
    <div class="space-y-1">
      {#each themeOptions as option (option.value)}
        <button
          onclick={() => handleThemeSelect(option.value)}
          class="w-full flex items-center justify-between px-2 py-1.5 rounded-sm hover:bg-accent/50 transition-colors group"
        >
          <div class="flex items-center gap-2">
            <!--
            TODO: accent circles look ugly
            <div 
              class="w-3 h-3 rounded-full transition-all duration-200 group-hover:scale-110"
              style="background-color: {option.accent}"
            ></div>
            -->
            <span class="text-xs font-medium">{option.label}</span>
          </div>

          {#if selectedTheme === option.value}
            <CheckIcon class="w-3 h-3 text-primary" />
          {/if}
        </button>
      {/each}
    </div>
  </PopoverContent>
</Popover>
