<script lang="ts">
	import { CalendarRange } from "@lucide/svelte";

	let { data } = $props();

	// TODO(patrik): Infinite scroll

	function formatDecade(decade: number): string {
		return `${decade}s`;
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

		<div class="flex items-center gap-3">
			<h1 class="text-4xl font-bold md:text-5xl">Top Decades</h1>
			<span class="text-2xl font-bold text-muted-foreground">All time</span>
		</div>

		<p class="text-sm text-muted-foreground">
			{data.page.totalItems.toLocaleString()}
			{data.decades.length === 1 ? "decade" : "decades"} ranked by plays
		</p>
	</div>

	<section>
		<div class="flex flex-col gap-2">
			{#if data.decades.length > 0}
				{#each data.decades as item (item.decade)}
					<div class="flex items-center gap-3 rounded-lg border bg-card p-2">
						<span
							class="w-7 shrink-0 text-center text-sm font-bold text-muted-foreground"
						>
							{item.rank}
						</span>
						<span class="min-w-0 flex-1">
							<span class="block truncate text-sm font-medium">
								{formatDecade(item.decade)}
							</span>
						</span>
						<span class="shrink-0 text-xs text-muted-foreground">
							{item.playCount.toLocaleString()} plays
						</span>
					</div>
				{/each}
			{:else}
				<div
					class="flex flex-col items-center gap-2 rounded-lg border bg-card p-8 text-center"
				>
					<CalendarRange class="size-6 text-muted-foreground" />
					<p class="text-sm text-muted-foreground">No decades to show.</p>
				</div>
			{/if}
		</div>
	</section>
</div>
