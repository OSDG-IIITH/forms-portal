<script lang="ts">
	import { Card, CardContent } from '$lib/components/ui/card';
	import { cn } from '$lib/utils';
	import { IconFileText, IconEyeOff, IconEdit } from '@tabler/icons-svelte';
	import { formatRelativeTime } from '$lib/utils/date';
	import { goto } from '$app/navigation';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import { writable } from 'svelte/store';

	interface Props {
		form: {
			id: number | string;
			title: string;
			slug?: string;
			owner?: string;
			createdAt: string;
			modified?: string;
			status?: string;
			responses?: number;
			anonymous?: boolean;
			editable_responses?: boolean;
		};
		variant?: 'created' | 'filled';
		viewMode?: string;
		userHandle?: string;
		class?: string;
	}

	let { form, variant = 'created', viewMode = 'grid', userHandle = 'nouser', class: className = '' }: Props = $props();

	const displayTime = form.modified ? formatRelativeTime(form.modified) : formatRelativeTime(form.createdAt);
	const slug = form.slug || 'nouser';
	function handleClick() {
		goto(`/${userHandle}/${slug}/edit`);
	}

	let openTooltipStore = writable<'anonymous' | 'editable' | null>(null);
</script>

{#if viewMode === 'list'}
	<!-- svelte-ignore a11y_click_events_have_key_events -->
	<!-- svelte-ignore a11y_no_static_element_interactions -->
	<div 
		class={cn("flex items-center gap-3 p-3 hover:bg-accent/50 rounded-md transition-colors text-sm cursor-pointer", className)}
		onclick={handleClick}
		role="button"
		tabindex={0}
		onkeydown={(e) => e.key === 'Enter' && handleClick()}
	>
		<div class="flex-shrink-0 w-8 h-8 rounded bg-accent flex items-center justify-center">
			<IconFileText class="w-4 h-4 text-muted-foreground" />
		</div>
		<div class="flex-1 min-w-0">
			<span class="font-medium truncate block">{form.title}</span>
		</div>
		<div class="flex items-center gap-4 ml-4">
			{#if form.anonymous}
				<Tooltip.Provider delayDuration={80}>
					<Tooltip.Root open={$openTooltipStore === 'anonymous'} onOpenChange={v => openTooltipStore.set(v ? 'anonymous' : null)}>
						<Tooltip.Trigger class="group">
							<IconEyeOff class="w-4 h-4 text-muted-foreground/80 transition-transform group-hover:scale-110 group-hover:text-primary" />
						</Tooltip.Trigger>
						<Tooltip.Content>
							<p>Anonymous responses enabled</p>
						</Tooltip.Content>
					</Tooltip.Root>
				</Tooltip.Provider>
			{/if}
			{#if form.editable_responses}
				<Tooltip.Provider delayDuration={80}>
					<Tooltip.Root open={$openTooltipStore === 'editable'} onOpenChange={v => openTooltipStore.set(v ? 'editable' : null)}>
						<Tooltip.Trigger class="group">
							<IconEdit class="w-4 h-4 text-muted-foreground/80 transition-transform group-hover:scale-110 group-hover:text-primary" />
						</Tooltip.Trigger>
						<Tooltip.Content>
							<p>Editable responses allowed</p>
						</Tooltip.Content>
					</Tooltip.Root>
				</Tooltip.Provider>
			{/if}
		</div>
		<div class="w-20 text-xs text-muted-foreground text-right tabular-nums">
			<Tooltip.Provider delayDuration={80}>
				<Tooltip.Root>
					<Tooltip.Trigger class="group">
						<span>{displayTime}</span>
					</Tooltip.Trigger>
					<Tooltip.Content>
						<p>{new Date(form.modified || form.createdAt).toLocaleString(undefined, { dateStyle: 'medium', timeStyle: 'short' })}</p>
					</Tooltip.Content>
				</Tooltip.Root>
			</Tooltip.Provider>
		</div>
	</div>
{:else}
	<!-- svelte-ignore a11y_click_events_have_key_events -->
	<!-- svelte-ignore a11y_no_static_element_interactions -->
	<Card 
		class={cn("p-0 bg-muted/10 hover:bg-muted/30 rounded-md shadow-xs border transition-all duration-300 ease-[cubic-bezier(.4,0,.2,1)] h-40 cursor-pointer will-change-transform will-change-shadow hover:shadow-lg hover:-translate-y-0 hover:scale-[1.012]", className)}
		onclick={handleClick}
		role="button"
		tabindex={0}
		onkeydown={(e) => e.key === 'Enter' && handleClick()}
	>
		<CardContent class="p-0 h-full">
			<div class="h-full flex flex-col">
				<div class="h-24 w-full bg-muted rounded-t-md flex items-center justify-center">
					<IconFileText class="w-8 h-8 text-muted-foreground" />
				</div>
				<div class="flex-1 px-4 py-4 flex flex-col justify-center">
					<h3 class="font-medium text-sm truncate mb-0.5">{form.title}</h3>
					<span class="text-xs text-muted-foreground">
						{#if variant === 'created'}
							{displayTime}
						{:else}
							Submitted {displayTime}
						{/if}
					</span>
				</div>
			</div>
		</CardContent>
	</Card>
{/if}
