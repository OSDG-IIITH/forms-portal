<script lang="ts">
  import { Input } from '$lib/components/ui/input';
  import { Label } from '$lib/components/ui/label';
  import { Checkbox } from '$lib/components/ui/checkbox';

  interface Props {
    question: {
      uid: string;
      text: string;
      required: boolean;
    };
  }

  let { question = $bindable() }: Props = $props();
</script>

<div class="space-y-4">
  <div class="space-y-2">
    <Label for="question-{question.uid}">Question Text</Label>
    <Input
      id="question-{question.uid}"
      placeholder="Enter your question here"
      bind:value={question.text}
    />
  </div>

  <div class="flex items-center space-x-2">
    <Checkbox 
      id="required-{question.uid}" 
      bind:checked={question.required}
    />
    <Label 
      for="required-{question.uid}" 
      class="text-sm font-normal cursor-pointer"
    >
      Required field
    </Label>
  </div>
</div>
