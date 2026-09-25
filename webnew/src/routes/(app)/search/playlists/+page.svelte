<script lang="ts">
	import { goto } from "$app/navigation";
	import { getApiClient, handleApiError } from "$lib";
	import type { Playlist } from "$lib/api/types";
	import Image from "$lib/components/Image.svelte";
	import InfiniteScroll from "$lib/components/InfiniteScroll.svelte";
	import TileGrid from "$lib/components/tiles/TileGrid.svelte";
	import { InfiniteScrollController } from "$lib/infinite-scroll.svelte";
	import { onMount } from "svelte";
	import SearchBarHeader from "../SearchBarHeader.svelte";

	let { data } = $props();

	const apiClient = getApiClient();

	async function doSearch(query: string) {
		await goto(`/search/playlists?query=${query}`, {
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

	const scroll = new InfiniteScrollController<Playlist>({
		initialLoad: () => {
			if (!data.page) return { items: [], hasMore: false, page: 0 };

			return {
				items: data.playlists,
				hasMore: data.page.page + 1 < data.page.totalPages,
				page: data.page.page,
			};
		},
		load: async (page) => {
			const res = await apiClient.searchPlaylists({
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
				items: res.data.playlists,
				hasMore: res.data.page.page + 1 < res.data.page.totalPages,
			};
		},
		itemKey: (playlist) => playlist.id,
	});
</script>

<svelte:head>
	<title>Search Playlists - Tunebook</title>
</svelte:head>

<div class="flex flex-col gap-6">
	<SearchBarHeader
		searchBarPlaceholder="Search playlists..."
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
			No playlists found for "{data.query}".
		</p>
	{/if}

	{#if scroll.items.length > 0}
		{#if data.page}
			<div class="flex items-baseline gap-2">
				<span class="text-sm text-muted-foreground">
					{data.page.totalItems} playlist(s)
				</span>
			</div>
		{/if}

		<InfiniteScroll controller={scroll}>
			<TileGrid>
				{#each scroll.items as playlist (playlist.id)}
					<div
						class="group relative flex flex-col overflow-hidden rounded-lg border bg-card transition-shadow hover:shadow-md"
					>
						<a href="/playlists/{playlist.id}">
							<Image
								class="aspect-square w-full rounded-none border-0"
								src={playlist.coverArt.medium}
								alt={playlist.name}
							/>
						</a>

						<div class="flex flex-col gap-0.5 p-2">
							<a
								href="/playlists/{playlist.id}"
								class="truncate text-sm font-medium hover:underline"
								title={playlist.name}
							>
								{playlist.name}
							</a>
							<p class="truncate text-xs text-muted-foreground">
								{playlist.trackCount} track{playlist.trackCount !== 1
									? "s"
									: ""}
							</p>
							<p
								class="flex items-center gap-1 truncate text-xs text-muted-foreground"
							>
								{#if playlist.ownerPicture}
									<img
										src={playlist.ownerPicture.small}
										alt=""
										class="h-4 w-4 rounded-full object-cover"
									/>
								{/if}
								{playlist.ownerDisplayName}
							</p>
						</div>
					</div>
				{/each}
			</TileGrid>
		</InfiniteScroll>
	{/if}
</div>
