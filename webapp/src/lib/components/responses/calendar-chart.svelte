<script lang="ts">
	import CalendarIcon from "@lucide/svelte/icons/calendar";
	import * as Chart from "$lib/components/ui/chart/index.js";
	import * as Card from "$lib/components/ui/card/index.js";
	import * as Popover from "$lib/components/ui/popover/index.js";
	import { Button } from "$lib/components/ui/button/index.js";
	import { BarChart, Highlight, type ChartContextValue } from "layerchart";
	import RangeCalendar from "$lib/components/ui/range-calendar/range-calendar.svelte";
	import type { DateRange } from "bits-ui";
	import {
		CalendarDate,
		getLocalTimeZone,
		today,
	} from "@internationalized/date";
	import { scaleBand } from "d3-scale";
	import { cubicInOut } from "svelte/easing";
	import type { Response } from "./columns.js";

	let { responses = [] }: { responses: Response[] } = $props();

	const chartData = $derived.by(() => {
		const dailyResponses = new Map<string, number>();
		for (const response of responses) {
			if (response.submitted) {
				const date = new Date(response.submitted)
					.toISOString()
					.split("T")[0];
				dailyResponses.set(date, (dailyResponses.get(date) || 0) + 1);
			}
		}
		return Array.from(dailyResponses.entries())
			.map(([date, count]) => ({
				date: new Date(date),
				responses: count,
			}))
			.sort((a, b) => a.date.getTime() - b.date.getTime());
	});

	const { minDate, maxDate } = $derived.by(() => {
		if (chartData.length === 0) {
			const now = today(getLocalTimeZone());
			return {
				minDate: now.subtract({ months: 1 }),
				maxDate: now,
			};
		}
		const min = chartData[0].date;
		const max = chartData[chartData.length - 1].date;
		return {
			minDate: new CalendarDate(
				min.getFullYear(),
				min.getMonth() + 1,
				min.getDate(),
			),
			maxDate: new CalendarDate(
				max.getFullYear(),
				max.getMonth() + 1,
				max.getDate(),
			),
		};
	});

	let value = $state<DateRange | undefined>();

	$effect(() => {
		value = { start: minDate, end: maxDate };
	});

	const filteredData = $derived.by(() => {
		if (!value?.start || !value?.end) return chartData;
		return chartData.filter(({ date }) => {
			const dateObj = new Date(date);
			if (!value) return true;
			const startDate = value.start!.toDate(getLocalTimeZone());
			startDate.setHours(0, 0, 0, 0);
			const endDate = value.end!.toDate(getLocalTimeZone());
			endDate.setHours(23, 59, 59, 999);
			return dateObj >= startDate && dateObj <= endDate;
		});
	});

	const yDomain = $derived.by(() => {
		const data = filteredData;
		if (!data.length) return [0, 1];
		const max = Math.max(...data.map((d) => d.responses));
		return [0, Math.ceil(max * 1.01)];
	});

	const total = $derived(
		filteredData.reduce((acc, curr) => acc + curr.responses, 0),
	);

	const chartConfig = {
		responses: {
			label: "Responses",
			color: "var(--color-primary)",
		},
	} satisfies Chart.ChartConfig;

	let context = $state<ChartContextValue>();
</script>

<Card.Root class="@container/card w-full h-full flex flex-col gap-0">
	<Card.Header class="@md/card:grid flex flex-col border-b">
		<Card.Title>Responses</Card.Title>
		<Card.Description
			>Showing total responses for the selected period.</Card.Description
		>
		<Card.Action class="@md/card:mt-0 mt-2">
			<Popover.Root>
				<Popover.Trigger>
					{#snippet child({ props })}
						<Button
							{...props}
							variant="outline"
							class="w-full justify-start text-left font-normal"
						>
							<CalendarIcon class="mr-2 h-4 w-4" />
							{value?.start && value?.end
								? `${value.start.toDate(getLocalTimeZone()).toLocaleDateString()} - ${value.end.toDate(getLocalTimeZone()).toLocaleDateString()}`
								: "Select Date Range"}
						</Button>
					{/snippet}
				</Popover.Trigger>
				<Popover.Content class="w-auto overflow-hidden p-0" align="end">
					<RangeCalendar
						class="w-full"
						fixedWeeks
						bind:value
						minValue={minDate}
						maxValue={maxDate}
					/>
				</Popover.Content>
			</Popover.Root>
		</Card.Action>
	</Card.Header>
<Card.Content class="p-0 m-0 flex-1 flex flex-col">
   <Chart.Container
	   config={chartConfig}
	   class="mb-2 aspect-auto flex-1 w-full"
   >
			<BarChart
				bind:context
				data={filteredData}
				xScale={scaleBand().padding(0.25)}
				x="date"
				axis="x"
				y="responses"
				
				{yDomain}
				props={{
					bars: {
						stroke: "none",
						rounded: "all",
						radius: 4,
						initialY: context?.height,
						initialHeight: 0,
						motion: {
							x: {
								type: "tween",
								duration: 500,
								easing: cubicInOut,
							},
							width: {
								type: "tween",
								duration: 500,
								easing: cubicInOut,
							},
							height: {
								type: "tween",
								duration: 500,
								easing: cubicInOut,
							},
							y: {
								type: "tween",
								duration: 500,
								easing: cubicInOut,
							},
						},
					},
					xAxis: {
						format: (d) =>
							d.toLocaleDateString("en-US", { day: "numeric" }),
					},
				}}
			>
				{#snippet belowMarks()}
					<Highlight area={{ class: "fill-muted" }} />
				{/snippet}
				{#snippet tooltip()}
					<Chart.Tooltip
						class="w-[150px]"
						nameKey="responses"
						labelFormatter={(d) =>
							d.toLocaleDateString("en-US", {
								month: "short",
								day: "numeric",
								year: "numeric",
							})}
					/>
				{/snippet}
			</BarChart>
		</Chart.Container>
	</Card.Content>
	<Card.Footer class="border-t">
		<div class="text-sm">
			You had
			<span class="font-semibold">{total.toLocaleString()}</span>
			responses in the selected period.
		</div>
	</Card.Footer>
</Card.Root>
