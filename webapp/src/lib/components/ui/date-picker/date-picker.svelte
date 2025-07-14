<script lang="ts">
import { DatePicker } from "bits-ui";
import { CalendarIcon, ChevronLeft, ChevronRight } from "@lucide/svelte";
import { parseDate, type DateValue } from "@internationalized/date";
import { cn } from "$lib/utils.js";

interface Props {
  value?: string;
  placeholder?: string;
  disabled?: boolean;
  class?: string;
}

let {
  value = $bindable(),
  placeholder = "Select a date",
  disabled = false,
  class: className
}: Props = $props();

let dateValue: DateValue | undefined = $state(value ? (() => {
  try {
    return parseDate(value);
  } catch {
    return undefined;
  }
})() : undefined);

$effect(() => {
  if (value === "" && dateValue !== undefined) {
    dateValue = undefined;
  }
});

function handleValueChange(newDateValue: DateValue | undefined) {
  dateValue = newDateValue;
  value = newDateValue ? newDateValue.toString() : "";
}
</script>

<DatePicker.Root value={dateValue} onValueChange={handleValueChange} weekdayFormat="short" fixedWeeks={true} {disabled} locale="en-GB" initialFocus={false}>
  <div class={cn("flex w-full max-w-sm flex-col gap-1.5", className)}>
    <DatePicker.Input
      class={cn(
        "flex h-9 w-full items-center rounded-md border border-input bg-background dark:bg-input/30 px-3 py-1 text-sm shadow-xs",
        disabled && "cursor-not-allowed opacity-50"
      )}
    >
      {#snippet children({ segments })}
        {#each segments as { part, value }, i (part + i)}
          <div class="inline-block select-none">
            {#if part === "literal"}
              <DatePicker.Segment {part} class="text-muted-foreground p-1">
                {value}
              </DatePicker.Segment>
            {:else}
              <DatePicker.Segment
                {part}
                class="rounded-sm hover:bg-accent aria-[valuetext=Empty]:text-muted-foreground px-1 py-1 focus:outline-none focus:ring-1 focus:ring-border"
              >
                {value}
              </DatePicker.Segment>
            {/if}
          </div>
        {/each}
        <DatePicker.Trigger
          class="text-muted-foreground hover:bg-accent hover:text-accent-foreground ml-auto inline-flex size-6 items-center justify-center rounded-sm transition-all focus:outline-none"
          tabindex={-1}
        >
          <CalendarIcon class="size-4" />
        </DatePicker.Trigger>
      {/snippet}
    </DatePicker.Input>
    <DatePicker.Content sideOffset={6} class="z-50">
      <DatePicker.Calendar
        class="rounded-md border bg-popover p-3 shadow-md"
      >
        {#snippet children({ months, weekdays })}
          <DatePicker.Header class="flex items-center justify-between">
            <DatePicker.PrevButton
              class="rounded-sm bg-transparent hover:bg-accent inline-flex size-8 items-center justify-center transition-all active:scale-95"
            >
              <ChevronLeft class="size-4" />
            </DatePicker.PrevButton>
            <DatePicker.Heading class="text-sm font-medium" />
            <DatePicker.NextButton
              class="rounded-sm bg-transparent hover:bg-accent inline-flex size-8 items-center justify-center transition-all active:scale-95"
            >
              <ChevronRight class="size-4" />
            </DatePicker.NextButton>
          </DatePicker.Header>
          <div class="flex flex-col space-y-4 pt-4 sm:flex-row sm:space-x-4 sm:space-y-0">
            {#each months as month (month.value)}
              <DatePicker.Grid class="w-full border-collapse select-none space-y-1">
                <DatePicker.GridHead>
                  <DatePicker.GridRow class="mb-1 flex w-full justify-between">
                    {#each weekdays as day (day)}
                      <DatePicker.HeadCell
                        class="text-muted-foreground w-8 rounded-md text-xs font-normal"
                      >
                        <div>{day.slice(0, 2)}</div>
                      </DatePicker.HeadCell>
                    {/each}
                  </DatePicker.GridRow>
                </DatePicker.GridHead>
                <DatePicker.GridBody>
                  {#each month.weeks as weekDates (weekDates)}
                    <DatePicker.GridRow class="flex w-full">
                      {#each weekDates as date (date)}
                        <DatePicker.Cell
                          {date}
                          month={month.value}
                          class="relative size-8 p-0 text-center text-sm"
                        >
                          <DatePicker.Day
                            class="rounded-sm text-foreground hover:bg-accent hover:text-accent-foreground data-selected:bg-primary data-selected:text-primary-foreground data-disabled:text-muted-foreground/50 data-disabled:pointer-events-none data-outside-month:pointer-events-none data-selected:font-medium data-unavailable:text-destructive data-unavailable:line-through relative inline-flex size-8 items-center justify-center whitespace-nowrap border border-transparent bg-transparent p-0 text-sm font-normal transition-all"
                          >
                            <div
                              class="bg-primary group-data-selected:bg-primary-foreground data-[today]:block absolute top-1 hidden size-1 rounded-full transition-all"
                            ></div>
                            {date.day}
                          </DatePicker.Day>
                        </DatePicker.Cell>
                      {/each}
                    </DatePicker.GridRow>
                  {/each}
                </DatePicker.GridBody>
              </DatePicker.Grid>
            {/each}
          </div>
        {/snippet}
      </DatePicker.Calendar>
    </DatePicker.Content>
  </div>
</DatePicker.Root>
