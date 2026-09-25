<script lang="ts">
	import { goto } from "$app/navigation";
	import { getApiClient, handleApiError } from "$lib";
	import type { Album } from "$lib/api/types";
	import {
		Button,
		buttonVariants,
		Dialog,
		DropdownMenu,
	} from "$lib/components/ui";
	import { EllipsisVertical, Info } from "@lucide/svelte";
	import InfiniteScroll from "$lib/components/InfiniteScroll.svelte";
	import TileGrid from "$lib/components/tiles/TileGrid.svelte";
	import { InfiniteScrollController } from "$lib/infinite-scroll.svelte";
	import { cn } from "$lib/utils";
	import { onMount } from "svelte";
	import SearchBarHeader from "../SearchBarHeader.svelte";

	let { data } = $props();

	const apiClient = getApiClient();

	const scroll = new InfiniteScrollController<Album>({
		initialLoad: () => {
			if (!data.page) return { items: [], hasMore: false, page: 0 };

			return {
				items: data.albums,
				hasMore: data.page.page + 1 < data.page.totalPages,
				page: data.page.page,
			};
		},
		load: async (page) => {
			const res = await apiClient.searchAlbums({
				query: {
					query: data.query,
					page: String(page),
					perPage: String(data.page?.perPage ?? 30),
				},
			});

			if (!res.success) {
				handleApiError(res.error);
				return null;
			}

			return {
				items: res.data.albums,
				hasMore: res.data.page.page + 1 < res.data.page.totalPages,
			};
		},
		itemKey: (album) => album.id,
	});

	async function doSearch(query: string) {
		await goto(`/search/albums?query=${query}`, {
			invalidateAll: true,
			keepFocus: true,
			replaceState: true,
		});
	}

	function clearSearch() {
		value = "";
		doSearch("");
	}

	let value = $state("");

	onMount(() => {
		value = data.query;
	});

	let infoAlbumId = $state<string | null>(null);
	let infoOpen = $state(false);

	let infoAlbum = $derived(
		infoAlbumId
			? (scroll.items.find((a) => a.id === infoAlbumId) ?? null)
			: null,
	);

	function showInfo(id: string) {
		infoAlbumId = id;
		infoOpen = true;
	}

	function formatDate(iso: string) {
		return new Date(iso).toLocaleDateString(undefined, {
			year: "numeric",
			month: "short",
			day: "numeric",
		});
	}
</script>

<svelte:head>
	<title>Search Albums - Tunebook</title>
</svelte:head>

<div class="flex flex-col gap-6">
	<SearchBarHeader
		searchBarPlaceholder="Search albums..."
		{value}
		setValue={(v) => {
			value = v;
		}}
		search={doSearch}
		searchWithValue={() => {
			doSearch(value);
		}}
		{clearSearch}
	/>

	{#if data.query && scroll.items.length === 0}
		<p class="py-12 text-center text-sm text-muted-foreground">
			No albums found for "{data.query}".
		</p>
	{/if}

	{#if scroll.items.length > 0}
		{#if data.page}
			<div class="flex items-baseline gap-2">
				<span class="text-sm text-muted-foreground">
					{data.page.totalItems} album(s)
				</span>
			</div>
		{/if}

		<InfiniteScroll controller={scroll}>
			<TileGrid>
				{#each scroll.items as album (album.id)}
					<div
						class="group relative flex flex-col overflow-hidden rounded-lg border bg-card transition-shadow hover:shadow-md"
					>
						<a href="/albums/{album.id}">
							<img
								src={album.coverArt.medium}
								alt={album.name}
								class="aspect-square w-full object-cover"
							/>
						</a>

						<div class="flex flex-col gap-0.5 p-2">
							<a
								href="/albums/{album.id}"
								class="truncate text-sm font-medium hover:underline"
								title={album.name}
							>
								{album.name}
							</a>
							<p
								class="truncate text-xs text-muted-foreground"
								title={album.artists.map((a) => a.name).join(", ")}
							>
								{album.artists.map((a) => a.name).join(", ")}
							</p>
						</div>

						<div class="absolute top-1.5 right-1.5">
							<DropdownMenu.Root>
								<DropdownMenu.Trigger
									class={cn(
										buttonVariants({ variant: "secondary", size: "icon" }),
										"rounded-full",
									)}
								>
									<EllipsisVertical size={14} />
								</DropdownMenu.Trigger>
								<DropdownMenu.Content align="end">
									<DropdownMenu.Group>
										<DropdownMenu.Item onclick={() => showInfo(album.id)}>
											<Info />
											Show more info
										</DropdownMenu.Item>
									</DropdownMenu.Group>
								</DropdownMenu.Content>
							</DropdownMenu.Root>
						</div>
					</div>
				{/each}
			</TileGrid>
		</InfiniteScroll>
	{/if}
</div>

<Dialog.Root open={infoOpen} onOpenChange={(v) => (infoOpen = v)}>
	<Dialog.Content class="">
		<Dialog.Header>
			<Dialog.Title>Album Info</Dialog.Title>
			<Dialog.Description>
				Detailed information about the album
			</Dialog.Description>
		</Dialog.Header>

		{#if infoAlbum}
			<div class="flex flex-col gap-4 sm:flex-row">
				<div class="flex shrink-0 justify-center sm:block">
					<img
						src={infoAlbum.coverArt.large}
						alt={infoAlbum.name}
						class="h-48 w-48 rounded-lg border object-cover sm:h-44 sm:w-44"
					/>
				</div>

				<div class="flex min-w-0 flex-1 flex-col gap-2">
					<div>
						<p class="text-lg leading-tight font-semibold">{infoAlbum.name}</p>
						<p class="text-sm text-muted-foreground">
							{#each infoAlbum.artists as artist, i (artist.id)}
								{#if i > 0}
									{", "}
								{/if}
								<a
									href="/artists/{artist.id}"
									class="hover:underline"
									title={artist.name}
								>
									{artist.name}
								</a>
							{/each}
						</p>
					</div>

					<div class="grid grid-cols-[auto_1fr] gap-x-3 gap-y-1 text-sm">
						{#if infoAlbum.year}
							<span class="text-muted-foreground">Year</span>
							<span>{infoAlbum.year}</span>
						{/if}

						{#if infoAlbum.tags.length > 0}
							<span class="text-muted-foreground">Tags</span>
							<div class="flex flex-wrap gap-1">
								{#each infoAlbum.tags as tag (tag)}
									<span class="rounded-md bg-secondary px-1.5 py-0.5 text-xs"
										>{tag}</span
									>
								{/each}
							</div>
						{/if}

						<span class="text-muted-foreground">Added</span>
						<span>{formatDate(infoAlbum.created)}</span>

						<span class="text-muted-foreground">Updated</span>
						<span>{formatDate(infoAlbum.updated)}</span>
					</div>
				</div>
			</div>
		{:else}
			<p class="text-sm text-muted-foreground">Album not found.</p>
		{/if}

		<Dialog.Footer>
			<Button variant="outline" onclick={() => (infoOpen = false)}>
				Close
			</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
