<script lang="ts">
	import { Input } from "$lib/components/ui/input";
	import { Label } from "$lib/components/ui/label";
	import { DatePicker } from "$lib/components/ui/date-picker";
	import { Checkbox } from "$lib/components/ui/checkbox";
	import { TimeField } from "bits-ui";
	import { getContext } from "svelte";
	import type { FormStore } from "./form-store.svelte";
	import type { Time } from "@internationalized/date";

	const store: FormStore = getContext("form-store");
	const { formData } = store;
</script>

<div class="grid flex-1 auto-rows-min gap-6 px-8 py-2 overflow-y-auto">
	<div class="grid grid-cols-1 gap-6 pt-2">
		<div class="space-y-2">
			<Label for="form-opens">Opens</Label>
			<div class="flex gap-2">
				<DatePicker
					bind:value={formData.opens}
					placeholder="Select open date"
					class="flex-1"
				/>
				<div class="w-22">
					<TimeField.Root
						bind:value={formData.opensTime}
						hourCycle={24}
					>
						<TimeField.Input
							class="flex h-9 w-full rounded-md border border-input bg-muted/50 px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
						>
							{#snippet children({ segments })}
								{#each segments as { part, value }, i (part + i)}
									<div class="inline-block select-none">
										{#if part === "literal"}
											<TimeField.Segment
												{part}
												class="text-muted-foreground p-1"
											>
												{value}
											</TimeField.Segment>
										{:else}
											<TimeField.Segment
												{part}
												class="rounded-sm px-1 py-1 hover:bg-muted focus:bg-muted focus:text-foreground aria-[valuetext=Empty]:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-border"
											>
												{value}
											</TimeField.Segment>
										{/if}
									</div>
								{/each}
							{/snippet}
						</TimeField.Input>
					</TimeField.Root>
				</div>
			</div>
		</div>
		<div class="space-y-2">
			<Label for="form-closes">Closes</Label>
			<div class="flex gap-2">
				<DatePicker
					bind:value={formData.closes}
					placeholder="Select close date"
					class="flex-1"
				/>
				<div class="w-22">
					<TimeField.Root
						bind:value={formData.closesTime}
						hourCycle={24}
					>
						<TimeField.Input
							class="flex h-9 w-full rounded-md border border-input bg-muted/50 px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
						>
							{#snippet children({ segments })}
								{#each segments as { part, value }, i (part + i)}
									<div class="inline-block select-none">
										{#if part === "literal"}
											<TimeField.Segment
												{part}
												class="text-muted-foreground p-1"
											>
												{value}
											</TimeField.Segment>
										{:else}
											<TimeField.Segment
												{part}
												class="rounded-sm px-1 py-1 hover:bg-muted focus:bg-muted focus:text-foreground aria-[valuetext=Empty]:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-border"
											>
												{value}
											</TimeField.Segment>
										{/if}
									</div>
								{/each}
							{/snippet}
						</TimeField.Input>
					</TimeField.Root>
				</div>
			</div>
		</div>
		<div class="space-y-3">
			<Label for="form-max-responses">Max Responses</Label>
			<Input
				id="form-max-responses"
				type="number"
				min="1"
				placeholder="Unlimited"
				bind:value={formData.max_responses}
			/>
		</div>
		<div class="space-y-3">
			<Label for="form-individual-limit">Individual Limit</Label>
			<Input
				id="form-individual-limit"
				type="number"
				min="1"
				bind:value={formData.individual_limit}
			/>
		</div>
		<div class="flex items-start gap-3 p-3 border rounded-md bg-muted/20">
			<Checkbox
				id="form-anonymous"
				bind:checked={formData.anonymous}
				class="mt-0.5"
				onclick={(e) => e.stopPropagation()}
			/>
			<div>
				<Label for="form-anonymous" class="text-sm font-medium"
					>Anonymous Responses</Label
				>
				<p class="text-xs text-muted-foreground">
					Allow users to submit responses without identification.
				</p>
			</div>
		</div>
		<div class="flex items-start gap-3 p-3 border rounded-md bg-muted/20">
			<Checkbox
				id="form-editable-responses"
				bind:checked={formData.editable_responses}
				class="mt-0.5"
				onclick={(e) => e.stopPropagation()}
			/>
			<div>
				<Label for="form-editable-responses" class="text-sm font-medium"
					>Editable Responses</Label
				>
				<p class="text-xs text-muted-foreground">
					Allow users to edit their submitted responses.
				</p>
			</div>
		</div>
	</div>
</div>
