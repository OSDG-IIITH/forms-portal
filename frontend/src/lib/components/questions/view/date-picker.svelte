<script lang="ts">
  import CalendarIcon from "@lucide/svelte/icons/calendar";
  import {
    type DateValue,
    DateFormatter,
    getLocalTimeZone,
    parseDate
  } from "@internationalized/date";
  import { cn } from "$lib/utils.js";
  import { Button } from "$lib/components/ui/button/index.js";
  import { Calendar } from "$lib/components/ui/calendar/index.js";
  import * as Popover from "$lib/components/ui/popover/index.js";

  interface Props {
    question: {
      uid: string;
      text: string;
      required: boolean;
    };
    value?: string;
  }

  let { question, value = $bindable('') }: Props = $props();

  const df = new DateFormatter("en-US", {
    dateStyle: "long",
  });

  let dateValue = $state<DateValue | undefined>();

  $effect(() => {
    if (value?.trim()) {
      try {
        dateValue = parseDate(value);
      } catch {
        dateValue = undefined;
      }
    } else {
      dateValue = undefined;
    }
  });

  $effect(() => {
    if (dateValue) {
      value = dateValue.toString();
    } else if (dateValue === undefined && value) {
      value = '';
    }
  });
</script>

<Popover.Root>
  <Popover.Trigger>
    {#snippet child({ props })}
      <Button
        variant="outline"
        class={cn(
          "w-full max-w-sm justify-start text-left font-normal",
          !dateValue && "text-muted-foreground"
        )}
        {...props}
      >
        <CalendarIcon class="mr-2 size-4" />
        {dateValue ? df.format(dateValue.toDate(getLocalTimeZone())) : "Select a date"}
      </Button>
    {/snippet}
  </Popover.Trigger>
  <Popover.Content class="w-auto p-0">
    <Calendar 
      bind:value={dateValue}
      type="single" 
      initialFocus 
    />
  </Popover.Content>
</Popover.Root>
