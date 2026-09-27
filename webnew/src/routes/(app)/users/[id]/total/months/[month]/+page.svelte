<script lang="ts">
	import { Activity, CalendarDays, ChevronLeft, ChevronRight, Heart } from "@lucide/svelte";
	import SectionHeader from "$lib/components/SectionHeader.svelte";
	import Spacer from "$lib/components/Spacer.svelte";
	import { Breadcrumb, Button } from "$lib/components/ui";
	import { formatPlayTime } from "$lib/utils";

	let { data } = $props();

	const monthNames = [
		"January",
		"February",
		"March",
		"April",
		"May",
		"June",
		"July",
		"August",
		"September",
		"October",
		"November",
		"December",
	];

	let month = $derived(data.month);
	let monthName = $derived(monthNames[data.monthNum - 1]);

	let skipRate = $derived(
		month && month.playCount > 0
			? (month.skipCount / month.playCount) * 100
			: 0,
	);
	let favoritesOverlap = $derived(
		month && month.playCount > 0
			? (month.favoritePlays / month.playCount) * 100
			: 0,
	);
	let repeatRate = $derived(
		month && month.uniqueTracks > 0 ? month.playCount / month.uniqueTracks : 0,
	);

	function formatPercent(value: number): string {
		return `${value.toLocaleString("en-US", { maximumFractionDigits: 1 })}%`;
	}
</script>

<div class="flex flex-col gap-6">
	<Breadcrumb.Root>
		<Breadcrumb.List>
			<Breadcrumb.Item>
				<Breadcrumb.Link href="/users/{data.userData.id}/total">
					All time
				</Breadcrumb.Link>
			</Breadcrumb.Item>
			<Breadcrumb.Separator />
			<Breadcrumb.Item>
				<Breadcrumb.Page>{monthName}</Breadcrumb.Page>
			</Breadcrumb.Item>
		</Breadcrumb.List>
	</Breadcrumb.Root>

	<SectionHeader count={month?.playCount ?? 0}>
		<CalendarDays />
		{monthName}

		{#snippet actions()}
			<div class="flex items-center gap-1">
				<Button
					variant="ghost"
					size="icon-sm"
					disabled={data.monthNum <= 1}
					href={
						data.monthNum > 1
							? `/users/${data.userData.id}/total/months/${data.monthNum - 1}`
							: undefined
					}
					aria-label="Previous month"
				>
					<ChevronLeft />
				</Button>
				<Button
					variant="ghost"
					size="icon-sm"
					disabled={data.monthNum >= 12}
					href={
						data.monthNum < 12
							? `/users/${data.userData.id}/total/months/${data.monthNum + 1}`
							: undefined
					}
					aria-label="Next month"
				>
					<ChevronRight />
				</Button>
			</div>
		{/snippet}
	</SectionHeader>

	{#if month}
		<p class="text-sm text-muted-foreground">
			{formatPlayTime(month.playTime)} listened across
			{month.uniqueTracks.toLocaleString()} unique tracks
		</p>

		<section>
			<SectionHeader>
				<Activity />
				Listening Stats
			</SectionHeader>
			<Spacer />

			<div class="grid grid-cols-2 gap-3 md:grid-cols-4">
				<div class="flex flex-col gap-1.5 rounded-lg border bg-card p-4">
					<span class="text-xs text-muted-foreground">Plays</span>
					<span class="text-2xl font-bold">
						{month.playCount.toLocaleString()}
					</span>
				</div>
				<div class="flex flex-col gap-1.5 rounded-lg border bg-card p-4">
					<span class="text-xs text-muted-foreground">Listening Time</span>
					<span class="text-2xl font-bold">
						{formatPlayTime(month.playTime)}
					</span>
				</div>
				<div class="flex flex-col gap-1.5 rounded-lg border bg-card p-4">
					<span class="text-xs text-muted-foreground">Unique Tracks</span>
					<span class="text-2xl font-bold">
						{month.uniqueTracks.toLocaleString()}
					</span>
				</div>
				<div class="flex flex-col gap-1.5 rounded-lg border bg-card p-4">
					<span class="text-xs text-muted-foreground">Skipped</span>
					<span class="text-2xl font-bold">
						{month.skipCount.toLocaleString()}
					</span>
				</div>
				<div class="flex flex-col gap-1.5 rounded-lg border bg-card p-4">
					<span class="text-xs text-muted-foreground">Avg Completion</span>
					<span class="text-2xl font-bold">
						{formatPercent(month.avgCompletion)}
					</span>
				</div>
				<div class="flex flex-col gap-1.5 rounded-lg border bg-card p-4">
					<span class="text-xs text-muted-foreground">Skip Rate</span>
					<span class="text-2xl font-bold">{formatPercent(skipRate)}</span>
				</div>
				<div class="flex flex-col gap-1.5 rounded-lg border bg-card p-4">
					<span class="text-xs text-muted-foreground">Repeat Rate</span>
					<span class="text-2xl font-bold">{repeatRate.toFixed(1)}x</span>
				</div>
				<div class="flex flex-col gap-1.5 rounded-lg border bg-card p-4">
					<span class="text-xs text-muted-foreground">Favorites Overlap</span>
					<div class="flex items-baseline gap-1.5">
						<span class="text-2xl font-bold">
							{formatPercent(favoritesOverlap)}
						</span>
						<Heart size={16} class="shrink-0 fill-primary stroke-primary" />
					</div>
				</div>
			</div>
		</section>
	{:else}
		<div
			class="flex flex-col items-center gap-2 rounded-lg border py-16 text-center"
		>
			<CalendarDays size={32} class="text-muted-foreground/40" />
			<p class="text-sm font-medium">No listening data for {monthName}</p>
			<p class="max-w-sm text-sm text-muted-foreground">
				This month has no plays in your all-time review. Use the arrows to browse
				the other months.
			</p>
		</div>
	{/if}
</div>
