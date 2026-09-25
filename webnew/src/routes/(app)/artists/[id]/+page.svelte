<script lang="ts">
	import Image from "$lib/components/Image.svelte";
	import TrackListItem from "$lib/components/track-list/TrackListItem.svelte";
	import { Separator } from "$lib/components/ui";
	import { Disc, Music, Play, TrendingUp, Users } from "@lucide/svelte";
	import SectionHeader from "$lib/components/SectionHeader.svelte";
	import TopTrackItem from "$lib/components/track-list/TopTrackItem.svelte";
	import Spacer from "$lib/components/Spacer.svelte";
	import TileGrid from "$lib/components/tiles/TileGrid.svelte";
	import AlbumTile from "$lib/components/tiles/AlbumTile.svelte";
	import TrackList from "$lib/components/track-list/TrackList.svelte";
	import { getMusicManager } from "$lib/music-manager.svelte";

	const { data } = $props();
	const musicManager = getMusicManager();

	let featuredCount = $derived(
		data.featuredAlbumPage.totalItems + data.featuredTrackPage.totalItems,
	);

	async function playAlbum(albumId: string) {
		await musicManager.queueRequest({ type: "addAlbum", albumId }, {});
	}

	function otherArtists(artists: { id: string; name: string }[]) {
		return artists.filter((a) => a.id !== data.artist.id).map((a) => a.name);
	}
</script>

<div class="flex flex-col gap-10">
	{#if data.userTopTracks.length > 0}
		<section>
			<SectionHeader
				count={data.userTopTrackPage?.totalItems ?? data.userTopTracks.length}
				viewAllHref={data.user
					? `/users/${data.user.id}/review/total/top-tracks`
					: undefined}
			>
				<TrendingUp />
				Your Top Tracks
			</SectionHeader>

			<Spacer />

			<div class="flex flex-col">
				{#each data.userTopTracks as track (track.id)}
					<TopTrackItem
						rank={track.rank}
						{track}
						playCount={track.playCount}
					/>
				{/each}
			</div>
		</section>
	{/if}

	{#if data.tracks.length > 0}
		<section>
			<SectionHeader
				count={data.trackPage.totalItems}
				viewAllHref={data.trackPage.totalItems > data.tracks.length
					? `/artists/${data.artist.id}/tracks`
					: undefined}
			>
				<Music />
				Songs
			</SectionHeader>

			<Spacer />

			<TrackList
				tracks={data.tracks}
				onPlay={async (trackId, shuffle) => {
					await musicManager.queueRequest(
						{ type: "addArtist", artistId: data.artist.id },
						{ queueIndexToTrackId: trackId, shuffle },
					);
				}}
			/>
		</section>
	{/if}

	{#if data.albums.length > 0}
		<section>
			<SectionHeader
				count={data.albumPage.totalItems}
				viewAllHref={data.albumPage.totalItems > data.albums.length
					? `/artists/${data.artist.id}/albums`
					: undefined}
			>
				<Disc />
				Albums
			</SectionHeader>

			<Spacer />

			<TileGrid>
				{#each data.albums as album (album.id)}
					<AlbumTile {album} />
				{/each}
			</TileGrid>
		</section>
	{/if}

	{#if featuredCount > 0}
		<section>
			<SectionHeader count={featuredCount}>
				<Users />
				Appears on
			</SectionHeader>

			<Spacer />

			<div class="flex flex-col gap-2">
				{#if data.featuredAlbums.length > 0}
					<div
						class="grid grid-cols-2 gap-3 px-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5"
					>
						{#each data.featuredAlbums as album (album.id)}
							<div class="group flex flex-col">
								<div class="relative overflow-hidden rounded-lg">
									<a
										href="/albums/{album.id}"
										class="block overflow-hidden rounded-lg"
										title={album.name}
									>
										<Image
											class="aspect-square w-full rounded-none transition-transform duration-300 group-hover:scale-105"
											src={album.coverArt.medium}
											alt={album.name}
										/>
									</a>

									<button
										class="absolute right-2 bottom-2 hidden h-10 w-10 translate-y-2 items-center justify-center rounded-full bg-primary text-primary-foreground opacity-0 shadow-lg transition-all duration-300 group-hover:translate-y-0 group-hover:scale-105 group-hover:opacity-100 hover:scale-110 sm:flex"
										title={`Play ${album.name}`}
										aria-label={`Play ${album.name}`}
										onclick={() => playAlbum(album.id)}
									>
										<Play size={18} />
									</button>
								</div>
								<div class="flex flex-col gap-0.5 pt-2">
									<a
										class="truncate text-sm font-medium group-hover:underline"
										href="/albums/{album.id}"
										title={album.name}
									>
										{album.name}
									</a>
									<p
										class="truncate text-xs text-muted-foreground"
										title={otherArtists(album.artists).join(", ")}
									>
										{otherArtists(album.artists).join(", ")}
									</p>
								</div>
							</div>
						{/each}
					</div>
				{/if}

				{#if data.featuredTracks.length > 0}
					<div class="flex flex-col">
						{#each data.featuredTracks as track, i (track.id)}
							<TrackListItem {track} />
							{#if i < data.featuredTracks.length - 1}
								<Separator />
							{/if}
						{/each}
					</div>
				{/if}
			</div>
		</section>
	{/if}
</div>
