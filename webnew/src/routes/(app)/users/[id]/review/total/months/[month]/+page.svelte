<script lang="ts">
	import { goto } from "$app/navigation";
	import { ChevronLeft, ChevronRight, Heart } from "@lucide/svelte";
	import SectionHeader from "$lib/components/SectionHeader.svelte";

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

	let skipRate = $derived(
		month && month.playCount > 0
			? (month.skipCount / month.playCount) * 100
			: 0,
	);
	let favoritesOverlap = $derived(
		month && month.playCount > 0
			? Math.round((month.favoritePlays / month.playCount) * 100)
			: 0,
	);
	let repeatRate = $derived(
		month && month.uniqueTracks > 0 ? month.playCount / month.uniqueTracks : 0,
	);

	function formatListeningTime(seconds: number): string {
		const hours = Math.floor(seconds / 3600);
		const minutes = Math.floor((seconds % 3600) / 60);
		return `${hours}h ${minutes}m`;
	}

	function goMonth(next: number) {
		const target = Math.min(12, Math.max(1, next));
		goto(`/users/${data.userData.id}/review/total/months/${target}`);
	}

	function formatPercent(rate: number): string {
		return `${rate.toLocaleString("en-US", { maximumFractionDigits: 1 })}%`;
	}
</script>

<div class="flex flex-col gap-10">
	<div
		class="flex flex-col gap-2 rounded-lg border bg-linear-to-b from-[oklch(0.93_0.045_75)] to-background p-6 sm:p-8 dark:from-[oklch(0.24_0.03_80)] dark:to-background"
	>
		<p
			class="text-xs font-semibold tracking-wider text-muted-foreground uppercase"
		>
			<a href="/users/{data.userData.id}/review/total" class="hover:underline">
				Total Review
			</a>
		</p>

		<div class="flex items-center justify-between gap-3">
			<div class="flex items-center gap-3">
				<h1 class="text-4xl font-bold md:text-5xl">
					{monthNames[data.monthNum - 1]}
				</h1>
				<span class="text-2xl font-bold text-muted-foreground">All time</span>
			</div>

			<div class="flex items-center gap-1">
				<button
					type="button"
					class="rounded-md border bg-card p-2 transition-colors hover:bg-accent disabled:opacity-40"
					disabled={data.monthNum <= 1}
					onclick={() => goMonth(data.monthNum - 1)}
					aria-label="Previous month"
				>
					<ChevronLeft size={18} />
				</button>
				<button
					type="button"
					class="rounded-md border bg-card p-2 transition-colors hover:bg-accent disabled:opacity-40"
					disabled={data.monthNum >= 12}
					onclick={() => goMonth(data.monthNum + 1)}
					aria-label="Next month"
				>
					<ChevronRight size={18} />
				</button>
			</div>
		</div>

		{#if month}
			<p class="text-sm text-muted-foreground">
				{month.playCount.toLocaleString()} plays &middot;
				{formatListeningTime(month.playTime)}
			</p>
		{:else}
			<p class="text-sm text-muted-foreground">
				No listening data this month.
			</p>
		{/if}
	</div>

	{#if month}
		<section>
			<SectionHeader>
				<Heart />
				At a Glance
			</SectionHeader>

			<div class="grid grid-cols-2 gap-4 md:grid-cols-4">
				<div class="flex flex-col gap-1.5 rounded-lg border bg-card p-4">
					<span class="text-xs text-muted-foreground">Avg Completion</span>
					<span class="text-2xl font-bold"
						>{formatPercent(month.avgCompletion)}</span
					>
				</div>

				<div class="flex flex-col gap-1.5 rounded-lg border bg-card p-4">
					<span class="text-xs text-muted-foreground">Skip Rate</span>
					<span class="text-2xl font-bold">{formatPercent(skipRate)}</span>
				</div>

				<div class="flex flex-col gap-1.5 rounded-lg border bg-card p-4">
					<span class="text-xs text-muted-foreground">Unique Tracks</span>
					<span class="text-2xl font-bold"
						>{month.uniqueTracks.toLocaleString()}</span
					>
				</div>

				<div class="flex flex-col gap-1.5 rounded-lg border bg-card p-4">
					<span class="text-xs text-muted-foreground">Repeat Rate</span>
					<span class="text-2xl font-bold">{repeatRate.toFixed(1)}x</span>
				</div>

				<div class="flex flex-col gap-1.5 rounded-lg border bg-card p-4">
					<span class="text-xs text-muted-foreground">Favorites Overlap</span>
					<div class="flex items-baseline gap-1.5">
						<span class="text-2xl font-bold">{favoritesOverlap}%</span>
						<Heart size={16} class="shrink-0 fill-primary stroke-primary" />
					</div>
				</div>
			</div>
		</section>
	{/if}
</div>
