<script lang="ts">
	import { goto } from "$app/navigation";
	import {
		ChartColumn,
		Clock,
		DiscAlbum,
		Heart,
		ListMusic,
		ListPlus,
		Play,
		Shuffle,
		Star,
		TrendingUp,
		User,
	} from "@lucide/svelte";
	import TopTrackItem from "$lib/components/track-list/TopTrackItem.svelte";
	import { DropdownMenu } from "$lib/components/ui";
	import { formatPlayTime } from "$lib/utils";
	import { getFavorites } from "$lib/favorites.svelte";
	import { getQuickPlaylist } from "$lib/quick-playlist.svelte";
	import SectionHeader from "$lib/components/SectionHeader.svelte";
	import { getMusicManager } from "$lib/music-manager.svelte";
	import { toast } from "svelte-sonner";
	import Spacer from "$lib/components/Spacer.svelte";
	import TileGrid from "$lib/components/tiles/TileGrid.svelte";
	import PlaylistTile from "$lib/components/tiles/PlaylistTile.svelte";

	const RECENT_REVIEW_YEARS = 5;

	let { data } = $props();
	const favoritesManager = getFavorites();
	const quickPlaylistManager = getQuickPlaylist();
	const musicManager = getMusicManager();

	const currentYear = new Date().getFullYear();

	let recentReviews = $derived(data.reviews.slice(0, RECENT_REVIEW_YEARS));

	let maxReviewTrackCount = $derived(
		Math.max(...recentReviews.map((r) => r.trackCount), 0),
	);

	let stats = $derived([
		{
			icon: Play,
			label: "Tracks Played",
			value: data.stats.numTracksPlayed.toLocaleString(),
		},
		{
			icon: Clock,
			label: "Listening Time",
			value: formatPlayTime(data.stats.listeningTime),
		},
		{
			icon: ListMusic,
			label: "Playlists",
			value: data.stats.numPlaylistsCreated.toLocaleString(),
		},
		{
			icon: Heart,
			label: "Favorites",
			value: data.stats.numFavoriteTracks.toLocaleString(),
		},
	]);

	function playPlaylist(playlistId: string, shuffle = false) {
		return musicManager.queueRequest(
			{ type: "addPlaylist", playlistId },
			{ shuffle },
		);
	}
</script>

<div class="flex flex-col gap-10">
	<div class="grid grid-cols-2 gap-4 lg:grid-cols-4">
		{#each stats as stat (stat.label)}
			<div class="flex flex-col gap-1.5 rounded-lg border bg-card p-4">
				<div class="flex items-center gap-2">
					<stat.icon size={16} class="text-muted-foreground" />
					<span class="text-xs text-muted-foreground">{stat.label}</span>
				</div>
				<span class="text-2xl font-bold">{stat.value}</span>
			</div>
		{/each}
	</div>

	{#if data.topTracks.length > 0}
		<section>
			<SectionHeader
				count={data.topTracks.length}
				viewAllHref="/users/{data.userData.id}/total/top-tracks"
			>
				<TrendingUp />
				Top Tracks
			</SectionHeader>

			<Spacer />

			<div class="flex flex-col">
				{#each data.topTracks as item (item.id)}
					<TopTrackItem
						rank={item.rank}
						track={item}
						playCount={item.playCount}
					>
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
									goto(`/albums/${item.albumId}`);
								}}
							>
								<DiscAlbum />
								Go to Album
							</DropdownMenu.Item>
							<DropdownMenu.Sub>
								<DropdownMenu.SubTrigger>
									<User />
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
										wasFav ? "Removed from favorites" : "Added to favorites",
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
							{#if quickPlaylistManager.playlist !== null}
								<DropdownMenu.Item
									onSelect={async () => {
										const wasIn = quickPlaylistManager.hasTrack(item.id);
										await quickPlaylistManager.toggleTrack(item.id);
										toast.success(
											wasIn
												? "Removed from quick playlist"
												: "Added to quick playlist",
										);
									}}
								>
									{#if quickPlaylistManager.hasTrack(item.id)}
										<Star class="fill-primary stroke-primary" />
										Remove from Quick
									{:else}
										<Star />
										Quick Add
									{/if}
								</DropdownMenu.Item>
							{/if}
						{/snippet}
					</TopTrackItem>
				{/each}
			</div>
		</section>
	{/if}

	{#if data.playlists.length > 0}
		<section>
			<SectionHeader
				count={data.playlists.length}
				viewAllHref="/users/{data.userData.id}/playlists"
			>
				<ListMusic />
				Playlists
			</SectionHeader>

			<Spacer />

			<TileGrid>
				{#each data.playlists as playlist (playlist.id)}
					<PlaylistTile {playlist} />
				{/each}
			</TileGrid>
		</section>
	{/if}

	{#if recentReviews.length > 0}
		<section>
			<SectionHeader viewAllHref="/users/{data.userData.id}/review">
				<ChartColumn />
				Year in Review
			</SectionHeader>

			<Spacer />

			<div class="flex flex-col">
				{#each recentReviews as review (review.year)}
					<a
						href="/users/{data.userData.id}/review/{review.year}"
						class="flex items-center gap-4 rounded-lg px-2 py-2 transition-colors hover:bg-accent hover:text-accent-foreground"
					>
						<span class="w-12 shrink-0 text-sm font-semibold tabular-nums">
							{review.year}
						</span>

						<div class="flex min-w-0 flex-1 flex-col gap-1.5">
							<div
								class="flex items-center justify-between gap-2 text-xs text-muted-foreground"
							>
								<span class="truncate">
									{review.trackCount.toLocaleString()} tracks
								</span>
								<span class="flex shrink-0 items-center gap-1.5">
									<span class="tabular-nums">
										{formatPlayTime(review.listeningTime)}
									</span>
									{#if review.year === currentYear}
										<span
											class="rounded-full bg-primary/10 px-2 py-0.5 font-medium text-primary ring-1 ring-primary/25"
										>
											In progress
										</span>
									{/if}
								</span>
							</div>

							<div class="h-1.5 w-full overflow-hidden rounded-full bg-muted">
								<div
									class="h-full rounded-full bg-linear-to-r from-logo-1 to-logo-3"
									style="width: {maxReviewTrackCount > 0
										? (review.trackCount / maxReviewTrackCount) * 100
										: 0}%"
								></div>
							</div>
						</div>
					</a>
				{/each}
			</div>
		</section>
	{/if}
</div>
