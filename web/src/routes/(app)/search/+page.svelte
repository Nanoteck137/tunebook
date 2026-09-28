<script lang="ts">
	import { goto } from "$app/navigation";
	import { Disc, Mic, Music } from "@lucide/svelte";
	import { onMount } from "svelte";
	import TrackList from "$lib/components/track-list/TrackList.svelte";
	import { getMusicManager } from "$lib/music-manager.svelte";
	import AlbumTile from "$lib/components/tiles/AlbumTile.svelte";
	import ArtistTile from "$lib/components/tiles/ArtistTile.svelte";
	import TileGrid from "$lib/components/tiles/TileGrid.svelte";
	import SearchBarHeader from "./SearchBarHeader.svelte";
	import SectionHeader from "$lib/components/SectionHeader.svelte";
	import Spacer from "$lib/components/Spacer.svelte";

	const { data } = $props();
	const musicManager = getMusicManager();

	async function search(query: string) {
		await goto(`?query=${query}`, {
			invalidateAll: true,
			keepFocus: true,
			replaceState: true,
		});
	}

	function clearSearch() {
		value = "";
		search("");
	}

	let value = $state("");

	onMount(() => {
		value = data.query;
	});

	function formatError(err: { type: string; code: number; message: string }) {
		return err.message;
	}

	let topTracks = $derived(data.tracks.slice(0, 6));

	async function playTrack(trackId: string) {
		await musicManager.addTracks({
			trackIds: topTracks.map((t) => t.id),
			trackId,
			clear: true,
		});
	}
</script>

<svelte:head>
	<title>Search - Tunebook</title>
</svelte:head>

<div class="flex flex-col gap-6">
	<SearchBarHeader
		searchBarPlaceholder="Search tracks, artists, albums..."
		{value}
		setValue={(v) => {
			value = v;
		}}
		{search}
		searchWithValue={() => {
			search(value);
		}}
		{clearSearch}
	/>

	{#if data.artistError || data.albumError || data.trackError}
		<div class="flex flex-col gap-1 text-sm text-red-400">
			{#if data.artistError}
				<p>Artists: {formatError(data.artistError)}</p>
			{/if}
			{#if data.albumError}
				<p>Albums: {formatError(data.albumError)}</p>
			{/if}
			{#if data.trackError}
				<p>Tracks: {formatError(data.trackError)}</p>
			{/if}
		</div>
	{/if}

	{#if data.query && data.artists.length === 0 && data.albums.length === 0 && data.tracks.length === 0}
		<p class="py-12 text-center text-sm text-muted-foreground">
			No results found for "{data.query}".
		</p>
	{/if}

	{#if topTracks.length > 0}
		<section>
			<SectionHeader
				count={data.tracks.length}
				viewAllHref="/search/tracks?query={data.query}"
			>
				<Music />
				Tracks
			</SectionHeader>

			<Spacer />

			<TrackList
				tracks={topTracks}
				onPlay={playTrack}
			/>
		</section>
	{/if}

	{#if data.artists.length > 0}
		<section>
			<SectionHeader
				count={data.artists.length}
				viewAllHref="/search/artists?query={data.query}"
			>
				<Mic />
				Artists
			</SectionHeader>

			<Spacer />

			<TileGrid>
				{#each data.artists.slice(0, 6) as artist (artist.id)}
					<ArtistTile {artist} />
				{/each}
			</TileGrid>
		</section>
	{/if}

	{#if data.albums.length > 0}
		<section>
			<SectionHeader
				count={data.albums.length}
				viewAllHref="/search/albums?query={data.query}"
			>
				<Disc />
				Albums
			</SectionHeader>

			<Spacer />

			<TileGrid>
				{#each data.albums.slice(0, 6) as album (album.id)}
					<AlbumTile {album} />
				{/each}
			</TileGrid>
		</section>
	{/if}
</div>
