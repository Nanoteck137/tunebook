<script lang="ts">
	import ReviewStats from "./ReviewStats.svelte";
	import SectionHeader from "$lib/components/SectionHeader.svelte";
	import { BarChart3, ChevronRight } from "@lucide/svelte";

	let { data } = $props();

	let currentYear = $derived(new Date().getFullYear());

	let maxTrackCount = $derived(
		Math.max(...data.reviews.map((r) => r.trackCount), 0),
	);
</script>

<div class="flex flex-col gap-6">
	<SectionHeader count={data.reviews.length}>
		<BarChart3 />
		Year in Review
	</SectionHeader>

	<a
		href="/users/{data.userData.id}/total"
		class="group relative flex flex-col gap-4 overflow-hidden rounded-lg border bg-card p-4 transition-colors hover:bg-accent hover:text-accent-foreground"
	>
		<div
			class="pointer-events-none absolute -top-16 -right-16 h-40 w-40 rounded-full bg-linear-to-tr from-logo-1/10 via-logo-2/10 to-logo-3/10 blur-xl"
		></div>

		<div class="relative flex items-center gap-2">
			<span class="text-2xl font-bold">All time</span>
			<ChevronRight
				size={18}
				class="ml-auto text-muted-foreground opacity-0 transition-opacity group-hover:opacity-100"
			/>
		</div>

		{#if data.totalReview}
			<div class="relative">
				<ReviewStats
					trackCount={data.totalReview.trackCount}
					uniqueTracks={data.totalReview.uniqueTracks}
					listeningTime={data.totalReview.listeningTime}
					avgCompletion={data.totalReview.avgCompletion}
					skipCount={data.totalReview.skipCount}
					favoritePlays={data.totalReview.favoritePlays}
				/>
			</div>
		{:else}
			<p class="relative text-sm text-muted-foreground">
				Every play, across every year.
			</p>
		{/if}
	</a>

	{#if data.reviews.length === 0}
		<div
			class="flex flex-col items-center gap-2 rounded-lg border py-16 text-center"
		>
			<BarChart3 size={32} class="text-muted-foreground/40" />
			<p class="text-sm font-medium">No yearly reviews yet</p>
			<p class="max-w-sm text-sm text-muted-foreground">
				A year in review is generated once there is a full year of listening
				data. Your all-time summary is available above.
			</p>
		</div>
	{:else}
		<div class="grid gap-3 sm:grid-cols-2">
			{#each data.reviews as review (review.year)}
				<a
					href="/users/{data.userData.id}/review/{review.year}"
					class="group flex flex-col gap-3 rounded-lg border bg-card p-4 transition-colors hover:bg-accent hover:text-accent-foreground"
				>
					<div class="flex items-center gap-2">
						<span class="text-2xl font-bold">{review.year}</span>
						{#if review.year === currentYear}
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

					<ReviewStats
						trackCount={review.trackCount}
						uniqueTracks={review.uniqueTracks}
						listeningTime={review.listeningTime}
						avgCompletion={review.avgCompletion}
						skipCount={review.skipCount}
						favoritePlays={review.favoritePlays}
					/>

					<div
						class="mt-auto h-1.5 w-full overflow-hidden rounded-full bg-muted"
					>
						<div
							class="h-full rounded-full bg-linear-to-r from-logo-1 to-logo-3"
							style="width: {maxTrackCount > 0
								? (review.trackCount / maxTrackCount) * 100
								: 0}%"
						></div>
					</div>
				</a>
			{/each}
		</div>
	{/if}
</div>
