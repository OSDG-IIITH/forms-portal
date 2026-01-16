<script lang="ts">
  import Textarea from "../../ui/textarea/textarea.svelte";
  export let question: {
    id: string;
    title: string;
    required: boolean;
    placeholder?: string;
    validations?: {
      "max-chars"?: number;
      "min-chars"?: number;
      regex?: string;
    };
  };
  export let value: string = "";
  export let disabled: boolean = false;

  import { createResponseSchema } from "$lib/utils/validation";

  let error: string = "";

  $: {
    const schema = createResponseSchema(question as any);
    const result = schema.safeParse(value);

    if (!result.success) {
      error = result.error.issues[0].message;
    } else {
      error = "";
    }
  }
</script>

<Textarea
  class="block w-full rounded-md border border-input bg-background px-3 py-2 text-base focus:outline-none disabled:opacity-50"
  id={question.id}
  name={question.id}
  placeholder={question.placeholder || "Your answer"}
  bind:value
  {disabled}
  aria-required={question.required}
  maxlength={question.validations?.["max-chars"]}
  minlength={question.validations?.["min-chars"]}
  rows={4}
/>
{#if value && error}
  <p class="text-sm text-destructive mt-1">{error}</p>
{/if}
