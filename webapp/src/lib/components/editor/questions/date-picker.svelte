<script lang="ts">
  import { Input } from '$lib/components/ui/input';
  import { Label } from '$lib/components/ui/label';
  import { Checkbox } from '$lib/components/ui/checkbox';

  interface Props {
    question: {
      id: string;
      title: string;
      required: boolean;
    };
  }

  let { question = $bindable() }: Props = $props();
</script>

<div class="space-y-4">
  <div class="space-y-2">
    <Label for="question-{question.id}">Question Text</Label>
    <Input
      id="question-{question.id}"
      placeholder="Enter your question here"
      bind:value={question.title}
    />
  </div>

  <div class="flex items-center space-x-2">
    <Checkbox id="required-{question.id}" bind:checked={question.required} />
    <Label for="required-{question.id}" class="text-sm font-normal cursor-pointer">
      Required field
    </Label>
  </div>
</div>