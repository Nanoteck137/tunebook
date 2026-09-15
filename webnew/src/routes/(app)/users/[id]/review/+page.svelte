<script lang="ts">
	import { BarChart3, ChevronRight } from "@lucide/svelte";
	import type { YearStat } from "$lib/api/types";

	// TODO(api): add `generatedAt` (ISO timestamp of when the report was
	// generated) to the GetUserYearStats response.
	type ReviewStat = YearStat & { generatedAt?: string };

	let { data } = $props();

	let yearStats = $derived((data.yearStats as ReviewStat[] | null) ?? []);

	let currentYear = $derived(new Date().getFullYear());

	let maxTrackCount = $derived(
		Math.max(...data.yearStats.map((s) => s.trackCount), 0),
	);

	function formatListeningTime(seconds: number): string {
		const hours = Math.floor(seconds / 3600);
		const minutes = Math.floor((seconds % 3600) / 60);
		return `${hours}h ${minutes}m`;
	}

	function formatReportDate(iso?: string): string {
		const d = iso ? new Date(iso) : new Date();
		return d.toLocaleDateString(undefined, {
			month: "short",
			day: "numeric",
			year: "numeric",
		});
	}
</script>

<div class="flex flex-col gap-6">
	<div class="flex flex-col gap-1">
		<h1 class="text-xl font-bold">Year in Review</h1>
		<p class="text-sm text-muted-foreground">
			A look back at {data.userData.displayName}'s listening habits, year by
			year.
		</p>
	</div>

	{#if data.yearStats.length === 0}
		<div class="flex flex-col items-center gap-2 rounded-lg border py-16">
			<BarChart3 size={32} class="text-muted-foreground/40" />
			<p class="text-sm text-muted-foreground">No listening data yet</p>
		</div>
	{:else}
		<div class="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
			{#each yearStats as stat (stat.year)}
				<a
					href="/users/{data.userData.id}/review/{stat.year}"
					class="group flex flex-col gap-3 rounded-lg border bg-card p-4 transition-colors hover:bg-accent hover:text-accent-foreground"
				>
					<div class="flex items-center gap-2">
						<span class="text-3xl font-bold">{stat.year}</span>
						{#if stat.year === currentYear}
							<span
								class="rounded-full bg-primary/10 px-2.5 py-0.5 text-xs font-medium text-primary ring-1 ring-primary/25"
							>
								In progress
							</span>
						{/if}
						<ChevronRight
							size={18}
							class="ml-auto text-muted-foreground opacity-0 transition-opacity group-hover:opacity-100"
						/>
					</div>

					<div class="flex flex-col gap-0.5">
						<span class="text-sm text-muted-foreground">
							{stat.trackCount.toLocaleString()} tracks &middot;
							{formatListeningTime(stat.listeningTime)}
						</span>
						<span class="text-xs text-muted-foreground/80">
							Report generated {formatReportDate(stat.generatedAt)}
						</span>
					</div>

					<div class="h-1.5 w-full overflow-hidden rounded-full bg-muted">
						<div
							class="h-full rounded-full bg-linear-to-r from-logo-1 to-logo-3"
							style="width: {maxTrackCount > 0
								? (stat.trackCount / maxTrackCount) * 100
								: 0}%"
						></div>
					</div>
				</a>
			{/each}
		</div>
	{/if}
</div>