<script lang="ts">
	import { goto } from "$app/navigation";
	import {
		Activity,
		BarChart3,
		CalendarDays,
		DiscAlbum,
		Heart,
		ListPlus,
		Music,
		Play,
		Repeat,
		Shuffle,
		SkipForward,
		Sprout,
		Tags,
		Users,
	} from "@lucide/svelte";
	import AlbumTile from "$lib/components/tiles/AlbumTile.svelte";
	import ArtistTile from "$lib/components/tiles/ArtistTile.svelte";
	import SectionHeader from "$lib/components/SectionHeader.svelte";
	import TopTrackItem from "$lib/components/track-list/TopTrackItem.svelte";
	import { DropdownMenu } from "$lib/components/ui";
	import { getFavorites } from "$lib/favorites.svelte";
	import { toast } from "svelte-sonner";

	let { data } = $props();
	const favoritesManager = getFavorites();

	let currentYear = $derived(new Date().getFullYear());
	let review = $derived(data.review);
	let topTracks = $derived(data.topTracks);

	let topArtist = $derived(data.topArtists[0] ?? null);
	let topAlbum = $derived(data.topAlbums[0] ?? null);
	let topTrack = $derived(topTracks[0] ?? null);

	let monthlyHours = $derived(data.months.map((m) => m.playTime / 3600));
	let monthLabels = $derived([
		"Jan", "Feb", "Mar", "Apr", "May", "Jun",
		"Jul", "Aug", "Sep", "Oct", "Nov", "Dec",
	]);
	let maxMonthlyHours = $derived(Math.max(...monthlyHours, 1 / 60));

	let barHeights = $derived(
		monthlyHours.map((h) =>
			h > 0 ? Math.max((h / maxMonthlyHours) * 100, 4) : 0,
		),
	);

	let hoveredMonth = $state(-1);

	let expandedArtist = $state<string | null>(null);
	let expandedAlbum = $state<string | null>(null);

	// function artistTracksFor(id: string) {
	// 	return review.artistTracks.find((a) => a.artist.id === id)?.tracks ?? [];
	// }
	//
	// function albumTracksFor(id: string) {
	// 	return review.albumTracks.find((a) => a.album.id === id)?.tracks ?? [];
	// }

	function toggleArtist(id: string) {
		expandedArtist = expandedArtist === id ? null : id;
	}
	
	function toggleAlbum(id: string) {
		expandedAlbum = expandedAlbum === id ? null : id;
	}

	function formatMonthlyHours(hours: number): string {
		if (hours >= 1) {
			return `${Math.round(hours)}h`;
		}

		return `${Math.round(hours * 60)}m`;
	}

	function formatListeningTime(seconds: number): string {
		const hours = Math.floor(seconds / 3600);
		const minutes = Math.floor((seconds % 3600) / 60);
		return `${hours}h ${minutes}m`;
	}

	function formatPercent(rate: number): string {
		return `${rate.toLocaleString("en-US", { maximumFractionDigits: 1 })}%`;
	}

	let repeatRate = $derived(
		review.uniqueTracks > 0 ? review.trackCount / review.uniqueTracks : 0,
	);
	let skipRate = $derived(
		review.trackCount > 0 ? (review.skipCount / review.trackCount) * 100 : 0,
	);
	let favoritesOverlap = $derived(
		review.trackCount > 0
			? Math.round((review.favoritePlays / review.trackCount) * 100)
			: 0,
	);

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
		<p
			class="text-xs font-semibold tracking-wider text-muted-foreground uppercase"
		>
			Year in Review
		</p>

		<div class="flex items-center gap-3">
			<h1 class="text-4xl font-bold md:text-5xl">{data.year}</h1>
			{#if data.year === currentYear}
				<span
					class="rounded-full bg-primary/10 px-2.5 py-0.5 text-xs font-medium text-primary ring-1 ring-primary/25"
				>
					In progress
				</span>
			{/if}
		</div>

		<!-- {#if review.review.trackCount > 0} -->
		<!-- 	<div class="flex flex-col gap-0.5"> -->
		<!-- 		<p class="text-sm text-muted-foreground"> -->
		<!-- 			{review.review.trackCount.toLocaleString()} plays &middot; -->
		<!-- 			{formatListeningTime(review.review.listeningTime)} -->
		<!-- 		</p> -->
		<!-- 	</div> -->
		<!-- {:else} -->
		<!-- 	<p class="text-sm text-muted-foreground">No listening data this year.</p> -->
		<!-- {/if} -->
	</div>

	<section>
		<SectionHeader>
			<BarChart3 />
			Summary
		</SectionHeader>

		<div class="grid grid-cols-2 gap-4 lg:grid-cols-4">
			<div class="flex flex-col gap-1.5 rounded-lg border bg-card p-4">
				<span class="text-xs text-muted-foreground">Top Artist</span>
				{#if topArtist}
					<a href="/artists/{topArtist.id}" class="flex items-center gap-2">
						<img
							src={topArtist.coverArt.small}
							alt=""
							class="h-10 w-10 rounded object-cover"
						/>
						<span class="truncate text-sm font-medium">
							{topArtist.name}
						</span>
					</a>
				{:else}
					<span class="truncate text-sm font-medium text-muted-foreground">
						—
					</span>
				{/if}
			</div>

			<div class="flex flex-col gap-1.5 rounded-lg border bg-card p-4">
				<span class="text-xs text-muted-foreground">Top Album</span>
				{#if topAlbum}
					<a href="/albums/{topAlbum.id}" class="flex items-center gap-2">
						<img
							src={topAlbum.coverArt.small}
							alt=""
							class="h-10 w-10 rounded object-cover"
						/>
						<span class="truncate text-sm font-medium">
							{topAlbum.name}
						</span>
					</a>
				{:else}
					<span class="truncate text-sm font-medium text-muted-foreground">
						—
					</span>
				{/if}
			</div>

			<div class="flex flex-col gap-1.5 rounded-lg border bg-card p-4">
				<span class="text-xs text-muted-foreground">Top Song</span>
				{#if topTrack}
					<span class="flex items-center gap-2">
						<img
							src={topTrack.coverArt.small}
							alt=""
							class="h-10 w-10 rounded object-cover"
						/>
						<span class="truncate text-sm font-medium">
							{topTrack.name}
						</span>
					</span>
				{:else}
					<span class="truncate text-sm font-medium text-muted-foreground">
						—
					</span>
				{/if}
			</div>
		</div>
	</section>

	<section>
		<SectionHeader>
			<Activity />
			At a Glance
		</SectionHeader>

		<div class="grid grid-cols-2 gap-4 md:grid-cols-4">
			<div class="flex flex-col gap-1.5 rounded-lg border bg-card p-4">
				<span class="text-xs text-muted-foreground">Avg Completion</span>
				<span class="text-2xl font-bold">
					{formatPercent(review.avgCompletion)}
				</span>
			</div>

			<div class="flex flex-col gap-1.5 rounded-lg border bg-card p-4">
				<span class="text-xs text-muted-foreground">Skip Rate</span>
				<span class="text-2xl font-bold">{formatPercent(skipRate)}</span>
			</div>

			<div class="flex flex-col gap-1.5 rounded-lg border bg-card p-4">
				<span class="text-xs text-muted-foreground">Unique Tracks</span>
				<span class="text-2xl font-bold">
					{review.uniqueTracks.toLocaleString()}
				</span>
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

	<section>
		<SectionHeader
			count={data.topTrackCount}
			viewAllHref="/users/{data.userData.id}/review/{data.year}/top-tracks"
		>
			<Music />
			Top Tracks
		</SectionHeader>

		<div class="flex flex-col">
			{#each topTracks as item (item.id)}
				<TopTrackItem rank={item.rank} track={item} playCount={item.playCount}>
					{#snippet menuItems()}
						<DropdownMenu.Group>
							<DropdownMenu.Item>
								<Play />
								Play
							</DropdownMenu.Item>
							<DropdownMenu.Item>
								<Shuffle />
								Shuffle play
							</DropdownMenu.Item>
						</DropdownMenu.Group>

						<DropdownMenu.Separator />

						<DropdownMenu.Item
							onSelect={() => {
								goto(`/albums/${item.track.albumId}`);
							}}
						>
							<DiscAlbum />
							Go to Album
						</DropdownMenu.Item>
						<DropdownMenu.Sub>
							<DropdownMenu.SubTrigger>
								<Users />
								Go to artist
							</DropdownMenu.SubTrigger>
							<DropdownMenu.SubContent>
								{#each item.artists as artist (artist.id)}
									<a
										href="/artists/{artist.id}"
										class="flex items-center gap-2 rounded-sm px-3 py-1.5 text-sm text-popover-foreground hover:bg-accent hover:text-accent-foreground"
									>
										{artist.name}
									</a>
								{/each}
							</DropdownMenu.SubContent>
						</DropdownMenu.Sub>
						<DropdownMenu.Separator />
						<DropdownMenu.Item>
							<ListPlus />
							Add to Playlist
						</DropdownMenu.Item>
						<DropdownMenu.Item
							onSelect={async () => {
								const wasFav = favoritesManager.hasTrack(item.id);
								await favoritesManager.toggleTrack(item.id);
								toast.success(
									wasFav
										? "Removed from favorites"
										: "Added to favorites",
								);
							}}
						>
							{#if favoritesManager.hasTrack(item.id)}
								<Heart class="fill-primary stroke-primary" />
								Unfavorite
							{:else}
								<Heart />
								Favorite
							{/if}
						</DropdownMenu.Item>
					{/snippet}
				</TopTrackItem>
			{/each}
		</div>
	</section>

	<section>
		<SectionHeader
			count={data.topAlbumCount}
			viewAllHref="/users/{data.userData.id}/review/{data.year}/top-albums"
		>
			<DiscAlbum />
			Top Albums
		</SectionHeader>

		<div class="flex items-start gap-4 overflow-x-auto pb-2">
			{#each data.topAlbums as item (item.id)}
				{@const expanded = expandedAlbum === item.id}
				<div
					class="flex w-40 shrink-0 flex-col transition-all"
					class:w-80={expanded}
				>
					<AlbumTile
						id={item.id}
						cover={item.coverArt.small}
						name={item.name}
						artists={item.artists}
						{expanded}
						onToggle={() => toggleAlbum(item.id)}
					/>

					<!-- {#if expanded} -->
					<!-- 	<div class="mt-1 flex w-full flex-col gap-1 rounded-lg border bg-card p-2"> -->
					<!-- 		{#each albumTracksFor(item.album.id) as entry (entry.track.id)} -->
					<!-- 			<a -->
					<!-- 				href="/tracks/{entry.track.id}" -->
					<!-- 				class="flex items-center gap-2 rounded-md p-1.5 hover:bg-accent" -->
					<!-- 			> -->
					<!-- 				<img -->
					<!-- 					src={entry.track.coverArt.small} -->
					<!-- 					alt="" -->
					<!-- 					class="h-8 w-8 shrink-0 rounded object-cover" -->
					<!-- 				/> -->
					<!-- 				<span class="min-w-0 flex-1 truncate text-sm"> -->
					<!-- 					{entry.track.name} -->
					<!-- 				</span> -->
					<!-- 				<span class="shrink-0 text-xs text-muted-foreground"> -->
					<!-- 					{entry.playCount.toLocaleString()} plays -->
					<!-- 				</span> -->
					<!-- 			</a> -->
					<!-- 		{/each} -->
					<!-- 	</div> -->
					<!-- {/if} -->
				</div>
			{/each}
		</div>
	</section>

	<section>
		<SectionHeader
			count={data.topArtistCount}
			viewAllHref="/users/{data.userData.id}/review/{data.year}/top-artists"
		>
			<Users />
			Top Artists
		</SectionHeader>

		<div class="flex items-start gap-4 overflow-x-auto pb-2">
			{#each data.topArtists as item (item.id)}
				{@const expanded = expandedArtist === item.id}
				<div
					class="flex w-40 shrink-0 flex-col transition-all"
					class:w-80={expanded}
				>
					<ArtistTile
						id={item.id}
						cover={item.coverArt.small}
						name={item.name}
						{expanded}
						onToggle={() => toggleArtist(item.id)}
					/>

					<!-- {#if expanded} -->
					<!-- 	<div class="mt-1 flex w-full flex-col gap-1 rounded-lg border bg-card p-2"> -->
					<!-- 		{#each artistTracksFor(item.artist.id) as entry (entry.track.id)} -->
					<!-- 			<a -->
					<!-- 				href="/tracks/{entry.track.id}" -->
					<!-- 				class="flex items-center gap-2 rounded-md p-1.5 hover:bg-accent" -->
					<!-- 			> -->
					<!-- 				<img -->
					<!-- 					src={entry.track.coverArt.small} -->
					<!-- 					alt="" -->
					<!-- 					class="h-8 w-8 shrink-0 rounded object-cover" -->
					<!-- 				/> -->
					<!-- 				<span class="min-w-0 flex-1 truncate text-sm"> -->
					<!-- 					{entry.track.name} -->
					<!-- 				</span> -->
					<!-- 				<span class="shrink-0 text-xs text-muted-foreground"> -->
					<!-- 					{entry.playCount.toLocaleString()} plays -->
					<!-- 				</span> -->
					<!-- 			</a> -->
					<!-- 		{/each} -->
					<!-- 	</div> -->
					<!-- {/if} -->
				</div>
			{/each}
		</div>
	</section>

	<section>
		<SectionHeader>
			<BarChart3 />
			Monthly Listening
		</SectionHeader>

		<div class="flex h-56 items-end gap-2 sm:h-48">
			{#each monthLabels as label, i}
				<div class="flex h-full flex-1 flex-col items-center gap-1.5">
					<span class="text-[10px] text-muted-foreground">
						{formatMonthlyHours(monthlyHours[i])}
					</span>
					<div
						class="flex w-full flex-1 items-end"
						role="button"
						tabindex="-1"
						onpointerenter={() => (hoveredMonth = i)}
						onpointerleave={() => (hoveredMonth = -1)}
					>
						<div class="relative w-full" style="height: {barHeights[i]}%">
							{#if hoveredMonth === i && monthlyHours[i] > 0}
								<div
									class="pointer-events-none absolute bottom-full left-1/2 z-10 mb-1 -translate-x-1/2 rounded-md bg-foreground px-2 py-1 text-xs font-medium whitespace-nowrap text-background shadow-md"
								>
									{label}: {formatMonthlyHours(monthlyHours[i])}
									&middot;
									{data.months[i].playCount.toLocaleString()} plays
								</div>
							{/if}
							<div
								class="h-full w-full rounded-t bg-linear-to-t from-logo-3 to-logo-1 transition-all hover:opacity-80"
							></div>
						</div>
					</div>
					<span class="text-[10px] text-muted-foreground">
						{label}
					</span>
				</div>
			{/each}
		</div>
	</section>

	<section>
		<SectionHeader>
			<CalendarDays />
			Monthly Breakdown
		</SectionHeader>

		<div class="grid grid-cols-4 gap-2 sm:flex sm:flex-wrap">
			{#each monthLabels as label, i (label)}
				{@const count = data.months[i]?.playCount ?? 0}
				<a
					href="/users/{data.userData.id}/review/{data.year}/months/{i + 1}"
					class="flex min-w-0 flex-col items-center gap-0.5 rounded-lg border bg-card px-2 py-2 transition-colors hover:bg-accent sm:min-w-[88px] sm:flex-none"
					class:pointer-events-none={count === 0}
					class:opacity-40={count === 0}
					aria-disabled={count === 0}
				>
					<span class="text-sm font-medium">{label}</span>
					<span class="text-xs text-muted-foreground">
						{count.toLocaleString()}
					</span>
				</a>
			{/each}
		</div>
	</section>

	<section>
		<SectionHeader count={data.topTagCount}>
			<Tags />
			Top Tags
		</SectionHeader>

		<div class="flex flex-wrap gap-2">
			{#each data.topTags as tag (tag.tagSlug)}
				<div
					class="flex items-center gap-2 rounded-full border bg-card px-3 py-1.5"
				>
					<span class="text-sm">
						{tag.rank}. {formatTagSlug(tag.tagSlug)}
					</span>
					<span class="text-xs text-muted-foreground">
						{tag.playCount.toLocaleString()} plays
					</span>
				</div>
			{/each}
		</div>
	</section>

	<section>
		<SectionHeader count={data.topDecadeCount}>
			<DiscAlbum />
			Decades
		</SectionHeader>

		<div class="flex flex-col gap-2">
			{#each data.topDecades as decade (decade.decade)}
				{@const count = decade.playCount}
				{@const maxDecadeCount = data.topDecades[0]?.playCount ?? 1}
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
</div>

