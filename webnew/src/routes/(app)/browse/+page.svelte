<script lang="ts">
	import { ChevronRight, DiscAlbum, FileMusic, Search, Users } from "@lucide/svelte";
	import AlbumTile from "$lib/components/tiles/AlbumTile.svelte";
	import { Button } from "$lib/components/ui";

	let { data } = $props();

	const categories = [
		{
			title: "Albums",
			description: "Browse your album collection",
			icon: DiscAlbum,
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
	<section>
		<h1 class="text-xl font-bold">Browse</h1>

		<div class="h-4"></div>

		<div class="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
			{#each categories as category (category.title)}
				<div
					class="flex flex-col gap-3 rounded-lg border bg-card p-4"
				>
					<div
						class="flex items-center gap-3 rounded border bg-muted p-3"
					>
						<category.icon class="h-6 w-6 shrink-0" />
						<div>
							<p class="font-semibold">{category.title}</p>
							<p class="text-xs text-muted-foreground">
								{category.description}
							</p>
						</div>
					</div>

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
</div>