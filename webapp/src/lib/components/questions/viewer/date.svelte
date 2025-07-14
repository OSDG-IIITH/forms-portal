<script lang="ts">
import { DatePicker } from "$lib/components/ui/date-picker/index.js";

export let question: {
  id: string;
  title: string;
  required: boolean;
  placeholder?: string;
};
export let value: string = "";
export let placeholder: string = "Select a date";
</script>

<DatePicker
  bind:value
  placeholder={question.placeholder || placeholder}
  class="w-full max-w-sm"
/>
