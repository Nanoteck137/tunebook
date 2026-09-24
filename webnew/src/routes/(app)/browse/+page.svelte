<script lang="ts">
	import {
		ChevronRight,
		Compass,
		Disc,
		Disc3,
		ListMusic,
		Mic,
		Music,
		Search,
	} from "@lucide/svelte";
	import AlbumTile from "$lib/components/tiles/AlbumTile.svelte";
	import TrackTile from "$lib/components/tiles/TrackTile.svelte";
	import PlaylistTile from "$lib/components/tiles/PlaylistTile.svelte";
	import SectionHeader from "$lib/components/SectionHeader.svelte";
	import HeroIcon from "$lib/components/HeroIcon.svelte";
	import { Button } from "$lib/components/ui";
	import Spacer from "$lib/components/Spacer.svelte";

	let { data } = $props();

	const categories = [
		{
			title: "Albums",
			description: "Browse your album collection",
			icon: Disc,
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
			icon: Mic,
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
			icon: Music,
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

<div class="section-browse flex flex-col gap-8">
	<section
		class="rounded-lg border bg-linear-to-b from-section-hero-from to-section-hero-to p-4 shadow-sm sm:p-6"
	>
		<div class="flex items-center gap-4">
			<HeroIcon>
				<Compass />
			</HeroIcon>
			<div class="flex min-w-0 flex-col">
				<h1 class="text-2xl font-bold">Browse</h1>
				<p class="text-sm text-muted-foreground">
					Albums, artists and tracks from your library
				</p>
			</div>
		</div>
	</section>

	<section>
		<SectionHeader>
			<Disc3 />
			Categories
		</SectionHeader>

		<Spacer />

		<div class="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
			{#each categories as category (category.title)}
				<div class="flex flex-col gap-3 rounded-lg border bg-card p-4">
					<a
						href={category.allHref}
						class="group flex items-center gap-3 rounded border bg-muted p-3 transition-colors hover:bg-accent"
						title={`Browse all ${category.title.toLowerCase()}`}
					>
						<category.icon class="h-6 w-6 shrink-0" />
						<div class="min-w-0 flex-1">
							<p class="font-semibold group-hover:underline">
								{category.title}
							</p>
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
		<SectionHeader viewAllHref="/albums?sort=created-new">
			<Disc />
			Recently Added Albums
		</SectionHeader>

		<Spacer />

		{#if data.recentAlbums.length > 0}
			<div class="flex gap-2 overflow-x-auto pb-4">
				{#each data.recentAlbums as album (album.id)}
					<AlbumTile class="w-40" size="sm" {album} />
				{/each}
			</div>
		{:else}
			<p class="px-2 text-sm text-muted-foreground">No albums yet</p>
		{/if}
	</section>

	<section>
		<SectionHeader viewAllHref="/tracks?sort=created-new">
			<Music />
			Recently Added Tracks
		</SectionHeader>

		<Spacer />

		{#if data.recentTracks.length > 0}
			<div class="flex gap-2 overflow-x-auto pb-4">
				{#each data.recentTracks as track (track.id)}
					<TrackTile size="sm" class="w-40" {track} />
				{/each}
			</div>
		{:else}
			<p class="px-2 text-sm text-muted-foreground">No tracks yet</p>
		{/if}
	</section>

	<section>
		<SectionHeader viewAllHref="/library/playlists?sort=created-new">
			<ListMusic />
			Recently Added Playlists
		</SectionHeader>

		<Spacer />

		{#if data.recentPlaylists.length > 0}
			<div class="flex gap-2 overflow-x-auto pb-4">
				{#each data.recentPlaylists as playlist (playlist.id)}
					<PlaylistTile size="sm" class="w-40" {playlist} />
				{/each}
			</div>
		{:else}
			<p class="px-2 text-sm text-muted-foreground">No playlists yet</p>
		{/if}
	</section>
</div>
