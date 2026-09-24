<script lang="ts">
	import { getMusicManager } from "$lib/music-manager.svelte";
	import PlaylistSkeletonTile from "$lib/components/tiles/PlaylistSkeletonTile.svelte";
	import PlaylistTile from "$lib/components/tiles/PlaylistTile.svelte";
	import TrackSkeletonTile from "$lib/components/tiles/TrackSkeletonTile.svelte";
	import TrackTile from "$lib/components/tiles/TrackTile.svelte";
	import SectionHeader from "$lib/components/SectionHeader.svelte";
	import { Button } from "$lib/components/ui";
	import {
		Clock,
		Compass,
		Disc3,
		Flame,
		Heart,
		History,
		Library,
		ListMusic,
		Pause,
		Play,
		Search,
		Server,
		Smartphone,
	} from "@lucide/svelte";
	import HeroCard from "$lib/components/HeroCard.svelte";
	import SectionImage from "$lib/components/SectionImage.svelte";

	let { data } = $props();

	const musicManager = getMusicManager();

	function greeting() {
		const h = new Date().getHours();
		if (h < 5) return "Working late";
		if (h < 12) return "Good morning";
		if (h < 17) return "Good afternoon";
		if (h < 21) return "Good evening";
		return "Good night";
	}

	function formatListeningTime(seconds: number): string {
		const hours = Math.floor(seconds / 3600);
		const minutes = Math.floor((seconds % 3600) / 60);
		if (hours > 0) return `${hours}h ${minutes}m`;
		return `${minutes}m`;
	}

	function timeAgo(millis: number): string {
		const diff = Date.now() - millis;
		const minutes = Math.floor(diff / 60000);
		if (minutes < 1) return "just now";
		if (minutes < 60) return `${minutes}m ago`;
		const hours = Math.floor(minutes / 60);
		if (hours < 24) return `${hours}h ago`;
		const days = Math.floor(hours / 24);
		return `${days}d ago`;
	}

	const quickActions = [
		{
			label: "Browse",
			icon: Compass,
			href: "/browse",
		},
		{
			label: "Search",
			icon: Search,
			href: "/search",
		},
		{
			label: "Playlists",
			icon: ListMusic,
			href: "/library/playlists",
		},
		{
			label: "Favorites",
			icon: Heart,
			href: "/library/favorites",
		},
		{
			label: "Library",
			icon: Library,
			href: "/library",
		},
	];
</script>

{#if !data.user}
	<div class="flex min-h-[70dvh] flex-col items-center justify-center p-4">
		<div
			class="w-full max-w-2xl rounded-2xl border bg-card p-8 text-center shadow-sm sm:p-12"
		>
			<div class="relative flex flex-col items-center gap-5">
				<div
					class="flex h-20 w-20 items-center justify-center rounded-full bg-linear-to-tr from-logo-1 via-logo-2 to-logo-3 shadow-lg"
				>
					<svg
						xmlns="http://www.w3.org/2000/svg"
						class="h-10 w-10 text-white"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="2"
						stroke-linecap="round"
						stroke-linejoin="round"
					>
						<circle cx="12" cy="12" r="10" /><circle cx="12" cy="12" r="3" />
					</svg>
				</div>

				<div class="flex flex-col items-center gap-2">
					<h1
						class="bg-linear-to-tr from-logo-1 via-logo-2 to-logo-3 bg-clip-text text-5xl font-bold text-transparent"
					>
						Tunebook
					</h1>
					<p class="text-base text-muted-foreground">
						Your personal music streaming server
					</p>
				</div>

				<div
					class="grid w-full max-w-md grid-cols-1 gap-2 text-left sm:grid-cols-3"
				>
					<div
						class="flex items-center gap-2 rounded-lg border p-3 text-sm text-foreground"
					>
						<Server class="h-4 w-4 shrink-0 text-primary" />
						<span>Self-host your library</span>
					</div>
					<div
						class="flex items-center gap-2 rounded-lg border p-3 text-sm text-foreground"
					>
						<Heart class="h-4 w-4 shrink-0 text-primary" />
						<span>Favorites & playlists</span>
					</div>
					<div
						class="flex items-center gap-2 rounded-lg border p-3 text-sm text-foreground"
					>
						<Smartphone class="h-4 w-4 shrink-0 text-primary" />
						<span>Play on any device</span>
					</div>
				</div>

				<Button size="lg" class="mt-2" href="/login">Login</Button>
			</div>
		</div>
	</div>
{:else}
	<div class="flex flex-col gap-8">
		<HeroCard
			class="section-home"
			innerClass="sm:flex-row sm:items-center sm:justify-between"
		>
			<div class="flex min-w-0 flex-col gap-4">
				<div class="flex items-center gap-4">
					<a href="/users/{data.user.id}" title={data.user.displayName}>
						<SectionImage
							variant="full"
							src={data.user.picture.small}
							alt={data.user.displayName}
							class="aspect-square w-16 rounded-full p-0.5 md:w-28"
							imgClass="rounded-full"
						/>
					</a>

					<div class="flex min-w-0 flex-col gap-0.5">
						<p
							class="text-xs font-semibold tracking-wider text-muted-foreground uppercase"
						>
							{greeting()}
						</p>
						<h1 class="line-clamp-1 text-xl font-bold sm:text-2xl">
							<a
								href="/users/{data.user.id}"
								class="hover:underline"
								title={data.user.displayName}>{data.user.displayName}</a
							>
						</h1>
					</div>
				</div>

				<p class="text-sm text-muted-foreground">
					{#await data.stats}
						Your personal music streaming server.
					{:then stats}
						{#if stats}
							You've played {stats.numTracksPlayed} track(s) for
							{formatListeningTime(stats.listeningTime)} total
							{#if stats.lastListenedAt}
								&middot; last listened {timeAgo(stats.lastListenedAt)}
							{/if}
							.
						{:else}
							Your personal music streaming server.
						{/if}
					{/await}
				</p>
			</div>

			<div
				class="flex w-full flex-col gap-2 sm:w-auto sm:flex-row sm:items-center"
			>
				{#if musicManager.currentItem}
					<Button
						class="w-full sm:w-auto"
						onclick={() => {
							if (musicManager.playing) {
								musicManager.pause();
							} else {
								musicManager.play();
							}
						}}
					>
						{#if musicManager.playing}
							<Pause />
							Pause
						{:else}
							<Play />
							Resume
						{/if}
					</Button>
				{/if}

				<Button variant="outline" class="w-full sm:w-auto" href="/browse">
					<Compass />
					Browse
				</Button>
			</div>
		</HeroCard>

		{#if musicManager.currentItem}
			<section>
				<div class="grid grid-cols-3 gap-2 sm:grid-cols-6">
					<button
						class="flex flex-col items-center gap-2 rounded-lg border bg-card p-4 transition-colors hover:bg-accent"
						onclick={() => {
							if (musicManager.playing) {
								musicManager.pause();
							} else {
								musicManager.play();
							}
						}}
					>
						{#if musicManager.playing}
							<Pause class="h-5 w-5" />
						{:else}
							<Play class="h-5 w-5" />
						{/if}
						<span class="text-xs font-medium">
							{musicManager.playing ? "Pause" : "Resume"}
						</span>
					</button>

					{#each quickActions as action (action.href)}
						<a
							class="flex flex-col items-center gap-2 rounded-lg border bg-card p-4 transition-colors hover:bg-accent"
							href={action.href}
						>
							<action.icon class="h-5 w-5" />
							<span class="text-xs font-medium">{action.label}</span>
						</a>
					{/each}
				</div>
			</section>
		{/if}

		{#await data.stats}
			<section>
				<div class="grid grid-cols-2 gap-2 sm:grid-cols-4">
					{#each Array(4) as _}
						<div class="h-28 animate-pulse rounded-lg border bg-card"></div>
					{/each}
				</div>
			</section>
		{:then stats}
			{#if stats}
				<section>
					<div class="grid grid-cols-2 gap-2 sm:grid-cols-4">
						<div class="flex flex-col gap-1 rounded-lg border bg-card p-4">
							<div class="flex items-center gap-2 text-muted-foreground">
								<Disc3 class="h-4 w-4" />
								<span class="text-xs tracking-wider uppercase">
									Tracks played
								</span>
							</div>
							<p class="text-2xl font-bold">
								{stats.numTracksPlayed}
							</p>
							<p class="text-xs text-muted-foreground">
								{stats.numTracksSkipped} skipped
							</p>
						</div>

						<div class="flex flex-col gap-1 rounded-lg border bg-card p-4">
							<div class="flex items-center gap-2 text-muted-foreground">
								<Clock class="h-4 w-4" />
								<span class="text-xs tracking-wider uppercase">
									Listening time
								</span>
							</div>
							<p class="text-2xl font-bold">
								{formatListeningTime(stats.listeningTime)}
							</p>
							<p class="text-xs text-muted-foreground">
								{#if stats.lastListenedAt}
									{timeAgo(stats.lastListenedAt)}
								{:else}
									No plays yet
								{/if}
							</p>
						</div>

						<a
							class="flex flex-col gap-1 rounded-lg border bg-card p-4 transition-colors hover:bg-accent"
							href="/library/favorites"
						>
							<div class="flex items-center gap-2 text-muted-foreground">
								<Heart class="h-4 w-4" />
								<span class="text-xs tracking-wider uppercase">
									Favorite tracks
								</span>
							</div>
							<p class="text-2xl font-bold">
								{stats.numFavoriteTracks}
							</p>
							<p class="text-xs text-muted-foreground">View favorites</p>
						</a>

						<a
							class="flex flex-col gap-1 rounded-lg border bg-card p-4 transition-colors hover:bg-accent"
							href="/library/playlists"
						>
							<div class="flex items-center gap-2 text-muted-foreground">
								<ListMusic class="h-4 w-4" />
								<span class="text-xs tracking-wider uppercase">
									Playlists
								</span>
							</div>
							<p class="text-2xl font-bold">
								{stats.numPlaylistsCreated}
							</p>
							<p class="text-xs text-muted-foreground">View playlists</p>
						</a>
					</div>
				</section>
			{/if}
		{/await}

		<section>
			<SectionHeader viewAllHref="/users/{data.user.id}/history">
				<History />
				Recently Played
			</SectionHeader>

			{#await data.recentlyPlayed}
				<div class="flex gap-2 overflow-x-auto pb-4">
					{#each Array(5) as _i}
						<TrackSkeletonTile />
					{/each}
				</div>
			{:then tracks}
				{#if tracks.length > 0}
					<div class="flex gap-2 overflow-x-auto pb-4">
						{#each tracks as track (track.id)}
							<TrackTile size="sm" class="w-40" {track} />
						{/each}
					</div>
				{:else}
					<p class="px-2 text-sm text-muted-foreground">
						Nothing played yet — start listening!
					</p>
				{/if}
			{/await}
		</section>

		<section>
			<SectionHeader>
				<Flame />
				Top Tracks
			</SectionHeader>

			{#await data.topTracks}
				<div class="flex gap-2 overflow-x-auto pb-4">
					{#each Array(5) as _i}
						<TrackSkeletonTile />
					{/each}
				</div>
			{:then tracks}
				{#if tracks.length > 0}
					<div class="flex gap-2 overflow-x-auto pb-4">
						{#each tracks as track (track.id)}
							<TrackTile size="sm" class="w-40" {track} />
						{/each}
					</div>
				{:else}
					<p class="px-2 text-sm text-muted-foreground">
						Your most-played tracks will appear here.
					</p>
				{/if}
			{/await}
		</section>

		<section>
			<SectionHeader viewAllHref="/library/playlists">
				<ListMusic />
				Your Playlists
			</SectionHeader>

			{#await data.playlists}
				<div class="flex gap-2 overflow-x-auto pb-4">
					{#each Array(5) as _i}
						<PlaylistSkeletonTile />
					{/each}
				</div>
			{:then playlists}
				{#if playlists.length > 0}
					<div class="flex gap-2 overflow-x-auto pb-4">
						{#each playlists as playlist (playlist.id)}
							<PlaylistTile size="sm" class="w-40" {playlist} />
						{/each}
					</div>
				{:else}
					<p class="px-2 text-sm text-muted-foreground">No playlists yet</p>
				{/if}
			{/await}
		</section>

		<section>
			<SectionHeader viewAllHref="/library/favorites">
				<Heart />
				Favorites
			</SectionHeader>

			{#await data.favorites}
				<div class="flex gap-2 overflow-x-auto pb-4">
					{#each Array(5) as _i}
						<TrackSkeletonTile />
					{/each}
				</div>
			{:then favorites}
				{#if favorites.length > 0}
					<div class="flex gap-2 overflow-x-auto pb-4">
						{#each favorites as track (track.id)}
							<TrackTile size="sm" class="w-40" {track} />
						{/each}
					</div>
				{:else}
					<p class="px-2 text-sm text-muted-foreground">
						Heart some tracks to see them here
					</p>
				{/if}
			{/await}
		</section>
	</div>
{/if}
