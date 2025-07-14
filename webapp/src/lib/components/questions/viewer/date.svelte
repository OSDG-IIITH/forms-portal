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

export let question: {
  id: string;
  title: string;
  required: boolean;
  placeholder?: string;
};
export let value: string = "";
export let placeholder: string = "Select a date";

const df = new DateFormatter("en-US", {
  dateStyle: "long",
});

let dateValue: DateValue | undefined = value ? parseDate(value) : undefined;
let lastValue = value;

$: {
  if (value !== lastValue) {
    if (value) {
      if (value !== dateValue?.toString()) {
        try {
          dateValue = parseDate(value);
        } catch {
          dateValue = undefined;
        }
      }
    } else {
      dateValue = undefined;
    }
    lastValue = value;
  }
  else if (dateValue?.toString() !== value) {
    value = dateValue ? dateValue.toString() : "";
    lastValue = value;
  }
}
</script>

<Popover.Root>
  <Popover.Trigger>
    <Button
      variant="outline"
      class={cn(
        "w-full max-w-sm justify-start text-left font-normal",
        !dateValue && "text-muted-foreground"
      )}
      type="button"
    >
      <CalendarIcon class="mr-2 size-4" />
      {dateValue ? df.format(dateValue.toDate(getLocalTimeZone())) : (question.placeholder || placeholder)}
    </Button>
  </Popover.Trigger>
  <Popover.Content class="w-auto p-0">
    <Calendar
      bind:value={dateValue}
      type="single"
      initialFocus
    />
  </Popover.Content>
</Popover.Root>
