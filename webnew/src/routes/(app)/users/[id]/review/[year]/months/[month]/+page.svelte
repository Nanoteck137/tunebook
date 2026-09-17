<script lang="ts">
	import { goto } from "$app/navigation";
	import { ChevronLeft, ChevronRight, DiscAlbum } from "@lucide/svelte";
	import { Heart, Music, Tags, Users } from "@lucide/svelte";
	import type { Track } from "$lib/api/types";
	import AlbumTile from "$lib/components/tiles/AlbumTile.svelte";
	import ArtistTile from "$lib/components/tiles/ArtistTile.svelte";
	import SectionHeader from "$lib/components/SectionHeader.svelte";

	let { data } = $props();

	const monthNames = [
		"January", "February", "March", "April", "May", "June",
		"July", "August", "September", "October", "November", "December",
	];
	const monthShort = [
		"Jan", "Feb", "Mar", "Apr", "May", "Jun",
		"Jul", "Aug", "Sep", "Oct", "Nov", "Dec",
	];

	let currentYear = $derived(new Date().getFullYear());
	let month = $derived(data.month);
	let monthData = $derived(data.monthData);

	let skipRate = $derived(
		monthData.playCount > 0
			? (monthData.skipCount / monthData.playCount) * 100
			: 0,
	);
	let favoritesOverlap = $derived(
		monthData.playCount > 0
			? Math.round((monthData.favoritePlays / monthData.playCount) * 100)
			: 0,
	);
	let repeatRate = $derived(
		monthData.uniqueTracks > 0 ? monthData.playCount / monthData.uniqueTracks : 0,
	);

	function formatListeningTime(seconds: number): string {
		const hours = Math.floor(seconds / 3600);
		const minutes = Math.floor((seconds % 3600) / 60);
		return `${hours}h ${minutes}m`;
	}

	function trackArtistNames(track: Track): string {
		return track.artists.map((a) => a.name).join(", ");
	}

	function goMonth(next: number) {
		const target = Math.min(12, Math.max(1, next));
		goto(`/users/${data.userData.id}/review/${data.year}/months/${target}`);
	}

	function formatPercent(rate: number): string {
		return `${rate.toLocaleString("en-US", { maximumFractionDigits: 1 })}%`;
	}

	function formatTagSlug(slug: string): string {
		const words = slug.split("-");
		const sentence = words.join(" ");
		return sentence.charAt(0).toUpperCase() + sentence.slice(1);
	}

	function formatDecade(decade: number): string {
		return `${decade}s`;
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

		<div class="flex items-center justify-between gap-3">
			<div class="flex items-center gap-3">
				<h1 class="text-4xl font-bold md:text-5xl">
					{monthNames[month - 1]}
				</h1>
				<span class="text-2xl font-bold text-muted-foreground">{data.year}</span>
				{#if data.year === currentYear && month === new Date().getMonth() + 1}
					<span
						class="rounded-full bg-primary/10 px-2.5 py-0.5 text-xs font-medium text-primary ring-1 ring-primary/25"
					>
						In progress
					</span>
				{/if}
			</div>

			<div class="flex items-center gap-1">
				<button
					type="button"
					class="rounded-md border bg-card p-2 transition-colors hover:bg-accent disabled:opacity-40"
					disabled={month <= 1}
					onclick={() => goMonth(month - 1)}
					aria-label="Previous month"
				>
					<ChevronLeft size={18} />
				</button>
				<button
					type="button"
					class="rounded-md border bg-card p-2 transition-colors hover:bg-accent disabled:opacity-40"
					disabled={month >= 12}
					onclick={() => goMonth(month + 1)}
					aria-label="Next month"
				>
					<ChevronRight size={18} />
				</button>
			</div>
		</div>

		<p class="text-sm text-muted-foreground">
			{monthData.playCount.toLocaleString()} plays &middot;
			{formatListeningTime(monthData.playTime)}
		</p>
	</div>

	<section>
		<SectionHeader>
			<Heart />
			At a Glance
		</SectionHeader>

		<div class="grid grid-cols-2 gap-4 md:grid-cols-4">
			<div class="flex flex-col gap-1.5 rounded-lg border bg-card p-4">
				<span class="text-xs text-muted-foreground">Avg Completion</span>
				<span class="text-2xl font-bold">{formatPercent(monthData.avgCompletion)}</span>
			</div>

			<div class="flex flex-col gap-1.5 rounded-lg border bg-card p-4">
				<span class="text-xs text-muted-foreground">Skip Rate</span>
				<span class="text-2xl font-bold">{formatPercent(skipRate)}</span>
			</div>

			<div class="flex flex-col gap-1.5 rounded-lg border bg-card p-4">
				<span class="text-xs text-muted-foreground">Unique Tracks</span>
				<span class="text-2xl font-bold">{monthData.uniqueTracks.toLocaleString()}</span>
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

			<div class="flex flex-col gap-1.5 rounded-lg border bg-card p-4">
				<span class="text-xs text-muted-foreground">Annual Share</span>
				<span class="text-2xl font-bold">
					{monthData.playCount > 0
						? Math.round(
								(monthData.playCount / data.review.review.trackCount) * 100,
							)
						: 0}%
				</span>
			</div>
		</div>
	</section>

	{#if monthData.tracks.length > 0}
		<section>
			<SectionHeader
				count={monthData.trackCount}
				viewAllHref="/users/{data.userData.id}/review/{data.year}/months/{month}/top-tracks"
			>
				<Music />
				Top Tracks
			</SectionHeader>

			<div class="flex flex-col gap-2">
				{#each monthData.tracks as item (item.track.id)}
					<a
						href="/tracks/{item.track.id}"
						class="flex items-center gap-3 rounded-lg border bg-card p-2 transition-colors hover:bg-accent"
					>
						<span class="w-6 shrink-0 text-center text-sm font-bold text-muted-foreground">
							{item.rank}
						</span>
						<img
							src={item.track.coverArt.small}
							alt=""
							class="h-10 w-10 shrink-0 rounded object-cover"
						/>
						<span class="min-w-0 flex-1">
							<span class="block truncate text-sm font-medium">
								{item.track.name}
							</span>
							<span class="block truncate text-xs text-muted-foreground">
								{trackArtistNames(item.track)}
							</span>
						</span>
						<span class="shrink-0 text-xs text-muted-foreground">
							{item.playCount.toLocaleString()} plays
						</span>
					</a>
				{/each}
			</div>
		</section>
	{/if}

	{#if monthData.albums.length > 0}
		<section>
			<SectionHeader
				count={monthData.albumCount}
				viewAllHref="/users/{data.userData.id}/review/{data.year}/months/{month}/top-albums"
			>
				<DiscAlbum />
				Top Albums
			</SectionHeader>

			<div class="flex items-start gap-4 overflow-x-auto pb-2">
				{#each monthData.albums as item (item.album.id)}
					<div class="w-40 shrink-0">
						<AlbumTile
							id={item.album.id}
							cover={item.album.coverArt.small}
							name={item.album.name}
							artists={item.album.artists}
						/>
					</div>
				{/each}
			</div>
		</section>
	{/if}

	{#if monthData.artists.length > 0}
		<section>
			<SectionHeader
				count={monthData.artistCount}
				viewAllHref="/users/{data.userData.id}/review/{data.year}/months/{month}/top-artists"
			>
				<Users />
				Top Artists
			</SectionHeader>

			<div class="flex items-start gap-4 overflow-x-auto pb-2">
				{#each monthData.artists as item (item.artist.id)}
					<div class="w-40 shrink-0">
						<ArtistTile
							id={item.artist.id}
							cover={item.artist.coverArt.small}
							name={item.artist.name}
						/>
					</div>
				{/each}
			</div>
		</section>
	{/if}

	{#if monthData.tags.length > 0}
		<section>
			<SectionHeader count={monthData.tags.length}>
				<Tags />
				Top Tags
			</SectionHeader>

			<div class="flex flex-wrap gap-2">
				{#each monthData.tags as tag, i (tag.tagSlug)}
					<div
						class="flex items-center gap-2 rounded-full border bg-card px-3 py-1.5"
					>
						<span class="text-sm">
							{i + 1}. {formatTagSlug(tag.tagSlug)}
						</span>
						<span class="text-xs text-muted-foreground">
							{tag.playCount.toLocaleString()} plays
						</span>
					</div>
				{/each}
			</div>
		</section>
	{/if}

	{#if monthData.decades.length > 0}
		<section>
			<SectionHeader count={monthData.decades.length}>
				<DiscAlbum />
				Decades
			</SectionHeader>

			<div class="flex flex-col gap-2">
				{#each monthData.decades as decade (decade.decade)}
					{@const count = decade.playCount}
					{@const maxDecadeCount = monthData.decades[0]?.playCount ?? 1}
					<div class="flex items-center gap-3">
						<span class="w-14 shrink-0 text-sm font-medium">
							{formatDecade(decade.decade)}
						</span>
						<div class="h-6 flex-1 overflow-hidden rounded bg-muted">
							<div
								class="h-full bg-linear-to-r from-logo-3 to-logo-1"
								style="width: {Math.max((count / maxDecadeCount) * 100, 4)}%"
							></div>
						</div>
						<span class="w-20 shrink-0 text-right text-xs text-muted-foreground">
							{count.toLocaleString()} plays
						</span>
					</div>
				{/each}
			</div>
		</section>
	{/if}
</div>