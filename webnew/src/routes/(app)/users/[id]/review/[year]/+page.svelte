<script lang="ts">
	import { goto } from "$app/navigation";
	import {
		BarChart3,
		CalendarCheck,
		DiscAlbum,
		Heart,
		ListPlus,
		Music,
		Play,
		Shuffle,
		Users,
	} from "@lucide/svelte";
	import type {
		Album,
		Artist,
		ArtistInfo,
		Track,
	} from "$lib/api/types";
	import AlbumTile from "$lib/components/tiles/AlbumTile.svelte";
	import ArtistTile from "$lib/components/tiles/ArtistTile.svelte";
	import SectionHeader from "$lib/components/SectionHeader.svelte";
	import TopTrackItem from "$lib/components/track-list/TopTrackItem.svelte";
	import { DropdownMenu } from "$lib/components/ui";
	import { getFavorites } from "$lib/favorites.svelte";
	import { getQuickPlaylist } from "$lib/quick-playlist.svelte";
	import { toast } from "svelte-sonner";

	let { data } = $props();
	const favoritesManager = getFavorites();
	const quickPlaylistManager = getQuickPlaylist();

	// NOTE(api): the pieces below are placeholder data until the year-review
	// endpoint exists (e.g. GET /api/v1/users/{id}/review/{year}). Replace each
	// block with real data from the API when ready.

	function placeholderCover(label: string, hue: number): string {
		const svg = `<svg xmlns='http://www.w3.org/2000/svg' width='160' height='160' viewBox='0 0 160 160'><defs><linearGradient id='g' x1='0' y1='0' x2='1' y2='1'><stop offset='0%' stop-color='hsl(${hue},70%,55%)'/><stop offset='100%' stop-color='hsl(${hue + 40},70%,35%)'/></linearGradient></defs><rect width='160' height='160' fill='url(#g)'/><text x='50%' y='55%' font-family='sans-serif' font-size='56' fill='rgba(255,255,255,0.9)' text-anchor='middle' dominant-baseline='middle'>${label
		.slice(0, 1)
		.toUpperCase()}</text></svg>`;
		return `data:image/svg+xml,${encodeURIComponent(svg)}`;
	}

	// TODO(api): `generatedAt` comes from the year-review endpoint
	// (e.g. GET /api/v1/users/{id}/review/{year}) — date the report was generated.
	let currentYear = $derived(new Date().getFullYear());
	let reportGeneratedAt = $derived(new Date().toISOString());

	function mockTrack(
		id: string,
		name: string,
		artists: ArtistInfo[],
		albumName: string,
	): Track {
		const cover = placeholderCover(name, (id.length * 47) % 360);
		return {
			id,
			name,
			order: null,
			duration: 240,
			number: null,
			year: data.year,
			coverArt: { original: cover, small: cover, medium: cover, large: cover },
			albumId: `${id}-album`,
			albumName,
			artists,
			tags: [],
			created: new Date().toISOString(),
			updated: new Date().toISOString(),
		};
	}

	function mockAlbum(
		id: string,
		name: string,
		artists: ArtistInfo[],
	): Album {
		const cover = placeholderCover(name, (id.length * 61) % 360);
		return {
			id,
			name,
			year: data.year,
			albumType: "album",
			coverArt: { original: cover, small: cover, medium: cover, large: cover },
			artists,
			tags: [],
			created: new Date().toISOString(),
			updated: new Date().toISOString(),
		};
	}

	function mockArtist(id: string, name: string): Artist {
		const cover = placeholderCover(name, (id.length * 83) % 360);
		return {
			id,
			name,
			coverArt: { original: cover, small: cover, medium: cover, large: cover },
			tags: [],
			created: new Date().toISOString(),
			updated: new Date().toISOString(),
		};
	}

	// TODO(api): top artist/album/track + most played day for the year
	const summary = {
		topArtist: { id: "artist-1", name: "Sample Artist", cover: placeholderCover("S", 210) },
		topAlbum: { id: "album-1", name: "Sample Album", cover: placeholderCover("A", 20), artists: [{ id: "artist-1", name: "Sample Artist" }] as ArtistInfo[] },
		topTrack: { id: "track-1", name: "Sample Song", cover: placeholderCover("T", 130), artists: [{ id: "artist-1", name: "Sample Artist" }] as ArtistInfo[] },
		mostPlayedDay: { label: "Wednesday", plays: 842 },
	};

	// TODO(api): top tracks for the year
	const placeholders: Track[] = [
		mockTrack("track-1", "Sample Song", [{ id: "artist-1", name: "Sample Artist" }], "Sample Album"),
		mockTrack("track-2", "Another Song", [{ id: "artist-2", name: "Another Artist" }], "Another Album"),
		mockTrack("track-3", "Third Song", [{ id: "artist-1", name: "Sample Artist" }], "Third Album"),
		mockTrack("track-4", "Fourth Song", [{ id: "artist-3", name: "Third Artist" }], "Fourth Album"),
		mockTrack("track-5", "Fifth Song", [{ id: "artist-1", name: "Sample Artist" }], "Fifth Album"),
	];

	// TODO(api): top albums for the year
	const topAlbums: Album[] = [
		mockAlbum("album-1", "Sample Album", [{ id: "artist-1", name: "Sample Artist" }]),
		mockAlbum("album-2", "Another Album", [{ id: "artist-2", name: "Another Artist" }]),
		mockAlbum("album-3", "Third Album", [{ id: "artist-1", name: "Sample Artist" }]),
		mockAlbum("album-4", "Fourth Album", [{ id: "artist-3", name: "Third Artist" }]),
		mockAlbum("album-5", "Fifth Album", [{ id: "artist-1", name: "Sample Artist" }]),
	];

	// TODO(api): top artists for the year
	const topArtists: Artist[] = [
		mockArtist("artist-1", "Sample Artist"),
		mockArtist("artist-2", "Another Artist"),
		mockArtist("artist-3", "Third Artist"),
		mockArtist("artist-4", "Fourth Artist"),
		mockArtist("artist-5", "Fifth Artist"),
	];

	// TODO(api): monthly listening time for the year (minutes per month)
	const monthlyMinutes = [184, 96, 212, 148, 260, 174, 128, 231, 95, 158, 203, 246];
	const monthLabels = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];
	const maxMonthly = Math.max(...monthlyMinutes);

	function formatListeningTime(seconds: number): string {
		const hours = Math.floor(seconds / 3600);
		const minutes = Math.floor((seconds % 3600) / 60);
		return `${hours}h ${minutes}m`;
	}

	function formatReportDate(iso: string): string {
		return new Date(iso).toLocaleDateString(undefined, {
			month: "short",
			day: "numeric",
			year: "numeric",
		});
	}
</script>

<div class="flex flex-col gap-10">
	<div
		class="flex flex-col gap-2 rounded-lg border bg-linear-to-b from-[oklch(0.93_0.045_75)] to-background p-6 sm:p-8 dark:from-[oklch(0.24_0.03_80)] dark:to-background"
	>
		<p class="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
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

		{#if data.yearStat}
			<div class="flex flex-col gap-0.5">
				<p class="text-sm text-muted-foreground">
					{data.yearStat.trackCount.toLocaleString()} tracks &middot;
					{formatListeningTime(data.yearStat.listeningTime)}
				</p>
				<p class="text-xs text-muted-foreground/80">
					Report generated {formatReportDate(reportGeneratedAt)}
				</p>
			</div>
		{:else}
			<p class="text-sm text-muted-foreground">No listening data this year.</p>
		{/if}
	</div>

	<!-- NOTE(api): this section needs a year-review endpoint -->
	<section>
		<SectionHeader>
			<BarChart3 />
			Summary
		</SectionHeader>

		<div class="grid grid-cols-2 gap-4 lg:grid-cols-4">
			<div class="flex flex-col gap-1.5 rounded-lg border bg-card p-4">
				<span class="text-xs text-muted-foreground">Top Artist</span>
				<div class="flex items-center gap-2">
					<img src={summary.topArtist.cover} alt="" class="h-10 w-10 rounded object-cover" />
					<span class="truncate text-sm font-medium">{summary.topArtist.name}</span>
				</div>
			</div>

			<div class="flex flex-col gap-1.5 rounded-lg border bg-card p-4">
				<span class="text-xs text-muted-foreground">Top Album</span>
				<div class="flex items-center gap-2">
					<img src={summary.topAlbum.cover} alt="" class="h-10 w-10 rounded object-cover" />
					<span class="truncate text-sm font-medium">{summary.topAlbum.name}</span>
				</div>
			</div>

			<div class="flex flex-col gap-1.5 rounded-lg border bg-card p-4">
				<span class="text-xs text-muted-foreground">Top Song</span>
				<div class="flex items-center gap-2">
					<img src={summary.topTrack.cover} alt="" class="h-10 w-10 rounded object-cover" />
					<span class="truncate text-sm font-medium">{summary.topTrack.name}</span>
				</div>
			</div>

			<div class="flex flex-col gap-1.5 rounded-lg border bg-card p-4">
				<span class="text-xs text-muted-foreground">Most Played Day</span>
				<div class="flex items-center gap-2">
					<CalendarCheck size={20} class="shrink-0 text-muted-foreground" />
					<span class="truncate text-sm font-medium">
						{summary.mostPlayedDay.label}
					</span>
				</div>
				<span class="text-xs text-muted-foreground">
					{summary.mostPlayedDay.plays.toLocaleString()} plays
				</span>
			</div>
		</div>
	</section>

	<!-- TODO(api): top tracks for the year, replace `placeholders` -->
	<section>
		<SectionHeader count={placeholders.length}>
			<Music />
			Top Tracks
		</SectionHeader>

		<div class="flex flex-col">
			{#each placeholders as track, i (track.id)}
				<TopTrackItem rank={i + 1} {track} playCount={i * 97 + 356}>
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
								goto(`/albums/${track.albumId}`);
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
								{#each track.artists as artist (artist.id)}
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
								const wasFav = favoritesManager.hasTrack(track.id);
								await favoritesManager.toggleTrack(track.id);
								toast.success(
									wasFav
										? "Removed from favorites"
										: "Added to favorites",
								);
							}}
						>
							{#if favoritesManager.hasTrack(track.id)}
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

	<!-- TODO(api): top albums for the year, replace `topAlbums` -->
	<section>
		<SectionHeader count={topAlbums.length}>
			<DiscAlbum />
			Top Albums
		</SectionHeader>

		<div class="flex gap-4 overflow-x-auto pb-2">
			{#each topAlbums as album (album.id)}
				<AlbumTile
					id={album.id}
					cover={album.coverArt.small}
					name={album.name}
					artists={album.artists}
				/>
			{/each}
		</div>
	</section>

	<!-- TODO(api): top artists for the year, replace `topArtists` -->
	<section>
		<SectionHeader count={topArtists.length}>
			<Users />
			Top Artists
		</SectionHeader>

		<div class="flex gap-4 overflow-x-auto pb-2">
			{#each topArtists as artist (artist.id)}
				<ArtistTile
					id={artist.id}
					cover={artist.coverArt.small}
					name={artist.name}
				/>
			{/each}
		</div>
	</section>

	<!-- TODO(api): monthly listening time for the year, replace `monthlyMinutes` -->
	<section>
		<SectionHeader>
			<BarChart3 />
			Monthly Listening
		</SectionHeader>

		<div class="flex h-40 items-end gap-2">
			{#each monthlyMinutes as minutes, i}
				<div class="flex flex-1 flex-col items-center gap-1.5">
					<span class="text-[10px] text-muted-foreground">
						{Math.round(minutes / 60)}h
					</span>
					<div
						class="w-full rounded-t bg-linear-to-t from-logo-3 to-logo-1 transition-all hover:opacity-80"
						title="{monthLabels[i]}: {minutes} minutes"
						style="height: {Math.max((minutes / maxMonthly) * 100, 4)}%"
					></div>
					<span class="text-[10px] text-muted-foreground">
						{monthLabels[i]}
					</span>
				</div>
			{/each}
		</div>
	</section>
</div>