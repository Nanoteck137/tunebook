<script lang="ts">
	import { goto } from "$app/navigation";
	import { page } from "$app/state";
	import { Button, InputGroup, Separator } from "$lib/components/ui";
	import { SearchIcon, XIcon } from "@lucide/svelte";
	import { onMount } from "svelte";
	import { cn } from "$lib/utils";
	import TrackList from "$lib/components/track-list/TrackList.svelte";
	import Image from "$lib/components/Image.svelte";
	import SearchBarHeader from "./SearchBarHeader.svelte";

	const { data } = $props();

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
</script>

<svelte:head>
	<title>Search - Tunebook</title>
</svelte:head>

<div class="flex flex-col gap-4">
	<SearchBarHeader
		searchBarPlaceholder="Search artists, albums, tracks, playlists, users..."
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

	<!-- <form -->
	<!-- 	action="" -->
	<!-- 	method="get" -->
	<!-- 	onsubmit={(e) => { -->
	<!-- 		e.preventDefault(); -->
	<!-- 		clearTimeout(timer); -->
	<!-- 		search(value); -->
	<!-- 	}} -->
	<!-- > -->
	<!-- 	<InputGroup.Root> -->
	<!-- 		<InputGroup.Input -->
	<!-- 			id="query" -->
	<!-- 			name="query" -->
	<!-- 			placeholder="Search artists, albums, tracks, playlists, users..." -->
	<!-- 			autocomplete="off" -->
	<!-- 			bind:value={initialValue} -->
	<!-- 			oninput={onInput} -->
	<!-- 		/> -->
	<!-- 		<InputGroup.Addon> -->
	<!-- 			<SearchIcon /> -->
	<!-- 		</InputGroup.Addon> -->
	<!-- 		<InputGroup.Addon align="inline-end"> -->
	<!-- 			<InputGroup.Button type="submit"> -->
	<!-- 				<SearchIcon /> -->
	<!-- 			</InputGroup.Button> -->
	<!---->
	<!-- 			<Separator class="min-h-4" orientation="vertical" /> -->
	<!---->
	<!-- 			<InputGroup.Button -->
	<!-- 				onclick={() => { -->
	<!-- 					const e = document.getElementById("query") as HTMLInputElement; -->
	<!-- 					e.value = ""; -->
	<!-- 					clearSearch(); -->
	<!-- 				}} -->
	<!-- 			> -->
	<!-- 				<XIcon /> -->
	<!-- 			</InputGroup.Button> -->
	<!-- 		</InputGroup.Addon> -->
	<!-- 	</InputGroup.Root> -->
	<!-- </form> -->
	<!---->
	<!-- <nav class="flex flex-wrap gap-1"> -->
	<!-- 	{#each tabs as { label, href }} -->
	<!-- 		<Button -->
	<!-- 			variant={page.url.pathname === href ? "default" : "outline"} -->
	<!-- 			href="{href}?query={data.query}" -->
	<!-- 		> -->
	<!-- 			{label} -->
	<!-- 		</Button> -->
	<!-- 	{/each} -->
	<!-- </nav> -->

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

	{#if data.query && data.artists.length === 0 && data.albums.length === 0 && data.tracks.length === 0 && data.playlists.length === 0 && data.users.length === 0}
		<p class="py-12 text-center text-sm text-muted-foreground">
			No results found for "{data.query}".
		</p>
	{/if}

	{#if data.tracks.length > 0}
		<div>
			<div class="mb-3 flex items-baseline justify-between">
				<h2 class="text-base font-semibold">Tracks</h2>
				<div class="flex items-center gap-2">
					<span class="text-xs text-muted-foreground"
						>{data.tracks.length}</span
					>
					<a
						href="/search/tracks?query={data.query}"
						class="text-xs font-medium text-primary hover:underline"
					>
						View all
					</a>
				</div>
			</div>

			<TrackList
				totalTracks={data.tracks.length}
				tracks={data.tracks}
				onPlay={() => {}}
			/>
		</div>
	{/if}

	{#if data.artists.length > 0}
		<div>
			<div class="mb-3 flex items-baseline justify-between">
				<h2 class="text-base font-semibold">Artists</h2>
				<div class="flex items-center gap-2">
					<span class="text-xs text-muted-foreground"
						>{data.artists.length}</span
					>
					<a
						href="/search/artists?query={data.query}"
						class="text-xs font-medium text-primary hover:underline"
					>
						View all
					</a>
				</div>
			</div>

			<div
				class="grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6"
			>
				{#each data.artists.slice(0, 6) as artist (artist.id)}
					<a
						href="/artists/{artist.id}"
						class="flex flex-col overflow-hidden rounded-lg border bg-card transition-shadow hover:shadow-md"
					>
						<img
							src={artist.coverArt.medium}
							alt={artist.name}
							class="aspect-square w-full object-cover"
						/>
						<div class="p-2">
							<p class="truncate text-sm font-medium" title={artist.name}>
								{artist.name}
							</p>
						</div>
					</a>
				{/each}
			</div>
		</div>
	{/if}

	{#if data.albums.length > 0}
		<div>
			<div class="mb-3 flex items-baseline justify-between">
				<h2 class="text-base font-semibold">Albums</h2>
				<div class="flex items-center gap-2">
					<span class="text-xs text-muted-foreground"
						>{data.albums.length}</span
					>
					<a
						href="/search/albums?query={data.query}"
						class="text-xs font-medium text-primary hover:underline"
					>
						View all
					</a>
				</div>
			</div>

			<div
				class="grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6"
			>
				{#each data.albums.slice(0, 6) as album (album.id)}
					<a
						href="/albums/{album.id}"
						class="flex flex-col overflow-hidden rounded-lg border bg-card transition-shadow hover:shadow-md"
					>
						<img
							src={album.coverArt.medium}
							alt={album.name}
							class="aspect-square w-full object-cover"
						/>
						<div class="p-2">
							<p class="truncate text-sm font-medium" title={album.name}>
								{album.name}
							</p>
							<p
								class="truncate text-xs text-muted-foreground"
								title={album.artists.map((a) => a.name).join(", ")}
							>
								{album.artists.map((a) => a.name).join(", ")}
							</p>
						</div>
					</a>
				{/each}
			</div>
		</div>
	{/if}

	{#if data.playlists.length > 0}
		<div>
			<div class="mb-3 flex items-baseline justify-between">
				<h2 class="text-base font-semibold">Playlists</h2>
				<div class="flex items-center gap-2">
					<span class="text-xs text-muted-foreground"
						>{data.playlists.length}</span
					>
					<a
						href="/search/playlists?query={data.query}"
						class="text-xs font-medium text-primary hover:underline"
					>
						View all
					</a>
				</div>
			</div>

			<div
				class="grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6"
			>
				{#each data.playlists.slice(0, 6) as playlist (playlist.id)}
					<a
						href="/playlists/{playlist.id}"
						class="flex flex-col overflow-hidden rounded-lg border bg-card transition-shadow hover:shadow-md"
					>
						<Image
							class="aspect-square w-full rounded-none border-0"
							src={playlist.coverArt.medium}
							alt={playlist.name}
						/>
						<div class="p-2">
							<p class="truncate text-sm font-medium" title={playlist.name}>
								{playlist.name}
							</p>
							<p class="truncate text-xs text-muted-foreground">
								{playlist.ownerDisplayName}
							</p>
						</div>
					</a>
				{/each}
			</div>
		</div>
	{/if}

	{#if data.users.length > 0}
		<div>
			<div class="mb-3 flex items-baseline justify-between">
				<h2 class="text-base font-semibold">Users</h2>
				<div class="flex items-center gap-2">
					<span class="text-xs text-muted-foreground">{data.users.length}</span
					>
					<a
						href="/search/users?query={data.query}"
						class="text-xs font-medium text-primary hover:underline"
					>
						View all
					</a>
				</div>
			</div>

			<div
				class="grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6"
			>
				{#each data.users.slice(0, 6) as user (user.id)}
					<a
						href="/users/{user.id}"
						class="flex flex-col overflow-hidden rounded-lg border bg-card transition-shadow hover:shadow-md"
					>
						<img
							src={user.picture.medium}
							alt={user.displayName}
							class="aspect-square w-full object-cover"
						/>
						<div class="p-2">
							<p class="truncate text-sm font-medium" title={user.displayName}>
								{user.displayName}
							</p>
							<p class="truncate text-xs text-muted-foreground">{user.role}</p>
						</div>
					</a>
				{/each}
			</div>
		</div>
	{/if}
</div>
