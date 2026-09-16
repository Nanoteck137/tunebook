<script lang="ts">
	import { Users } from "@lucide/svelte";

	let { data } = $props();

	const monthNames = [
		"January", "February", "March", "April", "May", "June",
		"July", "August", "September", "October", "November", "December",
	];
</script>

<div class="flex flex-col gap-10">
	<div
		class="flex flex-col gap-2 rounded-lg border bg-linear-to-b from-[oklch(0.93_0.045_75)] to-background p-6 sm:p-8 dark:from-[oklch(0.24_0.03_80)] dark:to-background"
	>
		<p class="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
			<a href="/users/{data.userData.id}/review/{data.year}" class="hover:underline">
				{data.year} in Review
			</a>
			&middot;
			<a
				href="/users/{data.userData.id}/review/{data.year}/months/{data.month}"
				class="hover:underline"
			>
				{monthNames[data.month - 1]}
			</a>
		</p>

		<div class="flex items-center gap-3">
			<h1 class="text-4xl font-bold md:text-5xl">Top Artists</h1>
			<span class="text-2xl font-bold text-muted-foreground">{data.year}</span>
		</div>

		<p class="text-sm text-muted-foreground">
			{data.artists.length.toLocaleString()}
			{data.artists.length === 1 ? "artist" : "artists"} ranked by plays
		</p>
	</div>

	<section>
		<div class="flex flex-col gap-2">
			{#if data.artists.length > 0}
				{#each data.artists as item (item.artist.id)}
					<a
						href="/artists/{item.artist.id}"
						class="flex items-center gap-3 rounded-lg border bg-card p-2 transition-colors hover:bg-accent"
					>
						<span
							class="w-7 shrink-0 text-center text-sm font-bold text-muted-foreground"
						>
							{item.rank}
						</span>
						<img
							src={item.artist.coverArt.small}
							alt=""
							class="h-10 w-10 shrink-0 rounded object-cover"
						/>
						<span class="min-w-0 flex-1">
							<span class="block truncate text-sm font-medium">
								{item.artist.name}
							</span>
						</span>
						<span class="shrink-0 text-xs text-muted-foreground">
							{item.playCount.toLocaleString()} plays
						</span>
					</a>
				{/each}
			{:else}
				<div
					class="flex flex-col items-center gap-2 rounded-lg border bg-card p-8 text-center"
				>
					<Users class="size-6 text-muted-foreground" />
					<p class="text-sm text-muted-foreground">
						No artists to show this month.
					</p>
				</div>
			{/if}
		</div>
	</section>
</div>