<script lang="ts">
	import { ChevronRight, DiscAlbum, FileMusic, Search, Users } from "@lucide/svelte";
	import collage from "$lib/assets/collage.png";
	import AlbumTile from "$lib/components/tiles/AlbumTile.svelte";
	import TrackTile from "$lib/components/tiles/TrackTile.svelte";
	import PlaylistTile from "$lib/components/tiles/PlaylistTile.svelte";
	import { Button } from "$lib/components/ui";

	let { data } = $props();

	const categories = [
		{
			title: "Albums",
			description: "Browse your album collection",
			icon: DiscAlbum,
			allHref: "/albums",
			search: { label: "Search Albums", href: "/search/albums" },
			links: [
				{ label: "All albums", href: "/albums" },
				{ label: "Recently added", href: "/albums?sort=created-new" },
				{ label: "Recently updated", href: "/albums?sort=updated-new" },
			],
		},
		{
			title: "Artists",
			description: "Browse your artists",
			icon: Users,
			allHref: "/artists",
			search: { label: "Search Artists", href: "/search/artists" },
			links: [
				{ label: "All artists", href: "/artists" },
				{ label: "Recently added", href: "/artists?sort=created-new" },
				{ label: "Recently updated", href: "/artists?sort=updated-new" },
			],
		},
		{
			title: "Tracks",
			description: "Browse your tracks",
			icon: FileMusic,
			allHref: "/tracks",
			search: { label: "Search Tracks", href: "/search/tracks" },
			links: [
				{ label: "All tracks", href: "/tracks" },
				{ label: "Recently added", href: "/tracks?sort=created-new" },
				{ label: "Recently updated", href: "/tracks?sort=updated-new" },
			],
		},
	];
</script>

<div class="flex flex-col gap-8">
	<section
		class="relative flex justify-center overflow-hidden rounded-lg bg-linear-to-tr from-logo-1 via-logo-2 to-logo-3 p-6 sm:justify-start sm:p-10"
	>
		<div class="absolute -top-16 -right-16 h-48 w-48 rounded-full bg-white/10 blur-2xl"></div>

		<img
			src={collage}
			alt=""
			class="pointer-events-none absolute inset-0 h-full w-full object-cover"
		/>

		<div
			class="relative flex min-w-0 flex-col items-center gap-2 rounded-lg border border-white/20 bg-black/30 p-4 text-center backdrop-blur-sm sm:items-start sm:text-left"
		>
			<h1 class="text-3xl font-bold text-white">Browse</h1>
			<p class="text-sm text-white/90">
				Albums, artists and tracks from your library
			</p>
		</div>
	</section>

	<section>
		<div class="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
			{#each categories as category (category.title)}
				<div
					class="flex flex-col gap-3 rounded-lg border bg-card p-4"
				>
<a
					href={category.allHref}
					class="group flex items-center gap-3 rounded border bg-muted p-3 transition-colors hover:bg-accent"
					title={`Browse all ${category.title.toLowerCase()}`}
				>
					<category.icon class="h-6 w-6 shrink-0" />
					<div class="min-w-0 flex-1">
						<p class="font-semibold group-hover:underline">{category.title}</p>
						<p
							class="text-xs text-muted-foreground group-hover:text-accent-foreground"
						>
							{category.description}
						</p>
					</div>
					<ChevronRight
						class="h-4 w-4 shrink-0 text-muted-foreground transition-transform group-hover:translate-x-0.5"
					/>
				</a>

					<div class="flex flex-col">
						{#each category.links as link (link.href)}
							<a
								class="flex items-center justify-between rounded px-2 py-1.5 text-sm text-foreground hover:bg-accent hover:text-accent-foreground"
								href={link.href}
							>
								{link.label}
								<ChevronRight class="h-4 w-4 text-muted-foreground" />
							</a>
						{/each}
					</div>

					<Button
						variant="outline"
						size="sm"
						class="w-full"
						href={category.search.href}
					>
						<Search />
						{category.search.label}
					</Button>
				</div>
			{/each}
		</div>
	</section>

	<section>
		<a
			class="flex items-center gap-1 text-xl font-semibold hover:cursor-pointer hover:underline"
			href="/albums?sort=created-new"
		>
			Recently Added Albums
			<ChevronRight />
		</a>

		<div class="h-4"></div>

		{#if data.recentAlbums.length > 0}
			<div class="flex gap-2 overflow-x-auto pb-4">
				{#each data.recentAlbums as album (album.id)}
					<AlbumTile
						id={album.id}
						cover={album.coverArt.medium}
						name={album.name}
						artists={album.artists}
					/>
				{/each}
			</div>
		{:else}
			<p class="text-sm text-muted-foreground">No albums yet</p>
		{/if}
	</section>

	<section>
		<a
			class="flex items-center gap-1 text-xl font-semibold hover:cursor-pointer hover:underline"
			href="/tracks"
		>
			Recently Added Tracks
			<ChevronRight />
		</a>

		<div class="h-4"></div>

		{#if data.recentTracks.length > 0}
			<div class="flex gap-2 overflow-x-auto pb-4">
				{#each data.recentTracks as track (track.id)}
					<TrackTile
						id={track.id}
						cover={track.coverArt.medium}
						name={track.name}
						artists={track.artists}
					/>
				{/each}
			</div>
		{:else}
			<p class="text-sm text-muted-foreground">No tracks yet</p>
		{/if}
	</section>

	<section>
		<a
			class="flex items-center gap-1 text-xl font-semibold hover:cursor-pointer hover:underline"
			href="/library/playlists?sort=created-new"
		>
			Recently Added Playlists
			<ChevronRight />
		</a>

		<div class="h-4"></div>

		{#if data.recentPlaylists.length > 0}
			<div class="flex gap-2 overflow-x-auto pb-4">
				{#each data.recentPlaylists as playlist (playlist.id)}
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
	</section>
</div>