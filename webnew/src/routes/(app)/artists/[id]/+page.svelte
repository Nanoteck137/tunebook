<script lang="ts">
	import Image from "$lib/components/Image.svelte";
	import TrackListItem from "$lib/components/track-list/TrackListItem.svelte";
	import { Separator } from "$lib/components/ui";
	import { Disc, Music, TrendingUp, Users } from "@lucide/svelte";
	import SectionHeader from "$lib/components/SectionHeader.svelte";
	import TopTrackItem from "$lib/components/track-list/TopTrackItem.svelte";
	import Spacer from "$lib/components/Spacer.svelte";

	const { data } = $props();

	let featuredCount = $derived(
		data.featuredAlbumPage.totalItems + data.featuredTrackPage.totalItems,
	);

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

			<div class="flex flex-col">
				{#each data.tracks as track, i (track.id)}
					<TrackListItem {track} />
					{#if i < data.tracks.length - 1}
						<Separator />
					{/if}
				{/each}
			</div>
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

			<div
				class="grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5"
			>
				{#each data.albums as album (album.id)}
					<a
						href="/albums/{album.id}"
						class="group flex flex-col"
						title={album.name}
					>
						<div class="relative overflow-hidden rounded-lg">
							<Image
								class="aspect-square w-full rounded-none transition-transform duration-300 group-hover:scale-105"
								src={album.coverArt.medium}
								alt={album.name}
							/>
						</div>
						<div class="flex flex-col gap-0.5 pt-2">
							<p
								class="truncate text-sm font-medium group-hover:underline"
								title={album.name}
							>
								{album.name}
							</p>
							<p
								class="truncate text-xs text-muted-foreground"
								title={album.artists.map((a) => a.name).join(", ")}
							>
								{#if album.year}
									{album.year}
								{:else}
									{album.albumType}
								{/if}
							</p>
						</div>
					</a>
				{/each}
			</div>
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
							<a
								href="/albums/{album.id}"
								class="group flex flex-col"
								title={album.name}
							>
								<div class="relative overflow-hidden rounded-lg">
									<Image
										class="aspect-square w-full rounded-none transition-transform duration-300 group-hover:scale-105"
										src={album.coverArt.medium}
										alt={album.name}
									/>
								</div>
								<div class="flex flex-col gap-0.5 pt-2">
									<p
										class="truncate text-sm font-medium group-hover:underline"
										title={album.name}
									>
										{album.name}
									</p>
									<p
										class="truncate text-xs text-muted-foreground"
										title={otherArtists(album.artists).join(", ")}
									>
										{otherArtists(album.artists).join(", ")}
									</p>
								</div>
							</a>
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
