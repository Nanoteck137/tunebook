<script lang="ts">
	import { goto } from "$app/navigation";
	import {
		BarChart3,
		Clock,
		DiscAlbum,
		EllipsisVertical,
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
	import { DropdownMenu, Separator, buttonVariants } from "$lib/components/ui";
	import { cn, formatPlayTime } from "$lib/utils";
	import { getFavorites } from "$lib/favorites.svelte";
	import { getQuickPlaylist } from "$lib/quick-playlist.svelte";
	import SectionHeader from "$lib/components/SectionHeader.svelte";
	import { getMusicManager } from "$lib/music-manager.svelte";
	import { toast } from "svelte-sonner";

	let { data } = $props();
	const favoritesManager = getFavorites();
	const quickPlaylistManager = getQuickPlaylist();
	const musicManager = getMusicManager();

	let stats = $derived([
		{
			icon: Play,
			label: "Tracks Played",
			value: data.stats.numTracksPlayed.toLocaleString(),
		},
		{
			icon: Clock,
			label: "Listening Time",
			value: formatListeningTime(data.stats.listeningTime),
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

	// let maxTrackCount = $derived(
	// 	Math.max(...data.yearStats.map((y) => y.trackCount), 0),
	// );

	function formatListeningTime(seconds: number): string {
		const hours = Math.floor(seconds / 3600);
		const minutes = Math.floor((seconds % 3600) / 60);
		return `${hours}h ${minutes}m`;
	}

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
				viewAllHref="/users/{data.userData.id}/top"
			>
				<TrendingUp />
				Top Tracks
			</SectionHeader>

			<div class="flex flex-col">
				{#each data.topTracks as track, i (track.id)}
					<TopTrackItem
						rank={i + 1}
						{track}
						playCount={Math.floor(Math.random() * 1000) + 100}
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
									goto(`/albums/${track.albumId}`);
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
										wasFav ? "Removed from favorites" : "Added to favorites",
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
							{#if quickPlaylistManager.playlist !== null}
								<DropdownMenu.Item
									onSelect={async () => {
										const wasIn = quickPlaylistManager.hasTrack(track.id);
										await quickPlaylistManager.toggleTrack(track.id);
										toast.success(
											wasIn
												? "Removed from quick playlist"
												: "Added to quick playlist",
										);
									}}
								>
									{#if quickPlaylistManager.hasTrack(track.id)}
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

			<div
				class="grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6"
			>
				{#each data.playlists as playlist (playlist.id)}
					<div class="group relative flex flex-col">
						<div class="relative">
							<a
								href="/playlists/{playlist.id}"
								class="block overflow-hidden rounded-lg"
							>
								<img
									src={playlist.coverArt.medium}
									alt={playlist.name}
									class="aspect-square w-full object-cover transition-transform duration-300 group-hover:scale-105"
								/>
							</a>

							<button
								class="absolute right-2 bottom-2 hidden h-10 w-10 translate-y-2 items-center justify-center rounded-full bg-primary text-primary-foreground opacity-0 shadow-lg transition-all duration-300 group-hover:translate-y-0 group-hover:scale-105 group-hover:opacity-100 hover:scale-110 sm:flex"
								title="Play playlist"
								aria-label={`Play ${playlist.name}`}
								onclick={() => playPlaylist(playlist.id)}
							>
								<Play size={18} />
							</button>
						</div>

						<div class="flex flex-col gap-0.5 pt-2">
							<div class="flex items-center gap-1">
								<a
									href="/playlists/{playlist.id}"
									class="min-w-0 flex-1 truncate text-sm font-medium hover:underline"
									title={playlist.name}
								>
									{playlist.name}
								</a>

								<DropdownMenu.Root>
									<DropdownMenu.Trigger
										class={cn(
											buttonVariants({ variant: "ghost", size: "icon-sm" }),
											"-mr-1 shrink-0 rounded-full text-muted-foreground",
										)}
										aria-label={`More options for ${playlist.name}`}
									>
										<EllipsisVertical size={14} />
									</DropdownMenu.Trigger>
									<DropdownMenu.Content align="end">
										<DropdownMenu.Group>
											<DropdownMenu.Item
												onSelect={() => playPlaylist(playlist.id)}
											>
												<Play size={14} />
												Play
											</DropdownMenu.Item>
											<DropdownMenu.Item
												onSelect={() => playPlaylist(playlist.id, true)}
											>
												<Shuffle size={14} />
												Shuffle play
											</DropdownMenu.Item>
										</DropdownMenu.Group>
									</DropdownMenu.Content>
								</DropdownMenu.Root>
							</div>

							<p class="truncate text-xs text-muted-foreground">
								{playlist.trackCount}
								{playlist.trackCount !== 1 ? "tracks" : "track"}
								{#if playlist.playTime > 0}
									&middot; {formatPlayTime(playlist.playTime)}
								{/if}
							</p>
						</div>
					</div>
				{/each}
			</div>
		</section>
	{/if}

	<!-- {#if data.yearStats.length > 0} -->
	<!-- 	<Separator /> -->
	<!---->
	<!-- 	<section> -->
	<!-- 		<SectionHeader -->
	<!-- 			count={data.yearStats.length} -->
	<!-- 			viewAllHref="/users/{data.userData.id}/review" -->
	<!-- 		> -->
	<!-- 			<BarChart3 /> -->
	<!-- 			Year in Review -->
	<!-- 		</SectionHeader> -->
	<!---->
	<!-- 		<div class="flex flex-col gap-1"> -->
	<!-- 			{#each data.yearStats as stat (stat.year)} -->
	<!-- 				<a -->
	<!-- 					href="/users/{data.userData.id}/review/{stat.year}" -->
	<!-- 					class="flex items-center gap-4 rounded-lg px-2 py-1.5 transition-colors hover:bg-accent hover:text-accent-foreground" -->
	<!-- 				> -->
	<!-- 					<span class="w-12 text-sm font-medium">{stat.year}</span> -->
	<!-- 					<div class="flex flex-1 flex-col gap-1"> -->
	<!-- 						<div -->
	<!-- 							class="flex items-center justify-between text-xs text-muted-foreground" -->
	<!-- 						> -->
	<!-- 							<span>{stat.trackCount.toLocaleString()} tracks</span> -->
	<!-- 							<span>{formatListeningTime(stat.listeningTime)}</span> -->
	<!-- 						</div> -->
	<!-- 						<div class="h-2 w-full rounded-full bg-muted"> -->
	<!-- 							<div -->
	<!-- 								class="h-2 rounded-full bg-linear-to-r from-logo-1 to-logo-3" -->
	<!-- 								style="width: {maxTrackCount > 0 -->
	<!-- 									? (stat.trackCount / maxTrackCount) * 100 -->
	<!-- 									: 0}%" -->
	<!-- 							></div> -->
	<!-- 						</div> -->
	<!-- 					</div> -->
	<!-- 				</a> -->
	<!-- 			{/each} -->
	<!-- 		</div> -->
	<!-- 	</section> -->
	<!-- {/if} -->
</div>
