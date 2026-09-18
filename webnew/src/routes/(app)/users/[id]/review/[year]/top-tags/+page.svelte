<script lang="ts">
	import { Tags } from "@lucide/svelte";

	let { data } = $props();

	// TODO(patrik): Infinite scroll 

	function formatTagSlug(slug: string): string {
		const words = slug.split("-");
		const sentence = words.join(" ");
		return sentence.charAt(0).toUpperCase() + sentence.slice(1);
	}
</script>

<div class="flex flex-col gap-10">
	<div
		class="flex flex-col gap-2 rounded-lg border bg-linear-to-b from-[oklch(0.93_0.045_75)] to-background p-6 sm:p-8 dark:from-[oklch(0.24_0.03_80)] dark:to-background"
	>
		<p class="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
			<a href="/users/{data.userData.id}/review/{data.year}" class="hover:underline">
				{data.year} in Review
			</a>
		</p>

		<div class="flex items-center gap-3">
			<h1 class="text-4xl font-bold md:text-5xl">Top Tags</h1>
			<span class="text-2xl font-bold text-muted-foreground">{data.year}</span>
		</div>

		<p class="text-sm text-muted-foreground">
			{data.page.totalItems.toLocaleString()}
			{data.tags.length === 1 ? "tag" : "tags"} ranked by plays
		</p>
	</div>

	<section>
		<div class="flex flex-col gap-2">
			{#if data.tags.length > 0}
				{#each data.tags as item (item.tagSlug)}
					<div
						class="flex items-center gap-3 rounded-lg border bg-card p-2"
					>
						<span
							class="w-7 shrink-0 text-center text-sm font-bold text-muted-foreground"
						>
							{item.rank}
						</span>
						<span class="min-w-0 flex-1">
							<span class="block truncate text-sm font-medium">
								{formatTagSlug(item.tagSlug)}
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
					<Tags class="size-6 text-muted-foreground" />
					<p class="text-sm text-muted-foreground">
						No tags to show this year.
					</p>
				</div>
			{/if}
		</div>
	</section>
</div>