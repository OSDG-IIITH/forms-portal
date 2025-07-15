<script lang="ts">
	import { Card, CardContent } from '$lib/components/ui/card';
	import { cn } from '$lib/utils';
	import { IconFileText, IconEyeOff, IconEdit } from '@tabler/icons-svelte';
	import CardFooter from '../ui/card/card-footer.svelte';
	import { formatRelativeTime } from '$lib/utils/date';
	import { goto } from '$app/navigation';

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
		<div class="w-32 text-xs text-muted-foreground text-right tabular-nums">
			{displayTime}
		</div>
		<div class="flex items-center gap-2 ml-4">
			{#if form.anonymous}
				<IconEyeOff class="w-4 h-4 text-muted-foreground" />
			{/if}
			{#if form.editable_responses}
				<IconEdit class="w-4 h-4 text-muted-foreground" />
			{/if}
		</div>
	</div>
{:else}
	<!-- svelte-ignore a11y_click_events_have_key_events -->
	<!-- svelte-ignore a11y_no_static_element_interactions -->
	<Card 
		class={cn("p-0 rounded-md shadow-xs border hover:shadow-sm h-40 cursor-pointer", className)}
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
