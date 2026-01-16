<script lang="ts">
  import Input from "../../ui/input/input.svelte";
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
    const schema = createResponseSchema(question as any); // Cast because Question type in components might be slightly different or need import, but let's assume structure matches. Better to fix type if needed, but for now 'as any' safe if structure aligns. Actually, let's look at the props. The prop 'question' is defined inline. It matches what createResponseSchema expects (Question).

    // Convert inline prop type to Question type or just pass it - wait, the prop definition on lines 3-13 matches Question interface roughly.
    // Let's pass it. Since strict types might mismatch slightly due to inline def potentially missing new Zod fields if not updated, but we are just using it for validation rules.

    const result = schema.safeParse(value);
    if (!result.success) {
      error = result.error.issues[0].message;
    } else {
      error = "";
    }
  }
</script>

<Input
  class="block w-full rounded-md border border-input bg-background px-3 py-2 text-base focus:outline-none disabled:opacity-50"
  type="text"
  id={question.id}
  name={question.id}
  placeholder={question.placeholder || "Your answer"}
  bind:value
  {disabled}
  aria-required={question.required}
  maxlength={question.validations?.["max-chars"]}
  minlength={question.validations?.["min-chars"]}
/>
{#if value && error}
  <p class="text-sm text-destructive mt-1">{error}</p>
{/if}
