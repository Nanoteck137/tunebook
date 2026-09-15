<script lang="ts">
	import { goto } from "$app/navigation";
	import {
		DiscAlbum,
		Heart,
		ListPlus,
		Play,
		Shuffle,
		Star,
		User,
	} from "@lucide/svelte";
	import TopTrackItem from "$lib/components/track-list/TopTrackItem.svelte";
	import { DropdownMenu } from "$lib/components/ui";
	import { getFavorites } from "$lib/favorites.svelte";
	import { getQuickPlaylist } from "$lib/quick-playlist.svelte";
	import { toast } from "svelte-sonner";

	let { data } = $props();
	const favoritesManager = getFavorites();
	const quickPlaylistManager = getQuickPlaylist();
</script>

<div class="flex flex-col gap-4">
	<div class="flex items-baseline justify-between">
		<div class="flex items-baseline gap-2">
			<h1 class="text-xl font-bold">Top Tracks</h1>
			<span class="text-sm text-muted-foreground">{data.topTracks.length}</span>
		</div>
	</div>
</div>

<div class="flex flex-col gap-1">
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
