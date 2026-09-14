<script lang="ts">
	import { getMusicManager } from "$lib/music-manager.svelte";
	import PlaylistSkeletonTile from "$lib/components/tiles/PlaylistSkeletonTile.svelte";
	import PlaylistTile from "$lib/components/tiles/PlaylistTile.svelte";
	import TrackSkeletonTile from "$lib/components/tiles/TrackSkeletonTile.svelte";
	import TrackTile from "$lib/components/tiles/TrackTile.svelte";
	import { Button } from "$lib/components/ui";
	import {
		ChevronRight,
		Clock,
		Compass,
		Disc3,
		Heart,
		Library,
		ListMusic,
		Pause,
		Play,
		Search,
		Server,
		Smartphone,
	} from "@lucide/svelte";
	import collage from "$lib/assets/collage.png";

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
	<div
		class="flex min-h-[70dvh] flex-col items-center justify-center p-4"
	>
		<div
			class="relative w-full max-w-2xl overflow-hidden rounded-2xl bg-linear-to-tr from-logo-1 via-logo-2 to-logo-3 p-8 text-center sm:p-12"
		>
			<div
				class="pointer-events-none absolute -top-20 -right-20 h-64 w-64 rounded-full bg-white/10 blur-2xl"
			></div>

			<img
				src={collage}
				alt=""
				class="pointer-events-none absolute inset-0 h-full w-full object-cover"
			/>

			<div
				class="relative flex flex-col items-center gap-5"
			>
				<div
					class="flex h-24 w-24 items-center justify-center rounded-full border border-white/30 bg-white/10 backdrop-blur-sm"
				>
					<svg
						xmlns="http://www.w3.org/2000/svg"
						class="h-12 w-12 text-white"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="2"
						stroke-linecap="round"
						stroke-linejoin="round"
					>
						<circle cx="12" cy="12" r="10" /><circle
							cx="12"
							cy="12"
							r="3"
						/>
					</svg>
				</div>

				<h1 class="text-5xl font-bold text-white">Tunebook</h1>
				<p class="text-base text-white/90">
					Your personal music streaming server
				</p>

				<div class="grid w-full max-w-md grid-cols-1 gap-2 text-left sm:grid-cols-3">
					<div
						class="flex items-center gap-2 rounded-lg border border-white/20 bg-black/25 p-3 text-sm text-white backdrop-blur-sm"
					>
						<Server class="h-4 w-4 shrink-0" />
						<span>Self-host your library</span>
					</div>
					<div
						class="flex items-center gap-2 rounded-lg border border-white/20 bg-black/25 p-3 text-sm text-white backdrop-blur-sm"
					>
						<Heart class="h-4 w-4 shrink-0" />
						<span>Favorites & playlists</span>
					</div>
					<div
						class="flex items-center gap-2 rounded-lg border border-white/20 bg-black/25 p-3 text-sm text-white backdrop-blur-sm"
					>
						<Smartphone class="h-4 w-4 shrink-0" />
						<span>Play on any device</span>
					</div>
				</div>

				<Button size="lg" class="mt-2" href="/login">Login</Button>
			</div>
		</div>
	</div>
{:else}
	<div class="flex flex-col gap-10">
		<!-- Hero -->
		<section
			class="relative flex justify-center overflow-hidden rounded-lg bg-linear-to-tr from-logo-1 via-logo-2 to-logo-3 p-6 sm:justify-start sm:p-10"
		>
			<div
				class="pointer-events-none absolute -top-16 -right-16 h-48 w-48 rounded-full bg-white/10 blur-2xl"
			></div>

			<img
				src={collage}
				alt=""
				class="pointer-events-none absolute inset-0 h-full w-full object-cover"
			/>

			<div
				class="relative flex min-w-0 flex-col items-center gap-4 rounded-lg border border-white/20 bg-black/30 p-4 text-center backdrop-blur-sm sm:items-start sm:p-6 sm:text-left"
			>
				<div class="flex items-center gap-3">
					<img
						class="h-14 w-14 rounded-full border-2 border-white/40"
						src={data.user.picture.small}
						alt={data.user.displayName}
					/>
					<div>
						<p class="text-sm font-medium uppercase tracking-widest text-white/80">
							{greeting()}
						</p>
						<h1 class="text-3xl font-bold text-white sm:text-4xl">
							{data.user.displayName}
						</h1>
					</div>
				</div>

				<p class="max-w-xl text-sm text-white/90">
					{#await data.stats}
					Your personal music streaming server.
					{:then stats}
						{#if stats}
							You've played {stats.numTracksPlayed} track(s) for
							{formatListeningTime(stats.listeningTime)} total
							{#if stats.lastListenedAt}
								· last listened {timeAgo(stats.lastListenedAt)}
							{/if}
							.
						{:else}
							Your personal music streaming server.
						{/if}
					{/await}
				</p>

				<div class="mt-1 flex items-center gap-2">
					{#if musicManager.currentItem}
						<Button
							size="lg"
							class="bg-white text-black hover:bg-white/90"
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

					<Button
						size="lg"
						variant="secondary"
						class="bg-black/25 text-white border-white/30 backdrop-blur-sm hover:bg-black/40"
						href="/browse"
					>
						<Compass />
						Browse
					</Button>
				</div>
			</div>
		</section>

		<!-- Quick actions -->
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
			<section class="flex flex-col gap-8">
				{#each Array(4) as _}
					<div
						class="h-6 w-40 animate-pulse rounded bg-muted"
					></div>
				{/each}
			</section>
		{:then stats}
			{#if stats}
				<section>
					<div class="grid grid-cols-2 gap-2 sm:grid-cols-4">
						<div
							class="flex flex-col gap-1 rounded-lg border bg-card p-4"
						>
							<div class="flex items-center gap-2 text-muted-foreground">
								<Disc3 class="h-4 w-4" />
								<span class="text-xs uppercase tracking-wider">
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

						<div
							class="flex flex-col gap-1 rounded-lg border bg-card p-4"
						>
							<div class="flex items-center gap-2 text-muted-foreground">
								<Clock class="h-4 w-4" />
								<span class="text-xs uppercase tracking-wider">
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
								<span class="text-xs uppercase tracking-wider">
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
								<span class="text-xs uppercase tracking-wider">
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

		<!-- Recently played -->
		<section>
			<a
				class="flex items-center gap-1 text-xl font-semibold hover:cursor-pointer hover:underline"
				href="/users/{data.user.id}/history"
			>
				Recently Played
				<ChevronRight />
			</a>

			<div class="h-4"></div>

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
							<TrackTile
								id={track.id}
								cover={track.coverArt.medium}
								name={track.name}
								artists={track.artists}
							/>
						{/each}
					</div>
				{:else}
					<p class="text-sm text-muted-foreground">
						Nothing played yet — start listening!
					</p>
				{/if}
			{/await}
		</section>

		<!-- Top tracks -->
		<section>
			<div class="flex items-center gap-1 text-xl font-semibold">
				Your Top Tracks
			</div>

			<div class="h-4"></div>

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
							<TrackTile
								id={track.id}
								cover={track.coverArt.medium}
								name={track.name}
								artists={track.artists}
							/>
						{/each}
					</div>
				{:else}
					<p class="text-sm text-muted-foreground">
						Your most-played tracks will appear here.
					</p>
				{/if}
			{/await}
		</section>

		<!-- Your playlists -->
		<section>
			<a
				class="flex items-center gap-1 text-xl font-semibold hover:cursor-pointer hover:underline"
				href="/library/playlists"
			>
				Your Playlists
				<ChevronRight />
			</a>

			<div class="h-4"></div>

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
							<PlaylistTile
								id={playlist.id}
								cover={playlist.coverArt.medium}
								name={playlist.name}
								trackCount={playlist.trackCount}
							/>
						{/each}
					</div>
				{:else}
					<p class="text-sm text-muted-foreground">No playlists yet</p>
				{/if}
			{/await}
		</section>

		<!-- Favorites -->
		<section>
			<a
				class="flex items-center gap-1 text-xl font-semibold hover:cursor-pointer hover:underline"
				href="/library/favorites"
			>
				Favorites
				<ChevronRight />
			</a>

			<div class="h-4"></div>

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
							<TrackTile
								id={track.id}
								cover={track.coverArt.medium}
								name={track.name}
								artists={track.artists}
							/>
						{/each}
					</div>
				{:else}
					<p class="text-sm text-muted-foreground">
						Heart some tracks to see them here
					</p>
				{/if}
			{/await}
		</section>
	</div>
{/if}