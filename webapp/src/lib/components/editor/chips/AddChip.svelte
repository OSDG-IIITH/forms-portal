<script lang="ts">
import PlusIcon from "@lucide/svelte/icons/plus";
import { tick, createEventDispatcher } from "svelte";
import * as Command from "$lib/components/ui/command";
import * as Popover from "$lib/components/ui/popover";
import { Button } from "$lib/components/ui/button";
import { cn } from "$lib/utils.js";

const dispatch = createEventDispatcher();
const { available, onAdd, totalValidations } = $props<{ available: Array<{ key: string; label: string; icon: any }>, onAdd: (key: string) => void, totalValidations?: number }>();

let open = $state(false);
let value = $state("");
let triggerRef = $state<HTMLButtonElement>(null!);

function closeAndFocusTrigger() {
  open = false;
  tick().then(() => {
    triggerRef?.focus();
  });
}

function handleSelect(key: string) {
  value = key;
  onAdd?.(key);
  closeAndFocusTrigger();
  value = "";
  tick().then(() => {
    dispatch("focusValidationInput", { key });
  });
}
</script>

<Popover.Root bind:open>
  <Popover.Trigger bind:ref={triggerRef} class="outline-none ring-0 focus:outline-none focus-visible:outline-none focus:ring-0">
    <Button
      variant="outline"
      class="rounded-md border hover:bg-muted/60 transition-all h-8 py-1.5 flex items-center justify-center text-sm min-w-0 outline-none ring-0 focus:outline-none focus-visible:outline-none focus:ring-0 gap-2"
      role="combobox"
      aria-expanded={open}
      aria-label="Add validation"
    >
      <PlusIcon class="size-4" />
      {#if totalValidations !== undefined && available.length === totalValidations}
        <span>Add</span>
      {/if}
    </Button>
  </Popover.Trigger>
  <Popover.Content class="w-[220px] p-0">
    <Command.Root bind:value>
      <Command.Input placeholder="Search validation..." />
      <Command.List>
        <Command.Empty>No validation found.</Command.Empty>
        <Command.Group>
          {#each available as v}
            <Command.Item
              value={v.key}
              onSelect={() => handleSelect(v.key)}
            >
              {#if v.icon}
                <v.icon class="size-4 mr-2" />
              {/if}
              {v.label}
            </Command.Item>
          {/each}
        </Command.Group>
      </Command.List>
    </Command.Root>
  </Popover.Content>
</Popover.Root>
